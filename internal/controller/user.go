package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/constant"
	"github.com/liuchen/gin-craft/internal/dto/user"
	"github.com/liuchen/gin-craft/internal/middleware"
	"github.com/liuchen/gin-craft/internal/service"
	"github.com/liuchen/gin-craft/pkg/errors"
	_ "github.com/liuchen/gin-craft/pkg/response"
	"go.uber.org/zap"
)

// UserController 用户控制器结构体
type UserController struct{}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口，创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response "注册成功"
// @Router /api/v1/user/register [post]
func (uc *UserController) Register(c *gin.Context, req *user.RegisterRequest) (interface{}, error) {
	// 获取应用Context
	appCtx := context.MustGetContext(c)

	// 记录操作开始
	appCtx.LogInfo("开始用户注册")

	// 设置自定义字段
	appCtx.SetCustomField("operation", "user_register")
	appCtx.SetCustomField("username", req.Username)
	appCtx.SetCustomField("email", req.Email)

	// 调用 service 层处理业务逻辑
	err := service.UserService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		appCtx.LogError("用户注册失败",
			zap.String("error", err.Error()),
			zap.String("username", req.Username),
		)
		return nil, err
	}

	// 记录操作成功
	appCtx.LogInfo("用户注册成功",
		zap.String("username", req.Username),
		zap.Duration("duration", appCtx.GetDuration()),
	)

	return nil, nil
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口，验证用户名和密码并返回访问令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "登录信息"
// @Success 200 {object} user.LoginResponse "登录成功"
// @Router /api/v1/user/login [post]
func (uc *UserController) Login(c *gin.Context, req *user.LoginRequest) (interface{}, error) {
	// 获取应用Context
	appCtx := context.MustGetContext(c)

	// 记录操作开始
	appCtx.LogInfo("开始用户登录")

	// 设置自定义字段
	appCtx.SetCustomField("operation", "user_login")
	appCtx.SetCustomField("username", req.Username)

	// 调用 service 层处理业务逻辑
	token, err := service.UserService.Login(req.Username, req.Password)
	if err != nil {
		appCtx.LogError("用户登录失败",
			zap.String("error", err.Error()),
			zap.String("username", req.Username),
		)
		return nil, err
	}

	// 记录操作成功
	appCtx.LogInfo("用户登录成功",
		zap.String("username", req.Username),
		zap.Duration("duration", appCtx.GetDuration()),
	)

	return user.LoginResponse{
		Token: token,
	}, nil
}

// Info 获取用户信息
// @Summary 获取用户信息
// @Description 根据访问令牌获取当前用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param Authorization header string true "访问令牌" default(Bearer {token})
// @Success 200 {object} user.User "获取成功"
// @Router /api/v1/user/info [get]
func (uc *UserController) Info(c *gin.Context, req *user.InfoRequest) (interface{}, error) {
	// 获取应用Context
	appCtx := context.MustGetContext(c)

	// 记录操作开始
	appCtx.LogInfo("开始获取用户信息")

	// 从请求头获取 token
	token := c.GetHeader("Authorization")
	if token == "" {
		appCtx.LogWarn("未提供Authorization头")
		return nil, errors.New(constant.Unauthorized)
	}

	// 设置自定义字段
	appCtx.SetCustomField("token_length", len(token))
	appCtx.SetCustomField("operation", "get_user_info")

	// 调用 service 层处理业务逻辑
	userInfo, err := service.UserService.GetUserInfo(token)
	if err != nil {
		appCtx.LogError("获取用户信息失败",
			zap.String("error", err.Error()),
			zap.String("token_prefix", token[:10]+"..."),
		)
		return nil, err
	}

	// 设置用户信息到Context
	appCtx.SetUser(
		strconv.FormatUint(uint64(userInfo.ID), 10),
		userInfo.Username,
		"user", // 默认角色
	)

	// 记录操作成功
	appCtx.LogInfo("用户信息获取成功",
		zap.String("user_id", strconv.FormatUint(uint64(userInfo.ID), 10)),
		zap.String("username", userInfo.Username),
		zap.Duration("duration", appCtx.GetDuration()),
	)

	return user.User{
		ID:        userInfo.ID,
		Username:  userInfo.Username,
		Email:     userInfo.Email,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
	}, nil
}
