package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	AdminID primitive.ObjectID `bson:"admin_id"`
	Members []TeamMember       `bson:"members"`
}

type TeamMember struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	DateAdded  time.Time          `bson:"date_added"`
	LastActive time.Time          `bson:"last_active"`
}
