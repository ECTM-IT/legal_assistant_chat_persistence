package services

import (
	"context"
	"errors"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	customErrors "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanService interface {
	GetPlans(ctx context.Context, planType string) (*dtos.PlanListResponse, error)
	TogglePlanType(ctx context.Context, req *dtos.TogglePlanTypeRequest) (*dtos.SubscriptionResponse, error)
	SelectPlan(ctx context.Context, req *dtos.SelectPlanRequest) (*dtos.SelectedPlanResponse, error)
	CreateSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
	GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*dtos.SubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, req *dtos.UpdateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
	DeleteSubscription(ctx context.Context, id primitive.ObjectID) error
}

type PlanServiceImpl struct {
	subscriptionRepo *repositories.SubscriptionRepositoryImpl
	planMapper       *mappers.PlanConversionServiceImpl
	subMapper        *mappers.SubscriptionConversionServiceImpl
	logger           logs.Logger
}

func NewPlanService(
	subscriptionRepo *repositories.SubscriptionRepositoryImpl,
	planMapper *mappers.PlanConversionServiceImpl,
	subMapper *mappers.SubscriptionConversionServiceImpl,
	logger logs.Logger,
) *PlanServiceImpl {
	return &PlanServiceImpl{
		subscriptionRepo: subscriptionRepo,
		planMapper:       planMapper,
		subMapper:        subMapper,
		logger:           logger,
	}
}

func (s *PlanServiceImpl) GetPlans(ctx context.Context, planType string) (*dtos.PlanListResponse, error) {
	s.logger.Info("Service: Retrieving plans")

	if planType != "monthly" && planType != "annual" {
		planType = "monthly" // Default to monthly if invalid type
	}

	plans := models.PredefinedPlans()[planType]
	planDTOs := s.planMapper.PlansToDTO(plans)

	response := &dtos.PlanListResponse{
		Plans: helpers.NewNullable(planDTOs),
	}

	s.logger.Info("Service: Successfully retrieved plans")
	return response, nil
}

func (s *PlanServiceImpl) TogglePlanType(ctx context.Context, req *dtos.TogglePlanTypeRequest) (*dtos.SubscriptionResponse, error) {
	s.logger.Info("Service: Toggling plan type")

	if !req.UserID.Present || !req.NewType.Present {
		s.logger.Error("Service: Missing required user ID or new type", errors.New("user ID and new type are required"))
		return nil, errors.New("user ID and new type are required")
	}

	// Get current subscription(s)
	subscriptions, err := s.subscriptionRepo.FindByUserID(ctx, req.UserID.Value)
	if err != nil {
		s.logger.Error("Service: Failed to find subscription", err)
		return nil, customErrors.NewDatabaseError("Failed to find subscription", "find_subscription_failed")
	}

	if len(subscriptions) == 0 {
		s.logger.Error("Service: No subscription found for user", errors.New("no subscription found for user"))
		return nil, errors.New("no subscription found for user")
	}

	// Get the active subscription
	var activeSubscription *models.Subscriptions
	for _, sub := range subscriptions {
		if sub.Status == "active" {
			activeSubscription = &sub
			break
		}
	}

	if activeSubscription == nil {
		s.logger.Error("Service: No active subscription found for user", errors.New("no active subscription found"))
		return nil, errors.New("no active subscription found")
	}

	update := bson.M{
		"type": string(req.NewType.Value),
	}

	updatedSubscription, err := s.subscriptionRepo.Update(ctx, activeSubscription.ID, update)
	if err != nil {
		s.logger.Error("Service: Failed to update subscription", err)
		return nil, customErrors.NewDatabaseError("Failed to update subscription", "update_subscription_failed")
	}

	response := s.subMapper.SubscriptionToDTO(updatedSubscription)
	s.logger.Info("Service: Successfully toggled plan type")
	return response, nil
}

func (s *PlanServiceImpl) SelectPlan(ctx context.Context, req *dtos.SelectPlanRequest) (*dtos.SelectedPlanResponse, error) {
	s.logger.Info("Service: Selecting plan")

	if !req.UserID.Present || !req.Plan.Present || !req.Type.Present {
		s.logger.Error("Service: Missing required user ID, plan, or type", errors.New("user ID, plan, and type are required"))
		return nil, errors.New("user ID, plan, and type are required")
	}

	// Validate plan exists
	plans := models.PredefinedPlans()[string(req.Type.Value)]
	var selectedPlan *models.Plan
	for _, p := range plans {
		if p.Name == req.Plan.Value {
			selectedPlan = &p
			break
		}
	}

	if selectedPlan == nil {
		s.logger.Error("Service: Invalid plan selection", errors.New("invalid plan selection"))
		return nil, errors.New("invalid plan selection")
	}

	// Get or create subscription
	subscriptions, err := s.subscriptionRepo.FindByUserID(ctx, req.UserID.Value)
	if err != nil {
		s.logger.Error("Service: Failed to find subscription", err)
		return nil, customErrors.NewDatabaseError("Failed to find subscription", "find_subscription_failed")
	}

	var subscription *models.Subscriptions
	if len(subscriptions) == 0 {
		// Create new subscription
		subscription = &models.Subscriptions{
			ID:                 primitive.NewObjectID(),
			UserID:             req.UserID.Value,
			Plan:               *selectedPlan,
			Status:             "active",
			CurrentPeriodStart: time.Now(),
			CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0), // default 1 month trial
		}

		subscription, err = s.subscriptionRepo.Create(ctx, subscription)
		if err != nil {
			s.logger.Error("Service: Failed to create subscription", err)
			return nil, customErrors.NewDatabaseError("Failed to create subscription", "create_subscription_failed")
		}
	} else {
		// Update existing subscription
		subscription = &subscriptions[0]
		update := bson.M{
			"plan": selectedPlan.Name,
			"type": selectedPlan.Type,
		}

		subscription, err = s.subscriptionRepo.Update(ctx, subscription.ID, update)
		if err != nil {
			s.logger.Error("Service: Failed to update subscription", err)
			return nil, customErrors.NewDatabaseError("Failed to update subscription", "update_subscription_failed")
		}
	}

	remainingDuration := time.Until(subscription.CurrentPeriodEnd)
	response := &dtos.SelectedPlanResponse{
		UserID:             helpers.NewNullable(subscription.UserID),
		Plan:               helpers.NewNullable(subscription.Plan.Name),
		Type:               helpers.NewNullable(dtos.PlanType(subscription.Plan.Type)),
		Price:              helpers.NewNullable(selectedPlan.Price),
		RemainingTrialDays: helpers.NewNullable(int(remainingDuration.Hours() / 24)),
	}

	s.logger.Info("Service: Successfully selected plan")
	return response, nil
}
