// handlers/subscription_handler.go

package handlers

import (
	"context"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionService interface {
	GetAllSubscriptions(ctx context.Context) ([]dtos.SubscriptionResponse, error)
	GetSubscriptionByID(ctx context.Context, id primitive.ObjectID) (*dtos.SubscriptionResponse, error)
	GetSubscriptionsByPlan(ctx context.Context, plan string) ([]dtos.SubscriptionResponse, error)
	CreateSubscription(ctx context.Context, req *dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, id primitive.ObjectID, req *dtos.UpdateSubscriptionRequest) (*dtos.SubscriptionResponse, error)
	DeleteSubscription(ctx context.Context, id primitive.ObjectID) error
	PurchaseSubscription(ctx context.Context, userID string, planID string, planType string, paymentMethodID string, billingInformations map[string]interface{}) (*models.Subscriptions, error)
	ReactivateSubscription(ctx context.Context, id primitive.ObjectID) (*dtos.SubscriptionResponse, error)
}

type SubscriptionHandler struct {
	BaseHandler
	service *services.SubscriptionServiceImpl
}

func NewSubscriptionHandler(service *services.SubscriptionServiceImpl) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.service.GetAllSubscriptions(r.Context())
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve subscriptions")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, subscriptions)
}

func (h *SubscriptionHandler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	subscription, err := h.service.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, "Subscription not found")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, subscription)
}

func (h *SubscriptionHandler) GetSubscriptionsByPlan(w http.ResponseWriter, r *http.Request) {
	plan := r.URL.Query().Get("plan")
	subscriptions, err := h.service.GetSubscriptionsByPlan(r.Context(), plan)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve subscriptions")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, subscriptions)
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateSubscriptionRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	subscription, err := h.service.CreateSubscription(r.Context(), &req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to create subscription")
		return
	}
	h.RespondWithJSON(w, http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) PurchaseSubscription(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID              string                 `json:"user_id" binding:"required"`
		PlanID              string                 `json:"plan_id" binding:"required"`
		PlanType            string                 `json:"plan_type" binding:"required"`
		PaymentMethodID     string                 `json:"payment_method_id" binding:"required"`
		BillingInformations map[string]interface{} `json:"billing_informations"`
	}

	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	subscription, err := h.service.PurchaseSubscription(r.Context(), req.UserID, req.PlanID, req.PlanType, req.PaymentMethodID, req.BillingInformations)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to purchase subscription")
		return
	}

	h.RespondWithJSON(w, http.StatusCreated, subscription)
}

func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	var req dtos.UpdateSubscriptionRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	subscription, err := h.service.UpdateSubscription(r.Context(), id, &req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to update subscription")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, subscription)
}

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	err = h.service.DeleteSubscription(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to cancel subscription")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, "subscription_canceled")
}

func (h *SubscriptionHandler) ReactivateSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid subscription ID")
		return
	}

	subscription, err := h.service.ReactivateSubscription(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to reactivate subscription")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, subscription)
}
