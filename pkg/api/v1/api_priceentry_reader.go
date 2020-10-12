package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/itiky/mdb-tutorial/pkg/model"
)

// List implements PriceEntryReaderServer interface.
func (s gRPCServer) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	// parse inputs
	paginationOption, err := NewPaginationOption(req.Pagination)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	peSortOptions := make([]model.PriceEntriesSortOption, 0)
	if order, ok := NewOrderOption(req.SortByName); ok {
		peSortOptions = append(peSortOptions, model.PriceEntrySortByProductName(order))
	}
	if order, ok := NewOrderOption(req.SortByPrice); ok {
		peSortOptions = append(peSortOptions, model.PriceEntrySortByPrice(order))
	}
	if order, ok := NewOrderOption(req.SortByTimestamp); ok {
		peSortOptions = append(peSortOptions, model.PriceEntrySortByImportTimestamp(order))
	}
	sortOptions := model.NewPriceEntriesSortOptions(peSortOptions...)

	// query and build response
	entries, err := s.service.PriceEntries().List(ctx, paginationOption, sortOptions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := &ListResponse{
		Entries: NewPriceEntries(entries),
	}

	return response, nil
}

// NewPriceEntries converts model.PriceEntries to gRPC PriceEntry list.
func NewPriceEntries(inEntries model.PriceEntries) (outEntries []*PriceEntry) {
	for _, inEntry := range inEntries {
		outEntries = append(outEntries, &PriceEntry{
			ProductName: inEntry.Name,
			Timestamp:   inEntry.Timestamp.Unix(),
			Price:       int32(inEntry.Price),
		})
	}

	return
}
