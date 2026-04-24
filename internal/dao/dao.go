package dao

import (
	"strings"

	"github.com/liuchen/gin-craft/internal/dto"
	pkgdb "github.com/liuchen/gin-craft/pkg/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// paginate 分页 scope，不回写 req。
func paginate(p *dto.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := p.NowPage
		if page <= 0 {
			page = 1
		}
		size := p.PerPage
		if size <= 0 {
			size = defaultPageSize
		}
		if size > maxPageSize {
			size = maxPageSize
		}
		return db.Offset((page - 1) * size).Limit(size)
	}
}

// defaultOrder 默认按 id 倒序
func defaultOrder() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC")
	}
}

// FirstByCondition 按条件查一条；找不到返回 gorm.ErrRecordNotFound
func FirstByCondition(db pkgdb.Database, m interface{}, cond map[string]interface{}) error {
	return db.GetDB().Where(cond).Scopes(defaultOrder()).First(m).Error
}

// FindAllByCondition 按条件查多条；orders 为空则默认 id DESC
func FindAllByCondition(db pkgdb.Database, data interface{}, cond map[string]interface{}, orders ...string) error {
	cur := db.GetDB().Where(cond)
	if len(orders) > 0 {
		cur = cur.Order(strings.Join(orders, ","))
	} else {
		cur = cur.Scopes(defaultOrder())
	}
	return errors.WithStack(cur.Find(data).Error)
}

// SaveModel 保存模型（存在则更新，不存在则创建）
func SaveModel(db pkgdb.Database, m interface{}) error {
	return db.GetDB().Save(m).Error
}

// CreateModel 创建模型
func CreateModel(db pkgdb.Database, m interface{}) error {
	return db.GetDB().Create(m).Error
}

// StartTransaction 开启事务
func StartTransaction(db pkgdb.Database, f func(tx *gorm.DB) error) error {
	return db.GetDB().Transaction(f)
}

// BatchCreateModel 批量创建
func BatchCreateModel(db pkgdb.Database, m interface{}, batchSize int) error {
	return db.GetDB().CreateInBatches(m, batchSize).Error
}
