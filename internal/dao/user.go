package dao

import (
	"sync"

	dtoUser "github.com/liuchen/gin-craft/internal/dto/user"
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

// GetUserDAO 获取 UserDAO 单例实例
func GetUserDAO() *UserDAO {
	userDAOOnce.Do(func() {
		userDAO = &UserDAO{}
	})
	return userDAO
}

// GetByID 根据 ID 获取用户；找不到返回 gorm.ErrRecordNotFound
func (d *UserDAO) GetByID(id uint) (*model.User, error) {
	var u model.User
	if err := database.GetDB().First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByUsername 根据用户名获取用户
func (d *UserDAO) GetByUsername(username string) (*model.User, error) {
	var u model.User
	if err := database.GetDB().Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmail 根据邮箱获取用户
func (d *UserDAO) GetByEmail(email string) (*model.User, error) {
	var u model.User
	if err := database.GetDB().Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Create 创建用户
func (d *UserDAO) Create(u *model.User) error {
	return database.GetDB().Create(u).Error
}

// Update 白名单字段更新；updates 为空则直接返回 nil
func (d *UserDAO) Update(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	res := database.GetDB().Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete 软删除用户
func (d *UserDAO) Delete(id uint) error {
	res := database.GetDB().Delete(&model.User{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ExistsByUsername 检查用户名是否存在
func (d *UserDAO) ExistsByUsername(username string) (bool, error) {
	var cnt int64
	err := database.GetDB().Model(&model.User{}).Where("username = ?", username).Count(&cnt).Error
	return cnt > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (d *UserDAO) ExistsByEmail(email string) (bool, error) {
	var cnt int64
	err := database.GetDB().Model(&model.User{}).Where("email = ?", email).Count(&cnt).Error
	return cnt > 0, err
}

// GetList 获取用户列表（支持用户名/邮箱模糊过滤 + 分页）
func (d *UserDAO) GetList(req *dtoUser.ListRequest) ([]model.User, error) {
	q := database.GetDB().Model(&model.User{})
	if req.Username != "" {
		q = q.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		q = q.Where("email LIKE ?", "%"+req.Email+"%")
	}

	if err := q.Count(&req.Total).Error; err != nil {
		return nil, err
	}

	var users []model.User
	if err := q.Scopes(paginate(&req.Pagination), defaultOrder()).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdatePassword 更新密码
func (d *UserDAO) UpdatePassword(id uint, password string) error {
	res := database.GetDB().Model(&model.User{}).Where("id = ?", id).Update("password", password)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
