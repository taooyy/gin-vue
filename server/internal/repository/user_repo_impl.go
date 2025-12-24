// server/internal/repository/user_repo_impl.go
package repository

import (
	"server/internal/model"

	"gorm.io/gorm"
)

// userRepository 是 IUserRepository 的真实数据库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建一个 userRepository 的实例
// 注意返回类型是接口 IUserRepository，这是一种依赖倒置的实践
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

// GetUserByUsername 根据用户名从数据库获取用户
func (r *userRepository) GetUserByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err // gorm.ErrRecordNotFound 会被上层 service 处理
	}
	return &user, nil
}

// CreateUser 在数据库中创建一个新用户
func (r *userRepository) CreateUser(user *model.SysUser) error {
	return r.db.Create(user).Error
}

// ListByCreator 根据创建者ID分页查询用户列表
func (r *userRepository) ListByCreator(creatorID uint, page int, pageSize int) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64

	// 1. 计算总数
	err := r.db.Model(&model.SysUser{}).Where("created_by = ?", creatorID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 2. 查询分页数据
	offset := (page - 1) * pageSize
	err = r.db.Where("created_by = ?", creatorID).Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID 根据用户ID获取用户
func (r *userRepository) GetUserByID(id uint) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.First(&user, id).Error
	return &user, err
}

// UpdateUser 更新用户信息
func (r *userRepository) UpdateUser(user *model.SysUser) error {
	return r.db.Save(user).Error
}

// DeleteUserByID 根据ID删除用户
func (r *userRepository) DeleteUserByID(id uint) error {
	return r.db.Delete(&model.SysUser{}, id).Error
}
