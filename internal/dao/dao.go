package dao

import (
	"github.com/liuchen/gin-craft/internal/dto"
	"github.com/liuchen/gin-craft/internal/pkg/database"
	"gorm.io/gorm"
	"time"
)

func paginate(pagination *dto.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination.NowPage >= 0 {
			if pagination.NowPage < 1 {
				pagination.NowPage = 1
			}

			if pagination.PerPage < 20 {
				pagination.PerPage = 20
			}
			offset := (pagination.NowPage - 1) * pagination.PerPage
			return db.Offset(offset).Limit(pagination.PerPage)
		}

		return db
	}
}

func order() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC")
	}
}

func FirstByConditionPro(m interface{}, condition map[string]interface{}) error {
	return database.GetDB().Where(condition).Order(order()).Limit(1).Find(m).Error
}

func SaveModel(m interface{}) error {
	return database.GetDB().Save(m).Error
}

func CreateModel(m interface{}) error {
	return database.GetDB().Create(m).Error
}

func StartTransaction(f func(tx *gorm.DB) error) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		return f(tx)
	})
}

func BatchCreateModel(m interface{}, batchSize int) error {
	return database.GetDB().CreateInBatches(m, batchSize).Error
}

func DeleteModelById(m interface{}, condition map[string]interface{}, operatorId uint) error {
	return database.GetDB().Model(m).
		Where(condition).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"updated_by": operatorId,
		}).Error
}
