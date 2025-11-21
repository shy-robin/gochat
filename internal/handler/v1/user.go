package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/internal/model"
	"github.com/shy-robin/gochat/internal/service"
	"github.com/shy-robin/gochat/pkg/common"
)

func Register(ctx *gin.Context) {
	var user model.User
	// 会经过 json 解析，binding 校验，如果不通过则会报错
	// ShouldBindJSON 可以自定义错误信息，而 BindJSON 会默认返回 400 状态码
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage(err.Error()),
		)
		return
	}

	err = service.UserSvc.Register(&user)

	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage(err.Error()),
		)
		return
	}

	common.SuccessResponse(
		ctx, common.WithSuccessResponseHttpCode(http.StatusCreated),
		common.WithSuccessResponseData(map[string]any{
			"username": user.Username,
			"uuid":     user.Uuid,
			"createAt": user.BaseModel.CreatedAt,
		}),
	)
}
