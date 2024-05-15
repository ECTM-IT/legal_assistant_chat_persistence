package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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

// GetUserByID retrieves a user by their ID
func (dao *UserDAO) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Error("User not found", err, zap.String("userID", id.Hex()))
			return nil, errors.New("user not found")
		}
		dao.logger.Error("Error retrieving user by ID", err, zap.String("userID", id.Hex()))
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func (dao *UserDAO) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Error("User not found", err, zap.String("email", email))
			return nil, errors.New("user not found")
		}
		dao.logger.Error("Error retrieving user by email", err, zap.String("email", email))
		return nil, err
	}
	return &user, nil
}

// GetUserByCaseID retrieves a user by their case ID
func (dao *UserDAO) GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"cases": caseID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			dao.logger.Error("User not found", err, zap.String("caseID", caseID.Hex()))
			return nil, errors.New("user not found")
		}
		dao.logger.Error("Error retrieving user by case ID", err, zap.String("caseID", caseID.Hex()))
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user
func (dao *UserDAO) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := dao.collection.InsertOne(ctx, user)
	if err != nil {
		dao.logger.Error("Error creating user", err)
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user by their ID
func (dao *UserDAO) UpdateUser(ctx context.Context, id primitive.ObjectID, user map[string]interface{}) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": user}
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		dao.logger.Error("Error updating user", err, zap.String("userID", id.Hex()))
		return nil, err
	}
	return result, nil
}

// DeleteUser deletes a user by their ID
func (dao *UserDAO) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		dao.logger.Error("Error deleting user", err, zap.String("userID", id.Hex()))
		return err
	}
	return nil
}

// GetAllUsers retrieves all users
func (dao *UserDAO) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		dao.logger.Error("Error retrieving all users", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		dao.logger.Error("Error decoding users", err)
		return nil, err
	}

	return users, nil
}
