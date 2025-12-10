package common

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shy-robin/gochat/pkg/global/log"
)

func SetupCustomValidator(router *gin.Engine) {
	// 通过 binding.Validator.Engine() 获取底层校验器实例
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 1. 注册自定义校验函数
		RegisterValidation(v, "username", ValidateUsername)
		RegisterValidation(v, "password", ValidatePassword)

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

		log.Logger.Info("自定义标签名函数注册成功")
	} else {
		log.Logger.Warn("无法获取底层的 go-playground/validator 实例")
	}
}

func RegisterValidation(v *validator.Validate, name string, fn validator.Func) {
	if err := v.RegisterValidation(name, fn); err != nil {
		log.Logger.Error(fmt.Sprintf("注册自定义校验函数 '%s' 失败", name), log.Any("error", err))
	} else {
		log.Logger.Info(fmt.Sprintf("自定义校验标签 '%s' 注册成功", name))
	}
}

// 校验用户名
func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9_-]{2,14})[a-zA-Z0-9]$")

	// 1. 正则表达式校验字符集和格式
	if !usernameRegex.MatchString(username) {
		return false
	}

	// 2. 保留字校验（简单示例）
	reservedNames := map[string]bool{
		"admin": true,
		"root":  true,
		"test":  true,
	}
	if reservedNames[username] {
		return false
	}

	return true
}

// 校验密码
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// 检查是否包含至少三种字符类型：
	// 1. 小写字母
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// 2. 大写字母
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// 3. 数字
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// 4. 特殊字符 (非字母或数字)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	// 计算包含的类型数量
	count := 0
	if hasLower {
		count++
	}
	if hasUpper {
		count++
	}
	if hasDigit {
		count++
	}
	if hasSpecial {
		count++
	}

	// 假设要求至少包含 3 种类型
	return count >= 3
}

var ValidateErrorMessages = map[string]map[string]string{
	"username": {
		"required": "用户昵称不能为空",
		"min":      "用户昵称长度必须大于2个字符",
		"max":      "用户昵称长度必须小于20个字符",
		"username": "用户昵称只能包含字母、数字、下划线，且以字母开头",
	},
	"password": {
		"required": "密码不能为空",
		"min":      "密码长度必须大于8个字符",
		"max":      "密码长度必须小于50个字符",
		"password": "密码必须包含大写字母、小写字母、数字、特殊字符中的至少三种类型",
	},
	"nickname": {
		"min": "昵称长度必须大于2个字符",
		"max": "昵称长度必须小于20个字符",
	},
	"avatar": {
		"url": "头像地址格式不正确",
	},
	"email": {
		"email": "邮箱格式不正确",
	},
}

// 辅助函数：将校验错误转换为更友好的格式
func FormatValidationErrors(errs validator.ValidationErrors) string {
	// 这是一个简化的示例，实际应用中应该更精细
	// 可以使用 Gin 社区提供的错误转换库，或自行实现中文映射
	// result := make(map[string]string)

	for _, err := range errs {
		fieldName := err.Field() // 结构体中的字段名 (如 Username)
		tag := err.Tag()         // 校验规则 (如 required, min, email)

		defaultMessage := fieldName + " 校验失败"
		// 根据字段名和规则返回中文错误提示
		// 示例：
		outerMap, ok := ValidateErrorMessages[fieldName]
		if !ok {
			return defaultMessage
		}

		value, ok := outerMap[tag]
		if !ok {
			return defaultMessage
		} else {
			return value
		}
	}

	return ""
}
