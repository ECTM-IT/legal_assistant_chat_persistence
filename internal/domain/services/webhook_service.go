package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/libs"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/stripe/stripe-go/v74"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// WebhookService defines the interface for Stripe webhook event processing
type WebhookService interface {
	ProcessWebhook(ctx context.Context, payload []byte, signature string) error
}

// WebhookServiceImpl implements the WebhookService interface
type WebhookServiceImpl struct {
	stripeService    libs.StripeService
	stripeEventRepo  *repositories.StripeEventRepositoryImpl
	subscriptionRepo *repositories.SubscriptionRepositoryImpl
	userRepo         *repositories.UserRepositoryImpl
	logger           logs.Logger
}

// NewWebhookService creates a new instance of the webhook service
func NewWebhookService(
	stripeService libs.StripeService,
	stripeEventRepo *repositories.StripeEventRepositoryImpl,
	subscriptionRepo *repositories.SubscriptionRepositoryImpl,
	userRepo *repositories.UserRepositoryImpl,
	logger logs.Logger,
) *WebhookServiceImpl {
	return &WebhookServiceImpl{
		stripeService:    stripeService,
		stripeEventRepo:  stripeEventRepo,
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
		logger:           logger,
	}
}

// ProcessWebhook processes a Stripe webhook event
func (s *WebhookServiceImpl) ProcessWebhook(ctx context.Context, payload []byte, signature string) error {
	// Verify webhook signature
	event, err := s.stripeService.VerifyWebhookSignature(payload, signature)
	if err != nil {
		s.logger.Error("Failed to verify webhook signature", err)
		return err
	}

	// Check if event has already been processed (idempotency)
	existingEvent, err := s.stripeEventRepo.FindByEventID(ctx, event.ID)
	if err != nil {
		s.logger.Error("Error checking for existing event", err)
		return err
	}

	if existingEvent != nil {
		s.logger.Info("Event already processed, skipping", zap.String("event_id", event.ID))
		return nil
	}

	// Store the event in the database
	eventData := make(map[string]interface{})
	if err := json.Unmarshal(payload, &eventData); err != nil {
		s.logger.Error("Failed to unmarshal event data", err)
		return err
	}

	stripeEvent := &models.StripeEvent{
		EventID:   event.ID,
		Type:      event.Type,
		CreatedAt: time.Now(),
		Data:      eventData,
		Processed: false,
	}

	_, err = s.stripeEventRepo.Create(ctx, stripeEvent)
	if err != nil {
		s.logger.Error("Failed to store event", err)
		return err
	}

	// Process the event based on its type
	err = s.handleEvent(ctx, event)
	if err != nil {
		s.logger.Error("Failed to handle event", err, zap.String("event_type", event.Type))
		return err
	}

	// Mark the event as processed
	_, err = s.stripeEventRepo.Update(ctx, stripeEvent.ID, bson.M{
		"processed": true,
	})
	if err != nil {
		s.logger.Error("Failed to mark event as processed", err)
		return err
	}

	return nil
}

// handleEvent processes the event based on its type
func (s *WebhookServiceImpl) handleEvent(ctx context.Context, event *stripe.Event) error {
	switch event.Type {
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	case "invoice.paid":
		return s.handleInvoicePaid(ctx, event)
	case "invoice.payment_failed":
		return s.handleInvoicePaymentFailed(ctx, event)
	default:
		s.logger.Info("Unhandled event type", zap.String("type", event.Type))
		return nil
	}
}

// handleSubscriptionCreated handles a subscription.created event
func (s *WebhookServiceImpl) handleSubscriptionCreated(ctx context.Context, event *stripe.Event) error {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		return err
	}

	s.logger.Info("Processing subscription created event",
		zap.String("subscription_id", subscription.ID),
		zap.String("customer_id", subscription.Customer.ID))

	// This event might be redundant if we already created the subscription in our database
	// during the PurchaseSubscription flow, but we should check here as well

	// Find the user with this customer ID
	user, err := s.findUserByStripeCustomerID(ctx, subscription.Customer.ID)
	if err != nil {
		return err
	}

	if user == nil {
		s.logger.Warn("No user found with Stripe customer ID", zap.String("customer_id", subscription.Customer.ID))
		return errors.New("no user found with provided Stripe customer ID")
	}

	// Check if a subscription already exists for this user
	subs, err := s.subscriptionRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	// If there are subscriptions, check if any have this Stripe subscription ID
	found := false
	for _, sub := range subs {
		if sub.StripeSubscriptionID == subscription.ID {
			found = true
			break
		}
	}

	if found {
		// Already have this subscription in our database
		return nil
	}

	// If not found, create a new subscription record
	// Note: We don't have the plan details here, so we'd need to update the subscription
	// with that information from another source.
	// For now, we'll just log this case
	s.logger.Warn("Found subscription in Stripe that's not in our database",
		zap.String("subscription_id", subscription.ID),
		zap.String("user_id", user.ID.Hex()))

	return nil
}

