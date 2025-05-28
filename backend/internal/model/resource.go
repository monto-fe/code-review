package model

// CategoryEnum 资源类别枚举
type CategoryEnum string

const (
	CategoryAction CategoryEnum = "Action"
	CategoryAPI    CategoryEnum = "API"
	CategoryMenu   CategoryEnum = "Menu"
	CategoryOther  CategoryEnum = "Other"
)

// Resource 资源模型
type Resource struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Namespace  string `gorm:"type:varchar(50);not null" json:"namespace"`
	Resource   string `gorm:"type:varchar(50);not null" json:"resource"`
	Name       string `gorm:"type:varchar(50);not null" json:"name"`
	Category   string `gorm:"type:varchar(50)" json:"category"`
	Operator   string `gorm:"type:varchar(50);not null" json:"operator"`
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (Resource) TableName() string {
	return TableResource
}

// ResourceReq 创建资源请求
type ResourceReq struct {
	Namespace  string `json:"namespace" binding:"required"`
	Resource   string `json:"resource" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Category   string `json:"category"`
	Operator   string `json:"operator" binding:"required"`
	CreateTime int64  `json:"create_time" binding:"required"`
	UpdateTime int64  `json:"update_time" binding:"required"`
}

// ResourceQuery 资源查询
type ResourceQuery struct {
	Namespace string   `json:"namespace" binding:"required"`
	Resource  string   `json:"resource,omitempty"`
	Name      string   `json:"name,omitempty"`
	Category  []string `json:"category,omitempty"`
	Limit     int      `json:"limit" binding:"required"`
	Offset    int      `json:"offset" binding:"required"`
}
