// server/internal/repository/organization_repo.go
package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

// IOrganizationRepository 定义组织仓库接口，封装了对 SysOrganization 表的数据库操作
type IOrganizationRepository interface {
	// GetDB 返回底层的 gorm.DB 实例，用于事务等高级操作
	GetDB() *gorm.DB
	// Create 创建一个新的组织
	Create(org *model.SysOrganization) error
	// GetByID 根据ID获取组织信息
	GetByID(id uint) (*model.SysOrganization, error)
	// List 分页列出组织，可按组织类型筛选
	List(page, pageSize int, orgTypes []int8) ([]model.SysOrganization, int64, error)
	// Update 更新一个已有的组织
	Update(org *model.SysOrganization) error
	// Delete 根据ID删除一个组织
	Delete(id uint) error
}
