package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/model"
)

// Storage provides MongoDB storage operations.
type Storage interface {
	// Product returns configured ProductStorage.
	Product() ProductStorage
	// PriceImport returns configured PriceImportStorage.
	PriceImport() PriceImportStorage
}

// ProductStorage provides "products" collection operation.
type ProductStorage interface {
	// GetByID loads product by ID.
	GetByID(ctx context.Context, id string) (model.Product, error)
	// GetByName loads product by name.
	GetByName(ctx context.Context, name string) (model.Product, error)
	// UpsertByName sets product by unique name.
	// Returns created ID (if created).
	UpsertByName(ctx context.Context, product model.Product) (primitive.ObjectID, error)
	// GetAll loads all product objects.
	GetAll(ctx context.Context) ([]model.Product, error)
}

// PriceImportStorage provides "price_imports" collection operation.
type PriceImportStorage interface {
	// UpsertByProductIDAndTimestamp sets price import by unique pair {timestamp, productID}.
	// Returns created ID (if created).
	UpsertByProductIDAndTimestamp(ctx context.Context, pricesImport model.PricesImport) (primitive.ObjectID, error)
	// GetAll loads all price import objects with optional filtering.
	GetAll(ctx context.Context, timestamp time.Time, productID string) ([]model.PricesImport, error)
	// GetPriceEntries returns merged Product and PriceImport collections with sort and pagination options.
	GetPriceEntries(ctx context.Context, sortOptions common.SortOptions, paginationOption common.PaginationOption) (model.PriceEntries, error)
}
