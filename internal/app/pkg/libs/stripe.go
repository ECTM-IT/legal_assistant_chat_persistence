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
)

// StripeService defines the interface for Stripe operations
type StripeService interface {
	CreateCustomer(ctx context.Context, email string) (string, error)
	CreateSubscription(ctx context.Context, customerID string, priceID string, paymentMethodID string) (string, error)
	GetProductsAndPrices(ctx context.Context) (map[string]map[string]string, error)
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

// GetProductsAndPrices fetches all products and their prices from Stripe
func (s *StripeServiceImpl) GetProductsAndPrices(ctx context.Context) (map[string]map[string]string, error) {
	s.logger.Debug("Fetching products and prices from Stripe")

	// Initialize the result map
	// First key is plan type (monthly/annual), second key is plan name (pro/team/enterprise)
	result := make(map[string]map[string]string)

	// Fetch all products
	params := &stripe.ProductListParams{}
	params.Filters.AddFilter("active", "", "true")
	params.Filters.AddFilter("expand[]", "", "data.default_price")

	products := product.List(params)
	for products.Next() {
		p := products.Product()

		// Get the default price for this product
		if p.DefaultPrice == nil {
			s.logger.Warn("Product has no default price")
			continue
		}

		// Determine plan type from price
		price, err := price.Get(p.DefaultPrice.ID, nil)
		if err != nil {
			s.logger.Error("Failed to get price details", err)
			continue
		}

		// Determine plan type based on interval
		planType := "monthly"
		if price.Recurring != nil && price.Recurring.Interval == "year" {
			planType = "annual"
		}

		// Initialize the inner map if it doesn't exist
		if _, exists := result[planType]; !exists {
			result[planType] = make(map[string]string)
		}

		// Map the product name to its price ID
		// Assuming product names in Stripe match our plan names (pro, team, enterprise)
		result[planType][p.Name] = price.ID
		s.logger.Debug("Mapped product to price")
	}

	if err := products.Err(); err != nil {
		s.logger.Error("Failed to fetch products", err)
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	return result, nil
}

// CreateCustomer creates a new customer in Stripe
func (s *StripeServiceImpl) CreateCustomer(ctx context.Context, email string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	customer, err := customer.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe customer", err)
		return "", fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	return customer.ID, nil
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
