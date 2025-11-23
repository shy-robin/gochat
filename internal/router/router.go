package router

import (
	"github.com/gin-gonic/gin"
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

// @host      localhost:8083
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func NewRouter() *gin.Engine {
	ginServer := gin.Default()

	{
		group1 := ginServer.Group("/api/v1")

		group1.POST("/users", v1.Register)
	}

	// programatically set swagger info
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v1"

	ginServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return ginServer
}
