package storage

import (
	"context"
	"sort"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/model"
	"github.com/itiky/mdb-tutorial/pkg/testutils"
	"github.com/itiky/mdb-tutorial/pkg/testutils/fixtures"
)

func (s *StorageTestSuite) TestStorage_PriceImportBasic() {
	t := s.T()
	ctx := context.Background()

	client := testutils.PrepareMongoDBFixtures(ctx, s.r, fixtures.NewEmptyMongoDBFixtures())
	storage, err := NewStorage(
		WithDatabase(testutils.TestMongoDBDatabase),
		WithMongoDBClient(client),
	)
	require.NoError(t, err)
	targetSt := storage.PriceImport()

	// check GetAll: empty
	{
		resp, err := targetSt.GetAll(ctx, time.Time{}, "")
		require.NoError(t, err)
		require.Empty(t, resp)
	}

	// check UpsertByProductIDAndTimestamp: invalid input
	{
		priceImport := fixtures.PriceImports[0]

		priceImport1 := priceImport
		priceImport1.ProductID = primitive.ObjectID{}
		_, err := targetSt.UpsertByProductIDAndTimestamp(ctx, priceImport1)
		require.Error(t, err)

		priceImport2 := priceImport
		priceImport2.Timestamp = time.Time{}
		_, err = targetSt.UpsertByProductIDAndTimestamp(ctx, priceImport2)
		require.Error(t, err)
	}

	// check UpsertByProductIDAndTimestamp: 1st one
	priceImport1 := fixtures.PriceImports[0]
	{
		id, err := targetSt.UpsertByProductIDAndTimestamp(ctx, priceImport1)
		require.NoError(t, err)
		require.False(t, id.IsZero())
		priceImport1.ID = id
	}

	// check UpsertByProductIDAndTimestamp: 1st one (no updates)
	{
		id, err := targetSt.UpsertByProductIDAndTimestamp(ctx, priceImport1)
		require.NoError(t, err)
		require.True(t, id.IsZero())
	}

	// check UpsertByProductIDAndTimestamp: 1st one (update)
	priceImport1.Prices = append(priceImport1.Prices, model.Price{Value: 1000})
	{
		id, err := targetSt.UpsertByProductIDAndTimestamp(ctx, priceImport1)
		require.NoError(t, err)
		require.True(t, id.IsZero())
	}

	// add one more import
	priceImport2 := fixtures.PriceImports[1]
	{
		id, err := targetSt.UpsertByProductIDAndTimestamp(ctx, priceImport2)
		require.NoError(t, err)
		require.False(t, id.IsZero())
		priceImport2.ID = id
	}

	// check GetAll: existing
	{
		resp, err := targetSt.GetAll(ctx, time.Time{}, "")
		require.NoError(t, err)
		require.Len(t, resp, 2)
	}

	// check GetAll: filter for the 1st one (also check if it was updated)
	{
		resp, err := targetSt.GetAll(ctx, priceImport1.Timestamp, priceImport1.ProductID.Hex())
		require.NoError(t, err)
		require.Len(t, resp, 1)
		require.EqualValues(t, priceImport1, resp[0])
	}

	// check GetAll: filter for the 2nd one
	{
		resp, err := targetSt.GetAll(ctx, priceImport2.Timestamp, priceImport2.ProductID.Hex())
		require.NoError(t, err)
		require.Len(t, resp, 1)
		require.EqualValues(t, priceImport2, resp[0])
	}
}

