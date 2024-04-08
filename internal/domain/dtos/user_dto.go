package dtos

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserRequest struct {
	Image          helpers.Nullable[string]               `json:"image"`
	Email          helpers.Nullable[string]               `json:"email"`
	FirstName      helpers.Nullable[string]               `json:"first_name"`
	LastName       helpers.Nullable[string]               `json:"last_name"`
	Phone          helpers.Nullable[string]               `json:"phone"`
	CaseIDs        helpers.Nullable[[]primitive.ObjectID] `json:"case_ids"`
	TeamID         helpers.Nullable[primitive.ObjectID]   `json:"team_id"`
	AgentIDs       helpers.Nullable[[]primitive.ObjectID] `json:"agent_ids"`
	SubscriptionID helpers.Nullable[primitive.ObjectID]   `json:"subscription_id"`
}

type UserResponse struct {
	ID             helpers.Nullable[primitive.ObjectID]   `json:"id"`
	Image          helpers.Nullable[string]               `json:"image"`
	Email          helpers.Nullable[string]               `json:"email"`
	FirstName      helpers.Nullable[string]               `json:"first_name"`
	LastName       helpers.Nullable[string]               `json:"last_name"`
	Phone          helpers.Nullable[string]               `json:"phone"`
	CaseIDs        helpers.Nullable[[]primitive.ObjectID] `json:"case_ids"`
	TeamID         helpers.Nullable[primitive.ObjectID]   `json:"team_id"`
	AgentIDs       helpers.Nullable[[]primitive.ObjectID] `json:"agent_ids"`
	SubscriptionID helpers.Nullable[primitive.ObjectID]   `json:"subscription_id"`
}

type UpdateUserRequest struct {
	Image          helpers.Nullable[string]   `json:"image"`
	Email          helpers.Nullable[string]   `json:"email"`
	FirstName      helpers.Nullable[string]   `json:"first_name"`
	LastName       helpers.Nullable[string]   `json:"last_name"`
	Phone          helpers.Nullable[string]   `json:"phone"`
	CaseIDs        helpers.Nullable[[]string] `json:"case_ids"`
	TeamID         helpers.Nullable[string]   `json:"team_id"`
	AgentIDs       helpers.Nullable[[]string] `json:"agent_ids"`
	SubscriptionID helpers.Nullable[string]   `json:"subscription_id"`
}

type DeleteUserRequest struct {
	ID helpers.Nullable[primitive.ObjectID] `json:"id"`
}
