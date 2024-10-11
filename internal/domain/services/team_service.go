package services

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
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
	AddMember(ctx context.Context, id primitive.ObjectID, request dtos.AddMemberRequest) (*dtos.TeamResponse, error)
	RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*mongo.UpdateResult, error)
}

// TeamServiceImpl implements the TeamService interface.
type TeamServiceImpl struct {
	teamRepo *repositories.TeamRepository
	mapper   *mappers.TeamConversionServiceImpl
	logger   logs.Logger
}

// NewTeamService creates a new instance of the team service.
func NewTeamService(teamRepo *repositories.TeamRepository, mapper *mappers.TeamConversionServiceImpl, logger logs.Logger) *TeamServiceImpl {
	return &TeamServiceImpl{
		teamRepo: teamRepo,
		mapper:   mapper,
		logger:   logger,
	}
}

// CreateTeam creates a new team.
func (s *TeamServiceImpl) CreateTeam(ctx context.Context, request dtos.CreateTeamRequest) (*dtos.TeamResponse, error) {
	s.logger.Info("Service Level: Attempting to create new team")
	team, err := s.mapper.DTOToTeam(&request)
	if err != nil {
		s.logger.Error("Service Level: Failed to convert DTO to team", err)
		return nil, err
	}

	createdTeam, err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		s.logger.Error("Service Level: Failed to create team", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to create team", "create_team_failed")
	}

	response := s.mapper.TeamToDTO(createdTeam)
	s.logger.Info("Service Level: Successfully created new team")
	return response, nil
}

// GetTeamByID retrieves a team by its ID.
func (s *TeamServiceImpl) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*dtos.TeamResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve team by ID")
	team, err := s.teamRepo.GetTeamByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warn("Team not found")
			return nil, errors.NewNotFoundError("Team not found", "team_not_found")
		}
		s.logger.Error("Service Level: Failed to get team", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get team", "get_team_failed")
	}

	response := s.mapper.TeamToDTO(team)
	s.logger.Info("Service Level: Successfully retrieved team by ID")
	return response, nil
}

// GetAllTeams retrieves all teams.
func (s *TeamServiceImpl) GetAllTeams(ctx context.Context) ([]dtos.TeamResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve all teams")
	teams, err := s.teamRepo.GetAllTeams(ctx)
	if err != nil {
		s.logger.Error("Service Level: Failed to get all teams", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get all teams", "get_all_teams_failed")
	}

	response := s.mapper.TeamsToDTO(teams)
	s.logger.Info("Service Level: Successfully retrieved all teams")
	return response, nil
}

// UpdateTeam updates an existing team.
func (s *TeamServiceImpl) UpdateTeam(ctx context.Context, id primitive.ObjectID, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error) {
	s.logger.Info("Service Level: Attempting to update team")
	updateFields := s.mapper.UpdateTeamFieldsToMap(request)

	updatedTeam, err := s.teamRepo.UpdateTeam(ctx, id, updateFields)
	if err != nil {
		s.logger.Error("Service Level: Failed to update team", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to update team", "update_team_failed")
	}

	response := s.mapper.TeamToDTO(updatedTeam)
	s.logger.Info("Service Level: Successfully updated team")
	return response, nil
}

// DeleteTeam deletes a team by its ID.
func (s *TeamServiceImpl) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	s.logger.Info("Service Level: Attempting to delete team")
	err := s.teamRepo.DeleteTeam(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to delete team", err)
		return errors.NewDatabaseError("Service Level: Failed to delete team", "delete_team_failed")
	}

	s.logger.Info("Service Level: Successfully deleted team")
	return nil
}

// GetTeamMember retrieves a team member by their ID.
func (s *TeamServiceImpl) GetTeamMember(ctx context.Context, id primitive.ObjectID) (*dtos.TeamMemberResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve team member")
	user, err := s.teamRepo.GetTeamMember(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warn("Team member not found")
			return nil, errors.NewNotFoundError("Team member not found", "team_member_not_found")
		}
		s.logger.Error("Service Level: Failed to get team member", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get team member", "get_team_member_failed")
	}

	response := s.mapper.TeamMemberToDTO(&models.TeamMember{
		ID:         user.ID,
		UserID:     user.ID,
		DateAdded:  user.SubscriptionID.Timestamp(),
		LastActive: time.Now(),
	})
	s.logger.Info("Service Level: Successfully retrieved team member")
	return response, nil
}

// ChangeAdmin changes the admin of a team.
func (s *TeamServiceImpl) ChangeAdmin(ctx context.Context, id primitive.ObjectID, request dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error) {
	s.logger.Info("Service Level: Attempting to change team admin")
	newAdmin, err := s.teamRepo.ChangeAdmin(ctx, id, request.Email.Value)
	if err != nil {
		s.logger.Error("Service Level: Failed to change admin", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to change admin", "change_admin_failed")
	}

	response := s.mapper.TeamMemberToDTO(&models.TeamMember{
		ID:         newAdmin.ID,
		UserID:     newAdmin.ID,
		DateAdded:  newAdmin.SubscriptionID.Timestamp(),
		LastActive: time.Now(),
	})
	s.logger.Info("Service Level: Successfully changed team admin")
	return response, nil
}

// AddMember adds a member to a team.
func (s *TeamServiceImpl) AddMember(ctx context.Context, id primitive.ObjectID, request dtos.AddMemberRequest) (*dtos.TeamResponse, error) {
	s.logger.Info("Service Level: Attempting to add member to team")
	_, err := s.teamRepo.AddMember(ctx, id, request.Email.Value)
	if err != nil {
		s.logger.Error("Service Level: Failed to add member to team", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to add member", "add_member_failed")
	}
	response, err := s.GetTeamByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to get updated team after adding member", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully added member to team")
	return response, nil
}

// RemoveMember removes a member from a team.
func (s *TeamServiceImpl) RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	s.logger.Info("Service Level: Attempting to remove member from team")
	result, err := s.teamRepo.RemoveMember(ctx, id, memberID)
	if err != nil {
		s.logger.Error("Service Level: Failed to remove member from team", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully removed member from team")
	return result, nil
}
