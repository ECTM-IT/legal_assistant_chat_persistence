package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentService struct {
	agentRepo *repositories.AgentRepository
}

func NewAgentService(agentRepo *repositories.AgentRepository) *AgentService {
	return &AgentService{
		agentRepo: agentRepo,
	}
}

func (s *AgentService) GetAllAgents(ctx context.Context) ([]dtos.AgentResponse, error) {
	return s.agentRepo.GetAllAgents(ctx)
}

func (s *AgentService) GetAgentByID(ctx context.Context, id string) (*dtos.AgentResponse, error) {
	return s.agentRepo.GetAgentByID(ctx, id)
}

func (s *AgentService) GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]dtos.AgentResponse, error) {
	return s.agentRepo.GetAgentsByUserID(ctx, userID)
}

func (s *AgentService) PurchaseAgent(ctx context.Context, userID primitive.ObjectID, agentID primitive.ObjectID) (*dtos.UserResponse, error) {
	return s.agentRepo.PurchaseAgent(ctx, userID, agentID)
}
