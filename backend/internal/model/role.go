package model

// PermissionItems 权限项
type PermissionItems struct {
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

// Role 角色模型
type Role struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Namespace  string `gorm:"type:varchar(50);not null" json:"namespace"`
	Role       string `gorm:"type:varchar(50);not null" json:"role"`
	Name       string `gorm:"type:varchar(50);not null" json:"name"`
	Describe   string `gorm:"type:varchar(200)" json:"describe"`
	Operator   string `gorm:"type:varchar(50);not null" json:"operator"`
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
	Resources  []uint `gorm:"-" json:"resources,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return TableRole
}

// RoleReq 创建角色请求
type RoleReq struct {
	Namespace  string `json:"namespace" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Describe   string `json:"describe" binding:"required"`
	Operator   string `json:"operator" binding:"required"`
	CreateTime int64  `json:"create_time" binding:"required"`
	UpdateTime int64  `json:"update_time" binding:"required"`
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	ID          uint              `json:"id" binding:"required"`
	Namespace   string            `json:"namespace" binding:"required"`
	Role        string            `json:"role" binding:"required"`
	Name        string            `json:"name" binding:"required"`
	Describe    string            `json:"describe" binding:"required"`
	Permissions []PermissionItems `json:"permissions" binding:"required"`
	Operator    string            `json:"operator" binding:"required"`
}

// RoleListQuery 角色列表查询
type RoleListQuery struct {
	Namespace string   `json:"namespace" binding:"required"`
	Role      string   `json:"role,omitempty"`
	Name      string   `json:"name,omitempty"`
	Resource  []string `json:"resource,omitempty"`
	Current   int      `json:"current" binding:"required"`
	PageSize  int      `json:"pageSize" binding:"required"`
}

// RoleParams 角色参数
type RoleParams struct {
	ID         uint   `json:"id" binding:"required"`
	Namespace  string `json:"namespace" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Describe   string `json:"describe" binding:"required"`
	Operator   string `json:"operator" binding:"required"`
	CreateTime int64  `json:"create_time" binding:"required"`
	UpdateTime int64  `json:"update_time" binding:"required"`
}
