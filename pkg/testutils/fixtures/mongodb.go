package fixtures

import (
	"github.com/itiky/mdb-tutorial/pkg/model"
)

// MongoDBCollection is a MongoDB fixtures collection interface.
type MongoDBCollection interface {
	GetCollection() string
	GetBSONObjects() []interface{}
}

// MongoDB keeps MongoDB fixtures.
type MongoDB struct {
	Collections []MongoDBCollection
}

// NewDefaultMongoDBFixtures returns prefilled fixtures.
func NewDefaultMongoDBFixtures() MongoDB {
	return MongoDB{
		Collections: []MongoDBCollection{
			MongoDBProduct{
				Products: Products,
			},
			MongoDBPriceImport{
				Imports: PriceImports,
			},
		},
	}
}

// NewEmptyMongoDBFixtures returns empty fixtures.
func NewEmptyMongoDBFixtures() MongoDB {
	return MongoDB{
		Collections: []MongoDBCollection{
			MongoDBProduct{
				Products: []model.Product{},
			},
			MongoDBPriceImport{
				Imports: []model.PricesImport{},
			},
		},
	}
}
