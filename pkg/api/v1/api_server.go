package v1

import (
	"crypto/tls"
	fmt "fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/itiky/mdb-tutorial/pkg/service"
)

var _ CSVFetcherServer = (*gRPCServer)(nil)
var _ PriceEntryReaderServer = (*gRPCServer)(nil)

// gRPCServer implement gRPC services.
type gRPCServer struct {
	service service.Service
	logger  *logrus.Logger
	//
	csvChunkSize   int
	tlsCertificate *tls.Certificate
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

// WithTLS enabled TLS encryption fro server.
func WithTLS(certificate *tls.Certificate) Option {
	return func(server *gRPCServer) error {
		server.tlsCertificate = certificate

		return nil
	}
}

// NewServer creates a new configured gRPCServer object.
func NewServer(options ...Option) (*grpc.Server, error) {
	var serverOptions []grpc.ServerOption

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

	// TLS option
	if s.tlsCertificate != nil {
		transportCreds := credentials.NewServerTLSFromCert(s.tlsCertificate)
		serverOptions = append(serverOptions, grpc.Creds(transportCreds))
		s.logger.Infof("gRPC server: using TLS")
	} else {
		s.logger.Infof("gRPC server: insecure")
	}

	gRPCServer := grpc.NewServer(serverOptions...)
	RegisterCSVFetcherServer(gRPCServer, s)
	RegisterPriceEntryReaderServer(gRPCServer, s)

	return gRPCServer, nil
}
