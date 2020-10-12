package service

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/itiky/mdb-tutorial/pkg/storage"
)

var _ Service = (*service)(nil)

// service implements Service interface.
type service struct {
	storage storage.Storage
	logger  *logrus.Logger
}

// CSVImporterService implements Service interface.
// nolint:gosimple
func (s service) CSVImporter() CSVImporterService {
	return csvImporterService{
		storage: s.storage,
		logger:  s.logger,
	}
}

// CSVProcessor implements Service interface.
// nolint:gosimple
func (s service) CSVProcessor() CSVProcessorService {
	return csvProcessorService{
		logger: s.logger,
	}
}

// PriceEntries implements Service interface.
// nolint:gosimple
func (s service) PriceEntries() PriceEntriesService {
	return priceEntriesService{
		storage: s.storage,
		logger:  s.logger,
	}
}

// Option specifies functional argument used by NewService function.
type Option func(service *service) error

// WithStorage sets storage for service.
func WithStorage(storage storage.Storage) Option {
	return func(service *service) error {
		if storage == nil {
			return fmt.Errorf("storage option: nil")
		}
		service.storage = storage

		return nil
	}
}

// WithLogger sets logger for service.
func WithLogger(logger *logrus.Logger) Option {
	return func(service *service) error {
		if logger == nil {
			return fmt.Errorf("logger option: nil")
		}
		service.logger = logger

		return nil
	}
}

// NewService creates a new configured Service object.
func NewService(options ...Option) (Service, error) {
	s := &service{}
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	if s.logger == nil {
		logger := logrus.New()
		logger.SetOutput(ioutil.Discard)
		s.logger = logger
	}

	return s, nil
}
