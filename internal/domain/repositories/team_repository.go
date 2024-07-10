package repositories

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamRepository struct {
	teamDAO *daos.TeamDAO
	userDAO *daos.UserDAO
}

func NewTeamRepository(teamDAO *daos.TeamDAO, userDAO *daos.UserDAO) *TeamRepository {
	return &TeamRepository{
		teamDAO: teamDAO,
		userDAO: userDAO,
	}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	return r.teamDAO.CreateTeam(ctx, team)
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*models.Team, error) {
	return r.teamDAO.GetTeamByID(ctx, id)
}

func (r *TeamRepository) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	return r.teamDAO.GetAllTeams(ctx)
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Team, error) {
	_, err := r.teamDAO.UpdateTeam(ctx, id, update)
	if err != nil {
		return nil, err
	}
	return r.teamDAO.GetTeamByID(ctx, id)
}

func (r *TeamRepository) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	return r.teamDAO.DeleteTeam(ctx, id)
}

func (r *TeamRepository) GetTeamMember(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return r.userDAO.GetUserByID(ctx, id)
}

func (r *TeamRepository) ChangeAdmin(ctx context.Context, id primitive.ObjectID, email string) (*models.User, error) {
	user, err := r.userDAO.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"admin_id": user.ID,
	}
	_, err = r.teamDAO.UpdateTeam(ctx, id, update)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *TeamRepository) AddMember(ctx context.Context, id primitive.ObjectID, email string) (*mongo.UpdateResult, error) {
	user, err := r.userDAO.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	member := models.TeamMember{
		UserID:     user.ID,
		DateAdded:  time.Now(),
		LastActive: time.Now(),
	}
	return r.teamDAO.AddMember(ctx, id, member)
}

func (r *TeamRepository) RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return r.teamDAO.RemoveMember(ctx, id, memberID)
}
