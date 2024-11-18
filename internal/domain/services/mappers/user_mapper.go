package mappers

import (
	"errors"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserConversionService handles conversions between User models and DTOs.
type UserConversionService interface {
	UserToDTO(user *models.User) *dtos.UserResponse
	UsersToDTO(users []models.User) []dtos.UserResponse
	DTOToUser(userDTO *dtos.CreateUserRequest) (*models.User, error)
	UpdateUserFieldsToMap(updateRequest dtos.UpdateUserRequest) map[string]interface{}
	ObjectIDsToDTO(ids []primitive.ObjectID) []string
	DTOToObjectIDs(idStrings []string) ([]primitive.ObjectID, error)
}

// UserConversionServiceImpl implements the UserConversionService interface.
type UserConversionServiceImpl struct {
	logger logs.Logger
}

// NewUserConversionService creates a new instance of UserConversionServiceImpl.
func NewUserConversionService(logger logs.Logger) *UserConversionServiceImpl {
	return &UserConversionServiceImpl{
		logger: logger,
	}
}

// UserToDTO converts a User model to a UserResponse DTO.
func (s *UserConversionServiceImpl) UserToDTO(user *models.User) *dtos.UserResponse {
	s.logger.Info("Converting User to DTO")
	if user == nil {
		s.logger.Warn("Attempted to convert nil User to DTO")
		return nil
	}

	dto := &dtos.UserResponse{
		ID:             helpers.NewNullable(user.ID),
		EncryptedName:  helpers.NewNullable(user.EncryptedName),
		EncryptedEmail: helpers.NewNullable(user.EncryptedEmail),
		Image:          helpers.NewNullable(user.Image),
		FirstName:      helpers.NewNullable(user.FirstName),
		LastName:       helpers.NewNullable(user.LastName),
		Phone:          helpers.NewNullable(user.Phone),
		CaseIDs:        helpers.NewNullable(user.CaseIDs),
		TeamID:         helpers.NewNullable(user.TeamID),
		AgentIDs:       helpers.NewNullable(user.AgentIDs),
		SubscriptionID: helpers.NewNullable(user.SubscriptionID),
		CreationDate:   helpers.NewNullable(user.CreationDate),
		LastEdit:       helpers.NewNullable(user.LastEdit),
		Share:          helpers.NewNullable(user.Share),
		IsArchived:     helpers.NewNullable(user.IsArchived),
	}
	s.logger.Info("Successfully converted User to DTO")
	return dto
}

// UsersToDTO converts a slice of User models to a slice of UserResponse DTOs.
func (s *UserConversionServiceImpl) UsersToDTO(users []models.User) []dtos.UserResponse {
	s.logger.Info("Converting multiple Users to DTOs")
	userResponses := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.UserToDTO(&user)
	}
	s.logger.Info("Successfully converted multiple Users to DTOs")
	return userResponses
}

// DTOToUser converts a CreateUserRequest DTO to a User model.
func (s *UserConversionServiceImpl) DTOToUser(userDTO *dtos.CreateUserRequest) (*models.User, error) {
	s.logger.Info("Converting DTO to User")
	if userDTO == nil {
		s.logger.Error("Failed to convert DTO to User: user DTO cannot be nil", errors.New("user DTO cannot be nil"))
		return nil, errors.New("user DTO cannot be nil")
	}

	if !userDTO.EncryptedEmail.Present {
		s.logger.Error("Failed to convert DTO to User: email is required", errors.New("email is required"))
		return nil, errors.New("email is required")
	}

	if !userDTO.EncryptedName.Present {
		s.logger.Error("Failed to convert DTO to User: name is required", errors.New("name is required"))
		return nil, errors.New("name is required")
	}

	user := &models.User{
		ID:             primitive.NewObjectID(),
		EncryptedName:  userDTO.EncryptedName.Value,
		EncryptedEmail: userDTO.EncryptedEmail.Value,
		Image:          userDTO.Image.OrElse(""),
		FirstName:      userDTO.FirstName.OrElse(""),
		LastName:       userDTO.LastName.OrElse(""),
		Phone:          userDTO.Phone.OrElse(""),
		CaseIDs:        userDTO.CaseIDs.OrElse([]primitive.ObjectID{}),
		TeamID:         userDTO.TeamID.OrElse(primitive.NilObjectID),
		AgentIDs:       userDTO.AgentIDs.OrElse([]primitive.ObjectID{}),
		SubscriptionID: userDTO.SubscriptionID.OrElse(primitive.NilObjectID),
		CreationDate:   userDTO.CreationDate.OrElse(time.Now()),
		LastEdit:       userDTO.LastEdit.OrElse(time.Now()),
		Share:          userDTO.Share.OrElse(false),
		IsArchived:     userDTO.IsArchived.OrElse(false),
	}
	s.logger.Info("Successfully converted DTO to User")
	return user, nil
}

// UpdateUserFieldsToMap converts an UpdateUserRequest DTO to a map for database updates.
func (s *UserConversionServiceImpl) UpdateUserFieldsToMap(updateRequest dtos.UpdateUserRequest) map[string]interface{} {
	s.logger.Info("Converting UpdateUserRequest to map")
	updateFields := make(map[string]interface{})

	if updateRequest.EncryptedName.Present {
		updateFields["encrypted_name"] = updateRequest.EncryptedName.Value
	}
	if updateRequest.EncryptedEmail.Present {
		updateFields["encrypted_email"] = updateRequest.EncryptedEmail.Value
	}
	if updateRequest.Image.Present {
		updateFields["image"] = updateRequest.Image.Value
	}
	if updateRequest.FirstName.Present {
		updateFields["first_name"] = updateRequest.FirstName.Value
	}
	if updateRequest.LastName.Present {
		updateFields["last_name"] = updateRequest.LastName.Value
	}
	if updateRequest.Phone.Present {
		updateFields["phone"] = updateRequest.Phone.Value
	}
	if updateRequest.CaseIDs.Present {
		updateFields["case_ids"] = updateRequest.CaseIDs.Value
	}
	if updateRequest.TeamID.Present {
		updateFields["team_id"] = updateRequest.TeamID.Value
	}
	if updateRequest.AgentIDs.Present {
		updateFields["agent_ids"] = updateRequest.AgentIDs.Value
	}
	if updateRequest.SubscriptionID.Present {
		updateFields["subscription_id"] = updateRequest.SubscriptionID.Value
	}

	s.logger.Info("Successfully converted UpdateUserRequest to map")
	return updateFields
}

// ObjectIDsToDTO converts a slice of ObjectIDs to a slice of their hexadecimal string representations.
func (s *UserConversionServiceImpl) ObjectIDsToDTO(ids []primitive.ObjectID) []string {
	s.logger.Info("Converting ObjectIDs to DTO")
	dtoIDs := make([]string, len(ids))
	for i, id := range ids {
		dtoIDs[i] = id.Hex()
	}
	s.logger.Info("Successfully converted ObjectIDs to DTO")
	return dtoIDs
}

// DTOToObjectIDs converts a slice of hexadecimal string representations to a slice of ObjectIDs.
func (s *UserConversionServiceImpl) DTOToObjectIDs(idStrings []string) ([]primitive.ObjectID, error) {
	s.logger.Info("Converting DTO to ObjectIDs")
	ids := make([]primitive.ObjectID, len(idStrings))
	for i, idString := range idStrings {
		id, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			s.logger.Error("Failed to convert DTO to ObjectID", err)
			return nil, err
		}
		ids[i] = id
	}
	s.logger.Info("Successfully converted DTO to ObjectIDs")
	return ids, nil
}
