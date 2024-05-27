package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateSubscriptionRequest represents a request to create a new subscription.
type CreateSubscriptionRequest struct {
	Plan                helpers.Nullable[string]                 `json:"plan"`
	Expiry              helpers.Nullable[time.Time]              `json:"expiry"`
	Type                helpers.Nullable[string]                 `json:"type"`
	BillingInformations helpers.Nullable[map[string]interface{}] `json:"billing_informations"`
}

// UpdateSubscriptionRequest represents a request to update an existing subscription.
type UpdateSubscriptionRequest struct {
	Plan                helpers.Nullable[string]                 `json:"plan"`
	Expiry              helpers.Nullable[time.Time]              `json:"expiry"`
	Type                helpers.Nullable[string]                 `json:"type"`
	BillingInformations helpers.Nullable[map[string]interface{}] `json:"billing_informations"`
}

// SubscriptionResponse represents a subscription response.
type SubscriptionResponse struct {
	ID                  helpers.Nullable[primitive.ObjectID]     `json:"id"`
	Plan                helpers.Nullable[string]                 `json:"plan"`
	Expiry              helpers.Nullable[time.Time]              `json:"expiry"`
	Type                helpers.Nullable[string]                 `json:"type"`
	BillingInformations helpers.Nullable[map[string]interface{}] `json:"billing_informations"`
}

// CommonBillingInfo contains fields common to multiple billing information types.
type CommonBillingInfo struct {
	VatNumber  helpers.Nullable[string] `json:"vat_number"`
	SdiCode    helpers.Nullable[string] `json:"sdi_code"`
	PecAddress helpers.Nullable[string] `json:"pec_address"`
}

// PersonalBillingInfo contains fields common to individuals (e.g., freelancers, physical persons).
type PersonalBillingInfo struct {
	CommonBillingInfo
	FirstName helpers.Nullable[string] `json:"firstname"`
	LastName  helpers.Nullable[string] `json:"lastname"`
}

// BillingInformationsFreelancer represents billing information for a freelancer.
type BillingInformationsFreelancer struct {
	PersonalBillingInfo
	ProfessionalAddress helpers.Nullable[string] `json:"professional_address"`
	TaxCode             helpers.Nullable[string] `json:"tax_code"`
}

// BillingInformationsIndividualEnterprise represents billing information for an individual enterprise.
type BillingInformationsIndividualEnterprise struct {
	PersonalBillingInfo
	CompanyAddress helpers.Nullable[string] `json:"company_address"`
	HolderTaxCode  helpers.Nullable[string] `json:"holder_tax_code"`
}

// CompanyBillingInfo contains fields common to company billing information.
type CompanyBillingInfo struct {
	CommonBillingInfo
	CompanyName    helpers.Nullable[string] `json:"company_name"`
	LegalAddress   helpers.Nullable[string] `json:"legal_address"`
	CompanyTaxCode helpers.Nullable[string] `json:"company_tax_code"`
}

// BillingInformationsCompany represents billing information for a company.
type BillingInformationsCompany struct {
	CompanyBillingInfo
}

// ProfessionalAssociationBillingInfo contains fields for professional associations.
type ProfessionalAssociationBillingInfo struct {
	CommonBillingInfo
	AssociationName helpers.Nullable[string] `json:"association_name"`
	Address         helpers.Nullable[string] `json:"address"`
	TaxCode         helpers.Nullable[string] `json:"tax_code"`
}

// BillingInformationsProfessionalAssociation represents billing information for a professional association.
type BillingInformationsProfessionalAssociation struct {
	ProfessionalAssociationBillingInfo
}

// BillingInformationsPhysicalPerson represents billing information for a physical person.
type BillingInformationsPhysicalPerson struct {
	PersonalBillingInfo
	ResidentialAddress helpers.Nullable[string] `json:"residential_address"`
	TaxCode            helpers.Nullable[string] `json:"tax_code"`
}
