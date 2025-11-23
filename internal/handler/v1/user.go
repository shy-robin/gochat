package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/internal/handler/v1/dto"
	"github.com/shy-robin/gochat/internal/model"
	"github.com/shy-robin/gochat/internal/service"
	"github.com/shy-robin/gochat/pkg/common"
)

// @Summary		注册用户
// @Description	传入参数，注册用户
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request	body		dto.CreateUserRequest	true	"请求参数"
// @Success		201		{object}	dto.CreateUserResponse
// @Failure		400		{object}	common.BadRequestResponse
// @Router			/users [post]
func Register(ctx *gin.Context) {
	var user dto.CreateUserRequest
	// 会经过 json 解析，binding 校验，如果不通过则会报错
	// ShouldBindJSON 可以自定义错误信息，而 BindJSON 会默认返回 400 状态码
	err := ctx.ShouldBindJSON(&user)

	// 参数校验失败
	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage(err.Error()),
		)
		return
	}

	userEntity := model.User{
		Username: user.Username,
		Password: user.Password,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Email:    user.Email,
	}
	err = service.UserSvc.Register(&userEntity)

	// 数据库操作失败
	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage(err.Error()),
		)
		return
	}

	// 注册成功
	common.SuccessResponse(
		ctx, common.WithSuccessResponseHttpCode(http.StatusCreated),
		common.WithSuccessResponseData(dto.CreateUserResponse{
			Username: userEntity.Username,
			Uuid:     userEntity.Uuid,
			CreateAt: userEntity.BaseModel.CreatedAt,
		}),
	)
}

func Login(ctx *gin.Context) {
	var params dto.LoginRequest

	err := ctx.ShouldBindJSON(&params)

	// 参数校验失败
	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage(err.Error()),
		)
		return
	}

	token, expireTime, loginErr := service.UserSvc.Login(&params)

	// 数据库操作失败
	if loginErr != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage(loginErr.Error()),
		)
		return
	}

	// 登录成功
	common.SuccessResponse(
		ctx, common.WithSuccessResponseHttpCode(http.StatusCreated),
		common.WithSuccessResponseData(dto.LoginResponse{
			Token:    token,
			ExpireAt: expireTime,
		}),
	)
}
