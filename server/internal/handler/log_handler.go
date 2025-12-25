// server/internal/handler/log_handler.go
package handler

import (
	"strconv"

	"server/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

// LogHandler 封装了日志相关的 HTTP 处理函数
type LogHandler struct {
	service service.ILogService
}

// NewLogHandler 创建一个新的 LogHandler 实例
func NewLogHandler(service service.ILogService) *LogHandler {
	return &LogHandler{service: service}
}

// List 处理列出日志的请求
func (h *LogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	orgIDQuery := c.Query("orgId")
	orgID, _ := strconv.ParseUint(orgIDQuery, 10, 32)

	logs, total, err := h.service.ListLogs(page, pageSize, uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  logs,
		"total": total,
	})
}
