package handler

import (
	"net/http"
	"server/internal/model"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 负责处理认证相关的 HTTP 请求
type AuthHandler struct {
	authService service.IAuthService
}

// NewAuthHandler 创建一个新的 AuthHandler 实例
func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary 用户登录
// @Description 根据用户名、密码和角色进行登录，成功返回 Token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   login  body   model.LoginRequest  true  "登录信息"
// @Success 200 {object} model.LoginResponse "登录成功"
// @Failure 400 {object} string "请求参数错误"
// @Failure 401 {object} string "认证失败"
// @Router /api/v1/system/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	// 1. 绑定并校验请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效", "details": err.Error()})
		return
	}

	// 2. 调用服务层处理登录逻辑
	resp, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 3. 登录成功，返回 Token 和用户信息
	c.JSON(http.StatusOK, resp)
}
