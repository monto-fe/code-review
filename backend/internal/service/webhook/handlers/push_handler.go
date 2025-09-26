package handlers

import (
	"fmt"
	"strings"

	"code-review-go/internal/dto"
	"code-review-go/internal/service/webhook/processors"
)

// PushHandler Push 事件处理器
type PushHandler struct {
	service *processors.PushService
	config  *PushConfig
}

// PushConfig Push 事件配置
type PushConfig struct {
	// 支持的分支列表
	SupportedBranches []string
	// 排除的分支列表
	ExcludedBranches []string
	// 是否启用分支过滤
	EnableBranchFilter bool
	// 最小提交数量
	MinCommitsCount int
	// 是否检查文件变更
	CheckFileChanges bool
}

// NewPushHandler 创建 Push 事件处理器
func NewPushHandler() *PushHandler {
	return &PushHandler{
		service: processors.NewPushService(),
		config: &PushConfig{
			SupportedBranches:  []string{"main", "master", "develop", "dev"},
			ExcludedBranches:   []string{"hotfix/*", "release/*"},
			EnableBranchFilter: false,
			MinCommitsCount:    1,
			CheckFileChanges:   true,
		},
	}
}

// Handle 处理 Push 事件
func (h *PushHandler) Handle(body dto.WebhookBody) error {
	fmt.Printf("开始处理Push事件: ProjectID=%d, Branch=%s\n",
		body.Project.ID, h.extractBranchFromRef(body.Ref))

	// 1. 验证 Push 事件
	if !h.shouldProcessPush(body) {
		fmt.Printf("跳过Push事件: ProjectID=%d, Branch=%s\n",
			body.Project.ID, h.extractBranchFromRef(body.Ref))
		return nil
	}

	// 2. 调用服务层处理
	return h.service.ProcessPush(body)
}

// GetEventType 获取事件类型
func (h *PushHandler) GetEventType() string {
	return "push"
}

// shouldProcessPush 判断是否应该处理这个Push事件
func (h *PushHandler) shouldProcessPush(body dto.WebhookBody) bool {
	// 1. 检查分支
	branch := h.extractBranchFromRef(body.Ref)
	if branch == "" {
		fmt.Printf("❌ 无法提取分支名称: %s\n", body.Ref)
		return false
	}

	// 2. 分支过滤
	if h.config.EnableBranchFilter && !h.shouldProcessBranch(branch) {
		fmt.Printf("❌ 分支不在支持列表中: %s\n", branch)
		return false
	}

	// 3. 检查提交数量
	if body.TotalCommitsCount < h.config.MinCommitsCount {
		fmt.Printf("❌ 提交数量不足: %d < %d\n", body.TotalCommitsCount, h.config.MinCommitsCount)
		return false
	}

	// 4. 检查是否有文件变更
	if h.config.CheckFileChanges && !h.hasFileChanges(body) {
		fmt.Printf("❌ 没有文件变更: %s\n", branch)
		return false
	}

	fmt.Printf("✅ Push 事件验证通过: 分支=%s, 提交数=%d\n", branch, body.TotalCommitsCount)
	return true
}

// shouldProcessBranch 判断是否应该处理该分支
func (h *PushHandler) shouldProcessBranch(branch string) bool {
	// 检查是否在支持列表中
	for _, supportedBranch := range h.config.SupportedBranches {
		if branch == supportedBranch {
			return true
		}
	}

	// 检查是否在排除列表中
	for _, pattern := range h.config.ExcludedBranches {
		if strings.Contains(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "/*")
			if strings.HasPrefix(branch, prefix) {
				return false
			}
		} else if branch == pattern {
			return false
		}
	}

	// 默认处理所有分支
	return true
}

// hasFileChanges 检查是否有文件变更
func (h *PushHandler) hasFileChanges(body dto.WebhookBody) bool {
	for _, commit := range body.Commits {
		if len(commit.Added) > 0 || len(commit.Modified) > 0 || len(commit.Removed) > 0 {
			return true
		}
	}
	return false
}

// extractBranchFromRef 从 ref 中提取分支名称
func (h *PushHandler) extractBranchFromRef(ref string) string {
	// 从 refs/heads/feature-branch 提取 feature-branch
	if strings.HasPrefix(ref, "refs/heads/") {
		return strings.TrimPrefix(ref, "refs/heads/")
	}
	return ref
}

// UpdateConfig 更新配置
func (h *PushHandler) UpdateConfig(config *PushConfig) {
	h.config = config
	fmt.Printf("更新Push处理器配置: %+v\n", config)
}
