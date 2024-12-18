package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TeamService defines the operations available for managing teams.
type TeamService interface {
	CreateTeam(ctx context.Context, request dtos.CreateTeamRequest) (*dtos.TeamResponse, error)
	GetTeamByID(ctx context.Context, id primitive.ObjectID) (*dtos.TeamResponse, error)
	GetAllTeams(ctx context.Context) ([]dtos.TeamResponse, error)
	UpdateTeam(ctx context.Context, id primitive.ObjectID, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error)
	SoftDeleteTeam(ctx context.Context, id primitive.ObjectID) error
	UndoTeamDeletion(ctx context.Context, id primitive.ObjectID) error
	AddTeamMember(ctx context.Context, teamID primitive.ObjectID, request dtos.AddTeamMemberRequest) error
	UpdateTeamMember(ctx context.Context, teamID, memberID primitive.ObjectID, request dtos.UpdateTeamMemberRequest) error
	SoftDeleteTeamMember(ctx context.Context, teamID, memberID primitive.ObjectID) error
	UndoTeamMemberDeletion(ctx context.Context, teamID, memberID primitive.ObjectID) error
	CreateInvitation(ctx context.Context, teamID primitive.ObjectID, request dtos.TeamInvitationRequest) (*dtos.TeamInvitationResponse, error)
	AcceptInvitation(ctx context.Context, request dtos.AcceptInvitationRequest) error
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
	team := &models.Team{
		Name:        request.Name,
		Description: request.Description,
		Members:     []models.TeamMember{},
	}

	createdTeam, err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		s.logger.Error("Service Level: Failed to create team", err)
		return nil, errors.NewDatabaseError("Failed to create team", "create_team_failed")
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
		s.logger.Error("Service Level: Failed to get team", err)
		return nil, errors.NewDatabaseError("Failed to get team", "get_team_failed")
	}

	response := s.mapper.TeamToDTO(team)
	s.logger.Info("Service Level: Successfully retrieved team")
	return response, nil
}

// GetAllTeams retrieves all teams.
func (s *TeamServiceImpl) GetAllTeams(ctx context.Context) ([]dtos.TeamResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve all teams")
	teams, err := s.teamRepo.GetAllTeams(ctx)
	if err != nil {
		s.logger.Error("Service Level: Failed to get all teams", err)
		return nil, errors.NewDatabaseError("Failed to get all teams", "get_all_teams_failed")
	}

	response := s.mapper.TeamsToDTO(teams)
	s.logger.Info("Service Level: Successfully retrieved all teams")
	return response, nil
}

// UpdateTeam updates an existing team.
func (s *TeamServiceImpl) UpdateTeam(ctx context.Context, id primitive.ObjectID, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error) {
	s.logger.Info("Service Level: Attempting to update team")
	update := bson.M{}
	if request.Name != nil {
		update["name"] = *request.Name
	}
	if request.Description != nil {
		update["description"] = *request.Description
	}

	updatedTeam, err := s.teamRepo.UpdateTeam(ctx, id, update)
	if err != nil {
		s.logger.Error("Service Level: Failed to update team", err)
		return nil, errors.NewDatabaseError("Failed to update team", "update_team_failed")
	}

	response := s.mapper.TeamToDTO(updatedTeam)
	s.logger.Info("Service Level: Successfully updated team")
	return response, nil
}

// SoftDeleteTeam soft deletes a team.
func (s *TeamServiceImpl) SoftDeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	s.logger.Info("Service Level: Attempting to soft delete team")
	err := s.teamRepo.SoftDeleteTeam(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to soft delete team", err)
		return errors.NewDatabaseError("Failed to delete team", "delete_team_failed")
	}
	s.logger.Info("Service Level: Successfully soft deleted team")
	return nil
}

// UndoTeamDeletion undoes a team deletion.
func (s *TeamServiceImpl) UndoTeamDeletion(ctx context.Context, id primitive.ObjectID) error {
	s.logger.Info("Service Level: Attempting to undo team deletion")
	err := s.teamRepo.UndoTeamDeletion(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to undo team deletion", err)
		return errors.NewDatabaseError("Failed to undo team deletion", "undo_team_deletion_failed")
	}
	s.logger.Info("Service Level: Successfully undid team deletion")
	return nil
}

// AddTeamMember adds a new team member.
func (s *TeamServiceImpl) AddTeamMember(ctx context.Context, teamID primitive.ObjectID, request dtos.AddTeamMemberRequest) error {
	s.logger.Info("Service Level: Attempting to add team member")
	member := models.TeamMember{
		ID:        primitive.NewObjectID(),
		Role:      request.Role,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
	}

	err := s.teamRepo.AddTeamMember(ctx, teamID, member)
	if err != nil {
		s.logger.Error("Service Level: Failed to add team member", err)
		return errors.NewDatabaseError("Failed to add team member", "add_team_member_failed")
	}
	s.logger.Info("Service Level: Successfully added team member")
	return nil
}

