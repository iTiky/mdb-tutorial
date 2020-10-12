package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PricesImport keeps prices import data.
type PricesImport struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	// Reference ID for Product
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	// Prices import DateTime
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	// Prices data
	Prices []Price `json:"prices" bson:"prices"`
}

// Price is an embedded PricesImport struct.
type Price struct {
	Value int `json:"value" bson:"value"`
}
