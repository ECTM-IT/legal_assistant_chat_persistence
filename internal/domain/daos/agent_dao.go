package daos

import (
	"context"
	"errors"
	"fmt"

	logs "github.com/ECTM-IT/legal_assistant_chat_persistence"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AgentDAOInterface defines the interface for the AgentDAO
type AgentDAOInterface interface {
	GetAllAgents(ctx context.Context) ([]models.Agent, error)
	GetAgentByID(ctx context.Context, id primitive.ObjectID) (*models.Agent, error)
	GetAgentsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Agent, error)
	CreateAgent(ctx context.Context, agent *models.Agent) (*mongo.InsertOneResult, error)
	UpdateAgent(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error)
	DeleteAgent(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
}

// AgentDAO implements the AgentDAOInterface
type AgentDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewAgentDAO creates a new AgentDAO
func NewAgentDAO(db *mongo.Database, logger logs.Logger) *AgentDAO {
	return &AgentDAO{
		collection: db.Collection("agents"),
		logger:     logger,
	}
}

// GetAllAgents retrieves all agents from the database
func (dao *AgentDAO) GetAllAgents(ctx context.Context) ([]models.Agent, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve all agents")
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve agents", err)
		return nil, fmt.Errorf("failed to retrieve agents: %w", err)
	}
	defer cursor.Close(ctx)

	var agents []models.Agent
	if err := cursor.All(ctx, &agents); err != nil {
		dao.logger.Error("DAO Level: Failed to decode agents", err)
		return nil, fmt.Errorf("failed to decode agents: %w", err)
	}

	dao.logger.Info("DAO Level: Successfully retrieved all agents")
	return agents, nil
}

// GetAgentByID retrieves an agent by its ID from the database
func (dao *AgentDAO) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve agent by ID")
	var agent models.Agent
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&agent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("Agent not found")
			return nil, errors.New("agent not found")
		}
		dao.logger.Error("DAO Level: Failed to retrieve agent", err)
		return nil, fmt.Errorf("failed to retrieve agent: %w", err)
	}
	dao.logger.Info("DAO Level: Successfully retrieved agent")
	return &agent, nil
}

// GetAgentsByIDs retrieves agents by their IDs from the database
func (dao *AgentDAO) GetAgentsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Agent, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve agents by IDs")
	cursor, err := dao.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve agents by IDs", err)
		return nil, fmt.Errorf("failed to retrieve agents by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var agents []models.Agent
	if err := cursor.All(ctx, &agents); err != nil {
		dao.logger.Error("DAO Level: Failed to decode agents", err)
		return nil, fmt.Errorf("failed to decode agents: %w", err)
	}

	dao.logger.Info("DAO Level: Successfully retrieved agents by IDs")
	return agents, nil
}

// CreateAgent creates a new agent in the database
func (dao *AgentDAO) CreateAgent(ctx context.Context, agent *models.Agent) (*mongo.InsertOneResult, error) {
	dao.logger.Info("DAO Level: Attempting to create new agent")
	result, err := dao.collection.InsertOne(ctx, agent)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to create agent", err)
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}
	dao.logger.Info("DAO Level: Successfully created new agent")
	return result, nil
}

// UpdateAgent updates an existing agent in the database
func (dao *AgentDAO) UpdateAgent(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to update agent")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to update agent", err)
		return nil, fmt.Errorf("failed to update agent: %w", err)
	}
	dao.logger.Info("DAO Level: Successfully updated agent")
	return result, nil
}

// DeleteAgent deletes an agent by its ID from the database
func (dao *AgentDAO) DeleteAgent(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	dao.logger.Info("DAO Level: Attempting to delete agent")
	result, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to delete agent", err)
		return nil, fmt.Errorf("failed to delete agent: %w", err)
	}
	dao.logger.Info("DAO Level: Successfully deleted agent")
	return result, nil
}
