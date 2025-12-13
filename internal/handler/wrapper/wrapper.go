package wrapper

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shy-robin/gochat/pkg/common"
	"github.com/shy-robin/gochat/pkg/global/log"
)

type RequestWithPassword interface {
	Password() string
}

// WrapGinHandler 是将自定义 Handler 转换为 gin.HandlerFunc 的包装器。
func WrapGinHandler[T any](
	// 自定义 Handler 签名：
	// 1. 接收 Gin Context
	// 2. 接收已解析且已校验的请求结构体 T
	// 3. 返回 ServiceError (nil 表示成功)
	handler func(*gin.Context, T) (*common.SuccessResponse, *common.ServiceError),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 初始化请求结构体
		var req T

		// 检查泛型类型 T 是否是 EmptyRequest 类型
		// 如果是，则跳过 c.ShouldBindJSON 和日志记录
		if reflect.TypeOf(req) != reflect.TypeOf(common.EmptyRequest{}) {
			// 2. JSON 解析和绑定
			// 会经过 json 解析，binding 校验，如果不通过则会报错
			// ShouldBindJSON 可以自定义错误信息，而 BindJSON 会默认返回 400 状态码
			// 1) 调用 ShouldBindJSON
			// Gin 会自动尝试解析 JSON 到 req 结构体，并根据 `validate` 标签进行校验。
			err := ctx.ShouldBindJSON(&req)
			//          ╭─────────────────────────────────────────────────────────╮
			//          │                    上报请求参数日志                     │
			//          ╰─────────────────────────────────────────────────────────╯
			// 数据脱敏
			reqLog := req
			// 注意：req 必须作为指针传递给 any() 进行断言，因为 SetPassword 是指针接收器方法。
			if setter, ok := any(&reqLog).(common.PasswordSetter); ok {
				// 调用后，指针对应的 reqLog 将会设置脱敏后的密码
				setter.SetPassword()
			}
			// TODO:
			log.Logger.Info("TODO", log.Any("传参", reqLog))

			if err != nil { // 3. 执行自定义结构体校验逻辑
				validationErr := translateValidationErrors(err)
				common.GenerateFailedResponse(ctx, validationErr)
				return
			}
		}

		// 4. 调用业务 Handler
		res, serviceErr := handler(ctx, req)
		if serviceErr != nil {
			// 业务逻辑执行失败，统一响应 ServiceError
			common.GenerateFailedResponse(ctx, serviceErr)
			return
		}

		// 业务 Handler 成功处理请求并已在内部调用 c.JSON/c.String 响应。
		common.GenerateSuccessResponse(ctx, res)
	}
}

// translateValidationErrors 将 validator.ValidationErrors 转换为一个 ServiceError
func translateValidationErrors(err error) *common.ServiceError {
	// 检查错误是否是 validator.ValidationErrors 类型
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// 通常只处理第一个错误，或者可以循环处理所有错误
		firstError := validationErrors[0]

		// 提取字段名和校验失败的规则
		fieldName := firstError.Field() // 结构体中的字段名 (如 Username)
		tag := firstError.Tag()         // 校验规则 (如 required, min, email)

		outerMap, ok := common.ValidateErrorMessages[fieldName]
		if !ok {
			return nil
		}

		value, ok := outerMap[tag]
		if ok {
			return value
		}

		return nil
	}

	// 如果不是校验错误，而是其他绑定错误（如 JSON 格式错误），使用通用错误
	return common.WrapServiceError(common.ErrInvalidInput, err)
}
