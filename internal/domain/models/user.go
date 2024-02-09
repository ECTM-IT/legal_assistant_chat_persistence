package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty"`
	Name  string               `bson:"name"`
	Plans primitive.M          `bson:"plans"`
	Cases []primitive.ObjectID `bson:"cases,omitempty"` // M-M relation with Cases
	Teams []primitive.ObjectID `bson:"teams,omitempty"` // M-M realation with Teams
}
