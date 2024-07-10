package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
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
	UpdateUser(ctx context.Context, userID primitive.ObjectID, user *dtos.UpdateUserRequest) (*dtos.UserResponse, error)
}

// UserServiceImpl implements the UserService interface.
type UserServiceImpl struct {
	userRepo *repositories.UserRepositoryImpl
	mapper   *mappers.UserConversionServiceImpl
}

// NewUserService creates a new instance of the user service.
func NewUserService(repo *repositories.UserRepositoryImpl, mapper *mappers.UserConversionServiceImpl) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: repo,
		mapper:   mapper,
	}
}

// GetUserByID retrieves details of a user by ID.
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("User not found", "user_not_found")
		}
		return nil, errors.NewDatabaseError("Failed to get user", "get_user_failed")
	}
	return s.mapper.UserToDTO(user), nil
}

// GetUserByEmail retrieves a user by their email.
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("User not found", "user_not_found")
		}
		return nil, errors.NewDatabaseError("Failed to get user", "get_user_failed")
	}
	return s.mapper.UserToDTO(user), nil
}

// GetUserByCaseID retrieves a user associated with a given case ID.
func (s *UserServiceImpl) GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindUserByCaseID(ctx, caseID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewNotFoundError("User not found", "user_not_found")
		}
		return nil, errors.NewDatabaseError("Failed to get user", "get_user_failed")
	}
	return s.mapper.UserToDTO(user), nil
}

// CreateUser creates a new user.
func (s *UserServiceImpl) CreateUser(ctx context.Context, userDTO *dtos.CreateUserRequest) (*dtos.UserResponse, error) {
	user, err := s.mapper.DTOToUser(userDTO)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create user", "create_user_failed")
	}

	return s.mapper.UserToDTO(createdUser), nil
}

// UpdateUser updates an existing user.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, userID primitive.ObjectID, userDTO *dtos.UpdateUserRequest) (*dtos.UserResponse, error) {
	updateFields := s.mapper.UpdateUserFieldsToMap(*userDTO)

	updatedUser, err := s.userRepo.UpdateUser(ctx, userID, updateFields)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to update user", "update_user_failed")
	}

	return s.mapper.UserToDTO(updatedUser), nil
}

// DeleteUserByID deletes an existing user by ID.
func (s *UserServiceImpl) DeleteUserByID(ctx context.Context, userID primitive.ObjectID) error {
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return errors.NewDatabaseError("Failed to delete user", "delete_user_failed")
	}
	return nil
}
