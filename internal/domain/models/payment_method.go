package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentMethod represents a payment method stored in the database
type PaymentMethod struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	StripeCustomerID string             `json:"stripe_customer_id" bson:"stripe_customer_id"`
	StripePaymentID  string             `json:"stripe_payment_id" bson:"stripe_payment_id"`
	Type             string             `json:"type" bson:"type"` // e.g., "card", "bank_account"
	CardBrand        string             `json:"card_brand" bson:"card_brand,omitempty"`
	LastFour         string             `json:"last_four" bson:"last_four,omitempty"`
	ExpirationMonth  int                `json:"expiration_month" bson:"expiration_month,omitempty"`
	ExpirationYear   int                `json:"expiration_year" bson:"expiration_year,omitempty"`
	BillingName      string             `json:"billing_name" bson:"billing_name,omitempty"`
	BillingEmail     string             `json:"billing_email" bson:"billing_email,omitempty"`
	IsDefault        bool               `json:"is_default" bson:"is_default"`
	CreatedAt        int64              `json:"created_at" bson:"created_at"`
}
