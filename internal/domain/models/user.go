package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	EncryptedName  string               `json:"encrypted_name" bson:"encrypted_name"`
	EncryptedEmail string               `json:"encrypted_email" bson:"encrypted_email"`
	Image          string               `json:"image" bson:"image"`
	FirstName      string               `json:"first_name" bson:"first_name"`
	LastName       string               `json:"last_name" bson:"last_name"`
	Phone          string               `json:"phone" bson:"phone"`
	CaseIDs        []primitive.ObjectID `json:"case_ids" bson:"case_ids"`
	TeamID         primitive.ObjectID   `json:"team_id" bson:"team_id"`
	AgentIDs       []primitive.ObjectID `json:"agent_ids" bson:"agent_ids"`
	SubscriptionID primitive.ObjectID   `json:"subscription_id" bson:"subscription_id"`
	CreationDate   time.Time            `json:"creation_date" bson:"creation_date"`
	LastEdit       time.Time            `json:"last_edit" bson:"last_edit"`
	Share          bool                 `json:"share" bson:"share"`
	IsArchived     bool                 `json:"is_archived" bson:"is_archived"`
}
