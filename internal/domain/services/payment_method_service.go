package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/libs"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// PaymentMethodService defines the interface for payment method operations
type PaymentMethodService interface {
	AddPaymentMethod(ctx context.Context, request *dtos.PaymentMethodRequest) (*dtos.PaymentMethodResponse, error)
	GetPaymentMethods(ctx context.Context, userID string) ([]dtos.PaymentMethodResponse, error)
	SetDefaultPaymentMethod(ctx context.Context, userID, paymentMethodID string) (*dtos.PaymentMethodResponse, error)
}

// PaymentMethodServiceImpl implements the PaymentMethodService interface
type PaymentMethodServiceImpl struct {
	repo          repositories.PaymentMethodRepository
	userRepo      repositories.UserRepository
	mapper        mappers.PaymentMethodConversionService
	stripeService libs.UpdatedStripeService
	logger        logs.Logger
}

// NewPaymentMethodService creates a new instance of PaymentMethodServiceImpl
func NewPaymentMethodService(
	repo repositories.PaymentMethodRepository,
	userRepo repositories.UserRepository,
	mapper mappers.PaymentMethodConversionService,
	stripeService libs.UpdatedStripeService,
	logger logs.Logger,
) *PaymentMethodServiceImpl {
	return &PaymentMethodServiceImpl{
		repo:          repo,
		userRepo:      userRepo,
		mapper:        mapper,
		stripeService: stripeService,
		logger:        logger,
	}
}

// AddPaymentMethod adds a new payment method to the user's account
func (s *PaymentMethodServiceImpl) AddPaymentMethod(ctx context.Context, request *dtos.PaymentMethodRequest) (*dtos.PaymentMethodResponse, error) {
	s.logger.Info(fmt.Sprintf("Adding payment method for user %s", request.UserID.Value))

	// Convert userID string to ObjectID
	userID, err := primitive.ObjectIDFromHex(request.UserID.Value)
	if err != nil {
		s.logger.Error("Invalid user ID format", err)
		return nil, errors.New("invalid user ID format")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to find user", err)
		return nil, errors.New("user not found")
	}

	// Retrieve payment method from Stripe
	pm, err := s.stripeService.GetPaymentMethod(ctx, request.PaymentMethodID.Value)
	if err != nil {
		s.logger.Error("Failed to retrieve payment method from Stripe", err)
		return nil, errors.New("invalid payment method")
	}

	// Attach payment method to customer
	err = s.stripeService.AttachPaymentMethod(ctx, request.PaymentMethodID.Value, user.StripeCustomerID)
	if err != nil {
		s.logger.Error("Failed to attach payment method to customer", err)
		return nil, errors.New("failed to attach payment method")
	}

	// Convert Stripe payment method to map
	pmMap := map[string]interface{}{
		"id":       pm.ID,
		"type":     pm.Type,
		"customer": user.StripeCustomerID,
	}

	if pm.Card != nil {
		pmMap["card"] = map[string]interface{}{
			"brand":     string(pm.Card.Brand),
			"last4":     pm.Card.Last4,
			"exp_month": float64(pm.Card.ExpMonth),
			"exp_year":  float64(pm.Card.ExpYear),
		}
	}

	// Set as default payment method if requested
	if request.IsDefault.Present && request.IsDefault.Value {
		err = s.stripeService.SetDefaultPaymentMethod(ctx, user.StripeCustomerID, pm.ID)
		if err != nil {
			s.logger.Warn("Failed to set default payment method in Stripe", zap.Error(err))
			// Continue with saving the payment method despite the error
		}
	}

	// Convert to domain model
	paymentMethod, err := s.mapper.DTOToPaymentMethod(request, pmMap)
	if err != nil {
		s.logger.Error("Failed to convert payment method DTO to model", err)
		return nil, errors.New("invalid payment method data")
	}

	// If this is set as default, unset any existing default payment methods
	if paymentMethod.IsDefault {
		err = s.unsetExistingDefaultPaymentMethods(ctx, paymentMethod.UserID)
		if err != nil {
			s.logger.Warn("Failed to unset existing default payment methods", zap.Error(err))
			// Continue with saving the new payment method despite the error
		}
	}

	// Save payment method to database
	newPaymentMethod, err := s.repo.Create(ctx, paymentMethod)
	if err != nil {
		s.logger.Error("Failed to save payment method to database", err)
		return nil, errors.New("failed to save payment method")
	}

	s.logger.Info(fmt.Sprintf("Successfully added payment method for user %s", request.UserID.Value))
	return s.mapper.PaymentMethodToDTO(newPaymentMethod), nil
}

