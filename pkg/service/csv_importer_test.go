package service

import (
	"context"
	"errors"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/model"
	"github.com/itiky/mdb-tutorial/pkg/storage"
	"github.com/itiky/mdb-tutorial/pkg/testutils"
	"github.com/itiky/mdb-tutorial/pkg/testutils/fixtures"
)

func (s *ServiceTestSuite) TestService_CSVImporter() {
	t := s.T()
	ctx := context.Background()

	client := testutils.PrepareMongoDBFixtures(ctx, s.r, fixtures.NewEmptyMongoDBFixtures())
	svcStorage, err := storage.NewStorage(
		storage.WithDatabase(testutils.TestMongoDBDatabase),
		storage.WithMongoDBClient(client),
	)
	require.NoError(t, err)

	service, err := NewService(
		WithStorage(svcStorage),
	)
	require.NoError(t, err)
	targetSvc := service.CSVImporter()

	// check ImportPrices: invalid timestamp
	{
		csvImport := model.CSVImport{
			Timestamp: time.Time{},
			Entries:   model.CSVEntries{model.CSVEntry{ProductName: "name", Price: 100}},
		}

		err := targetSvc.ImportPrices(ctx, csvImport)
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrInvalidInput))
	}

	// check ImportPrices: no entries
	{
		csvImport := model.CSVImport{
			Timestamp: time.Now(),
			Entries:   nil,
		}

		err := targetSvc.ImportPrices(ctx, csvImport)
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrInvalidInput))
	}

	// check ImportPrices: empty productName
	{
		csvImport := model.CSVImport{
			Timestamp: time.Now(),
			Entries:   model.CSVEntries{model.CSVEntry{ProductName: "", Price: 100}},
		}

		err := targetSvc.ImportPrices(ctx, csvImport)
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrInvalidInput))
	}

	// check ImportPrices: invalid price
	{
		csvImport := model.CSVImport{
			Timestamp: time.Now(),
			Entries:   model.CSVEntries{model.CSVEntry{ProductName: "name", Price: -1}},
		}

		err := targetSvc.ImportPrices(ctx, csvImport)
		require.Error(t, err)
		require.True(t, errors.Is(err, common.ErrInvalidInput))
	}

	// check ImportPrices: ok
	{
		csvImport := model.CSVImport{
			Timestamp: time.Now(),
			Entries: model.CSVEntries{
				model.CSVEntry{ProductName: "Product Z", Price: 0},
				model.CSVEntry{ProductName: "Product X", Price: 5},
				model.CSVEntry{ProductName: "Product Y", Price: 50},
				model.CSVEntry{ProductName: "Product X", Price: 10},
				model.CSVEntry{ProductName: "Product Y", Price: 55},
				model.CSVEntry{ProductName: "Product Y", Price: 60},
				model.CSVEntry{ProductName: "Product Z", Price: 100},
			},
		}

		err := targetSvc.ImportPrices(ctx, csvImport)
		require.NoError(t, err)

		products, err := svcStorage.Product().GetAll(ctx)
		require.NoError(t, err)
		require.Len(t, products, 3)

		priceImports, err := svcStorage.PriceImport().GetAll(ctx, time.Time{}, "")
		require.NoError(t, err)
		require.Len(t, priceImports, 3)
	}
}
