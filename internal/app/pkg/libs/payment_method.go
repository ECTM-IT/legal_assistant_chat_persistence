package libs

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	"go.uber.org/zap"
)

// UpdateStripeService extends StripeService with payment method operations
type UpdatedStripeService interface {
	StripeService
	AttachPaymentMethod(ctx context.Context, paymentMethodID, customerID string) error
	DetachPaymentMethod(ctx context.Context, paymentMethodID string) error
	SetDefaultPaymentMethod(ctx context.Context, customerID, paymentMethodID string) error
	GetPaymentMethod(ctx context.Context, paymentMethodID string) (*stripe.PaymentMethod, error)
	ListPaymentMethods(ctx context.Context, customerID string) ([]*stripe.PaymentMethod, error)
}

// AttachPaymentMethod attaches a payment method to a customer
func (s *StripeServiceImpl) AttachPaymentMethod(ctx context.Context, paymentMethodID, customerID string) error {
	s.logger.Debug("Attaching payment method to customer",
		zap.String("paymentMethodID", paymentMethodID),
		zap.String("customerID", customerID))

	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}

	_, err := paymentmethod.Attach(paymentMethodID, params)
	if err != nil {
		s.logger.Error("Failed to attach payment method", err)
		return fmt.Errorf("failed to attach payment method: %w", err)
	}

	return nil
}

// DetachPaymentMethod detaches a payment method from a customer
func (s *StripeServiceImpl) DetachPaymentMethod(ctx context.Context, paymentMethodID string) error {
	s.logger.Debug("Detaching payment method", zap.String("paymentMethodID", paymentMethodID))

	_, err := paymentmethod.Detach(paymentMethodID, nil)
	if err != nil {
		s.logger.Error("Failed to detach payment method", err)
		return fmt.Errorf("failed to detach payment method: %w", err)
	}

	return nil
}

// SetDefaultPaymentMethod sets a payment method as the default for a customer
func (s *StripeServiceImpl) SetDefaultPaymentMethod(ctx context.Context, customerID, paymentMethodID string) error {
	s.logger.Debug("Setting default payment method",
		zap.String("customerID", customerID),
		zap.String("paymentMethodID", paymentMethodID))

	params := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentMethodID),
		},
	}

	_, err := customer.Update(customerID, params)
	if err != nil {
		s.logger.Error("Failed to set default payment method", err)
		return fmt.Errorf("failed to set default payment method: %w", err)
	}

	return nil
}

// GetPaymentMethod retrieves a payment method by ID
func (s *StripeServiceImpl) GetPaymentMethod(ctx context.Context, paymentMethodID string) (*stripe.PaymentMethod, error) {
	s.logger.Debug("Retrieving payment method", zap.String("paymentMethodID", paymentMethodID))

	pm, err := paymentmethod.Get(paymentMethodID, nil)
	if err != nil {
		s.logger.Error("Failed to retrieve payment method", err)
		return nil, fmt.Errorf("failed to retrieve payment method: %w", err)
	}

	return pm, nil
}

// ListPaymentMethods lists all payment methods for a customer
func (s *StripeServiceImpl) ListPaymentMethods(ctx context.Context, customerID string) ([]*stripe.PaymentMethod, error) {
	s.logger.Debug("Listing payment methods for customer", zap.String("customerID", customerID))

	var paymentMethods []*stripe.PaymentMethod

	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String("card"),
	}

	i := paymentmethod.List(params)
	for i.Next() {
		pm := i.PaymentMethod()
		paymentMethods = append(paymentMethods, pm)
	}

	if err := i.Err(); err != nil {
		s.logger.Error("Failed to list payment methods", err)
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}

	return paymentMethods, nil
}
