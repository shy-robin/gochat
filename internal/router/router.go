package router

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	ginServer := gin.Default()

	ginServer.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	ginServer.GET("/userInfo", func(ctx *gin.Context) {
		userId := ctx.Query("id")
		userName := ctx.Query("name")

		ctx.JSON(http.StatusOK, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})

	ginServer.GET("/user/info/:id/:name", func(ctx *gin.Context) {
		userId := ctx.Param("id")
		userName := ctx.Param("name")

		ctx.JSON(http.StatusOK, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})

	ginServer.POST("/json", func(ctx *gin.Context) {
		data, _ := ctx.GetRawData()
		var res map[string]any

		_ = json.Unmarshal(data, &res)

		ctx.JSON(http.StatusOK, res)
	})

	ginServer.GET("/redirect", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	return ginServer
}
