package mappers

import (
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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

type UserConversionServiceImpl struct{}

func NewUserConversionService() *UserConversionServiceImpl {
	return &UserConversionServiceImpl{}
}

func (s *UserConversionServiceImpl) UserToDTO(user *models.User) *dtos.UserResponse {
	if user == nil {
		return nil
	}

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

func (s *UserConversionServiceImpl) UsersToDTO(users []models.User) []dtos.UserResponse {
	userResponses := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.UserToDTO(&user)
	}
	return userResponses
}

func (s *UserConversionServiceImpl) DTOToUser(userDTO *dtos.CreateUserRequest) (*models.User, error) {
	if userDTO == nil {
		return nil, errors.New("user DTO cannot be nil")
	}

	if !userDTO.Email.Present {
		return nil, errors.New("email is required")
	}

	return &models.User{
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
	}, nil
}

func (s *UserConversionServiceImpl) UpdateUserFieldsToMap(updateRequest dtos.UpdateUserRequest) map[string]interface{} {
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

	return updateFields
}
