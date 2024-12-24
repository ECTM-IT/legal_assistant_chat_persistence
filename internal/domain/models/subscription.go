package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscriptions struct {
	ID                   primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID               primitive.ObjectID     `json:"user_id" bson:"user_id"`
	Plan                 Plan                   `json:"plan" bson:"plan"`
	Expiry               time.Time              `json:"expiry" bson:"expiry"`
	Status               string                 `json:"status" bson:"status"` // e.g., "active", "canceled", "past_due"
	StripeCustomerID     string                 `json:"stripe_customer_id" bson:"stripe_customer_id"`
	StripeSubscriptionID string                 `json:"stripe_subscription_id" bson:"stripe_subscription_id"`
	CurrentPeriodStart   time.Time              `json:"current_period_start" bson:"current_period_start"`
	CurrentPeriodEnd     time.Time              `json:"current_period_end" bson:"current_period_end"`
	CancelAtPeriodEnd    bool                   `json:"cancel_at_period_end" bson:"cancel_at_period_end"`
	BillingInformations  map[string]interface{} `json:"billing_informations" bson:"billing_informations"`
}
