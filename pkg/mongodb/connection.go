package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ConnectTimeout = 2 * time.Second
	PintTimeout    = 1 * time.Second
)

// Configuration keeps MongoDB connection options.
type Configuration struct {
	Url  string
	Port string
}

// Connect creates a new MongoDB client.
func Connect(config Configuration) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s/", config.Url, config.Port)

	client, err := mongo.NewClient(
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}

	connectCtx, connectCancel := context.WithTimeout(context.Background(), ConnectTimeout)
	defer connectCancel()
	if err := client.Connect(connectCtx); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), PintTimeout)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return client, nil
}
