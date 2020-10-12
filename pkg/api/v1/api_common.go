package v1

import (
	"fmt"

	"github.com/itiky/mdb-tutorial/pkg/common"
)

// NewPaginationOption converts gRPC PaginationParams to common.PaginationOption.
func NewPaginationOption(apiParams *PaginationParams) (common.PaginationOption, error) {
	if apiParams == nil {
		return common.PaginationOption{}, fmt.Errorf("apiParams is nil")
	}

	opt := common.NewPaginationOption(int(apiParams.Skip), int(apiParams.Limit))
	if err := opt.Validate(); err != nil {
		return common.PaginationOption{}, fmt.Errorf("pagination params validation failed: %v", err)
	}

	return opt, nil
}

// NewOrderOption converts gRPC SortOrder to common.Order and returns ok if sort is specified.
func NewOrderOption(apiParam SortOrder) (common.Order, bool) {
	switch apiParam {
	case SortOrder_Asc:
		return common.AscOrder, true
	case SortOrder_Desc:
		return common.DescOrder, true
	default:
		return common.AscOrder, false
	}
}
