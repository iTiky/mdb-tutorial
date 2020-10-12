package storage

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DefaultDB              = "db"
	ProductsCollection     = "products"
	PriceImportsCollection = "price_imports"
)

var _ Storage = (*storage)(nil)

// storageCommon keeps common dependencies for all Storage objects.
type storageCommon struct {
	client *mongo.Client
	logger *logrus.Logger
}

// storage implements Storage interface.
type storage struct {
	storageCommon
	db string
}

// Product implements Storage interface.
// nolint:gosimple
func (s storage) Product() ProductStorage {
	return productStorage{
		s.storageCommon,
		s.client.Database(s.db).Collection(ProductsCollection),
	}
}

// PriceImport implements Storage interface.
// nolint:gosimple
func (s storage) PriceImport() PriceImportStorage {
	return priceImportStorage{
		s.storageCommon,
		s.client.Database(s.db).Collection(PriceImportsCollection),
		ProductsCollection,
	}
}

// Option specifies functional argument used by NewStorage function.
type Option func(storage *storage) error

// WithMongoDBClient sets MongoDB client for storage.
func WithMongoDBClient(client *mongo.Client) Option {
	return func(storage *storage) error {
		if client == nil {
			return fmt.Errorf("mongoDB client option: nil")
		}
		storage.client = client

		return nil
	}
}

// WithLogger sets logger for storage.
func WithLogger(logger *logrus.Logger) Option {
	return func(storage *storage) error {
		if logger == nil {
			return fmt.Errorf("logger option: nil")
		}
		storage.logger = logger

		return nil
	}
}

// WithDatabase sets database name for storage.
func WithDatabase(db string) Option {
	return func(storage *storage) error {
		if db == "" {
			return fmt.Errorf("db option: empty")
		}
		storage.db = db

		return nil
	}
}

// NewStorage creates a new configured Storage object.
func NewStorage(options ...Option) (Storage, error) {
	s := &storage{
		db: DefaultDB,
	}
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	if s.client == nil {
		client, err := mongo.NewClient()
		if err != nil {
			return nil, fmt.Errorf("default MongoDB client: %v", err)
		}
		s.client = client
	}
	if s.logger == nil {
		logger := logrus.New()
		logger.SetOutput(ioutil.Discard)
		s.logger = logger
	}

	return s, nil
}
