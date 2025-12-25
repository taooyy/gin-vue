// server/internal/handler/account_handler.go
package handler

import (
	"net/http"
	"strconv"

	"server/internal/model"
	"server/internal/router/middleware"
	"server/internal/service"
	"server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// UpdateStatusRequest 定义了更新状态的请求体
type UpdateStatusRequest struct {
	Status int8 `json:"status" binding:"required,oneof=1 2"` // 1=正常, 2=锁定
}

// AccountHandler 负责处理账号相关的 HTTP 请求
type AccountHandler struct {
	svc service.IAccountService
}

// NewAccountHandler 创建一个新的 AccountHandler
func NewAccountHandler(svc service.IAccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// CreateAccount godoc
// @Summary 创建子账号
// @Description 由授权用户创建一个权限更低的子账号
// @Tags Accounts
// @Accept json
// @Produce json
// @Param account body model.CreateAccountRequest true "账号信息"
// @Success 201 {object} object "{"message":"账号创建成功"}"
// @Failure 400 {object} object "{"error":"请求参数错误"}"
// @Failure 403 {object} object "{"error":"无权限操作"}"
// @Failure 409 {object} object "{"error":"用户名已存在"}"
// @Failure 500 {object} object "{"error":"内部服务器错误"}"
// @Router /api/v1/accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	// 1. 绑定和验证请求参数
	var req model.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 2. 从 Gin Context 获取用户信息
	claims, exists := c.Get(middleware.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法获取用户信息，禁止访问"})
		return
	}
	creatorClaims, ok := claims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "用户信息格式错误，禁止访问"})
		return
	}

	// 3. 调用 Service 层来处理业务逻辑
	if err := h.svc.CreateAccount(&req, creatorClaims); err != nil {
		// 根据错误类型返回不同的 HTTP 状态码
		// 这里只是一个简单的示例，实际应用中可能需要更精细的错误处理
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusCreated, gin.H{"message": "账号创建成功"})
}

// ListAccounts godoc
// @Summary 获取子账号列表
// @Description 获取由当前用户创建的子账号列表（分页）
// @Tags Accounts
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} object "{"list":[]model.SysUser, "total": 0}"
// @Failure 403 {object} object "{"error":"无权限操作"}"
// @Failure 500 {object} object "{"error":"内部服务器错误"}"
// @Router /api/v1/accounts [get]
func (h *AccountHandler) ListAccounts(c *gin.Context) {
	// 1. 解析分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 2. 从 Gin Context 获取用户信息
	claims, exists := c.Get(middleware.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法获取用户信息，禁止访问"})
		return
	}
	creatorClaims, ok := claims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "用户信息格式错误，禁止访问"})
		return
	}

	// 3. 调用 Service 层
	list, total, err := h.svc.ListAccounts(creatorClaims, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// UpdateAccountStatus godoc
// @Summary 更新子账号状态
// @Description 禁用或启用一个由当前用户创建的子账号
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Param status body UpdateStatusRequest true "新的状态"
// @Success 200 {object} object "{"message":"状态更新成功"}"
// @Failure 400 {object} object "{"error":"请求参数错误"}"
// @Failure 403 {object} object "{"error":"无权限操作"}"
// @Failure 500 {object} object "{"error":"内部服务器错误"}"
// @Router /api/v1/accounts/{id}/status [put]
func (h *AccountHandler) UpdateAccountStatus(c *gin.Context) {
	// 1. 解析路径参数 ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 2. 绑定请求体
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 3. 从 Gin Context 获取用户信息
	claims, exists := c.Get(middleware.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法获取用户信息，禁止访问"})
		return
	}
	actorClaims, ok := claims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "用户信息格式错误，禁止访问"})
		return
	}

	// 4. 调用 Service 层
	if err := h.svc.UpdateAccountStatus(uint(id), req.Status, actorClaims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// DeleteAccount godoc
// @Summary 删除子账号
// @Description 删除一个由当前用户创建的子账号
// @Tags Accounts
// @Produce json
// @Param id path int true "账号ID"
// @Success 200 {object} object "{"message":"账号删除成功"}"
// @Failure 400 {object} object "{"error":"无效的用户ID"}"
// @Failure 403 {object} object "{"error":"无权限操作"}"
// @Failure 500 {object} object "{"error":"内部服务器错误"}"
// @Router /api/v1/accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	// 1. 解析路径参数 ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 2. 从 Gin Context 获取用户信息
	claims, exists := c.Get(middleware.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法获取用户信息，禁止访问"})
		return
	}
	actorClaims, ok := claims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "用户信息格式错误，禁止访问"})
		return
	}

	// 3. 调用 Service 层
	if err := h.svc.DeleteAccount(uint(id), actorClaims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "账号删除成功"})
}

// UpdateAccount godoc
// @Summary 更新子账号基本信息
// @Description 更新一个由当前用户创建的子账号的基本信息
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Param account body model.UpdateAccountRequest true "要更新的账号信息"
// @Success 200 {object} object "{"message":"账号更新成功"}"
// @Failure 400 {object} object "{"error":"请求参数错误"}"
// @Failure 403 {object} object "{"error":"无权限操作"}"
// @Failure 500 {object} object "{"error":"内部服务器错误"}"
// @Router /api/v1/accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	// 1. 解析路径参数 ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 2. 绑定请求体
	var req model.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 3. 从 Gin Context 获取用户信息
	claims, exists := c.Get(middleware.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法获取用户信息，禁止访问"})
		return
	}
	actorClaims, ok := claims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "用户信息格式错误，禁止访问"})
		return
	}

	// 4. 调用 Service 层
	if err := h.svc.UpdateAccount(uint(id), &req, actorClaims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "账号更新成功"})
}

// ResetPassword godoc
// @Summary 重置子账号密码
// @Description 重置一个由当前用户创建的子账号的密码
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "账号ID"
// @Param password body model.ResetPasswordRequest true "新密码"
// @Success 200 {object} object "{"message":"密码重置成功"}"
// @Failure 400 {object} object "{"error":"请求参数错误"}"
// @Failure 403 {object} object "{"error":"无权限操作"}"
// @Failure 500 {object} object "{"error":"内部服务器错误"}"
// @Router /api/v1/accounts/{id}/password [put]
func (h *AccountHandler) ResetPassword(c *gin.Context) {
	// 1. 解析路径参数 ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 2. 绑定请求体
	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 3. 从 Gin Context 获取用户信息
	claims, exists := c.Get(middleware.ContextUserClaimsKey)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "无法获取用户信息，禁止访问"})
		return
	}
	actorClaims, ok := claims.(*jwt.CustomClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "用户信息格式错误，禁止访问"})
		return
	}

	// 4. 调用 Service 层
	if err := h.svc.ResetPassword(uint(id), &req, actorClaims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}
