package repositories

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StripeEventRepository defines the operations available on a stripe event repository.
type StripeEventRepository interface {
	FindByEventID(ctx context.Context, eventID string) (*models.StripeEvent, error)
	Create(ctx context.Context, event *models.StripeEvent) (*models.StripeEvent, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.StripeEvent, error)
}

// StripeEventRepositoryImpl implements the StripeEventRepository interface.
type StripeEventRepositoryImpl struct {
	stripeEventDAO daos.StripeEventDAOInterface
}

// NewStripeEventRepository creates a new instance of the stripe event repository.
func NewStripeEventRepository(stripeEventDAO daos.StripeEventDAOInterface) *StripeEventRepositoryImpl {
	return &StripeEventRepositoryImpl{
		stripeEventDAO: stripeEventDAO,
	}
}

// FindByEventID retrieves a stripe event by its Stripe event ID.
func (r *StripeEventRepositoryImpl) FindByEventID(ctx context.Context, eventID string) (*models.StripeEvent, error) {
	return r.stripeEventDAO.FindByEventID(ctx, eventID)
}

// Create creates a new stripe event.
func (r *StripeEventRepositoryImpl) Create(ctx context.Context, event *models.StripeEvent) (*models.StripeEvent, error) {
	event.ID = primitive.NewObjectID()

	result, err := r.stripeEventDAO.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	// Set the ID from the result
	event.ID = result.InsertedID.(primitive.ObjectID)

	return event, nil
}

// Update updates an existing stripe event.
func (r *StripeEventRepositoryImpl) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.StripeEvent, error) {
	_, err := r.stripeEventDAO.UpdateEvent(ctx, id, update)
	if err != nil {
		return nil, err
	}

	// This is a simplification - we'd normally retrieve the updated document
	// but for our purposes we don't need the updated document
	return nil, nil
}
