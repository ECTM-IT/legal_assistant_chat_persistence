package daos

import (
	"context"
	"errors"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// InvitationDAOInterface defines the interface for the InvitationDAO
type InvitationDAOInterface interface {
	CreateInvitation(ctx context.Context, invitation *models.TeamInvitation) (*models.TeamInvitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*models.TeamInvitation, error)
	UpdateInvitation(ctx context.Context, id primitive.ObjectID, update bson.M) error
	DeleteExpiredInvitations(ctx context.Context) error
}

// InvitationDAO implements the InvitationDAOInterface
type InvitationDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewInvitationDAO creates a new InvitationDAO
func NewInvitationDAO(db *mongo.Database, logger logs.Logger) *InvitationDAO {
	return &InvitationDAO{
		collection: db.Collection("team_invitations"),
		logger:     logger,
	}
}

// CreateInvitation creates a new team invitation in the database
func (dao *InvitationDAO) CreateInvitation(ctx context.Context, invitation *models.TeamInvitation) (*models.TeamInvitation, error) {
	dao.logger.Info("DAO Level: Attempting to create new team invitation")
	_, err := dao.collection.InsertOne(ctx, invitation)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to create team invitation", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully created new team invitation")
	return invitation, nil
}

// GetInvitationByToken retrieves a team invitation by its token
func (dao *InvitationDAO) GetInvitationByToken(ctx context.Context, token string) (*models.TeamInvitation, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve team invitation by token")
	var invitation models.TeamInvitation
	err := dao.collection.FindOne(ctx, bson.M{
		"token":      token,
		"is_used":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&invitation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("Team invitation not found or expired")
			return nil, errors.New("team invitation not found or expired")
		}
		dao.logger.Error("DAO Level: Failed to retrieve team invitation", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved team invitation")
	return &invitation, nil
}

// UpdateInvitation updates a team invitation in the database
func (dao *InvitationDAO) UpdateInvitation(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	dao.logger.Info("DAO Level: Attempting to update team invitation")
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to update team invitation", err)
		return err
	}
	if result.MatchedCount == 0 {
		dao.logger.Warn("Team invitation not found")
		return errors.New("team invitation not found")
	}
	dao.logger.Info("DAO Level: Successfully updated team invitation")
	return nil
}

// DeleteExpiredInvitations deletes all expired invitations from the database
func (dao *InvitationDAO) DeleteExpiredInvitations(ctx context.Context) error {
	dao.logger.Info("DAO Level: Attempting to delete expired team invitations")
	_, err := dao.collection.DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
		"is_used":    false,
	})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to delete expired team invitations", err)
		return err
	}
	dao.logger.Info("DAO Level: Successfully deleted expired team invitations")
	return nil
}
