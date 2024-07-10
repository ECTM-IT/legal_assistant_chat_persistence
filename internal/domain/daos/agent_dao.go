package daos

import (
	"context"
	"errors"

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
}

// NewAgentDAO creates a new AgentDAO
func NewAgentDAO(db *mongo.Database) *AgentDAO {
	return &AgentDAO{
		collection: db.Collection("agents"),
	}
}

// GetAllAgents retrieves all agents from the database
func (dao *AgentDAO) GetAllAgents(ctx context.Context) ([]models.Agent, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var agents []models.Agent
	if err := cursor.All(ctx, &agents); err != nil {
		return nil, err
	}

	return agents, nil
}

// GetAgentByID retrieves an agent by its ID from the database
func (dao *AgentDAO) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	var agent models.Agent
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&agent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("agent not found")
		}
		return nil, err
	}
	return &agent, nil
}

// GetAgentsByIDs retrieves agents by their IDs from the database
func (dao *AgentDAO) GetAgentsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Agent, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var agents []models.Agent
	if err := cursor.All(ctx, &agents); err != nil {
		return nil, err
	}

	return agents, nil
}

// CreateAgent creates a new agent in the database
func (dao *AgentDAO) CreateAgent(ctx context.Context, agent *models.Agent) (*mongo.InsertOneResult, error) {
	result, err := dao.collection.InsertOne(ctx, agent)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateAgent updates an existing agent in the database
func (dao *AgentDAO) UpdateAgent(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteAgent deletes an agent by its ID from the database
func (dao *AgentDAO) DeleteAgent(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return result, nil
}
