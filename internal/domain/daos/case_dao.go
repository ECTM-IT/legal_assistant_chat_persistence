package daos

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
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

func (d *CaseDAO) FindAll() ([]dtos.CaseResponse, error) {
	ctx := context.Background()
	var cases []dtos.CaseResponse
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var caseResponse dtos.CaseResponse
		if err := cursor.Decode(&caseResponse); err != nil {
			return nil, err
		}
		cases = append(cases, caseResponse)
	}
	return cases, nil
}

func (d *CaseDAO) FindByID(id primitive.ObjectID) (dtos.CaseResponse, error) {
	ctx := context.Background()
	var caseResponse dtos.CaseResponse
	err := d.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&caseResponse)
	if err != nil {
		return dtos.CaseResponse{}, err
	}
	return caseResponse, nil
}

func (d *CaseDAO) FindByCreatorID(creatorID primitive.ObjectID) ([]dtos.CaseResponse, error) {
	ctx := context.Background()
	var cases []dtos.CaseResponse
	cursor, err := d.collection.Find(ctx, bson.M{"creator_id": creatorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var caseResponse dtos.CaseResponse
		if err := cursor.Decode(&caseResponse); err != nil {
			return nil, err
		}
		cases = append(cases, caseResponse)
	}
	return cases, nil
}

func (d *CaseDAO) Create(caseRequest dtos.CreateCaseRequest) error {
	ctx := context.Background()
	_, err := d.collection.InsertOne(ctx, caseRequest)
	return err
}

func (d *CaseDAO) Update(id primitive.ObjectID, updates dtos.UpdateCaseRequest) error {
	ctx := context.Background()
	_, err := d.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}

func (d *CaseDAO) Delete(id primitive.ObjectID) error {
	ctx := context.Background()
	_, err := d.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (d *CaseDAO) AddCollaborator(caseID, collaboratorID primitive.ObjectID) error {
	ctx := context.Background()
	_, err := d.collection.UpdateOne(ctx, bson.M{"_id": caseID}, bson.M{"$addToSet": bson.M{"collaborator_ids": collaboratorID}})
	return err
}

func (d *CaseDAO) RemoveCollaborator(caseID, collaboratorID primitive.ObjectID) error {
	ctx := context.Background()
	_, err := d.collection.UpdateOne(ctx, bson.M{"_id": caseID}, bson.M{"$pull": bson.M{"collaborator_ids": collaboratorID}})
	return err
}
