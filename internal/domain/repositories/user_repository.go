package repositories

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the operations available on a user repository.
type UserRepository interface {
	FindUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*models.User, error)
	TotalUsers(ctx context.Context) ([]*models.User, error)
	DeleteUser(ctx context.Context, userID primitive.ObjectID) error
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}) (*models.User, error)
}

// UserRepositoryImpl implements the UserRepository interface.
type UserRepositoryImpl struct {
	userDAO *daos.UserDAO
}

// NewUserRepository creates a new instance of the user repository.
func NewUserRepository(userDAO *daos.UserDAO) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		userDAO: userDAO,
	}
}

// FindUserByID retrieves a user by their ID.
func (r *UserRepositoryImpl) FindUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	return r.userDAO.GetUserByID(ctx, userID)
}

// FindUserByEmail retrieves a user by their email.
func (r *UserRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.userDAO.GetUserByEmail(ctx, email)
}

// FindUserByCaseID retrieves a user by a case ID.
func (r *UserRepositoryImpl) FindUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*models.User, error) {
	return r.userDAO.GetUserByCaseID(ctx, caseID)
}

// TotalUsers retrieves all users.
func (r *UserRepositoryImpl) TotalUsers(ctx context.Context) ([]*models.User, error) {
	return r.userDAO.GetAllUsers(ctx)
}

// DeleteUser deletes a user by their ID.
func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	return r.userDAO.DeleteUser(ctx, userID)
}

// CreateUser creates a new user.
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return r.userDAO.CreateUser(ctx, user)
}

// UpdateUser updates an existing user.
func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, userID primitive.ObjectID, updates map[string]interface{}) (*models.User, error) {
	_, err := r.userDAO.UpdateUser(ctx, userID, updates)
	if err != nil {
		return nil, err
	}
	return r.FindUserByID(ctx, userID)
}
