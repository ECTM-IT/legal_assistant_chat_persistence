package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/security"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
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
	logger   logs.Logger
	encrypt  security.EncryptionService
}

// NewUserService creates a new instance of the user service.
func NewUserService(repo *repositories.UserRepositoryImpl, mapper *mappers.UserConversionServiceImpl, logger logs.Logger, encrypt security.EncryptionService) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: repo,
		mapper:   mapper,
		logger:   logger,
		encrypt:  encrypt,
	}
}

// GetUserByID retrieves details of a user by ID.
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*dtos.UserResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve user by ID")
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warn("User not found")
			return nil, errors.NewNotFoundError("User not found", "user_not_found")
		}
		s.logger.Error("Service Level: Failed to get user", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get user", "get_user_failed")
	}

	// Decrypt name and email
	decryptedName, err := s.encrypt.Decrypt(user.EncryptedName)
	if err != nil {
		s.logger.Error("Service Level: Failed to decrypt name", err)
		return nil, errors.NewDatabaseError("failed to decrypt name", "decrypt error")
	}

	decryptedEmail, err := s.encrypt.Decrypt(user.EncryptedEmail)
	if err != nil {
		s.logger.Error("Service Level: Failed to decrypt email", err)
		return nil, errors.NewDatabaseError("failed to decrypt email", "encryption_error")
	}

	user.EncryptedName = decryptedName
	user.EncryptedEmail = decryptedEmail

	s.logger.Info("Service Level: Successfully retrieved user by ID")
	return s.mapper.UserToDTO(user), nil
}

// GetUserByEmail retrieves a user by their email.
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve user by email")

	// Encrypt the email to match the stored encrypted email
	encryptedEmail, err := s.encrypt.Encrypt(email)
	if err != nil {
		s.logger.Error("Service Level: Failed to encrypt email for search", err)
		return nil, errors.NewDatabaseError("failed to encrypt email", "encrypt error")
	}

	user, err := s.userRepo.FindUserByEmail(ctx, encryptedEmail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warn("User not found")
			return nil, errors.NewNotFoundError("User not found", "user_not_found")
		}
		s.logger.Error("Service Level: Failed to get user", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get user", "get_user_failed")
	}

	// Decrypt name and email
	decryptedName, err := s.encrypt.Decrypt(user.EncryptedName)
	if err != nil {
		s.logger.Error("Service Level: Failed to decrypt name", err)
		return nil, errors.NewDatabaseError("failed to decrypt name", "encryption_error")
	}

	decryptedEmail, err := s.encrypt.Decrypt(user.EncryptedEmail)
	if err != nil {
		s.logger.Error("Service Level: Failed to decrypt email", err)
		return nil, errors.NewDatabaseError("failed to decrypt email", "encryption_error")
	}

	user.EncryptedName = decryptedName
	user.EncryptedEmail = decryptedEmail

	s.logger.Info("Service Level: Successfully retrieved user by email")
	return s.mapper.UserToDTO(user), nil
}

// GetUserByCaseID retrieves a user associated with a given case ID.
func (s *UserServiceImpl) GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*dtos.UserResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve user by case ID")
	user, err := s.userRepo.FindUserByCaseID(ctx, caseID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Warn("User not found")
			return nil, errors.NewNotFoundError("User not found", "user_not_found")
		}
		s.logger.Error("Service Level: Failed to get user", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to get user", "get_user_failed")
	}
	s.logger.Info("Service Level: Successfully retrieved user by case ID")
	return s.mapper.UserToDTO(user), nil
}

