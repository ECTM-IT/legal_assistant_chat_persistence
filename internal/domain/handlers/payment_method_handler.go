package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/gorilla/mux"
)

// PaymentMethodHandler handles HTTP requests for payment methods
type PaymentMethodHandler struct {
	paymentMethodService services.PaymentMethodService
	logger               logs.Logger
}

// NewPaymentMethodHandler creates a new instance of PaymentMethodHandler
func NewPaymentMethodHandler(paymentMethodService services.PaymentMethodService, logger logs.Logger) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		paymentMethodService: paymentMethodService,
		logger:               logger,
	}
}

// AddPaymentMethod handles the POST /payment-methods endpoint
func (h *PaymentMethodHandler) AddPaymentMethod(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling request to add payment method")

	// Decode request body
	var request dtos.PaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode request body", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if !request.UserID.Present || !request.PaymentMethodID.Present {
		h.logger.Error("Missing required fields", nil)
		http.Error(w, "Missing required fields: userID or paymentMethodID", http.StatusBadRequest)
		return
	}

	// Call service to add payment method
	response, err := h.paymentMethodService.AddPaymentMethod(r.Context(), &request)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to add payment method: %v", err), err)
		http.Error(w, fmt.Sprintf("Failed to add payment method: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetPaymentMethods handles the GET /users/{userId}/payment-methods endpoint
func (h *PaymentMethodHandler) GetPaymentMethods(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	vars := mux.Vars(r)
	userID := vars["userId"]
	if userID == "" {
		h.logger.Error("Missing user ID in request path", nil)
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	h.logger.Info(fmt.Sprintf("Handling request to get payment methods for user %s", userID))

	// Call service to get payment methods
	response, err := h.paymentMethodService.GetPaymentMethods(r.Context(), userID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to get payment methods: %v", err), err)
		http.Error(w, fmt.Sprintf("Failed to get payment methods: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetDefaultPaymentMethod handles the POST /users/{userId}/payment-methods/{paymentMethodId}/set-default endpoint
func (h *PaymentMethodHandler) SetDefaultPaymentMethod(w http.ResponseWriter, r *http.Request) {
	// Extract user ID and payment method ID from URL path
	vars := mux.Vars(r)
	userID := vars["userId"]
	paymentMethodID := vars["paymentMethodId"]

	if userID == "" || paymentMethodID == "" {
		h.logger.Error("Missing required parameters in request path", nil)
		http.Error(w, "Missing required parameters: userId or paymentMethodId", http.StatusBadRequest)
		return
	}

	h.logger.Info(fmt.Sprintf("Handling request to set payment method %s as default for user %s", paymentMethodID, userID))

	// Call service to set default payment method
	response, err := h.paymentMethodService.SetDefaultPaymentMethod(r.Context(), userID, paymentMethodID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to set default payment method: %v", err), err)
		http.Error(w, fmt.Sprintf("Failed to set default payment method: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
