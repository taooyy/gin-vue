package service

import (
	"errors"
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/jwt"
)

// IAuthService 定义认证服务接口
type IAuthService interface {
	Login(req model.LoginRequest) (*model.LoginResponse, error)
}

// AuthService 实现了 IAuthService 接口
type AuthService struct {
	userRepo repository.IUserRepository
}

// NewAuthService 创建一个新的 AuthService 实例
func NewAuthService(userRepo repository.IUserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Login 处理用户登录逻辑
func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	// 1. 根据用户名（这里用前端传来的角色作为用户名）从仓库获取用户信息
	// 注意：这里的登录逻辑是模拟的，实际场景中用户名是独立于角色的
	user, err := s.userRepo.GetUserByUsername(req.Role)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 校验密码（实际项目中密码应该是加密存储的）
	if user.Password != req.Password {
		return nil, errors.New("密码不正确")
	}

	// 3. 校验角色是否匹配（这里用作双重检查）
	// 在模拟数据中，我们让用户名和角色名一致，所以 user.Username 就是角色
	if user.Username != req.Role {
		return nil, errors.New("角色选择不匹配")
	}

	// 4. 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.Username, req.Role)
	if err != nil {
		return nil, errors.New("生成 Token 失败")
	}

	// 5. 准备返回给前端的用户信息
	userInfo := model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Role:     req.Role,
		OrgID:    user.OrgID,
	}

	// 6. 构造并返回响应
	return &model.LoginResponse{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}
