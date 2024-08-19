package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageResponse struct {
	Content      helpers.Nullable[string] `json:"content" bson:"content"`
	Sender       helpers.Nullable[string] `json:"sender_id" bson:"sender_id"`
	Recipient    helpers.Nullable[string] `json:"recipient_id" bson:"recipient_id"`
	FunctionCall helpers.Nullable[bool]   `json:"function_call" bson:"function_call"`
}

type CollaboratorResponse struct {
	ID   helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id,omitempty"`
	Edit helpers.Nullable[bool]               `json:"edit" bson:"edit"`
}

type CreateCaseRequest struct {
	Name          helpers.Nullable[string]                 `json:"name" bson:"name"`
	CreatorID     helpers.Nullable[primitive.ObjectID]     `json:"creator_id" bson:"creator_id"`
	Messages      helpers.Nullable[[]MessageResponse]      `json:"messages" bson:"messages"`
	Collaborators helpers.Nullable[[]CollaboratorResponse] `json:"collaborators" bson:"collaborators"`
	Action        helpers.Nullable[string]                 `json:"action" bson:"action"`
	AgentID       helpers.Nullable[primitive.ObjectID]     `json:"agent_id" bson:"agent_id"`
	LastEdit      helpers.Nullable[time.Time]              `json:"last_edit" bson:"last_edit"`
	Share         helpers.Nullable[bool]                   `json:"share" bson:"share"`
	IsArchived    helpers.Nullable[bool]                   `json:"is_archived" bson:"is_archived"`
}

type CaseResponse struct {
	ID            helpers.Nullable[primitive.ObjectID]     `json:"id" bson:"_id,omitempty"`
	Name          helpers.Nullable[string]                 `json:"name" bson:"name"`
	CreatorID     helpers.Nullable[primitive.ObjectID]     `json:"creator_id" bson:"creator_id"`
	Messages      helpers.Nullable[[]MessageResponse]      `json:"messages" bson:"messages"`
	Collaborators helpers.Nullable[[]CollaboratorResponse] `json:"collaborators" bson:"collaborators"`
	Action        helpers.Nullable[string]                 `json:"action" bson:"action"`
	AgentID       helpers.Nullable[primitive.ObjectID]     `json:"agent_id" bson:"agent_id"`
	CreationDate  helpers.Nullable[time.Time]              `json:"creation_date" bson:"creation_date"`
	LastEdit      helpers.Nullable[time.Time]              `json:"last_edit" bson:"last_edit"`
	Share         helpers.Nullable[bool]                   `json:"share" bson:"share"`
	IsArchived    helpers.Nullable[bool]                   `json:"is_archived" bson:"is_archived"`
}

type UpdateCaseRequest struct {
	Name          helpers.Nullable[string]                 `json:"name" bson:"name,omitempty"`
	Messages      helpers.Nullable[[]MessageResponse]      `json:"messages" bson:"messages,omitempty"`
	Collaborators helpers.Nullable[[]CollaboratorResponse] `json:"collaborators" bson:"collaborators,omitempty"`
	Action        helpers.Nullable[string]                 `json:"action" bson:"action,omitempty"`
	AgentID       helpers.Nullable[primitive.ObjectID]     `json:"agent_id" bson:"agent_id,omitempty"`
	LastEdit      helpers.Nullable[time.Time]              `json:"last_edit" bson:"last_edit,omitempty"`
	Share         helpers.Nullable[bool]                   `json:"share" bson:"share,omitempty"`
	IsArchived    helpers.Nullable[bool]                   `json:"is_archived" bson:"is_archived,omitempty"`
}

type AddCollaboratorToCase struct {
	Edit  helpers.Nullable[bool]   `json:"edit" bson:"edit,omitempty"`
	Email helpers.Nullable[string] `json:"email" bson:"email,omitempty"`
}
type DeleteCaseRequest struct {
	ID helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id"`
}
