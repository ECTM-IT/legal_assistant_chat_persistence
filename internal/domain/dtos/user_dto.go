package dtos

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserRequest struct {
	Image          helpers.Nullable[string]               `json:"image" bson:"image"`
	Email          helpers.Nullable[string]               `json:"email" bson:"email"`
	FirstName      helpers.Nullable[string]               `json:"first_name" bson:"first_name"`
	LastName       helpers.Nullable[string]               `json:"last_name" bson:"last_name"`
	Phone          helpers.Nullable[string]               `json:"phone" bson:"phone"`
	CaseIDs        helpers.Nullable[[]primitive.ObjectID] `json:"case_ids" bson:"case_ids"`
	TeamID         helpers.Nullable[primitive.ObjectID]   `json:"team_id" bson:"team_id"`
	AgentIDs       helpers.Nullable[[]primitive.ObjectID] `json:"agent_ids" bson:"agent_ids"`
	SubscriptionID helpers.Nullable[primitive.ObjectID]   `json:"subscription_id" bson:"subscription_id"`
}

type UserResponse struct {
	ID             helpers.Nullable[primitive.ObjectID]   `json:"id" bson:"_id,omitempty"`
	Image          helpers.Nullable[string]               `json:"image" bson:"image"`
	Email          helpers.Nullable[string]               `json:"email" bson:"email"`
	FirstName      helpers.Nullable[string]               `json:"first_name" bson:"first_name"`
	LastName       helpers.Nullable[string]               `json:"last_name" bson:"last_name"`
	Phone          helpers.Nullable[string]               `json:"phone" bson:"phone"`
	CaseIDs        helpers.Nullable[[]primitive.ObjectID] `json:"case_ids" bson:"case_ids"`
	TeamID         helpers.Nullable[primitive.ObjectID]   `json:"team_id" bson:"team_id"`
	AgentIDs       helpers.Nullable[[]primitive.ObjectID] `json:"agent_ids" bson:"agent_ids"`
	SubscriptionID helpers.Nullable[primitive.ObjectID]   `json:"subscription_id" bson:"subscription_id"`
}

type UpdateUserRequest struct {
	Image          helpers.Nullable[string]               `json:"image" bson:"image,omitempty"`
	Email          helpers.Nullable[string]               `json:"email" bson:"email,omitempty"`
	FirstName      helpers.Nullable[string]               `json:"first_name" bson:"first_name,omitempty"`
	LastName       helpers.Nullable[string]               `json:"last_name" bson:"last_name,omitempty"`
	Phone          helpers.Nullable[string]               `json:"phone" bson:"phone,omitempty"`
	CaseIDs        helpers.Nullable[[]primitive.ObjectID] `json:"case_ids" bson:"case_ids,omitempty"`
	TeamID         helpers.Nullable[primitive.ObjectID]   `json:"team_id" bson:"team_id,omitempty"`
	AgentIDs       helpers.Nullable[[]primitive.ObjectID] `json:"agent_ids" bson:"agent_ids,omitempty"`
	SubscriptionID helpers.Nullable[primitive.ObjectID]   `json:"subscription_id" bson:"subscription_id"`
}

type DeleteUserRequest struct {
	ID helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id"`
}
