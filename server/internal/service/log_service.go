// server/internal/service/log_service.go
package service

import (
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/jwt"
)

// ILogService 定义日志服务接口
type ILogService interface {
	CreateLog(claims *jwt.CustomClaims, module, action, params string)
	ListLogs(page, pageSize int, orgID uint) ([]model.SysOpLog, int64, error)
}

type logService struct {
	logRepo repository.ILogRepository
}

// NewLogService 创建一个新的 logService 实例
func NewLogService(logRepo repository.ILogRepository) ILogService {
	return &logService{logRepo: logRepo}
}

// CreateLog 创建日志
func (s *logService) CreateLog(claims *jwt.CustomClaims, module, action, params string) {
	if claims == nil {
		return // 无法获取用户信息，不记录日志
	}

	log := &model.SysOpLog{
		UserID:   claims.UserID,
		Username: claims.Username,
		OrgID:    claims.OrgID,
		Module:   module,
		Action:   action,
		Params:   params,
	}

	// 异步保存日志
	go func() {
		err := s.logRepo.Create(log)
		if err != nil {
			// 在实际项目中，这里应该使用一个合适的日志库来记录错误
		}
	}()
}

// ListLogs 获取日志列表
func (s *logService) ListLogs(page, pageSize int, orgID uint) ([]model.SysOpLog, int64, error) {
	return s.logRepo.List(page, pageSize, orgID)
}
