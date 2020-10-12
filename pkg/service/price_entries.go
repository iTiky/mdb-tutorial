package service

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/model"
	"github.com/itiky/mdb-tutorial/pkg/storage"
)

var _ PriceEntriesService = (*priceEntriesService)(nil)

type priceEntriesService struct {
	storage storage.Storage
	logger  *logrus.Logger
}

// List implements PriceEntriesService interface.
func (s priceEntriesService) List(ctx context.Context, paginationOpt common.PaginationOption, sortOpts common.SortOptions) (model.PriceEntries, error) {
	return s.storage.PriceImport().GetPriceEntries(ctx, sortOpts, paginationOpt)
}
