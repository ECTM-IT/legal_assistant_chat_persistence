package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
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
	PurchaseSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
}

// SubscriptionServiceImpl implements the SubscriptionService interface.
type SubscriptionServiceImpl struct {
	repo     *repositories.SubscriptionRepositoryImpl
	userRepo *repositories.UserRepositoryImpl
	mapper   *mappers.SubscriptionConversionServiceImpl
	logger   logs.Logger
}

// NewSubscriptionService creates a new instance of the subscription service.
func NewSubscriptionService(repo *repositories.SubscriptionRepositoryImpl, userRepo *repositories.UserRepositoryImpl, mapper *mappers.SubscriptionConversionServiceImpl, logger logs.Logger) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{repo: repo, userRepo: userRepo, mapper: mapper, logger: logger}
}

// CreateSubscription handles the business logic for creating a subscription.
func (s *SubscriptionServiceImpl) CreateSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	s.logger.Info("Service Level: Attempting to create new subscription")
	subscription, err := s.mapper.DTOToSubscription(req)
	if err != nil {
		s.logger.Error("Service Level: Failed to convert DTO to subscription", err)
		return nil, err
	}

	createdSubscription, err := s.repo.Create(ctx, subscription)
	if err != nil {
		s.logger.Error("Service Level: Failed to create subscription", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to create subscription", "create_subscription_failed")
	}

	response := s.mapper.SubscriptionToDTO(createdSubscription)
	s.logger.Info("Service Level: Successfully created new subscription")
	return response, nil
}

func (s *SubscriptionServiceImpl) PurchaseSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	s.logger.Info("Attempting to purchase subscription")

	// Create subscription
	subscription, err := s.CreateSubscription(ctx, req)

	if err != nil {
		s.logger.Error("Service Level: Failed to create subscription", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to create subscription", "create_subscription_failed")
	}

	// Update user model
	_, err = s.userRepo.UpdateUser(ctx, req.UserID.Value, map[string]interface{}{"subscription_id": subscription.ID.Value})
	if err != nil {
		s.logger.Error("Service Level: Failed to update user", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to update user", "update_user_failed")
	}

	return subscription, nil
}

// UpdateSubscription handles the business logic for updating a subscription.
func (s *SubscriptionServiceImpl) UpdateSubscription(ctx context.Context, id primitive.ObjectID, req *dtos.UpdateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	s.logger.Info("Service Level: Attempting to update subscription")
	updateFields := s.mapper.UpdateSubscriptionFieldsToMap(*req)

	_, err := s.repo.Update(ctx, id, updateFields)
	if err != nil {
		s.logger.Error("Service Level: Failed to update subscription", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to update subscription", "update_subscription_failed")
	}

	updatedSubscription, err := s.GetSubscriptionByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated subscription", err)
		return nil, err
	}

	s.logger.Info("Service Level: Successfully updated subscription")
	return updatedSubscription, nil
}

// GetAllSubscriptions retrieves all subscriptions.
func (s *SubscriptionServiceImpl) GetAllSubscriptions(ctx context.Context) ([]dtos.SubscriptionResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve all subscriptions")
	subscriptions, err := s.repo.FindAll(ctx)
	if err != nil {
		s.logger.Error("Service Level: Failed to get all subscriptions", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get all subscriptions", "get_all_subscriptions_failed")
	}

	response := s.mapper.SubscriptionsToDTO(subscriptions)
	s.logger.Info("Service Level: Successfully retrieved all subscriptions")
	return response, nil
}

// GetSubscriptionByID retrieves a subscription by its ID.
func (s *SubscriptionServiceImpl) GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*dtos.SubscriptionResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve subscription by ID")
	subscription, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to get subscription", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get subscription", "get_subscription_failed")
	}

	response := s.mapper.SubscriptionToDTO(subscription)
	s.logger.Info("Service Level: Successfully retrieved subscription by ID")
	return response, nil
}

// GetSubscriptionsByPlan retrieves subscriptions by their plan.
func (s *SubscriptionServiceImpl) GetSubscriptionsByPlan(ctx context.Context, plan string) ([]dtos.SubscriptionResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve subscriptions by plan")
	subscriptions, err := s.repo.FindByPlan(ctx, plan)
	if err != nil {
		s.logger.Error("Service Level: Failed to get subscriptions by plan", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get subscriptions by plan", "get_subscriptions_by_plan_failed")
	}

	response := s.mapper.SubscriptionsToDTO(subscriptions)
	s.logger.Info("Service Level: Successfully retrieved subscriptions by plan")
	return response, nil
}

// DeleteSubscription deletes a subscription by its ID.
func (s *SubscriptionServiceImpl) DeleteSubscription(ctx context.Context, id primitive.ObjectID) error {
	s.logger.Info("Service Level: Attempting to delete subscription")
	_, err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to delete subscription", err)
		return errors.NewDatabaseError("Service Level: Failed to delete subscription", "delete_subscription_failed")
	}

	s.logger.Info("Service Level: Successfully deleted subscription")
	return nil
}
