package dto

import (
	"time"
)

// CreateUserRequest 是创建用户接口的请求体 (DTO)
// NOTE: 在 Gin 中，使用 binding 添加校验规则而不是 validate
type CreateUserRequest struct {
	Username string `json:"username" form:"username" example:"robin" binding:"required,min=2,max=20,username"`
	Password string `json:"password" form:"password" example:"123456" binding:"required,min=8,password"`
	Nickname string `json:"nickname" example:"robin"`
	Avatar   string `json:"avatar" example:"https://avatars.githubusercontent.com/u/123456?v=4"`
	Email    string `json:"email" example:"robin@qq.com"`
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
