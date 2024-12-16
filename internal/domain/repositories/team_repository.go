package repositories

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamRepository struct {
	teamDAO       *daos.TeamDAO
	userDAO       *daos.UserDAO
	invitationDAO *daos.InvitationDAO
	logger        logs.Logger
}

func NewTeamRepository(teamDAO *daos.TeamDAO, userDAO *daos.UserDAO, invitationDAO *daos.InvitationDAO, logger logs.Logger) *TeamRepository {
	return &TeamRepository{
		teamDAO:       teamDAO,
		userDAO:       userDAO,
		invitationDAO: invitationDAO,
		logger:        logger,
	}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	r.logger.Info("Repository Level: Attempting to create new team")
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()
	createdTeam, err := r.teamDAO.CreateTeam(ctx, team)
	if err != nil {
		r.logger.Error("Repository Level: Failed to create team", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully created new team")
	return createdTeam, nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id primitive.ObjectID) (*models.Team, error) {
	r.logger.Info("Repository Level: Attempting to retrieve team by ID")
	team, err := r.teamDAO.GetTeamByID(ctx, id)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve team", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved team")
	return team, nil
}

func (r *TeamRepository) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	r.logger.Info("Repository Level: Attempting to retrieve all teams")
	filter := bson.M{"is_deleted": false}
	teams, err := r.teamDAO.GetAllTeams(ctx, filter)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve teams", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved all teams")
	return teams, nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.Team, error) {
	r.logger.Info("Repository Level: Attempting to update team")
	update["updated_at"] = time.Now()
	_, err := r.teamDAO.UpdateTeam(ctx, id, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to update team", err)
		return nil, err
	}
	updatedTeam, err := r.teamDAO.GetTeamByID(ctx, id)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve updated team", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully updated team")
	return updatedTeam, nil
}

func (r *TeamRepository) SoftDeleteTeam(ctx context.Context, id primitive.ObjectID) error {
	r.logger.Info("Repository Level: Attempting to soft delete team")
	now := time.Now()
	update := bson.M{
		"is_deleted": true,
		"deleted_at": now,
		"updated_at": now,
	}
	_, err := r.teamDAO.UpdateTeam(ctx, id, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to soft delete team", err)
		return err
	}
	r.logger.Info("Repository Level: Successfully soft deleted team")
	return nil
}

func (r *TeamRepository) UndoTeamDeletion(ctx context.Context, id primitive.ObjectID) error {
	r.logger.Info("Repository Level: Attempting to undo team deletion")
	now := time.Now()
	update := bson.M{
		"is_deleted": false,
		"deleted_at": nil,
		"updated_at": now,
	}
	_, err := r.teamDAO.UpdateTeam(ctx, id, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to undo team deletion", err)
		return err
	}
	r.logger.Info("Repository Level: Successfully undid team deletion")
	return nil
}

func (r *TeamRepository) AddTeamMember(ctx context.Context, teamID primitive.ObjectID, member models.TeamMember) error {
	r.logger.Info("Repository Level: Attempting to add team member")
	member.DateAdded = time.Now()
	member.LastActive = time.Now()

	// Check if there's already an admin if the new member is to be an admin
	if member.Role == models.RoleAdmin {
		hasAdmin, err := r.teamDAO.HasAdminMember(ctx, teamID)
		if err != nil {
			r.logger.Error("Repository Level: Failed to check for existing admin", err)
			return err
		}
		if hasAdmin {
			r.logger.Error("Repository Level: Cannot add another admin, one already exists", nil)
			return mongo.ErrNoDocuments
		}
	}

	_, err := r.teamDAO.AddMember(ctx, teamID, member)
	if err != nil {
		r.logger.Error("Repository Level: Failed to add team member", err)
		return err
	}
	r.logger.Info("Repository Level: Successfully added team member")
	return nil
}

func (r *TeamRepository) UpdateTeamMember(ctx context.Context, teamID, memberID primitive.ObjectID, update bson.M) error {
	r.logger.Info("Repository Level: Attempting to update team member")

	// If updating to admin role, check if there's already an admin
	if role, ok := update["role"].(models.Role); ok && role == models.RoleAdmin {
		hasAdmin, err := r.teamDAO.HasAdminMember(ctx, teamID)
		if err != nil {
			r.logger.Error("Repository Level: Failed to check for existing admin", err)
			return err
		}
		if hasAdmin {
			r.logger.Error("Repository Level: Cannot update to admin role, one already exists", nil)
			return mongo.ErrNoDocuments
		}
	}

	_, err := r.teamDAO.UpdateMember(ctx, teamID, memberID, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to update team member", err)
		return err
	}
	r.logger.Info("Repository Level: Successfully updated team member")
	return nil
}

func (r *TeamRepository) SoftDeleteTeamMember(ctx context.Context, teamID, memberID primitive.ObjectID) error {
	r.logger.Info("Repository Level: Attempting to soft delete team member")
	now := time.Now()
	update := bson.M{
		"is_deleted": true,
		"deleted_at": now,
	}
	_, err := r.teamDAO.UpdateMember(ctx, teamID, memberID, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to soft delete team member", err)
		return err
	}
	r.logger.Info("Repository Level: Successfully soft deleted team member")
	return nil
}

func (r *TeamRepository) UndoTeamMemberDeletion(ctx context.Context, teamID, memberID primitive.ObjectID) error {
	r.logger.Info("Repository Level: Attempting to undo team member deletion")
	update := bson.M{
		"is_deleted": false,
		"deleted_at": nil,
	}
	_, err := r.teamDAO.UpdateMember(ctx, teamID, memberID, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to undo team member deletion", err)
		return err
	}
	r.logger.Info("Repository Level: Successfully undid team member deletion")
	return nil
}

func (r *TeamRepository) CreateInvitation(ctx context.Context, invitation *models.TeamInvitation) (*models.TeamInvitation, error) {
	r.logger.Info("Repository Level: Attempting to create team invitation")
	invitation.CreatedAt = time.Now()
	invitation.ExpiresAt = time.Now().Add(24 * time.Hour) // 24-hour expiration
	createdInvitation, err := r.invitationDAO.CreateInvitation(ctx, invitation)
	if err != nil {
		r.logger.Error("Repository Level: Failed to create team invitation", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully created team invitation")
	return createdInvitation, nil
}

func (r *TeamRepository) GetInvitationByToken(ctx context.Context, token string) (*models.TeamInvitation, error) {
	r.logger.Info("Repository Level: Attempting to retrieve invitation by token")
	invitation, err := r.invitationDAO.GetInvitationByToken(ctx, token)
	if err != nil {
		r.logger.Error("Repository Level: Failed to retrieve invitation", err)
		return nil, err
	}
	r.logger.Info("Repository Level: Successfully retrieved invitation")
	return invitation, nil
}

func (r *TeamRepository) MarkInvitationAsUsed(ctx context.Context, id primitive.ObjectID) error {
	r.logger.Info("Repository Level: Attempting to mark invitation as used")
	update := bson.M{"is_used": true}
	err := r.invitationDAO.UpdateInvitation(ctx, id, update)
	if err != nil {
		r.logger.Error("Repository Level: Failed to mark invitation as used", err)
		return err
	}
	r.logger.Info("Repository Level: Successfully marked invitation as used")
	return nil
}
