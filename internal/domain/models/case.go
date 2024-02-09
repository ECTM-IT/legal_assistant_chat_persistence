package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Case struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	Messages    bson.D              `bson:"messages"`
	Type        string              `bson:"type"`
	Title       string              `bson:"title"`
	Description string              `bson:"description"`
	CreatedAt   time.Time           `bson:"created_at"`
	LastUpdate  primitive.Timestamp `bson:"last_updated"`
}

// "lastUpdate": primitive.Timestamp{T:uint32(time.Now().Unix())}
// or primitive.NewDateTimeFromTime
