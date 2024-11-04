package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Start of Selection
type MessageResponse struct {
	Content      helpers.Nullable[string]     `json:"content,omitempty" bson:"content"`
	Sender       helpers.Nullable[string]     `json:"sender,omitempty" bson:"sender"`
	Recipient    helpers.Nullable[string]     `json:"recipient,omitempty" bson:"recipient"`
	FunctionCall helpers.Nullable[bool]       `json:"function_call,omitempty" bson:"function_call"`
	DocumentPath helpers.Nullable[string]     `json:"document_path,omitempty" bson:"document_path"`
	Feedbacks    helpers.Nullable[[]Feedback] `json:"feedbacks,omitempty" bson:"feedbacks"`
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

type Feedback struct {
	ID           helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id,omitempty"`
	CaseID       helpers.Nullable[primitive.ObjectID] `json:"case_id" bson:"case_id" validate:"required"`
	MessageID    helpers.Nullable[primitive.ObjectID] `json:"message_id" bson:"message_id,omitempty"`
	CreatorID    helpers.Nullable[primitive.ObjectID] `json:"creator_id" bson:"creator_id,omitempty"`
	Score        helpers.Nullable[string]             `json:"score" bson:"score,omitempty"`
	Reasons      helpers.Nullable[[]string]           `json:"reasons" bson:"reasons,omitempty"`
	Comment      helpers.Nullable[string]             `json:"comment" bson:"comment,omitempty"`
	CreationDate helpers.Nullable[time.Time]          `json:"creation_date" bson:"creation_date,omitempty"`
}

type AddFeedbackRequest struct {
	CaseID       primitive.ObjectID `json:"case_id" bson:"case_id" validate:"required"`
	MessageID    primitive.ObjectID `json:"message_id" bson:"message_id" validate:"required"`
	CreatorID    primitive.ObjectID `json:"creator_id" bson:"creator_id" validate:"required"`
	Score        string             `json:"score" bson:"score" validate:"required"`
	Reasons      []string           `json:"reasons" bson:"reasons"`
	Comment      string             `json:"comment" bson:"comment"`
	CreationDate time.Time          `json:"creation_date" bson:"creation_date"`
}
