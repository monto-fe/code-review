package model

// Namespace 命名空间模型
type Namespace struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Namespace  string `gorm:"type:varchar(50);not null" json:"namespace"`
	Parent     string `gorm:"type:varchar(50)" json:"parent"`
	Name       string `gorm:"type:varchar(50);not null" json:"name"`
	Describe   string `gorm:"type:varchar(200)" json:"describe"`
	Operator   string `gorm:"type:varchar(50);not null" json:"operator"`
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (Namespace) TableName() string {
	return TableNamespace
}

// TimeModel 时间模型
type TimeModel struct {
	CreateTime int64 `json:"create_time" binding:"required"`
	UpdateTime int64 `json:"update_time" binding:"required"`
}

// NamespaceReqList 命名空间列表请求
type NamespaceReqList struct {
	Page      int    `json:"page" binding:"required"`
	PageSize  int    `json:"page_size" binding:"required"`
	Namespace string `json:"namespace,omitempty"`
}

// NamespaceReq 创建命名空间请求
type NamespaceReq struct {
	Namespace  string `json:"namespace" binding:"required"`
	Parent     string `json:"parent"`
	Name       string `json:"name" binding:"required"`
	Describe   string `json:"describe"`
	Operator   string `json:"operator" binding:"required"`
	CreateTime int64  `json:"create_time" binding:"required"`
	UpdateTime int64  `json:"update_time" binding:"required"`
}

// NamespaceQuery 命名空间查询
type NamespaceQuery struct {
	Namespace string `json:"namespace,omitempty"`
	Limit     int    `json:"limit" binding:"required"`
	Offset    int    `json:"offset" binding:"required"`
}

// NamespaceCheck 命名空间检查
type NamespaceCheck struct {
	Namespace string `json:"namespace" binding:"required"`
	Name      string `json:"name" binding:"required"`
}
