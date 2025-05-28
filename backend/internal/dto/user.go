package dto

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Namespace string `json:"namespace" binding:"required" example:"acl"`
	User      string `json:"user" binding:"required" example:"admin"`
	Password  string `json:"password" binding:"required" example:"12345678"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	JWTToken string `json:"jwt_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User     string `json:"user" example:"admin"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID          uint     `json:"id" example:"1"`
	Namespace   string   `json:"namespace" example:"default"`
	User        string   `json:"user" example:"admin"`
	Name        string   `json:"name" example:"管理员"`
	Job         string   `json:"job" example:"系统管理员"`
	PhoneNumber string   `json:"phone_number" example:"13800138000"`
	Email       string   `json:"email" example:"admin@example.com"`
	RoleList    []string `json:"roleList" example:"['admin','user']"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	ID        uint   `form:"id" example:"1"`
	User      string `form:"user" example:"admin"`
	UserName  string `form:"userName" example:"管理员"`
	RoleName  string `form:"roleName" example:"admin"`
	Current   int    `form:"current" binding:"required" example:"1"`
	PageSize  int    `form:"pageSize" binding:"required" example:"10"`
	Namespace string `form:"namespace" example:"default"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Data  []UserInfoResponse `json:"data"`
	Count int64              `json:"count" example:"100"`
}

// DeleteUserRequest 删除用户请求
type DeleteUserRequest struct {
	ID        uint   `json:"id" binding:"required" example:"1"`
	Namespace string `json:"namespace" binding:"required" example:"default"`
	User      string `json:"user" binding:"required" example:"test"`
}

// CreateInnerUserRequest 创建内部用户请求
type CreateInnerUserRequest struct {
	Namespace   string `json:"namespace" binding:"required"`
	User        string `json:"user" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Job         string `json:"job"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RoleIDs     []uint `json:"role_ids"`
	Operator    string `json:"operator"`
}

// UpdateInnerUserRequest 更新内部用户请求
type UpdateInnerUserRequest struct {
	ID          uint    `json:"id" binding:"required"`
	User        string  `json:"user"`
	Namespace   string  `json:"namespace" binding:"required"`
	Name        string  `json:"name"`
	Job         string  `json:"job"`
	Password    string  `json:"password"`
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	RoleIDs     []uint  `json:"role_ids"`
}
