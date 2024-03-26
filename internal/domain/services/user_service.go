package services

import (
	dto "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserService defines the operations available on a user
type UserService interface {
	GetUserByID(userID string) (*dto.UserResponse, error)
	GetUserByCaseID(caseID string) (*dto.UserResponse, error)
	DeleteUserByID(userID string) error
	CreateUser(user *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(userID string, user *dto.UpdateUserRequest) error
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
func (s *UserServiceImpl) GetUserByID(userID primitive.ObjectID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByCaseID retrieves a user associated with a given case ID
func (s *UserServiceImpl) GetUserByCaseID(caseID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByCasesID(caseID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user
func (s *UserServiceImpl) CreateUser(user *dto.CreateUserRequest) (*dto.UserResponse, error) {
	createdUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Retrieve the created user
	newUser, err := s.userRepo.FindUserByID(createdUser.ID)
	if err != nil {
		return nil, err // Handle error during new user retrieval
	}
	return newUser, nil
}

// UpdateUser updates an existing user
func (s *UserServiceImpl) UpdateUser(userID string, user *dto.UpdateUserRequest) error {
	return s.userRepo.UpdateUser(userID, user)
}

// DeleteUserByID deletes an existing user by identifier
func (s *UserServiceImpl) DeleteUserByID(userID string) error {
	return s.userRepo.DeleteUser(userID)
}
