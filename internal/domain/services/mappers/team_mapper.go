package mappers

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamConversionService interface {
	TeamToDTO(team *models.Team) *dtos.TeamResponse
	TeamsToDTO(teams []models.Team) []dtos.TeamResponse
	TeamMemberToDTO(member *models.TeamMember) *dtos.TeamMemberResponse
	TeamMembersToDTO(members []models.TeamMember) []dtos.TeamMemberResponse
	TeamInvitationToDTO(invitation *models.TeamInvitation) *dtos.TeamInvitationResponse
	TeamInvitationsToDTO(invitations []models.TeamInvitation) []dtos.TeamInvitationResponse
}

type TeamConversionServiceImpl struct {
	logger logs.Logger
}

func NewTeamConversionService(logger logs.Logger) *TeamConversionServiceImpl {
	return &TeamConversionServiceImpl{
		logger: logger,
	}
}

func (s *TeamConversionServiceImpl) TeamToDTO(team *models.Team) *dtos.TeamResponse {
	s.logger.Info("Converting Team to DTO")
	if team == nil {
		s.logger.Warn("Attempted to convert nil Team to DTO")
		return nil
	}

	dto := &dtos.TeamResponse{
		ID:          helpers.Nullable[primitive.ObjectID]{Value: team.ID},
		Name:        helpers.Nullable[string]{Value: team.Name},
		Description: helpers.Nullable[string]{Value: team.Description},
		Members:     helpers.Nullable[[]dtos.TeamMemberResponse]{Value: s.TeamMembersToDTO(team.Members)},
		CreatedAt:   helpers.Nullable[time.Time]{Value: team.CreatedAt},
		UpdatedAt:   helpers.Nullable[time.Time]{Value: team.UpdatedAt},
	}
	s.logger.Info("Successfully converted Team to DTO")
	return dto
}

func (s *TeamConversionServiceImpl) TeamsToDTO(teams []models.Team) []dtos.TeamResponse {
	s.logger.Info("Converting multiple Teams to DTOs")
	teamResponses := make([]dtos.TeamResponse, len(teams))
	for i, team := range teams {
		teamCopy := team // Create a copy to avoid issues with pointer references
		if response := s.TeamToDTO(&teamCopy); response != nil {
			teamResponses[i] = *response
		}
	}
	s.logger.Info("Successfully converted multiple Teams to DTOs")
	return teamResponses
}

func (s *TeamConversionServiceImpl) TeamMemberToDTO(member *models.TeamMember) *dtos.TeamMemberResponse {
	s.logger.Info("Converting TeamMember to DTO")
	if member == nil {
		s.logger.Warn("Attempted to convert nil TeamMember to DTO")
		return nil
	}

	dto := &dtos.TeamMemberResponse{
		ID:         helpers.Nullable[primitive.ObjectID]{Value: member.ID},
		UserID:     helpers.Nullable[primitive.ObjectID]{Value: member.UserID},
		Role:       helpers.Nullable[models.Role]{Value: member.Role},
		FirstName:  helpers.Nullable[string]{Value: member.FirstName},
		LastName:   helpers.Nullable[string]{Value: member.LastName},
		Email:      helpers.Nullable[string]{Value: member.Email},
		DateAdded:  helpers.Nullable[time.Time]{Value: member.DateAdded},
		LastActive: helpers.Nullable[time.Time]{Value: member.LastActive},
	}
	s.logger.Info("Successfully converted TeamMember to DTO")
	return dto
}

func (s *TeamConversionServiceImpl) TeamMembersToDTO(members []models.TeamMember) []dtos.TeamMemberResponse {
	s.logger.Info("Converting multiple TeamMembers to DTOs")
	memberResponses := make([]dtos.TeamMemberResponse, 0, len(members))
	for _, member := range members {
		if !member.IsDeleted {
			memberCopy := member // Create a copy to avoid issues with pointer references
			if response := s.TeamMemberToDTO(&memberCopy); response != nil {
				memberResponses = append(memberResponses, *response)
			}
		}
	}
	s.logger.Info("Successfully converted multiple TeamMembers to DTOs")
	return memberResponses
}

func (s *TeamConversionServiceImpl) TeamInvitationToDTO(invitation *models.TeamInvitation) *dtos.TeamInvitationResponse {
	s.logger.Info("Converting TeamInvitation to DTO")
	if invitation == nil {
		s.logger.Warn("Attempted to convert nil TeamInvitation to DTO")
		return nil
	}

	dto := &dtos.TeamInvitationResponse{
		ID:        helpers.Nullable[primitive.ObjectID]{Value: invitation.ID},
		Email:     helpers.Nullable[string]{Value: invitation.Email},
		Role:      helpers.Nullable[models.Role]{Value: invitation.Role},
		CreatedAt: helpers.Nullable[time.Time]{Value: invitation.CreatedAt},
		ExpiresAt: helpers.Nullable[time.Time]{Value: invitation.ExpiresAt},
		IsUsed:    helpers.Nullable[bool]{Value: invitation.IsUsed},
	}

	s.logger.Info("Successfully converted TeamInvitation to DTO")
	return dto
}

func (s *TeamConversionServiceImpl) TeamInvitationsToDTO(invitations []models.TeamInvitation) []dtos.TeamInvitationResponse {
	s.logger.Info("Converting multiple TeamInvitations to DTOs")
	invitationResponses := make([]dtos.TeamInvitationResponse, 0, len(invitations))
	for _, invitation := range invitations {
		if response := s.TeamInvitationToDTO(&invitation); response != nil {
			invitationResponses = append(invitationResponses, *response)
		}
	}
	s.logger.Info("Successfully converted multiple TeamInvitations to DTOs")
	return invitationResponses
}

func (s *TeamConversionServiceImpl) AcceptInvitationToDTO(invitation *models.TeamInvitation) *dtos.AcceptInvitationRequest {
	s.logger.Info("Converting AcceptInvitation to DTO")
	if invitation == nil {
		s.logger.Warn("Attempted to convert nil AcceptInvitation to DTO")
		return nil
	}

	dto := &dtos.AcceptInvitationRequest{
		Token: helpers.Nullable[string]{Value: invitation.Token},
	}

	s.logger.Info("Successfully converted AcceptInvitation to DTO")
	return dto
}
