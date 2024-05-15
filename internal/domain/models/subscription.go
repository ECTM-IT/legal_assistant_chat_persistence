package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscriptions struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Plan   string             `json:"plan" bson:"plan"`
	Expiry time.Time          `json:"expiry" bson:"expiry"`
}
