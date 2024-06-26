package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateSubscriptionRequest struct {
	Plan                helpers.Nullable[string]                 `json:"plan" bson:"plan"`
	Expiry              helpers.Nullable[time.Time]              `json:"expiry" bson:"expiry"`
	Type                helpers.Nullable[string]                 `json:"type" bson:"type"`
	BillingInformations helpers.Nullable[map[string]interface{}] `json:"billing_informations" bson:"billing_informations"`
}

type UpdateSubscriptionRequest struct {
	Plan                helpers.Nullable[string]                 `json:"plan" bson:"plan,omitempty"`
	Expiry              helpers.Nullable[time.Time]              `json:"expiry" bson:"expiry,omitempty"`
	Type                helpers.Nullable[string]                 `json:"type" bson:"type,omitempty"`
	BillingInformations helpers.Nullable[map[string]interface{}] `json:"billing_informations" bson:"billing_informations,omitempty"`
}

type SubscriptionResponse struct {
	ID                  helpers.Nullable[primitive.ObjectID]     `json:"id" bson:"_id,omitempty"`
	Plan                helpers.Nullable[string]                 `json:"plan" bson:"plan"`
	Expiry              helpers.Nullable[time.Time]              `json:"expiry" bson:"expiry"`
	Type                helpers.Nullable[string]                 `json:"type" bson:"type"`
	BillingInformations helpers.Nullable[map[string]interface{}] `json:"billing_informations" bson:"billing_informations"`
}

type CommonBillingInfo struct {
	VatNumber  helpers.Nullable[string] `json:"vat_number" bson:"vat_number"`
	SdiCode    helpers.Nullable[string] `json:"sdi_code" bson:"sdi_code"`
	PecAddress helpers.Nullable[string] `json:"pec_address" bson:"pec_address"`
}

type PersonalBillingInfo struct {
	CommonBillingInfo
	FirstName helpers.Nullable[string] `json:"firstname" bson:"firstname"`
	LastName  helpers.Nullable[string] `json:"lastname" bson:"lastname"`
}

type BillingInformationsFreelancer struct {
	PersonalBillingInfo
	ProfessionalAddress helpers.Nullable[string] `json:"professional_address" bson:"professional_address"`
	TaxCode             helpers.Nullable[string] `json:"tax_code" bson:"tax_code"`
}

type BillingInformationsIndividualEnterprise struct {
	PersonalBillingInfo
	CompanyAddress helpers.Nullable[string] `json:"company_address" bson:"company_address"`
	HolderTaxCode  helpers.Nullable[string] `json:"holder_tax_code" bson:"holder_tax_code"`
}

type CompanyBillingInfo struct {
	CommonBillingInfo
	CompanyName    helpers.Nullable[string] `json:"company_name" bson:"company_name"`
	LegalAddress   helpers.Nullable[string] `json:"legal_address" bson:"legal_address"`
	CompanyTaxCode helpers.Nullable[string] `json:"company_tax_code" bson:"company_tax_code"`
}

type BillingInformationsCompany struct {
	CompanyBillingInfo
}

type ProfessionalAssociationBillingInfo struct {
	CommonBillingInfo
	AssociationName helpers.Nullable[string] `json:"association_name" bson:"association_name"`
	Address         helpers.Nullable[string] `json:"address" bson:"address"`
	TaxCode         helpers.Nullable[string] `json:"tax_code" bson:"tax_code"`
}

type BillingInformationsProfessionalAssociation struct {
	ProfessionalAssociationBillingInfo
}

type BillingInformationsPhysicalPerson struct {
	PersonalBillingInfo
	ResidentialAddress helpers.Nullable[string] `json:"residential_address" bson:"residential_address"`
	TaxCode            helpers.Nullable[string] `json:"tax_code" bson:"tax_code"`
}
