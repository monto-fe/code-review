package dto

type AiMessage struct {
	ID        uint   `json:"id"`
	Namespace string `json:"namespace"`
	Message   string `json:"message"`
}

// AIMessage AI消息模型
type AIMessage struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	ProjectID  uint   `gorm:"index;not null" json:"project_id"`
	Content    string `gorm:"type:text;not null" json:"content"`
	Status     int8   `gorm:"default:1" json:"status"` // 1: 正常, -1: 删除
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
}

// AIMessageResponse AI消息响应模型
type AIMessageResponse struct {
	ID         uint   `json:"id"`
	ProjectID  uint   `json:"project_id"`
	Content    string `json:"content"`
	Status     int8   `json:"status"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

// AIMessageCreateRequest AI消息创建请求
type AIMessageCreateRequest struct {
	ProjectID uint   `json:"project_id" binding:"required"` // 项目ID
	MergeURL  string `json:"merge_url" binding:"required"`  // 合并请求URL
	MergeID   uint   `json:"merge_id" binding:"required"`   // 合并请求ID
	AIModel   string `json:"ai_model" binding:"required"`   // AI模型
	Rule      int8   `json:"rule" binding:"required"`       // 规则类型
	RuleID    uint   `json:"rule_id" binding:"required"`    // 规则ID
	Result    string `json:"result" binding:"required"`     // AI评论结果
}

// AIMessageCreateResponse AI消息创建响应
type AIMessageCreateResponse struct {
	ID       uint   `json:"id"`       // 消息ID
	Comments string `json:"comments"` // 评论内容
}

// AIMessageListRequest AI消息列表请求
type AIMessageListRequest struct {
	ProjectNamespace string `form:"project_namespace"` // 项目组
	ProjectName      string `form:"project_name"`      // 项目名称
	Passed           bool   `form:"passed"`            // 是否通过
	CreateTime       int64  `form:"create_time"`       // 创建时间
	ProjectID        uint   `form:"project_id"`        // 项目ID
	ID               uint   `form:"id"`                // 消息ID
	Current          int    `form:"current"`           // 当前页码
	PageSize         int    `form:"page_size"`         // 每页数量
}

// AIMessageListResponse AI消息列表响应
type AIMessageListResponse struct {
	Data  []AIMessageItem `json:"data"`  // 消息列表
	Total int64           `json:"total"` // 总数
}

// AIMessageItem AI消息项
type AIMessageItem struct {
	ID          uint   `json:"id"`           // 消息ID
	ProjectID   uint   `json:"project_id"`   // 项目ID
	MergeURL    string `json:"merge_url"`    // 合并请求URL
	MergeID     uint   `json:"merge_id"`     // 合并请求ID
	AIModel     string `json:"ai_model"`     // AI模型
	Rule        int8   `json:"rule"`         // 规则类型
	RuleID      uint   `json:"rule_id"`      // 规则ID
	Result      string `json:"result"`       // AI评论结果
	HumanRating int8   `json:"human_rating"` // 人工评分
	Remark      string `json:"remark"`       // 备注
	Passed      bool   `json:"passed"`       // 是否通过
	CheckedBy   string `json:"checked_by"`   // 检查人
	Status      int8   `json:"status"`       // 状态
	CreateTime  int64  `json:"create_time"`  // 创建时间
	UpdateTime  int64  `json:"update_time"`  // 更新时间
}

// AIMessageUpdateRequest 更新AI消息请求
type AIMessageUpdateRequest struct {
	ID          uint   `json:"id" binding:"required"`
	HumanRating int    `json:"human_rating" binding:"required"`
	Remark      string `json:"remark" binding:"required"`
}
