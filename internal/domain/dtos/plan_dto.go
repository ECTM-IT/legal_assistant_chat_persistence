package dtos

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanResponse struct {
	Name          helpers.Nullable[string]   `json:"name" bson:"name"`
	Type          helpers.Nullable[string]   `json:"type" bson:"type"`
	Price         helpers.Nullable[float64]  `json:"price" bson:"price"`
	Description   helpers.Nullable[string]   `json:"description" bson:"description"`
	Features      helpers.Nullable[[]string] `json:"features" bson:"features"`
	StripePriceID helpers.Nullable[string]   `json:"stripe_price_id" bson:"stripe_price_id"`
}

type TogglePlanTypeRequest struct {
	UserID  helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	NewType helpers.Nullable[PlanType]           `json:"new_type"`
}

type SelectPlanRequest struct {
	Plan   helpers.Nullable[string]             `json:"plan"`
	UserID helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	Type   helpers.Nullable[PlanType]           `json:"type"`
}

type PlanListResponse struct {
	Plans helpers.Nullable[[]PlanResponse] `json:"plans"`
}

type SelectedPlanResponse struct {
	UserID             helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	Plan               helpers.Nullable[string]             `json:"plan"`
	Type               helpers.Nullable[PlanType]           `json:"type"`
	Price              helpers.Nullable[float64]            `json:"price"`
	RemainingTrialDays helpers.Nullable[int]                `json:"remaining_trial_days"`
}
