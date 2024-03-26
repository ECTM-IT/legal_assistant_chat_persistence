// dao/user.go
package daos

import (
	"context"

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

func (dao *UserDAO) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	ctx := context.Background()

	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetUserByCaseID(caseID primitive.ObjectID) (*models.User, error) {
	ctx := context.Background()

	var user models.User
	err := dao.collection.FindOne(ctx, bson.M{"cases": caseID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) CreateUser(user *models.User) (*models.User, error) {
	ctx := context.Background()

	_, err := dao.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDAO) UpdateUser(id primitive.ObjectID, user *models.User) error {
	ctx := context.Background()

	_, err := dao.collection.ReplaceOne(ctx, bson.M{"_id": id}, user)
	return err
}

func (dao *UserDAO) DeleteUser(id primitive.ObjectID) error {
	ctx := context.Background()

	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (dao *UserDAO) GetAllUsers() ([]*models.User, error) {
	ctx := context.Background()

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
