package common

import "net/http"

// ServiceError 定义了统一的业务错误结构体
type ServiceError struct {
	// Code 是自定义的业务错误码 (例如: 30001, 40001)
	Code int `json:"code"`

	// Message 是面向用户的错误信息
	Message string `json:"message"`

	// HTTPStatus 是这个错误对应的 HTTP 状态码，用于 Handler 层响应 (JSON 忽略)
	HTTPStatus int `json:"-"`

	// InternalError 是底层的 Go error，用于日志记录，不返回给客户端 (JSON 忽略)
	InternalError error `json:"-"`
}

// Error 实现 Go 的 error 接口
func (this *ServiceError) Error() string {
	return this.Message
}

// Wrap 创建一个新的 ServiceError，通常用于封装底层的错误
func WrapServiceError(this *ServiceError, internalErr error) *ServiceError {
	if this == nil {
		return nil
	}
	// 复制错误，避免修改预定义的全局错误常量
	newErr := *this
	newErr.InternalError = internalErr
	return &newErr
}

// --- 预定义的常见 Service 错误 ---
var (
	// 400 Bad Request
	ErrInvalidInput = &ServiceError{
		Code:       20001,
		Message:    "请求参数不正确",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrWrongPassword = &ServiceError{
		Code:       20002,
		Message:    "密码错误",
		HTTPStatus: http.StatusBadRequest,
	}

	// 404 Not Found
	ErrUserNotFound = &ServiceError{
		Code:       30001,
		Message:    "用户不存在",
		HTTPStatus: http.StatusNotFound,
	}

	// 409 Conflict
	ErrUsernameConflict = &ServiceError{
		Code:       40001,
		Message:    "该用户名已被占用",
		HTTPStatus: http.StatusConflict,
	}

	// 500 Internal Server Error
	ErrDatabaseFailed = &ServiceError{
		Code:       10001,
		Message:    "系统繁忙，请稍后重试",
		HTTPStatus: http.StatusInternalServerError,
	}
)
