package mappers

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
)

type PlanConversionService interface {
	PlanToDTO(plan *models.Plan) *dtos.PlanResponse
	PlansToDTO(plans []models.Plan) []dtos.PlanResponse
}

type PlanConversionServiceImpl struct {
	logger logs.Logger
}

func NewPlanConversionService(logger logs.Logger) *PlanConversionServiceImpl {
	return &PlanConversionServiceImpl{
		logger: logger,
	}
}

func (s *PlanConversionServiceImpl) PlanToDTO(plan *models.Plan) *dtos.PlanResponse {
	s.logger.Info("Converting Plan to DTO")
	if plan == nil {
		s.logger.Warn("Attempted to convert nil Plan to DTO")
		return nil
	}

	dto := &dtos.PlanResponse{
		Name:        helpers.NewNullable(plan.Name),
		Type:        helpers.NewNullable(plan.Type),
		Price:       helpers.NewNullable(plan.Price),
		Description: helpers.NewNullable(plan.Description),
		Features:    helpers.NewNullable(plan.Features),
	}

	s.logger.Info("Successfully converted Plan to DTO")
	return dto
}

func (s *PlanConversionServiceImpl) PlansToDTO(plans []models.Plan) []dtos.PlanResponse {
	s.logger.Info("Converting multiple Plans to DTOs")
	responseList := make([]dtos.PlanResponse, len(plans))
	for i, plan := range plans {
		responseList[i] = *s.PlanToDTO(&plan)
	}
	s.logger.Info("Successfully converted multiple Plans to DTOs")
	return responseList
}
