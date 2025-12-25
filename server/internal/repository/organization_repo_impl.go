// server/internal/repository/organization_repo_impl.go
package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository 创建一个新的 organizationRepository 实例
func NewOrganizationRepository(db *gorm.DB) IOrganizationRepository {
	return &organizationRepository{db: db}
}

func (repo *organizationRepository) GetDB() *gorm.DB {
	return repo.db
}

func (repo *organizationRepository) Create(org *model.SysOrganization) error {
	return repo.db.Create(org).Error
}

func (repo *organizationRepository) GetByID(id uint) (*model.SysOrganization, error) {
	var org model.SysOrganization
	err := repo.db.First(&org, id).Error
	return &org, err
}

func (repo *organizationRepository) List(page, pageSize int, orgTypes []int8, parentID *uint) ([]model.SysOrganization, int64, error) {
	var orgs []model.SysOrganization
	var total int64

	query := repo.db.Model(&model.SysOrganization{})
	if len(orgTypes) > 0 {
		query = query.Where("org_type IN ?", orgTypes)
	}
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&orgs).Error
	if err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

func (repo *organizationRepository) Update(org *model.SysOrganization) error {
	return repo.db.Save(org).Error
}

func (repo *organizationRepository) Delete(id uint) error {
	return repo.db.Delete(&model.SysOrganization{}, id).Error
}
