package repositories

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	FindUserByID(userID primitive.ObjectID) (*dtos.UserResponse, error)
	FindUserByCasesID(casesID string) (*dtos.UserResponse, error)
	TotalUsers() ([]*models.User, error)
	DeleteUser(userID string) error
	CreateUser(user *dtos.CreateUserRequest) (*models.User, error)
	UpdateUser(userID string, user *dtos.UpdateUserRequest) error
}

type UserRepositoryImpl struct {
	userDAO *daos.UserDAO
}

func NewUserRepository(userDAO *daos.UserDAO) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		userDAO: userDAO,
	}
}

func (r *UserRepositoryImpl) FindUserByID(userID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := r.userDAO.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return r.toUserResponse(user), nil
}

func (r *UserRepositoryImpl) FindUserByCasesID(casesID string) (*dtos.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(casesID)
	if err != nil {
		return nil, err
	}

	user, err := r.userDAO.GetUserByCaseID(objectID)
	if err != nil {
		return nil, err
	}

	return r.toUserResponse(user), nil
}

func (r *UserRepositoryImpl) TotalUsers() ([]*models.User, error) {
	return r.userDAO.GetAllUsers()
}

func (r *UserRepositoryImpl) DeleteUser(userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	return r.userDAO.DeleteUser(objectID)
}

func (r *UserRepositoryImpl) CreateUser(user *dtos.CreateUserRequest) (*models.User, error) {
	userModel := &models.User{
		ID:             primitive.NewObjectID(),
		Image:          user.Image,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Phone:          user.Phone,
		CaseIDs:        user.CaseIDs,
		TeamID:         user.TeamID,
		AgentIDs:       user.AgentIDs,
		SubscriptionID: user.SubscriptionID,
	}
	return r.userDAO.CreateUser(userModel)
}

func (r *UserRepositoryImpl) UpdateUser(userID string, user *dtos.UpdateUserRequest) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	userModel := &models.User{}

	if user.Image != nil {
		userModel.Image = *user.Image
	}
	if user.Email != nil {
		userModel.Email = *user.Email
	}
	if user.FirstName != nil {
		userModel.FirstName = *user.FirstName
	}
	if user.LastName != nil {
		userModel.LastName = *user.LastName
	}
	if user.Phone != nil {
		userModel.Phone = *user.Phone
	}

	if len(user.CaseIDs) > 0 {
		caseIDs := make([]primitive.ObjectID, len(user.CaseIDs))
		for i, caseID := range user.CaseIDs {
			caseObjectID, err := primitive.ObjectIDFromHex(caseID)
			if err != nil {
				return err
			}
			caseIDs[i] = caseObjectID
		}
		userModel.CaseIDs = caseIDs
	}

	if user.TeamID != nil {
		teamID, err := primitive.ObjectIDFromHex(*user.TeamID)
		if err != nil {
			return err
		}
		userModel.TeamID = teamID
	}

	if len(user.AgentIDs) > 0 {
		agentIDs := make([]primitive.ObjectID, len(user.AgentIDs))
		for i, agentID := range user.AgentIDs {
			agentObjectID, err := primitive.ObjectIDFromHex(agentID)
			if err != nil {
				return err
			}
			agentIDs[i] = agentObjectID
		}
		userModel.AgentIDs = agentIDs
	}

	if user.SubscriptionID != nil {
		subscriptionID, err := primitive.ObjectIDFromHex(*user.SubscriptionID)
		if err != nil {
			return err
		}
		userModel.SubscriptionID = subscriptionID
	}

	return r.userDAO.UpdateUser(objectID, userModel)
}

func (r *UserRepositoryImpl) toUserResponse(user *models.User) *dtos.UserResponse {
	return &dtos.UserResponse{
		ID:             user.ID,
		Image:          user.Image,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Phone:          user.Phone,
		CaseIDs:        user.CaseIDs,
		TeamID:         user.TeamID,
		AgentIDs:       user.AgentIDs,
		SubscriptionID: user.SubscriptionID,
	}
}
