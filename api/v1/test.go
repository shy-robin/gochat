package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func GetUserInfo(ctx *gin.Context) {
	userId := ctx.Query("id")
	userName := ctx.Query("name")

	ctx.JSON(http.StatusOK, gin.H{
		"userId":   userId,
		"userName": userName,
	})
}

func GetUserInfoRestful(ctx *gin.Context) {
	userId := ctx.Param("id")
	userName := ctx.Param("name")

	ctx.JSON(http.StatusOK, gin.H{
		"userId":   userId,
		"userName": userName,
	})
}

func TestJson(ctx *gin.Context) {
	data, _ := ctx.GetRawData()
	var res map[string]any

	_ = json.Unmarshal(data, &res)

	ctx.JSON(http.StatusOK, res)
}

func TestRedirect(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
}
