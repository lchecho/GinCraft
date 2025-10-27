package service

import (
	"context"
	dtoUser "github.com/liuchen/gin-craft/internal/dto/user"
	pkgCtx "github.com/liuchen/gin-craft/internal/pkg/context"
	"github.com/liuchen/gin-craft/internal/pkg/database"
	"go.uber.org/zap"
	"time"

	"github.com/liuchen/gin-craft/internal/constant"
	"github.com/liuchen/gin-craft/internal/dao"
	"github.com/liuchen/gin-craft/internal/model"
	"github.com/liuchen/gin-craft/pkg/errors"
)

// UserService 用户服务
type userService struct {
	userDAO *dao.UserDAO
}

var UserService = &userService{
	userDAO: dao.GetUserDAO(),
}

// Register 用户注册
func (s *userService) Register(req *dtoUser.RegisterRequest) error {
	// 检查用户名是否已存在
	usernameExists, err := s.userDAO.ExistsByUsername(req.Username)
	if err != nil {
		return err
	}
	if usernameExists {
		return errors.New(constant.UsernameAlreadyExist)
	}

	// 检查邮箱是否已存在
	emailExists, err := s.userDAO.ExistsByEmail(req.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New(constant.EmailAlreadyExist)
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: req.Password, // 实际应用中应该对密码进行加密
		Email:    req.Email,
	}

	return s.userDAO.Create(user)
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *dtoUser.LoginRequest) (string, error) {
	appCtx := pkgCtx.MustGetContext(ctx)

	user, err := s.userDAO.GetByUsername(req.Username)
	if err != nil {
		return "", err
	}
	// 实际应用中应该验证密码
	if user.Password != req.Password {
		return "", errors.New(constant.PasswordError)
	}
	// 生成 token
	token := "mock_token_" + req.Username + "_" + time.Now().Format("20060102150405")

	appCtx.LogInfo("用户登录成功", zap.String("username", req.Username))

	return token, nil
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(req *dtoUser.ListRequest) (*dtoUser.ListResponse, error) {
	users, err := s.userDAO.GetList(req)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	userList := make([]dtoUser.User, 0, len(users))
	for _, u := range users {
		userList = append(userList, dtoUser.User{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return &dtoUser.ListResponse{
		List:       userList,
		Pagination: req.Pagination,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(req *dtoUser.InfoRequest) (*model.User, error) {
	// 获取用户
	var user model.User
	err := dao.FirstByCondition(database.GetDatabase(), &user, map[string]interface{}{"id": req.ID})
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, errors.New(constant.UserNotExist)
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(ctx context.Context, req *dtoUser.UpdateRequest) error {
	appCtx := pkgCtx.MustGetContext(ctx)

	// 获取用户
	var user model.User
	err := dao.FirstByCondition(database.GetDatabase(), &user, map[string]interface{}{"id": req.ID})
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New(constant.UserNotExist)
	}

	updates := make(map[string]interface{})
	err = s.userDAO.Update(req.ID, updates)
	if err != nil {
		return err
	}
	appCtx.LogInfo("更新用户信息", zap.Uint("user_id", req.ID))

	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, req *dtoUser.InfoRequest) error {
	appCtx := pkgCtx.MustGetContext(ctx)

	// 获取用户
	var user model.User
	err := dao.FirstByCondition(database.GetDatabase(), &user, map[string]interface{}{"id": req.ID})
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New(constant.UserNotExist)
	}

	err = s.userDAO.Delete(req.ID)
	if err != nil {
		return err
	}
	appCtx.LogInfo("删除用户信息", zap.Uint("user_id", req.ID))

	return nil
}
