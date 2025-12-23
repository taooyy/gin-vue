package repository

import (
	"errors"
	"server/internal/model"
)

// IUserRepository 定义用户仓库接口
type IUserRepository interface {
	GetUserByUsername(username string) (*model.SysUser, error)
}

// MockUserRepository 是一个模拟的用户仓库，用于测试和开发
type MockUserRepository struct {
	users map[string]*model.SysUser
}

// NewMockUserRepository 创建一个新的模拟用户仓库实例
func NewMockUserRepository() *MockUserRepository {
	// 初始化一些模拟用户数据
	// 在真实的DDD或分层架构中，SysUser 应该是一个领域对象或PO，这里为了简化直接使用
	users := make(map[string]*model.SysUser)
	users["platform_admin"] = &model.SysUser{ID: 1, OrgID: 1, Username: "platform_admin", Password: "password123", RealName: "平台管理员", RoleID: 1}
	users["school_admin"] = &model.SysUser{ID: 2, OrgID: 101, Username: "school_admin", Password: "password123", RealName: "学校管理员", RoleID: 2}
	users["supplier"] = &model.SysUser{ID: 3, OrgID: 201, Username: "supplier", Password: "password123", RealName: "供应商代表", RoleID: 3}
	users["canteen_admin"] = &model.SysUser{ID: 4, OrgID: 101, Username: "canteen_admin", Password: "password123", RealName: "食堂管理员", RoleID: 4}

	return &MockUserRepository{users: users}
}

// GetUserByUsername 根据用户名获取用户
func (r *MockUserRepository) GetUserByUsername(username string) (*model.SysUser, error) {
	if user, ok := r.users[username]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}
