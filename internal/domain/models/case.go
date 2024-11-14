package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Case struct {
	ID            primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name          string               `json:"name" bson:"name"`
	CreatorID     primitive.ObjectID   `json:"creator_id" bson:"creator_id"`
	Messages      []Message            `json:"messages" bson:"messages"`
	Collaborators []Collaborators      `json:"collaborators" bson:"collaborators"`
	AgentSkills   []AgentSkillResponse `json:"agent_skills" bson:"agent_skills"`
	Action        string               `json:"action" bson:"action"`
	AgentID       primitive.ObjectID   `json:"agent_id" bson:"agent_id"`
	CreationDate  time.Time            `json:"creation_date" bson:"creation_date"`
	LastEdit      time.Time            `json:"last_edit" bson:"last_edit"`
	Share         bool                 `json:"share" bson:"share"`
	IsArchived    bool                 `json:"is_archived" bson:"is_archived"`
}

type Collaborators struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Edit bool               `json:"edit" bson:"edit"`
}

type AgentSkillResponse struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AgentID primitive.ObjectID `json:"agent_id" bson:"agent_id"`
	Name    string             `json:"name" bson:"name,omitempty"`
}

type Message struct {
	Sender       string `json:"sender" bson:"sender"`
	Recipient    string `json:"recipient" bson:"recipient"`
	Content      string `json:"content" bson:"content"`
	DocumentPath string `json:"document_path" bson:"document_path"`
	FunctionCall bool   `json:"function_call" bson:"function_call"`
}
