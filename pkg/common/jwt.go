package common

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shy-robin/gochat/config"
)

// Claims 定义了 JWT 的载荷信息
type Claims struct {
	jwt.RegisteredClaims
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

// GenerateToken 生成一个新的 JWT
func GenerateToken(userId string, username string) (string, int64, error) {
	jwtConfig := config.GetConfig().Jwt

	expireTime := time.Now().Add(time.Duration(jwtConfig.ExpireTime) * time.Hour).Unix()

	claims := &Claims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// exp: 设置过期时间 (必须使用 NewNumericDate)
			ExpiresAt: jwt.NewNumericDate(time.Unix(expireTime, 0)),
			// iat: 设置签发时间 (推荐设置)
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// iss: 设置签发者
			Issuer: "gochat-api-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))

	return tokenString, expireTime, err
}

// ValidateToken 验证 JWT 的有效性 (未在登录接口中使用，但用于后续接口)
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	jwtConfig := config.GetConfig().Jwt

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			// 需要将密码转换为 []byte 类型
			return []byte(jwtConfig.Secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
