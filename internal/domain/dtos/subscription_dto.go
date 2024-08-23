package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubscriptionType represents the type of subscription
type SubscriptionType string

const (
	Monthly SubscriptionType = "monthly"
	Annual  SubscriptionType = "annual"
)

// BillingType represents the type of billing information
type BillingType string

const (
	Freelancer              BillingType = "freelancer"
	IndividualEnterprise    BillingType = "individual_enterprise"
	Company                 BillingType = "company"
	ProfessionalAssociation BillingType = "professional_association"
	PhysicalPerson          BillingType = "physical_person"
)

type CreateSubscriptionRequest struct {
	Plan                helpers.Nullable[string]             `json:"plan" bson:"plan"`
	Expiry              helpers.Nullable[time.Time]          `json:"expiry" bson:"expiry"`
	Type                helpers.Nullable[SubscriptionType]   `json:"type" bson:"type"`
	BillingInformations helpers.Nullable[BillingInformation] `json:"billing_informations" bson:"billing_informations"`
}

type UpdateSubscriptionRequest struct {
	Plan                helpers.Nullable[string]             `json:"plan" bson:"plan,omitempty"`
	Expiry              helpers.Nullable[time.Time]          `json:"expiry" bson:"expiry,omitempty"`
	Type                helpers.Nullable[SubscriptionType]   `json:"type" bson:"type,omitempty"`
	BillingInformations helpers.Nullable[BillingInformation] `json:"billing_informations" bson:"billing_informations,omitempty"`
}

type SubscriptionResponse struct {
	ID                  helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id,omitempty"`
	Plan                helpers.Nullable[string]             `json:"plan" bson:"plan"`
	Expiry              helpers.Nullable[time.Time]          `json:"expiry" bson:"expiry"`
	Type                helpers.Nullable[SubscriptionType]   `json:"type" bson:"type"`
	BillingInformations helpers.Nullable[BillingInformation] `json:"billing_informations" bson:"billing_informations"`
}

type CommonBillingInfo struct {
	VatNumber  helpers.Nullable[string] `json:"vat_number" bson:"vat_number"`
	SdiCode    helpers.Nullable[string] `json:"sdi_code" bson:"sdi_code"`
	PecAddress helpers.Nullable[string] `json:"pec_address" bson:"pec_address"`
}

type BillingInformation struct {
	Type BillingType `json:"type" bson:"type"`
	Info interface{} `json:"info" bson:"info"`
}

type FreelancerInfo struct {
	CommonBillingInfo
	FirstName           helpers.Nullable[string] `json:"firstname" bson:"firstname"`
	LastName            helpers.Nullable[string] `json:"lastname" bson:"lastname"`
	ProfessionalAddress helpers.Nullable[string] `json:"professional_address" bson:"professional_address"`
	TaxCode             helpers.Nullable[string] `json:"tax_code" bson:"tax_code"`
}

type IndividualEnterpriseInfo struct {
	CommonBillingInfo
	FirstName      helpers.Nullable[string] `json:"firstname" bson:"firstname"`
	LastName       helpers.Nullable[string] `json:"lastname" bson:"lastname"`
	CompanyAddress helpers.Nullable[string] `json:"company_address" bson:"company_address"`
	HolderTaxCode  helpers.Nullable[string] `json:"holder_tax_code" bson:"holder_tax_code"`
}

type CompanyInfo struct {
	CommonBillingInfo
	CompanyName    helpers.Nullable[string] `json:"company_name" bson:"company_name"`
	LegalAddress   helpers.Nullable[string] `json:"legal_address" bson:"legal_address"`
	CompanyTaxCode helpers.Nullable[string] `json:"company_tax_code" bson:"company_tax_code"`
}

type ProfessionalAssociationInfo struct {
	CommonBillingInfo
	AssociationName helpers.Nullable[string] `json:"association_name" bson:"association_name"`
	Address         helpers.Nullable[string] `json:"address" bson:"address"`
	TaxCode         helpers.Nullable[string] `json:"tax_code" bson:"tax_code"`
}

type PhysicalPersonInfo struct {
	CommonBillingInfo
	FirstName          helpers.Nullable[string] `json:"firstname" bson:"firstname"`
	LastName           helpers.Nullable[string] `json:"lastname" bson:"lastname"`
	ResidentialAddress helpers.Nullable[string] `json:"residential_address" bson:"residential_address"`
	TaxCode            helpers.Nullable[string] `json:"tax_code" bson:"tax_code"`
}
