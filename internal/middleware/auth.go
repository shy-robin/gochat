package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/pkg/common"
	"github.com/shy-robin/gochat/pkg/global/log"
)

// JWTAuthMiddleware 是一个 Gin 中间件，用于验证 JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 从请求头中提取 Token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			common.FailResponse(
				ctx,
				common.WithFailResponseHttpCode(http.StatusUnauthorized),
				common.WithFailResponseMessage("缺失请求头 Authorization"),
			)
			log.Logger.Error("鉴权失败", log.Any("缺失请求头 Authorization", authHeader))
			return
		}

		// 2. 检查格式是否为 "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.FailResponse(
				ctx,
				common.WithFailResponseHttpCode(http.StatusUnauthorized),
				common.WithFailResponseMessage("Authorization 格式错误"),
			)
			log.Logger.Error("鉴权失败", log.Any("Authorization 格式错误", authHeader))
			return
		}
		tokenString := parts[1]

		// 3. 解析和验证 Token
		claims, err := common.ValidateToken(tokenString)

		// 4. 处理验证错误
		if err != nil {
			common.FailResponse(
				ctx,
				common.WithFailResponseHttpCode(http.StatusUnauthorized),
				common.WithFailResponseMessage("无效的 Token"),
			)
			log.Logger.Error("鉴权失败", log.Any("无效的 Token", err))
			return
		}

		// 5. 验证成功，将用户信息存入 Context
		ctx.Set("userId", claims.UserId)
		ctx.Set("username", claims.Username)
		ctx.Next() // 放行，请求继续执行后续的 Handler
	}
}
