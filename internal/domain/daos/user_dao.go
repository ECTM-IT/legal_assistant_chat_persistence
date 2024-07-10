package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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
}

// NewUserDAO creates a new UserDAO
func NewUserDAO(db *mongo.Database) *UserDAO {
	return &UserDAO{
		collection: db.Collection("users"),
	}
}

// GetUserByID retrieves a user by their ID from the database
func (dao *UserDAO) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email from the database
func (dao *UserDAO) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByCaseID retrieves a user by their case ID from the database
func (dao *UserDAO) GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"cases": caseID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func (dao *UserDAO) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := dao.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user by their ID in the database
func (dao *UserDAO) UpdateUser(ctx context.Context, id primitive.ObjectID, user map[string]interface{}) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": user}
	result, err := dao.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteUser deletes a user by their ID from the database
func (dao *UserDAO) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

// GetAllUsers retrieves all users from the database
func (dao *UserDAO) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
