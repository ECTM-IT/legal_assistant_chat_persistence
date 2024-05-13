package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Image          string               `json:"image" bson:"image"`
	Email          string               `json:"email" bson:"email"`
	FirstName      string               `json:"first_name" bson:"first_name"`
	LastName       string               `json:"last_name" bson:"last_name"`
	Phone          string               `json:"phone" bson:"phone"`
	CaseIDs        []primitive.ObjectID `json:"case_ids" bson:"case_ids"`
	TeamID         primitive.ObjectID   `json:"team_id" bson:"team_id"`
	AgentIDs       []primitive.ObjectID `json:"agent_ids" bson:"agent_ids"`
	SubscriptionID primitive.ObjectID   `json:"subscription_id" bson:"subscription_id"` //choose encoding strategy to handle payments
}
