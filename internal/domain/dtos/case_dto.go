package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Start of Selection
type MessageResponse struct {
	Content      helpers.Nullable[string] `json:"content,omitempty" bson:"content"`
	Sender       helpers.Nullable[string] `json:"sender,omitempty" bson:"sender"`
	Recipient    helpers.Nullable[string] `json:"recipient,omitempty" bson:"recipient"`
	FunctionCall helpers.Nullable[bool]   `json:"function_call,omitempty" bson:"function_call"`
	DocumentPath helpers.Nullable[string] `json:"document_path,omitempty" bson:"document_path"`
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
	Documents     helpers.Nullable[[]DocumentResponse]     `json:"documents" bson:"documents"`
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
	Documents     helpers.Nullable[[]DocumentResponse]     `json:"documents" bson:"documents"`
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

type DocumentResponse struct {
	ID          helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id,omitempty"`
	FileName    helpers.Nullable[string]             `json:"file_name" bson:"file_name"`
	FileType    helpers.Nullable[string]             `json:"file_type" bson:"file_type"`
	FileContent helpers.Nullable[[]byte]             `json:"file_content" bson:"file_content"`
	UploadDate  helpers.Nullable[time.Time]          `json:"upload_date" bson:"upload_date"`
}

type AddDocumentToCase struct {
	FileName    string `json:"file_name" validate:"required"`
	FileType    string `json:"file_type" validate:"required"`    // e.g., "pdf", "docx", "xlsx"
	FileContent []byte `json:"file_content" validate:"required"` // The actual file content (e.g., in base64 format if sending as JSON)
}
