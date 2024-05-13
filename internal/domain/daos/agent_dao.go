package daos

import (
	"context"
	"errors"

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

func (dao *AgentDAO) GetAgentsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Agent, error) { //todo: fix
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

func (dao *AgentDAO) CreateAgent(ctx context.Context, agent *models.Agent) error {
	_, err := dao.collection.InsertOne(ctx, agent)
	return err
}

func (dao *AgentDAO) UpdateAgent(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	return dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
}

func (dao *AgentDAO) DeleteAgent(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	return dao.collection.DeleteOne(ctx, bson.M{"_id": id})
}
