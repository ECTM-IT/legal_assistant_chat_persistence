package repositories

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository defines the operations available on a user repository.
type UserRepository interface {
	FindUserByID(ctx context.Context, userID primitive.ObjectID) (*dtos.UserResponse, error)
	FindUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error)
	FindUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*dtos.UserResponse, error)
	TotalUsers(ctx context.Context) ([]*models.User, error)
	DeleteUser(ctx context.Context, userID primitive.ObjectID) error
	CreateUser(ctx context.Context, user *dtos.CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, userID primitive.ObjectID, user *dtos.UpdateUserRequest) (*mongo.UpdateResult, error)
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
func (r *UserRepositoryImpl) FindUserByID(ctx context.Context, userID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := r.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return r.toUserResponse(user), nil
}

// FindUserByEmail retrieves a user by their email.
func (r *UserRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error) {
	user, err := r.userDAO.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return r.toUserResponse(user), nil
}

// FindUserByCaseID retrieves a user by a case ID.
func (r *UserRepositoryImpl) FindUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := r.userDAO.GetUserByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	return r.toUserResponse(user), nil
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
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *dtos.CreateUserRequest) (*models.User, error) {
	userModel := &models.User{
		ID:             primitive.NewObjectID(),
		Image:          user.Image.OrElse(""),
		Email:          user.Email.OrElse(""),
		FirstName:      user.FirstName.OrElse(""),
		LastName:       user.LastName.OrElse(""),
		Phone:          user.Phone.OrElse(""),
		CaseIDs:        user.CaseIDs.OrElse([]primitive.ObjectID{}),
		TeamID:         user.TeamID.Value,
		AgentIDs:       user.AgentIDs.OrElse([]primitive.ObjectID{}),
		SubscriptionID: user.SubscriptionID.Value,
	}
	return r.userDAO.CreateUser(ctx, userModel)
}

// UpdateUser updates an existing user.
func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, userID primitive.ObjectID, user map[string]interface{}) (*mongo.UpdateResult, error) {
	return r.userDAO.UpdateUser(ctx, userID, user)
}

// toUserResponse converts a User model to a UserResponse DTO.
func (r *UserRepositoryImpl) toUserResponse(user *models.User) *dtos.UserResponse {
	return &dtos.UserResponse{
		ID:             helpers.NewNullable(user.ID),
		Image:          helpers.NewNullable(user.Image),
		Email:          helpers.NewNullable(user.Email),
		FirstName:      helpers.NewNullable(user.FirstName),
		LastName:       helpers.NewNullable(user.LastName),
		Phone:          helpers.NewNullable(user.Phone),
		CaseIDs:        helpers.NewNullable(user.CaseIDs),
		TeamID:         helpers.NewNullable(user.TeamID),
		AgentIDs:       helpers.NewNullable(user.AgentIDs),
		SubscriptionID: helpers.NewNullable(user.SubscriptionID),
	}
}
