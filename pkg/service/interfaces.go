package service

import (
	"context"
	"io"
	"time"

	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/model"
)

type Service interface {
	// CSVImporter returns configured CSVImporter service.
	CSVImporter() CSVImporterService
	// CSVProcessor returns configured CSVProcessor service.
	CSVProcessor() CSVProcessorService
	// PriceEntries returns configured PriceEntries service.
	PriceEntries() PriceEntriesService
}

// CSVImporterService processes product-price data CSV-file import.
type CSVImporterService interface {
	// ImportPrices imports CSV file data containing price changes per product.
	// If import already exists, its data would be overwritten.
	ImportPrices(ctx context.Context, csvImport model.CSVImport) error
}

// CSVProcessorService downloads and parses product-price data CSV-file.
type CSVProcessorService interface {
	// Download download a CSV-file to temp dir and returns filePath and download timestamp.
	Download(inputPath string) (string, time.Time, error)
	// Process processed downloaded CSV-file sequentially in chunks.
	Process(ctx context.Context, reader io.Reader, importTimestamp time.Time, chunkSize int, chunkWorker csvChunkWorker) error
}

// PriceEntriesService provides product-price entries operations.
type PriceEntriesService interface {
	// List queries price entries with pagination and sorting options.
	List(ctx context.Context, paginationOpt common.PaginationOption, sortOpts common.SortOptions) (model.PriceEntries, error)
}
