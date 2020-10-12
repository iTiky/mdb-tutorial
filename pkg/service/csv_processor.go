package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/itiky/mdb-tutorial/pkg/common"
)

var _ CSVProcessorService = (*csvProcessorService)(nil)

type csvProcessorService struct {
	logger *logrus.Logger
}

// Download implements CSVDownloaderService interface.
func (s csvProcessorService) Download(inputPath string) (outputFilePath string, downloadTimestamp time.Time, retErr error) {
	// input check
	if _, err := url.Parse(inputPath); err != nil {
		retErr = fmt.Errorf("%w: path invalid: %v", common.ErrInvalidInput, err)
		return
	}

	// download
	resp, err := http.Get(inputPath)
	if err != nil {
		retErr = fmt.Errorf("GET failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		retErr = fmt.Errorf("%w: GET failed: %s", common.ErrNotFound, resp.Status)
		return
	}

	// create a tmp file to avoid potentially big RAM usage
	downloadTimestamp = time.Now().UTC()
	outputFileName := fmt.Sprintf("prices_import_%s.csv", downloadTimestamp.Format("2006-01-02T15-04-05"))
	outputFilePath = path.Join(os.TempDir(), outputFileName)

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		retErr = fmt.Errorf("creating tmpFile %s: %v", outputFilePath, err)
		return
	}
	defer outputFile.Close()

	// copy response body
	n, err := io.Copy(outputFile, resp.Body)
	if err != nil {
		retErr = fmt.Errorf("body to file copy failed: %v", err)
		return
	}
	if n == 0 {
		retErr = fmt.Errorf("body to file copy failed: written bytes (%d)", n)
		return
	}

	s.logger.Infof("file %s downloaded: %s", inputPath, outputFilePath)

	return
}

// ParseCSV implements CSVDownloaderService interface.
func (s csvProcessorService) Process(
	ctx context.Context,
	reader io.Reader, importTimestamp time.Time,
	chunkSize int, chunkWorker csvChunkWorker,
) error {

	// input check
	if reader == nil {
		return fmt.Errorf("%w: reader is nil", common.ErrInvalidInput)
	}
	if chunkSize <= 0 {
		return fmt.Errorf("%w: chunkSize should be GT 0", common.ErrInvalidInput)
	}
	if chunkWorker == nil {
		return fmt.Errorf("%w: chunkWorker is nil", common.ErrInvalidInput)
	}

	// configure CSV-reader
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	// processChunk executes chunk and accumulates errors
	retErrStrings := make([]string, 0)
	processFailedChunk := func(chunk *csvChunk) {
		errStr := chunk.getErrorString()
		retErrStrings = append(retErrStrings, errStr)
		s.logger.Errorf("processing CSV: %s", errStr)
	}
	processChunk := func(chunk *csvChunk) {
		chunk.execute(ctx, chunkWorker)
		if chunk.isFailed() {
			processFailedChunk(chunk)
		} else {
			s.logger.Infof("processing CSV: %s", chunk.getStateString())
		}
	}

	// read file line by line
	curLineNumber, curChunkID := 0, 1
	curChunk := newCSVChunk(curChunkID, chunkSize, importTimestamp)
	for {
		curLineNumber++

		// read line and process reading error
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			curChunk.addParsingError(curLineNumber, err)
			break
		}

		// parse row
		if len(row) != 2 {
			curChunk.addParsingError(curLineNumber, fmt.Errorf("invalid row length (%d)", len(row)))
			continue
		}
		productName := row[0]
		price, err := strconv.ParseInt(strings.TrimSpace(row[1]), 10, 32)
		if err != nil {
			curChunk.addParsingError(curLineNumber, fmt.Errorf("price convertion failed (%s)", row[1]))
			continue
		}

		// append entry and process the current chunk
		curChunk.addEntry(productName, int(price))
		if curChunk.isFull() {
			processChunk(curChunk)
			curChunk = newCSVChunk(curChunkID, chunkSize, importTimestamp)
		}
	}

	// check the last chunk (not full)
	if !curChunk.isEmpty() {
		processChunk(curChunk)
	} else if curChunk.isFailed() {
		processFailedChunk(curChunk)
	}

	// check chunk errors
	if len(retErrStrings) > 0 {
		return fmt.Errorf("partially processed: %s", strings.Join(retErrStrings, ", "))
	}

	return nil
}
