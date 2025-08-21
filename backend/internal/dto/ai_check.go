package dto

type ProjectInfo struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PathWithNamespace string `json:"path_with_namespace"`
}

type ObjectAttributes struct {
	IID               int    `json:"iid"`
	URL               string `json:"url"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Action            string `json:"action"`
	ProjectID         int    `json:"project_id"`
	MergeURL          string `json:"merge_url"`
	State             string `json:"state"`
	SourceBranch      string `json:"source_branch"`
	TargetBranch      string `json:"target_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
	Note              string `json:"note,omitempty"`
}

type WebhookBody struct {
	ObjectKind       string           `json:"object_kind,omitempty"`
	Project          ProjectInfo      `json:"project"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
	MergeRequest     MergeRequest     `json:"merge_request,omitempty"`
}

type MergeRequest struct {
	State string `json:"state"`
}

// ManualCheckRequest 手动触发代码审核请求
type ManualCheckRequest struct {
	ProjectID int    `json:"project_id" binding:"required" example:"123"` // 项目ID
	MergeID   int    `json:"merge_id" binding:"required" example:"456"`   // 合并请求ID
	AIModel   string `json:"ai_model" example:"gpt-4"`                    // AI模型，可选
	RuleID    uint   `json:"rule_id" example:"1"`                         // 规则ID，可选
}

// ManualCheckResponse 手动触发代码审核响应
type ManualCheckResponse struct {
	TaskID        uint   `json:"task_id"`        // 任务ID
	Status        string `json:"status"`         // 任务状态
	EstimatedTime int    `json:"estimated_time"` // 预估完成时间（秒）
}

// ManualCheckHistoryRequest 手动审核历史查询请求
type ManualCheckHistoryRequest struct {
	Current   int    `json:"current" example:"1"`             // 当前页码
	PageSize  int    `json:"page_size" example:"20"`          // 每页数量
	Status    int    `json:"status" example:"0"`              // 状态筛选：0-全部，1-进行中，2-完成，3-失败
	StartDate string `json:"start_date" example:"1703001600"` // 开始日期
	EndDate   string `json:"end_date" example:"1705680000"`   // 结束日期
}

// ManualCheckHistoryResponse 手动审核历史响应
type ManualCheckHistoryResponse struct {
	List  []ManualCheckHistoryItem `json:"list"`  // 任务列表
	Total int64                    `json:"total"` // 总数
}

// ManualCheckHistoryItem 手动审核历史项
type ManualCheckHistoryItem struct {
	ID          uint   `json:"id"`           // 任务ID
	ProjectID   int    `json:"project_id"`   // 项目ID
	ProjectName string `json:"project_name"` // 项目名称
	MergeID     int    `json:"merge_id"`     // 合并请求ID
	MergeTitle  string `json:"merge_title"`  // 合并请求标题
	Status      int    `json:"status"`       // 状态
	StatusText  string `json:"status_text"`  // 状态文本
	CreateTime  int64  `json:"create_time"`  // 创建时间
	UpdateTime  int64  `json:"update_time"`  // 更新时间
}

// ManualCheckResultResponse 手动审核结果响应
type ManualCheckResultResponse struct {
	ID          uint   `json:"id"`           // 任务ID
	ProjectID   int    `json:"project_id"`   // 项目ID
	ProjectName string `json:"project_name"` // 项目名称
	MergeID     int    `json:"merge_id"`     // 合并请求ID
	MergeTitle  string `json:"merge_title"`  // 合并请求标题
	MergeURL    string `json:"merge_url"`    // 合并请求URL
	Status      int    `json:"status"`       // 状态
	Result      string `json:"result"`       // 审核结果
	AIModel     string `json:"ai_model"`     // AI模型
	RuleName    string `json:"rule_name"`    // 规则名称
	CreateTime  int64  `json:"create_time"`  // 创建时间
	UpdateTime  int64  `json:"update_time"`  // 更新时间
}
