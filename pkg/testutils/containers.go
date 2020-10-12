package testutils

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	ContainerPortMongoDB = "27017/tcp"
)

type MongoDBContainer struct {
	tc.Container
	Db   string
	Port nat.Port
}

// Container is a wrapper object for test containers operations.
type Container struct{}

// MongoDBRequest creates a test container request for MongoDB start up.
func (c Container) MongoDBRequest(dbName string) tc.ContainerRequest {
	return tc.ContainerRequest{
		Image:        "mongo:4.4",
		ExposedPorts: []string{ContainerPortMongoDB},
		Env: map[string]string{
			"MONGO_INITDB_DATABASE": dbName,
		},
		WaitingFor: wait.ForListeningPort(ContainerPortMongoDB),
	}
}

// RunMongoDB start a new MongoDB test container.
func (c Container) RunMongoDB(ctx context.Context, req tc.ContainerRequest) (*MongoDBContainer, error) {
	cont, err := c.run(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("running container: %w", err)
	}

	resCont := MongoDBContainer{
		Container: cont,
		Db:        req.Env["MONGO_INITDB_DATABASE"],
	}

	port, err := cont.MappedPort(context.Background(), ContainerPortMongoDB)
	if err != nil {
		return nil, fmt.Errorf("getting mapped port %q: %w", ContainerPortMongoDB, err)
	}
	resCont.Port = port

	return &resCont, nil
}

// run starts a generic test container.
func (c Container) run(ctx context.Context, req tc.ContainerRequest) (tc.Container, error) {
	return tc.GenericContainer(ctx, tc.GenericContainerRequest{ContainerRequest: req, Started: true})
}
