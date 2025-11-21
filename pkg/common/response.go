package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是通用的 API 响应结构
type Response struct {
	// 业务状态标识: "success" 或 "error"
	Status string `json:"status"`

	// 成功时返回的业务数据
	Data any `json:"data,omitempty"`

	// 元数据: 用于分页、总数等信息
	Meta any `json:"meta,omitempty"`

	// ------------------------------------------
	// 以下字段仅在 Status 为 "error" 时填充

	// 内部业务错误码 (最关键的字段，用于前端精确判断)
	ErrorCode int `json:"error_code,omitempty"`

	// 供前端或用户显示的提示信息
	Message string `json:"message,omitempty"`

	// 详细错误信息 (例如表单验证失败的字段细节)
	Details any `json:"details,omitempty"`
}

type SuccessResponseConfig struct {
	HttpCode int
	Data     any
}

type SuccessResponseOption func(*SuccessResponseConfig)

func WithSuccessResponseHttpCode(httpCode int) SuccessResponseOption {
	return func(config *SuccessResponseConfig) {
		config.HttpCode = httpCode
	}
}

func WithSuccessResponseData(data any) SuccessResponseOption {
	return func(config *SuccessResponseConfig) {
		config.Data = data
	}
}

// Success 响应成功，用于 GET/PUT/POST 等操作
func SuccessResponse(ctx *gin.Context, opts ...SuccessResponseOption) {
	// 设置默认值
	config := &SuccessResponseConfig{
		HttpCode: http.StatusOK,
	}

	// 依次应用传入的选项函数，覆盖默认值
	for _, option := range opts {
		option(config)
	}

	ctx.JSON(config.HttpCode, Response{
		Status: "success",
		Data:   config.Data,
	})
}

type SuccessListResponseConfig struct {
	HttpCode int
	Data     any
	Meta     any
}

type SuccessListResponseOption func(*SuccessListResponseConfig)

func WithSuccessListResponseHttpCode(httpCode int) SuccessListResponseOption {
	return func(config *SuccessListResponseConfig) {
		config.HttpCode = httpCode
	}
}

func WithSuccessListResponseData(data any) SuccessListResponseOption {
	return func(config *SuccessListResponseConfig) {
		config.Data = data
	}
}

func WithSuccessListResponseMeta(meta any) SuccessListResponseOption {
	return func(config *SuccessListResponseConfig) {
		config.Meta = meta
	}
}

// SuccessList 响应列表成功，包含分页元数据
func SuccessListResponse(ctx *gin.Context, opts ...SuccessListResponseOption) {
	// 设置默认值
	config := &SuccessListResponseConfig{
		HttpCode: http.StatusOK,
	}

	// 依次应用传入的选项函数，覆盖默认值
	for _, option := range opts {
		option(config)
	}

	ctx.JSON(config.HttpCode, Response{
		Status: "success",
		Data:   config.Data,
		Meta:   config.Meta,
	})
}

// NoContent 响应成功，无内容 (204)
func NoContentResponse(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNoContent)
}

type FailResponseConfig struct {
	HttpCode int
	ErrCode  int
	Message  string
	Details  any
}

type FailResponseOption func(*FailResponseConfig)

func WithFailResponseHttpCode(httpCode int) FailResponseOption {
	return func(config *FailResponseConfig) {
		config.HttpCode = httpCode
	}
}

func WithFailResponseErrCode(errCode int) FailResponseOption {
	return func(config *FailResponseConfig) {
		config.ErrCode = errCode
	}
}
func WithFailResponseMessage(message string) FailResponseOption {
	return func(config *FailResponseConfig) {
		config.Message = message
	}
}

func WithFailResponseDetails(details any) FailResponseOption {
	return func(config *FailResponseConfig) {
		config.Details = details
	}
}

// Fail 响应失败，统一处理所有错误
func FailResponse(ctx *gin.Context, opts ...FailResponseOption) {
	// 设置默认值
	config := &FailResponseConfig{
		HttpCode: http.StatusInternalServerError,
	}

	// 依次应用传入的选项函数，覆盖默认值
	for _, option := range opts {
		option(config)
	}

	ctx.JSON(config.HttpCode, Response{
		Status:    "error",
		ErrorCode: config.ErrCode,
		Message:   config.Message,
		Details:   config.Details,
	})
	ctx.Abort() // 终止后续 Handler 执行
}

// ErrorCode 定义，例如:
const (
	CodeValidationFailed = 1001 // 客户端参数错误
	CodeNotFound         = 1002 // 资源不存在
	CodeInternalError    = 9999 // 服务器内部错误
)
