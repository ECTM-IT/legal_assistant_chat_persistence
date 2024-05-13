package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CaseDAO struct {
	collection *mongo.Collection
}

func NewCaseDAO(db *mongo.Database) *CaseDAO {
	return &CaseDAO{
		collection: db.Collection("cases"),
	}
}

func (d *CaseDAO) FindAll(ctx context.Context) ([]models.Case, error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cases []models.Case
	if err := cursor.All(ctx, &cases); err != nil {
		return nil, err
	}

	return cases, nil
}

func (d *CaseDAO) FindByID(ctx context.Context, id primitive.ObjectID) (models.Case, error) {
	var caseResponse models.Case
	err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&caseResponse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Case{}, errors.New("case not found")
		}
		return models.Case{}, err
	}
	return caseResponse, nil
}

func (d *CaseDAO) FindByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]models.Case, error) {
	cursor, err := d.collection.Find(ctx, bson.M{"creator_id": creatorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cases []models.Case
	if err := cursor.All(ctx, &cases); err != nil {
		return nil, err
	}

	return cases, nil
}

func (d *CaseDAO) Create(ctx context.Context, caseRequest *models.Case) (*mongo.InsertOneResult, error) {
	return d.collection.InsertOne(ctx, caseRequest)
}

func (d *CaseDAO) Update(ctx context.Context, id primitive.ObjectID, updates dtos.UpdateCaseRequest) (*mongo.UpdateResult, error) {
	return d.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
}

func (d *CaseDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := d.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (d *CaseDAO) AddCollaborator(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return d.collection.UpdateOne(ctx, bson.M{"_id": caseID}, bson.M{"$addToSet": bson.M{"collaborator_ids": collaboratorID}})
}

func (d *CaseDAO) RemoveCollaborator(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return d.collection.UpdateOne(ctx, bson.M{"_id": caseID}, bson.M{"$pull": bson.M{"collaborator_ids": collaboratorID}})
}
