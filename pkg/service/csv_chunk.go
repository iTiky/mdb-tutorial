package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/itiky/mdb-tutorial/pkg/model"
)

const (
	csvChunkMaxParsingErrors = 10
)

// csvChunkWorker is a CSV-file chunk writer/executor interface.
type csvChunkWorker func(ctx context.Context, csvImport model.CSVImport) error

// csvChunk keeps CSV file chunk data.
type csvChunk struct {
	id             int
	timestamp      time.Time
	entries        model.CSVEntries
	parsingErrors  []error
	executionError error
}

// addEntry appends a new CSV entry to chunk.
func (c *csvChunk) addEntry(productName string, price int) {
	c.entries = append(c.entries, model.CSVEntry{
		ProductName: productName,
		Price:       price,
	})
}

// addParsingError adds a new parsing error to chunk.
func (c *csvChunk) addParsingError(line int, err error) {
	if len(c.parsingErrors) == cap(c.parsingErrors) {
		return
	}

	c.parsingErrors = append(c.parsingErrors, fmt.Errorf("parsing line [%d]: %v", line, err))
}

// isEmpty checks if chunk is empty.
func (c *csvChunk) isEmpty() bool {
	return len(c.entries) == 0
}

// isFull checks if chunk is full and ready to be processed.
func (c *csvChunk) isFull() bool {
	return len(c.entries) == cap(c.entries)
}

// isFailed checks if chunk has any errors.
func (c *csvChunk) isFailed() bool {
	return len(c.parsingErrors) != 0 || c.executionError != nil
}

// execute prepares a model.CSVImport object and runs the csvChunkWorker.
func (c *csvChunk) execute(ctx context.Context, worker csvChunkWorker) {
	csvImport := model.CSVImport{
		Timestamp: c.timestamp,
		Entries:   c.entries,
	}

	c.executionError = worker(ctx, csvImport)
}

// getError build an accumulated error string based on parsing and execution errors.
func (c *csvChunk) getErrorString() string {
	str := strings.Builder{}

	str.WriteString(fmt.Sprintf("chunkID: %d; ", c.id))

	str.WriteString("parsing: ")
	for i, err := range c.parsingErrors {
		str.WriteString(err.Error())
		if i < len(c.parsingErrors)-1 {
			str.WriteString(", ")
		}
	}
	if len(c.parsingErrors) == cap(c.parsingErrors) {
		str.WriteString(",..")
	}
	str.WriteString("; ")

	str.WriteString(fmt.Sprintf("execution: %v", c.executionError))

	return str.String()
}

// getStateString returns current chunk state.
func (c *csvChunk) getStateString() string {
	return fmt.Sprintf("chunkID %d: %d", c.id, len(c.entries))
}

// newCSVChunk creates a new csvChunk objects with limited number of parsing errors and entries.
func newCSVChunk(id, size int, timestamp time.Time) *csvChunk {
	return &csvChunk{
		id:             id,
		timestamp:      timestamp,
		entries:        make(model.CSVEntries, 0, size),
		parsingErrors:  make([]error, 0, csvChunkMaxParsingErrors),
		executionError: nil,
	}
}
