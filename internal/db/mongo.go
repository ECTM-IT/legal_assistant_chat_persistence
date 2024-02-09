package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to the MongoDB database.
func Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use the uri parameter to connect to MongoDB
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err // return here if connection fails
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(ctx) // Attempt to disconnect if ping fails
		return nil, err
	}
	// Return the client to be used elsewhere, don't disconnect here
	return client, err
}

func StartSession(uri string) (mongo.Session, error) {
	client, err := Connect(uri)

	if err != nil {
		return nil, err
	}

	return client.StartSession()
}
