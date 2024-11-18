package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Case struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	CreatorID     primitive.ObjectID `json:"creator_id" bson:"creator_id"`
	Messages      []Message          `json:"messages" bson:"messages"`
	Collaborators []Collaborators    `json:"collaborators" bson:"collaborators"`
	Documents     []Document         `json:"documents" bson:"documents"`
	Action        string             `json:"action" bson:"action"`
	AgentID       primitive.ObjectID `json:"agent_id" bson:"agent_id"`
	CreationDate  time.Time          `json:"creation_date" bson:"creation_date"`
	LastEdit      time.Time          `json:"last_edit" bson:"last_edit"`
	Share         bool               `json:"share" bson:"share"`
	IsArchived    bool               `json:"is_archived" bson:"is_archived"`
}

type Collaborators struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Edit bool               `json:"edit" bson:"edit"`
}

type Message struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Sender       string             `json:"sender" bson:"sender"`
	Recipient    string             `json:"recipient" bson:"recipient"`
	Content      string             `json:"content" bson:"content"`
	DocumentPath string             `json:"document_path" bson:"document_path"`
	FunctionCall bool               `json:"function_call" bson:"function_call"`
	Feedbacks    []Feedback         `json:"feedbacks" bson:"feedbacks"`
}

type Feedback struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CaseID       primitive.ObjectID `json:"case_id" bson:"case_id" validate:"required"`
	MessageID    primitive.ObjectID `json:"message_id" bson:"message_id"`
	CreatorID    primitive.ObjectID `json:"creator_id" bson:"creator_id"`
	Score        string             `json:"score" bson:"score"`
	Reasons      []string           `json:"reasons" bson:"reasons"`
	Comment      string             `json:"comment" bson:"comment"`
	CreationDate time.Time          `json:"creation_date" bson:"creation_date"`
}

type Document struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileName    string             `json:"file_name" bson:"file_name"`
	FileType    string             `json:"file_type" bson:"file_type"` // e.g., pdf, docx, xls
	FileContent []byte             `json:"file_content" bson:"file_content"`
	UploadDate  time.Time          `json:"upload_date" bson:"upload_date"`
}
