package mongo

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *UserRepository {
	return &UserRepository{
		db: client.Database(dbName).Collection(collectionName),
	}
}

// Example of a method that uses the MongoDB connection
func (r *UserRepository) CreateUser(user *daos.UserDao) error {
	// _, err := r.db.InsertOne(context.Background(), user)
	err := user.SaveUser(user.SaveUser())
	return err
}
