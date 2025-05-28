package model

// RolePermission 角色权限关联
type RolePermission struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Namespace  string `gorm:"type:varchar(50);not null" json:"namespace"`
	RoleID     uint   `gorm:"not null" json:"role_id"`
	ResourceID uint   `gorm:"not null" json:"resource_id"`
	Describe   string `gorm:"type:varchar(200)" json:"describe"`
	Operator   string `gorm:"type:varchar(50);not null" json:"operator"`
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return TableRolePermission
}

// RoleSinglePermissionReq 单个角色权限请求
type RoleSinglePermissionReq struct {
	Namespace  string `json:"namespace" binding:"required"`
	RoleID     uint   `json:"role_id" binding:"required"`
	ResourceID uint   `json:"resource_id" binding:"required"`
	Operator   string `json:"operator" binding:"required"`
	Describe   string `json:"describe"`
}

// UserRoleReq 用户角色请求
type UserRoleReq struct {
	Namespace  string `json:"namespace" binding:"required"`
	User       string `json:"user" binding:"required"`
	RoleIDs    []uint `json:"role_ids" binding:"required"`
	Operator   string `json:"operator" binding:"required"`
	CreateTime int64  `json:"create_time" binding:"required"`
	UpdateTime int64  `json:"update_time" binding:"required"`
}

// PermissionQuery 权限查询
type PermissionQuery struct {
	Namespace string `json:"namespace" binding:"required"`
	RoleID    uint   `json:"role_id" binding:"required"`
	Limit     int    `json:"limit" binding:"required"`
	Offset    int    `json:"offset" binding:"required"`
}

// SelfPermissionQuery 自身权限查询
type SelfPermissionQuery struct {
	Namespace string `json:"namespace" binding:"required"`
	User      string `json:"user" binding:"required"`
	Category  string `json:"category,omitempty"`
}
