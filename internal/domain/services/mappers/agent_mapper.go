package mappers

import (
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentConversionService interface {
	AgentToDTO(agent *models.Agent) *dtos.AgentResponse
	AgentsToDTO(agents []models.Agent) []dtos.AgentResponse
	DTOToAgent(agentDTO dtos.CreateAgentRequest) (*models.Agent, error)
	UpdateAgentFieldsToMap(updateRequest dtos.UpdateAgentRequest) map[string]interface{}
}

type AgentConversionServiceImpl struct {
	logger logs.Logger
}

func NewAgentConversionService(logger logs.Logger) *AgentConversionServiceImpl {
	return &AgentConversionServiceImpl{
		logger: logger,
	}
}

func (s *AgentConversionServiceImpl) AgentToDTO(agent *models.Agent) *dtos.AgentResponse {
	s.logger.Info("Converting Agent to DTO")
	if agent == nil {
		s.logger.Warn("Attempted to convert nil Agent to DTO")
		return nil
	}

	dto := &dtos.AgentResponse{
		ID:           helpers.NewNullable(agent.ID),
		ProfileImage: helpers.NewNullable(agent.ProfileImage),
		Name:         helpers.NewNullable(agent.Name),
		Description:  helpers.NewNullable(agent.Description),
		Skills:       helpers.NewNullable(agent.Skills),
		Price:        helpers.NewNullable(agent.Price),
		Code:         helpers.NewNullable(agent.Code),
	}
	s.logger.Info("Successfully converted Agent to DTO")
	return dto
}

func (s *AgentConversionServiceImpl) AgentsToDTO(agents []models.Agent) []dtos.AgentResponse {
	s.logger.Info("Converting multiple Agents to DTOs")
	agentResponses := make([]dtos.AgentResponse, len(agents))
	for i, agent := range agents {
		agentResponses[i] = *s.AgentToDTO(&agent)
	}
	s.logger.Info("Successfully converted multiple Agents to DTOs")
	return agentResponses
}

func (s *AgentConversionServiceImpl) DTOToAgent(agentDTO dtos.CreateAgentRequest) (*models.Agent, error) {
	s.logger.Info("Converting DTO to Agent")
	if !agentDTO.Name.Present || !agentDTO.Description.Present {
		s.logger.Error("Failed to convert DTO to Agent: name and description are required", errors.New("bad Request"))
		return nil, errors.New("name and description are required")
	}

	agent := &models.Agent{
		ID:           primitive.NewObjectID(),
		ProfileImage: agentDTO.ProfileImage.OrElse(""),
		Name:         agentDTO.Name.Value,
		Description:  agentDTO.Description.Value,
		Skills:       agentDTO.Skills.OrElse([]string{}),
		Price:        agentDTO.Price.OrElse(0),
		Code:         agentDTO.Code.OrElse(""),
	}
	s.logger.Info("Successfully converted DTO to Agent")
	return agent, nil
}

func (s *AgentConversionServiceImpl) UpdateAgentFieldsToMap(updateRequest dtos.UpdateAgentRequest) map[string]interface{} {
	s.logger.Info("Converting UpdateAgentRequest to map")
	updateFields := make(map[string]interface{})

	if updateRequest.ProfileImage.Present {
		updateFields["profile_image"] = updateRequest.ProfileImage.Value
	}
	if updateRequest.Name.Present {
		updateFields["name"] = updateRequest.Name.Value
	}
	if updateRequest.Description.Present {
		updateFields["description"] = updateRequest.Description.Value
	}
	if updateRequest.Skills.Present {
		updateFields["skills"] = updateRequest.Skills.Value
	}
	if updateRequest.Price.Present {
		updateFields["price"] = updateRequest.Price.Value
	}
	if updateRequest.Code.Present {
		updateFields["code"] = updateRequest.Code.Value
	}

	s.logger.Info("Successfully converted UpdateAgentRequest to map")
	return updateFields
}
