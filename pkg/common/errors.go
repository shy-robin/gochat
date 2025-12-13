package common

import "net/http"

// ServiceError 定义了统一的业务错误结构体
type ServiceError struct {
	// 业务状态标识: "success" 或 "error"
	Status string `json:"status"`

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
		Status:     "error",
		Message:    "请求参数不正确",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrWrongPassword = &ServiceError{
		Code:       20002,
		Status:     "error",
		Message:    "密码错误",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrTokenUserIdNotFound = &ServiceError{
		Code:       20003,
		Status:     "error",
		Message:    "鉴权失败",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrMissAuthorizationHeader = &ServiceError{
		Code:       20004,
		Status:     "error",
		Message:    "缺失请求头Authorization",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidAuthorizationHeader = &ServiceError{
		Code:       20005,
		Status:     "error",
		Message:    "Authorization格式错误",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidToken = &ServiceError{
		Code:       20006,
		Status:     "error",
		Message:    "Token无效",
		HTTPStatus: http.StatusUnauthorized,
	}

	// 404 Not Found
	ErrUserNotFound = &ServiceError{
		Code:       30001,
		Status:     "error",
		Message:    "用户不存在",
		HTTPStatus: http.StatusNotFound,
	}

	ErrUserNameEmpty = &ServiceError{
		Code:       30002,
		Status:     "error",
		Message:    "用户名不能为空",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrUserNameTooShort = &ServiceError{
		Code:       30003,
		Status:     "error",
		Message:    "用户名长度必须大于2个字符",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrUserNameTooLong = &ServiceError{
		Code:       30004,
		Status:     "error",
		Message:    "用户名长度必须小于20个字符",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrUserNameInvalid = &ServiceError{
		Code:       30005,
		Status:     "error",
		Message:    "用户名只能包含字母、数字、下划线，且以字母开头",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPsswordEmpty = &ServiceError{
		Code:       30006,
		Status:     "error",
		Message:    "密码不能为空",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPsswordTooShort = &ServiceError{
		Code:       30007,
		Status:     "error",
		Message:    "密码长度必须大于8个字符",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPsswordTooLong = &ServiceError{
		Code:       30008,
		Status:     "error",
		Message:    "密码长度必须小于50个字符",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrPsswordInvalid = &ServiceError{
		Code:       30009,
		Status:     "error",
		Message:    "密码必须包含大写字母、小写字母、数字、特殊字符中的至少三种类型",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrNicknameTooShort = &ServiceError{
		Code:       30010,
		Status:     "error",
		Message:    "昵称长度必须大于2个字符",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrNicknameTooLong = &ServiceError{
		Code:       30011,
		Status:     "error",
		Message:    "昵称长度必须小于20个字符",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrAvatarInvalid = &ServiceError{
		Code:       30012,
		Status:     "error",
		Message:    "头像地址格式不正确",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrEmailInvalid = &ServiceError{
		Code:       30013,
		Status:     "error",
		Message:    "邮箱格式不正确",
		HTTPStatus: http.StatusBadRequest,
	}

	// 409 Conflict
	ErrUsernameConflict = &ServiceError{
		Code:       40001,
		Status:     "error",
		Message:    "该用户名已被占用",
		HTTPStatus: http.StatusConflict,
	}

	// 500 Internal Server Error
	ErrDatabaseFailed = &ServiceError{
		Code:       10001,
		Status:     "error",
		Message:    "系统繁忙，请稍后重试",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// 用于存储校验错误消息
var ValidateErrorMessages = map[string]map[string]*ServiceError{
	"username": {
		"required": ErrUserNameEmpty,
		"min":      ErrUserNameTooShort,
		"max":      ErrUserNameTooLong,
		"username": ErrUserNameInvalid,
	},
	"password": {
		"required": ErrPsswordEmpty,
		"min":      ErrPsswordTooShort,
		"max":      ErrPsswordTooLong,
		"password": ErrPsswordInvalid,
	},
	"nickname": {
		"min": ErrNicknameTooShort,
		"max": ErrNicknameTooLong,
	},
	"avatar": {
		"url": ErrAvatarInvalid,
	},
	"email": {
		"email": ErrEmailInvalid,
	},
}
