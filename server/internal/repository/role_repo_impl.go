// server/internal/repository/role_repo_impl.go
package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

// roleRepository 是 IRoleRepository 的真实数据库实现
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建一个 roleRepository 的实例
func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &roleRepository{db: db}
}

// FindRoleByRoleKey 根据角色标识从数据库获取角色信息
func (r *roleRepository) FindRoleByRoleKey(roleKey string) (*model.SysRole, error) {
	var role model.SysRole
	err := r.db.Where("role_key = ?", roleKey).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindRoleByID 根据角色ID从数据库获取角色信息
func (r *roleRepository) FindRoleByID(id uint) (*model.SysRole, error) {
	var role model.SysRole
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
