package fixtures

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/itiky/mdb-tutorial/pkg/model"
)

var Products = []model.Product{
	{
		ID:   primitive.ObjectID([12]byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1}),
		Name: "Product A",
	},
	{
		ID:   primitive.ObjectID([12]byte{0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 1}),
		Name: "Product B",
	},
}

type MongoDBProduct struct {
	Products []model.Product
}

// GetCollection implements MongoDBCollection interface.
func (f MongoDBProduct) GetCollection() string {
	return "products"
}

// GetBSONObjects implements MongoDBCollection interface.
func (f MongoDBProduct) GetBSONObjects() []interface{} {
	output := make([]interface{}, 0, len(f.Products))
	for _, product := range f.Products {
		output = append(output, bson.M{
			"_id":  product.ID,
			"name": product.Name,
		})
	}

	return output
}

// CheckFindProductByID is a test helper which find fixture product by ID.
func CheckFindProductByID(t *testing.T, id primitive.ObjectID) model.Product {
	for _, product := range Products {
		if product.ID.String() == id.String() {
			return product
		}
	}
	t.Fatalf("product with ID %s: not found in fixtures", id)

	return model.Product{}
}
