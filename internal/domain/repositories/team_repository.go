package repositories

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamRepository struct {
	teamDAO *daos.TeamDAO
	userDAO *daos.UserDAO
	logger  logs.Logger
}

func NewTeamRepository(teamDAO *daos.TeamDAO, userDAO *daos.UserDAO, logger logs.Logger) *TeamRepository {
	return &TeamRepository{
		teamDAO: teamDAO,
		userDAO: userDAO,
		logger:  logger,
	}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	r.logger.Info("Repository Level: Attempting to create new team")
	createdTeam, err := r.teamDAO.CreateTeam(ctx, team)
	if err != nil {
		r.logger.Error("Repository Level: Failed to create team", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully created new team")
	return createdTeam, nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*models.Team, error) {
	r.logger.Info("Repository Level: Attempting to retrieve team by ID")
	team, err := r.teamDAO.GetTeamByID(ctx, id)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve team", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved team")
	return team, nil
}

func (r *TeamRepository) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	r.logger.Info("Repository Level: Attempting to retrieve all teams")
	teams, err := r.teamDAO.GetAllTeams(ctx)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve teams", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved all teams")
	return teams, nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Team, error) {
	r.logger.Info("Repository Level: Attempting to update team")
	_, err := r.teamDAO.UpdateTeam(ctx, id, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to update team", err)
		return nil, err
	}
	updatedTeam, err := r.teamDAO.GetTeamByID(ctx, id)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve updated team", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully updated team")
	return updatedTeam, nil
}

func (r *TeamRepository) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	r.logger.Info("Repository Level: Attempting to delete team")
	err := r.teamDAO.DeleteTeam(ctx, id)
	if err != nil {
		r.logger.Error("Repository Level: Failed to delete team", err)
		return err
	}
	r.logger.Info("Repository Level: Successfully deleted team")
	return nil
}

func (r *TeamRepository) GetTeamMember(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	r.logger.Info("Repository Level: Attempting to retrieve team member")
	user, err := r.userDAO.GetUserByID(ctx, id)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve team member", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved team member")
	return user, nil
}

func (r *TeamRepository) ChangeAdmin(ctx context.Context, id primitive.ObjectID, email string) (*models.User, error) {
	r.logger.Info("Repository Level: Attempting to change team admin")
	user, err := r.userDAO.GetUserByEmail(ctx, email)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve user by email", err)
		return nil, err
	}
	update := bson.M{
		"admin_id": user.ID,
	}
	_, err = r.teamDAO.UpdateTeam(ctx, id, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to update team admin", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully changed team admin")
	return user, nil
}

func (r *TeamRepository) AddMember(ctx context.Context, id primitive.ObjectID, email string) (*mongo.UpdateResult, error) {
	r.logger.Info("Repository Level: Attempting to add team member")
	user, err := r.userDAO.GetUserByEmail(ctx, email)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve user by email", err)
		return nil, err
	}
	member := models.TeamMember{
		UserID:     user.ID,
		DateAdded:  time.Now(),
		LastActive: time.Now(),
	}
	result, err := r.teamDAO.AddMember(ctx, id, member)
	if err != nil {
		r.logger.Error("Repository Level: Failed to add team member", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully added team member")
	return result, nil
}

func (r *TeamRepository) RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	r.logger.Info("Repository Level: Attempting to remove team member")
	result, err := r.teamDAO.RemoveMember(ctx, id, memberID)
	if err != nil {
		r.logger.Error("Repository Level: Failed to remove team member", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully removed team member")
	return result, nil
}
