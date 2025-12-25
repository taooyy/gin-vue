// server/internal/handler/school_handler.go
package handler

import (
	"strconv"

	"server/internal/model"
	"server/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

// SchoolHandler 封装了学校相关的 HTTP 处理函数
type SchoolHandler struct {
	service service.ISchoolService
}

// NewSchoolHandler 创建一个新的 SchoolHandler 实例
func NewSchoolHandler(service service.ISchoolService) *SchoolHandler {
	return &SchoolHandler{service: service}
}

// Create 处理创建学校的请求
func (h *SchoolHandler) Create(c *gin.Context) {
	var req model.CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateSchool(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "学校创建成功"})
}

// List 处理列出学校的请求
func (h *SchoolHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	schools, total, err := h.service.ListSchools(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  schools,
		"total": total,
	})
}

// GetByID 处理根据ID获取学校的请求
func (h *SchoolHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	school, err := h.service.GetSchoolByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "学校未找到"})
		return
	}

	c.JSON(http.StatusOK, school)
}

// Update 处理更新学校的请求
func (h *SchoolHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req model.UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSchool(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "学校更新成功"})
}

// Delete 处理删除学校的请求
func (h *SchoolHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.service.DeleteSchool(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "学校删除成功"})
}
