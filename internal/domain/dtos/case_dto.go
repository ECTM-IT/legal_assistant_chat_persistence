package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageResponse struct {
	Content   helpers.Nullable[string] `json:"content" bson:"content"`
	Sender    helpers.Nullable[string] `json:"sender_id" bson:"sender_id"`
	Recipient helpers.Nullable[string] `json:"recipient_id" bson:"recipient_id"`
	Skill     helpers.Nullable[string] `json:"skill" bson:"skill"`
}

type CollaboratorResponse struct {
	ID   helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id,omitempty"`
	Edit helpers.Nullable[bool]               `json:"edit" bson:"edit"`
}

type CreateCaseRequest struct {
	Name          helpers.Nullable[string]                 `json:"name" bson:"name"`
	Description   helpers.Nullable[string]                 `json:"description" bson:"description"`
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
	Description   helpers.Nullable[string]                 `json:"description" bson:"description"`
	CreatorID     helpers.Nullable[primitive.ObjectID]     `json:"creator_id" bson:"creator_id"`
	Messages      helpers.Nullable[[]MessageResponse]      `json:"messages" bson:"messages"`
	Collaborators helpers.Nullable[[]CollaboratorResponse] `json:"collaborators" bson:"collaborators"`
	Action        helpers.Nullable[string]                 `json:"action" bson:"action"`
	AgentID       helpers.Nullable[primitive.ObjectID]     `json:"agent_id" bson:"agent_id"`
	LastEdit      helpers.Nullable[time.Time]              `json:"last_edit" bson:"last_edit"`
	Share         helpers.Nullable[bool]                   `json:"share" bson:"share"`
	IsArchived    helpers.Nullable[bool]                   `json:"is_archived" bson:"is_archived"`
}

type UpdateCaseRequest struct {
	Name          helpers.Nullable[string]                 `json:"name" bson:"name,omitempty"`
	Description   helpers.Nullable[string]                 `json:"description" bson:"description,omitempty"`
	Messages      helpers.Nullable[[]MessageResponse]      `json:"messages" bson:"messages,omitempty"`
	Collaborators helpers.Nullable[[]CollaboratorResponse] `json:"collaborators" bson:"collaborators,omitempty"`
	Action        helpers.Nullable[string]                 `json:"action" bson:"action,omitempty"`
	AgentID       helpers.Nullable[primitive.ObjectID]     `json:"agent_id" bson:"agent_id,omitempty"`
	LastEdit      helpers.Nullable[time.Time]              `json:"last_edit" bson:"last_edit,omitempty"`
	Share         helpers.Nullable[bool]                   `json:"share" bson:"share,omitempty"`
	IsArchived    helpers.Nullable[bool]                   `json:"is_archived" bson:"is_archived,omitempty"`
}

type DeleteCaseRequest struct {
	ID helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id"`
}
