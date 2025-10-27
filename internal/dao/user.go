package dao

import (
	dtoUser "github.com/liuchen/gin-craft/internal/dto/user"
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

// GetByUsername 根据用户名获取用户
func (d *UserDAO) GetByUsername(username string) (*model.User, error) {
	var user model.User
	db := database.GetDB()

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (d *UserDAO) GetByEmail(email string) (*model.User, error) {
	var user model.User
	db := database.GetDB()

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (d *UserDAO) Create(user *model.User) error {
	return database.GetDB().Create(user).Error
}

// Update 更新用户
func (d *UserDAO) Update(id uint, updates map[string]interface{}) error {
	return database.GetDB().Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除用户
func (d *UserDAO) Delete(id uint) error {
	return database.GetDB().Delete(&model.User{}, id).Error
}

// ExistsByUsername 检查用户名是否存在
func (d *UserDAO) ExistsByUsername(username string) (bool, error) {
	var count int64
	db := database.GetDB()

	err := db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (d *UserDAO) ExistsByEmail(email string) (bool, error) {
	var count int64
	db := database.GetDB()

	err := db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// GetList 获取用户列表
func (d *UserDAO) GetList(req *dtoUser.ListRequest) ([]model.User, error) {
	var users []model.User

	db := database.GetDB()

	// 计算总数
	if err := db.Model(&model.User{}).Count(&req.Total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	if err := db.Scopes(paginate(&req.Pagination)).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// UpdatePassword 更新密码
func (d *UserDAO) UpdatePassword(id int, password string) error {
	db := database.GetDB()

	result := db.Model(&model.User{}).Where("id = ?", id).Update("password", password)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
