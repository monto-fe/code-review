package model

// GitlabInfo Gitlab 信息模型
type GitlabInfo struct {
	ID               uint   `json:"id"`
	API              string `json:"api"`
	Token            string `json:"token"`
	WebhookURL       string `json:"webhook_url,omitempty"`
	WebhookName      string `json:"webhook_name,omitempty"`
	Status           int8   `json:"status"`
	GitlabVersion    string `json:"gitlab_version,omitempty"`
	Expired          int64  `json:"expired,omitempty"`
	GitlabURL        string `json:"gitlab_url,omitempty"`
	SourceBranch     string `json:"source_branch,omitempty"`
	TargetBranch     string `json:"target_branch,omitempty"`
	Prompt           string `json:"prompt,omitempty"`
	WebhookStatus    int8   `json:"webhook_status,omitempty"` // 1: 启用, 2: 禁用
	ProjectIds       string `json:"project_ids,omitempty" gorm:"type:text"`
	ProjectIdsSynced int8   `json:"project_ids_synced" gorm:"default:2"` // 1: 缓存失败, 2: 缓存中, 3: 缓存成功
	RuleCheckStatus  int8   `json:"rule_check_status,omitempty"`         // 1: 启用, 2: 禁用
	CreateTime       int64  `json:"create_time"`
	UpdateTime       int64  `json:"update_time"`
}

// TableName 指定表名
func (GitlabInfo) TableName() string {
	return TableGitlab
}

// GitlabInfoCreate Gitlab 信息创建请求
type GitlabInfoCreate struct {
	API             string `json:"api" binding:"required"`
	Token           string `json:"token" binding:"required"`
	WebhookURL      string `json:"webhook_url,omitempty"`
	WebhookName     string `json:"webhook_name,omitempty"`
	Status          int8   `json:"status" binding:"required"` // 1: 启用, 2: 禁用
	GitlabVersion   string `json:"gitlab_version,omitempty"`
	Expired         int64  `json:"expired,omitempty"`
	GitlabURL       string `json:"gitlab_url,omitempty"`
	SourceBranch    string `json:"source_branch,omitempty"`
	TargetBranch    string `json:"target_branch,omitempty"`
	Prompt          string `json:"prompt,omitempty" gorm:"type:text"`
	WebhookStatus   int8   `json:"webhook_status,omitempty"`    // 1: 启用, 2: 禁用
	RuleCheckStatus int8   `json:"rule_check_status,omitempty"` // 1: 启用, 2: 禁用
}

type GitlabInfoUpdate struct {
	ID              uint   `json:"id" binding:"required"`
	API             string `json:"api,omitempty"`
	Token           string `json:"token,omitempty"`
	WebhookURL      string `json:"webhook_url,omitempty"`
	WebhookName     string `json:"webhook_name,omitempty"`
	Status          int8   `json:"status,omitempty"` // 1: 启用, -1: 禁用
	GitlabVersion   string `json:"gitlab_version,omitempty"`
	Expired         int64  `json:"expired,omitempty"`
	GitlabURL       string `json:"gitlab_url,omitempty"`
	SourceBranch    string `json:"source_branch,omitempty"`
	TargetBranch    string `json:"target_branch,omitempty"`
	Prompt          string `json:"prompt,omitempty" gorm:"type:text"`
	WebhookStatus   int8   `json:"webhook_status,omitempty"`
	RuleCheckStatus int8   `json:"rule_check_status,omitempty"`
}

// GitlabCacheItem Gitlab 缓存项
type GitlabCacheItem struct {
	ProjectIDs      []string `json:"projectids"`
	GitlabAPI       string   `json:"gitlabAPI"`
	WebhookURL      string   `json:"webhook_url"`
	Prompt          string   `json:"prompt"`
	WebhookStatus   int8     `json:"webhook_status"`
	RuleCheckStatus int8     `json:"rule_check_status"`
}

// GitlabCache Gitlab 缓存
type GitlabCache map[string]GitlabCacheItem