// GetPaymentMethods retrieves all payment methods for a user
func (s *PaymentMethodServiceImpl) GetPaymentMethods(ctx context.Context, userID string) ([]dtos.PaymentMethodResponse, error) {
	s.logger.Info(fmt.Sprintf("Retrieving payment methods for user %s", userID))

	// Convert userID string to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error("Invalid user ID format", err)
		return nil, errors.New("invalid user ID format")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userObjID)
	if err != nil {
		s.logger.Error("Failed to find user", err)
		return nil, errors.New("user not found")
	}

	// Retrieve payment methods from database
	paymentMethods, err := s.repo.FindByUserID(ctx, userObjID)
	if err != nil {
		s.logger.Error("Failed to retrieve payment methods from database", err)
		return nil, errors.New("failed to retrieve payment methods")
	}

	// If no payment methods in database, try fetching from Stripe
	if len(paymentMethods) == 0 {
		s.logger.Info("No payment methods found in database, fetching from Stripe")

		// Fetch payment methods from Stripe
		stripePMs, err := s.stripeService.ListPaymentMethods(ctx, user.StripeCustomerID)
		if err != nil {
			s.logger.Warn("Failed to fetch payment methods from Stripe", zap.Error(err))
			// Continue with empty result despite Stripe API error
		}

		// Process each Stripe payment method
		for _, pm := range stripePMs {
			// Convert to a map for processing
			pmMap := map[string]interface{}{
				"id":       pm.ID,
				"type":     pm.Type,
				"customer": user.StripeCustomerID,
			}

			if pm.Card != nil {
				pmMap["card"] = map[string]interface{}{
					"brand":     string(pm.Card.Brand),
					"last4":     pm.Card.Last4,
					"exp_month": float64(pm.Card.ExpMonth),
					"exp_year":  float64(pm.Card.ExpYear),
				}
			}

			// Create a minimal request DTO
			request := &dtos.PaymentMethodRequest{
				UserID:          helpers.Nullable[string]{Value: userID, Present: true},
				PaymentMethodID: helpers.Nullable[string]{Value: pm.ID, Present: true},
				IsDefault:       helpers.Nullable[bool]{Value: false, Present: true},
			}

			// Convert to domain model
			paymentMethod, err := s.mapper.DTOToPaymentMethod(request, pmMap)
			if err != nil {
				s.logger.Warn("Failed to convert Stripe payment method to model", zap.Error(err))
				continue
			}

			// Save to database
			_, err = s.repo.Create(ctx, paymentMethod)
			if err != nil {
				s.logger.Warn("Failed to save Stripe payment method to database", zap.Error(err))
				continue
			}

			paymentMethods = append(paymentMethods, *paymentMethod)
		}
	}

	s.logger.Info(fmt.Sprintf("Successfully retrieved %d payment methods for user %s", len(paymentMethods), userID))
	return s.mapper.PaymentMethodsToDTO(paymentMethods), nil
}

// SetDefaultPaymentMethod sets a payment method as the default for a user
func (s *PaymentMethodServiceImpl) SetDefaultPaymentMethod(ctx context.Context, userID, paymentMethodID string) (*dtos.PaymentMethodResponse, error) {
	s.logger.Info(fmt.Sprintf("Setting payment method %s as default for user %s", paymentMethodID, userID))

	// Convert userID string to ObjectID
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error("Invalid user ID format", err)
		return nil, errors.New("invalid user ID format")
	}

	// Convert paymentMethodID string to ObjectID
	paymentMethodObjID, err := primitive.ObjectIDFromHex(paymentMethodID)
	if err != nil {
		s.logger.Error("Invalid payment method ID format", err)
		return nil, errors.New("invalid payment method ID format")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userObjID)
	if err != nil {
		s.logger.Error("Failed to find user", err)
		return nil, errors.New("user not found")
	}

	// Find payment method in database
	paymentMethod, err := s.repo.FindByID(ctx, paymentMethodObjID)
	if err != nil {
		s.logger.Error("Failed to find payment method", err)
		return nil, errors.New("payment method not found")
	}

	// Verify that the payment method belongs to the user
	if paymentMethod.UserID != userObjID {
		s.logger.Error("Payment method does not belong to user", nil)
		return nil, errors.New("payment method does not belong to user")
	}

	// Update default payment method in Stripe
	err = s.stripeService.SetDefaultPaymentMethod(ctx, user.StripeCustomerID, paymentMethod.StripePaymentID)
	if err != nil {
		s.logger.Error("Failed to set default payment method in Stripe", err)
		return nil, errors.New("failed to set default payment method in Stripe")
	}

	// Unset any existing default payment methods
	err = s.unsetExistingDefaultPaymentMethods(ctx, userObjID)
	if err != nil {
		s.logger.Warn("Failed to unset existing default payment methods", zap.Error(err))
		// Continue with setting the new default despite the error
	}

	// Set payment method as default in database
	paymentMethod.IsDefault = true
	updates := bson.M{"is_default": true}
	updatedPaymentMethod, err := s.repo.Update(ctx, paymentMethodObjID, updates)
	if err != nil {
		s.logger.Error("Failed to update payment method in database", err)
		return nil, errors.New("failed to update payment method")
	}

	s.logger.Info(fmt.Sprintf("Successfully set payment method %s as default for user %s", paymentMethodID, userID))
	return s.mapper.PaymentMethodToDTO(updatedPaymentMethod), nil
}

// unsetExistingDefaultPaymentMethods unsets any existing default payment methods for a user
func (s *PaymentMethodServiceImpl) unsetExistingDefaultPaymentMethods(ctx context.Context, userID primitive.ObjectID) error {
	// Find all payment methods for the user
	paymentMethods, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Unset the default flag for each payment method
	for _, pm := range paymentMethods {
		if pm.IsDefault {
			updates := bson.M{"is_default": false}
			_, err := s.repo.Update(ctx, pm.ID, updates)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
