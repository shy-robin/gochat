package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
func Register(ctx *gin.Context) {
	var user dto.CreateUserRequest
	// 会经过 json 解析，binding 校验，如果不通过则会报错
	// ShouldBindJSON 可以自定义错误信息，而 BindJSON 会默认返回 400 状态码
	// 1. 调用 ShouldBindJSON
	// Gin 会自动尝试解析 JSON 到 req 结构体，并根据 `validate` 标签进行校验。
	err := ctx.ShouldBindJSON(&user)

	// 数据脱敏
	logUser := user
	logUser.Password = "******"
	log.Logger.Info("注册用户", log.Any("传参", logUser))

	// 参数校验失败
	if err != nil {
		// 2. 处理错误
		// err 可能来自 JSON 解析失败 (如格式错误)，也可能来自校验失败。

		// 尝试将错误转换为 validator.ValidationErrors (如果校验失败)
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// 校验失败

			// 3. 提取并返回用户友好的错误信息
			// 编写一个函数来处理这些错误，以便返回给客户端更清晰的信息。
			// formattedErrors := common.FormatValidationErrors(validationErrors)
			// jsonBytes, _ := json.Marshal(formattedErrors)
			// errMsg := string(jsonBytes)

			common.FailResponse(
				ctx,
				common.WithFailResponseHttpCode(http.StatusBadRequest),
				common.WithFailResponseMessage("参数校验失败"),
			)

			log.Logger.Error("注册用户", log.Any("参数校验失败", validationErrors))
			return
		}

		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage("参数校验失败"),
		)

		log.Logger.Error("注册用户", log.Any("参数校验失败", err))
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

		log.Logger.Error("注册用户", log.Any("数据库操作失败", err))
		return
	}

	// 注册成功
	common.SuccessResponse(
		ctx, common.WithSuccessResponseHttpCode(http.StatusCreated),
		common.WithSuccessResponseData(dto.CreateUserResponseData{
			Username: userEntity.Username,
			Uuid:     userEntity.Uuid,
			CreateAt: userEntity.BaseModel.CreatedAt,
		}),
	)

	// 数据脱敏
	userEntity.Password = "******"
	log.Logger.Info("注册用户", log.Any("注册成功", userEntity))
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
func Login(ctx *gin.Context) {
	var params dto.LoginRequest

	err := ctx.ShouldBindJSON(&params)

	logParams := params
	logParams.Password = "******"
	log.Logger.Info("登录", log.Any("传参", logParams))

	// 参数校验失败
	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage("参数校验失败"),
		)

		log.Logger.Error("登录", log.Any("参数校验失败", err))
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

		log.Logger.Error("登录", log.Any("数据库操作失败", loginErr))
		return
	}

	// 登录成功
	common.SuccessResponse(
		ctx, common.WithSuccessResponseHttpCode(http.StatusCreated),
		common.WithSuccessResponseData(dto.LoginResponseData{
			Token:    token,
			ExpireAt: expireTime,
		}),
	)

	log.Logger.Info("登录", log.Any("登录成功", expireTime))
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
func GetUsersMe(ctx *gin.Context) {
	userIdValue, ok := ctx.Get("userId")

	if !ok {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusUnauthorized),
			common.WithFailResponseMessage("鉴权失败"),
		)

		log.Logger.Error("获取当前用户信息", log.Any("鉴权失败", "用户 ID 不存在"))
		return
	}

	userId := userIdValue.(string)

	userInfo, err := service.UserSvc.GetUserInfo(userId)

	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage(err.Error()),
		)

		log.Logger.Error("获取当前用户信息", log.Any("数据库操作失败", err))
		return
	}

	common.SuccessResponse(
		ctx,
		common.WithSuccessResponseData(userInfo),
	)

	log.Logger.Info("获取当前用户信息", log.Any("获取成功", userInfo))
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
func GetUsers(ctx *gin.Context) {
	id := ctx.Param("id")

	userInfo, err := service.UserSvc.GetUserInfo(id)

	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage(err.Error()),
		)

		log.Logger.Error("获取当前用户信息", log.Any("数据库操作失败", err))
		return
	}

	if userInfo == nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage("用户不存在"),
		)

		log.Logger.Error("获取当前用户信息", log.Any("参数校验失败", "用户不存在"))
		return
	}

	common.SuccessResponse(
		ctx,
		common.WithSuccessResponseData(userInfo),
	)

	log.Logger.Info("获取当前用户信息", log.Any("获取成功", userInfo))
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
func ModifyUsersMe(ctx *gin.Context) {
	var params dto.ModifyUserInfoRequest

	if err := ctx.ShouldBindJSON(&params); err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusBadRequest),
			common.WithFailResponseMessage("参数校验失败"),
		)

		log.Logger.Error("修改当前用户信息", log.Any("参数校验失败", err))
		return
	}

	userIdValue, ok := ctx.Get("userId")

	if !ok {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusUnauthorized),
			common.WithFailResponseMessage("鉴权失败"),
		)

		log.Logger.Error("修改当前用户信息", log.Any("鉴权失败", "用户 ID 不存在"))
		return
	}

	userId := userIdValue.(string)

	userInfo, err := service.UserSvc.ModifyUserInfo(userId, params)

	if err != nil {
		common.FailResponse(
			ctx,
			common.WithFailResponseHttpCode(http.StatusInternalServerError),
			common.WithFailResponseMessage(err.Error()),
		)

		log.Logger.Error("修改当前用户信息", log.Any("数据库操作失败", err))
		return
	}

	common.SuccessResponse(
		ctx,
		common.WithSuccessResponseData(userInfo),
	)

	log.Logger.Info("修改当前用户信息", log.Any("修改成功", userInfo))
}
