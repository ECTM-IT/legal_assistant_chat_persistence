package mappers

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentMethodConversionService defines the interface for payment method conversion operations
type PaymentMethodConversionService interface {
	PaymentMethodToDTO(paymentMethod *models.PaymentMethod) *dtos.PaymentMethodResponse
	PaymentMethodsToDTO(paymentMethods []models.PaymentMethod) []dtos.PaymentMethodResponse
	DTOToPaymentMethod(request *dtos.PaymentMethodRequest, stripePaymentDetails map[string]interface{}) (*models.PaymentMethod, error)
}

// PaymentMethodConversionServiceImpl implements the PaymentMethodConversionService interface
type PaymentMethodConversionServiceImpl struct {
	logger logs.Logger
}

// NewPaymentMethodConversionService creates a new instance of PaymentMethodConversionServiceImpl
func NewPaymentMethodConversionService(logger logs.Logger) *PaymentMethodConversionServiceImpl {
	return &PaymentMethodConversionServiceImpl{
		logger: logger,
	}
}

// PaymentMethodToDTO converts a payment method model to a DTO
func (s *PaymentMethodConversionServiceImpl) PaymentMethodToDTO(paymentMethod *models.PaymentMethod) *dtos.PaymentMethodResponse {
	s.logger.Info("Converting PaymentMethod to DTO")
	if paymentMethod == nil {
		s.logger.Warn("Attempted to convert nil PaymentMethod to DTO")
		return nil
	}

	return &dtos.PaymentMethodResponse{
		ID:              helpers.NewNullable(paymentMethod.ID),
		UserID:          helpers.NewNullable(paymentMethod.UserID),
		Type:            helpers.NewNullable(paymentMethod.Type),
		CardBrand:       helpers.NewNullable(paymentMethod.CardBrand),
		LastFour:        helpers.NewNullable(paymentMethod.LastFour),
		ExpirationMonth: helpers.NewNullable(paymentMethod.ExpirationMonth),
		ExpirationYear:  helpers.NewNullable(paymentMethod.ExpirationYear),
		BillingName:     helpers.NewNullable(paymentMethod.BillingName),
		BillingEmail:    helpers.NewNullable(paymentMethod.BillingEmail),
		IsDefault:       helpers.NewNullable(paymentMethod.IsDefault),
		CreatedAt:       helpers.NewNullable(paymentMethod.CreatedAt),
	}
}

// PaymentMethodsToDTO converts a slice of payment method models to DTOs
func (s *PaymentMethodConversionServiceImpl) PaymentMethodsToDTO(paymentMethods []models.PaymentMethod) []dtos.PaymentMethodResponse {
	s.logger.Info("Converting multiple PaymentMethods to DTOs")
	responseList := make([]dtos.PaymentMethodResponse, len(paymentMethods))
	for i, paymentMethod := range paymentMethods {
		responseList[i] = *s.PaymentMethodToDTO(&paymentMethod)
	}
	s.logger.Info("Successfully converted multiple PaymentMethods to DTOs")
	return responseList
}

// DTOToPaymentMethod converts a payment method request DTO to a model
func (s *PaymentMethodConversionServiceImpl) DTOToPaymentMethod(request *dtos.PaymentMethodRequest, stripePaymentDetails map[string]interface{}) (*models.PaymentMethod, error) {
	s.logger.Info("Converting PaymentMethod DTO to model")

	// Extract user ID from request
	userID, err := primitive.ObjectIDFromHex(request.UserID.Value)
	if err != nil {
		s.logger.Error("Invalid user ID in payment method request", err)
		return nil, err
	}

	// Set billing details
	billingName := ""
	billingEmail := ""
	if request.BillingDetails.Present {
		if request.BillingDetails.Value.Name.Present {
			billingName = request.BillingDetails.Value.Name.Value
		}
		if request.BillingDetails.Value.Email.Present {
			billingEmail = request.BillingDetails.Value.Email.Value
		}
	}

	// Extract card details from Stripe payment method
	paymentType := getStringValue(stripePaymentDetails, "type", "card")
	cardBrand := ""
	lastFour := ""
	expirationMonth := 0
	expirationYear := 0

	if cardDetails, ok := stripePaymentDetails["card"].(map[string]interface{}); ok {
		cardBrand = getStringValue(cardDetails, "brand", "")
		lastFour = getStringValue(cardDetails, "last4", "")

		if month, ok := cardDetails["exp_month"].(float64); ok {
			expirationMonth = int(month)
		}

		if year, ok := cardDetails["exp_year"].(float64); ok {
			expirationYear = int(year)
		}
	}

	// Create payment method model
	paymentMethod := &models.PaymentMethod{
		ID:               primitive.NewObjectID(),
		UserID:           userID,
		StripeCustomerID: getStringValue(stripePaymentDetails, "customer", ""),
		StripePaymentID:  request.PaymentMethodID.Value,
		Type:             paymentType,
		CardBrand:        cardBrand,
		LastFour:         lastFour,
		ExpirationMonth:  expirationMonth,
		ExpirationYear:   expirationYear,
		BillingName:      billingName,
		BillingEmail:     billingEmail,
		IsDefault:        request.IsDefault.Present && request.IsDefault.Value,
		CreatedAt:        time.Now().Unix(),
	}

	s.logger.Info("Successfully converted PaymentMethod DTO to model")
	return paymentMethod, nil
}

// Helper function to get string values from a map
func getStringValue(data map[string]interface{}, key, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return defaultValue
}
