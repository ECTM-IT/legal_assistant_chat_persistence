package daos

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserDao struct {
	db *mongo.Database
}

func NewUserDao(db *mongo.Database) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) collection() *mongo.Collection {
	return dao.db.Collection("users")
}

// FindUserById finds the user with the provided user_id.
func (dao *UserDao) FindUserById(userId string) (*models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := dao.collection().FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindUserByCasesId finds the user with the provided cases_id.
func (dao *UserDao) FindUserByCasesId(casesId string) (*models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := dao.collection().FindOne(ctx, bson.M{"cases_id": casesId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// TotalUsers returns the number of existing user records.
func (dao *UserDao) TotalUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	total, err := dao.collection().CountDocuments(ctx, bson.M{})
	return total, err
}

// DeleteUser deletes the provided User model.
func (dao *UserDao) DeleteUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := dao.collection().DeleteOne(ctx, bson.M{"_id": user.ID})
	return err
}

// SaveUser upserts the provided User model.
func (dao *UserDao) SaveUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	_, err := dao.collection().UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user}, opts)
	return err
}
