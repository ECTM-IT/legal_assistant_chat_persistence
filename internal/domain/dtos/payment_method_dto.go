package dtos

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentMethodRequest represents a request to add a new payment method
type PaymentMethodRequest struct {
	UserID          helpers.Nullable[string]            `json:"user_id"`
	PaymentMethodID helpers.Nullable[string]            `json:"payment_method_id"`
	IsDefault       helpers.Nullable[bool]              `json:"is_default"`
	BillingDetails  helpers.Nullable[BillingDetailsDTO] `json:"billing_details"`
}

// BillingDetailsDTO represents billing details for a payment method
type BillingDetailsDTO struct {
	Name  helpers.Nullable[string] `json:"name"`
	Email helpers.Nullable[string] `json:"email"`
}

// PaymentMethodResponse represents a payment method returned to the client
type PaymentMethodResponse struct {
	ID              helpers.Nullable[primitive.ObjectID] `json:"id"`
	UserID          helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	Type            helpers.Nullable[string]             `json:"type"`
	CardBrand       helpers.Nullable[string]             `json:"card_brand"`
	LastFour        helpers.Nullable[string]             `json:"last_four"`
	ExpirationMonth helpers.Nullable[int]                `json:"expiration_month"`
	ExpirationYear  helpers.Nullable[int]                `json:"expiration_year"`
	BillingName     helpers.Nullable[string]             `json:"billing_name"`
	BillingEmail    helpers.Nullable[string]             `json:"billing_email"`
	IsDefault       helpers.Nullable[bool]               `json:"is_default"`
	CreatedAt       helpers.Nullable[int64]              `json:"created_at"`
}
