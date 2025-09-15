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

// BotConfig 机器人配置
type BotConfig struct {
	BotName   string `json:"bot_name" binding:"required" example:"代码审查机器人"`                // 检测机器人名称，必填
	BotPrompt string `json:"bot_prompt" binding:"required" example:"请仔细审查代码，重点关注安全性和性能问题"` // 机器人提示词，必填
}

// ManualCheckRequest 手动触发代码审核请求
type ManualCheckRequest struct {
	MergeURL   string      `json:"merge_url" binding:"required" example:"https://gitlab.example.com/project/repo/-/merge_requests/123"` // 合并请求链接，必填
	AIModelID  uint        `json:"ai_model_id" binding:"required" example:"1"`                                                          // AI模型ID，必填
	BotConfigs []BotConfig `json:"bot_configs"`                                                                                         // 机器人配置列表，可选
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
	ID          uint        `json:"id"`           // 任务ID
	ProjectID   int         `json:"project_id"`   // 项目ID
	ProjectName string      `json:"project_name"` // 项目名称
	MergeID     int         `json:"merge_id"`     // 合并请求ID
	MergeTitle  string      `json:"merge_title"`  // 合并请求标题
	MergeURL    string      `json:"merge_url"`    // 合并请求URL
	Status      int         `json:"status"`       // 状态
	Result      string      `json:"result"`       // 审核结果
	AIModel     string      `json:"ai_model"`     // AI模型
	BotConfigs  []BotConfig `json:"bot_configs"`  // 机器人配置列表
	RuleName    string      `json:"rule_name"`    // 规则名称
	CreateTime  int64       `json:"create_time"`  // 创建时间
	UpdateTime  int64       `json:"update_time"`  // 更新时间
}

// BotReviewRequest 机器人审查请求
type BotReviewRequest struct {
	BotName          string `json:"bot_name" binding:"required" example:"security_reviewer"`           // 机器人名称，必填
	Type             string `json:"type" binding:"required" example:"OpenAI"`                          // AI类型，必填
	Model            string `json:"model" binding:"required" example:"gpt-4"`                          // AI模型，必填
	APIURL           string `json:"api_url" binding:"required" example:"https://api.openai.com/v1"`    // API地址，必填
	CodeContent      string `json:"code_content" binding:"required" example:"function test() { ... }"` // 代码内容，必填
	AdditionalPrompt string `json:"additional_prompt" example:"请特别关注性能问题"`                             // 额外提示词，可选
}

// BotReviewResult 机器人审查结果
type BotReviewResult struct {
	BotName     string `json:"bot_name"`     // 机器人名称
	BotCategory string `json:"bot_category"` // 机器人分类
	AIModel     string `json:"ai_model"`     // AI模型
	Result      string `json:"result"`       // 审查结果
	CreateTime  int64  `json:"create_time"`  // 创建时间
}

// BotRoleResponse 机器人角色响应
type BotRoleResponse struct {
	Name        string `json:"name"`        // 机器人名称
	Description string `json:"description"` // 机器人描述
	Category    string `json:"category"`    // 机器人分类
}

// BotRolesListResponse 机器人角色列表响应
type BotRolesListResponse struct {
	Roles []BotRoleResponse `json:"roles"` // 机器人角色列表
}
