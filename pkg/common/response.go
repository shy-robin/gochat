package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/pkg/global/log"
)

type SuccessResponse struct {
	// 业务状态标识: "success" 或 "error"
	Status string `json:"status"`

	// Code 是自定义的业务错误码 (例如: 30001, 40001)
	Code int `json:"code"`

	// Message 是面向用户的错误信息
	Message string `json:"message"`

	// 成功时返回的业务数据
	Data any `json:"data,omitempty"`

	// HTTPStatus 是这个错误对应的 HTTP 状态码，用于 Handler 层响应 (JSON 忽略)
	HTTPStatus int `json:"-"`

	// InternalError 是底层的 Go error，用于日志记录，不返回给客户端 (JSON 忽略)
	InternalError error `json:"-"`
}

// Wrap 创建一个新的 SuccessResponse，通常用于设置响应数据
func WrapSuccessResponse(this *SuccessResponse, data any) *SuccessResponse {
	if this == nil {
		return nil
	}
	// 复制错误，避免修改预定义的全局错误常量
	newRes := *this
	newRes.Data = data
	return &newRes
}

var (
	ResOk = &SuccessResponse{
		Code:       0,
		Status:     "success",
		Message:    "ok",
		HTTPStatus: http.StatusOK,
	}
	ResCreated = &SuccessResponse{
		Code:       0,
		Status:     "success",
		Message:    "ok",
		HTTPStatus: http.StatusCreated,
	}
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

type BadRequestResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"参数错误"`
}

type UnauthorizedResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"鉴权失败"`
}

type SuccessResponseConfig struct {
	HttpCode int
	Data     any
}

// Success 响应成功，用于 GET/PUT/POST 等操作
func GenerateSuccessResponse(ctx *gin.Context, res *SuccessResponse) {
	ctx.JSON(res.HTTPStatus, res)
	// TODO:
	log.Logger.Error("TODO", log.Any("成功", res))
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

// Fail 响应失败，统一处理所有错误
func GenerateFailedResponse(ctx *gin.Context, err *ServiceError) {
	ctx.JSON(err.HTTPStatus, err)
	ctx.Abort() // 终止后续 Handler 执行
	// TODO:
	log.Logger.Error("TODO", log.Any("参数校验失败", err))
}

// PasswordSetter 定义了设置 Password 字段的方法
type PasswordSetter interface {
	SetPassword()
}

type EmptyRequest struct{}
