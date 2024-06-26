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

func (r *TeamRepository) CreateTeam(ctx context.Context, request dtos.CreateTeamRequest) (*dtos.TeamResponse, error) {
	membersRequest := request.Members
	members := []models.TeamMember{}
	if membersRequest.Present {
		for _, singleMember := range membersRequest.Value {
			member := models.TeamMember{
				ID:         singleMember.ID.OrElse(primitive.NilObjectID),
				UserID:     singleMember.UserID.OrElse(primitive.NilObjectID),
				DateAdded:  singleMember.DateAdded.OrElse(time.Time{}),
				LastActive: singleMember.LastActive.OrElse(time.Time{}),
			}
			members = append(members, member)
		}
	}
	team := &models.Team{
		ID:      primitive.NewObjectID(),
		AdminID: request.AdminID.OrElse(primitive.NilObjectID),
		Members: members,
	}
	team, err := r.teamDAO.CreateTeam(ctx, team)
	if err != nil {
		return nil, err
	}
	return r.toTeamResponse(team), nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*dtos.TeamResponse, error) {
	team, err := r.teamDAO.GetTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toTeamResponse(team), nil
}

func (r *TeamRepository) GetAllTeams(ctx context.Context) ([]*dtos.TeamResponse, error) {
	teams, err := r.teamDAO.GetAllTeams(ctx)
	if err != nil {
		return nil, err
	}
	var teamResponses []*dtos.TeamResponse
	for _, team := range teams {
		teamResponses = append(teamResponses, r.toTeamResponse(&team))
	}
	return teamResponses, nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, id primitive.ObjectID, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error) {
	update := bson.M{}
	if request.AdminID.Present {
		update["admin_id"] = request.AdminID.Value
	}
	if request.Members.Present {
		update["members"] = request.Members.Value
	}
	_, err := r.teamDAO.UpdateTeam(ctx, id, update)
	if err != nil {
		return nil, err
	}
	team, err := r.teamDAO.GetTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toTeamResponse(team), nil
}

func (r *TeamRepository) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	return r.teamDAO.DeleteTeam(ctx, id)
}

func (r *TeamRepository) GetTeamMember(ctx context.Context, id primitive.ObjectID) (*dtos.TeamMemberResponse, error) {
	user, err := r.userDAO.GetUserByID(ctx, id)
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

func (r *TeamRepository) ChangeAdmin(ctx context.Context, id primitive.ObjectID, request dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error) {
	user, err := r.userDAO.GetUserByEmail(ctx, request.Email.OrElse(""))
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
	return &dtos.TeamMemberResponse{
		ID:         helpers.NewNullable(user.ID),
		UserID:     helpers.NewNullable(user.ID),
		DateAdded:  helpers.NewNullable(time.Now()),
		LastActive: helpers.NewNullable(time.Now()),
	}, nil
}

func (r *TeamRepository) AddMember(ctx context.Context, id primitive.ObjectID, request dtos.AddMemberRequest) (*dtos.TeamMemberResponse, error) {
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
	_, err = r.teamDAO.AddMember(ctx, id, member)
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

func (r *TeamRepository) RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return r.teamDAO.RemoveMember(ctx, id, memberID)
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
