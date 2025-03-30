package services

import (
	"context"
	"errors"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/libs"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// SubscriptionService defines the subscription service interface.
type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, id primitive.ObjectID, req *dtos.UpdateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
	GetAllSubscriptions(ctx context.Context) ([]dtos.SubscriptionResponse, error)
	GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*dtos.SubscriptionResponse, error)
	GetSubscriptionsByPlan(ctx context.Context, plan string) ([]dtos.SubscriptionResponse, error)
	DeleteSubscription(ctx context.Context, id primitive.ObjectID) error
	PurchaseSubscription(ctx context.Context, userID string, planID string, planType string, paymentMethodID string) (*models.Subscriptions, error)
}

// SubscriptionServiceImpl implements the SubscriptionService interface.
type SubscriptionServiceImpl struct {
	repo          *repositories.SubscriptionRepositoryImpl
	userRepo      *repositories.UserRepositoryImpl
	mapper        *mappers.SubscriptionConversionServiceImpl
	planService   *PlanServiceImpl
	stripeService libs.StripeService
	logger        logs.Logger
}

// NewSubscriptionService creates a new instance of the subscription service.
func NewSubscriptionService(
	repo *repositories.SubscriptionRepositoryImpl,
	userRepo *repositories.UserRepositoryImpl,
	mapper *mappers.SubscriptionConversionServiceImpl,
	planService *PlanServiceImpl,
	stripeService libs.StripeService,
	logger logs.Logger,
) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{
		repo:          repo,
		userRepo:      userRepo,
		mapper:        mapper,
		planService:   planService,
		stripeService: stripeService,
		logger:        logger,
	}
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
		return nil, errors.New("failed to create subscription")
	}

	response := s.mapper.SubscriptionToDTO(createdSubscription)
	s.logger.Info("Service Level: Successfully created new subscription")
	return response, nil
}

func (s *SubscriptionServiceImpl) PurchaseSubscription(ctx context.Context, userID string, planID string, planType string, paymentMethodID string) (*models.Subscriptions, error) {
	s.logger.Info("Attempting to purchase subscription",
		zap.String("userID", userID),
		zap.String("planID", planID),
		zap.String("planType", planType))

	// Convert userID string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error("Invalid user ID", err)
		return nil, errors.New("invalid user ID")
	}

	// Get user details
	user, err := s.userRepo.FindByID(ctx, userObjectID)
	if err != nil {
		s.logger.Error("Failed to get user", err)
		return nil, errors.New("failed to get user")
	}

	// Get the Stripe price ID from the predefined plans
	predefinedPlans := models.PredefinedPlans()
	planTypePlans, exists := predefinedPlans[planType]
	if !exists {
		s.logger.Error("Invalid plan type", errors.New("invalid plan type"))
		return nil, errors.New("invalid plan type")
	}

	var selectedPlan models.Plan
	found := false
	for _, p := range planTypePlans {
		if p.Name == planID {
			selectedPlan = p
			found = true
			break
		}
	}

	if !found {
		s.logger.Error("Plan not found", errors.New("plan not found"))
		return nil, errors.New("plan not found")
	}

	// Fetch products and prices from Stripe
	stripeProducts, err := s.stripeService.GetProductsAndPrices(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch Stripe products", err)
		return nil, errors.New("failed to fetch Stripe products")
	}

	// Get the Stripe price ID for the selected plan
	stripePlanPrices, exists := stripeProducts[planType]
	if !exists {
		s.logger.Error("Plan type not found in Stripe", errors.New("plan type not found in Stripe"))
		return nil, errors.New("plan type not found in Stripe")
	}

	stripePriceID, exists := stripePlanPrices[planID]
	if !exists {
		s.logger.Error("Plan not found in Stripe", errors.New("plan not found in Stripe"))
		return nil, errors.New("plan not found in Stripe")
	}

	// Create Stripe customer
	stripeCustomerID, err := s.stripeService.CreateCustomer(ctx, user.Email)
	if err != nil {
		s.logger.Error("Failed to create Stripe customer", err)
		return nil, errors.New("failed to create Stripe customer")
	}

	// Create Stripe subscription using the fetched price ID and payment method
	stripeSubscriptionID, err := s.stripeService.CreateSubscription(ctx, stripeCustomerID, stripePriceID, paymentMethodID)
	if err != nil {
		s.logger.Error("Failed to create Stripe subscription", err)
		return nil, errors.New("failed to create Stripe subscription")
	}

	// Create subscription in database
	subscription := &models.Subscriptions{
		UserID:               userObjectID,
		Plan:                 selectedPlan,
		Expiry:               time.Now().AddDate(0, 1, 0), // 1 month from now
		Status:               "active",
		StripeCustomerID:     stripeCustomerID,
		StripeSubscriptionID: stripeSubscriptionID,
		CurrentPeriodStart:   time.Now(),
		CurrentPeriodEnd:     time.Now().AddDate(0, 1, 0),
		CancelAtPeriodEnd:    false,
		BillingInformations:  make(map[string]interface{}),
	}

	createdSubscription, err := s.repo.Create(ctx, subscription)
	if err != nil {
		s.logger.Error("Failed to create subscription", err)
		return nil, errors.New("failed to create subscription")
	}

	return createdSubscription, nil
}

// UpdateSubscription handles the business logic for updating a subscription.
func (s *SubscriptionServiceImpl) UpdateSubscription(ctx context.Context, id primitive.ObjectID, req *dtos.UpdateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	s.logger.Info("Service Level: Attempting to update subscription")
	updateFields := s.mapper.UpdateSubscriptionFieldsToMap(*req)

	_, err := s.repo.Update(ctx, id, updateFields)
	if err != nil {
		s.logger.Error("Service Level: Failed to update subscription", err)
		return nil, errors.New("failed to update subscription")
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
		return nil, errors.New("failed to get all subscriptions")
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
		return nil, errors.New("failed to get subscription")
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
		return nil, errors.New("failed to get subscriptions by plan")
	}

	response := s.mapper.SubscriptionsToDTO(subscriptions)
	s.logger.Info("Service Level: Successfully retrieved subscriptions by plan")
	return response, nil
}

// DeleteSubscription deletes a subscription by its ID.
func (s *SubscriptionServiceImpl) DeleteSubscription(ctx context.Context, id primitive.ObjectID) error {
	s.logger.Info("Service Level: Attempting to delete subscription")
	_, err := s.userRepo.UpdateUser(ctx, id, map[string]interface{}{"subscription_id": primitive.NilObjectID})
	if err != nil {
		s.logger.Error("Service Level: Failed to delete user", err)
		return errors.New("failed to delete user")
	}
	_, err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to delete subscription", err)
		return errors.New("failed to delete subscription")
	}

	s.logger.Info("Service Level: Successfully deleted subscription")
	return nil
}
