package service

import (
	"github.com/shy-robin/gochat/internal/db"
	"github.com/shy-robin/gochat/internal/model"
)

type userService struct {
}

func (this *userService) Register(user *model.User) {
	db.GetDB()
}

// 分配内存，初始化零值并返回指针
var UserService = new(userService)
