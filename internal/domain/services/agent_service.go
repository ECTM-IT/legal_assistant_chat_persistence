package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
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
	agentRepo *repositories.AgentRepository
}

// NewAgentService creates a new AgentService.
func NewAgentService(agentRepo *repositories.AgentRepository) *AgentServiceImpl {
	return &AgentServiceImpl{
		agentRepo: agentRepo,
	}
}

// GetAllAgents retrieves all agents.
func (s *AgentServiceImpl) GetAllAgents(ctx context.Context) ([]dtos.AgentResponse, error) {
	return s.agentRepo.GetAllAgents(ctx)
}

// GetAgentByID retrieves an agent by its ID.
func (s *AgentServiceImpl) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*dtos.AgentResponse, error) {
	return s.agentRepo.GetAgentByID(ctx, id)
}

// GetAgentsByUserID retrieves agents by the user ID.
func (s *AgentServiceImpl) GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]dtos.AgentResponse, error) {
	return s.agentRepo.GetAgentsByUserID(ctx, userID)
}

// PurchaseAgent allows a user to purchase an agent.
func (s *AgentServiceImpl) PurchaseAgent(ctx context.Context, userID, agentID primitive.ObjectID) (*dtos.UserResponse, error) {
	return s.agentRepo.PurchaseAgent(ctx, userID, agentID)
}
