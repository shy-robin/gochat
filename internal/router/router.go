package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/config"
	"github.com/shy-robin/gochat/docs"
	"github.com/shy-robin/gochat/internal/handler/v1"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           GoChat Swagger API
// @version         1.0
// @description     This is the Swagger API docs for the GoChat API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func NewRouter() *gin.Engine {
	ginServer := gin.Default()

	ginServer.Use(enableCORS())

	{
		group1 := ginServer.Group("/api/v1")

		group1.POST("/users", v1.Register)
	}

	// programatically set swagger info
	apiConfig := config.GetConfig().Api
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", apiConfig.Host, apiConfig.Port)
	docs.SwaggerInfo.BasePath = apiConfig.Prefix

	ginServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return ginServer
}

func enableCORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			ctx.Header(
				"Access-Control-Allow-Origin",
				"http://localhost:8083",
			) // 可将将 * 替换为指定的域名
			ctx.Header(
				"Access-Control-Allow-Methods",
				"POST,GET, OPTIONS, PUT, DELETE, UPDATE",
			)
			ctx.Header(
				"Access-Control-Allow-Headers",
				"Origin, X-Requested-With, Content-Type, Accept, Authorization",
			)
			ctx.Header(
				"Access-Control-Expose-Headers",
				"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type",
			)
			ctx.Header(
				"Access-Control-Allow-Credentials",
				"true",
			)
		}

		// 如果是预检请求 (OPTIONS)，直接返回 200/204 状态码
		// 正确的预检响应：应该只有 CORS 头部，响应体为空，状态码为 200 OK 或 204 No Content
		// 预检请求只是一个“问路”的请求，它不应该消耗任何业务资源
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "ok")
		}

		ctx.Next()
	}
}
