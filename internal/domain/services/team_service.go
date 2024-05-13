package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamService struct {
	teamRepo *repositories.TeamRepository
}

func NewTeamService(teamRepo *repositories.TeamRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, request dtos.CreateTeamRequest) (*dtos.TeamResponse, error) {
	return s.teamRepo.CreateTeam(ctx, request)
}

func (s *TeamService) GetTeamByID(ctx context.Context, id string) (*dtos.TeamResponse, error) {
	return s.teamRepo.GetTeamByID(ctx, id)
}

func (s *TeamService) GetAllTeams(ctx context.Context) ([]*dtos.TeamResponse, error) {
	return s.teamRepo.GetAllTeams(ctx)
}

func (s *TeamService) UpdateTeam(ctx context.Context, id string, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error) {
	return s.teamRepo.UpdateTeam(ctx, id, request)
}

func (s *TeamService) DeleteTeam(ctx context.Context, id string) error {
	return s.teamRepo.DeleteTeam(ctx, id)
}

func (s *TeamService) GetTeamMember(ctx context.Context, id string) (*dtos.TeamMemberResponse, error) {
	return s.teamRepo.GetTeamMember(ctx, id)
}

func (s *TeamService) ChangeAdmin(ctx context.Context, id string, request dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error) {
	return s.teamRepo.ChangeAdmin(ctx, id, request)
}

func (s *TeamService) AddMember(ctx context.Context, id string, request dtos.AddMemberRequest) (*dtos.TeamMemberResponse, error) {
	return s.teamRepo.AddMember(ctx, id, request)
}

func (s *TeamService) RemoveMember(ctx context.Context, id string, memberID string) (*mongo.UpdateResult, error) {
	return s.teamRepo.RemoveMember(ctx, id, memberID)
}
