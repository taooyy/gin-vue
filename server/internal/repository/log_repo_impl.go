// server/internal/repository/log_repo_impl.go
package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

type logRepository struct {
	db *gorm.DB
}

// NewLogRepository 创建一个新的 logRepository 实例
func NewLogRepository(db *gorm.DB) ILogRepository {
	return &logRepository{db: db}
}

// Create 创建一条操作日志
func (r *logRepository) Create(log *model.SysOpLog) error {
	return r.db.Create(log).Error
}

// List 分页列出操作日志，可按 OrgID 筛选
func (r *logRepository) List(page, pageSize int, orgID uint) ([]model.SysOpLog, int64, error) {
	var logs []model.SysOpLog
	var total int64

	query := r.db.Model(&model.SysOpLog{})

	// 如果 orgID 不为0，则按组织ID筛选
	if orgID != 0 {
		query = query.Where("org_id = ?", orgID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
