package service

import (
	"errors"
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/password"
)

// CreateSupplierRequest 定义了创建供应商及其管理员的请求结构
type CreateSupplierRequest struct {
	Name         string `json:"name" binding:"required"`
	ContactName  string `json:"contactName" binding:"required"`
	ContactPhone string `json:"contactPhone" binding:"required"`
	Address      string `json:"address"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
	RealName     string `json:"realName" binding:"required"`
}

// UpdateSupplierRequest 定义了更新供应商信息的请求结构
type UpdateSupplierRequest struct {
	Name         string `json:"name" binding:"required"`
	ContactName  string `json:"contactName" binding:"required"`
	ContactPhone string `json:"contactPhone" binding:"required"`
	Address      string `json:"address"`
	RealName     string `json:"realName" binding:"required"`
}

// SupplierDetailsResponse 定义了获取供应商详情的响应结构
type SupplierDetailsResponse struct {
	*model.SysOrganization
	AdminUser *struct {
		Username string `json:"username"`
		RealName string `json:"realName"`
	} `json:"adminUser"`
}

// ISupplierService 定义供应商服务接口
type ISupplierService interface {
	CreateSupplierWithAdmin(req *CreateSupplierRequest, creatorID uint, parentOrgID uint) (*model.SysOrganization, *model.SysUser, error)
	ListSuppliers(page, pageSize int, schoolID uint) ([]model.SysOrganization, int64, error)
	GetSupplierByID(id uint) (*SupplierDetailsResponse, error)
	UpdateSupplier(id uint, req *UpdateSupplierRequest) (*model.SysOrganization, error)
	UpdateSupplierStatus(id uint, isEnabled bool) error
}

// supplierService 供应商服务实现
type supplierService struct {
	orgRepo  repository.IOrganizationRepository
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
}

// NewSupplierService 创建一个新的供应商服务实例
func NewSupplierService(orgRepo repository.IOrganizationRepository, userRepo repository.IUserRepository, roleRepo repository.IRoleRepository) ISupplierService {
	return &supplierService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateSupplierWithAdmin 创建供应商及关联的管理员账号
func (s *supplierService) CreateSupplierWithAdmin(req *CreateSupplierRequest, creatorID uint, parentOrgID uint) (*model.SysOrganization, *model.SysUser, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetUserByUsername(req.Username); err == nil {
		return nil, nil, errors.New("用户名已存在")
	}

	// 查找 "供应商管理员" 角色
	supplierAdminRole, err := s.roleRepo.FindRoleByRoleKey("supplier_admin")
	if err != nil {
		return nil, nil, errors.New("未能找到 'supplier_admin' 角色，请先初始化角色")
	}

	// 哈希密码
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, nil, errors.New("密码加密失败")
	}

	// 准备组织和用户数据
	org := &model.SysOrganization{
		Name:         req.Name,
		OrgType:      2, // 2: 供应商
		ParentID:     parentOrgID,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Address:      req.Address,
		IsEnabled:    true,
	}

	user := &model.SysUser{
		Username:  req.Username,
		Password:  hashedPassword,
		RealName:  req.RealName,
		RoleID:    supplierAdminRole.ID,
		Status:    1, // 1: 正常
		CreatedBy: creatorID,
	}

	// 使用事务确保原子性
	tx := s.orgRepo.GetDB().Begin()
	if tx.Error != nil {
		return nil, nil, errors.New("开启事务失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建组织
	if err := tx.Create(org).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("创建供应商组织失败: " + err.Error())
	}

	// 2. 关联组织并创建用户
	user.OrgID = org.ID
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("创建供应商管理员失败: " + err.Error())
	}

	// 3. 更新组织信息中的主管理员ID
	if err := tx.Model(org).Update("admin_user_id", user.ID).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("更新供应商主管理员ID失败: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, errors.New("提交事务失败: " + err.Error())
	}

	return org, user, nil
}

// ListSuppliers 列出属于某个学校的所有供应商
func (s *supplierService) ListSuppliers(page, pageSize int, schoolID uint) ([]model.SysOrganization, int64, error) {
	// 供应商的 OrgType 为 3，ParentID 是其所属的学校ID
	return s.orgRepo.List(page, pageSize, []int8{3}, &schoolID)
}

// GetSupplierByID 根据ID获取供应商及其管理员信息
func (s *supplierService) GetSupplierByID(id uint) (*SupplierDetailsResponse, error) {
	org, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if org.OrgType != 2 {
		return nil, errors.New("该组织不是供应商")
	}

	user, err := s.userRepo.GetUserByID(org.AdminUserID)
	if err != nil {
		// 如果找不到用户，可以只返回组织信息，或者返回错误，这里选择返回错误
		return nil, errors.New("找不到关联的管理员账户")
	}

	resp := &SupplierDetailsResponse{
		SysOrganization: org,
		AdminUser: &struct {
			Username string `json:"username"`
			RealName string `json:"realName"`
		}{
			Username: user.Username,
			RealName: user.RealName,
		},
	}
	return resp, nil
}

// UpdateSupplier 更新供应商信息
func (s *supplierService) UpdateSupplier(id uint, req *UpdateSupplierRequest) (*model.SysOrganization, error) {
	org, err := s.orgRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if org.OrgType != 2 {
		return nil, errors.New("该组织不是供应商")
	}

	// 更新组织信息
	org.Name = req.Name
	org.ContactName = req.ContactName
	org.ContactPhone = req.ContactPhone
	org.Address = req.Address

	// 更新关联的管理员信息
	user, err := s.userRepo.GetUserByID(org.AdminUserID)
	if err != nil {
		return nil, errors.New("找不到关联的管理员账户")
	}
	user.RealName = req.RealName

	// 使用事务确保原子性
	tx := s.orgRepo.GetDB().Begin()
	if err := tx.Error; err != nil {
		return nil, errors.New("开启事务失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(org).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return org, nil
}

// UpdateSupplierStatus 更新供应商状态
func (s *supplierService) UpdateSupplierStatus(id uint, isEnabled bool) error {
	org, err := s.orgRepo.GetByID(id)
	if err != nil {
		return err
	}
	if org.OrgType != 2 {
		return errors.New("该组织不是供应商")
	}
	org.IsEnabled = isEnabled
	return s.orgRepo.Update(org)
}
