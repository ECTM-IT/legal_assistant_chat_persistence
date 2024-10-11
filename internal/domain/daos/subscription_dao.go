package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SubscriptionsDAOInterface defines the interface for the SubscriptionsDAO
type SubscriptionsDAOInterface interface {
	GetAllSubscriptions(ctx context.Context) ([]models.Subscriptions, error)
	GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*models.Subscriptions, error)
	GetSubscriptionsByPlan(ctx context.Context, plan string) ([]models.Subscriptions, error)
	CreateSubscription(ctx context.Context, subscription *models.Subscriptions) (*mongo.InsertOneResult, error)
	UpdateSubscription(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error)
	DeleteSubscription(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
}

// SubscriptionsDAO implements the SubscriptionsDAOInterface
type SubscriptionsDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewSubscriptionsDAO creates a new SubscriptionsDAO
func NewSubscriptionsDAO(db *mongo.Database, logger logs.Logger) *SubscriptionsDAO {
	return &SubscriptionsDAO{
		collection: db.Collection("subscriptions"),
		logger:     logger,
	}
}

// GetAllSubscriptions retrieves all subscriptions from the database
func (dao *SubscriptionsDAO) GetAllSubscriptions(ctx context.Context) ([]models.Subscriptions, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve all subscriptions")
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve subscriptions", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []models.Subscriptions
	if err := cursor.All(ctx, &subscriptions); err != nil {
		dao.logger.Error("DAO Level: Failed to decode subscriptions", err)
		return nil, err
	}

	dao.logger.Info("DAO Level: Successfully retrieved all subscriptions")
	return subscriptions, nil
}

// GetSubscriptionByID retrieves a subscription by its ID from the database
func (dao *SubscriptionsDAO) GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*models.Subscriptions, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve subscription by ID")
	var subscription models.Subscriptions
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("Subscription not found")
			return nil, errors.New("subscription not found")
		}
		dao.logger.Error("DAO Level: Failed to retrieve subscription", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved subscription")
	return &subscription, nil
}

// GetSubscriptionsByPlan retrieves subscriptions by their plan from the database
func (dao *SubscriptionsDAO) GetSubscriptionsByPlan(ctx context.Context, plan string) ([]models.Subscriptions, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve subscriptions by plan")
	cursor, err := dao.collection.Find(ctx, bson.M{"plan": plan})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve subscriptions by plan", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []models.Subscriptions
	if err := cursor.All(ctx, &subscriptions); err != nil {
		dao.logger.Error("DAO Level: Failed to decode subscriptions", err)
		return nil, err
	}

	dao.logger.Info("DAO Level: Successfully retrieved subscriptions by plan")
	return subscriptions, nil
}

// CreateSubscription creates a new subscription in the database
func (dao *SubscriptionsDAO) CreateSubscription(ctx context.Context, subscription *models.Subscriptions) (*mongo.InsertOneResult, error) {
	dao.logger.Info("DAO Level: Attempting to create new subscription")
	result, err := dao.collection.InsertOne(ctx, subscription)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to create subscription", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully created new subscription")
	return result, nil
}

// UpdateSubscription updates an existing subscription in the database
func (dao *SubscriptionsDAO) UpdateSubscription(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to update subscription")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to update subscription", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully updated subscription")
	return result, nil
}

// DeleteSubscription deletes a subscription by its ID from the database
func (dao *SubscriptionsDAO) DeleteSubscription(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	dao.logger.Info("DAO Level: Attempting to delete subscription")
	result, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to delete subscription", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully deleted subscription")
	return result, nil
}
