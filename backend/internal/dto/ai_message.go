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
	Passed           int    `form:"passed"`            // 是否通过 -1未通过 1通过
	CreateTime       int64  `form:"create_time"`       // 创建时间
	ProjectID        uint   `form:"project_id"`        // 项目ID
	ID               uint   `form:"id"`                // 消息ID
	Current          int    `form:"current"`           // 当前页码
	PageSize         int    `form:"page_size"`         // 每页数量
}

// AIMessageListResponse AI消息列表响应
type AIMessageListResponse struct {
	Data  []AIMessageItem `json:"data"`  // 消息列表
	Count int64           `json:"count"` // 总数
}

// AIMessageItem AI消息项
type AIMessageItem struct {
	ID               uint   `json:"id"`                // 消息ID
	ProjectID        uint   `json:"project_id"`        // 项目ID
	ProjectName      string `json:"project_name"`      // 项目名称
	ProjectNamespace string `json:"project_namespace"` // 项目组
	MergeDescription string `json:"merge_description"` // 合并请求描述
	MergeURL         string `json:"merge_url"`         // 合并请求URL
	MergeID          uint   `json:"merge_id"`          // 合并请求ID
	AIModel          string `json:"ai_model"`          // AI模型
	Rule             int8   `json:"rule"`              // 规则类型
	RuleID           uint   `json:"rule_id"`           // 规则ID
	Result           string `json:"result"`            // AI评论结果
	HumanRating      int8   `json:"human_rating"`      // 人工评分
	Remark           string `json:"remark"`            // 备注
	Passed           int    `json:"passed"`            // 是否通过
	CheckedBy        string `json:"checked_by"`        // 检查人
	Status           int8   `json:"status"`            // 状态
	CreateTime       int64  `json:"create_time"`       // 创建时间
	UpdateTime       int64  `json:"update_time"`       // 更新时间
}

// AIMessageUpdateRequest 更新AI消息请求
type AIMessageUpdateRequest struct {
	ID          uint   `json:"id" binding:"required"`
	HumanRating int    `json:"human_rating"`
	Remark      string `json:"remark"`
}

type CheckCountRequest struct {
	Passed           int    `form:"passed"` // -1未通过 1通过
	StartTime        int64  `form:"start_time"`
	EndTime          int64  `form:"end_time"`
	ProjectID        uint   `form:"project_id"`
	ProjectNamespace string `form:"project_namespace"`
	ProjectName      string `form:"project_name"`
}

type CheckCountResponse struct {
	Count int `json:"count"`
}

type AIProblemCountRequest struct {
	StartTime        int64  `form:"start_time"`
	EndTime          int64  `form:"end_time"`
	ProjectID        uint   `form:"project_id"`
	ProjectNamespace string `form:"project_namespace"`
	ProjectName      string `form:"project_name"`
}

type HumanRatingStat struct {
	Level int   `json:"level"`
	Count int64 `json:"count"`
}

type AIProblemCountResponse struct {
	Data   []HumanRatingStat `json:"data"`
	Period string            `json:"period"`
}

type AIMessageWrite struct {
	ProjectID        uint   `gorm:"not null" json:"project_id"`
	ProjectName      string `gorm:"type:varchar(200);not null" json:"project_name"`
	ProjectNamespace string `gorm:"type:varchar(200);not null" json:"project_namespace"`
	MergeDescription string `gorm:"type:text;not null" json:"merge_description"`
	MergeURL         string `gorm:"type:varchar(200);not null" json:"merge_url"`
	MergeID          string `gorm:"type:varchar(50);not null" json:"merge_id"`
	AIModel          string `gorm:"type:varchar(50);not null" json:"ai_model"`
	Rule             int    `gorm:"type:tinyint;not null" json:"rule"`
	RuleID           uint   `gorm:"not null" json:"rule_id"`
	Result           string `gorm:"type:text;not null" json:"result"`
	Passed           int    `gorm:"type:tinyint" json:"passed,omitempty"`
	CreateTime       int64  `gorm:"not null" json:"create_time"`
}
