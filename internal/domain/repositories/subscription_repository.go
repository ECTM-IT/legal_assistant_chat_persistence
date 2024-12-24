package repositories

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SubscriptionRepository defines the operations available on a subscription repository.
type SubscriptionRepository interface {
	FindAll(ctx context.Context) ([]models.Subscriptions, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Subscriptions, error)
	FindByPlan(ctx context.Context, plan string) ([]models.Subscriptions, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.Subscriptions, error)
	Create(ctx context.Context, subscription *models.Subscriptions) (*models.Subscriptions, error)
	Update(ctx context.Context, id primitive.ObjectID, subscription *bson.M) (*models.Subscriptions, error)
	Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
}

// SubscriptionRepositoryImpl implements the SubscriptionRepository interface.
type SubscriptionRepositoryImpl struct {
	subscriptionDAO daos.SubscriptionsDAOInterface
}

// NewSubscriptionRepository creates a new instance of the subscription repository.
func NewSubscriptionRepository(subscriptionDAO daos.SubscriptionsDAOInterface) *SubscriptionRepositoryImpl {
	return &SubscriptionRepositoryImpl{
		subscriptionDAO: subscriptionDAO,
	}
}

// FindAll retrieves all subscriptions.
func (r *SubscriptionRepositoryImpl) FindAll(ctx context.Context) ([]models.Subscriptions, error) {
	return r.subscriptionDAO.GetAllSubscriptions(ctx)
}

// FindByID retrieves a subscription by its ID.
func (r *SubscriptionRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Subscriptions, error) {
	return r.subscriptionDAO.GetSubscriptionByID(ctx, id)
}

// FindByPlan retrieves subscriptions by their plan.
func (r *SubscriptionRepositoryImpl) FindByPlan(ctx context.Context, plan string) ([]models.Subscriptions, error) {
	return r.subscriptionDAO.GetSubscriptionsByPlan(ctx, plan)
}

// FindByUserID retrieves subscriptions by user ID.
func (r *SubscriptionRepositoryImpl) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.Subscriptions, error) {
	return r.subscriptionDAO.GetSubscriptionsByUserID(ctx, userID)
}

// Create creates a new subscription.
func (r *SubscriptionRepositoryImpl) Create(ctx context.Context, req *models.Subscriptions) (*models.Subscriptions, error) {
	_, err := r.subscriptionDAO.CreateSubscription(ctx, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Update updates an existing subscription.
func (r *SubscriptionRepositoryImpl) Update(ctx context.Context, id primitive.ObjectID, req bson.M) (*models.Subscriptions, error) {
	_, err := r.subscriptionDAO.UpdateSubscription(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

// Delete deletes a subscription by its ID.
func (r *SubscriptionRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return r.subscriptionDAO.DeleteSubscription(ctx, id)
}
