package handler

import (
	"net/http"
	"server/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SupplierHandler 封装了供应商相关的 HTTP 处理逻辑
type SupplierHandler struct {
	supplierService service.ISupplierService
}

// NewSupplierHandler 创建一个新的 SupplierHandler 实例
func NewSupplierHandler(supplierService service.ISupplierService) *SupplierHandler {
	return &SupplierHandler{
		supplierService: supplierService,
	}
}

// CreateSupplier 处理创建供应商的请求
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req service.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
		return
	}

	// 从认证中间件获取当前用户信息
	creatorIDVal, _ := c.Get("userID")
	creatorID, ok := creatorIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户ID"})
		return
	}

	parentOrgIDVal, _ := c.Get("orgID")
	parentOrgID, ok := parentOrgIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的组织ID"})
		return
	}

	org, user, err := h.supplierService.CreateSupplierWithAdmin(&req, creatorID, parentOrgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "供应商及管理员创建成功",
		"supplier": org,
		"admin":    user,
	})
}

// GetSupplierByID handles fetching a supplier by its ID
func (h *SupplierHandler) GetSupplierByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	supplier, err := h.supplierService.GetSupplierByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商未找到: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// UpdateSupplier handles updating a supplier's details
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req service.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
		return
	}

	supplier, err := h.supplierService.UpdateSupplier(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "供应商更新成功",
		"supplier": supplier,
	})
}

// UpdateSupplierStatus handles updating a supplier's enabled status
func (h *SupplierHandler) UpdateSupplierStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		IsEnabled bool `json:"isEnabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
		return
	}

	err = h.supplierService.UpdateSupplierStatus(uint(id), req.IsEnabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// ListSuppliers 处理列出供应商的请求
func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	orgIDVal, _ := c.Get("orgID")
	orgID, ok := orgIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的组织ID"})
		return
	}

	suppliers, total, err := h.supplierService.ListSuppliers(page, pageSize, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": suppliers,
		"total": total,
	})
}
