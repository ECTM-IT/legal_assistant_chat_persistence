package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateUserRequest struct {
	Image          string               `json:"image"`
	Email          string               `json:"email"`
	FirstName      string               `json:"first_name"`
	LastName       string               `json:"last_name"`
	Phone          string               `json:"phone"`
	CaseIDs        []primitive.ObjectID `json:"case_ids"`
	TeamID         primitive.ObjectID   `json:"team_id"`
	AgentIDs       []primitive.ObjectID `json:"agent_ids"`
	SubscriptionID primitive.ObjectID   `json:"subscription_id"`
}

type UserResponse struct {
	ID             primitive.ObjectID   `json:"id"`
	Image          string               `json:"image"`
	Email          string               `json:"email"`
	FirstName      string               `json:"first_name"`
	LastName       string               `json:"last_name"`
	Phone          string               `json:"phone"`
	CaseIDs        []primitive.ObjectID `json:"case_ids"`
	TeamID         primitive.ObjectID   `json:"team_id"`
	AgentIDs       []primitive.ObjectID `json:"agent_ids"`
	SubscriptionID primitive.ObjectID   `json:"subscription_id"`
}

type UpdateUserRequest struct {
	Image          *string  `json:"image"`
	Email          *string  `json:"email"`
	FirstName      *string  `json:"first_name"`
	LastName       *string  `json:"last_name"`
	Phone          *string  `json:"phone"`
	CaseIDs        []string `json:"case_ids"`
	TeamID         *string  `json:"team_id"`
	AgentIDs       []string `json:"agent_ids"`
	SubscriptionID *string  `json:"subscription_id"`
}

type DeleteUserRequest struct {
	ID primitive.ObjectID `json:"id"`
}
