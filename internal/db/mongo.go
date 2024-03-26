package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to the MongoDB database.
func Connect(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = pingDatabase(ctx, client)
	if err != nil {
		return disconnectClient(ctx, client)
	}

	return client, nil
}

func pingDatabase(ctx context.Context, client *mongo.Client) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := client.Ping(pingCtx, nil)
	if err != nil {
		return err
	}

	return nil
}

func CreateDB(client *mongo.Client) *mongo.Database {
	return client.Database("laDB")
}

func disconnectClient(ctx context.Context, client *mongo.Client) (*mongo.Client, error) {
	disconnectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Disconnect(disconnectCtx); err != nil {
		return nil, err
	}
	return client, nil
}

func StartSession(uri string, timeout time.Duration) (mongo.Session, error) {
	client, err := Connect(uri, timeout)
	if err != nil {
		return nil, err
	}

	return client.StartSession()
}
