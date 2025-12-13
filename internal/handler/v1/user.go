package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/shy-robin/gochat/internal/handler/v1/dto"
	"github.com/shy-robin/gochat/internal/model"
	"github.com/shy-robin/gochat/internal/service"
	"github.com/shy-robin/gochat/pkg/common"
	"github.com/shy-robin/gochat/pkg/global/log"
)

// @Summary		注册用户
// @Description	传入参数，注册用户
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request	body		dto.CreateUserRequest		true	"请求参数"
// @Success		201		{object}	dto.CreateUserResponse		"注册成功"
// @Failure		400		{object}	common.BadRequestResponse	"参数错误"
// @Router			/users [post]
func Register(
	ctx *gin.Context,
	req dto.CreateUserRequest,
) (*common.SuccessResponse, *common.ServiceError) {
	userEntity := model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Email:    req.Email,
	}
	userInfo, err := service.UserSvc.Register(&userEntity)

	// 数据库操作失败
	if err != nil {
		return nil, err
	}

	// 注册成功
	return common.WrapSuccessResponse(
		common.ResCreated,
		userInfo,
	), nil
}

// @Summary		用户登录
// @Description	传入参数，用户登录
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request	body		dto.LoginRequest			true	"请求参数"
// @Success		201		{object}	dto.LoginResponse			"登录成功"
// @Failure		400		{object}	common.BadRequestResponse	"参数错误"
// @Router			/sessions [post]
func Login(
	ctx *gin.Context,
	req dto.LoginRequest,
) (*common.SuccessResponse, *common.ServiceError) {
	res, loginErr := service.UserSvc.Login(&req)

	// 数据库操作失败
	if loginErr != nil {
		return nil, loginErr
	}

	// 登录成功
	return common.WrapSuccessResponse(
		common.ResOk,
		res,
	), nil
}

// @Summary		获取当前用户信息
// @Description	传入参数，获取当前用户信息
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		201	{object}	dto.GetUserInfoResponse		"获取成功"
// @Failure		400	{object}	common.BadRequestResponse	"参数错误"
// @Failure		401	{object}	common.UnauthorizedResponse	"鉴权失败"
// @Router			/users/me [get]
func GetUsersMe(
	ctx *gin.Context,
	req common.EmptyRequest,
) (*common.SuccessResponse, *common.ServiceError) {
	userIdValue, ok := ctx.Get("userId")

	if !ok {
		return nil, common.ErrTokenUserIdNotFound
	}

	userId := userIdValue.(string)

	log.Logger.Info("获取当前用户信息", log.Any("传参", userId))

	userInfo, err := service.UserSvc.GetUserInfo(userId)

	if err != nil {
		return nil, err
	}

	return common.WrapSuccessResponse(
		common.ResOk,
		userInfo,
	), nil
}

// @Summary		获取用户信息
// @Description	传入参数，获取用户信息
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		201	{object}	dto.GetUserInfoResponse		"获取成功"
// @Failure		400	{object}	common.BadRequestResponse	"参数错误"
// @Failure		401	{object}	common.UnauthorizedResponse	"鉴权失败"
// @Router			/users/:id [get]
func GetUsers(
	ctx *gin.Context,
	req common.EmptyRequest,
) (*common.SuccessResponse, *common.ServiceError) {
	id := ctx.Param("id")

	log.Logger.Info("获取用户信息", log.Any("传参", id))

	userInfo, err := service.UserSvc.GetUserInfo(id)

	if err != nil {
		return nil, err
	}

	return common.WrapSuccessResponse(
		common.ResOk,
		userInfo,
	), nil
}

// @Summary		修改当前用户信息
// @Description	传入参数，修改当前信息
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request	body		dto.ModifyUserInfoRequest	true	"请求参数"
// @Success		201		{object}	dto.ModifyUserInfoResponse	"获取成功"
// @Failure		400		{object}	common.BadRequestResponse	"参数错误"
// @Failure		401		{object}	common.UnauthorizedResponse	"鉴权失败"
// @Router			/users/me [patch]
func ModifyUsersMe(
	ctx *gin.Context,
	req dto.ModifyUserInfoRequest,
) (*common.SuccessResponse, *common.ServiceError) {
	var params dto.ModifyUserInfoRequest

	userIdValue, ok := ctx.Get("userId")

	if !ok {
		return nil, common.ErrTokenUserIdNotFound
	}

	userId := userIdValue.(string)

	userInfo, err := service.UserSvc.ModifyUserInfo(userId, params)

	if err != nil {
		return nil, err
	}

	return common.WrapSuccessResponse(
		common.ResOk,
		userInfo,
	), nil
}
