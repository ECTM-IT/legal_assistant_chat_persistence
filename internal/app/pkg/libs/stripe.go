package libs

import (
	"context"
	"fmt"
	"os"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
	"github.com/stripe/stripe-go/v74/subscription"
	"github.com/stripe/stripe-go/v74/webhook"
	"go.uber.org/zap"
)

// StripeService defines the interface for Stripe operations
type StripeService interface {
	CreateCustomer(ctx context.Context, email string) (string, error)
	CreateSubscription(ctx context.Context, customerID string, priceID string, paymentMethodID string) (string, error)
	GetProductsAndPrices(ctx context.Context) (map[string]map[string]string, error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	ReactivateSubscription(ctx context.Context, subscriptionID string) error
	VerifyWebhookSignature(payload []byte, signature string) (*stripe.Event, error)
	GetSubscriptionStatus(ctx context.Context, subscriptionID string) (string, error)
}

// StripeServiceImpl implements the StripeService interface
type StripeServiceImpl struct {
	logger logs.Logger
}

// NewStripeService creates a new instance of the Stripe service
func NewStripeService(logger logs.Logger) StripeService {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	return &StripeServiceImpl{
		logger: logger,
	}
}

// VerifyWebhookSignature verifies the signature of a Stripe webhook event
func (s *StripeServiceImpl) VerifyWebhookSignature(payload []byte, signature string) (*stripe.Event, error) {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		s.logger.Error("Stripe webhook secret not found in environment variables", fmt.Errorf("missing STRIPE_WEBHOOK_SECRET"))
		return nil, fmt.Errorf("stripe webhook secret not configured")
	}

	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		webhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		s.logger.Error("Failed to verify webhook signature", err)
		return nil, fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	return &event, nil
}

// GetProductsAndPrices fetches all products and their prices from Stripe
func (s *StripeServiceImpl) GetProductsAndPrices(ctx context.Context) (map[string]map[string]string, error) {
	s.logger.Debug("Fetching products and prices from Stripe")

	// Initialize the result map
	// First key is plan type (monthly/annual), second key is plan name (pro/team/enterprise)
	result := make(map[string]map[string]string)
	result["monthly"] = make(map[string]string)
	result["annual"] = make(map[string]string)

	// Fetch all products
	productParams := &stripe.ProductListParams{}
	productParams.Filters.AddFilter("active", "", "true")

	products := product.List(productParams)
	for products.Next() {
		p := products.Product()

		// Fetch all prices for this product
		priceParams := &stripe.PriceListParams{
			Product: stripe.String(p.ID),
			Active:  stripe.Bool(true),
		}

		prices := price.List(priceParams)
		for prices.Next() {
			pr := prices.Price()

			// Skip if not a recurring price
			if pr.Recurring == nil {
				continue
			}

			// Determine plan type based on interval
			planType := "monthly"
			if pr.Recurring.Interval == "year" {
				planType = "annual"
			}

			// Map the product name to its price ID
			// Assuming product names in Stripe match our plan names (pro, team, enterprise)
			result[planType][p.Name] = pr.ID
			s.logger.Debug("Mapped product to price",
				zap.String("product", p.Name),
				zap.String("planType", planType),
				zap.String("priceID", pr.ID))
		}

		if err := prices.Err(); err != nil {
			s.logger.Error("Failed to fetch prices for product", err,
				zap.String("productID", p.ID),
				zap.String("productName", p.Name))
			continue
		}
	}

	if err := products.Err(); err != nil {
		s.logger.Error("Failed to fetch products", err)
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	// Log the result for debugging
	for planType, plans := range result {
		for planName, priceID := range plans {
			s.logger.Debug("Product mapping",
				zap.String("planType", planType),
				zap.String("planName", planName),
				zap.String("priceID", priceID))
		}
	}

	return result, nil
}

// CreateCustomer creates a new customer in Stripe or returns an existing one with the same email
func (s *StripeServiceImpl) CreateCustomer(ctx context.Context, email string) (string, error) {
	// First, try to find if a customer with this email already exists
	s.logger.Debug("Searching for existing customer with email", zap.String("email", email))

	params := &stripe.CustomerListParams{
		Email: stripe.String(email),
	}
	params.Filters.AddFilter("limit", "", "1")

	customers := customer.List(params)
	for customers.Next() {
		// Found a customer with this email
		existingCustomer := customers.Customer()
		s.logger.Info("Found existing Stripe customer",
			zap.String("id", existingCustomer.ID),
			zap.String("email", email))
		return existingCustomer.ID, nil
	}

	if err := customers.Err(); err != nil {
		s.logger.Error("Error searching for existing customer", err)
		// Continue with creating a new customer despite the search error
	}

	// No existing customer found, create a new one
	s.logger.Info("No existing customer found, creating new Stripe customer", zap.String("email", email))
	newCustomerParams := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	newCustomer, err := customer.New(newCustomerParams)
	if err != nil {
		s.logger.Error("Failed to create Stripe customer", err)
		return "", fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	return newCustomer.ID, nil
}

// CreateSubscription creates a new subscription in Stripe
func (s *StripeServiceImpl) CreateSubscription(ctx context.Context, customerID string, priceID string, paymentMethodID string) (string, error) {
	// First, attach the payment method to the customer
	attachParams := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}

	_, err := paymentmethod.Attach(paymentMethodID, attachParams)
	if err != nil {
		s.logger.Error("Failed to attach payment method", err)
		return "", fmt.Errorf("failed to attach payment method: %w", err)
	}

	// Set the payment method as the default for the customer
	customerParams := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentMethodID),
		},
	}

	_, err = customer.Update(customerID, customerParams)
	if err != nil {
		s.logger.Error("Failed to set default payment method", err)
		return "", fmt.Errorf("failed to set default payment method: %w", err)
	}

	// Create the subscription
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			PaymentMethodTypes: []*string{
				stripe.String("card"),
			},
		},
	}

	subscription, err := subscription.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe subscription", err)
		return "", fmt.Errorf("failed to create Stripe subscription: %w", err)
	}

	return subscription.ID, nil
}

