package v1

import (
	fmt "fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/itiky/mdb-tutorial/pkg/service"
)

var _ CSVFetcherServer = (*gRPCServer)(nil)
var _ PriceEntryReaderServer = (*gRPCServer)(nil)

// gRPCServer implement gRPC services.
type gRPCServer struct {
	service service.Service
	logger  *logrus.Logger
	//
	csvChunkSize int
}

func (s gRPCServer) mustEmbedUnimplementedCSVFetcherServer()       {}
func (s gRPCServer) mustEmbedUnimplementedPriceEntryReaderServer() {}

// Option specifies functional argument used by NewServer function.
type Option func(server *gRPCServer) error

// WithService sets service for server.
func WithService(service service.Service) Option {
	return func(server *gRPCServer) error {
		if service == nil {
			return fmt.Errorf("service option: nil")
		}
		server.service = service

		return nil
	}
}

// WithLogger sets logger for server.
func WithLogger(logger *logrus.Logger) Option {
	return func(server *gRPCServer) error {
		if logger == nil {
			return fmt.Errorf("logger option: nil")
		}
		server.logger = logger

		return nil
	}
}

// WithCSVChunkSize sets CSV-processor chunk size for server.
func WithCSVChunkSize(size int) Option {
	return func(server *gRPCServer) error {
		server.csvChunkSize = size

		return nil
	}
}

// NewServer creates a new configured gRPCServer object.
func NewServer(options ...Option) (*grpc.Server, error) {
	s := &gRPCServer{
		csvChunkSize: 1000,
	}
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

	gRPCServer := grpc.NewServer()
	RegisterCSVFetcherServer(gRPCServer, s)
	RegisterPriceEntryReaderServer(gRPCServer, s)

	return gRPCServer, nil
}
