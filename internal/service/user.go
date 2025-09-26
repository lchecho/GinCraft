package service

import (
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
func (s *userService) Register(username, password, email string) error {
	// 检查用户名是否已存在
	usernameExists, err := s.userDAO.ExistsByUsername(username)
	if err != nil {
		return err
	}
	if usernameExists {
		return errors.New(constant.UsernameAlreadyExist)
	}

	// 检查邮箱是否已存在
	emailExists, err := s.userDAO.ExistsByEmail(email)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New(constant.EmailAlreadyExist)
	}

	// 创建用户
	user := &model.User{
		Username: username,
		Password: password, // 实际应用中应该对密码进行加密
		Email:    email,
	}

	return s.userDAO.Create(user)
}

// Login 用户登录
func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userDAO.GetByUsername(username)
	if err != nil {
		return "", err
	}

	// 实际应用中应该验证密码
	if user.Password != password {
		return "", errors.New(constant.PasswordError)
	}

	// 生成 token
	token := "mock_token_" + username + "_" + time.Now().Format("20060102150405")
	return token, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(token string) (*model.User, error) {
	// 实际应用中应该解析 token 获取用户 ID
	// 这里暂时返回 ID 为 1 的用户
	user, err := s.userDAO.GetByID(1)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id int) (*model.User, error) {
	user, err := s.userDAO.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	user, err := s.userDAO.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(id int, updates map[string]interface{}) error {
	return s.userDAO.Update(id, updates)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id int) error {
	return s.userDAO.Delete(id)
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(page, pageSize int) ([]model.User, int64, error) {
	return s.userDAO.GetList(page, pageSize)
}

// UpdatePassword 更新密码
func (s *userService) UpdatePassword(id int, password string) error {
	// 实际应用中应该对密码进行加密
	return s.userDAO.UpdatePassword(id, password)
}
