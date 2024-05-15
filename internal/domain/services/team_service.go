package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TeamService defines the operations available for managing teams.
type TeamService interface {
	CreateTeam(ctx context.Context, request dtos.CreateTeamRequest) (*dtos.TeamResponse, error)
	GetTeamByID(ctx context.Context, id primitive.ObjectID) (*dtos.TeamResponse, error)
	GetAllTeams(ctx context.Context) ([]*dtos.TeamResponse, error)
	UpdateTeam(ctx context.Context, id primitive.ObjectID, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error)
	DeleteTeam(ctx context.Context, id primitive.ObjectID) error
	GetTeamMember(ctx context.Context, id primitive.ObjectID) (*dtos.TeamMemberResponse, error)
	ChangeAdmin(ctx context.Context, id primitive.ObjectID, request dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error)
	AddMember(ctx context.Context, id primitive.ObjectID, request dtos.AddMemberRequest) (*dtos.TeamMemberResponse, error)
	RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*mongo.UpdateResult, error)
}

// TeamServiceImpl implements the TeamService interface.
type TeamServiceImpl struct {
	teamRepo *repositories.TeamRepository
}

// NewTeamService creates a new instance of the team service.
func NewTeamService(teamRepo *repositories.TeamRepository) *TeamServiceImpl {
	return &TeamServiceImpl{
		teamRepo: teamRepo,
	}
}

// CreateTeam creates a new team.
func (s *TeamServiceImpl) CreateTeam(ctx context.Context, request dtos.CreateTeamRequest) (*dtos.TeamResponse, error) {
	return s.teamRepo.CreateTeam(ctx, request)
}

// GetTeamByID retrieves a team by its ID.
func (s *TeamServiceImpl) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*dtos.TeamResponse, error) {
	return s.teamRepo.GetTeamByID(ctx, id)
}

// GetAllTeams retrieves all teams.
func (s *TeamServiceImpl) GetAllTeams(ctx context.Context) ([]*dtos.TeamResponse, error) {
	return s.teamRepo.GetAllTeams(ctx)
}

// UpdateTeam updates an existing team.
func (s *TeamServiceImpl) UpdateTeam(ctx context.Context, id primitive.ObjectID, request dtos.UpdateTeamRequest) (*dtos.TeamResponse, error) {
	return s.teamRepo.UpdateTeam(ctx, id, request)
}

// DeleteTeam deletes a team by its ID.
func (s *TeamServiceImpl) DeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	return s.teamRepo.DeleteTeam(ctx, id)
}

// GetTeamMember retrieves a team member by their ID.
func (s *TeamServiceImpl) GetTeamMember(ctx context.Context, id primitive.ObjectID) (*dtos.TeamMemberResponse, error) {
	return s.teamRepo.GetTeamMember(ctx, id)
}

// ChangeAdmin changes the admin of a team.
func (s *TeamServiceImpl) ChangeAdmin(ctx context.Context, id primitive.ObjectID, request dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error) {
	return s.teamRepo.ChangeAdmin(ctx, id, request)
}

// AddMember adds a member to a team.
func (s *TeamServiceImpl) AddMember(ctx context.Context, id primitive.ObjectID, request dtos.AddMemberRequest) (*dtos.TeamMemberResponse, error) {
	return s.teamRepo.AddMember(ctx, id, request)
}

// RemoveMember removes a member from a team.
func (s *TeamServiceImpl) RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return s.teamRepo.RemoveMember(ctx, id, memberID)
}
