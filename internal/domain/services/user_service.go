package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService defines the operations available on a user.
type UserService interface {
	GetUserByID(ctx context.Context, userID primitive.ObjectID) (*dtos.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error)
	GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*dtos.UserResponse, error)
	DeleteUserByID(ctx context.Context, userID primitive.ObjectID) error
	CreateUser(ctx context.Context, user *dtos.CreateUserRequest) (*dtos.UserResponse, error)
	UpdateUser(ctx context.Context, userID primitive.ObjectID, user *dtos.UpdateUserRequest) (*mongo.UpdateResult, error)
}

// userServiceImpl implements the UserService interface.
type UserServiceImpl struct {
	userRepo *repositories.UserRepositoryImpl
}

// NewUserService creates a new instance of the user service.
func NewUserService(repo *repositories.UserRepositoryImpl) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: repo,
	}
}

// GetUserByID retrieves details of a user by ID.
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by their email.
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByCaseID retrieves a user associated with a given case ID.
func (s *UserServiceImpl) GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindUserByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user.
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *dtos.CreateUserRequest) (*dtos.UserResponse, error) {
	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	// Retrieve the created user
	newUser, err := s.userRepo.FindUserByID(ctx, createdUser.ID)
	if err != nil {
		return nil, err // Handle error during new user retrieval
	}
	return newUser, nil
}

// UpdateUser updates an existing user.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, userID primitive.ObjectID, user map[string]interface{}) (*dtos.UserResponse, error) {
	_, err := s.userRepo.UpdateUser(ctx, userID, user)
	if err != nil {
		return nil, err
	}
	return s.userRepo.FindUserByID(ctx, userID)
}

// DeleteUserByID deletes an existing user by ID.
func (s *UserServiceImpl) DeleteUserByID(ctx context.Context, userID primitive.ObjectID) error {
	return s.userRepo.DeleteUser(ctx, userID)
}
