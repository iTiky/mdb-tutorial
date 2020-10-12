package common

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	AscOrder  Order = true
	DescOrder Order = false
	//
	MaxLimit = 1000
)

// Order is a sort order type.
type Order bool

// MongoDBOrder returns MongoDB sort order option value.
func (o Order) MongoDBOrder() int {
	if o == AscOrder {
		return 1
	}

	return -1
}

// SortOption keeps sort request info.
type SortOption struct {
	FieldName string
	Order     Order
}

// ToBSON converts SortOption to BSON for aggregation usage.
func (o SortOption) ToBSON() bson.E {
	return bson.E{
		Key:   o.FieldName,
		Value: o.Order.MongoDBOrder(),
	}
}

// SortOptions keeps multiple sort request.
// Slice index dictates the sort priority (1st - highest).
type SortOptions []SortOption

// ToBSON converts SortOptions to BSON slice for aggregation usage.
func (o SortOptions) ToBSON() []bson.E {
	e := make([]bson.E, 0, len(o))
	for _, opt := range o {
		e = append(e, opt.ToBSON())
	}

	return e
}

// PaginationOption keeps pagination request info.
type PaginationOption struct {
	Skip  int
	Limit int
}

// Validate validates PaginationOption.
func (o PaginationOption) Validate() error {
	if o.Skip < 0 {
		return fmt.Errorf("skip: should be GTE 0")
	}
	if o.Limit <= 0 {
		return fmt.Errorf("limit: should be GTE 1")
	}
	if o.Limit > MaxLimit {
		return fmt.Errorf("limit: should be LT %d", MaxLimit)
	}

	return nil
}

// NewPaginationOption creates a new PaginationOption object.
func NewPaginationOption(skip, limit int) PaginationOption {
	return PaginationOption{
		Skip:  skip,
		Limit: limit,
	}
}
