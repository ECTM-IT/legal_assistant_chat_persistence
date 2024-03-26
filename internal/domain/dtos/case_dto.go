package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCaseRequest struct {
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	CreatorID       primitive.ObjectID   `json:"creator_id"`
	AgentID         primitive.ObjectID   `json:"agent_id"`
	CollaboratorIDs []primitive.ObjectID `json:"collaborator_ids"`
	Action          string               `json:"action"`
	Skill           string               `json:"skill"`
	Share           bool                 `json:"share"`
	IsArchived      bool                 `json:"is_archived"`
}

type CaseResponse struct {
	ID              primitive.ObjectID   `json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	CreatorID       primitive.ObjectID   `json:"creator_id"`
	AgentID         primitive.ObjectID   `json:"agent_id"`
	CollaboratorIDs []primitive.ObjectID `json:"collaborator_ids"`
	Action          string               `json:"action"`
	Skill           string               `json:"skill"`
	Share           bool                 `json:"share"`
	IsArchived      bool                 `json:"is_archived"`
	Messages        []MessageResponse    `json:"messages"`
	LastEdit        time.Time            `json:"last_edit"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

type UpdateCaseRequest struct {
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	AgentID         primitive.ObjectID   `json:"agent_id"`
	CollaboratorIDs []primitive.ObjectID `json:"collaborator_ids"`
	Action          string               `json:"action"`
	Skill           string               `json:"skill"`
	Share           bool                 `json:"share"`
	IsArchived      bool                 `json:"is_archived"`
}

type DeleteCaseRequest struct {
	ID primitive.ObjectID `json:"id"`
}

type MessageResponse struct {
	Content   string             `json:"content"`
	SenderID  primitive.ObjectID `json:"sender_id"`
	Recipient string             `json:"recipient"`
}
