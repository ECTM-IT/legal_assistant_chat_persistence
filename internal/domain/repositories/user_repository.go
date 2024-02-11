package mongo

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	dto "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
)

type UserRepository interface {
	FindUserById(userId string) (*dto.UserResponse, error)
	FindUserByCasesId(casesId string) (*dto.UserResponse, error)
	TotalUsers() (int64, error)
	DeleteUser(userid string) error
	SaveUser(user *dto.CreateUserRequest) error
}

type UserRepositoryImpl struct {
	userDao daos.UserDao
}

func NewUserRepository(userDao daos.UserDao) *UserRepositoryImpl {
	return &UserRepositoryImpl{userDao: userDao}
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	return r.userDao.SaveUser(user)
}
