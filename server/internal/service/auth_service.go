package service

import (
	"errors"
	"fmt"
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/jwt"
	"server/pkg/password"
)

// IAuthService 定义认证服务接口
type IAuthService interface {
	Login(req model.LoginRequest) (*model.LoginResponse, error)
}

// AuthService 实现了 IAuthService 接口
type AuthService struct {
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
}

// NewAuthService 创建一个新的 AuthService 实例
func NewAuthService(userRepo repository.IUserRepository, roleRepo repository.IRoleRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// Login 处理用户登录逻辑
func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	// 1. 根据用户名从仓库获取用户信息
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码不正确") // 出于安全，不明确指出是用户名错了
	}

	// 2. 校验哈希密码
	if !password.Check(req.Password, user.Password) {
		return nil, errors.New("用户名或密码不正确")
	}

	// 3. 获取用户的角色信息
	role, err := s.roleRepo.FindRoleByID(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("无法获取用户角色信息: %w", err)
	}

	// 4. 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.OrgID, user.Username, role.RoleKey)
	if err != nil {
		return nil, errors.New("生成 Token 失败")
	}

	// 5. 准备返回给前端的用户信息
	userInfo := model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Role:     role.RoleKey, // 使用数据库中的权威角色
		OrgID:    user.OrgID,
	}

	// 6. 构造并返回响应
	return &model.LoginResponse{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}
