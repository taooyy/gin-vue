// server/internal/model/account.go
package model

// CreateAccountRequest 定义了创建新账号时，handler 层接收的前端请求体
type CreateAccountRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"realName" binding:"required"`
	Mobile   string `json:"mobile"`
	// 注意：新用户的 OrgID 和 RoleID 将由后端服务根据创建者的身份自动确定，
	// 而不是由前端传递，以确保安全。
}
