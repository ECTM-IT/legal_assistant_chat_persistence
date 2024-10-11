package libs

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/sub"
	"go.uber.org/zap"
)

// StripeService defines the interface for Stripe operations.
type StripeService interface {
	CreateCustomer(ctx context.Context, email string) (*stripe.Customer, error)
	CreateSubscription(ctx context.Context, customerID, priceID string) (*stripe.Subscription, error)
	GetCustomer(ctx context.Context, customerID string) (*stripe.Customer, error)
	GetSubscription(ctx context.Context, subscriptionID string) (*stripe.Subscription, error)
}

// stripeService is the concrete implementation of StripeService.
type stripeService struct {
	logger *zap.Logger
}

// NewStripeService creates a new instance of StripeService.
func NewStripeService(logger *zap.Logger) StripeService {
	stripe.Key = helpers.GetStripeAPIKey() // Ensure you have a helper to retrieve the API key securely.
	return &stripeService{
		logger: logger,
	}
}

// CreateCustomer creates a new Stripe customer.
func (s *stripeService) CreateCustomer(ctx context.Context, email string) (*stripe.Customer, error) {
	s.logger.Info("Creating new Stripe customer", zap.String("email", email))

	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	// Avoid naming conflict with the 'customer' package
	cust, err := customer.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe customer", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Stripe customer created successfully", zap.String("customer_id", cust.ID))
	return cust, nil
}

// CreateSubscription creates a new subscription for a customer.
func (s *stripeService) CreateSubscription(ctx context.Context, customerID, priceID string) (*stripe.Subscription, error) {
	s.logger.Info("Creating new Stripe subscription", zap.String("customer_id", customerID), zap.String("price_id", priceID))

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
	}

	// Retry logic for transient errors (up to 3 retries)
	var subscriptionResult *stripe.Subscription
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		subscriptionResult, err = sub.New(params)
		if err == nil {
			break
		}
		s.logger.Warn("Retrying subscription creation due to error", zap.Int("attempt", attempt+1), zap.Error(err))
		if attempt < 2 {
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		s.logger.Error("Failed to create Stripe subscription", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Stripe subscription created successfully", zap.String("subscription_id", subscriptionResult.ID))
	return subscriptionResult, nil
}

// GetCustomer retrieves a Stripe customer by ID.
func (s *stripeService) GetCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	s.logger.Info("Retrieving Stripe customer", zap.String("customer_id", customerID))

	cust, err := customer.Get(customerID, nil)
	if err != nil {
		s.logger.Error("Failed to retrieve Stripe customer", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Stripe customer retrieved successfully", zap.String("customer_id", cust.ID))
	return cust, nil
}

// GetSubscription retrieves a Stripe subscription by ID.
func (s *stripeService) GetSubscription(ctx context.Context, subscriptionID string) (*stripe.Subscription, error) {
	s.logger.Info("Retrieving Stripe subscription", zap.String("subscription_id", subscriptionID))

	subscriptionResult, err := sub.Get(subscriptionID, nil)
	if err != nil {
		s.logger.Error("Failed to retrieve Stripe subscription", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Stripe subscription retrieved successfully", zap.String("subscription_id", subscriptionResult.ID))
	return subscriptionResult, nil
}
