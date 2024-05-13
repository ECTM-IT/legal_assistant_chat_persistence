package services

import (
	"context"

	dto "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService defines the operations available on a user
type UserService interface {
	GetUserByID(ctx context.Context, userID primitive.ObjectID) (*dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	GetUserByCaseID(ctx context.Context, caseID string) (*dto.UserResponse, error)
	DeleteUserByID(ctx context.Context, userID string) error
	CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, userID string, user *dto.UpdateUserRequest) (*mongo.UpdateResult, error)
}

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	userRepo *repositories.UserRepositoryImpl
}

// NewUserService creates a new instance of the user service
func NewUserService(repo *repositories.UserRepositoryImpl) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: repo,
	}
}

// GetUserByID retrieves details of a user by ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByCaseID retrieves a user associated with a given case ID
func (s *UserServiceImpl) GetUserByCaseID(ctx context.Context, caseID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByCasesID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*dto.UserResponse, error) {
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

// UpdateUser updates an existing user
func (s *UserServiceImpl) UpdateUser(ctx context.Context, userID string, user *dto.UpdateUserRequest) (*mongo.UpdateResult, error) {
	return s.userRepo.UpdateUser(ctx, userID, user)
}

// DeleteUserByID deletes an existing user by identifier
func (s *UserServiceImpl) DeleteUserByID(ctx context.Context, userID string) error {
	return s.userRepo.DeleteUser(ctx, userID)
}
