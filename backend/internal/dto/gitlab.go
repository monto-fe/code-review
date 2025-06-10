package dto

type GitlabToken struct {
	ID        uint   `json:"id"`
	Namespace string `json:"namespace"`
	Token     string `json:"token"`
}

// GitlabInfoResponse Gitlab 信息响应模型
type GitlabInfoResponse struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	API              string `json:"api"`
	WebhookURL       string `json:"webhook_url"`
	WebhookName      string `json:"webhook_name"`
	Status           int8   `json:"status"`
	GitlabVersion    string `json:"gitlab_version"`
	Expired          int64  `json:"expired"`
	GitlabURL        string `json:"gitlab_url"`
	SourceBranch     string `json:"source_branch"`
	TargetBranch     string `json:"target_branch"`
	Prompt           string `json:"prompt"`
	ProjectIdsSynced int8   `json:"project_ids_synced"`
	WebhookStatus    int8   `json:"webhook_status"`
	RuleCheckStatus  int8   `json:"rule_check_status"`
	CreateTime       int64  `json:"create_time"`
	UpdateTime       int64  `json:"update_time"`
}

// GitlabDeleteRequest Gitlab 删除请求
type GitlabDeleteRequest struct {
	ID uint `json:"id" binding:"required"`
}
