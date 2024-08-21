package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// GetTeamByID retrieves a team by its ID from the database
func (dao *TeamDAO) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*models.Team, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve team by ID")
	var team models.Team
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("Team not found")
			return nil, errors.New("team not found")
		}
		dao.logger.Error("DAO Level: Failed to retrieve team", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved team")
	return &team, nil
}

// GetAllTeams retrieves all teams from the database
func (dao *TeamDAO) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve all teams")
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve teams", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var teams []models.Team
	if err := cursor.All(ctx, &teams); err != nil {
		dao.logger.Error("DAO Level: Failed to decode teams", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved all teams")
	return teams, nil
}

// CreateTeam creates a new team in the database
func (dao *TeamDAO) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	dao.logger.Info("DAO Level: Attempting to create new team")
	_, err := dao.collection.InsertOne(ctx, team)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to create team", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully created new team")
	return team, nil
}

// UpdateTeam updates an existing team in the database
func (dao *TeamDAO) UpdateTeam(ctx context.Context, id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to update team")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to update team", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully updated team")
	return result, nil
}

// DeleteTeam deletes a team by its ID from the database
func (dao *TeamDAO) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	dao.logger.Info("DAO Level: Attempting to delete team")
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to delete team", err)
		return err
	}
	dao.logger.Info("DAO Level: Successfully deleted team")
	return nil
}

// AddMember adds a member to a team in the database
func (dao *TeamDAO) AddMember(ctx context.Context, id primitive.ObjectID, member models.TeamMember) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to add member to team")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$push": bson.M{"members": member}})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to add member to team", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully added member to team")
	return result, nil
}

// RemoveMember removes a member from a team in the database
func (dao *TeamDAO) RemoveMember(ctx context.Context, id primitive.ObjectID, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to remove member from team")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"members": bson.M{"_id": memberID}}})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to remove member from team", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully removed member from team")
	return result, nil
}
