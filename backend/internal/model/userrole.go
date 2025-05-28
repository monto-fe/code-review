package model

// UserRole 用户角色关联
type UserRole struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Namespace  string `gorm:"type:varchar(50);not null" json:"namespace"`
	User       string `gorm:"type:varchar(50);not null" json:"user"`
	RoleID     uint   `gorm:"not null" json:"role_id"`
	Status     int    `gorm:"not null" json:"status"`
	Operator   string `gorm:"type:varchar(50);not null" json:"operator"`
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return TableUserRole
}
