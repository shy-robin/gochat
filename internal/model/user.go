package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Uuid     string `json:"uuid" gorm:"type:varchar(150);not null;uniqueIndex:idx_uuid;comment:'uuid'"`
	Username string `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'用户名'"`
	Password string `json:"password" form:"password" binding:"required" gorm:"type:varchar(150);not null; comment:'密码'"`
	Nickname string `json:"nickname" gorm:"comment:'昵称'"`
	Avatar   string `json:"avatar" gorm:"type:varchar(150);comment:'头像'"`
	Email    string `json:"email" gorm:"type:varchar(80);column:email;comment:'邮箱'"`
}

// BeforeCreate 是 GORM 的 Hook 函数。
// 在执行 DB.Create() 时，数据被写入数据库之前，GORM 会自动调用此方法
// tx -> Transaction (它代表了一个 GORM 数据库事务连接)
func (this *User) BeforeCreate(tx *gorm.DB) (err error) {
	// 创建用户时，如果没有指定 UUID，则自动生成一个 UUID
	if this.Uuid == "" {
		this.Uuid = uuid.NewString()
	}

	// 返回 nil 表示操作成功，GORM 继续执行插入操作
	return nil
}
