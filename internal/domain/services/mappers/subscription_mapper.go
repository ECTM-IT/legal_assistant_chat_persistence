package mappers

import (
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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

type SubscriptionConversionServiceImpl struct{}

func NewSubscriptionConversionService() *SubscriptionConversionServiceImpl {
	return &SubscriptionConversionServiceImpl{}
}

func (s *SubscriptionConversionServiceImpl) SubscriptionToDTO(subscription *models.Subscriptions) *dtos.SubscriptionResponse {
	if subscription == nil {
		return nil
	}

	return &dtos.SubscriptionResponse{
		ID:                  helpers.NewNullable(subscription.ID),
		Plan:                helpers.NewNullable(subscription.Plan),
		Expiry:              helpers.NewNullable(subscription.Expiry),
		Type:                helpers.NewNullable(subscription.Type),
		BillingInformations: helpers.NewNullable(s.BillingInfoToDTO(subscription.BillingInformations)),
	}
}

func (s *SubscriptionConversionServiceImpl) SubscriptionsToDTO(subscriptions []models.Subscriptions) []dtos.SubscriptionResponse {
	responseList := make([]dtos.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responseList[i] = *s.SubscriptionToDTO(&subscription)
	}
	return responseList
}

func (s *SubscriptionConversionServiceImpl) DTOToSubscription(req *dtos.CreateSubscriptionRequest) (*models.Subscriptions, error) {
	if req == nil {
		return nil, errors.New("subscription request cannot be nil")
	}

	if !req.Plan.Present || !req.Expiry.Present || !req.Type.Present {
		return nil, errors.New("plan, expiry, and type are required fields")
	}

	billingInfo, err := s.DTOToBillingInfo(req.BillingInformations.OrElse(map[string]interface{}{}))
	if err != nil {
		return nil, err
	}

	return &models.Subscriptions{
		ID:                  primitive.NewObjectID(),
		Plan:                req.Plan.Value,
		Expiry:              req.Expiry.Value,
		Type:                req.Type.Value,
		BillingInformations: billingInfo,
	}, nil
}

func (s *SubscriptionConversionServiceImpl) UpdateSubscriptionFieldsToMap(updateRequest dtos.UpdateSubscriptionRequest) map[string]interface{} {
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
		billingInfo, _ := s.DTOToBillingInfo(updateRequest.BillingInformations.Value)
		updateFields["billing_informations"] = billingInfo
	}

	return updateFields
}

func (s *SubscriptionConversionServiceImpl) BillingInfoToDTO(billingInfo map[string]interface{}) map[string]interface{} {
	// This method might need more complex logic depending on your specific billing information structure
	return billingInfo
}

func (s *SubscriptionConversionServiceImpl) DTOToBillingInfo(billingInfoDTO map[string]interface{}) (map[string]interface{}, error) {
	// This method might need more complex logic depending on your specific billing information structure
	// For example, you might want to validate certain fields or convert types
	return billingInfoDTO, nil
}
