package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/itiky/mdb-tutorial/pkg/testutils"
)

// StorageTestSuite implements testify.Suite interface and used for Storage tests.
type StorageTestSuite struct {
	suite.Suite
	r *testutils.Resources
}

func (s *StorageTestSuite) SetupSuite() {
	s.r = testutils.NewResources()
}

// nolint:errcheck
func (s *StorageTestSuite) TearDownSuite() {
	if s.r != nil {
		s.r.Shutdown(context.Background())
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
