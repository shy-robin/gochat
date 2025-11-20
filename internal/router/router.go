package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/api/v1"
)

func NewRouter() *gin.Engine {
	ginServer := gin.Default()

	{
		group1 := ginServer.Group("/api/v1")
		group1.GET("/ping", v1.Ping)
		group1.GET("/userInfo", v1.GetUserInfo)
		group1.GET("/user/info/:id/:name", v1.GetUserInfoRestful)
		group1.POST("/json", v1.TestJson)
		group1.GET("/redirect", v1.TestRedirect)

		group1.POST("/user/register", v1.Register)
	}

	return ginServer
}