// handleSubscriptionUpdated handles a subscription.updated event
func (s *WebhookServiceImpl) handleSubscriptionUpdated(ctx context.Context, event *stripe.Event) error {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		return err
	}

	s.logger.Info("Processing subscription updated event",
		zap.String("subscription_id", subscription.ID),
		zap.String("status", string(subscription.Status)),
		zap.Bool("cancel_at_period_end", subscription.CancelAtPeriodEnd))

	// Find the subscription in our database
	subs, err := s.findSubscriptionsByStripeID(ctx, subscription.ID)
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		s.logger.Warn("No subscription found with Stripe ID", zap.String("subscription_id", subscription.ID))
		return nil
	}

	// Update the subscription status
	sub := subs[0]

	updates := bson.M{}

	// Update subscription status
	statusMap := map[stripe.SubscriptionStatus]string{
		stripe.SubscriptionStatusActive:            "active",
		stripe.SubscriptionStatusCanceled:          "canceled",
		stripe.SubscriptionStatusPastDue:           "past_due",
		stripe.SubscriptionStatusUnpaid:            "unpaid",
		stripe.SubscriptionStatusIncomplete:        "incomplete",
		stripe.SubscriptionStatusIncompleteExpired: "incomplete_expired",
		stripe.SubscriptionStatusTrialing:          "trialing",
	}

	if status, ok := statusMap[subscription.Status]; ok {
		updates["status"] = status
	}

	// Important: If cancel_at_period_end is true, we should mark the subscription as canceled in our database
	// This is the key fix - Stripe keeps status as "active" but we want to show "canceled" to users
	if subscription.CancelAtPeriodEnd {
		updates["cancel_at_period_end"] = true

		// Override the status if this is a cancellation - critical fix!
		updates["status"] = "canceled"

		// Set canceled_at if not already set
		if sub.CanceledAt.IsZero() {
			updates["canceled_at"] = time.Now()
		}
	} else {
		updates["cancel_at_period_end"] = false
	}

	// Update current period
	if subscription.CurrentPeriodStart > 0 {
		updates["current_period_start"] = time.Unix(subscription.CurrentPeriodStart, 0)
	}
	if subscription.CurrentPeriodEnd > 0 {
		updates["current_period_end"] = time.Unix(subscription.CurrentPeriodEnd, 0)
		// Also update expiry to match current period end
		updates["expiry"] = time.Unix(subscription.CurrentPeriodEnd, 0)
	}

	// If canceled, update canceled_at
	if subscription.Status == stripe.SubscriptionStatusCanceled && subscription.CanceledAt > 0 {
		updates["canceled_at"] = time.Unix(subscription.CanceledAt, 0)
	}

	// Only update if we have changes to make
	if len(updates) > 0 {
		// Apply updates
		_, err = s.subscriptionRepo.Update(ctx, sub.ID, updates)
		if err != nil {
			s.logger.Error("Failed to update subscription", err)
			return err
		}

		// Log the changes we made
		updateFields := make([]string, 0, len(updates))
		for key := range updates {
			updateFields = append(updateFields, key)
		}
		s.logger.Info("Updated subscription from webhook event",
			zap.String("subscription_id", sub.ID.Hex()),
			zap.Strings("updated_fields", updateFields))
	} else {
		s.logger.Info("No updates needed for subscription",
			zap.String("subscription_id", sub.ID.Hex()))
	}

	return nil
}

