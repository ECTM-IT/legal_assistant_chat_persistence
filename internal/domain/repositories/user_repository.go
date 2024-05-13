package repositories

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindUserByID(ctx context.Context, userID helpers.Nullable[primitive.ObjectID]) (*dtos.UserResponse, error)
	FindUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error)
	FindUserByCasesID(ctx context.Context, casesID helpers.Nullable[string]) (*dtos.UserResponse, error)
	TotalUsers(ctx context.Context) ([]*models.User, error)
	DeleteUser(ctx context.Context, userID helpers.Nullable[string]) error
	CreateUser(ctx context.Context, user *dtos.CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, userID helpers.Nullable[string], user *dtos.UpdateUserRequest) (*mongo.UpdateResult, error)
}

type UserRepositoryImpl struct {
	userDAO *daos.UserDAO
}

func NewUserRepository(userDAO *daos.UserDAO) *UserRepositoryImpl {
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

func (r *UserRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error) {
	user, err := r.userDAO.GetUserByEmail(ctx, email)

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
		Image:          user.Image.OrElse(""),
		Email:          user.Email.OrElse(""),
		FirstName:      user.FirstName.OrElse(""),
		LastName:       user.LastName.OrElse(""),
		Phone:          user.Phone.OrElse(""),
		CaseIDs:        user.CaseIDs.OrElse([]primitive.ObjectID{}),
		TeamID:         user.TeamID.Val,
		AgentIDs:       user.AgentIDs.OrElse([]primitive.ObjectID{}),
		SubscriptionID: user.SubscriptionID.Val,
	}
	return r.userDAO.CreateUser(ctx, userModel)
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, userID string, user *dtos.UpdateUserRequest) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	userModel := &models.User{}

	if user.Image.Valid {
		userModel.Image = user.Image.Val
	}
	if user.Email.Valid {
		userModel.Email = user.Email.Val
	}
	if user.FirstName.Valid {
		userModel.FirstName = user.FirstName.Val
	}
	if user.LastName.Valid {
		userModel.LastName = user.LastName.Val
	}
	if user.Phone.Valid {
		userModel.Phone = user.Phone.Val
	}

	if user.CaseIDs.Valid {
		caseIDs := make([]primitive.ObjectID, len(user.CaseIDs.Val))
		for i, caseID := range user.CaseIDs.Val {
			caseObjectID, err := primitive.ObjectIDFromHex(caseID)
			if err != nil {
				return nil, err
			}
			caseIDs[i] = caseObjectID
		}
		userModel.CaseIDs = caseIDs
	}

	if user.TeamID.Valid {
		teamID, err := primitive.ObjectIDFromHex(user.TeamID.Val)
		if err != nil {
			return nil, err
		}
		userModel.TeamID = teamID
	}

	if user.AgentIDs.Valid {
		agentIDs := make([]primitive.ObjectID, len(user.AgentIDs.Val))
		for i, agentID := range user.AgentIDs.Val {
			agentObjectID, err := primitive.ObjectIDFromHex(agentID)
			if err != nil {
				return nil, err
			}
			agentIDs[i] = agentObjectID
		}
		userModel.AgentIDs = agentIDs
	}

	if user.SubscriptionID.Valid {
		subscriptionID, err := primitive.ObjectIDFromHex(user.SubscriptionID.Val)
		if err != nil {
			return nil, err
		}
		userModel.SubscriptionID = subscriptionID
	}

	return r.userDAO.UpdateUser(ctx, objectID, userModel)
}

func (r *UserRepositoryImpl) toUserResponse(user *models.User) *dtos.UserResponse {
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
