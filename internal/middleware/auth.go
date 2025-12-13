package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/pkg/common"
)

// JWTAuthMiddleware 是一个 Gin 中间件，用于验证 JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 从请求头中提取 Token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			common.GenerateFailedResponse(ctx, common.ErrMissAuthorizationHeader)
			return
		}

		// 2. 检查格式是否为 "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.GenerateFailedResponse(ctx, common.ErrInvalidAuthorizationHeader)
			return
		}
		tokenString := parts[1]

		// 3. 解析和验证 Token
		claims, err := common.ValidateToken(tokenString)

		// 4. 处理验证错误
		if err != nil {
			common.GenerateFailedResponse(ctx, common.ErrInvalidToken)
			return
		}

		// 5. 验证成功，将用户信息存入 Context
		ctx.Set("userId", claims.UserId)
		ctx.Set("username", claims.Username)
		ctx.Next() // 放行，请求继续执行后续的 Handler
	}
}
