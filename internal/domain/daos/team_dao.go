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

// TeamDAOInterface defines the interface for the TeamDAO
type TeamDAOInterface interface {
	GetTeamByID(ctx context.Context, id primitive.ObjectID) (*models.Team, error)
	GetAllTeams(ctx context.Context) ([]models.Team, error)
	CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	UpdateTeam(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error)
	DeleteTeam(ctx context.Context, id primitive.ObjectID) error
	AddMember(ctx context.Context, id primitive.ObjectID, member models.TeamMember) (*mongo.UpdateResult, error)
	RemoveMember(ctx context.Context, id primitive.ObjectID, memberID primitive.ObjectID) (*mongo.UpdateResult, error)
}

// TeamDAO implements the TeamDAOInterface
type TeamDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewTeamDAO creates a new TeamDAO
func NewTeamDAO(db *mongo.Database, logger logs.Logger) *TeamDAO {
	return &TeamDAO{
		collection: db.Collection("teams"),
		logger:     logger,
	}
}

// GetTeamByID retrieves a team by its ID
func (dao *TeamDAO) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*models.Team, error) {
	var team models.Team
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Error("Team not found", err, zap.String("teamID", id.Hex()))
			return nil, errors.New("team not found")
		}
		dao.logger.Error("Error retrieving team by ID", err, zap.String("teamID", id.Hex()))
		return nil, err
	}
	return &team, nil
}

// GetAllTeams retrieves all teams
func (dao *TeamDAO) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("Error retrieving all teams", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var teams []models.Team
	if err := cursor.All(ctx, &teams); err != nil {
		dao.logger.Error("Error decoding teams", err)
		return nil, err
	}
	return teams, nil
}

// CreateTeam creates a new team
func (dao *TeamDAO) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	_, err := dao.collection.InsertOne(ctx, team)
	if err != nil {
		dao.logger.Error("Error creating team", err)
		return nil, err
	}
	return team, nil
}

// UpdateTeam updates an existing team
func (dao *TeamDAO) UpdateTeam(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		dao.logger.Error("Error updating team", err, zap.String("teamID", id.Hex()))
		return nil, err
	}
	return result, nil
}

// DeleteTeam deletes a team by its ID
func (dao *TeamDAO) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("Error deleting team", err, zap.String("teamID", id.Hex()))
		return err
	}
	return nil
}

// AddMember adds a member to a team
func (dao *TeamDAO) AddMember(ctx context.Context, id primitive.ObjectID, member models.TeamMember) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$push": bson.M{"members": member}})
	if err != nil {
		dao.logger.Error("Error adding member to team", err, zap.String("teamID", id.Hex()))
		return nil, err
	}
	return result, nil
}

// RemoveMember removes a member from a team
func (dao *TeamDAO) RemoveMember(ctx context.Context, id primitive.ObjectID, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"members": bson.M{"_id": memberID}}})
	if err != nil {
		dao.logger.Error("Error removing member from team", err, zap.String("teamID", id.Hex()), zap.String("memberID", memberID.Hex()))
		return nil, err
	}
	return result, nil
}
