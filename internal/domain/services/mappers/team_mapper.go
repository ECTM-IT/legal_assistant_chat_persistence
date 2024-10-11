package mappers

import (
	"errors"
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
	DTOToTeam(teamDTO *dtos.CreateTeamRequest) (*models.Team, error)
	UpdateTeamFieldsToMap(updateRequest dtos.UpdateTeamRequest) map[string]interface{}
	TeamMemberToDTO(member models.TeamMember) *dtos.TeamMemberResponse
	DTOToTeamMember(memberDTO dtos.TeamMemberResponse) (models.TeamMember, error)
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
		ID:      helpers.NewNullable(team.ID),
		AdminID: helpers.NewNullable(team.AdminID),
		Members: helpers.NewNullable(*s.TeamMembersToDTO(team.Members)),
	}
	s.logger.Info("Successfully converted Team to DTO")
	return dto
}

func (s *TeamConversionServiceImpl) TeamsToDTO(teams []models.Team) []dtos.TeamResponse {
	s.logger.Info("Converting multiple Teams to DTOs")
	teamResponses := make([]dtos.TeamResponse, len(teams))
	for i, team := range teams {
		teamResponses[i] = *s.TeamToDTO(&team)
	}
	s.logger.Info("Successfully converted multiple Teams to DTOs")
	return teamResponses
}

func (s *TeamConversionServiceImpl) DTOToTeam(teamDTO *dtos.CreateTeamRequest) (*models.Team, error) {
	s.logger.Info("Converting DTO to Team")
	if teamDTO == nil {
		s.logger.Error("Failed to convert DTO to Team: team DTO cannot be nil", errors.New("team DTO cannot be nil"))
		return nil, errors.New("team DTO cannot be nil")
	}

	if !teamDTO.AdminID.Present {
		s.logger.Error("Failed to convert DTO to Team: admin ID is required", errors.New("admin ID is required"))
		return nil, errors.New("admin ID is required")
	}

	members, err := s.DTOToTeamMembers(teamDTO.Members.Value)
	if err != nil {
		s.logger.Error("Failed to convert team members", err)
		return nil, err
	}

	team := &models.Team{
		ID:      primitive.NewObjectID(),
		AdminID: teamDTO.AdminID.Value,
		Members: members,
	}
	s.logger.Info("Successfully converted DTO to Team")
	return team, nil
}

func (s *TeamConversionServiceImpl) UpdateTeamFieldsToMap(updateRequest dtos.UpdateTeamRequest) map[string]interface{} {
	s.logger.Info("Converting UpdateTeamRequest to map")
	updateFields := make(map[string]interface{})

	if updateRequest.AdminID.Present {
		updateFields["admin_id"] = updateRequest.AdminID.Value
	}
	if updateRequest.Members.Present {
		members, err := s.DTOToTeamMembers(updateRequest.Members.Value)
		if err != nil {
			s.logger.Error("Failed to convert team members for update", err)
		} else {
			updateFields["members"] = members
		}
	}

	s.logger.Info("Successfully converted UpdateTeamRequest to map")
	return updateFields
}

func (s *TeamConversionServiceImpl) TeamMemberToDTO(member *models.TeamMember) *dtos.TeamMemberResponse {
	s.logger.Info("Converting TeamMember to DTO")
	dto := &dtos.TeamMemberResponse{
		ID:         helpers.NewNullable(member.ID),
		UserID:     helpers.NewNullable(member.UserID),
		DateAdded:  helpers.NewNullable(member.DateAdded),
		LastActive: helpers.NewNullable(member.LastActive),
	}
	s.logger.Info("Successfully converted TeamMember to DTO")
	return dto
}

func (s *TeamConversionServiceImpl) DTOToTeamMember(memberDTO dtos.TeamMemberResponse) (models.TeamMember, error) {
	s.logger.Info("Converting DTO to TeamMember")
	if !memberDTO.UserID.Present {
		s.logger.Error("Failed to convert DTO to TeamMember: user ID is required for team member", errors.New("user ID is required for team member"))
		return models.TeamMember{}, errors.New("user ID is required for team member")
	}

	member := models.TeamMember{
		ID:         memberDTO.ID.OrElse(primitive.NewObjectID()),
		UserID:     memberDTO.UserID.Value,
		DateAdded:  memberDTO.DateAdded.OrElse(time.Now()),
		LastActive: memberDTO.LastActive.OrElse(time.Now()),
	}
	s.logger.Info("Successfully converted DTO to TeamMember")
	return member, nil
}

func (s *TeamConversionServiceImpl) TeamMembersToDTO(members []models.TeamMember) *[]dtos.TeamMemberResponse {
	s.logger.Info("Converting multiple TeamMembers to DTOs")
	memberResponses := make([]dtos.TeamMemberResponse, len(members))
	for i, member := range members {
		memberResponses[i] = *s.TeamMemberToDTO(&member)
	}
	s.logger.Info("Successfully converted multiple TeamMembers to DTOs")
	return &memberResponses
}

func (s *TeamConversionServiceImpl) DTOToTeamMembers(membersDTO []dtos.TeamMemberResponse) ([]models.TeamMember, error) {
	s.logger.Info("Converting DTOs to TeamMembers")
	members := make([]models.TeamMember, len(membersDTO))
	for i, memberDTO := range membersDTO {
		member, err := s.DTOToTeamMember(memberDTO)
		if err != nil {
			s.logger.Error("Failed to convert DTO to TeamMember", err)
			return nil, err
		}
		members[i] = member
	}
	s.logger.Info("Successfully converted DTOs to TeamMembers")
	return members, nil
}
