package mappers

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
)

type TeamConversionService interface {
	TeamToDTO(team *models.Team) *dtos.TeamResponse
	TeamsToDTO(teams []models.Team) []dtos.TeamResponse
	TeamMemberToDTO(member *models.TeamMember) *dtos.TeamMemberResponse
	TeamMembersToDTO(members []models.TeamMember) []dtos.TeamMemberResponse
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
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		Members:     s.TeamMembersToDTO(team.Members),
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
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
		ID:         member.ID,
		UserID:     member.UserID,
		Role:       member.Role,
		FirstName:  member.FirstName,
		LastName:   member.LastName,
		Email:      member.Email,
		DateAdded:  member.DateAdded,
		LastActive: member.LastActive,
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
