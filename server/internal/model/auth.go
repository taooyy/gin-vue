package model

// LoginRequest 定义了登录请求的结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// LoginResponse 定义了登录成功后返回的结构
type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"user_info"`
}

// UserInfo 定义了返回给前端的用户基本信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Role     string `json:"role"`
	OrgID    uint   `json:"org_id"`
}
