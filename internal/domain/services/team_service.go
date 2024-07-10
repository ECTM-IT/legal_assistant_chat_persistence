package services

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TeamService defines the operations available for managing teams.
type TeamService interface {
	CreateTeam(ctx context.Context, request dtos.CreateTeamRequest) (*dtos.TeamResponse, error)
	GetTeamByID(ctx context.Context, id primitive.ObjectID) (*dtos.TeamResponse, error)
	GetAllTeams(ctx context.Context) ([]dtos.TeamResponse, error)
	UpdateTeam(ctx context.Context, id primitive.ObjectID, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error)
	DeleteTeam(ctx context.Context, id primitive.ObjectID) error
	GetTeamMember(ctx context.Context, id primitive.ObjectID) (*dtos.TeamMemberResponse, error)
	ChangeAdmin(ctx context.Context, id primitive.ObjectID, request dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error)
	AddMember(ctx context.Context, id primitive.ObjectID, request dtos.AddMemberRequest) (*dtos.TeamMemberResponse, error)
	RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) error
}

// TeamServiceImpl implements the TeamService interface.
type TeamServiceImpl struct {
	teamRepo *repositories.TeamRepository
	mapper   *mappers.TeamConversionServiceImpl
}

// NewTeamService creates a new instance of the team service.
func NewTeamService(teamRepo *repositories.TeamRepository, mapper *mappers.TeamConversionServiceImpl) *TeamServiceImpl {
	return &TeamServiceImpl{
		teamRepo: teamRepo,
		mapper:   mapper,
	}
}

// CreateTeam creates a new team.
func (s *TeamServiceImpl) CreateTeam(ctx context.Context, request dtos.CreateTeamRequest) (*dtos.TeamResponse, error) {
	team, err := s.mapper.DTOToTeam(&request)
	if err != nil {
		return nil, err
	}

	createdTeam, err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create team", "create_team_failed")
	}

	return s.mapper.TeamToDTO(createdTeam), nil
}

// GetTeamByID retrieves a team by its ID.
func (s *TeamServiceImpl) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*dtos.TeamResponse, error) {
	team, err := s.teamRepo.GetTeamByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("Team not found", "team_not_found")
		}
		return nil, errors.NewDatabaseError("Failed to get team", "get_team_failed")
	}

	return s.mapper.TeamToDTO(team), nil
}

// GetAllTeams retrieves all teams.
func (s *TeamServiceImpl) GetAllTeams(ctx context.Context) ([]dtos.TeamResponse, error) {
	teams, err := s.teamRepo.GetAllTeams(ctx)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get all teams", "get_all_teams_failed")
	}

	return s.mapper.TeamsToDTO(teams), nil
}

// UpdateTeam updates an existing team.
func (s *TeamServiceImpl) UpdateTeam(ctx context.Context, id primitive.ObjectID, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error) {
	updateFields := s.mapper.UpdateTeamFieldsToMap(request)

	updatedTeam, err := s.teamRepo.UpdateTeam(ctx, id, updateFields)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to update team", "update_team_failed")
	}

	return s.mapper.TeamToDTO(updatedTeam), nil
}

// DeleteTeam deletes a team by its ID.
func (s *TeamServiceImpl) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	err := s.teamRepo.DeleteTeam(ctx, id)
	if err != nil {
		return errors.NewDatabaseError("Failed to delete team", "delete_team_failed")
	}

	return nil
}

// GetTeamMember retrieves a team member by their ID.
func (s *TeamServiceImpl) GetTeamMember(ctx context.Context, id primitive.ObjectID) (*dtos.TeamMemberResponse, error) {
	user, err := s.teamRepo.GetTeamMember(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("Team member not found", "team_member_not_found")
		}
		return nil, errors.NewDatabaseError("Failed to get team member", "get_team_member_failed")
	}

	return s.mapper.TeamMemberToDTO(&models.TeamMember{
		ID:         user.ID,
		UserID:     user.ID,
		DateAdded:  user.SubscriptionID.Timestamp(),
		LastActive: time.Now(),
	}), nil
}

// ChangeAdmin changes the admin of a team.
func (s *TeamServiceImpl) ChangeAdmin(ctx context.Context, id primitive.ObjectID, request dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error) {
	newAdmin, err := s.teamRepo.ChangeAdmin(ctx, id, request.Email.Value)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to change admin", "change_admin_failed")
	}

	return s.mapper.TeamMemberToDTO(&models.TeamMember{
		ID:         newAdmin.ID,
		UserID:     newAdmin.ID,
		DateAdded:  newAdmin.SubscriptionID.Timestamp(),
		LastActive: time.Now(),
	}), nil
}

// AddMember adds a member to a team.
func (s *TeamServiceImpl) AddMember(ctx context.Context, id primitive.ObjectID, request dtos.AddMemberRequest) (*dtos.TeamResponse, error) {
	_, err := s.teamRepo.AddMember(ctx, id, request.Email.Value)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to add member", "add_member_failed")
	}
	return s.GetTeamByID(ctx, id)
}

// RemoveMember removes a member from a team.
func (s *TeamServiceImpl) RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return s.teamRepo.RemoveMember(ctx, id, memberID)
}