// CreateUser creates a new user with encrypted name and email.
func (s *UserServiceImpl) CreateUser(ctx context.Context, userDTO *dtos.CreateUserRequest) (*dtos.UserResponse, error) {
	s.logger.Info("Service Level: Attempting to create new user")

	mappedUser, err := s.mapper.DTOToUser(userDTO)
	if err != nil {
		s.logger.Error("Service Level: Failed to convert DTO to user", err)
		return nil, errors.NewDatabaseError("failed to convert DTO to user", "conversion_error")
	}
	// Encrypt username and email
	encryptedName, err := s.encrypt.Encrypt(mappedUser.FirstName)
	if err != nil {
		s.logger.Error("Service Level: Failed to encrypt name", err)
		return nil, errors.NewDatabaseError("failed to encrypt name", "encryption_error")
	}

	encryptedEmail, err := s.encrypt.Encrypt(mappedUser.Email)
	if err != nil {
		s.logger.Error("Service Level: Failed to encrypt email", err)
		return nil, errors.NewDatabaseError("failed to encrypt email", "encryption_error")
	}

	mappedUser.EncryptedName = encryptedName
	mappedUser.EncryptedEmail = encryptedEmail
	mappedUser.Email = ""

	createdUser, err := s.userRepo.CreateUser(ctx, mappedUser)
	if err != nil {
		s.logger.Error("Service Level: Failed to create user", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to create user", "create_user_failed")
	}

	s.logger.Info("Service Level: Successfully created new user")
	return s.mapper.UserToDTO(createdUser), nil
}

// UpdateUser updates an existing user with encrypted fields.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, userID primitive.ObjectID, userDTO *dtos.UpdateUserRequest) (*dtos.UserResponse, error) {
	s.logger.Info("Service Level: Attempting to update user")

	updateFields := make(map[string]interface{})

	if userDTO.EncryptedName.Present {
		encryptedName, err := s.encrypt.Encrypt(userDTO.EncryptedName.Value)
		if err != nil {
			s.logger.Error("Service Level: Failed to encrypt name during update", err)
			return nil, errors.NewDatabaseError("failed to encrypt name", "encryption_error")
		}
		updateFields["encrypted_name"] = encryptedName
	}

	if userDTO.EncryptedEmail.Present {
		encryptedEmail, err := s.encrypt.Encrypt(userDTO.EncryptedEmail.Value)
		if err != nil {
			s.logger.Error("Service Level: Failed to encrypt email during update", err)
			return nil, errors.NewDatabaseError("failed to encrypt email", "encryption_error")
		}
		updateFields["encrypted_email"] = encryptedEmail
	}

	// Process other fields
	otherUpdates := s.mapper.UpdateUserFieldsToMap(*userDTO)
	for k, v := range otherUpdates {
		updateFields[k] = v
	}

	updatedUser, err := s.userRepo.UpdateUser(ctx, userID, updateFields)
	if err != nil {
		s.logger.Error("Service Level: Failed to update user", err)
		return nil, errors.NewDatabaseError("Service Level: Failed to update user", "update_user_failed")
	}

	// Decrypt name and email
	decryptedName, err := s.encrypt.Decrypt(updatedUser.EncryptedName)
	if err != nil {
		s.logger.Error("Service Level: Failed to decrypt name after update", err)
		return nil, errors.NewDatabaseError("failed to decrypt name after update", "encrypt_error")
	}

	decryptedEmail, err := s.encrypt.Decrypt(updatedUser.EncryptedEmail)
	if err != nil {
		s.logger.Error("Service Level: Failed to decrypt email after update", err)
		return nil, errors.NewDatabaseError("failed to decrypt email after update", "encryption_error")
	}

	updatedUser.EncryptedName = decryptedName
	updatedUser.EncryptedEmail = decryptedEmail

	s.logger.Info("Service Level: Successfully updated user")
	return s.mapper.UserToDTO(updatedUser), nil
}

// DeleteUserByID deletes an existing user by ID.
func (s *UserServiceImpl) DeleteUserByID(ctx context.Context, userID primitive.ObjectID) error {
	s.logger.Info("Service Level: Attempting to delete user")
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		s.logger.Error("Service Level: Failed to delete user", err)
		return errors.NewDatabaseError("Service Level: Failed to delete user", "delete_user_failed")
	}
	s.logger.Info("Service Level: Successfully deleted user")
	return nil
}
