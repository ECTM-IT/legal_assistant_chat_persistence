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

// CaseDAOInterface defines the interface for the CaseDAO
type CaseDAOInterface interface {
	FindAll(ctx context.Context) ([]models.Case, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (models.Case, error)
	FindByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]models.Case, error)
	Create(ctx context.Context, caseRequest *models.Case) (*mongo.InsertOneResult, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddCollaborator(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error)
	RemoveCollaborator(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error)
}

// CaseDAO implements the CaseDAOInterface
type CaseDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewCaseDAO creates a new CaseDAO
func NewCaseDAO(db *mongo.Database, logger logs.Logger) *CaseDAO {
	return &CaseDAO{
		collection: db.Collection("cases"),
		logger:     logger,
	}
}

// FindAll retrieves all cases
func (dao *CaseDAO) FindAll(ctx context.Context) ([]models.Case, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("Error retrieving all cases", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var cases []models.Case
	if err := cursor.All(ctx, &cases); err != nil {
		dao.logger.Error("Error decoding cases", err)
		return nil, err
	}

	return cases, nil
}

// FindByID retrieves a case by its ID
func (dao *CaseDAO) FindByID(ctx context.Context, id primitive.ObjectID) (models.Case, error) {
	var caseResponse models.Case
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&caseResponse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Error("Case not found", err, zap.String("caseID", id.Hex()))
			return models.Case{}, errors.New("case not found")
		}
		dao.logger.Error("Error retrieving case by ID", err, zap.String("caseID", id.Hex()))
		return models.Case{}, err
	}
	return caseResponse, nil
}

// FindByCreatorID retrieves cases by creator ID
func (dao *CaseDAO) FindByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]models.Case, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{"creator_id": creatorID})
	if err != nil {
		dao.logger.Error("Error retrieving cases by creator ID", err, zap.String("creatorID", creatorID.Hex()))
		return nil, err
	}
	defer cursor.Close(ctx)

	var cases []models.Case
	if err := cursor.All(ctx, &cases); err != nil {
		dao.logger.Error("Error decoding cases", err)
		return nil, err
	}

	return cases, nil
}

// Create creates a new case
func (dao *CaseDAO) Create(ctx context.Context, caseRequest *models.Case) (*mongo.InsertOneResult, error) {
	result, err := dao.collection.InsertOne(ctx, caseRequest)
	if err != nil {
		dao.logger.Error("Error creating case", err)
		return nil, err
	}
	return result, nil
}

// Update updates an existing case
func (dao *CaseDAO) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	if err != nil {
		dao.logger.Error("Error updating case", err, zap.String("caseID", id.Hex()))
		return nil, err
	}
	return result, nil
}

// Delete deletes a case by its ID
func (dao *CaseDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("Error deleting case", err, zap.String("caseID", id.Hex()))
		return err
	}
	return nil
}

// AddCollaborator adds a collaborator to a case
func (dao *CaseDAO) AddCollaborator(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": caseID}, bson.M{"$addToSet": bson.M{"collaborator_ids": collaboratorID}})
	if err != nil {
		dao.logger.Error("Error adding collaborator to case", err, zap.String("caseID", caseID.Hex()), zap.String("collaboratorID", collaboratorID.Hex()))
		return nil, err
	}
	return result, nil
}

// RemoveCollaborator removes a collaborator from a case
func (dao *CaseDAO) RemoveCollaborator(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": caseID}, bson.M{"$pull": bson.M{"collaborator_ids": collaboratorID}})
	if err != nil {
		dao.logger.Error("Error removing collaborator from case", err, zap.String("caseID", caseID.Hex()), zap.String("collaboratorID", collaboratorID.Hex()))
		return nil, err
	}
	return result, nil
}
