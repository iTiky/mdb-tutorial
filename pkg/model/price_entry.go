package model

import (
	"time"

	"github.com/itiky/mdb-tutorial/pkg/common"
)

// PriceEntry is an output for Product / PricesImport aggregate.
type PriceEntry struct {
	Name      string    `json:"name" bson:"name"`
	Price     int       `json:"price" bson:"price"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// PriceEntries is a slice of PriceEntry objects.
type PriceEntries []PriceEntry

// PriceEntriesSortOption defines PriceEntry sorter interface.
type PriceEntriesSortOption func(s common.SortOptions) common.SortOptions

// PriceEntrySortByProductName returns PriceEntry sort by product name option.
func PriceEntrySortByProductName(order common.Order) PriceEntriesSortOption {
	return func(s common.SortOptions) common.SortOptions {
		s = append(s, common.SortOption{
			FieldName: "name",
			Order:     order,
		})
		return s
	}
}

// PriceEntrySortByPrice returns PriceEntry sort by price name option.
func PriceEntrySortByPrice(order common.Order) PriceEntriesSortOption {
	return func(s common.SortOptions) common.SortOptions {
		s = append(s, common.SortOption{
			FieldName: "price",
			Order:     order,
		})
		return s
	}
}

// PriceEntrySortByImportTimestamp returns PriceEntry sort by import timestamp option.
func PriceEntrySortByImportTimestamp(order common.Order) PriceEntriesSortOption {
	return func(s common.SortOptions) common.SortOptions {
		s = append(s, common.SortOption{
			FieldName: "timestamp",
			Order:     order,
		})
		return s
	}
}

// NewPriceEntriesSortOptions creates a new SortOptions for PriceEntry objects.
func NewPriceEntriesSortOptions(options ...PriceEntriesSortOption) common.SortOptions {
	s := common.SortOptions{}
	for _, option := range options {
		s = option(s)
	}

	return s
}
