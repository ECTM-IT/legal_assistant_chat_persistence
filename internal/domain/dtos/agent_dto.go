package dtos

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAgentRequest struct {
	ProfileImage helpers.Nullable[string]   `json:"profile_image" bson:"profile_image"`
	Name         helpers.Nullable[string]   `json:"name" bson:"name"`
	Description  helpers.Nullable[string]   `json:"description" bson:"description"`
	Skills       helpers.Nullable[[]string] `json:"skills" bson:"skills"`
	Price        helpers.Nullable[float64]  `json:"price" bson:"price"`
	Code         helpers.Nullable[string]   `json:"code" bson:"code"`
}

type AgentResponse struct {
	ID           helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id,omitempty"`
	ProfileImage helpers.Nullable[string]             `json:"profile_image" bson:"profile_image"`
	Name         helpers.Nullable[string]             `json:"name" bson:"name"`
	Description  helpers.Nullable[string]             `json:"description" bson:"description"`
	Skills       helpers.Nullable[[]string]           `json:"skills" bson:"skills"`
	Price        helpers.Nullable[float64]            `json:"price" bson:"price"`
	Code         helpers.Nullable[string]             `json:"code" bson:"code"`
}

type UpdateAgentRequest struct {
	ProfileImage helpers.Nullable[string]   `json:"profile_image" bson:"profile_image,omitempty"`
	Name         helpers.Nullable[string]   `json:"name" bson:"name,omitempty"`
	Description  helpers.Nullable[string]   `json:"description" bson:"description,omitempty"`
	Skills       helpers.Nullable[[]string] `json:"skills" bson:"skills,omitempty"`
	Price        helpers.Nullable[float64]  `json:"price" bson:"price,omitempty"`
	Code         helpers.Nullable[string]   `json:"code" bson:"code,omitempty"`
}

type DeleteAgentRequest struct {
	ID helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id"`
}