func (s *StorageTestSuite) TestStorage_PriceImportAggregation() {
	t := s.T()
	ctx := context.Background()

	client := testutils.PrepareMongoDBFixtures(ctx, s.r, fixtures.NewDefaultMongoDBFixtures())
	storage, err := NewStorage(
		WithDatabase(testutils.TestMongoDBDatabase),
		WithMongoDBClient(client),
	)
	require.NoError(t, err)
	targetSt := storage.PriceImport()

	// build expected entries
	expEntries := model.PriceEntries{}
	for _, priceImport := range fixtures.PriceImports {
		product := fixtures.CheckFindProductByID(t, priceImport.ProductID)
		for _, price := range priceImport.Prices {
			expEntries = append(expEntries, model.PriceEntry{
				Name:      product.Name,
				Price:     price.Value,
				Timestamp: priceImport.Timestamp,
			})
		}
	}

	// helper which counts received entries which exists in the expEntries variable
	countRcvEntries := func(rcvEntries model.PriceEntries) int {
		count := 0
		for _, expEntry := range expEntries {
			for _, rcvEntry := range rcvEntries {
				if expEntry.Name != rcvEntry.Name {
					continue
				}
				if expEntry.Price != rcvEntry.Price {
					continue
				}
				if !expEntry.Timestamp.Equal(rcvEntry.Timestamp) {
					continue
				}

				count++
				break
			}
		}

		return count
	}

	// check GetPriceEntries: sort by price ASC
	{
		sortOpts := model.NewPriceEntriesSortOptions(
			model.PriceEntrySortByPrice(common.AscOrder),
		)
		pageOpt := common.NewPaginationOption(0, 100)

		rcvEntries, err := targetSt.GetPriceEntries(ctx, sortOpts, pageOpt)
		require.NoError(t, err)
		require.NotEmpty(t, rcvEntries)
		require.Equal(t, len(expEntries), countRcvEntries(rcvEntries))

		// check sorting
		sorted := sort.SliceIsSorted(rcvEntries, func(i, j int) bool {
			return rcvEntries[i].Price < rcvEntries[j].Price
		})
		require.True(t, sorted)
	}

	// check GetPriceEntries: sort by product name ASC
	{
		sortOpts := model.NewPriceEntriesSortOptions(
			model.PriceEntrySortByProductName(common.AscOrder),
		)
		pageOpt := common.NewPaginationOption(0, 100)

		rcvEntries, err := targetSt.GetPriceEntries(ctx, sortOpts, pageOpt)
		require.NoError(t, err)
		require.NotEmpty(t, rcvEntries)
		require.Equal(t, len(expEntries), countRcvEntries(rcvEntries))

		// check sorting
		sorted := sort.SliceIsSorted(rcvEntries, func(i, j int) bool {
			return rcvEntries[i].Name < rcvEntries[j].Name
		})
		require.True(t, sorted)
	}

	// check GetPriceEntries: sort by timestamp DESC
	{
		sortOpts := model.NewPriceEntriesSortOptions(
			model.PriceEntrySortByImportTimestamp(common.DescOrder),
		)
		pageOpt := common.NewPaginationOption(0, 100)

		rcvEntries, err := targetSt.GetPriceEntries(ctx, sortOpts, pageOpt)
		require.NoError(t, err)
		require.NotEmpty(t, rcvEntries)
		require.Equal(t, len(expEntries), countRcvEntries(rcvEntries))

		// check sorting
		sorted := sort.SliceIsSorted(rcvEntries, func(i, j int) bool {
			return rcvEntries[j].Timestamp.Before(rcvEntries[i].Timestamp)
		})
		require.True(t, sorted)
	}

	// check GetPriceEntries: sort by timestamp DESC and price DESC
	{
		sortOpts := model.NewPriceEntriesSortOptions(
			model.PriceEntrySortByImportTimestamp(common.DescOrder),
			model.PriceEntrySortByPrice(common.DescOrder),
		)
		pageOpt := common.NewPaginationOption(0, 100)

		rcvEntries, err := targetSt.GetPriceEntries(ctx, sortOpts, pageOpt)
		require.NoError(t, err)
		require.NotEmpty(t, rcvEntries)
		require.Equal(t, len(expEntries), countRcvEntries(rcvEntries))

		// check sorting
		sorted := sort.SliceIsSorted(rcvEntries, func(i, j int) bool {
			tsLess := rcvEntries[j].Timestamp.Before(rcvEntries[i].Timestamp)
			priceLess := rcvEntries[j].Price < rcvEntries[i].Price
			return tsLess && priceLess
		})
		require.True(t, sorted)
	}

	// check GetPriceEntries: limit
	{
		pageOpt := common.NewPaginationOption(0, 3)

		rcvEntries, err := targetSt.GetPriceEntries(ctx, nil, pageOpt)
		require.NoError(t, err)
		require.NotEmpty(t, rcvEntries)
		require.Equal(t, 3, countRcvEntries(rcvEntries))
	}

	// check GetPriceEntries: skip
	{
		pageOpt := common.NewPaginationOption(1, 100)

		rcvEntries, err := targetSt.GetPriceEntries(ctx, nil, pageOpt)
		require.NoError(t, err)
		require.NotEmpty(t, rcvEntries)
		require.Equal(t, len(expEntries)-1, countRcvEntries(rcvEntries))
	}
}
