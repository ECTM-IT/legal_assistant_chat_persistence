package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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
}

// NewSubscriptionsDAO creates a new SubscriptionsDAO
func NewSubscriptionsDAO(db *mongo.Database) *SubscriptionsDAO {
	return &SubscriptionsDAO{
		collection: db.Collection("subscriptions"),
	}
}

// GetAllSubscriptions retrieves all subscriptions from the database
func (dao *SubscriptionsDAO) GetAllSubscriptions(ctx context.Context) ([]models.Subscriptions, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []models.Subscriptions
	if err := cursor.All(ctx, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// GetSubscriptionByID retrieves a subscription by its ID from the database
func (dao *SubscriptionsDAO) GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*models.Subscriptions, error) {
	var subscription models.Subscriptions
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("subscription not found")
		}
		return nil, err
	}
	return &subscription, nil
}

// GetSubscriptionsByPlan retrieves subscriptions by their plan from the database
func (dao *SubscriptionsDAO) GetSubscriptionsByPlan(ctx context.Context, plan string) ([]models.Subscriptions, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{"plan": plan})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []models.Subscriptions
	if err := cursor.All(ctx, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// CreateSubscription creates a new subscription in the database
func (dao *SubscriptionsDAO) CreateSubscription(ctx context.Context, subscription *models.Subscriptions) (*mongo.InsertOneResult, error) {
	result, err := dao.collection.InsertOne(ctx, subscription)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateSubscription updates an existing subscription in the database
func (dao *SubscriptionsDAO) UpdateSubscription(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteSubscription deletes a subscription by its ID from the database
func (dao *SubscriptionsDAO) DeleteSubscription(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return result, nil
}
