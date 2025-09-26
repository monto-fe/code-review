package processors

import (
	"fmt"
	"strings"
	"time"

	"code-review-go/internal/cache"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service/gitlab_service"
	"code-review-go/internal/service/webhook/helpers"
)

// PushService Push 事件服务
type PushService struct {
	// 这里应该注入实际的依赖服务
	// 暂时使用 nil，实际使用时需要注入
}

// NewPushService 创建 Push 事件服务
func NewPushService() *PushService {
	return &PushService{}
}

// ProcessPush 处理 Push 事件
func (s *PushService) ProcessPush(body dto.WebhookBody) error {
	startTime := time.Now()
	fmt.Printf("开始处理Push事件: ProjectID=%d, Branch=%s\n",
		body.Project.ID, extractBranchFromRef(body.Ref))

	result, err := CheckPushRequestWithAI(body)
	if err != nil {
		return fmt.Errorf("AI检查失败: %v", err)
	}

	fmt.Printf("Push事件处理完成 (耗时: %v): %s\n", time.Since(startTime), result)
	return nil
}

// ProcessPushEvent 处理 Push 事件（对外接口）
func ProcessPushEvent(body dto.WebhookBody) error {
	service := NewPushService()
	return service.ProcessPush(body)
}

// CheckPushRequestWithAI 使用指定的管理器执行Push事件检查
func CheckPushRequestWithAI(body dto.WebhookBody) (string, error) {
	// 获取主服务管理器
	manager := GetMainServiceManager()
	if err := manager.Initialize(); err != nil {
		fmt.Printf("服务初始化失败: %v\n", err)
		return "", err
	}

	// 前置验证（Push事件专用）
	if err := validatePushRequest(body); err != nil {
		sendFailureNotification(manager, body, fmt.Sprintf("请求验证失败: %v", err))
		return "", err
	}

	// 获取必要数据（Push事件专用）
	data, err := preparePushData(body)
	if err != nil {
		sendFailureNotification(manager, body, fmt.Sprintf("数据准备失败: %v", err))
		return "", err
	}

	// 1. 执行RAG分析
	prompt := ""
	ragResult, err := manager.PerformRAGAnalysis(data)
	if err == nil {
		prompt = generatePushEnhancedPrompt(ragResult, data, body)
		fmt.Printf("RAG分析成功，提示词: %s\n", prompt)
	} else {
		// 2. 如果RAG分析失败，则使用AI检查
		prompt = generatePushAIPrompt(data, body)
		fmt.Printf("AI Push检查分析的提示词: %s\n", prompt)
	}

	// AI检查
	result, err := manager.PerformAIEnhancement(prompt, data)
	if err != nil {
		sendFailureNotification(manager, body, fmt.Sprintf("AI Push检查分析失败: %v", err))
		return "", err
	}

	// 保存结果并发送成功通知
	if err := saveAndNotify(manager, body, result, data); err != nil {
		// 即使保存失败，也发送结果通知
		manager.SendNotifications(body, result, data, 0)
		return result, err
	}

	return result, nil
}

// validatePushRequest 验证Push事件请求
func validatePushRequest(body dto.WebhookBody) error {
	// 验证GitLab配置
	gitlabCache := cache.GetGitlabCache()
	fmt.Printf("GitLab缓存项数量: %d\n", len(gitlabCache))

	_, gitlabInfo, ok := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)
	if !ok {
		return fmt.Errorf("请配置gitlab Token，项目ID: %d", body.Project.ID)
	}

	// 验证分支匹配（Push事件使用Ref字段）
	branch := extractBranchFromRef(body.Ref)
	if !helpers.BranchMatch(gitlabInfo.Config.TargetBranch, branch) {
		return fmt.Errorf("分支不匹配，期望: %s, 实际: %s", gitlabInfo.Config.TargetBranch, branch)
	}

	// 验证AI配置
	aiConfig, ok := cache.GetAIConfigCache()
	if !ok {
		return fmt.Errorf("初始化AI配置管理器失败")
	}
	if aiConfig.APIURL == "" {
		return fmt.Errorf("AI服务URL未配置")
	}

	return nil
}

// preparePushData 准备Push事件分析所需的数据
func preparePushData(body dto.WebhookBody) (*AnalysisData, error) {
	// 获取GitLab配置
	gitlabCache := cache.GetGitlabCache()
	gitlabToken, gitlabInfo, _ := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)

	// 为Push事件构建MergeRequest信息
	mergeRequest := &model.MergeRequestInfo{
		Title:        fmt.Sprintf("Push to %s", extractBranchFromRef(body.Ref)),
		Description:  buildPushDescription(body),
		SourceBranch: extractBranchFromRef(body.Ref),
		TargetBranch: extractBranchFromRef(body.Ref),
		WebURL:       fmt.Sprintf("%s/-/commits/%s", gitlabInfo.Config.WebhookURL, body.After),
	}

	// 获取Push事件的代码内容
	diff := getPushCodeContent(gitlabInfo.Config.API, body.Project.ID, body.After, gitlabToken)

	// 获取AI配置
	aiConfig, _ := cache.GetAIConfigCache()

	// 获取规则
	finalRule, err := buildPushFinalRule(mergeRequest, gitlabInfo, gitlabToken)
	if err != nil {
		return nil, err
	}

	// 构建Git URL
	gitURL := buildPushGitURL(gitlabInfo.Config.WebhookURL, extractBranchFromRef(body.Ref))

	// 构建差异字符串
	diffStr := buildDiffString(diff)

	return &AnalysisData{
		GitlabInfo:   &gitlabInfo,
		GitlabToken:  gitlabToken,
		CommentType:  gitlabInfo.CommentType,
		MergeRequest: mergeRequest,
		Diff:         diff,
		AIConfig:     aiConfig,
		FinalRule:    finalRule,
		GitURL:       gitURL,
		DiffStr:      diffStr,
	}, nil
}

