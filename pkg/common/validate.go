package common

import (
	"reflect"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupCustomValidator(router *gin.Engine) {
	// 通过 binding.Validator.Engine() 获取底层校验器实例
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 1. 注册自定义校验函数
		v.RegisterValidation("password", ValidatePassword)
		// TODO: log
		// log.Println("自定义校验标签 'complexity' 注册成功。")

		// 2. 注册自定义标签名函数
		// 告诉 validator 库，当生成校验错误时
		// 字段名（err.Field()）应该取自结构体标签中 json 键的值
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "" {
				name = fld.Name
			}
			return name
		})
		// TODO: log
		// log.Println("自定义标签名函数注册成功。")
	} else {
		// TODO: log
		//log.Println("警告：无法获取底层的 go-playground/validator 实例。")
	}
}

// 校验密码
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}
	return hasDigit
}

// 辅助函数：将校验错误转换为更友好的格式
func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	// 这是一个简化的示例，实际应用中应该更精细
	// 可以使用 Gin 社区提供的错误转换库，或自行实现中文映射
	result := make(map[string]string)
	for _, err := range errs {
		fieldName := err.Field() // 结构体中的字段名 (如 Username)
		tag := err.Tag()         // 校验规则 (如 required, min, email)

		// 根据字段名和规则返回中文错误提示
		// 示例：
		if fieldName == "Username" && tag == "required" {
			result["username"] = "用户昵称不能为空"
		} else if fieldName == "Email" && tag == "email" {
			result["email"] = "邮箱格式不正确"
		} else {
			// 默认错误信息
			result[fieldName] = fieldName + " 校验失败，规则: " + tag
		}
	}
	return result
}
