package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/api/v1"
)

func NewRouter() *gin.Engine {
	ginServer := gin.Default()

	{
		group1 := ginServer.Group("/api/v1")

		group1.POST("/user/register", v1.Register)
	}

	return ginServer
}
