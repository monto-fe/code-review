package handlers

import (
	"fmt"
	"strings"

	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service/webhook/helpers"
	"code-review-go/internal/service/webhook/processors"
)

// MergeRequestHandler Merge Request 事件处理器
type MergeRequestHandler struct {
	service *processors.MergeRequestService
	config  *MergeRequestConfig
}

// MergeRequestConfig Merge Request 配置
type MergeRequestConfig struct {
	// 支持的分支列表
	SupportedBranches []string
	// 排除的分支列表
	ExcludedBranches []string
	// 是否启用分支过滤
	EnableBranchFilter bool
}

// NewMergeRequestHandler 创建 Merge Request 处理器
func NewMergeRequestHandler() *MergeRequestHandler {
	return &MergeRequestHandler{
		service: processors.NewMergeRequestService(),
		config: &MergeRequestConfig{
			SupportedBranches:  []string{"main", "master", "develop", "dev"},
			ExcludedBranches:   []string{"hotfix/*", "release/*"},
			EnableBranchFilter: true,
		},
	}
}

// Handle 处理 Merge Request 事件
func (h *MergeRequestHandler) Handle(body dto.WebhookBody) error {

	// 1. 防重复处理（缓存机制）
	if utils.IsDuplicateWebhook(body) {
		fmt.Printf("重复webhook请求，跳过处理: ProjectID=%d, MergeRequestID=%d\n",
			body.Project.ID, body.ObjectAttributes.IID)
		return nil
	}

	// 2. 状态验证
	if !helpers.ShouldProcessState(body) {
		fmt.Printf("跳过非opened状态的合并请求: ProjectID=%d, State=%s\n",
			body.Project.ID, body.ObjectAttributes.State)
		return nil
	}

	// 3. 分支过滤
	if h.config.EnableBranchFilter && !h.shouldProcessBranch(body) {
		fmt.Printf("跳过不匹配分支的合并请求: ProjectID=%d, TargetBranch=%s\n",
			body.Project.ID, body.ObjectAttributes.TargetBranch)
		return nil
	}

	// 4. 调用服务层处理
	return h.service.ProcessMergeRequest(body)
}

// GetEventType 获取事件类型
func (h *MergeRequestHandler) GetEventType() string {
	return "merge_request"
}

// shouldProcessBranch 判断是否应该处理该分支
func (h *MergeRequestHandler) shouldProcessBranch(body dto.WebhookBody) bool {
	targetBranch := body.ObjectAttributes.TargetBranch

	// 检查是否在支持列表中
	for _, branch := range h.config.SupportedBranches {
		if targetBranch == branch {
			return true
		}
	}

	// 检查是否在排除列表中
	for _, pattern := range h.config.ExcludedBranches {
		if strings.Contains(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "/*")
			if strings.HasPrefix(targetBranch, prefix) {
				return false
			}
		} else if targetBranch == pattern {
			return false
		}
	}

	// 默认处理所有分支
	return true
}

// UpdateConfig 更新配置
func (h *MergeRequestHandler) UpdateConfig(config *MergeRequestConfig) {
	h.config = config
	fmt.Printf("更新Merge Request处理器配置: %+v\n", config)
}
