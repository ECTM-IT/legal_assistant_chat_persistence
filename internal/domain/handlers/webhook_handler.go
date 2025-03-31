package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/stripe/stripe-go/v74/webhook"
	"go.uber.org/zap"
)

// WebhookService defines the interface for the webhook service
type WebhookService interface {
	ProcessWebhook(ctx context.Context, payload []byte, signature string) error
}

// WebhookHandler handles webhook events from Stripe
type WebhookHandler struct {
	BaseHandler
	service       services.WebhookServiceImpl
	logger        logs.Logger
	webhookSecret string
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(service *services.WebhookServiceImpl, logger logs.Logger, webhookSecret string) *WebhookHandler {
	return &WebhookHandler{
		service:       *service,
		logger:        logger,
		webhookSecret: webhookSecret,
	}
}

// HandleStripeWebhook processes incoming webhook events from Stripe
func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read webhook request body", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Get the signature header
	signatureHeader := r.Header.Get("Stripe-Signature")

	// Construct event with options to ignore API version mismatch
	event, err := webhook.ConstructEventWithOptions(
		body,
		signatureHeader,
		h.webhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		h.logger.Error("Failed to construct event", err)
		http.Error(w, "Failed to construct event", http.StatusBadRequest)
		return
	}

	h.logger.Info("Received Stripe webhook event",
		zap.String("type", event.Type),
		zap.String("signature", signatureHeader[:10]+"..."))

	// Process the webhook event asynchronously
	go func() {
		ctx := context.Background()
		err := h.service.ProcessWebhook(ctx, body, signatureHeader)
		if err != nil {
			h.logger.Error("Failed to process webhook event", err)
		}
	}()

	// Respond immediately with a 200 OK
	h.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "received"})
}
