package service

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/itiky/mdb-tutorial/pkg/model"
)

var mockCSV = `Product_1;1
Product_2;1
Product_1;2
Product_1;3
Product_2;2
Product_2;3
Product_3;1
Product_3;2
Product_1;4
Product_2;4
`

func (s *ServiceTestSuite) TestService_CSVProcessor_Download() {
	t := s.T()

	// as it uses an external resource, should be executed manually
	t.Skip()

	service, err := NewService()
	require.NoError(t, err)
	targetSvc := service.CSVProcessor()

	// check Download: invalid url
	{
		_, _, err := targetSvc.Download("")
		require.Error(t, err)
	}

	// check Download: ok
	{
		filePath, timestamp, err := targetSvc.Download("https://people.sc.fsu.edu/~jburkardt/data/csv/addresses.csv")
		require.NoError(t, err)
		require.NotEmpty(t, filePath)
		require.False(t, timestamp.IsZero())

		fStat, err := os.Stat(filePath)
		require.NoError(t, err)
		require.NotZero(t, fStat.Size())

		os.Remove(filePath)
	}
}

func (s *ServiceTestSuite) TestService_CSVProcessor_Process() {
	t := s.T()
	ctx := context.Background()

	service, err := NewService()
	require.NoError(t, err)
	targetSvc := service.CSVProcessor()

	reader := strings.NewReader(mockCSV)
	timestamp := time.Now()
	chunkSize := 3

	// mockChunkWorker save import for later check
	processedImports := make([]model.CSVImport, 0)
	mockChunkWorker := func(ctx context.Context, csvImport model.CSVImport) error {
		processedImports = append(processedImports, csvImport)
		return nil
	}

	// check Process: nil reader
	{
		err := targetSvc.Process(ctx, nil, timestamp, chunkSize, mockChunkWorker)
		require.Error(t, err)
	}

	// check Process: nil worker
	{
		err := targetSvc.Process(ctx, reader, timestamp, chunkSize, nil)
		require.Error(t, err)
	}

	// check Process: invalid chunk size
	{
		err := targetSvc.Process(ctx, reader, timestamp, 0, mockChunkWorker)
		require.Error(t, err)
	}

	// check Process: ok
	{
		err := targetSvc.Process(ctx, reader, timestamp, chunkSize, mockChunkWorker)
		require.NoError(t, err)

		// check imports
		totalEntries := 0
		require.Len(t, processedImports, 4)
		for _, processedImport := range processedImports {
			require.True(t, processedImport.Timestamp.Equal(timestamp))
			require.NotEmpty(t, processedImport.Entries)
			for _, entry := range processedImport.Entries {
				require.NotEmpty(t, entry.ProductName)
				require.Greater(t, entry.Price, 0)
				totalEntries++
			}
		}
		require.Equal(t, totalEntries, 10)
	}
}
