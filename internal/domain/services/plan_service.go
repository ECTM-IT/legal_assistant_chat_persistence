package services

import (
	"context"
	"errors"
	"fmt"
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

type SubscriptionConversionService interface {
	SubscriptionToDTO(sub *models.Subscriptions) *dtos.SubscriptionResponse
}

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

	validPlanType := s.validatePlanType(planType)
	plans, err := s.GetPlansByType(ctx, validPlanType)
	if err != nil {
		s.logger.Error("Service: Failed to get plans by type", err)
		return nil, err
	}
	planDTOs := s.planMapper.PlansToDTO(plans)

	response := &dtos.PlanListResponse{
		Plans: helpers.NewNullable(planDTOs),
	}

	s.logger.Info("Service: Successfully retrieved plans")
	return response, nil
}

func (s *PlanServiceImpl) TogglePlanType(ctx context.Context, req *dtos.TogglePlanTypeRequest) (*dtos.SubscriptionResponse, error) {
	s.logger.Info("Service: Toggling plan type")

	if err := s.validateTogglePlanTypeRequest(req); err != nil {
		s.logger.Error("Service: Validation failed for TogglePlanTypeRequest", err)
		return nil, err
	}

	subscriptions, err := s.subscriptionRepo.FindByUserID(ctx, req.UserID.Value)
	if err != nil {
		s.logger.Error("Service: Failed to find subscriptions by userID", err)
		return nil, customErrors.NewDatabaseError("Failed to find subscriptions", "find_subscription_failed")
	}

	activeSubscription := s.findActiveSubscription(subscriptions)
	if activeSubscription == nil {
		errMsg := "no active subscription found"
		s.logger.Error("Service: "+errMsg, errors.New(errMsg))
		return nil, errors.New(errMsg)
	}

	update := bson.M{"type": string(req.NewType.Value)}
	updatedSubscription, err := s.subscriptionRepo.Update(ctx, activeSubscription.ID, update)
	if err != nil {
		s.logger.Error("Service: Failed to update subscription type", err)
		return nil, customErrors.NewDatabaseError("Failed to update subscription", "update_subscription_failed")
	}

	response := s.subMapper.SubscriptionToDTO(updatedSubscription)
	s.logger.Info("Service: Successfully toggled plan type")
	return response, nil
}

func (s *PlanServiceImpl) SelectPlan(ctx context.Context, req *dtos.SelectPlanRequest) (*dtos.SelectedPlanResponse, error) {
	s.logger.Info("Service: Selecting plan")

	if err := s.validateSelectPlanRequest(req); err != nil {
		s.logger.Error("Service: Validation failed for SelectPlanRequest", err)
		return nil, err
	}

	selectedPlan, err := s.GetPlanByNameAndType(ctx, string(req.Type.Value), string(req.Plan.Value))
	if err != nil {
		s.logger.Error("Service: Failed to get plan by name and type", err)
		return nil, err
	}

	subscription, err := s.getOrCreateSubscription(ctx, req.UserID.Value, selectedPlan)
	if err != nil {
		s.logger.Error("Service: Failed to get or create subscription", err)
		return nil, err
	}

	remainingDuration := time.Until(subscription.CurrentPeriodEnd)
	response := &dtos.SelectedPlanResponse{
		UserID:             helpers.NewNullable(subscription.UserID),
		Plan:               helpers.NewNullable(subscription.Plan.Name),
		Type:               helpers.NewNullable(dtos.PlanType(subscription.Plan.Type)),
		Price:              helpers.NewNullable(subscription.Plan.Price),
		RemainingTrialDays: helpers.NewNullable(int(remainingDuration.Hours() / 24)),
	}

	s.logger.Info("Service: Successfully selected plan")
	return response, nil
}

// ------------------------- Private Helper Methods --------------------------

func (s *PlanServiceImpl) validatePlanType(planType string) string {
	if planType != "monthly" && planType != "annual" {
		s.logger.Info(fmt.Sprintf("Invalid planType '%s', defaulting to 'monthly'", planType))
		return "monthly"
	}
	return planType
}

func (s *PlanServiceImpl) validateTogglePlanTypeRequest(req *dtos.TogglePlanTypeRequest) error {
	if !req.UserID.Present || !req.NewType.Present {
		return errors.New("user ID and new type are required")
	}
	return nil
}

func (s *PlanServiceImpl) validateSelectPlanRequest(req *dtos.SelectPlanRequest) error {
	if !req.UserID.Present || !req.Plan.Present || !req.Type.Present {
		return errors.New("user ID, plan, and type are required")
	}
	return nil
}

func (s *PlanServiceImpl) findActiveSubscription(subscriptions []models.Subscriptions) *models.Subscriptions {
	for _, sub := range subscriptions {
		if sub.Status == "active" {
			return &sub
		}
	}
	return nil
}

// GetAllPlans returns all predefined plans keyed by their plan type.
func (p *PlanServiceImpl) GetAllPlans() map[string][]models.Plan {
	return models.PredefinedPlans()
}

// GetPlansByType returns all plans for a given plan type.
// If the plan type does not exist, it returns an empty slice.
func (p *PlanServiceImpl) GetPlansByType(ctx context.Context, planType string) ([]models.Plan, error) {
	if plans, exists := models.PredefinedPlans()[planType]; exists {
		return plans, nil
	}
	return []models.Plan{}, nil
}

// GetPlanByNameAndType searches for a plan by its type and name.
// Returns the plan if found, otherwise returns an error.
func (p *PlanServiceImpl) GetPlanByNameAndType(ctx context.Context, planType, planName string) (*models.Plan, error) {
	plans, err := p.GetPlansByType(ctx, planType)
	if err != nil {
		return nil, err
	}
	for _, plan := range plans {
		if plan.Name == planName {
			return &plan, nil
		}
	}
	return nil, fmt.Errorf("no plan named '%s' found for type '%s'", planName, planType)
}

func (s *PlanServiceImpl) getOrCreateSubscription(ctx context.Context, userID primitive.ObjectID, selectedPlan *models.Plan) (*models.Subscriptions, error) {
	subs, err := s.subscriptionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, customErrors.NewDatabaseError("Failed to find subscription", "find_subscription_failed")
	}

	if len(subs) == 0 {
		// Create new subscription
		subscription := &models.Subscriptions{
			ID:                 primitive.NewObjectID(),
			UserID:             userID,
			Plan:               *selectedPlan,
			Status:             "active",
			CurrentPeriodStart: time.Now(),
			CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0), // default trial of 1 month
		}

		return s.subscriptionRepo.Create(ctx, subscription)
	}

	// Update existing subscription
	existingSub := &subs[0]

	return s.subscriptionRepo.Update(ctx, existingSub.ID, bson.M{
		"plan": selectedPlan,
	})
}