// UpdateTeamMember updates a team member.
func (s *TeamServiceImpl) UpdateTeamMember(ctx context.Context, teamID, memberID primitive.ObjectID, request dtos.UpdateTeamMemberRequest) error {
	s.logger.Info("Service Level: Attempting to update team member")
	update := bson.M{}
	if request.Role != nil {
		update["role"] = *request.Role
	}
	if request.FirstName != nil {
		update["first_name"] = *request.FirstName
	}
	if request.LastName != nil {
		update["last_name"] = *request.LastName
	}
	if request.Email != nil {
		update["email"] = *request.Email
	}

	err := s.teamRepo.UpdateTeamMember(ctx, teamID, memberID, update)
	if err != nil {
		s.logger.Error("Service Level: Failed to update team member", err)
		return errors.NewDatabaseError("Failed to update team member", "update_team_member_failed")
	}
	s.logger.Info("Service Level: Successfully updated team member")
	return nil
}

// SoftDeleteTeamMember soft deletes a team member.
func (s *TeamServiceImpl) SoftDeleteTeamMember(ctx context.Context, teamID, memberID primitive.ObjectID) error {
	s.logger.Info("Service Level: Attempting to soft delete team member")
	err := s.teamRepo.SoftDeleteTeamMember(ctx, teamID, memberID)
	if err != nil {
		s.logger.Error("Service Level: Failed to soft delete team member", err)
		return errors.NewDatabaseError("Failed to delete team member", "delete_team_member_failed")
	}
	s.logger.Info("Service Level: Successfully soft deleted team member")
	return nil
}

// UndoTeamMemberDeletion undoes a team member deletion.
func (s *TeamServiceImpl) UndoTeamMemberDeletion(ctx context.Context, teamID, memberID primitive.ObjectID) error {
	s.logger.Info("Service Level: Attempting to undo team member deletion")
	err := s.teamRepo.UndoTeamMemberDeletion(ctx, teamID, memberID)
	if err != nil {
		s.logger.Error("Service Level: Failed to undo team member deletion", err)
		return errors.NewDatabaseError("Failed to undo team member deletion", "undo_team_member_deletion_failed")
	}
	s.logger.Info("Service Level: Successfully undid team member deletion")
	return nil
}

// generateInvitationToken generates a secure random token for team invitations.
func (s *TeamServiceImpl) generateInvitationToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CreateInvitation creates a new team invitation.
func (s *TeamServiceImpl) CreateInvitation(ctx context.Context, teamID primitive.ObjectID, request dtos.TeamInvitationRequest) (*dtos.TeamInvitationResponse, error) {
	s.logger.Info("Service Level: Attempting to create team invitation")
	token, err := s.generateInvitationToken()
	if err != nil {
		s.logger.Error("Service Level: Failed to generate invitation token", err)
		return nil, errors.NewDatabaseError("Failed to generate invitation token", "generate_token_failed")
	}

	invitation := &models.TeamInvitation{
		ID:     primitive.NewObjectID(),
		TeamID: teamID,
		Email:  request.Email,
		Role:   request.Role,
		Token:  token,
	}

	createdInvitation, err := s.teamRepo.CreateInvitation(ctx, invitation)
	if err != nil {
		s.logger.Error("Service Level: Failed to create team invitation", err)
		return nil, errors.NewDatabaseError("Failed to create team invitation", "create_invitation_failed")
	}

	response := &dtos.TeamInvitationResponse{
		ID:        createdInvitation.ID,
		Email:     createdInvitation.Email,
		Role:      createdInvitation.Role,
		CreatedAt: createdInvitation.CreatedAt,
		ExpiresAt: createdInvitation.ExpiresAt,
		IsUsed:    createdInvitation.IsUsed,
	}

	s.logger.Info("Service Level: Successfully created team invitation")
	return response, nil
}

// AcceptInvitation accepts a team invitation.
func (s *TeamServiceImpl) AcceptInvitation(ctx context.Context, request dtos.AcceptInvitationRequest) error {
	s.logger.Info("Service Level: Attempting to accept team invitation")
	invitation, err := s.teamRepo.GetInvitationByToken(ctx, request.Token)
	if err != nil {
		s.logger.Error("Service Level: Failed to get invitation", err)
		return errors.NewNotFoundError("Invalid or expired invitation token", "invalid_token")
	}

	member := models.TeamMember{
		ID:        primitive.NewObjectID(),
		Role:      invitation.Role,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     invitation.Email,
	}

	err = s.teamRepo.AddTeamMember(ctx, invitation.TeamID, member)
	if err != nil {
		s.logger.Error("Service Level: Failed to add team member from invitation", err)
		return errors.NewDatabaseError("Failed to add team member", "add_team_member_failed")
	}

	err = s.teamRepo.MarkInvitationAsUsed(ctx, invitation.ID)
	if err != nil {
		s.logger.Error("Service Level: Failed to mark invitation as used", err)
		return errors.NewDatabaseError("Failed to mark invitation as used", "mark_invitation_used_failed")
	}

	s.logger.Info("Service Level: Successfully accepted team invitation")
	return nil
}
