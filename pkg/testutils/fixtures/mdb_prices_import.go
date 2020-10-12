package fixtures

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/itiky/mdb-tutorial/pkg/model"
)

var PriceImports = []model.PricesImport{
	{
		ID:        primitive.ObjectID([12]byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 2}),
		ProductID: Products[0].ID,
		Timestamp: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Prices: []model.Price{
			{Value: 50},
			{Value: 100},
		},
	},
	{
		ID:        primitive.ObjectID([12]byte{0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 2}),
		ProductID: Products[1].ID,
		Timestamp: time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC),
		Prices: []model.Price{
			{Value: 75},
			{Value: 100},
			{Value: 150},
		},
	},
}

type MongoDBPriceImport struct {
	Imports []model.PricesImport
}

// GetCollection implements MongoDBCollection interface.
func (f MongoDBPriceImport) GetCollection() string {
	return "price_imports"
}

// GetBSONObjects implements MongoDBCollection interface.
func (f MongoDBPriceImport) GetBSONObjects() []interface{} {
	output := make([]interface{}, 0, len(f.Imports))
	for _, priceImport := range f.Imports {
		mPrices := make([]bson.M, 0, len(priceImport.Prices))
		for _, price := range priceImport.Prices {
			mPrices = append(mPrices, bson.M{
				"value": price.Value,
			})
		}

		output = append(output, bson.M{
			"_id":        priceImport.ID,
			"product_id": priceImport.ProductID,
			"timestamp":  priceImport.Timestamp,
			"prices":     mPrices,
		})
	}

	return output
}
