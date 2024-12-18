package handlers

import (
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
)

type PlanHandler struct {
	BaseHandler
	service *services.PlanServiceImpl
}

func NewPlanHandler(service *services.PlanServiceImpl) *PlanHandler {
	return &PlanHandler{service: service}
}

// GetPlanOptions handles the retrieval of plan options
func (h *PlanHandler) GetPlanOptions(w http.ResponseWriter, r *http.Request) {
	// Default to monthly if not specified
	planType := r.URL.Query().Get("type")
	if planType == "" {
		planType = "monthly"
	}

	plans, err := h.service.GetPlans(r.Context(), planType)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch plans")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, plans)
}

// TogglePlanType handles switching between monthly and annual plans
func (h *PlanHandler) TogglePlanType(w http.ResponseWriter, r *http.Request) {
	var req dtos.TogglePlanTypeRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	subscription, err := h.service.TogglePlanType(r.Context(), &req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to toggle plan type")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, subscription)
}

// SelectPlan handles the plan selection process
func (h *PlanHandler) SelectPlan(w http.ResponseWriter, r *http.Request) {
	var req dtos.SelectPlanRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	selectedPlan, err := h.service.SelectPlan(r.Context(), &req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to select plan")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, selectedPlan)
}
