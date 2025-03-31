package daos

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// StripeEventDAOInterface defines the interface for the StripeEventDAO
type StripeEventDAOInterface interface {
	FindByEventID(ctx context.Context, eventID string) (*models.StripeEvent, error)
	CreateEvent(ctx context.Context, event *models.StripeEvent) (*mongo.InsertOneResult, error)
	UpdateEvent(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error)
}

// StripeEventDAO implements the StripeEventDAOInterface
type StripeEventDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewStripeEventDAO creates a new StripeEventDAO
func NewStripeEventDAO(db *mongo.Database, logger logs.Logger) *StripeEventDAO {
	return &StripeEventDAO{
		collection: db.Collection("stripe_events"),
		logger:     logger,
	}
}

// FindByEventID retrieves a stripe event by its Stripe event ID
func (dao *StripeEventDAO) FindByEventID(ctx context.Context, eventID string) (*models.StripeEvent, error) {
	filter := bson.M{"event_id": eventID}

	var event models.StripeEvent
	err := dao.collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No event found with this ID
		}
		dao.logger.Error("Failed to find stripe event by event ID", err)
		return nil, err
	}

	return &event, nil
}

// CreateEvent stores a new stripe event in the database
func (dao *StripeEventDAO) CreateEvent(ctx context.Context, event *models.StripeEvent) (*mongo.InsertOneResult, error) {
	result, err := dao.collection.InsertOne(ctx, event)
	if err != nil {
		dao.logger.Error("Failed to create stripe event", err)
		return nil, err
	}

	return result, nil
}

// UpdateEvent updates a stripe event in the database
func (dao *StripeEventDAO) UpdateEvent(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	updateDoc := bson.M{"$set": update}

	result, err := dao.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		dao.logger.Error("Failed to update stripe event", err)
		return nil, err
	}

	return result, nil
}
