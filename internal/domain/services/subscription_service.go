package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubscriptionService defines the subscription service interface.
type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, id primitive.ObjectID, req *dtos.UpdateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
	GetAllSubscriptions(ctx context.Context) ([]dtos.SubscriptionResponse, error)
	GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*dtos.SubscriptionResponse, error)
	GetSubscriptionsByPlan(ctx context.Context, plan string) ([]dtos.SubscriptionResponse, error)
	DeleteSubscription(ctx context.Context, id primitive.ObjectID) error
}

// SubscriptionServiceImpl implements the SubscriptionService interface.
type SubscriptionServiceImpl struct {
	repo   *repositories.SubscriptionRepositoryImpl
	mapper *mappers.SubscriptionConversionServiceImpl
}

// NewSubscriptionService creates a new instance of the subscription service.
func NewSubscriptionService(repo *repositories.SubscriptionRepositoryImpl, mapper *mappers.SubscriptionConversionServiceImpl) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{repo: repo, mapper: mapper}
}

// CreateSubscription handles the business logic for creating a subscription.
func (s *SubscriptionServiceImpl) CreateSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	subscription, err := s.mapper.DTOToSubscription(req)
	if err != nil {
		return nil, err
	}

	createdSubscription, err := s.repo.Create(ctx, subscription)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create subscription", "create_subscription_failed")
	}

	return s.mapper.SubscriptionToDTO(createdSubscription), nil
}

// UpdateSubscription handles the business logic for updating a subscription.
func (s *SubscriptionServiceImpl) UpdateSubscription(ctx context.Context, id primitive.ObjectID, req *dtos.UpdateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	updateFields := s.mapper.UpdateSubscriptionFieldsToMap(*req)

	_, err := s.repo.Update(ctx, id, updateFields)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to update subscription", "update_subscription_failed")
	}

	return s.GetSubscriptionByID(ctx, id)
}

// GetAllSubscriptions retrieves all subscriptions.
func (s *SubscriptionServiceImpl) GetAllSubscriptions(ctx context.Context) ([]dtos.SubscriptionResponse, error) {
	subscriptions, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get all subscriptions", "get_all_subscriptions_failed")
	}

	return s.mapper.SubscriptionsToDTO(subscriptions), nil
}

// GetSubscriptionByID retrieves a subscription by its ID.
func (s *SubscriptionServiceImpl) GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*dtos.SubscriptionResponse, error) {
	subscription, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get subscription", "get_subscription_failed")
	}

	return s.mapper.SubscriptionToDTO(subscription), nil
}

// GetSubscriptionsByPlan retrieves subscriptions by their plan.
func (s *SubscriptionServiceImpl) GetSubscriptionsByPlan(ctx context.Context, plan string) ([]dtos.SubscriptionResponse, error) {
	subscriptions, err := s.repo.FindByPlan(ctx, plan)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get subscriptions by plan", "get_subscriptions_by_plan_failed")
	}

	return s.mapper.SubscriptionsToDTO(subscriptions), nil
}

// DeleteSubscription deletes a subscription by its ID.
func (s *SubscriptionServiceImpl) DeleteSubscription(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.repo.Delete(ctx, id)
	if err != nil {
		return errors.NewDatabaseError("Failed to delete subscription", "delete_subscription_failed")
	}

	return nil
}
