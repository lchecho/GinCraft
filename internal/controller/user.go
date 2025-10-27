package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/dto/user"
	"github.com/liuchen/gin-craft/internal/service"
	_ "github.com/liuchen/gin-craft/pkg/response"
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
	return nil, service.UserService.Register(req)
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
	return service.UserService.Login(c, req)
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 根据请求条件获取用户列表信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.ListRequest true "用户列表信息"
// @Success 200 {object} user.ListResponse "获取成功"
// @Router /api/v1/user/info [post]
func (uc *UserController) List(c *gin.Context, req *user.ListRequest) (interface{}, error) {
	return service.UserService.GetUserList(req)
}

// Info 获取用户信息
// @Summary 获取用户信息
// @Description 根据ID获取用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.InfoRequest true "请求信息"
// @Success 200 {object} user.User "获取成功"
// @Router /api/v1/user/info [get]
func (uc *UserController) Info(c *gin.Context, req *user.InfoRequest) (interface{}, error) {
	return service.UserService.GetUserInfo(req)
}

// Update 更新用户
// @Summary 更新用户
// @Description 根据用户ID更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.UpdateRequest true "更新信息"
// @Success 200 {object} response.Response "更新成功"
// @Router /api/v1/user/update [post]
func (uc *UserController) Update(c *gin.Context, req *user.UpdateRequest) (interface{}, error) {
	return nil, service.UserService.UpdateUser(c, req)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 根据用户ID删除用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body user.InfoRequest true "删除信息"
// @Success 200 {object} response.Response "删除成功"
// @Router /api/v1/user/delete [post]
func (uc *UserController) Delete(c *gin.Context, req *user.InfoRequest) (interface{}, error) {
	return nil, service.UserService.DeleteUser(c, req)
}
