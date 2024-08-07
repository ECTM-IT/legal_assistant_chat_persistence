package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Case struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	Description   string             `json:"description" bson:"description"`
	CreatorID     primitive.ObjectID `json:"creator_id" bson:"creator_id"`
	Messages      []Message          `json:"messages" bson:"messages"`
	Collaborators []Collaborators    `json:"collaborator_ids" bson:"collaborator_ids"`
	Action        string             `json:"action" bson:"action"`
	AgentID       primitive.ObjectID `json:"agent_id" bson:"agent_id"`
	LastEdit      time.Time          `json:"last_edit" bson:"last_edit"`
	Share         bool               `json:"share" bson:"share"`
	IsArchived    bool               `json:"is_archived" bson:"is_archived"`
}

type Collaborators struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Edit bool               `json:"edit" bson:"edit"`
}

type Message struct {
	Content      string `json:"content" bson:"content"`
	SenderID     string `json:"sender_id" bson:"sender_id"`
	RecipientID  string `json:"recipient_id" bson:"recipient_id"`
	FunctionCall bool   `json:"function_call" bson:"function_call"`
}
