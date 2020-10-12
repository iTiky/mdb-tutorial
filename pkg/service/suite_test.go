package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/itiky/mdb-tutorial/pkg/testutils"
)

// ServiceTestSuite implements testify.Suite interface and used for Service tests.
type ServiceTestSuite struct {
	suite.Suite
	r *testutils.Resources
}

func (s *ServiceTestSuite) SetupSuite() {
	s.r = testutils.NewResources()
}

// nolint:errcheck
func (s *ServiceTestSuite) TearDownSuite() {
	if s.r != nil {
		s.r.Shutdown(context.Background())
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
