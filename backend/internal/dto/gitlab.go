package dto

import "code-review-go/internal/model"

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
	CommentType      int8   `json:"comment_type"`
	CreateTime       int64  `json:"create_time"`
	UpdateTime       int64  `json:"update_time"`
}

// GitlabDeleteRequest Gitlab 删除请求
type GitlabDeleteRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GitlabCacheItem struct {
	Token           string
	Config          model.GitlabInfo
	ProjectIDs      []string
	Prompt          string
	WebhookURL      string
	WebhookStatus   int8
	RuleCheckStatus int8
	// 评论类型：0-普通评论，1-行级评论
	CommentType int8
}

// GitlabTokenProjectResponse Gitlab Token项目信息响应
type GitlabTokenProjectResponse struct {
	ID         uint   `json:"id"`          // Token ID
	Name       string `json:"name"`        // Token名称
	ProjectIds string `json:"project_ids"` // 项目ID列表
}

// GitlabTokenProjectListResponse Gitlab Token项目信息列表响应
type GitlabTokenProjectListResponse struct {
	Data []GitlabTokenProjectResponse `json:"data"` // Token项目信息列表
}

// GitlabCreateUpdateResponse Gitlab 创建/更新响应模型（不包含敏感信息）
type GitlabCreateUpdateResponse struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	API              string `json:"api"`
	WebhookURL       string `json:"webhook_url,omitempty"`
	WebhookName      string `json:"webhook_name,omitempty"`
	Status           int8   `json:"status"`
	GitlabVersion    string `json:"gitlab_version,omitempty"`
	Expired          int64  `json:"expired,omitempty"`
	GitlabURL        string `json:"gitlab_url,omitempty"`
	SourceBranch     string `json:"source_branch,omitempty"`
	TargetBranch     string `json:"target_branch,omitempty"`
	Prompt           string `json:"prompt,omitempty"`
	WebhookStatus    int8   `json:"webhook_status,omitempty"`
	ProjectIds       string `json:"project_ids,omitempty"`
	ProjectIdsSynced int8   `json:"project_ids_synced"`
	RuleCheckStatus  int8   `json:"rule_check_status,omitempty"`
	CommentType      int8   `json:"comment_type,omitempty"`
	CreateTime       int64  `json:"create_time"`
	UpdateTime       int64  `json:"update_time"`
}
