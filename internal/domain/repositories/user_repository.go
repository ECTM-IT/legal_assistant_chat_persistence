package repositories

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	FindUserByID(ctx context.Context, userID primitive.ObjectID) (*dtos.UserResponse, error)
	FindUserByCasesID(ctx context.Context, casesID string) (*dtos.UserResponse, error)
	TotalUsers(ctx context.Context) (int64, error)
	DeleteUser(ctx context.Context, userID string) error
	CreateUser(ctx context.Context, user *dtos.CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, user *dtos.UpdateUserRequest) error
}

type UserRepositoryImpl struct {
	userDAO daos.UserDAO
}

func NewUserRepository(userDAO daos.UserDAO) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		userDAO: userDAO,
	}
}

func (r *UserRepositoryImpl) FindUserByID(ctx context.Context, userID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := r.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return r.toUserResponse(user), nil
}

func (r *UserRepositoryImpl) FindUserByCasesID(ctx context.Context, casesID string) (*dtos.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(casesID)
	if err != nil {
		return nil, err
	}

	user, err := r.userDAO.GetUserByCaseID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return r.toUserResponse(user), nil
}

func (r *UserRepositoryImpl) TotalUsers(ctx context.Context) ([]*models.User, error) {
	return r.userDAO.GetAllUsers(ctx)
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	return r.userDAO.DeleteUser(ctx, objectID)
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *dtos.CreateUserRequest) (*models.User, error) {
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
	return r.userDAO.CreateUser(ctx, userModel)
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, userID string, user *dtos.UpdateUserRequest) error {
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

	return r.userDAO.UpdateUser(ctx, objectID, userModel)
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
