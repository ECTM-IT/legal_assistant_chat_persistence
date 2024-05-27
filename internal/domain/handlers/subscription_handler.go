package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubscriptionHandler handles HTTP requests for subscription management.
type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionServiceImpl
}

// NewSubscriptionHandler creates a new instance of SubscriptionHandler.
func NewSubscriptionHandler(subscriptionService *services.SubscriptionServiceImpl) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// GetAllSubscriptions handles the retrieval of all subscriptions.
func (h *SubscriptionHandler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subscriptions, err := h.subscriptionService.GetAllSubscriptions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

// GetSubscriptionByID handles the retrieval of a subscription by its ID.
func (h *SubscriptionHandler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	subscription, err := h.subscriptionService.GetSubscriptionByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// GetSubscriptionsByPlan handles the retrieval of subscriptions by their plan.
func (h *SubscriptionHandler) GetSubscriptionsByPlan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	plan := strings.TrimSpace(r.URL.Query().Get("plan"))
	subscriptions, err := h.subscriptionService.GetSubscriptionsByPlan(ctx, plan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

// CreateSubscription handles the creation of a new subscription.
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dtos.CreateSubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subscription, err := h.subscriptionService.CreateSubscription(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// UpdateSubscription handles the update of an existing subscription.
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req dtos.UpdateSubscriptionRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subscription, err := h.subscriptionService.UpdateSubscription(ctx, id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// DeleteSubscription handles the deletion of a subscription.
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	deletedResult, err := h.subscriptionService.DeleteSubscription(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedResult)
}
