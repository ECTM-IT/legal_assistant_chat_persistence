package dtos

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanResponse struct {
	Name        helpers.Nullable[string]   `json:"name"`
	Type        helpers.Nullable[string]   `json:"type"`
	Price       helpers.Nullable[float64]  `json:"price"`
	Description helpers.Nullable[string]   `json:"description"`
	Features    helpers.Nullable[[]string] `json:"features"`
}

type TogglePlanTypeRequest struct {
	UserID  helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	NewType helpers.Nullable[SubscriptionType]   `json:"new_type"`
}

type SelectPlanRequest struct {
	UserID helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	Plan   helpers.Nullable[string]             `json:"plan"`
	Type   helpers.Nullable[SubscriptionType]   `json:"type"`
}

type PlanListResponse struct {
	Plans helpers.Nullable[[]PlanResponse] `json:"plans"`
}

type SelectedPlanResponse struct {
	UserID             helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	Plan               helpers.Nullable[string]             `json:"plan"`
	Type               helpers.Nullable[SubscriptionType]   `json:"type"`
	Price              helpers.Nullable[float64]            `json:"price"`
	RemainingTrialDays helpers.Nullable[int]                `json:"remaining_trial_days"`
}
