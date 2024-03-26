package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agent struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProfileImage string             `json:"profile_image" bson:"profile_image"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	Skills       []Skill            `json:"skills" bson:"skills"`
	Price        float64            `json:"price" bson:"price"`
	Code         string             `json:"code" bson:"code"`
}

type Skill struct {
	Name         string   `json:"name" bson:"name"`
	Descriptions []string `json:"descriptions" bson:"descriptions"`
}
