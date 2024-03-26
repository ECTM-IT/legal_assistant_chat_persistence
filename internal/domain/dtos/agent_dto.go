package dtos

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAgentRequest struct {
	ProfileImage string   `json:"profile_image"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Skills       []string `json:"skills"`
	Price        float64  `json:"price"`
	Code         string   `json:"code"`
}

type AgentResponse struct {
	ID           primitive.ObjectID `json:"id"`
	ProfileImage string             `json:"profile_image"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Skills       []SkillResponse    `json:"skills"`
	Price        float64            `json:"price"`
	Code         string             `json:"code"`
}

type SkillResponse struct {
	Name         string   `json:"name"`
	Descriptions []string `json:"descriptions"`
}

type UpdateAgentRequest struct {
	ProfileImage string   `json:"profile_image"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Skills       []string `json:"skills"`
	Price        float64  `json:"price"`
	Code         string   `json:"code"`
}
