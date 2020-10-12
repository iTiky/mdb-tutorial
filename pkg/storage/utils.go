package storage

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/itiky/mdb-tutorial/pkg/common"
)

// singleResultDecode checks for mongo.SingleResult errors and decodes a result object.
func singleResultDecode(res *mongo.SingleResult, obj interface{}) error {
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return common.ErrNotFound
		}

		return fmt.Errorf("result: %w", res.Err())
	}

	if err := res.Decode(obj); err != nil {
		return fmt.Errorf("decoding: %w", err)
	}

	return nil
}

// cursorIterateAndDecode iterates over mongo.Cursor, decodes result and emits the handler.
func cursorIterateAndDecode(ctx context.Context, cursor *mongo.Cursor, handler func(cursor *mongo.Cursor) error) error {
	idx := 0
	for cursor.Next(ctx) {
		if err := handler(cursor); err != nil {
			return fmt.Errorf("handling [%d]: %w", idx, err)
		}
		idx++
	}

	if cursor.Err() != nil {
		return cursor.Err()
	}

	return nil
}

// addSortAggregationStage adds sort stage to the aggregation pipeline if needed.
// nolint:govet
func addSortAggregationStage(pipeline mongo.Pipeline, opt common.SortOptions) mongo.Pipeline {
	sortBson := opt.ToBSON()
	if len(sortBson) == 0 {
		return pipeline
	}

	return append(pipeline, bson.D{
		{"$sort", sortBson},
	})
}

// addPaginationAggregationStage adds pagination stage to the aggregation pipeline.
// nolint:govet
func addPaginationAggregationStage(pipeline mongo.Pipeline, opt common.PaginationOption) mongo.Pipeline {
	pageStages := []bson.D{
		{{"$skip", opt.Skip}},
		{{"$limit", opt.Limit}},
	}

	return append(pipeline, pageStages...)
}
