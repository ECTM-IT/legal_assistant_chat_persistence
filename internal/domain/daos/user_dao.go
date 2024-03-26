// dao/user.go
package daos

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetUserByCaseID(ctx context.Context, caseID primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"cases": caseID}).Decode(&user)
	if err != nil {
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

func (dao *UserDAO) UpdateUser(ctx context.Context, id primitive.ObjectID, user *models.User) error {
	_, err := dao.collection.ReplaceOne(ctx, bson.M{"_id": id}, user)
	return err
}

func (dao *UserDAO) DeleteUser(ctx context.Context, id primitive.ObjectID) error {

	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (dao *UserDAO) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := dao.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