// handleSubscriptionDeleted handles a subscription.deleted event
func (s *WebhookServiceImpl) handleSubscriptionDeleted(ctx context.Context, event *stripe.Event) error {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		return err
	}

	s.logger.Info("Processing subscription deleted event", zap.String("subscription_id", subscription.ID))

	// Find the subscription in our database
	subs, err := s.findSubscriptionsByStripeID(ctx, subscription.ID)
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		s.logger.Warn("No subscription found with Stripe ID", zap.String("subscription_id", subscription.ID))
		return nil
	}

	// Update the subscription to mark as canceled
	sub := subs[0]

	updates := bson.M{
		"status":      "canceled",
		"canceled_at": time.Now(),
	}

	// Apply updates
	_, err = s.subscriptionRepo.Update(ctx, sub.ID, updates)
	if err != nil {
		s.logger.Error("Failed to update subscription", err)
		return err
	}

	return nil
}

// handleInvoicePaid handles an invoice.paid event
func (s *WebhookServiceImpl) handleInvoicePaid(ctx context.Context, event *stripe.Event) error {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		return err
	}

	// Only process subscription invoices
	if invoice.Subscription == nil {
		return nil
	}

	s.logger.Info("Processing invoice paid event",
		zap.String("invoice_id", invoice.ID),
		zap.String("subscription_id", invoice.Subscription.ID))

	// Find the subscription in our database
	subs, err := s.findSubscriptionsByStripeID(ctx, invoice.Subscription.ID)
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		s.logger.Warn("No subscription found with Stripe ID", zap.String("subscription_id", invoice.Subscription.ID))
		return nil
	}

	// Update the subscription to ensure it's active
	sub := subs[0]

	// If the subscription is past_due, mark it as active again
	if sub.Status == "past_due" || sub.Status == "unpaid" {
		updates := bson.M{
			"status": "active",
		}

		// Apply updates
		_, err = s.subscriptionRepo.Update(ctx, sub.ID, updates)
		if err != nil {
			s.logger.Error("Failed to update subscription", err)
			return err
		}
	}

	return nil
}

// handleInvoicePaymentFailed handles an invoice.payment_failed event
func (s *WebhookServiceImpl) handleInvoicePaymentFailed(ctx context.Context, event *stripe.Event) error {
	var invoice stripe.Invoice
	err := json.Unmarshal(event.Data.Raw, &invoice)
	if err != nil {
		return err
	}

	// Only process subscription invoices
	if invoice.Subscription == nil {
		return nil
	}

	s.logger.Info("Processing invoice payment failed event",
		zap.String("invoice_id", invoice.ID),
		zap.String("subscription_id", invoice.Subscription.ID))

	// Find the subscription in our database
	subs, err := s.findSubscriptionsByStripeID(ctx, invoice.Subscription.ID)
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		s.logger.Warn("No subscription found with Stripe ID", zap.String("subscription_id", invoice.Subscription.ID))
		return nil
	}

	// Update the subscription to mark as past_due
	sub := subs[0]

	updates := bson.M{
		"status": "past_due",
	}

	// Apply updates
	_, err = s.subscriptionRepo.Update(ctx, sub.ID, updates)
	if err != nil {
		s.logger.Error("Failed to update subscription", err)
		return err
	}

	return nil
}

// findUserByStripeCustomerID finds a user by their Stripe customer ID
func (s *WebhookServiceImpl) findUserByStripeCustomerID(ctx context.Context, customerID string) (*models.User, error) {
	// Get all users (this is inefficient but works for a small number of users)
	allUsers, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range allUsers {
		if user.StripeCustomerID == customerID {
			return &user, nil
		}
	}

	return nil, nil
}

// findSubscriptionsByStripeID finds subscriptions by their Stripe subscription ID
func (s *WebhookServiceImpl) findSubscriptionsByStripeID(ctx context.Context, subscriptionID string) ([]models.Subscriptions, error) {
	// This is assuming you have a method to find by Stripe subscription ID
	// If not, we'd need to search through all subscriptions

	// Get all subscriptions (this is inefficient but works for a small number of subscriptions)
	allSubs, err := s.subscriptionRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []models.Subscriptions
	for _, sub := range allSubs {
		if sub.StripeSubscriptionID == subscriptionID {
			result = append(result, sub)
		}
	}

	return result, nil
}
