package daos

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
)

type AgentDAO struct {
	collection *mongo.Collection
}

func NewAgentDAO(db *mongo.Database) *AgentDAO {
	return &AgentDAO{
		collection: db.Collection("agents"),
	}
}

func (dao *AgentDAO) GetAllAgents(ctx context.Context) ([]models.Agent, error) {
	var agents []models.Agent
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var agent models.Agent
		if err := cursor.Decode(&agent); err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

func (dao *AgentDAO) GetAgentByID(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	var agent models.Agent
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&agent)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (dao *AgentDAO) GetAgentsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Agent, error) {
	var agents []models.Agent
	cursor, err := dao.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var agent models.Agent
		if err := cursor.Decode(&agent); err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}
