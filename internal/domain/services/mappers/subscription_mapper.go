package mappers

import (
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionConversionService interface {
	SubscriptionToDTO(subscription *models.Subscriptions) *dtos.SubscriptionResponse
	SubscriptionsToDTO(subscriptions []models.Subscriptions) []dtos.SubscriptionResponse
	DTOToSubscription(req *dtos.CreateSubscriptionRequest) (*models.Subscriptions, error)
	UpdateSubscriptionFieldsToMap(updateRequest dtos.UpdateSubscriptionRequest) map[string]interface{}
	BillingInfoToDTO(billingInfo map[string]interface{}) map[string]interface{}
	DTOToBillingInfo(billingInfoDTO map[string]interface{}) (map[string]interface{}, error)
}

type SubscriptionConversionServiceImpl struct {
	logger logs.Logger
}

func NewSubscriptionConversionService(logger logs.Logger) *SubscriptionConversionServiceImpl {
	return &SubscriptionConversionServiceImpl{
		logger: logger,
	}
}

func (s *SubscriptionConversionServiceImpl) SubscriptionToDTO(subscription *models.Subscriptions) *dtos.SubscriptionResponse {
	s.logger.Info("Converting Subscription to DTO")
	if subscription == nil {
		s.logger.Warn("Attempted to convert nil Subscription to DTO")
		return nil
	}

	dto := &dtos.SubscriptionResponse{
		ID:                  helpers.NewNullable(subscription.ID),
		Plan:                helpers.NewNullable(subscription.Plan),
		Expiry:              helpers.NewNullable(subscription.Expiry),
		Type:                helpers.NewNullable(subscription.Type),
		BillingInformations: helpers.NewNullable(s.BillingInfoToDTO(subscription.BillingInformations)),
	}
	s.logger.Info("Successfully converted Subscription to DTO")
	return dto
}

func (s *SubscriptionConversionServiceImpl) SubscriptionsToDTO(subscriptions []models.Subscriptions) []dtos.SubscriptionResponse {
	s.logger.Info("Converting multiple Subscriptions to DTOs")
	responseList := make([]dtos.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responseList[i] = *s.SubscriptionToDTO(&subscription)
	}
	s.logger.Info("Successfully converted multiple Subscriptions to DTOs")
	return responseList
}

func (s *SubscriptionConversionServiceImpl) DTOToSubscription(req *dtos.CreateSubscriptionRequest) (*models.Subscriptions, error) {
	s.logger.Info("Converting DTO to Subscription")
	if req == nil {
		s.logger.Error("Failed to convert DTO to Subscription: subscription request cannot be nil", errors.New("subscription request cannot be nil"))
		return nil, errors.New("subscription request cannot be nil")
	}

	if !req.Plan.Present || !req.Expiry.Present || !req.Type.Present {
		s.logger.Error("Failed to convert DTO to Subscription: plan, expiry, and type are required fields", errors.New("plan, expiry, and type are required fields"))
		return nil, errors.New("plan, expiry, and type are required fields")
	}

	billingInfo, err := s.DTOToBillingInfo(req.BillingInformations.OrElse(map[string]interface{}{}))
	if err != nil {
		s.logger.Error("Failed to convert billing information", err)
		return nil, err
	}

	subscription := &models.Subscriptions{
		ID:                  primitive.NewObjectID(),
		Plan:                req.Plan.Value,
		Expiry:              req.Expiry.Value,
		Type:                req.Type.Value,
		BillingInformations: billingInfo,
	}
	s.logger.Info("Successfully converted DTO to Subscription")
	return subscription, nil
}

func (s *SubscriptionConversionServiceImpl) UpdateSubscriptionFieldsToMap(updateRequest dtos.UpdateSubscriptionRequest) map[string]interface{} {
	s.logger.Info("Converting UpdateSubscriptionRequest to map")
	updateFields := make(map[string]interface{})

	if updateRequest.Plan.Present {
		updateFields["plan"] = updateRequest.Plan.Value
	}
	if updateRequest.Expiry.Present {
		updateFields["expiry"] = updateRequest.Expiry.Value
	}
	if updateRequest.Type.Present {
		updateFields["type"] = updateRequest.Type.Value
	}
	if updateRequest.BillingInformations.Present {
		billingInfo, err := s.DTOToBillingInfo(updateRequest.BillingInformations.Value)
		if err != nil {
			s.logger.Error("Failed to convert billing information for update", err)
		} else {
			updateFields["billing_informations"] = billingInfo
		}
	}

	s.logger.Info("Successfully converted UpdateSubscriptionRequest to map")
	return updateFields
}

func (s *SubscriptionConversionServiceImpl) BillingInfoToDTO(billingInfo map[string]interface{}) map[string]interface{} {
	s.logger.Info("Converting BillingInfo to DTO")
	// This method might need more complex logic depending on your specific billing information structure
	s.logger.Info("Successfully converted BillingInfo to DTO")
	return billingInfo
}

func (s *SubscriptionConversionServiceImpl) DTOToBillingInfo(billingInfoDTO map[string]interface{}) (map[string]interface{}, error) {
	s.logger.Info("Converting DTO to BillingInfo")
	// This method might need more complex logic depending on your specific billing information structure
	// For example, you might want to validate certain fields or convert types
	s.logger.Info("Successfully converted DTO to BillingInfo")
	return billingInfoDTO, nil
}
