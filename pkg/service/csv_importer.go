package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/itiky/mdb-tutorial/pkg/common"
	"github.com/itiky/mdb-tutorial/pkg/model"
	"github.com/itiky/mdb-tutorial/pkg/storage"
)

var _ CSVImporterService = (*csvImporterService)(nil)

// csvImporterService keeps CSVImporterService dependencies.
type csvImporterService struct {
	storage storage.Storage
	logger  *logrus.Logger
}

// ImportPrices implements CSVImporterService interface.
func (s csvImporterService) ImportPrices(ctx context.Context, csvImport model.CSVImport) error {
	// sanity check
	if csvImport.Timestamp.IsZero() {
		return fmt.Errorf("%w: csvImport.Timestmap: zero", common.ErrInvalidInput)
	}
	if len(csvImport.Entries) == 0 {
		return fmt.Errorf("%w: csvImport.Entries: empty", common.ErrInvalidInput)
	}

	// prepare map-reduce data (products with prices set)
	productsMap := make(map[string][]model.Price)
	for i, entry := range csvImport.Entries {
		// sanity check
		if entry.ProductName == "" {
			return fmt.Errorf("%w: csvImport.Entries[%d].ProductName (%s)", common.ErrInvalidInput, i, entry.ProductName)
		}
		if entry.Price < 0 {
			return fmt.Errorf("%w: csvImport.Entries[%d].Price (%d)", common.ErrInvalidInput, i, entry.Price)
		}

		// update set
		prices := productsMap[entry.ProductName]
		prices = append(prices, model.Price{
			Value: entry.Price,
		})
		productsMap[entry.ProductName] = prices
	}

	// reduce
	wg := sync.WaitGroup{}
	errsCh := make(chan error, len(productsMap))
	for productName, productPrices := range productsMap {
		wg.Add(1)
		go func(name string, prices []model.Price) {
			defer wg.Done()

			// upsert product
			productID, err := s.storage.Product().UpsertByName(ctx, model.Product{Name: name})
			if err != nil {
				errsCh <- fmt.Errorf("product %s: product object upsert failed: %w", name, err)
				return
			}

			// if product already exists (not inserted), fetch its ID
			if productID.IsZero() {
				product, err := s.storage.Product().GetByName(ctx, name)
				if err != nil {
					errsCh <- fmt.Errorf("product %s: fetch object failed: %w", name, err)
					return
				}
				productID = product.ID
			}

			// prepare PricesImport object and upsert it
			pricesImport := model.PricesImport{
				ProductID: productID,
				Timestamp: csvImport.Timestamp,
				Prices:    prices,
			}

			pricesImportID, err := s.storage.PriceImport().UpsertByProductIDAndTimestamp(ctx, pricesImport)
			if err != nil {
				errsCh <- fmt.Errorf("product %s: pricesImport object upsert failed: %w", name, err)
				return
			}

			if !pricesImportID.IsZero() {
				s.logger.Infof("CSV import: product %s import prices set: %s", name, csvImport.Timestamp)
			} else {
				s.logger.Warnf("CSV import: product %s import prices updated: %s", name, csvImport.Timestamp)
			}
		}(productName, productPrices)
	}

	// wait for workers to finish and checks for accumulated errors
	wg.Wait()
	close(errsCh)

	errs := make([]string, 0)
	for err := range errsCh {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return fmt.Errorf("partially failed (%d / %d): %s", len(productsMap), len(errs), strings.Join(errs, ", "))
	}

	return nil
}
