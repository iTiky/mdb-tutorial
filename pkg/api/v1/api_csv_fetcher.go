package v1

import (
	"context"
	"errors"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/itiky/mdb-tutorial/pkg/common"
)

// Fetch implements CSVFetcherServer interface.
func (s gRPCServer) Fetch(ctx context.Context, req *CSVFetchRequest) (*CSVFetchResponse, error) {
	// download file
	tmpFilePath, importTimestamp, err := s.service.CSVProcessor().Download(req.Url)
	if err != nil {
		if errors.Is(err, common.ErrInvalidInput) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, common.ErrNotFound) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// cleanup
	defer func() {
		os.Remove(tmpFilePath)
	}()

	// open tmp file
	file, err := os.Open(tmpFilePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "tmp file open failed: %v", err)
	}
	defer file.Close()

	// process with CSVImporter service handler
	err = s.service.CSVProcessor().Process(ctx, file, importTimestamp, s.csvChunkSize, s.service.CSVImporter().ImportPrices)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "processing: %v", err)
	}

	return &CSVFetchResponse{}, nil
}
