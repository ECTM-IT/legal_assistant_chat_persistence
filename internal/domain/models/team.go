package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty"`
	Name  string               `bson:"name"`
	Admin primitive.ObjectID   `bson:"admin"`           // 1-M "relation" with Users
	Users []primitive.ObjectID `bson:"users,omitempty"` // M-M relation with Users
	Cases []primitive.ObjectID `bson:"cases,omitempty"` // M-M relation with Cases
}
