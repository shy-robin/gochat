package service

import (
	"fmt"

	"github.com/shy-robin/gochat/internal/db"
	"github.com/shy-robin/gochat/internal/handler/v1/dto"
	"github.com/shy-robin/gochat/internal/model"
	"github.com/shy-robin/gochat/internal/repository"
	"github.com/shy-robin/gochat/pkg/common"
	"gorm.io/gorm"
)

type UserService struct {
}

func (this *UserService) Register(user *model.User) error {
	existingUser, err := repository.UserRepo.FindByUsername(user.Username)

	if err != nil {
		return common.WrapServiceError(common.ErrDatabaseFailed, fmt.Errorf("repo find by username failed: %w", err))
	}

	if existingUser != nil {
		return common.ErrUsernameConflict
	}

	db := db.GetDB()

	// txErr -> Transaction Error (事物错误)
	// 确保所有 DB 操作要么全部成功，要么全部失败
	txErr := db.Transaction(func(tx *gorm.DB) error {
		if err := repository.UserRepo.CreateUser(user); err != nil {
			// 如果创建失败，返回错误，GORM 自动回滚 (ROLLBACK)
			return common.WrapServiceError(common.ErrDatabaseFailed, fmt.Errorf("repo create user failed: %w", err))
		}
		// 事务成功，GORM 自动提交 (COMMIT)
		return nil
	})

	return common.WrapServiceError(common.ErrDatabaseFailed, fmt.Errorf("repo transaction failed: %w", txErr))
}

func (this *UserService) Login(params *dto.LoginRequest) (string, int64, error) {
	existingUser, err := repository.UserRepo.FindByUsername(params.Username)
	defaultToken := ""
	defaultExpireTime := int64(0)

	if err != nil {
		return defaultToken, defaultExpireTime, common.WrapServiceError(common.ErrDatabaseFailed, fmt.Errorf("repo find by username failed: %w", err))
	}

	if existingUser == nil {
		return defaultToken, defaultExpireTime, common.ErrUserNotFound
	}

	isPasswordCorrect := existingUser.CheckPassword(params.Password)

	if !isPasswordCorrect {
		return defaultToken, defaultExpireTime, common.ErrWrongPassword
	}

	// 生成 Token
	token, expireTime, tokenErr := common.GenerateToken(existingUser.Uuid, existingUser.Username)

	if tokenErr != nil {
		return defaultToken, defaultExpireTime, common.WrapServiceError(common.ErrDatabaseFailed, fmt.Errorf("generate token failed: %w", tokenErr))
	}

	return token, expireTime, nil
}

func (this *UserService) GetUserInfo(uuid string) (*dto.GetUserInfoData, error) {
	user, err := repository.UserRepo.FindByUuid(uuid)

	if err != nil {
		return nil, common.WrapServiceError(common.ErrDatabaseFailed, fmt.Errorf("repo find by uuid failed: %w", err))
	}

	if user == nil {
		return nil, common.ErrUserNotFound
	}

	return &dto.GetUserInfoData{
		Username: user.Username,
		Uuid:     user.Uuid,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}, nil
}

func (this *UserService) ModifyUserInfo(
	uuid string,
	updates dto.ModifyUserInfoRequest,
) (*dto.ModifyUserInfoData, error) {
	userInfo, err := repository.UserRepo.UpdatesByUuid(uuid, updates)

	if err != nil {
		return nil, common.WrapServiceError(common.ErrDatabaseFailed, fmt.Errorf("repo updates by uuid failed: %w", err))
	}

	if userInfo == nil {
		return nil, common.ErrUserNotFound
	}

	return &dto.ModifyUserInfoData{
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.Avatar,
		Email:    userInfo.Email,
	}, nil
}

// 分配内存，初始化零值并返回指针
// var UserSvc = new(UserService)
var UserSvc = &UserService{}
