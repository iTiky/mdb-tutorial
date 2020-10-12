package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/itiky/mdb-tutorial/pkg/model"
)

var _ ProductStorage = (*productStorage)(nil)

// productStorage keeps ProductStorage dependencies.
type productStorage struct {
	storageCommon
	mdbCollection *mongo.Collection
}

// GetByID implements ProductStorage interface.
func (s productStorage) GetByID(ctx context.Context, id string) (retObj model.Product, retErr error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		retErr = err
		return
	}

	filter := bson.M{"_id": objectID}
	res := s.mdbCollection.FindOne(ctx, filter)
	if err := singleResultDecode(res, &retObj); err != nil {
		retErr = err
		return
	}

	return
}

// GetByName implements ProductStorage interface.
func (s productStorage) GetByName(ctx context.Context, name string) (retObj model.Product, retErr error) {
	filter := bson.M{"name": name}
	res := s.mdbCollection.FindOne(ctx, filter)
	if err := singleResultDecode(res, &retObj); err != nil {
		retErr = err
		return
	}

	return
}

// UpsertByName implements ProductStorage interface.
func (s productStorage) UpsertByName(ctx context.Context, product model.Product) (createdID primitive.ObjectID, retErr error) {
	filter := bson.M{"name": product.Name}
	update := bson.M{"$set": bson.M{
		"name": product.Name,
	}}

	res, err := s.mdbCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		retErr = err
		return
	}

	if res.UpsertedCount > 0 {
		id, ok := res.UpsertedID.(primitive.ObjectID)
		if !ok {
			retErr = fmt.Errorf("res.UpsertedID type convertion failed: %T", res.UpsertedID)
		}
		createdID = id
	}

	return
}

// GetAll implements ProductStorage interface.
func (s productStorage) GetAll(ctx context.Context) (retObjs []model.Product, retErr error) {
	filter := bson.D{}
	cursor, err := s.mdbCollection.Find(ctx, filter)
	if err != nil {
		retErr = err
		return
	}

	err = cursorIterateAndDecode(ctx, cursor, func(curCursor *mongo.Cursor) error {
		var product model.Product
		if err := curCursor.Decode(&product); err != nil {
			return err
		}
		retObjs = append(retObjs, product)

		return nil
	})
	if err != nil {
		retErr = err
		return
	}

	return
}
