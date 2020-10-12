package testutils

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/itiky/mdb-tutorial/pkg/mongodb"
	"github.com/itiky/mdb-tutorial/pkg/testutils/fixtures"
)

const (
	TestMongoDBDatabase = "db"
)

// Resources keeps test environment data.
type Resources struct {
	sync.Mutex

	// MongoDB
	mdbOnce     sync.Once
	mdbCont     *MongoDBContainer
	mdbClient   *mongo.Client
	mdbFixtures fixtures.MongoDB

	// Logger
	logLvl logrus.Level
}

// MongoDB starts a MongoDB container and returns a client connection.
func (r *Resources) MongoDB(ctx context.Context) (*MongoDBContainer, *mongo.Client) {
	r.mdbOnce.Do(func() {
		c := Container{}

		cont, err := c.RunMongoDB(ctx, c.MongoDBRequest(TestMongoDBDatabase))
		if err != nil {
			r.crash(ctx, fmt.Errorf("staring MongoDB container: %w", err))
		}
		r.mdbCont = cont

		client, err := mongodb.Connect(mongodb.Configuration{
			Url:  "localhost",
			Port: r.mdbCont.Port.Port(),
		})
		if err != nil {
			r.crash(ctx, fmt.Errorf("connecting to MongoDB container: %w", err))
		}
		r.mdbClient = client
	})

	return r.mdbCont, r.mdbClient
}

// Shutdown shutdowns all test container.
func (r *Resources) Shutdown(ctx context.Context) {
	r.Lock()
	defer r.Unlock()

	if r.mdbClient != nil {
		_ = r.mdbClient.Disconnect(ctx)
	}
	if r.mdbCont != nil {
		_ = r.mdbCont.Terminate(ctx)
	}
}

// logger returns a configured logger.
func (r *Resources) Logger() *logrus.Logger {
	l := logrus.StandardLogger()
	l.SetLevel(r.logLvl)

	return l
}

// crash crashes the current test environment.
func (r *Resources) crash(ctx context.Context, err error) {
	r.Shutdown(ctx)
	r.Logger().Fatal(err)
}

// ResourcesOption is a functional argument used by NewResources func.
type ResourcesOption func(r *Resources)

// WithLogLevel sets logger log level.
func WithLogLevel(logLvl logrus.Level) ResourcesOption {
	return func(r *Resources) {
		r.logLvl = logLvl
	}
}

// NewResources creates a new Resources object.
func NewResources(options ...ResourcesOption) *Resources {
	r := &Resources{
		logLvl: logrus.DebugLevel,
	}
	for _, option := range options {
		option(r)
	}

	return r
}
