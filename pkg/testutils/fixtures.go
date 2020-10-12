package testutils

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/itiky/mdb-tutorial/pkg/testutils/fixtures"
)

// PrepareMongoDBFixtures loads MongoDB fixtures and returns connection client.
func PrepareMongoDBFixtures(ctx context.Context, r *Resources, fixtures fixtures.MongoDB) *mongo.Client {
	cont, client := r.MongoDB(ctx)

	db := client.Database(cont.Db)
	for _, collection := range fixtures.Collections {
		// drop collection
		if err := db.Collection(collection.GetCollection()).Drop(ctx); err != nil {
			r.crash(ctx, fmt.Errorf("mdb: collection %s: drop: %w", collection.GetCollection(), err))
		}

		// insert collection
		insertObjs := collection.GetBSONObjects()
		if len(insertObjs) > 0 {
			res, err := db.Collection(collection.GetCollection()).InsertMany(ctx, insertObjs)
			if err != nil {
				r.crash(ctx, fmt.Errorf("mdb: collection %s: insert: %w", collection.GetCollection(), err))
			}
			if len(res.InsertedIDs) != len(insertObjs) {
				r.crash(ctx, fmt.Errorf("mdb: collection %s: insert: length mismatch %d / %d", collection.GetCollection(), len(res.InsertedIDs), len(insertObjs)))
			}
		}
	}

	return client
}
