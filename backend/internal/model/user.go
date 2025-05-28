package model

// User 用户模型
type User struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	OID         uint   `gorm:"column:o_id" json:"o_id"`
	Namespace   string `gorm:"type:varchar(50);not null;index:idx_user_namespace" json:"namespace"`
	User        string `gorm:"type:varchar(50);not null;index:idx_user_namespace" json:"user"`
	Name        string `gorm:"type:varchar(50);not null" json:"name"`
	Job         string `gorm:"type:varchar(50);not null" json:"job"`
	Password    string `gorm:"type:varchar(100);not null" json:"password"`
	PhoneNumber string `gorm:"type:varchar(20)" json:"phone_number"`
	Email       string `gorm:"type:varchar(100)" json:"email"`
	CreateTime  int64  `gorm:"not null" json:"create_time"`
	UpdateTime  int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (User) TableName() string {
	return TableUser
}

// UserAndRoles 用户和角色关联
type UserAndRoles struct {
	User  User        `json:"User"`
	Roles interface{} `json:"Roles"`
}

// UserResponse 用户响应
type UserResponse struct {
	Data  []User `json:"data"`
	Count int64  `json:"count"`
}

// UserLogin 用户登录请求
type UserLogin struct {
	Namespace string `json:"namespace" binding:"required"`
	User      string `json:"user" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// LoginParams 登录参数
type LoginParams struct {
	ID        uint   `json:"id,omitempty"`
	User      string `json:"user" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// UserQuery 用户查询
type UserQuery struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserListReq 用户列表请求
type UserListReq struct {
	Page      int      `json:"page" binding:"required"`
	PageSize  int      `json:"page_size" binding:"required"`
	ID        uint     `json:"id,omitempty"`
	User      []string `json:"user,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
	Name      string   `json:"name,omitempty"`
	RoleName  string   `json:"roleName,omitempty"`
}

// UserReq 用户请求
type UserReq struct {
	ID          uint    `json:"id,omitempty"`
	User        string  `json:"user" binding:"required"`
	Namespace   string  `json:"namespace" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Job         string  `json:"job" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Email       *string `json:"email,omitempty"`
	RoleIDs     []uint  `json:"role_ids" binding:"required"`
	Operator    string  `json:"operator" binding:"required"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	ID          uint    `json:"id" binding:"required"`
	User        string  `json:"user" binding:"required"`
	Namespace   string  `json:"namespace" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Job         string  `json:"job" binding:"required"`
	Password    string  `json:"password" binding:"required"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Email       *string `json:"email,omitempty"`
	RoleIDs     []uint  `json:"role_ids" binding:"required"`
	Operator    string  `json:"operator" binding:"required"`
}

// UserListQuery 用户列表查询
type UserListQuery struct {
	ID        uint   `json:"id" binding:"required"`
	User      string `json:"user" binding:"required"`
	UserName  string `json:"userName" binding:"required"`
	RoleName  string `json:"roleName" binding:"required"`
	Current   int    `json:"current" binding:"required"`
	PageSize  int    `json:"pageSize" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
}

// CreateUserParams 创建用户参数
type CreateUserParams struct {
	Namespace   string
	User        string
	Name        string
	Job         string
	Password    string
	Email       string
	PhoneNumber string
	RoleIDs     []uint
	Operator    string
}
