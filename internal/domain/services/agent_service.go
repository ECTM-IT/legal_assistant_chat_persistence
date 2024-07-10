package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AgentService defines the interface for agent-related operations.
type AgentService interface {
	GetAllAgents(ctx context.Context) ([]dtos.AgentResponse, error)
	GetAgentByID(ctx context.Context, id string) (*dtos.AgentResponse, error)
	GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]dtos.AgentResponse, error)
	PurchaseAgent(ctx context.Context, userID, agentID primitive.ObjectID) (*dtos.UserResponse, error)
}

// AgentServiceImpl implements the AgentService interface.
type AgentServiceImpl struct {
	agentRepo  *repositories.AgentRepository
	mapper     *mappers.AgentConversionServiceImpl
	userMapper *mappers.UserConversionServiceImpl
}

// NewAgentService creates a new AgentService.
func NewAgentService(agentRepo *repositories.AgentRepository, mapper *mappers.AgentConversionServiceImpl, userMapper *mappers.UserConversionServiceImpl) *AgentServiceImpl {
	return &AgentServiceImpl{
		agentRepo:  agentRepo,
		mapper:     mapper,
		userMapper: userMapper,
	}
}

// GetAllAgents retrieves all agents.
func (s *AgentServiceImpl) GetAllAgents(ctx context.Context) ([]dtos.AgentResponse, error) {
	agents, err := s.agentRepo.GetAllAgents(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.AgentsToDTO(agents), nil
}

// GetAgentByID retrieves an agent by its ID.
func (s *AgentServiceImpl) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*dtos.AgentResponse, error) {
	agents, err := s.agentRepo.GetAgentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapper.AgentToDTO(agents), nil
}

// GetAgentsByUserID retrieves agents by the user ID.
func (s *AgentServiceImpl) GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]dtos.AgentResponse, error) {
	agents, err := s.agentRepo.GetAgentsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.mapper.AgentsToDTO(agents), nil
}

// PurchaseAgent allows a user to purchase an agent.
func (s *AgentServiceImpl) PurchaseAgent(ctx context.Context, userID, agentID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := s.agentRepo.PurchaseAgent(ctx, userID, agentID)
	if err != nil {
		return nil, err
	}
	return s.userMapper.UserToDTO(user), nil
}
