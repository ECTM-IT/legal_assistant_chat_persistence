package service

import (
	"context"

	dto "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	repository "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
)

// UserService defines the operations available on a user
type UserService interface {
	GetUserByUserID(ctx context.Context, userID string) (*dto.UserResponse, error)
	GetUserByCaseID(ctx context.Context, caseID string) (*dto.UserResponse, error)
	DeleteUserByID(ctx context.Context, userID string) error
	CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*dto.UserResponse, error)
	// ... potentially more service methods (e.g., for updates)
}

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of the user service
func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: repo}
}

// GetUserDetails retrieves details of a user
func (s *UserServiceImpl) GetUserDetails(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserById(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByCaseID retrieves a user associated with a given case ID
func (s *UserServiceImpl) GetUserByCaseID(ctx context.Context, caseID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByCasesId(caseID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*dto.UserResponse, error) {
	err := s.userRepo.SaveUser(user)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.userRepo.FindUserById(user.Firebase_id)
	if err != nil {
		return nil, err // Handle error during new user retrieval
	}
	return createdUser, nil
}

// DeleteUserByID deletes an existing user by identifier
func (s *UserServiceImpl) DeleteUserByID(ctx context.Context, userID string) error {
	return s.userRepo.DeleteUser(userID)
}
