package db

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to the MongoDB database.
func Connect(uri string, timeout time.Duration, logger logs.Logger) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		logger.Error("Failed to connect to MongoDB", err)
		return nil, err
	}

	err = pingDatabase(ctx, client, logger)
	if err != nil {
		logger.Error("Failed to ping MongoDB", err)
		return disconnectClient(ctx, client, logger)
	}

	logger.Info("Successfully connected to MongoDB")
	return client, nil
}

func pingDatabase(ctx context.Context, client *mongo.Client, logger logs.Logger) error {
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		logger.Error("Failed to ping MongoDB", err)
		return err
	}

	logger.Info("Successfully pinged MongoDB")
	return nil
}

func CreateDB(client *mongo.Client, logger logs.Logger) *mongo.Database {
	logger.Info("Creating database: laDB")
	return client.Database("laDB")
}

func disconnectClient(ctx context.Context, client *mongo.Client, logger logs.Logger) (*mongo.Client, error) {
	disconnectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Disconnect(disconnectCtx); err != nil {
		logger.Error("Failed to disconnect from MongoDB", err)
		return nil, err
	}

	logger.Info("Successfully disconnected from MongoDB")
	return client, nil
}

func StartSession(uri string, timeout time.Duration, logger logs.Logger) (mongo.Session, error) {
	client, err := Connect(uri, timeout, logger)
	if err != nil {
		logger.Error("Failed to start MongoDB session", err)
		return nil, err
	}

	logger.Info("Starting MongoDB session")
	return client.StartSession()
}
