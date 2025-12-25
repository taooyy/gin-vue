// server/internal/repository/log_repo.go
package repository

import "server/internal/model"

// ILogRepository 定义日志仓库接口
type ILogRepository interface {
	Create(log *model.SysOpLog) error
	List(page, pageSize int, orgID uint) ([]model.SysOpLog, int64, error)
}
