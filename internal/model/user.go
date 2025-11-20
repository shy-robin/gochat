package model

type User struct {
	BaseModel
	Uuid     string `json:"uuid" gorm:"type:varchar(150);not null;uniqueIndex:idx_uuid;comment:'uuid'"`
	Username string `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'用户名'"`
	Password string `json:"password" form:"password" binding:"required" gorm:"type:varchar(150);not null; comment:'密码'"`
	Nickname string `json:"nickname" gorm:"comment:'昵称'"`
	Avatar   string `json:"avatar" gorm:"type:varchar(150);comment:'头像'"`
	Email    string `json:"email" gorm:"type:varchar(80);column:email;comment:'邮箱'"`
}
