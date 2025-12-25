// server/internal/model/school.go
package model

// CreateSchoolRequest 定义了创建学校时的请求体
type CreateSchoolRequest struct {
	Name          string `json:"name" binding:"required"`
	ContactName   string `json:"contactName"`
	ContactPhone  string `json:"contactPhone"`
	Address       string `json:"address"`
	AdminUsername string `json:"adminUsername" binding:"required"`
	AdminPassword string `json:"adminPassword" binding:"required"`
	AdminRealName string `json:"adminRealName" binding:"required"`
}

// UpdateSchoolRequest 定义了更新学校时的请求体
type UpdateSchoolRequest struct {
	Name         string `json:"name" binding:"required"`
	ContactName  string `json:"contactName"`
	ContactPhone string `json:"contactPhone"`
	Address      string `json:"address"`
	IsEnabled    bool   `json:"isEnabled"`
}

// SchoolListItem 定义了学校列表返回的结构
type SchoolListItem struct {
	SysOrganization
	AdminUsername string `json:"adminUsername"`
}
