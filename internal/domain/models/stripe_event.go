package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StripeEvent represents a processed Stripe webhook event
type StripeEvent struct {
	ID        primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	EventID   string                 `json:"event_id" bson:"event_id"`
	Type      string                 `json:"type" bson:"type"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	Data      map[string]interface{} `json:"data" bson:"data"`
	Processed bool                   `json:"processed" bson:"processed"`
}
