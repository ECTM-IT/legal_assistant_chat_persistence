package mappers

import (
	"errors"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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

type TeamConversionServiceImpl struct{}

func NewTeamConversionService() *TeamConversionServiceImpl {
	return &TeamConversionServiceImpl{}
}

func (s *TeamConversionServiceImpl) TeamToDTO(team *models.Team) *dtos.TeamResponse {
	if team == nil {
		return nil
	}

	return &dtos.TeamResponse{
		ID:      helpers.NewNullable(team.ID),
		AdminID: helpers.NewNullable(team.AdminID),
		Members: helpers.NewNullable(*s.TeamMembersToDTO(team.Members)),
	}
}

func (s *TeamConversionServiceImpl) TeamsToDTO(teams []models.Team) []dtos.TeamResponse {
	teamResponses := make([]dtos.TeamResponse, len(teams))
	for i, team := range teams {
		teamResponses[i] = *s.TeamToDTO(&team)
	}
	return teamResponses
}

func (s *TeamConversionServiceImpl) DTOToTeam(teamDTO *dtos.CreateTeamRequest) (*models.Team, error) {
	if teamDTO == nil {
		return nil, errors.New("team DTO cannot be nil")
	}

	if !teamDTO.AdminID.Present {
		return nil, errors.New("admin ID is required")
	}

	members, err := s.DTOToTeamMembers(teamDTO.Members.Value)
	if err != nil {
		return nil, err
	}

	return &models.Team{
		ID:      primitive.NewObjectID(),
		AdminID: teamDTO.AdminID.Value,
		Members: members,
	}, nil
}

func (s *TeamConversionServiceImpl) UpdateTeamFieldsToMap(updateRequest dtos.UpdateTeamRequest) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if updateRequest.AdminID.Present {
		updateFields["admin_id"] = updateRequest.AdminID.Value
	}
	if updateRequest.Members.Present {
		members, _ := s.DTOToTeamMembers(updateRequest.Members.Value)
		updateFields["members"] = members
	}

	return updateFields
}

func (s *TeamConversionServiceImpl) TeamMemberToDTO(member *models.TeamMember) *dtos.TeamMemberResponse {
	return &dtos.TeamMemberResponse{
		ID:         helpers.NewNullable(member.ID),
		UserID:     helpers.NewNullable(member.UserID),
		DateAdded:  helpers.NewNullable(member.DateAdded),
		LastActive: helpers.NewNullable(member.LastActive),
	}
}

func (s *TeamConversionServiceImpl) DTOToTeamMember(memberDTO dtos.TeamMemberResponse) (models.TeamMember, error) {
	if !memberDTO.UserID.Present {
		return models.TeamMember{}, errors.New("user ID is required for team member")
	}

	return models.TeamMember{
		ID:         memberDTO.ID.OrElse(primitive.NewObjectID()),
		UserID:     memberDTO.UserID.Value,
		DateAdded:  memberDTO.DateAdded.OrElse(time.Now()),
		LastActive: memberDTO.LastActive.OrElse(time.Now()),
	}, nil
}

func (s *TeamConversionServiceImpl) TeamMembersToDTO(members []models.TeamMember) *[]dtos.TeamMemberResponse {
	memberResponses := make([]dtos.TeamMemberResponse, len(members))
	for i, member := range members {
		memberResponses[i] = *s.TeamMemberToDTO(&member)
	}
	return &memberResponses
}

func (s *TeamConversionServiceImpl) DTOToTeamMembers(membersDTO []dtos.TeamMemberResponse) ([]models.TeamMember, error) {
	members := make([]models.TeamMember, len(membersDTO))
	for i, memberDTO := range membersDTO {
		member, err := s.DTOToTeamMember(memberDTO)
		if err != nil {
			return nil, err
		}
		members[i] = member
	}
	return members, nil
}
