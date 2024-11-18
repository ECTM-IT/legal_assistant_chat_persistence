package mappers

import (
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserConversionService interface {
	UserToDTO(user *models.User) *dtos.UserResponse
	UsersToDTO(users []models.User) []dtos.UserResponse
	DTOToUser(userDTO *dtos.CreateUserRequest) (*models.User, error)
	UpdateUserFieldsToMap(updateRequest dtos.UpdateUserRequest) map[string]interface{}
	ObjectIDsToDTO(ids []primitive.ObjectID) []string
	DTOToObjectIDs(idStrings []string) ([]primitive.ObjectID, error)
}

type UserConversionServiceImpl struct {
	logger logs.Logger
}

func NewUserConversionService(logger logs.Logger) *UserConversionServiceImpl {
	return &UserConversionServiceImpl{
		logger: logger,
	}
}

func (s *UserConversionServiceImpl) UserToDTO(user *models.User) *dtos.UserResponse {
	s.logger.Info("Converting User to DTO")
	if user == nil {
		s.logger.Warn("Attempted to convert nil User to DTO")
		return nil
	}

	dto := &dtos.UserResponse{
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
	s.logger.Info("Successfully converted User to DTO")
	return dto
}

func (s *UserConversionServiceImpl) UsersToDTO(users []models.User) []dtos.UserResponse {
	s.logger.Info("Converting multiple Users to DTOs")
	userResponses := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.UserToDTO(&user)
	}
	s.logger.Info("Successfully converted multiple Users to DTOs")
	return userResponses
}

func (s *UserConversionServiceImpl) DTOToUser(userDTO *dtos.CreateUserRequest) (*models.User, error) {
	s.logger.Info("Converting DTO to User")
	if userDTO == nil {
		s.logger.Error("Failed to convert DTO to User: user DTO cannot be nil", errors.New("user DTO cannot be nil"))
		return nil, errors.New("user DTO cannot be nil")
	}

	if !userDTO.Email.Present {
		s.logger.Error("Failed to convert DTO to User: email is required", errors.New("email is required"))
		return nil, errors.New("email is required")
	}

	user := &models.User{
		ID:             primitive.NewObjectID(),
		Image:          userDTO.Image.OrElse(""),
		Email:          userDTO.Email.Value,
		FirstName:      userDTO.FirstName.OrElse(""),
		LastName:       userDTO.LastName.OrElse(""),
		Phone:          userDTO.Phone.OrElse(""),
		CaseIDs:        userDTO.CaseIDs.OrElse([]primitive.ObjectID{}),
		TeamID:         userDTO.TeamID.OrElse(primitive.NilObjectID),
		AgentIDs:       userDTO.AgentIDs.OrElse([]primitive.ObjectID{}),
		SubscriptionID: userDTO.SubscriptionID.OrElse(primitive.NilObjectID),
	}
	s.logger.Info("Successfully converted DTO to User")
	return user, nil
}

func (s *UserConversionServiceImpl) UpdateUserFieldsToMap(updateRequest dtos.UpdateUserRequest) map[string]interface{} {
	s.logger.Info("Converting UpdateUserRequest to map")
	updateFields := make(map[string]interface{})

	if updateRequest.Image.Present {
		updateFields["image"] = updateRequest.Image.Value
	}
	if updateRequest.Email.Present {
		updateFields["email"] = updateRequest.Email.Value
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

func (s *UserConversionServiceImpl) ObjectIDsToDTO(ids []primitive.ObjectID) []string {
	s.logger.Info("Converting ObjectIDs to DTO")
	dtoIDs := make([]string, len(ids))
	for i, id := range ids {
		dtoIDs[i] = id.Hex()
	}
	s.logger.Info("Successfully converted ObjectIDs to DTO")
	return dtoIDs
}

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
