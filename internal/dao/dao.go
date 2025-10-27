package dao

import (
	"github.com/liuchen/gin-craft/internal/dto"
	"github.com/liuchen/gin-craft/pkg/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
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

func FirstByCondition(db database.Database, m interface{}, condition map[string]interface{}) error {
	return db.GetDB().Where(condition).Order(order()).Limit(1).Find(m).Error
}

func FindAllByCondition(db database.Database, data interface{}, condition map[string]interface{}, orders ...string) error {
	cur := db.GetDB().Where(condition)
	if len(orders) > 0 {
		cur.Order(strings.Join(orders, ","))
	} else {
		cur.Scopes(order())
	}
	return errors.WithStack(cur.Find(data).Error)
}

func SaveModel(db database.Database, m interface{}) error {
	return db.GetDB().Save(m).Error
}

func CreateModel(db database.Database, m interface{}) error {
	return db.GetDB().Create(m).Error
}

func StartTransaction(db database.Database, f func(tx *gorm.DB) error) error {
	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		return f(tx)
	})
}

func BatchCreateModel(db database.Database, m interface{}, batchSize int) error {
	return db.GetDB().CreateInBatches(m, batchSize).Error
}

func DeleteModelById(db database.Database, m interface{}, condition map[string]interface{}, operatorId uint) error {
	return db.GetDB().Model(m).
		Where(condition).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"updated_by": operatorId,
		}).Error
}
