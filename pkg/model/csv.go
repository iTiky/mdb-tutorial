package model

import "time"

// CSVImport contains price data to be imported.
type CSVImport struct {
	Timestamp time.Time
	Entries   CSVEntries
}

// CSVEntry contains one CSV import file row data.
type CSVEntry struct {
	ProductName string
	Price       int
}

// CSVEntries is a slice of CSVEntry objects.
type CSVEntries []CSVEntry
