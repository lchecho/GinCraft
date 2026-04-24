package service

import (
	"context"
	"errors"
	"time"

	"github.com/liuchen/gin-craft/internal/constant"
	"github.com/liuchen/gin-craft/internal/dao"
	dtoUser "github.com/liuchen/gin-craft/internal/dto/user"
	"github.com/liuchen/gin-craft/internal/model"
	pkgCtx "github.com/liuchen/gin-craft/internal/pkg/context"
	apperr "github.com/liuchen/gin-craft/internal/pkg/errors"
	"github.com/liuchen/gin-craft/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// userService 用户服务
type userService struct {
	userDAO *dao.UserDAO
}

// NewUserService 构造函数
func NewUserService() *userService {
	return &userService{userDAO: dao.GetUserDAO()}
}

// UserService 全局默认实例（兼容旧调用方）
var UserService = NewUserService()

// Register 用户注册
func (s *userService) Register(req *dtoUser.RegisterRequest) error {
	if exists, err := s.userDAO.ExistsByUsername(req.Username); err != nil {
		return err
	} else if exists {
		return apperr.New(constant.UsernameAlreadyExist)
	}
	if exists, err := s.userDAO.ExistsByEmail(req.Email); err != nil {
		return err
	} else if exists {
		return apperr.New(constant.EmailAlreadyExist)
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return apperr.New(constant.SystemError, err.Error())
	}
	return s.userDAO.Create(&model.User{
		Username: req.Username,
		Password: hashed,
		Email:    req.Email,
	})
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *dtoUser.LoginRequest) (*dtoUser.LoginResponse, error) {
	appCtx := pkgCtx.MustGetContext(ctx)

	user, err := s.userDAO.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(constant.UserNotExist)
		}
		return nil, err
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, apperr.New(constant.PasswordError)
	}

	token := "mock_token_" + req.Username + "_" + time.Now().Format("20060102150405")
	appCtx.LogInfo("用户登录成功", zap.String("username", req.Username))
	return &dtoUser.LoginResponse{Token: token}, nil
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(req *dtoUser.ListRequest) (*dtoUser.ListResponse, error) {
	users, err := s.userDAO.GetList(req)
	if err != nil {
		return nil, err
	}
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
	return &dtoUser.ListResponse{List: userList, Pagination: req.Pagination}, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(req *dtoUser.InfoRequest) (*dtoUser.User, error) {
	u, err := s.userDAO.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(constant.UserNotExist)
		}
		return nil, err
	}
	return &dtoUser.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

// UpdateUser 更新用户信息（白名单字段）
func (s *userService) UpdateUser(ctx context.Context, req *dtoUser.UpdateRequest) error {
	appCtx := pkgCtx.MustGetContext(ctx)

	if _, err := s.userDAO.GetByID(req.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.New(constant.UserNotExist)
		}
		return err
	}

	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if len(updates) == 0 {
		return nil
	}
	if err := s.userDAO.Update(req.ID, updates); err != nil {
		return err
	}
	appCtx.LogInfo("更新用户信息", zap.Uint("user_id", req.ID))
	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, req *dtoUser.InfoRequest) error {
	appCtx := pkgCtx.MustGetContext(ctx)
	if _, err := s.userDAO.GetByID(req.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.New(constant.UserNotExist)
		}
		return err
	}
	if err := s.userDAO.Delete(req.ID); err != nil {
		return err
	}
	appCtx.LogInfo("删除用户信息", zap.Uint("user_id", req.ID))
	return nil
}
