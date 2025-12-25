// server/internal/repository/organization_repo.go
package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

// IOrganizationRepository 定义组织仓库接口
type IOrganizationRepository interface {
	GetDB() *gorm.DB
	Create(org *model.SysOrganization) error
	GetByID(id uint) (*model.SysOrganization, error)
	List(page, pageSize int, orgTypes []int8) ([]model.SysOrganization, int64, error)
	Update(org *model.SysOrganization) error
	Delete(id uint) error
}
