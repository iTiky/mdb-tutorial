package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product keeps product data.
type Product struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	// Product unique name
	Name string `json:"name" bson:"name"`
}
