// dao/user.go
package daos

import (
	"context"
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDAO struct {
	collection *mongo.Collection
}

func NewUserDAO(db *mongo.Database) *UserDAO {
	return &UserDAO{
		collection: db.Collection("users"),
	}
}

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

func (dao *UserDAO) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := dao.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDAO) UpdateUser(ctx context.Context, id primitive.ObjectID, user *models.User) (*mongo.UpdateResult, error) {
	return dao.collection.ReplaceOne(ctx, bson.M{"_id": id}, user)
}

func (dao *UserDAO) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

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
