package storage

import (
	"context"
	"errors"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/model"
	"github.com/itiky/mdb-tutorial/pkg/testutils"
	"github.com/itiky/mdb-tutorial/pkg/testutils/fixtures"
)

func (s *StorageTestSuite) TestStorage_Product() {
	t := s.T()
	ctx := context.Background()

	client := testutils.PrepareMongoDBFixtures(ctx, s.r, fixtures.NewEmptyMongoDBFixtures())
	storage, err := NewStorage(
		WithDatabase(testutils.TestMongoDBDatabase),
		WithMongoDBClient(client),
	)
	require.NoError(t, err)
	targetSt := storage.Product()

	// check GetAll: empty
	{
		resp, err := targetSt.GetAll(ctx)
		require.NoError(t, err)
		require.Empty(t, resp)
	}

	// check GetByID: invalid ObjectID
	{
		_, err := targetSt.GetByID(ctx, "invalid")
		require.Error(t, err)
	}

	// check GetByID: non-existing ObjectID
	{
		_, err := targetSt.GetByID(ctx, primitive.NewObjectIDFromTimestamp(time.Now()).Hex())
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrNotFound))
	}

	// check GetByName: non-existing name
	{
		_, err := targetSt.GetByName(ctx, "non-existing")
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrNotFound))
	}

	// check: UpsertByName (new)
	product := model.Product{Name: "NewProduct"}
	{
		id, err := targetSt.UpsertByName(ctx, product)
		require.NoError(t, err)
		require.False(t, id.IsZero())
		product.ID = id
	}

	// check GetByID: existing
	{
		rcvProduct, err := targetSt.GetByID(ctx, product.ID.Hex())
		require.NoError(t, err)
		require.EqualValues(t, product, rcvProduct)
	}

	// check GetByName: existing
	{
		rcvProduct, err := targetSt.GetByName(ctx, product.Name)
		require.NoError(t, err)
		require.EqualValues(t, product, rcvProduct)
	}

	// check GetAll: existing
	{
		companies, err := targetSt.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, companies, 1)
		require.EqualValues(t, product, companies[0])
	}

	// check: UpsertByName (existing)
	{
		id, err := targetSt.UpsertByName(ctx, product)
		require.NoError(t, err)
		require.True(t, id.IsZero())

		// check only one exists
		companies, err := targetSt.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, companies, 1)
	}

	// create one more
	{
		id, err := targetSt.UpsertByName(ctx, model.Product{Name: "OtherProduct"})
		require.NoError(t, err)
		require.False(t, id.IsZero())

		// check only one exists
		companies, err := targetSt.GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, companies, 2)
	}
}
