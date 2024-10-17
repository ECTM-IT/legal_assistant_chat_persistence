package db

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Connect establishes a connection to the MongoDB database.
func Connect(uri string, timeout time.Duration, logger logs.Logger) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", err)
		return nil, err
	}

	if err = pingDatabase(ctx, client, logger); err != nil {
		logger.Error("Failed to ping MongoDB", err)
		return nil, disconnectClient(ctx, client, logger)
	}

	logger.Info("Successfully connected to MongoDB")
	return client, nil
}

// pingDatabase pings the MongoDB database to ensure the connection is established.
func pingDatabase(ctx context.Context, client *mongo.Client, logger logs.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		logger.Error("Failed to ping MongoDB", err)
		return err
	}

	logger.Info("Successfully pinged MongoDB")
	return nil
}

// CreateDB returns a reference to the specified database.
func CreateDB(client *mongo.Client, dbName string, logger logs.Logger) *mongo.Database {
	logger.Info("Creating database", zap.String("database", dbName))
	return client.Database(dbName)
}

// disconnectClient disconnects from the MongoDB database.
func disconnectClient(ctx context.Context, client *mongo.Client, logger logs.Logger) error {
	disconnectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Disconnect(disconnectCtx); err != nil {
		logger.Error("Failed to disconnect from MongoDB", err)
		return err
	}

	logger.Info("Successfully disconnected from MongoDB")
	return nil
}

// StartSession starts a new session with MongoDB.
func StartSession(uri string, timeout time.Duration, logger logs.Logger) (mongo.Session, error) {
	client, err := Connect(uri, timeout, logger)
	if err != nil {
		logger.Error("Failed to start MongoDB session", err)
		return nil, err
	}

	session, err := client.StartSession()
	if err != nil {
		logger.Error("Failed to start MongoDB session", err)
		return nil, err
	}

	logger.Info("Successfully started MongoDB session")
	return session, nil
}
