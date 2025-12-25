// server/internal/service/school_service.go
package service

import (
	"errors"
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/password"

	"gorm.io/gorm"
)

const OrgTypeSchool int8 = 2
const RoleSchoolAdmin = "school_admin"

// ISchoolService 定义学校管理服务接口
type ISchoolService interface {
	CreateSchool(req *model.CreateSchoolRequest) error
	ListSchools(page, pageSize int) ([]model.SchoolListItem, int64, error)
	GetSchoolByID(id uint) (*model.SysOrganization, error)
	UpdateSchool(id uint, req *model.UpdateSchoolRequest) error
	DeleteSchool(id uint) error
}

// schoolService 实现了 ISchoolService 接口
type schoolService struct {
	orgRepo  repository.IOrganizationRepository
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
}

// NewSchoolService 创建一个新的 schoolService 实例
func NewSchoolService(orgRepo repository.IOrganizationRepository, userRepo repository.IUserRepository, roleRepo repository.IRoleRepository) ISchoolService {
	return &schoolService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateSchool 创建一所新学校及其管理员账号。

// 此操作是事务性的，确保学校和其管理员账号要么都成功创建，要么都不创建。

func (s *schoolService) CreateSchool(req *model.CreateSchoolRequest) error {

	// 使用 GORM 的事务功能，确保数据一致性

	return s.orgRepo.GetDB().Transaction(func(tx *gorm.DB) error {

		// 在事务中，我们需要用 tx 重新初始化 repo, 以确保所有操作都在同一个事务中

		txUserRepo := repository.NewUserRepository(tx)

		txRoleRepo := repository.NewRoleRepository(tx)

		// 1. 检查管理员用户名是否已存在

		_, err := txUserRepo.GetUserByUsername(req.AdminUsername)

		if err == nil {

			return errors.New("管理员用户名已存在")

		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {

			return err // 如果是除了“未找到”之外的其他数据库错误，则回滚事务

		}

		// 2. 获取 school_admin 角色

		schoolAdminRole, err := txRoleRepo.FindRoleByRoleKey(RoleSchoolAdmin)

		if err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {

				return errors.New("关键角色 'school_admin' 不存在，请先初始化数据")

			}

			return err

		}

		// 3. 哈希密码，增强安全性

		hashedPassword, err := password.Hash(req.AdminPassword)

		if err != nil {

			return errors.New("密码加密失败")

		}

		// 4. 创建学校组织

		school := &model.SysOrganization{

			Name: req.Name,

			OrgType: OrgTypeSchool,

			ContactName: req.ContactName,

			ContactPhone: req.ContactPhone,

			Address: req.Address,

			IsEnabled: true, // 默认启用

		}

		if err := tx.Create(school).Error; err != nil {

			return err

		}

		// 5. 创建管理员用户

		adminUser := &model.SysUser{

			OrgID: school.ID,

			Username: req.AdminUsername,

			Password: hashedPassword,

			RealName: req.AdminRealName,

			RoleID: schoolAdminRole.ID,

			Status: 1, // 正常状态

		}

		if err := tx.Create(adminUser).Error; err != nil {

			return err

		}

		// 6. 回填学校组织记录中的管理员ID

		if err := tx.Model(school).Update("admin_user_id", adminUser.ID).Error; err != nil {

			return err

		}

		return nil // 返回 nil 意味着事务成功，将被自动提交

	})

}

// ListSchools 列出所有学校，并附带管理员信息。

// 为了避免 N+1 查询问题，此函数先查询出所有学校，然后一次性批量查询出所有相关的管理员用户信息。

func (s *schoolService) ListSchools(page, pageSize int) ([]model.SchoolListItem, int64, error) {

	orgs, total, err := s.orgRepo.List(page, pageSize, []int8{OrgTypeSchool})

	if err != nil {

		return nil, 0, err

	}

	var schoolListItems []model.SchoolListItem

	if len(orgs) == 0 {

		return schoolListItems, 0, nil

	}

	// 1. 提取所有需要的管理员用户ID

	adminUserIDs := make([]uint, 0, len(orgs))

	for _, org := range orgs {

		if org.AdminUserID != 0 {

			adminUserIDs = append(adminUserIDs, org.AdminUserID)

		}

	}

	// 2. 一次性批量查询所有相关的管理员用户

	adminUsers, err := s.userRepo.FindUsersByIDs(adminUserIDs)

	if err != nil {

		return nil, 0, err

	}

	// 3. 将用户列表转换为 Map, 以便通过ID快速查找，优化性能

	adminUserMap := make(map[uint]model.SysUser)

	for _, user := range adminUsers {

		adminUserMap[user.ID] = user

	}

	// 4. 组装最终的返回列表

	for _, org := range orgs {

		item := model.SchoolListItem{

			SysOrganization: org,
		}

		if adminUser, ok := adminUserMap[org.AdminUserID]; ok {

			item.AdminUsername = adminUser.Username

		}

		schoolListItems = append(schoolListItems, item)

	}

	return schoolListItems, total, nil

}

// GetSchoolByID 根据ID获取学校信息

func (s *schoolService) GetSchoolByID(id uint) (*model.SysOrganization, error) {

	return s.orgRepo.GetByID(id)

}

// UpdateSchool 更新学校信息

func (s *schoolService) UpdateSchool(id uint, req *model.UpdateSchoolRequest) error {

	school, err := s.orgRepo.GetByID(id)

	if err != nil {

		return err

	}

	school.Name = req.Name

	school.ContactName = req.ContactName

	school.ContactPhone = req.ContactPhone

	school.Address = req.Address

	school.IsEnabled = req.IsEnabled

	return s.orgRepo.Update(school)

}

// DeleteSchool 删除一所学校及其关联的管理员账号（事务性）

func (s *schoolService) DeleteSchool(id uint) error {

	return s.orgRepo.GetDB().Transaction(func(tx *gorm.DB) error {

		// 使用临时的 repo 以确保在事务内操作

		school, err := repository.NewOrganizationRepository(tx).GetByID(id)

		if err != nil {

			return err // 学校不存在或查询错误

		}

		// 删除学校记录

		if err := tx.Delete(&model.SysOrganization{}, id).Error; err != nil {

			return err

		}

		// 如果学校有关联的管理员，也一并删除

		if school.AdminUserID != 0 {

			if err := tx.Delete(&model.SysUser{}, school.AdminUserID).Error; err != nil {

				// 如果用户记录因为某些原因已经不存在了，我们不认为这是一个需要回滚的错误

				if !errors.Is(err, gorm.ErrRecordNotFound) {

					return err

				}

			}

		}

		return nil

	})

}
