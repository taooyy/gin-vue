// server/internal/service/account_service.go
package service

import (
	"errors"
	"fmt"
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/jwt"
	"server/pkg/password"

	"gorm.io/gorm"
)

// subordinateRoleMap 定义了管理员角色及其对应的下属员工角色
var subordinateRoleMap = map[string]string{
	"platform_admin": "platform_staff",
	"school_admin":   "school_staff",
	"supplier_admin": "supplier_staff",
}

// IAccountService 定义账号服务接口
type IAccountService interface {
	CreateAccount(req *model.CreateAccountRequest, creatorClaims *jwt.CustomClaims) error
	ListAccounts(creatorClaims *jwt.CustomClaims, page int, pageSize int) ([]model.SysUser, int64, error)
	UpdateAccountStatus(id uint, status int8, actorClaims *jwt.CustomClaims) error
	DeleteAccount(id uint, actorClaims *jwt.CustomClaims) error
	UpdateAccount(id uint, req *model.UpdateAccountRequest, actorClaims *jwt.CustomClaims) error
	ResetPassword(id uint, req *model.ResetPasswordRequest, actorClaims *jwt.CustomClaims) error
}

// accountService 实现了 IAccountService 接口
type accountService struct {
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
}

// NewAccountService 创建一个新的 accountService 实例
func NewAccountService(userRepo repository.IUserRepository, roleRepo repository.IRoleRepository) IAccountService {
	return &accountService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateAccount 创建一个新账号
func (s *accountService) CreateAccount(req *model.CreateAccountRequest, creatorClaims *jwt.CustomClaims) error {
	// 1. 检查用户名是否已存在
	_, err := s.userRepo.GetUserByUsername(req.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果是除了 "未找到" 之外的其他数据库错误
		return fmt.Errorf("检查用户名失败: %w", err)
	}

	// 2. 根据创建者的角色，确定新用户的角色
	subordinateRoleKey, ok := subordinateRoleMap[creatorClaims.Role]
	if !ok {
		return fmt.Errorf("无法为您的角色 [%s] 创建子账号", creatorClaims.Role)
	}

	// 3. 获取新用户角色的 RoleID
	subordinateRole, err := s.roleRepo.FindRoleByRoleKey(subordinateRoleKey)
	if err != nil {
		return fmt.Errorf("无法找到角色 [%s] 的定义", subordinateRoleKey)
	}

	// 4. 哈希密码
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 5. 创建 SysUser 对象
	newUser := &model.SysUser{
		OrgID:     creatorClaims.OrgID, // 将子账号与创建者组织绑定
		Username:  req.Username,
		Password:  hashedPassword,
		RealName:  req.RealName,
		Mobile:    req.Mobile,
		RoleID:    subordinateRole.ID, // 分配子角色的ID
		Status:    1,                  // 默认状态正常
		CreatedBy: creatorClaims.UserID,
	}

	// 6. 保存到数据库
	if err := s.userRepo.CreateUser(newUser); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// ListAccounts 获取由当前用户创建的账号列表
func (s *accountService) ListAccounts(creatorClaims *jwt.CustomClaims, page int, pageSize int) ([]model.SysUser, int64, error) {
	return s.userRepo.ListByCreator(creatorClaims.UserID, page, pageSize)
}

// UpdateAccountStatus 更新用户状态
func (s *accountService) UpdateAccountStatus(id uint, status int8, actorClaims *jwt.CustomClaims) error {
	// 1. 获取要更新的用户
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 2. 权限检查：确保操作者是该用户的创建者
	if user.CreatedBy != actorClaims.UserID {
		return errors.New("无权操作此账号")
	}

	// 3. 更新状态
	user.Status = status

	// 4. 保存到数据库
	return s.userRepo.UpdateUser(user)
}

// DeleteAccount 删除一个子账号
func (s *accountService) DeleteAccount(id uint, actorClaims *jwt.CustomClaims) error {
	// 1. 获取要删除的用户
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 2. 权限检查：确保操作者是该用户的创建者
	if user.CreatedBy != actorClaims.UserID {
		return errors.New("无权操作此账号")
	}

	// 3. 执行删除
	return s.userRepo.DeleteUserByID(id)
}

// UpdateAccount 更新用户基本信息
func (s *accountService) UpdateAccount(id uint, req *model.UpdateAccountRequest, actorClaims *jwt.CustomClaims) error {
	// 1. 获取要更新的用户
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 2. 权限检查：确保操作者是该用户的创建者
	if user.CreatedBy != actorClaims.UserID {
		return errors.New("无权操作此账号")
	}

	// 3. 更新字段
	user.RealName = req.RealName
	user.Mobile = req.Mobile

	// 4. 保存到数据库
	return s.userRepo.UpdateUser(user)
}

// ResetPassword 重置用户密码
func (s *accountService) ResetPassword(id uint, req *model.ResetPasswordRequest, actorClaims *jwt.CustomClaims) error {
	// 1. 获取要更新的用户
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 2. 权限检查：确保操作者是该用户的创建者
	if user.CreatedBy != actorClaims.UserID {
		return errors.New("无权操作此账号")
	}

	// 3. 哈希新密码
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return errors.New("新密码加密失败")
	}

	// 4. 更新密码
	user.Password = hashedPassword

	// 5. 保存到数据库
	return s.userRepo.UpdateUser(user)
}
