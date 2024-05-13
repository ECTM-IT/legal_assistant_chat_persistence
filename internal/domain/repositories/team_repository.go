package repositories

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
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

func (r *TeamRepository) GetTeamByID(ctx context.Context, id string) (*dtos.TeamResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	team, err := r.teamDAO.GetTeamByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	return r.toTeamResponse(team), nil
}

func (r *TeamRepository) GetTeamMember(ctx context.Context, id string) (*dtos.TeamMemberResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := r.userDAO.GetUserByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	return &dtos.TeamMemberResponse{
		ID:         helpers.NewNullable(user.ID),
		UserID:     helpers.NewNullable(user.ID),
		DateAdded:  helpers.NewNullable(time.Now()),
		LastActive: helpers.NewNullable(time.Now()),
	}, nil
}

func (r *TeamRepository) ChangeAdmin(ctx context.Context, id string, request dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error) {
	teamObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := r.userDAO.GetUserByEmail(ctx, request.Email.OrElse(""))
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"admin_id": user.ID,
	}
	_, err = r.teamDAO.UpdateTeam(ctx, teamObjectID, update)
	if err != nil {
		return nil, err
	}
	return &dtos.TeamMemberResponse{
		ID:         helpers.NewNullable(user.ID),
		UserID:     helpers.NewNullable(user.ID),
		DateAdded:  helpers.NewNullable(time.Now()),
		LastActive: helpers.NewNullable(time.Now()),
	}, nil
}

func (r *TeamRepository) AddMember(ctx context.Context, id string, request dtos.AddMemberRequest) (*dtos.TeamMemberResponse, error) {
	teamObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := r.userDAO.GetUserByEmail(ctx, request.Email.OrElse(""))
	if err != nil {
		return nil, err
	}
	member := models.TeamMember{
		ID:         primitive.NewObjectID(),
		UserID:     user.ID,
		DateAdded:  time.Now(),
		LastActive: time.Now(),
	}
	_, err = r.teamDAO.AddMember(ctx, teamObjectID, member)
	if err != nil {
		return nil, err
	}
	return &dtos.TeamMemberResponse{
		ID:         helpers.NewNullable(member.ID),
		UserID:     helpers.NewNullable(member.UserID),
		DateAdded:  helpers.NewNullable(member.DateAdded),
		LastActive: helpers.NewNullable(member.LastActive),
	}, nil
}

func (r *TeamRepository) RemoveMember(ctx context.Context, id string, memberID string) (*mongo.UpdateResult, error) {
	teamObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	memberObjectID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, err
	}
	return r.teamDAO.RemoveMember(ctx, teamObjectID, memberObjectID)
}

func (r *TeamRepository) toTeamResponse(team *models.Team) *dtos.TeamResponse {
	var memberResponses []dtos.TeamMemberResponse
	for _, member := range team.Members {
		memberResponses = append(memberResponses, dtos.TeamMemberResponse{
			ID:         helpers.NewNullable(member.ID),
			UserID:     helpers.NewNullable(member.UserID),
			DateAdded:  helpers.NewNullable(member.DateAdded),
			LastActive: helpers.NewNullable(member.LastActive),
		})
	}
	return &dtos.TeamResponse{
		ID:      helpers.NewNullable(team.ID),
		AdminID: helpers.NewNullable(team.AdminID),
		Members: helpers.NewNullable(memberResponses),
	}
}
