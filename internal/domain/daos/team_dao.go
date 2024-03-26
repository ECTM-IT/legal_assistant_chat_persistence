package daos

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamDAO struct {
	collection *mongo.Collection
}

func NewTeamDAO(db *mongo.Database) *TeamDAO {
	return &TeamDAO{
		collection: db.Collection("teams"),
	}
}

func (dao *TeamDAO) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*models.Team, error) {
	var team models.Team
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (dao *TeamDAO) UpdateTeam(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (dao *TeamDAO) AddMember(ctx context.Context, id primitive.ObjectID, member models.TeamMember) error {
	_, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$push": bson.M{"members": member}})
	return err
}

func (dao *TeamDAO) RemoveMember(ctx context.Context, id primitive.ObjectID, memberID primitive.ObjectID) error {
	_, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"members": bson.M{"_id": memberID}}})
	return err
}
