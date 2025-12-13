package dto

import (
	"time"
)

// CreateUserRequest 是创建用户接口的请求体 (DTO)
// NOTE: 在 Gin 中，使用 binding 添加校验规则而不是 validate
// 内置校验方法参考：https://github.com/go-playground/validator
type CreateUserRequest struct {
	Username string `json:"username" form:"username" example:"robin" binding:"required,min=2,max=20,username"`
	Password string `json:"password" form:"password" example:"123456" binding:"required,min=8,max=50,password"`
	Nickname string `json:"nickname" example:"robin" binding:"omitempty,min=2,max=20"`
	Avatar   string `json:"avatar" example:"https://avatars.githubusercontent.com/u/123456?v=4" binding:"omitempty,url"`
	Email    string `json:"email" example:"robin@qq.com" binding:"omitempty,email"`
}

func (this *CreateUserRequest) SetPassword() {
	// 数据脱敏
	this.Password = "******"
}

// CreateUserResponse  是返回给客户端的用户响应 (DTO)
type CreateUserResponse struct {
	Status string `json:"status" example:"success"`
	Data   CreateUserResponseData
}

type CreateUserResponseData struct {
	// 注意: 响应中不包含 PasswordHash
	Username string    `json:"username" example:"robin"`
	Uuid     string    `json:"uuid" example:"db376853-8f93-41f9-9a44-3c5ad8eedbbb"`
	CreateAt time.Time `json:"createAt" example:"2025-11-23T15:53:56.811"` // 可能需要格式化为字符串
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"robin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

func (this *LoginRequest) SetPassword() {
	// 数据脱敏
	this.Password = "******"
}

type LoginResponse struct {
	Status string `json:"status" example:"success"`
	Data   LoginResponseData
}

type LoginResponseData struct {
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2NoYXQtYXBpLXNlcnZpY2UiLCJleHAiOjE3NjM5OTc1MTksImlhdCI6MTc2MzkxMTExOSwidXNlcklkIjoiY2YzZDE1ZmMtM2ZlYS00NDkzLWFjMDMtYzBiMjkzYTBjNjc4IiwidXNlcm5hbWUiOiJyb2JpbjYifQ.Rddl8mxWIPBWYCIO5TQYTG8uvyyPbP3FF9ozGdfytwg"`
	ExpireAt int64  `json:"expireAt" example:"1763910483465"`
}

type GetUserInfoResponse struct {
	Status string `json:"status" example:"success"`
	Data   GetUserInfoData
}

type GetUserInfoData struct {
	Username string `json:"username" example:"robin"`
	Uuid     string `json:"uuid" example:"db376853-8f93-41f9-9a44-3c5ad8eedbbb"`
	Nickname string `json:"nickname" example:"robin"`
	Avatar   string `json:"avatar" example:"https://avatars.githubusercontent.com/u/123456?v=4"`
	Email    string `json:"email" example:"robin@test.com"`
}

type ModifyUserInfoRequest struct {
	// NOTE: omitempty 标签告诉 JSON 编码器在 nil 时忽略该字段
	// 如果参数需要可选，否则参数一定会校验。
	Password string `json:"password" example:"123456" binding:"omitempty,min=8,max=50,password"`
	Nickname string `json:"nickname" example:"robin" binding:"omitempty,min=2,max=20"`
	Avatar   string `json:"avatar" example:"https://avatars.githubusercontent.com/u/123456?v=4" binding:"omitempty,url"`
	Email    string `json:"email" example:"robin@test.com" binding:"omitempty,email"`
}

func (this *ModifyUserInfoRequest) SetPassword() {
	// 数据脱敏
	this.Password = "******"
}

type ModifyUserInfoResponse struct {
	Status string `json:"status" example:"success"`
	Data   ModifyUserInfoData
}

type ModifyUserInfoData struct {
	Nickname string `json:"nickname" example:"robin"`
	Avatar   string `json:"avatar" example:"https://avatars.githubusercontent.com/u/123456?v=4"`
	Email    string `json:"email" example:"robin@test.com"`
}
