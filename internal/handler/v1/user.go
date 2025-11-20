package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/internal/model"
	"github.com/shy-robin/gochat/internal/service"
)

func Register(ctx *gin.Context) {
	var user model.User
	// 会经过 json 解析，binding 校验，如果不通过则会报错
	// ShouldBindJSON 可以自定义错误信息，而 BindJSON 会默认返回 400 状态码
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = service.UserSvc.Register(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}
