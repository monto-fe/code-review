package model

// AIMessage AI消息模型
type AIMessage struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	HumanRating int    `gorm:"column:human_rating" json:"human_rating"`
	Remark      string `gorm:"type:text" json:"remark"`
	CreateTime  int64  `gorm:"not null" json:"create_time"`
	UpdateTime  int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (AIMessage) TableName() string {
	return TableAIMessage
}

// AIMessageUpdate AI消息更新请求
type AIMessageUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	HumanRating int    `json:"human_rating" binding:"required"`
	Remark      string `json:"remark" binding:"required"`
}

// HumanRating 人工评分枚举
type HumanRating int

const (
	HumanRatingExcellent        HumanRating = 1 // 发现Bug
	HumanRatingPartial          HumanRating = 2 // 部分有效：AI 建议部分正确
	HumanRatingNeedsImprovement HumanRating = 3 // 需要改进：AI 建议需要进一步优化
	HumanRatingInvalid          HumanRating = 4 // 无效
)

// RuleType 规则类型枚举
type RuleType int

const (
	RuleTypeCommon RuleType = 1 // 通用规则
	RuleTypeCustom RuleType = 2 // 自定义规则
)

// AImessage AI 消息模型
type AImessage struct {
	ID          uint        `gorm:"primarykey" json:"id"`
	ProjectID   uint        `gorm:"not null" json:"project_id"`
	MergeURL    string      `gorm:"type:varchar(200);not null" json:"merge_url"`
	MergeID     string      `gorm:"type:varchar(50);not null" json:"merge_id"`
	AIModel     string      `gorm:"type:varchar(50);not null" json:"ai_model"`
	Rule        RuleType    `gorm:"type:tinyint;not null" json:"rule"` // 1: common 2: custom
	RuleID      uint        `gorm:"not null" json:"rule_id"`
	Result      string      `gorm:"type:text;not null" json:"result"`
	HumanRating HumanRating `gorm:"type:tinyint;not null" json:"human_rating"`
	Remark      string      `gorm:"type:varchar(200)" json:"remark"`
	Passed      int         `gorm:"type:tinyint" json:"passed,omitempty"` // 默认 -1 未通过 1 通过
	CheckedBy   string      `gorm:"type:varchar(50)" json:"checked_by,omitempty"`
	CreateTime  int64       `gorm:"not null" json:"create_time"`
}

// TableName 指定表名
func (AImessage) TableName() string {
	return TableAIMessage
}
