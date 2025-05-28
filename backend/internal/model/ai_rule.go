package model

// CommonRule 通用规则模型
type CommonRule struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Language    string `gorm:"type:varchar(50);not null" json:"language"`
	Rule        string `gorm:"type:text;not null" json:"rule"`
	Description string `gorm:"type:text" json:"description"`
	CreateTime  int64  `gorm:"not null" json:"create_time"`
	UpdateTime  int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (CommonRule) TableName() string {
	return TableCommonRule
}

// CustomRule 自定义规则模型
type CustomRule struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	ProjectName string `gorm:"type:varchar(100);not null" json:"project_name"`
	ProjectID   string `gorm:"type:varchar(50);not null" json:"project_id"`
	Rule        string `gorm:"type:text;not null" json:"rule"`
	Status      int    `gorm:"not null;default:1" json:"status"`
	Operator    string `gorm:"type:varchar(50);not null" json:"operator"`
	CreateTime  int64  `gorm:"not null" json:"create_time"`
	UpdateTime  int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (CustomRule) TableName() string {
	return TableCustomRule
}

// CommonRuleCreate 创建通用规则请求
type CommonRuleCreate struct {
	Name        string `json:"name" binding:"required"`
	Language    string `json:"language" binding:"required"`
	Rule        string `json:"rule" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CommonRuleUpdate 更新通用规则请求
type CommonRuleUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Language    string `json:"language" binding:"required"`
	Rule        string `json:"rule" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CustomRuleCreate 创建自定义规则请求
type CustomRuleCreate struct {
	ProjectName string `json:"project_name" binding:"required"`
	ProjectID   string `json:"project_id" binding:"required"`
	Rule        string `json:"rule" binding:"required"`
	Status      int    `json:"status" binding:"required"`
	Operator    string `json:"operator" binding:"required"`
}

// CustomRuleUpdate 更新自定义规则请求
type CustomRuleUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	ProjectName string `json:"project_name" binding:"required"`
	ProjectID   string `json:"project_id" binding:"required"`
	Rule        string `json:"rule" binding:"required"`
	Status      int    `json:"status" binding:"required"`
	Operator    string `json:"operator" binding:"required"`
}
