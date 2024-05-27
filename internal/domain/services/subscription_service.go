package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	repo *repositories.SubscriptionRepositoryImpl
}

// NewSubscriptionService creates a new instance of the subscription service.
func NewSubscriptionService(repo *repositories.SubscriptionRepositoryImpl) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{repo: repo}
}

// CreateSubscription handles the business logic for creating a subscription.
func (s *SubscriptionServiceImpl) CreateSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	// Transform DTO to entity
	subscription := &models.Subscriptions{
		ID:                  primitive.NewObjectID(),
		Plan:                req.Plan.Val,
		Expiry:              req.Expiry.Val,
		Type:                req.Type.Val,
		BillingInformations: req.BillingInformations.Val,
	}
	// Call repository to save the entity
	_, err := s.repo.Create(ctx, subscription)
	if err != nil {
		return nil, err
	}
	return s.toSubscriptionResponse(subscription), nil
}

// UpdateSubscription handles the business logic for updating a subscription.
func (s *SubscriptionServiceImpl) UpdateSubscription(ctx context.Context, id primitive.ObjectID, req *dtos.UpdateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	// Transform DTO to entity updates
	update := bson.M{
		"plan":                 req.Plan.Val,
		"expiry":               req.Expiry.Val,
		"type":                 req.Type.Val,
		"billing_informations": req.BillingInformations.Val,
	}
	// Call repository to update the entity
	_, err := s.repo.Update(ctx, id, update)
	if err != nil {
		return nil, err
	}
	// Fetch updated subscription to return
	updatedSubscription, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toSubscriptionResponse(updatedSubscription), nil
}

// GetAllSubscriptions retrieves all subscriptions.
func (s *SubscriptionServiceImpl) GetAllSubscriptions(ctx context.Context) ([]dtos.SubscriptionResponse, error) {
	subscriptions, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return s.toSubscriptionResponseList(subscriptions), nil
}

// GetSubscriptionByID retrieves a subscription by its ID.
func (s *SubscriptionServiceImpl) GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*dtos.SubscriptionResponse, error) {
	subscription, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toSubscriptionResponse(subscription), nil
}

// GetSubscriptionsByPlan retrieves subscriptions by their plan.
func (s *SubscriptionServiceImpl) GetSubscriptionsByPlan(ctx context.Context, plan string) ([]dtos.SubscriptionResponse, error) {
	subscriptions, err := s.repo.FindByPlan(ctx, plan)
	if err != nil {
		return nil, err
	}
	return s.toSubscriptionResponseList(subscriptions), nil
}

// DeleteSubscription deletes a subscription by its ID.
func (s *SubscriptionServiceImpl) DeleteSubscription(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return s.repo.Delete(ctx, id)
}

// toSubscriptionResponse converts a Subscriptions model to a SubscriptionResponse DTO.
func (s *SubscriptionServiceImpl) toSubscriptionResponse(subscription *models.Subscriptions) *dtos.SubscriptionResponse {
	return &dtos.SubscriptionResponse{
		ID:                  helpers.NewNullable(subscription.ID),
		Plan:                helpers.NewNullable(subscription.Plan),
		Expiry:              helpers.NewNullable(subscription.Expiry),
		Type:                helpers.NewNullable(subscription.Type),
		BillingInformations: helpers.NewNullable(subscription.BillingInformations),
	}
}

// toSubscriptionResponseList converts a list of Subscriptions models to a list of SubscriptionResponse DTOs.
func (s *SubscriptionServiceImpl) toSubscriptionResponseList(subscriptions []models.Subscriptions) []dtos.SubscriptionResponse {
	responseList := make([]dtos.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responseList[i] = *s.toSubscriptionResponse(&subscription)
	}
	return responseList
}
