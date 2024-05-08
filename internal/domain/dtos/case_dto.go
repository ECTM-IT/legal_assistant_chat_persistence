package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageResponse struct {
	Content     helpers.Nullable[string]             `json:"content"`
	SenderID    helpers.Nullable[primitive.ObjectID] `json:"sender_id"`
	RecipientID helpers.Nullable[primitive.ObjectID] `json:"recipient_id"`
	Skill       helpers.Nullable[string]             `json:"skill"`
}

type CreateCaseRequest struct {
	Name            helpers.Nullable[string]               `json:"name"`
	Description     helpers.Nullable[string]               `json:"description"`
	CreatorID       helpers.Nullable[primitive.ObjectID]   `json:"creator_id"`
	Messages        helpers.Nullable[[]MessageResponse]    `json:"messages"`
	CollaboratorIDs helpers.Nullable[[]primitive.ObjectID] `json:"collaborator_ids"`
	Action          helpers.Nullable[string]               `json:"action"`
	AgentID         helpers.Nullable[primitive.ObjectID]   `json:"agent_id"`
	LastEdit        helpers.Nullable[time.Time]            `json:"last_edit"`
	Share           helpers.Nullable[bool]                 `json:"share"`
	IsArchived      helpers.Nullable[bool]                 `json:"is_archived"`
}

type CaseResponse struct {
	ID              helpers.Nullable[primitive.ObjectID]   `json:"id"`
	Name            helpers.Nullable[string]               `json:"name"`
	Description     helpers.Nullable[string]               `json:"description"`
	CreatorID       helpers.Nullable[primitive.ObjectID]   `json:"creator_id"`
	AgentID         helpers.Nullable[primitive.ObjectID]   `json:"agent_id"`
	CollaboratorIDs helpers.Nullable[[]primitive.ObjectID] `json:"collaborator_ids"`
	Action          helpers.Nullable[string]               `json:"action"`
	Skill           helpers.Nullable[string]               `json:"skill"`
	Share           helpers.Nullable[bool]                 `json:"share"`
	IsArchived      helpers.Nullable[bool]                 `json:"is_archived"`
	Messages        helpers.Nullable[[]MessageResponse]    `json:"messages"`
	LastEdit        helpers.Nullable[time.Time]            `json:"last_edit"`
	CreatedAt       helpers.Nullable[time.Time]            `json:"created_at"`
	UpdatedAt       helpers.Nullable[time.Time]            `json:"updated_at"`
}

type UpdateCaseRequest struct {
	Name            helpers.Nullable[string]               `json:"name"`
	Description     helpers.Nullable[string]               `json:"description"`
	AgentID         helpers.Nullable[primitive.ObjectID]   `json:"agent_id"`
	CollaboratorIDs helpers.Nullable[[]primitive.ObjectID] `json:"collaborator_ids"`
	Action          helpers.Nullable[string]               `json:"action"`
	Skill           helpers.Nullable[string]               `json:"skill"`
	Share           helpers.Nullable[bool]                 `json:"share"`
	IsArchived      helpers.Nullable[bool]                 `json:"is_archived"`
}

type DeleteCaseRequest struct {
	ID helpers.Nullable[primitive.ObjectID] `json:"id"`
}
