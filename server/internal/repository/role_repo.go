// server/internal/repository/role_repo.go
package repository

import "server/internal/model"

// IRoleRepository 定义角色仓库接口
type IRoleRepository interface {
	FindRoleByRoleKey(roleKey string) (*model.SysRole, error)
	FindRoleByID(id uint) (*model.SysRole, error)
}
