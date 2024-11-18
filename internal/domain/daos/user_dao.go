package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserDAOInterface defines the interface for the UserDAO
type UserDAOInterface interface {
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, user map[string]interface{}) (*mongo.UpdateResult, error)
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}

// UserDAO implements the UserDAOInterface
type UserDAO struct {
	collection *mongo.Collection
	logger     logs.Logger
}

// NewUserDAO creates a new UserDAO
func NewUserDAO(db *mongo.Database, logger logs.Logger) *UserDAO {
	return &UserDAO{
		collection: db.Collection("users"),
		logger:     logger,
	}
}

// GetUserByID retrieves a user by their ID from the database
func (dao *UserDAO) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve user by ID")
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("User not found")
			return nil, errors.New("user not found")
		}
		dao.logger.Error("DAO Level: Failed to retrieve user", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved user")
	return &user, nil
}

// GetUserByEmail retrieves a user by their email from the database
func (dao *UserDAO) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve user by email")
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("User not found")
			return nil, errors.New("user not found")
		}
		dao.logger.Error("DAO Level: Failed to retrieve user", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved user")
	return &user, nil
}

// GetUserByCaseID retrieves a user by their case ID from the database
func (dao *UserDAO) GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*models.User, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve user by case ID")
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"cases": caseID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Warn("User not found")
			return nil, errors.New("user not found")
		}
		dao.logger.Error("DAO Level: Failed to retrieve user", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved user")
	return &user, nil
}

// CreateUser creates a new user in the database
func (dao *UserDAO) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	dao.logger.Info("DAO Level: Attempting to create new user")
	_, err := dao.collection.InsertOne(ctx, user)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to create user", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully created new user")
	return user, nil
}

// UpdateUser updates a user by their ID in the database
func (dao *UserDAO) UpdateUser(ctx context.Context, id primitive.ObjectID, user map[string]interface{}) (*mongo.UpdateResult, error) {
	dao.logger.Info("DAO Level: Attempting to update user")
	update := bson.M{"$set": user}
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		dao.logger.Error("DAO Level: Failed to update user", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully updated user")
	return result, nil
}

// DeleteUser deletes a user by their ID from the database
func (dao *UserDAO) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	dao.logger.Info("DAO Level: Attempting to delete user")
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to delete user", err)
		return err
	}
	dao.logger.Info("DAO Level: Successfully deleted user")
	return nil
}

// GetAllUsers retrieves all users from the database
func (dao *UserDAO) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	dao.logger.Info("DAO Level: Attempting to retrieve all users")
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("DAO Level: Failed to retrieve users", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		dao.logger.Error("DAO Level: Failed to decode users", err)
		return nil, err
	}
	dao.logger.Info("DAO Level: Successfully retrieved all users")
	return users, nil
}
