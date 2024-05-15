package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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

// GetAllAgents retrieves all agents
func (dao *AgentDAO) GetAllAgents(ctx context.Context) ([]models.Agent, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("Error retrieving all agents", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var agents []models.Agent
	if err := cursor.All(ctx, &agents); err != nil {
		dao.logger.Error("Error decoding agents", err)
		return nil, err
	}

	return agents, nil
}

// GetAgentByID retrieves an agent by its ID
func (dao *AgentDAO) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	var agent models.Agent
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&agent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Error("Agent not found", err, zap.String("agentID", id.Hex()))
			return nil, errors.New("agent not found")
		}
		dao.logger.Error("Error retrieving agent by ID", err, zap.String("agentID", id.Hex()))
		return nil, err
	}
	return &agent, nil
}

// GetAgentsByIDs retrieves agents by their IDs
func (dao *AgentDAO) GetAgentsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Agent, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		dao.logger.Error("Error retrieving agents by IDs", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var agents []models.Agent
	if err := cursor.All(ctx, &agents); err != nil {
		dao.logger.Error("Error decoding agents", err)
		return nil, err
	}

	return agents, nil
}

// CreateAgent creates a new agent
func (dao *AgentDAO) CreateAgent(ctx context.Context, agent *models.Agent) (*mongo.InsertOneResult, error) {
	result, err := dao.collection.InsertOne(ctx, agent)
	if err != nil {
		dao.logger.Error("Error creating agent", err)
		return nil, err
	}
	return result, nil
}

// UpdateAgent updates an existing agent
func (dao *AgentDAO) UpdateAgent(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		dao.logger.Error("Error updating agent", err, zap.String("agentID", id.Hex()))
		return nil, err
	}
	return result, nil
}

// DeleteAgent deletes an agent by its ID
func (dao *AgentDAO) DeleteAgent(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("Error deleting agent", err, zap.String("agentID", id.Hex()))
		return nil, err
	}
	return result, nil
}
