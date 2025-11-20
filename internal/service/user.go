package service

import (
	"errors"

	"github.com/shy-robin/gochat/internal/db"
	"github.com/shy-robin/gochat/internal/model"
	"github.com/shy-robin/gochat/internal/repository"
	"gorm.io/gorm"
)

type UserService struct {
}

func (this *UserService) Register(user *model.User) error {
	existingUser, err := repository.UserRepo.FindByUsername(user.Username)

	if err != nil {
		return errors.New("查询用户失败")
	}

	if existingUser != nil {
		return errors.New("用户名已存在")
	}

	db := db.GetDB()

	// txErr -> Transaction Error (事物错误)
	// 确保所有 DB 操作要么全部成功，要么全部失败
	txErr := db.Transaction(func(tx *gorm.DB) error {
		if err := repository.UserRepo.CreateUser(user); err != nil {
			// 如果创建失败，返回错误，GORM 自动回滚 (ROLLBACK)
			return err
		}
		// 事务成功，GORM 自动提交 (COMMIT)
		return nil
	})

	return txErr
}

// 分配内存，初始化零值并返回指针
// var UserSvc = new(UserService)
var UserSvc = &UserService{}
