package dtos

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAgentRequest struct {
	ProfileImage helpers.Nullable[string]   `json:"profile_image"`
	Name         helpers.Nullable[string]   `json:"name"`
	Description  helpers.Nullable[string]   `json:"description"`
	Skills       helpers.Nullable[[]string] `json:"skills"`
	Price        helpers.Nullable[float64]  `json:"price"`
	Code         helpers.Nullable[string]   `json:"code"`
}

type AgentResponse struct {
	ID           helpers.Nullable[primitive.ObjectID] `json:"id"`
	ProfileImage helpers.Nullable[string]             `json:"profile_image"`
	Name         helpers.Nullable[string]             `json:"name"`
	Description  helpers.Nullable[string]             `json:"description"`
	Skills       helpers.Nullable[[]SkillResponse]    `json:"skills"`
	Price        helpers.Nullable[float64]            `json:"price"`
	Code         helpers.Nullable[string]             `json:"code"`
}

type SkillResponse struct {
	Name         helpers.Nullable[string]   `json:"name"`
	Descriptions helpers.Nullable[[]string] `json:"descriptions"`
}

type UpdateAgentRequest struct {
	ProfileImage helpers.Nullable[string]   `json:"profile_image"`
	Name         helpers.Nullable[string]   `json:"name"`
	Description  helpers.Nullable[string]   `json:"description"`
	Skills       helpers.Nullable[[]string] `json:"skills"`
	Price        helpers.Nullable[float64]  `json:"price"`
	Code         helpers.Nullable[string]   `json:"code"`
}
