package libs

import (
	"context"
	"errors"
	"os"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/subscription"
	"go.uber.org/zap"
)

type StripeService interface {
	CreateCustomer(ctx context.Context, email string) (string, error)
	CreateSubscription(ctx context.Context, customerID, priceID string) (string, error)
}

type StripeServiceImpl struct {
	logger logs.Logger
}

func NewStripeService(logger logs.Logger) StripeService {
	// Initialize Stripe API key
	stripeKey := os.Getenv("STRIPE_API_KEY")
	if stripeKey == "" {
		logger.Error("Stripe API key not found in environment variables", errors.New("missing STRIPE_API_KEY"))
		return nil
	}

	// Log the first few characters of the key for debugging (safely)
	if len(stripeKey) > 4 {
		logger.Info("Initializing Stripe service with API key", zap.String("key_prefix", stripeKey[:4]+"..."))
	} else {
		logger.Info("Initializing Stripe service with API key", zap.String("key_prefix", stripeKey))
	}

	// Set the API key
	stripe.Key = stripeKey

	// Create and return the service
	service := &StripeServiceImpl{
		logger: logger,
	}

	// Test the connection by making a simple API call
	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("limit", "", "1")
	iter := customer.List(params)
	if iter.Err() != nil {
		logger.Error("Failed to connect to Stripe API", iter.Err())
		return nil
	}

	logger.Info("Successfully initialized Stripe service")
	return service
}

func (s *StripeServiceImpl) CreateCustomer(ctx context.Context, email string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	customer, err := customer.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe customer", err)
		return "", errors.New("failed to create Stripe customer")
	}

	return customer.ID, nil
}

func (s *StripeServiceImpl) CreateSubscription(ctx context.Context, customerID, priceID string) (string, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
	}

	sub, err := subscription.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe subscription", err)
		return "", errors.New("failed to create Stripe subscription")
	}

	return sub.ID, nil
}