// CancelSubscription cancels a subscription in Stripe
func (s *StripeServiceImpl) CancelSubscription(ctx context.Context, subscriptionID string) error {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := subscription.Update(subscriptionID, params)
	if err != nil {
		s.logger.Error("Failed to cancel Stripe subscription", err)
		return fmt.Errorf("failed to cancel Stripe subscription: %w", err)
	}

	return nil
}

// ReactivateSubscription reactivates a previously canceled subscription in Stripe
func (s *StripeServiceImpl) ReactivateSubscription(ctx context.Context, subscriptionID string) error {
	// Retrieve the subscription to ensure it exists
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		s.logger.Error("Failed to retrieve Stripe subscription", err)
		return fmt.Errorf("failed to retrieve Stripe subscription: %w", err)
	}

	// If the subscription was already canceled (not just pending cancellation),
	// we can't reactivate it and need to create a new one
	if sub.Status == "canceled" {
		s.logger.Error("Cannot reactivate a fully canceled subscription", nil)
		return fmt.Errorf("cannot reactivate a fully canceled subscription")
	}

	// If subscription is set to cancel at period end, update it to remove cancellation
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}

	_, err = subscription.Update(subscriptionID, params)
	if err != nil {
		s.logger.Error("Failed to reactivate Stripe subscription", err)
		return fmt.Errorf("failed to reactivate Stripe subscription: %w", err)
	}

	return nil
}

// GetSubscriptionStatus gets the current status of a subscription from Stripe
func (s *StripeServiceImpl) GetSubscriptionStatus(ctx context.Context, subscriptionID string) (string, error) {
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		s.logger.Error("Failed to retrieve Stripe subscription", err)
		return "", fmt.Errorf("failed to retrieve Stripe subscription: %w", err)
	}

	return string(sub.Status), nil
}
