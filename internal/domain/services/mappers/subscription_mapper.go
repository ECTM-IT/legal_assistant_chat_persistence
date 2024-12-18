package mappers

import (
	"errors"
	"fmt"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionConversionService interface {
	SubscriptionToDTO(subscription *models.Subscriptions) *dtos.SubscriptionResponse
	SubscriptionsToDTO(subscriptions []models.Subscriptions) []dtos.SubscriptionResponse
	DTOToSubscription(req *dtos.CreateSubscriptionRequest) (*models.Subscriptions, error)
	UpdateSubscriptionFieldsToMap(updateRequest dtos.UpdateSubscriptionRequest) map[string]interface{}
	BillingInfoToDTO(billingInfo map[string]interface{}) dtos.BillingInformation
	DTOToBillingInfo(billingInfoDTO dtos.BillingInformation) (map[string]interface{}, error)
}

type SubscriptionConversionServiceImpl struct {
	logger logs.Logger
}

func NewSubscriptionConversionService(logger logs.Logger) *SubscriptionConversionServiceImpl {
	return &SubscriptionConversionServiceImpl{
		logger: logger,
	}
}

func (s *SubscriptionConversionServiceImpl) SubscriptionToDTO(subscription *models.Subscriptions) *dtos.SubscriptionResponse {
	s.logger.Info("Converting Subscription to DTO")
	if subscription == nil {
		s.logger.Warn("Attempted to convert nil Subscription to DTO")
		return nil
	}

	dto := &dtos.SubscriptionResponse{
		ID:                   helpers.NewNullable(subscription.ID),
		UserID:               helpers.NewNullable(subscription.UserID),
		Plan:                 helpers.NewNullable(subscription.Plan.Name),
		Renewal:              helpers.NewNullable(dtos.PlanType(subscription.Plan.Type)),
		Status:               helpers.NewNullable(dtos.SubscriptionStatus(subscription.Status)),
		CurrentPeriodStart:   helpers.NewNullable(subscription.CurrentPeriodStart),
		CurrentPeriodEnd:     helpers.NewNullable(subscription.CurrentPeriodEnd),
		CancelAtPeriodEnd:    helpers.NewNullable(subscription.CancelAtPeriodEnd),
		StripeCustomerID:     helpers.NewNullable(subscription.StripeCustomerID),
		StripeSubscriptionID: helpers.NewNullable(subscription.StripeSubscriptionID),
		BillingInformations:  helpers.NewNullable(s.BillingInfoToDTO(subscription.BillingInformations)),
	}
	s.logger.Info("Successfully converted Subscription to DTO")
	return dto
}

func (s *SubscriptionConversionServiceImpl) SubscriptionsToDTO(subscriptions []models.Subscriptions) []dtos.SubscriptionResponse {
	s.logger.Info("Converting multiple Subscriptions to DTOs")
	responseList := make([]dtos.SubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responseList[i] = *s.SubscriptionToDTO(&subscription)
	}
	s.logger.Info("Successfully converted multiple Subscriptions to DTOs")
	return responseList
}

func (s *SubscriptionConversionServiceImpl) DTOToSubscription(req *dtos.CreateSubscriptionRequest) (*models.Subscriptions, error) {
	s.logger.Info("Converting DTO to Subscription")
	if req == nil {
		s.logger.Error("Failed to convert DTO to Subscription: subscription request cannot be nil", errors.New("subscription request cannot be nil"))
		return nil, errors.New("subscription request cannot be nil")
	}

	if !req.UserID.Present || !req.Plan.Present || !req.Type.Present {
		s.logger.Error("Failed to convert DTO to Subscription: userID, plan, and type are required fields", errors.New("userID, plan, and type are required fields"))
		return nil, errors.New("userID, plan, and type are required fields")
	}

	billingInfo, err := s.DTOToBillingInfo(req.BillingInformations.Value)
	if err != nil {
		s.logger.Error("Failed to convert billing information", err)
		return nil, err
	}

	subscription := &models.Subscriptions{
		ID:     primitive.NewObjectID(),
		UserID: req.UserID.Value,
		Plan: models.Plan{
			Name: req.Plan.Value,
			Type: string(req.Type.Value),
		},
		Status:              string(dtos.Incomplete), // Set initial status
		BillingInformations: billingInfo,
		// Other fields will be set when the Stripe subscription is created
	}
	s.logger.Info("Successfully converted DTO to Subscription")
	return subscription, nil
}

func (s *SubscriptionConversionServiceImpl) UpdateSubscriptionFieldsToMap(updateRequest dtos.UpdateSubscriptionRequest) map[string]interface{} {
	s.logger.Info("Converting UpdateSubscriptionRequest to map")
	updateFields := make(map[string]interface{})

	if updateRequest.Plan.Present {
		updateFields["plan"] = updateRequest.Plan.Value
	}
	if updateRequest.Type.Present {
		updateFields["type"] = string(updateRequest.Type.Value)
	}
	if updateRequest.BillingInformations.Present {
		billingInfo, err := s.DTOToBillingInfo(updateRequest.BillingInformations.Value)
		if err != nil {
			s.logger.Error("Failed to convert billing information for update", err)
		} else {
			updateFields["billing_informations"] = billingInfo
		}
	}
	if updateRequest.PaymentMethodID.Present {
		updateFields["payment_method_id"] = updateRequest.PaymentMethodID.Value
	}

	s.logger.Info("Successfully converted UpdateSubscriptionRequest to map")
	return updateFields
}

func (s *SubscriptionConversionServiceImpl) BillingInfoToDTO(billingInfo map[string]interface{}) dtos.BillingInformation {
	s.logger.Info("Converting BillingInfo to DTO")

	billingType, ok := billingInfo["type"].(string)
	if !ok {
		s.logger.Error("Failed to convert BillingInfo to DTO: missing or invalid type", errors.New("BillingInfo to DTO: missing or invalid type"))
		return dtos.BillingInformation{}
	}

	var info interface{}
	var err error

	switch dtos.BillingType(billingType) {
	case dtos.Freelancer:
		info, err = s.mapToFreelancerInfo(billingInfo)
	case dtos.IndividualEnterprise:
		info, err = s.mapToIndividualEnterpriseInfo(billingInfo)
	case dtos.Company:
		info, err = s.mapToCompanyInfo(billingInfo)
	case dtos.ProfessionalAssociation:
		info, err = s.mapToProfessionalAssociationInfo(billingInfo)
	case dtos.PhysicalPerson:
		info, err = s.mapToPhysicalPersonInfo(billingInfo)
	default:
		s.logger.Error(fmt.Sprintf("Unknown billing type: %s", billingType), errors.New("unknown billing type"))
		return dtos.BillingInformation{}
	}

	if err != nil {
		s.logger.Error("Failed to convert BillingInfo to DTO: ", err)
		return dtos.BillingInformation{}
	}

	s.logger.Info("Successfully converted BillingInfo to DTO")
	return dtos.BillingInformation{
		Type: dtos.BillingType(billingType),
		Info: info,
	}
}

func (s *SubscriptionConversionServiceImpl) DTOToBillingInfo(billingInfoDTO dtos.BillingInformation) (map[string]interface{}, error) {
	s.logger.Info("Converting DTO to BillingInfo")

	result := make(map[string]interface{})
	result["type"] = string(billingInfoDTO.Type)

	var err error
	switch billingInfoDTO.Type {
	case dtos.Freelancer:
		err = s.freelancerInfoToMap(billingInfoDTO.Info.(dtos.FreelancerInfo), result)
	case dtos.IndividualEnterprise:
		err = s.individualEnterpriseInfoToMap(billingInfoDTO.Info.(dtos.IndividualEnterpriseInfo), result)
	case dtos.Company:
		err = s.companyInfoToMap(billingInfoDTO.Info.(dtos.CompanyInfo), result)
	case dtos.ProfessionalAssociation:
		err = s.professionalAssociationInfoToMap(billingInfoDTO.Info.(dtos.ProfessionalAssociationInfo), result)
	case dtos.PhysicalPerson:
		err = s.physicalPersonInfoToMap(billingInfoDTO.Info.(dtos.PhysicalPersonInfo), result)
	default:
		return nil, fmt.Errorf("unknown billing type: %s", billingInfoDTO.Type)
	}

	if err != nil {
		s.logger.Error("Failed to convert DTO to BillingInfo: ", err)
		return nil, err
	}

	s.logger.Info("Successfully converted DTO to BillingInfo")
	return result, nil
}

// Helper methods for BillingInfoToDTO

func (s *SubscriptionConversionServiceImpl) mapToFreelancerInfo(m map[string]interface{}) (dtos.FreelancerInfo, error) {
	var info dtos.FreelancerInfo
	info.CommonBillingInfo = s.mapToCommonBillingInfo(m)
	info.FirstName = helpers.NewNullable(m["firstname"].(string))
	info.LastName = helpers.NewNullable(m["lastname"].(string))
	info.ProfessionalAddress = helpers.NewNullable(m["professional_address"].(string))
	info.TaxCode = helpers.NewNullable(m["tax_code"].(string))
	return info, nil
}

func (s *SubscriptionConversionServiceImpl) mapToIndividualEnterpriseInfo(m map[string]interface{}) (dtos.IndividualEnterpriseInfo, error) {
	var info dtos.IndividualEnterpriseInfo
	info.CommonBillingInfo = s.mapToCommonBillingInfo(m)
	info.FirstName = helpers.NewNullable(m["firstname"].(string))
	info.LastName = helpers.NewNullable(m["lastname"].(string))
	info.CompanyAddress = helpers.NewNullable(m["company_address"].(string))
	info.HolderTaxCode = helpers.NewNullable(m["holder_tax_code"].(string))
	return info, nil
}

func (s *SubscriptionConversionServiceImpl) mapToCompanyInfo(m map[string]interface{}) (dtos.CompanyInfo, error) {
	var info dtos.CompanyInfo
	info.CommonBillingInfo = s.mapToCommonBillingInfo(m)
	info.CompanyName = helpers.NewNullable(m["company_name"].(string))
	info.LegalAddress = helpers.NewNullable(m["legal_address"].(string))
	info.CompanyTaxCode = helpers.NewNullable(m["company_tax_code"].(string))
	return info, nil
}

func (s *SubscriptionConversionServiceImpl) mapToProfessionalAssociationInfo(m map[string]interface{}) (dtos.ProfessionalAssociationInfo, error) {
	var info dtos.ProfessionalAssociationInfo
	info.CommonBillingInfo = s.mapToCommonBillingInfo(m)
	info.AssociationName = helpers.NewNullable(m["association_name"].(string))
	info.Address = helpers.NewNullable(m["address"].(string))
	info.TaxCode = helpers.NewNullable(m["tax_code"].(string))
	return info, nil
}

func (s *SubscriptionConversionServiceImpl) mapToPhysicalPersonInfo(m map[string]interface{}) (dtos.PhysicalPersonInfo, error) {
	var info dtos.PhysicalPersonInfo
	info.CommonBillingInfo = s.mapToCommonBillingInfo(m)
	info.FirstName = helpers.NewNullable(m["firstname"].(string))
	info.LastName = helpers.NewNullable(m["lastname"].(string))
	info.ResidentialAddress = helpers.NewNullable(m["residential_address"].(string))
	info.TaxCode = helpers.NewNullable(m["tax_code"].(string))
	return info, nil
}

func (s *SubscriptionConversionServiceImpl) mapToCommonBillingInfo(m map[string]interface{}) dtos.CommonBillingInfo {
	return dtos.CommonBillingInfo{
		VatNumber:  helpers.NewNullable(m["vat_number"].(string)),
		SdiCode:    helpers.NewNullable(m["sdi_code"].(string)),
		PecAddress: helpers.NewNullable(m["pec_address"].(string)),
	}
}

// Helper methods for DTOToBillingInfo

func (s *SubscriptionConversionServiceImpl) freelancerInfoToMap(info dtos.FreelancerInfo, m map[string]interface{}) error {
	s.commonBillingInfoToMap(info.CommonBillingInfo, m)
	m["firstname"] = info.FirstName.Value
	m["lastname"] = info.LastName.Value
	m["professional_address"] = info.ProfessionalAddress.Value
	m["tax_code"] = info.TaxCode.Value
	return nil
}

func (s *SubscriptionConversionServiceImpl) individualEnterpriseInfoToMap(info dtos.IndividualEnterpriseInfo, m map[string]interface{}) error {
	s.commonBillingInfoToMap(info.CommonBillingInfo, m)
	m["firstname"] = info.FirstName.Value
	m["lastname"] = info.LastName.Value
	m["company_address"] = info.CompanyAddress.Value
	m["holder_tax_code"] = info.HolderTaxCode.Value
	return nil
}

func (s *SubscriptionConversionServiceImpl) companyInfoToMap(info dtos.CompanyInfo, m map[string]interface{}) error {
	s.commonBillingInfoToMap(info.CommonBillingInfo, m)
	m["company_name"] = info.CompanyName.Value
	m["legal_address"] = info.LegalAddress.Value
	m["company_tax_code"] = info.CompanyTaxCode.Value
	return nil
}

func (s *SubscriptionConversionServiceImpl) professionalAssociationInfoToMap(info dtos.ProfessionalAssociationInfo, m map[string]interface{}) error {
	s.commonBillingInfoToMap(info.CommonBillingInfo, m)
	m["association_name"] = info.AssociationName.Value
	m["address"] = info.Address.Value
	m["tax_code"] = info.TaxCode.Value
	return nil
}

func (s *SubscriptionConversionServiceImpl) physicalPersonInfoToMap(info dtos.PhysicalPersonInfo, m map[string]interface{}) error {
	s.commonBillingInfoToMap(info.CommonBillingInfo, m)
	m["firstname"] = info.FirstName.Value
	m["lastname"] = info.LastName.Value
	m["residential_address"] = info.ResidentialAddress.Value
	m["tax_code"] = info.TaxCode.Value
	return nil
}

func (s *SubscriptionConversionServiceImpl) commonBillingInfoToMap(info dtos.CommonBillingInfo, m map[string]interface{}) {
	m["vat_number"] = info.VatNumber.Value
	m["sdi_code"] = info.SdiCode.Value
	m["pec_address"] = info.PecAddress.Value
}
