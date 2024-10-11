package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AgentService defines the interface for agent-related operations.
type AgentService interface {
	GetAllAgents(ctx context.Context) ([]dtos.AgentResponse, error)
	GetAgentByID(ctx context.Context, id primitive.ObjectID) (*dtos.AgentResponse, error)
	GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]dtos.AgentResponse, error)
	PurchaseAgent(ctx context.Context, userID, agentID primitive.ObjectID) (*dtos.UserResponse, error)
}

// AgentServiceImpl implements the AgentService interface.
type AgentServiceImpl struct {
	agentRepo  *repositories.AgentRepository
	mapper     *mappers.AgentConversionServiceImpl
	userMapper *mappers.UserConversionServiceImpl
	logger     logs.Logger
}

// NewAgentService creates a new AgentService.
func NewAgentService(agentRepo *repositories.AgentRepository, mapper *mappers.AgentConversionServiceImpl, userMapper *mappers.UserConversionServiceImpl, logger logs.Logger) *AgentServiceImpl {
	return &AgentServiceImpl{
		agentRepo:  agentRepo,
		mapper:     mapper,
		userMapper: userMapper,
		logger:     logger,
	}
}

// GetAllAgents retrieves all agents.
func (s *AgentServiceImpl) GetAllAgents(ctx context.Context) ([]dtos.AgentResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve all agents")
	agents, err := s.agentRepo.GetAllAgents(ctx)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve all agents", err)
		return nil, err
	}
	agentResponses := s.mapper.AgentsToDTO(agents)
	s.logger.Info("Service Level: Successfully retrieved all agents")
	return agentResponses, nil
}

// GetAgentByID retrieves an agent by its ID.
func (s *AgentServiceImpl) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*dtos.AgentResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve agent by ID")
	agent, err := s.agentRepo.GetAgentByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve agent by ID", err)
		return nil, err
	}
	agentResponse := s.mapper.AgentToDTO(agent)
	s.logger.Info("Service Level: Successfully retrieved agent by ID")
	return agentResponse, nil
}

// GetAgentsByUserID retrieves agents by the user ID.
func (s *AgentServiceImpl) GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]dtos.AgentResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve agents by user ID")
	agents, err := s.agentRepo.GetAgentsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve agents by user ID", err)
		return nil, err
	}
	agentResponses := s.mapper.AgentsToDTO(agents)
	s.logger.Info("Service Level: Successfully retrieved agents by user ID")
	return agentResponses, nil
}

// PurchaseAgent allows a user to purchase an agent.
func (s *AgentServiceImpl) PurchaseAgent(ctx context.Context, userID, agentID primitive.ObjectID) (*dtos.UserResponse, error) {
	s.logger.Info("Service Level: Attempting to purchase agent for user")
	user, err := s.agentRepo.PurchaseAgent(ctx, userID, agentID)
	if err != nil {
		s.logger.Error("Service Level: Failed to purchase agent for user", err)
		return nil, err
	}
	userResponse := s.userMapper.UserToDTO(user)
	s.logger.Info("Service Level: Successfully purchased agent for user")
	return userResponse, nil
}
