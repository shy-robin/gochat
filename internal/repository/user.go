package repository

import (
	"errors"

	"github.com/shy-robin/gochat/internal/db"
	"github.com/shy-robin/gochat/internal/handler/v1/dto"
	"github.com/shy-robin/gochat/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// NOTE: Repository 层负责隔离数据库操作
// 错误处理：
// - 职责：基础设施错误处理
// - 错误载体：Go error
// - 核心作用：仅传递数据库或外部服务错误

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

func (this *UserRepository) FindByUuid(uuid string) (*model.User, error) {
	db := db.GetDB()
	user := &model.User{}

	result := db.Where("uuid = ?", uuid).First(user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, result.Error
}

func (this *UserRepository) UpdatesByUuid(
	uuid string,
	updates dto.ModifyUserInfoRequest,
) (*model.User, error) {
	db := db.GetDB()

	if updates.Password != "" {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(updates.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, hashErr
		}
		// 替换明文密码
		updates.Password = string(hashedPassword)
	}
	// 需要加 Model 才能更新关联关系，否则无法找到对应的表
	result := db.Model(&model.User{}).Where("uuid = ?", uuid).Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	user := &model.User{}
	res := db.Where("uuid = ?", uuid).First(user)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, res.Error
}
