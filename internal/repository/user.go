package repository

import (
	"errors"

	"github.com/shy-robin/gochat/internal/db"
	"github.com/shy-robin/gochat/internal/model"
	"gorm.io/gorm"
)

// Repository 层负责隔离数据库操作

type UserRepository struct {
}

// 等效于 var UserRepo = new(UserRepository)
// 但是这种方式可以自定义初始化参数
var UserRepo = &UserRepository{}

func (this *UserRepository) FindByUsername(username string) (*model.User, error) {
	db := db.GetDB()
	user := &model.User{}

	result := db.Where("username = ?", username).First(user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, result.Error
}

func (this *UserRepository) CreateUser(user *model.User) error {
	db := db.GetDB()
	result := db.Create(user)

	return result.Error
}
