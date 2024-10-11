package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/fatih/structs"
)

type AgentRepository struct {
	agentDAO *daos.AgentDAO
	userDAO  *daos.UserDAO
	logger   logs.Logger
}

// NewAgentRepository initializes a new AgentRepository with the given DAOs.
func NewAgentRepository(agentDAO *daos.AgentDAO, userDAO *daos.UserDAO, logger logs.Logger) *AgentRepository {
	return &AgentRepository{
		agentDAO: agentDAO,
		userDAO:  userDAO,
		logger:   logger,
	}
}

// GetAllAgents retrieves all agents from the database.
func (r *AgentRepository) GetAllAgents(ctx context.Context) ([]models.Agent, error) {
	r.logger.Info("Repository Level: Attempting to retrieve all agents")
	agents, err := r.agentDAO.GetAllAgents(ctx)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve agents", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved all agents")
	return agents, nil
}

// GetAgentByID retrieves an agent by its ID from the database.
func (r *AgentRepository) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	r.logger.Info("Repository Level: Attempting to retrieve agent by ID")
	agent, err := r.agentDAO.GetAgentByID(ctx, id)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve agent", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved agent")
	return agent, nil
}

// GetAgentsByUserID retrieves all agents associated with a specific user ID.
func (r *AgentRepository) GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.Agent, error) {
	r.logger.Info("Repository Level: Attempting to retrieve agents by user ID")
	user, err := r.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve user", err)
		return nil, err
	}

	if len(user.AgentIDs) == 0 {
		r.logger.Info("User has no agents")
		return []models.Agent{}, nil
	}

	agents, err := r.agentDAO.GetAgentsByIDs(ctx, user.AgentIDs)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve agents", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved agents for user")
	return agents, nil
}

// PurchaseAgent associates an agent with a user by adding the agent's ID to the user's AgentIDs.
func (r *AgentRepository) PurchaseAgent(ctx context.Context, userID primitive.ObjectID, agentID primitive.ObjectID) (*models.User, error) {
	r.logger.Info("Repository Level: Attempting to purchase agent for user")
	user, err := r.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve user", err)
		return nil, err
	}

	if contains(user.AgentIDs, agentID) {
		r.logger.Warn("Agent already added to the user")
		return nil, errors.New("agent already added to the user")
	}

	_, err = r.agentDAO.GetAgentByID(ctx, agentID)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve agent", err)
		return nil, err
	}

	user.AgentIDs = append(user.AgentIDs, agentID)
	userMap := structs.Map(user)

	_, err = r.userDAO.UpdateUser(ctx, userID, userMap)
	if err != nil {
		r.logger.Error("Repository Level: Failed to update user", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully purchased agent for user")
	return user, nil
}

// contains checks if a slice of ObjectIDs contains a specific ObjectID.
func contains(slice []primitive.ObjectID, element primitive.ObjectID) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
