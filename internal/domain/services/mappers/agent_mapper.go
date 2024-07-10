package mappers

import (
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentConversionService interface {
	AgentToDTO(agent *models.Agent) *dtos.AgentResponse
	AgentsToDTO(agents []models.Agent) []dtos.AgentResponse
	DTOToAgent(agentDTO dtos.CreateAgentRequest) (*models.Agent, error)
	UpdateAgentFieldsToMap(updateRequest dtos.UpdateAgentRequest) map[string]interface{}
}

type AgentConversionServiceImpl struct{}

func NewAgentConversionService() *AgentConversionServiceImpl {
	return &AgentConversionServiceImpl{}
}

func (s *AgentConversionServiceImpl) AgentToDTO(agent *models.Agent) *dtos.AgentResponse {
	if agent == nil {
		return nil
	}

	return &dtos.AgentResponse{
		ID:           helpers.NewNullable(agent.ID),
		ProfileImage: helpers.NewNullable(agent.ProfileImage),
		Name:         helpers.NewNullable(agent.Name),
		Description:  helpers.NewNullable(agent.Description),
		Skills:       helpers.NewNullable(agent.Skills),
		Price:        helpers.NewNullable(agent.Price),
		Code:         helpers.NewNullable(agent.Code),
	}
}

func (s *AgentConversionServiceImpl) AgentsToDTO(agents []models.Agent) []dtos.AgentResponse {
	agentResponses := make([]dtos.AgentResponse, len(agents))
	for i, agent := range agents {
		agentResponses[i] = *s.AgentToDTO(&agent)
	}
	return agentResponses
}

func (s *AgentConversionServiceImpl) DTOToAgent(agentDTO dtos.CreateAgentRequest) (*models.Agent, error) {
	if !agentDTO.Name.Present || !agentDTO.Description.Present {
		return nil, errors.New("name and description are required")
	}

	return &models.Agent{
		ID:           primitive.NewObjectID(),
		ProfileImage: agentDTO.ProfileImage.OrElse(""),
		Name:         agentDTO.Name.Value,
		Description:  agentDTO.Description.Value,
		Skills:       agentDTO.Skills.OrElse([]string{}),
		Price:        agentDTO.Price.OrElse(0),
		Code:         agentDTO.Code.OrElse(""),
	}, nil
}

func (s *AgentConversionServiceImpl) UpdateAgentFieldsToMap(updateRequest dtos.UpdateAgentRequest) map[string]interface{} {
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

	return updateFields
}