// generatePushEnhancedPrompt 生成Push事件的增强提示词
func generatePushEnhancedPrompt(ragResult string, data *AnalysisData, body dto.WebhookBody) string {
	if data.CommentType == utils.CommentTypeCommon {
		return utils.GenerateRAGEnhancedCommonPrompt(
			ragResult,
			data.FinalRule,
			fmt.Sprintf("Push to %s", extractBranchFromRef(body.Ref)),
			buildPushDescription(body),
			data.DiffStr,
		)
	} else {
		return utils.GenerateRAGEnhancedInlinePrompt(
			ragResult,
			data.FinalRule,
			fmt.Sprintf("Push to %s", extractBranchFromRef(body.Ref)),
			buildPushDescription(body),
			data.DiffStr,
		)
	}
}

// generatePushAIPrompt 生成Push事件的AI提示词
func generatePushAIPrompt(data *AnalysisData, body dto.WebhookBody) string {
	gitlabPrompt := data.GitlabInfo.Config.Prompt
	// 如果 gitlabPrompt 为空，使用默认的 gitlabPrompt, 2关闭自定义配置
	if gitlabPrompt == "" || data.GitlabInfo.RuleCheckStatus == 2 {
		gitlabPrompt = utils.CodeReviewPrompt
	}

	// 生成提示词
	// 读取gitlab中的评论类型，然后选择不同的提示词
	if data.GitlabInfo.CommentType == utils.CommentTypeCommon {
		return utils.GenerateAICheckCommonPrompt(
			gitlabPrompt,
			fmt.Sprintf("Push to %s", extractBranchFromRef(body.Ref)),
			buildPushDescription(body),
			data.DiffStr,
		)
	} else {
		return utils.GenerateAICheckInlinePrompt(
			gitlabPrompt,
			fmt.Sprintf("Push to %s", extractBranchFromRef(body.Ref)),
			buildPushDescription(body),
			data.DiffStr,
		)
	}
}

// buildPushDescription 构建Push事件的描述
func buildPushDescription(body dto.WebhookBody) string {
	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("Push事件: %s -> %s\n", body.Before, body.After))
	desc.WriteString(fmt.Sprintf("提交数量: %d\n", body.TotalCommitsCount))
	desc.WriteString("提交列表:\n")

	for i, commit := range body.Commits {
		if i >= 5 { // 最多显示5个提交
			desc.WriteString(fmt.Sprintf("  ... 还有 %d 个提交\n", len(body.Commits)-5))
			break
		}
		commitID := commit.ID
		if len(commitID) > 8 {
			commitID = commitID[:8]
		}
		desc.WriteString(fmt.Sprintf("  %s: %s\n", commitID, commit.Message))
	}

	return desc.String()
}

// buildPushFinalRule 构建Push事件的最终规则
func buildPushFinalRule(mergeRequest *model.MergeRequestInfo, gitlabInfo dto.GitlabCacheItem, gitlabToken string) (string, error) {
	// 简化实现，返回默认规则
	return "请分析Push事件中的代码变更，重点关注：\n1. 代码质量和规范\n2. 潜在的安全问题\n3. 性能优化建议\n4. 最佳实践建议", nil
}

// buildPushGitURL 构建Push事件的Git URL
func buildPushGitURL(webhookURL, branch string) string {
	// 从webhook URL中提取项目路径
	if strings.Contains(webhookURL, "/-/") {
		baseURL := strings.Split(webhookURL, "/-/")[0]
		return fmt.Sprintf("%s.git", baseURL)
	}
	return webhookURL
}

// getPushCodeContent 获取Push事件的代码内容
func getPushCodeContent(api string, projectID int, commitSHA, token string) []model.Change {
	// 调用GitLab API获取当前commit的代码内容
	if api != "" && projectID != 0 && commitSHA != "" {
		codeContent, err := gitlab_service.GetCommitCodeContent(api, projectID, commitSHA, token)
		if err != nil {
			fmt.Printf("获取Push代码内容失败: %v\n", err)
			return []model.Change{}
		}
		fmt.Printf("获取Push代码内容成功: ProjectID=%d, CommitSHA=%s, 文件数=%d\n",
			projectID, commitSHA, len(codeContent))
		return codeContent
	}
	fmt.Printf("获取Push代码内容: 参数不完整, ProjectID=%d, CommitSHA=%s\n", projectID, commitSHA)
	return []model.Change{}
}

// buildDiffString 构建差异字符串
func buildDiffString(diff []model.Change) string {
	var diffStr strings.Builder
	for _, change := range diff {
		diffStr.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", change.OldPath, change.NewPath))
		diffStr.WriteString(change.Diff)
		diffStr.WriteString("\n")
	}
	return diffStr.String()
}

// extractBranchFromRef 从ref中提取分支名称
func extractBranchFromRef(ref string) string {
	if strings.HasPrefix(ref, "refs/heads/") {
		return strings.TrimPrefix(ref, "refs/heads/")
	}
	return ref
}
