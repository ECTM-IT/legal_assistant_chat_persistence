package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/fatih/structs"
)

type AgentRepository struct {
	agentDAO *daos.AgentDAO
	userDAO  *daos.UserDAO
}

// NewAgentRepository initializes a new AgentRepository with the given DAOs.
func NewAgentRepository(agentDAO *daos.AgentDAO, userDAO *daos.UserDAO) *AgentRepository {
	return &AgentRepository{
		agentDAO: agentDAO,
		userDAO:  userDAO,
	}
}

// GetAllAgents retrieves all agents from the database.
// It returns a slice of Agent models and an error if the operation fails.
func (r *AgentRepository) GetAllAgents(ctx context.Context) ([]models.Agent, error) {
	agents, err := r.agentDAO.GetAllAgents(ctx)
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// GetAgentByID retrieves an agent by its ID from the database.
// It returns the Agent model and an error if the operation fails or the agent is not found.
func (r *AgentRepository) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	agent, err := r.agentDAO.GetAgentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// GetAgentsByUserID retrieves all agents associated with a specific user ID.
// It returns a slice of Agent models and an error if the operation fails.
func (r *AgentRepository) GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.Agent, error) {
	user, err := r.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(user.AgentIDs) == 0 {
		return []models.Agent{}, nil
	}

	agents, err := r.agentDAO.GetAgentsByIDs(ctx, user.AgentIDs)
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// PurchaseAgent associates an agent with a user by adding the agent's ID to the user's AgentIDs.
// It returns the updated User model and an error if the operation fails or the agent is already added.
func (r *AgentRepository) PurchaseAgent(ctx context.Context, userID primitive.ObjectID, agentID primitive.ObjectID) (*models.User, error) {
	user, err := r.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if contains(user.AgentIDs, agentID) {
		return nil, errors.New("agent already added to the user")
	}

	_, err = r.agentDAO.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	user.AgentIDs = append(user.AgentIDs, agentID)
	userMap := structs.Map(user)

	_, err = r.userDAO.UpdateUser(ctx, userID, userMap)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// contains checks if a slice of ObjectIDs contains a specific ObjectID.
// It returns true if the element is found, otherwise false.
func contains(slice []primitive.ObjectID, element primitive.ObjectID) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
