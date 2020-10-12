package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/model"
)

var _ PriceImportStorage = (*priceImportStorage)(nil)

// priceImportStorage keeps PriceImportStorage dependencies.
type priceImportStorage struct {
	storageCommon
	mdbCollection      *mongo.Collection
	productsCollection string
}

// UpsertByProductIDAndTimestamp implements PriceImportStorage interface.
func (s priceImportStorage) UpsertByProductIDAndTimestamp(ctx context.Context, pricesImport model.PricesImport) (createdID primitive.ObjectID, retErr error) {
	if pricesImport.Timestamp.IsZero() {
		retErr = fmt.Errorf("%w: timestamp: can not be empty", common.ErrInvalidInput)
		return
	}
	if pricesImport.ProductID.IsZero() {
		retErr = fmt.Errorf("%w: product_id: can not be empty", common.ErrInvalidInput)
		return
	}

	filter := bson.M{
		"timestamp":  pricesImport.Timestamp,
		"product_id": pricesImport.ProductID,
	}
	update := bson.M{"$set": bson.M{
		"prices": pricesImport.Prices,
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

// GetAll implements PriceImportStorage interface.
func (s priceImportStorage) GetAll(ctx context.Context, timestamp time.Time, productID string) (retObjs []model.PricesImport, retErr error) {
	filter := bson.M{}

	if !timestamp.IsZero() {
		filter["timestamp"] = timestamp
	}
	if productID != "" {
		id, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			retErr = fmt.Errorf("productID filter: %w", err)
			return
		}
		filter["product_id"] = id
	}

	cursor, err := s.mdbCollection.Find(ctx, filter)
	if err != nil {
		retErr = err
		return
	}

	err = cursorIterateAndDecode(ctx, cursor, func(curCursor *mongo.Cursor) error {
		var pricesImport model.PricesImport
		if err := curCursor.Decode(&pricesImport); err != nil {
			return err
		}
		retObjs = append(retObjs, pricesImport)

		return nil
	})
	if err != nil {
		retErr = err
		return
	}

	return
}

// GetPriceEntries implements PriceImportStorage interface.
// nolint:govet
func (s priceImportStorage) GetPriceEntries(
	ctx context.Context,
	sortOptions common.SortOptions, paginationOption common.PaginationOption,
) (retObjs model.PriceEntries, retErr error) {

	// define required stages
	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", s.productsCollection},
			{"localField", "product_id"},
			{"foreignField", "_id"},
			{"as", "fromProducts"},
		}},
	}
	replaceRootStage := bson.D{
		{"$replaceRoot", bson.D{
			{"newRoot", bson.D{
				{"$mergeObjects", bson.A{
					bson.D{{"$arrayElemAt", bson.A{"$fromProducts", 0}}},
					"$$ROOT",
				}},
			}},
		}},
	}
	unwindStage := bson.D{
		{"$unwind", "$prices"},
	}
	addFieldsStage := bson.D{
		{"$addFields", bson.D{
			{"price", "$prices.value"},
		}},
	}
	projectStage := bson.D{
		{"$project", bson.D{
			{"fromProducts", 0},
			{"_id", 0},
			{"product_id", 0},
			{"prices", 0},
		}},
	}

	// build pipeline
	pipeline := mongo.Pipeline{lookupStage, replaceRootStage, unwindStage, addFieldsStage, projectStage}
	pipeline = addSortAggregationStage(pipeline, sortOptions)
	pipeline = addPaginationAggregationStage(pipeline, paginationOption)

	cursor, err := s.mdbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		retErr = err
		return
	}

	err = cursorIterateAndDecode(ctx, cursor, func(curCursor *mongo.Cursor) error {
		var priceEntry model.PriceEntry
		if err := curCursor.Decode(&priceEntry); err != nil {
			return err
		}
		retObjs = append(retObjs, priceEntry)

		return nil
	})
	if err != nil {
		retErr = err
		return
	}

	return
}
