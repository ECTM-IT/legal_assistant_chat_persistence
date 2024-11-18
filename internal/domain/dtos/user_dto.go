package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserRequest struct {
	EncryptedName  helpers.Nullable[string]               `json:"encrypted_name" bson:"encrypted_name"`
	EncryptedEmail helpers.Nullable[string]               `json:"encrypted_email" bson:"encrypted_email"`
	Image          helpers.Nullable[string]               `json:"image" bson:"image"`
	FirstName      helpers.Nullable[string]               `json:"first_name" bson:"first_name"`
	LastName       helpers.Nullable[string]               `json:"last_name" bson:"last_name"`
	Phone          helpers.Nullable[string]               `json:"phone" bson:"phone"`
	CaseIDs        helpers.Nullable[[]primitive.ObjectID] `json:"case_ids" bson:"case_ids"`
	TeamID         helpers.Nullable[primitive.ObjectID]   `json:"team_id" bson:"team_id"`
	AgentIDs       helpers.Nullable[[]primitive.ObjectID] `json:"agent_ids" bson:"agent_ids"`
	SubscriptionID helpers.Nullable[primitive.ObjectID]   `json:"subscription_id" bson:"subscription_id"`
	CreationDate   helpers.Nullable[time.Time]            `json:"creation_date" bson:"creation_date"`
	LastEdit       helpers.Nullable[time.Time]            `json:"last_edit" bson:"last_edit"`
	Share          helpers.Nullable[bool]                 `json:"share" bson:"share"`
	IsArchived     helpers.Nullable[bool]                 `json:"is_archived" bson:"is_archived"`
}

type UserResponse struct {
	ID             helpers.Nullable[primitive.ObjectID]   `json:"id" bson:"_id,omitempty"`
	EncryptedName  helpers.Nullable[string]               `json:"encrypted_name" bson:"encrypted_name"`
	EncryptedEmail helpers.Nullable[string]               `json:"encrypted_email" bson:"encrypted_email"`
	Image          helpers.Nullable[string]               `json:"image" bson:"image"`
	FirstName      helpers.Nullable[string]               `json:"first_name" bson:"first_name"`
	LastName       helpers.Nullable[string]               `json:"last_name" bson:"last_name"`
	Phone          helpers.Nullable[string]               `json:"phone" bson:"phone"`
	CaseIDs        helpers.Nullable[[]primitive.ObjectID] `json:"case_ids" bson:"case_ids"`
	TeamID         helpers.Nullable[primitive.ObjectID]   `json:"team_id" bson:"team_id"`
	AgentIDs       helpers.Nullable[[]primitive.ObjectID] `json:"agent_ids" bson:"agent_ids"`
	SubscriptionID helpers.Nullable[primitive.ObjectID]   `json:"subscription_id" bson:"subscription_id"`
	CreationDate   helpers.Nullable[time.Time]            `json:"creation_date" bson:"creation_date"`
	LastEdit       helpers.Nullable[time.Time]            `json:"last_edit" bson:"last_edit"`
	Share          helpers.Nullable[bool]                 `json:"share" bson:"share"`
	IsArchived     helpers.Nullable[bool]                 `json:"is_archived" bson:"is_archived"`
}

type UpdateUserRequest struct {
	EncryptedName  helpers.Nullable[string]               `json:"encrypted_name" bson:"encrypted_name,omitempty"`
	EncryptedEmail helpers.Nullable[string]               `json:"encrypted_email" bson:"encrypted_email,omitempty"`
	Image          helpers.Nullable[string]               `json:"image" bson:"image,omitempty"`
	FirstName      helpers.Nullable[string]               `json:"first_name" bson:"first_name,omitempty"`
	LastName       helpers.Nullable[string]               `json:"last_name" bson:"last_name,omitempty"`
	Phone          helpers.Nullable[string]               `json:"phone" bson:"phone,omitempty"`
	CaseIDs        helpers.Nullable[[]primitive.ObjectID] `json:"case_ids" bson:"case_ids,omitempty"`
	TeamID         helpers.Nullable[primitive.ObjectID]   `json:"team_id" bson:"team_id,omitempty"`
	AgentIDs       helpers.Nullable[[]primitive.ObjectID] `json:"agent_ids" bson:"agent_ids,omitempty"`
	SubscriptionID helpers.Nullable[primitive.ObjectID]   `json:"subscription_id" bson:"subscription_id,omitempty"`
}

type DeleteUserRequest struct {
	ID helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id"`
}
