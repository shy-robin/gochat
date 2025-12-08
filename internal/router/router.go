package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/config"
	"github.com/shy-robin/gochat/docs"
	v1 "github.com/shy-robin/gochat/internal/handler/v1"
	"github.com/shy-robin/gochat/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			GoChat Swagger API
//	@version		1.0
//	@description	This is the Swagger API docs for the GoChat API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func NewRouter() *gin.Engine {
	ginServer := gin.Default()

	// 添加 Recovery 中间件，捕获并恢复程序，防止整个 Web 服务崩溃
	// 您的 Go 进程不会终止，而是会向客户端返回 HTTP 500 错误
	// 同时，您会在控制台中看到 Recovery 中间件打印的堆栈跟踪日志
	// ginServer.Use(gin.Recovery()) // NOTE: gin.Default() 默认包含了 gin.Recovery() 中间件，因此不需要重复添加

	// ginServer.Use(enableCORS())

	ginServer.Use(cors.New(cors.Config{
		// 必须配置项
		AllowOrigins: []string{
			"http://localhost:8083",
			"http://localhost:3000",
			// "https://your-frontend-domain.com", // 允许的前端域名
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的方法
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}, // 允许的头部

		// 核心配置项：设置预检请求的缓存时间为 12 小时 (43200 秒)
		MaxAge: 12 * time.Hour,

		// 可选配置项
		AllowCredentials: true, // 如果需要发送 Cookie 或 HTTP 认证信息
	}))

	{
		group1 := ginServer.Group("/api/v1")

		group1.POST("/users", v1.Register)
		group1.POST("/sessions", v1.Login)

		{
			userGroup := group1.Group("/users")
			// 公开访问：获取单个用户详情
			userGroup.GET("/:id", v1.GetUsers)
			// 私有访问：获取当前用户信息（需要认证）
			userGroup.GET("/me", middleware.JWTAuthMiddleware(), v1.GetUsersMe)
			userGroup.PATCH("/me", middleware.JWTAuthMiddleware(), v1.ModifyUsersMe)
		}
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
		// 只有非简单请求，才会发送预检请求，并且预检请求会有缓存机制（Access-Control-Max-Age）
		// 核心原则： 当您更改 CORS 策略时，所有用户的浏览器必须等待 Max-Age 时间过后才能获取最新的策略。长时间缓存可能导致用户在缓存期间无法访问新配置的 API。
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "ok")
		}

		ctx.Next()
	}
}
