package daos

import (
	"context"
	"errors"

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

// GetTeamByID retrieves a team by its ID
func (dao *TeamDAO) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*models.Team, error) {
	var team models.Team
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("team not found")
		}
		return nil, err
	}
	return &team, nil
}

// GetAllTeams retrieves all teams
func (dao *TeamDAO) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var teams []models.Team
	if err := cursor.All(ctx, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// CreateTeam creates a new team
func (dao *TeamDAO) CreateTeam(ctx context.Context, team *models.Team) error {
	_, err := dao.collection.InsertOne(ctx, team)
	return err
}

// UpdateTeam updates an existing team
func (dao *TeamDAO) UpdateTeam(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	return dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
}

// DeleteTeam deletes a team by its ID
func (dao *TeamDAO) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AddMember adds a member to a team
func (dao *TeamDAO) AddMember(ctx context.Context, id primitive.ObjectID, member models.TeamMember) (*mongo.UpdateResult, error) {
	return dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$push": bson.M{"members": member}})
}

// RemoveMember removes a member from a team
func (dao *TeamDAO) RemoveMember(ctx context.Context, id primitive.ObjectID, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"members": bson.M{"_id": memberID}}})
}
