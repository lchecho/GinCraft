package dao

import (
	"errors"
	"sync"

	"github.com/liuchen/gin-craft/internal/model"
	"github.com/liuchen/gin-craft/internal/pkg/database"
	"gorm.io/gorm"
)

// UserDAO 用户数据访问对象
type UserDAO struct{}

var (
	userDAO     *UserDAO
	userDAOOnce sync.Once
)

// GetUserDAO 获取UserDAO单例实例
func GetUserDAO() *UserDAO {
	userDAOOnce.Do(func() {
		userDAO = &UserDAO{}
	})
	return userDAO
}

// GetByID 根据ID获取用户
func (d *UserDAO) GetByID(id int) (*model.User, error) {
	var user model.User
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (d *UserDAO) GetByUsername(username string) (*model.User, error) {
	var user model.User
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (d *UserDAO) GetByEmail(email string) (*model.User, error) {
	var user model.User
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not connected")
	}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (d *UserDAO) Create(user *model.User) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	return db.Create(user).Error
}

// Update 更新用户
func (d *UserDAO) Update(id int, updates map[string]interface{}) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	result := db.Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete 删除用户
func (d *UserDAO) Delete(id int) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	result := db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Exists 检查用户是否存在
func (d *UserDAO) Exists(id int) (bool, error) {
	var count int64
	db := database.GetDB()
	if db == nil {
		return false, errors.New("database not connected")
	}

	err := db.Model(&model.User{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// ExistsByUsername 检查用户名是否存在
func (d *UserDAO) ExistsByUsername(username string) (bool, error) {
	var count int64
	db := database.GetDB()
	if db == nil {
		return false, errors.New("database not connected")
	}

	err := db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (d *UserDAO) ExistsByEmail(email string) (bool, error) {
	var count int64
	db := database.GetDB()
	if db == nil {
		return false, errors.New("database not connected")
	}

	err := db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// GetList 获取用户列表
func (d *UserDAO) GetList(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not connected")
	}

	// 计算总数
	if err := db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdatePassword 更新密码
func (d *UserDAO) UpdatePassword(id int, password string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	result := db.Model(&model.User{}).Where("id = ?", id).Update("password", password)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Transaction 执行事务
func (d *UserDAO) Transaction(fn func(*gorm.DB) error) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not connected")
	}

	return db.Transaction(fn)
}
