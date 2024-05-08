package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Case struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name            string               `json:"name" bson:"name"`
	Description     string               `json:"description" bson:"description"`
	CreatorID       primitive.ObjectID   `json:"creator_id" bson:"creator_id"`
	Messages        []Message            `json:"messages" bson:"messages"`
	CollaboratorIDs []primitive.ObjectID `json:"collaborator_ids" bson:"collaborator_ids"`
	Action          string               `json:"action" bson:"action"`
	AgentID         primitive.ObjectID   `json:"agent_id" bson:"agent_id"`
	LastEdit        time.Time            `json:"last_edit" bson:"last_edit"`
	Share           bool                 `json:"share" bson:"share"`
	IsArchived      bool                 `json:"is_archived" bson:"is_archived"`
}

type Message struct {
	Content     string             `json:"content" bson:"content"`
	SenderID    primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	RecipientID primitive.ObjectID `json:"recipient_id" bson:"recipient_id"`
	Skill       string             `json:"skill" bson:"skill"`
}
