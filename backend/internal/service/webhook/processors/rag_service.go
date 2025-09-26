package processors

import (
	"code-review-go/config"
	"code-review-go/internal/cache"
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service/ai"
	"code-review-go/internal/service/common"
	"code-review-go/internal/service/gitlab_service"
	"code-review-go/internal/service/webhook/helpers"
	"code-review-go/internal/service/webhook/manager"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MainServiceManager 主服务管理器，单例模式
type MainServiceManager struct {
	coordinator *manager.ServiceCoordinator
	mu          sync.RWMutex
	initialized bool
}

// AnalysisData 分析所需的数据结构
type AnalysisData struct {
	GitlabInfo   *dto.GitlabCacheItem
	GitlabToken  string
	MergeRequest *model.MergeRequestInfo
	Diff         []model.Change
	AIConfig     *model.AIConfig
	FinalRule    string
	GitURL       string
	DiffStr      string
	CommentType  int8
}

var (
	mainManager *MainServiceManager
	mainOnce    sync.Once
)

// GetMainServiceManager 获取主服务管理器单例
func GetMainServiceManager() *MainServiceManager {
	mainOnce.Do(func() {
		mainManager = &MainServiceManager{
			coordinator: manager.NewServiceCoordinator(),
		}
	})
	return mainManager
}

// NewMainServiceManager 创建主服务管理器实例（用于测试）
func NewMainServiceManager(coordinator *manager.ServiceCoordinator) *MainServiceManager {
	return &MainServiceManager{
		coordinator: coordinator,
	}
}

// RAGServiceImpl RAG服务实现
type RAGServiceImpl struct {
	ragService  common.RAGServiceInterface
	config      *config.Config
	mu          sync.RWMutex
	initialized bool
}

// NewRAGService 创建RAG服务实例
func NewRAGService() common.RAGServiceInterface {
	return &RAGServiceImpl{}
}

// Initialize 初始化RAG服务
func (r *RAGServiceImpl) Initialize() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.initialized {
		return nil
	}

	// 加载配置
	r.config = config.LoadConfig()
	if r.config.RAGServiceURL == "" {
		return fmt.Errorf("RAG服务URL未配置，请设置RAG_SERVICE_URL环境变量")
	}

	// 创建RAG客户端
	ragService, err := gitlab_service.CreateDefaultRAGClient(r.config.RAGServiceURL)
	if err != nil {
		return fmt.Errorf("创建RAG客户端失败: %v", err)
	}

	r.ragService = ragService
	r.initialized = true
	return nil
}

// AnalyzeCodeWithRequest 分析代码
func (r *RAGServiceImpl) AnalyzeCodeWithRequest(req *dto.CodeReviewRequest) (*dto.CodeReviewResponse, error) {
	r.mu.RLock()
	if !r.initialized {
		r.mu.RUnlock()
		return nil, fmt.Errorf("RAG服务未初始化")
	}
	ragService := r.ragService
	r.mu.RUnlock()

	return ragService.AnalyzeCodeWithRequest(req)
}

// IsInitialized 检查是否已初始化
func (r *RAGServiceImpl) IsInitialized() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.initialized
}

// saveAndNotify 保存结果并发送通知
func saveAndNotify(manager *MainServiceManager, body dto.WebhookBody, result string, data *AnalysisData) error {
	// 保存结果
	aiMessage, err := manager.SaveResult(body, result, data)
	if err != nil {
		return err
	}

	// 发送通知
	var aiMessageID uint
	if aiMessage != nil {
		aiMessageID = aiMessage.ID
	}
	manager.SendNotifications(body, result, data, aiMessageID)

	return nil
}

// sendFailureNotification 发送失败通知
func sendFailureNotification(manager *MainServiceManager, body dto.WebhookBody, errorMsg string) {
	// 尝试获取数据用于通知，如果失败则使用基本信息
	data, err := manager.PrepareData(body)
	if err != nil {
		// 如果无法获取完整数据，创建一个基本的数据结构用于通知
		data = &AnalysisData{
			GitlabInfo: &dto.GitlabCacheItem{
				Config: model.GitlabInfo{
					API:           "",
					WebhookURL:    "",
					WebhookStatus: 0,
				},
			},
			CommentType: 1, // 默认普通评论
		}
	}

	// 发送失败通知
	manager.SendNotifications(body, fmt.Sprintf("❌ 代码审查失败\n\n%s", errorMsg), data, 0)
}

// Initialize 初始化主服务管理器
func (m *MainServiceManager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return nil
	}

	// 初始化协调器
	if err := m.coordinator.Initialize(); err != nil {
		return fmt.Errorf("初始化服务协调器失败: %v", err)
	}

	m.initialized = true
	return nil
}

// ValidateRequest 前置验证
func (m *MainServiceManager) ValidateRequest(body dto.WebhookBody) error {
	// 验证GitLab配置
	gitlabCache := cache.GetGitlabCache()
	fmt.Printf("GitLab缓存项数量: %d\n", len(gitlabCache))

	_, gitlabInfo, ok := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)
	if !ok {
		return fmt.Errorf("请配置gitlab Token，项目ID: %d", body.Project.ID)
	}

	// 验证分支匹配
	if !helpers.BranchMatch(gitlabInfo.Config.SourceBranch, body.ObjectAttributes.SourceBranch) {
		return fmt.Errorf("源分支不匹配")
	}
	if !helpers.BranchMatch(gitlabInfo.Config.TargetBranch, body.ObjectAttributes.TargetBranch) {
		return fmt.Errorf("目标分支不匹配")
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

// PrepareData 准备分析所需的数据
func (m *MainServiceManager) PrepareData(body dto.WebhookBody) (*AnalysisData, error) {
	// 获取GitLab配置
	gitlabCache := cache.GetGitlabCache()
	gitlabToken, gitlabInfo, _ := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)

	// 获取合并请求信息
	mergeRequest, err := gitlab_service.GetMergeRequestInfo(gitlabInfo.Config.API, strconv.Itoa(body.Project.ID), gitlabToken)
	if err != nil {
		return nil, fmt.Errorf("获取合并请求信息失败: %v", err)
	}

	// 获取差异信息
	diff := helpers.GetMergeDiff(gitlabInfo.Config.API, body.Project.ID, body.ObjectAttributes.IID, gitlabToken)

	// 获取AI配置
	aiConfig, _ := cache.GetAIConfigCache()

	// 获取规则
	finalRule, err := m.buildFinalRule(mergeRequest, gitlabInfo, gitlabToken)
	if err != nil {
		return nil, err
	}

	// 构建Git URL
	gitURL := m.buildGitURL(mergeRequest.WebURL)

	// 构建差异字符串
	diffStr := m.buildDiffString(diff)

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

// buildFinalRule 构建最终规则
func (m *MainServiceManager) buildFinalRule(mergeRequest *model.MergeRequestInfo, gitlabInfo dto.GitlabCacheItem, gitlabToken string) (string, error) {
	// 获取AI规则服务
	dbService := m.coordinator.GetDatabaseService()
	aiRuleServiceInterface := dbService.GetAIRuleService()
	if aiRuleServiceInterface == nil {
		return "", fmt.Errorf("AI规则服务未初始化")
	}

	// 类型断言
	aiRuleService, ok := aiRuleServiceInterface.(*ai.AIRuleService)
	if !ok {
		return "", fmt.Errorf("AI规则服务类型断言失败")
	}

	// 获取项目规则
	customRule, err := aiRuleService.GetCustomRuleByProjectID(mergeRequest.ProjectID)
	if err != nil {
		return "", fmt.Errorf("获取自定义规则失败: %v", err)
	}

	var currentRule string
	if customRule != nil && customRule.Rule != "" {
		currentRule = customRule.Rule
	} else {
		// 获取项目主要语言
		language, err := gitlab_service.GetDominantLanguage(gitlabInfo.Config.API, fmt.Sprintf("%d", mergeRequest.ProjectID), gitlabToken)
		if err != nil {
			return "", fmt.Errorf("获取项目语言失败: %v", err)
		}

		// 获取通用规则
		commonRule, err := aiRuleService.GetCommonRule(language)
		if err != nil {
			return "", fmt.Errorf("获取通用规则失败: %v", err)
		}
		if len(commonRule) > 0 {
			currentRule = commonRule[0].Rule
		}
	}

	// 获取GitLab配置中的Prompt
	gitlabPrompt := gitlabInfo.Config.Prompt
	if gitlabPrompt == "" {
		gitlabPrompt = "请检查代码是否符合以下要求：\n1. 代码风格是否规范\n2. 是否有潜在的安全问题\n3. 是否有性能优化空间\n4. 是否有重复代码\n5. 是否有未使用的代码\n6. 是否有更好的实现方式"
	}

	// 拼接规则
	if currentRule != "" {
		return fmt.Sprintf("%s\n%s", currentRule, gitlabPrompt), nil
	}
	return gitlabPrompt, nil
}

// buildGitURL 构建Git URL
func (m *MainServiceManager) buildGitURL(webURL string) string {
	if strings.Contains(webURL, "/-/merge_requests/") {
		parts := strings.Split(webURL, "/-/merge_requests/")
		if len(parts) > 0 {
			return parts[0] + ".git"
		}
	}
	return webURL
}

// buildDiffString 构建差异字符串
func (m *MainServiceManager) buildDiffString(diff []model.Change) string {
	var diffStr strings.Builder
	for _, change := range diff {
		diffStr.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", change.OldPath, change.NewPath))
		diffStr.WriteString(change.Diff)
		diffStr.WriteString("\n")
	}
	return diffStr.String()
}

// PerformRAGAnalysis 执行RAG分析
func (m *MainServiceManager) PerformRAGAnalysis(data *AnalysisData) (string, error) {
	// 准备RAG服务请求
	req := &dto.CodeReviewRequest{
		GitURL:      data.GitURL,
		Branch:      data.MergeRequest.SourceBranch,
		DiffContent: data.DiffStr,
		Query:       data.FinalRule,
		GitlabToken: data.GitlabToken,
	}

	// 调用RAG服务
	ragService := m.coordinator.GetRAGService()
	analysis, err := ragService.AnalyzeCodeWithRequest(req)
	if err != nil {
		return "", fmt.Errorf("调用RAG服务失败: %v", err)
	}

	// 处理RAG结果
	comments := analysis.Review
	if comments == "" {
		comments = "RAG服务未返回审查结果"
	}

	return comments, nil
}

// GenerateEnhancedPrompt 生成增强的提示词
func (m *MainServiceManager) GenerateEnhancedPrompt(ragResult string, data *AnalysisData) string {
	fmt.Printf("ragResult: %v\n", data.CommentType)
	if data.CommentType == utils.CommentTypeCommon {
		return utils.GenerateRAGEnhancedCommonPrompt(
			ragResult,
			data.FinalRule,
			data.MergeRequest.Title,
			data.MergeRequest.Description,
			data.DiffStr,
		)
	} else {
		return utils.GenerateRAGEnhancedInlinePrompt(
			ragResult,
			data.FinalRule,
			data.MergeRequest.Title,
			data.MergeRequest.Description,
			data.DiffStr,
		)
	}
}

// PerformAIEnhancement 执行AI增强分析
func (m *MainServiceManager) PerformAIEnhancement(prompt string, data *AnalysisData) (string, error) {
	// 调用AI服务
	aiService := m.coordinator.GetAIService()
	comments, err := aiService.CallAI(prompt)
	if err != nil {
		return "", fmt.Errorf("调用AI服务失败: %v", err)
	}

	return comments, nil
}

// SaveResult 保存分析结果
func (m *MainServiceManager) SaveResult(body dto.WebhookBody, comments string, data *AnalysisData) (*model.AIMessage, error) {
	// 判断是否通过
	passed := -1
	if strings.Contains(comments, "未发现Bug") || strings.Contains(comments, "通过") {
		passed = 1
	}

	// 创建AI消息记录
	aiMessage := &model.AIMessage{
		ProjectID:        uint(body.Project.ID),
		MergeURL:         body.ObjectAttributes.URL,
		ProjectName:      body.Project.Name,
		ProjectNamespace: body.Project.PathWithNamespace,
		MergeDescription: body.ObjectAttributes.Description,
		MergeID:          fmt.Sprintf("%d", body.ObjectAttributes.IID),
		AIModel:          data.AIConfig.Type,
		Rule:             model.RuleType(1),
		RuleID:           1,
		Result:           comments,
		Passed:           passed,
		CreateTime:       time.Now().Unix(),
	}

	// 保存到数据库
	dbService := m.coordinator.GetDatabaseService()
	db := dbService.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}
	if err := db.Create(aiMessage).Error; err != nil {
		return nil, fmt.Errorf("保存AI消息失败: %v", err)
	}

	return aiMessage, nil
}

// SendNotifications 发送通知
func (m *MainServiceManager) SendNotifications(body dto.WebhookBody, comments string, data *AnalysisData, aiMessageID uint) {
	// 获取通知服务
	notificationService := m.coordinator.GetNotificationService()

	// 发送GitLab评论
	apiURL := data.GitlabInfo.Config.API
	gitlabToken := data.GitlabToken

	if apiURL == "" {
		fmt.Printf("警告: GitLab API URL 为空，尝试从缓存重新获取\n")
		// 尝试从缓存重新获取
		gitlabCache := cache.GetGitlabCache()
		token, gitlabInfo, ok := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)
		if ok && gitlabInfo.Config.API != "" {
			apiURL = gitlabInfo.Config.API
			gitlabToken = token // 同时更新 token
			fmt.Printf("从缓存重新获取到 API URL: %s\n", apiURL)
		} else {
			fmt.Printf("错误: 无法获取有效的 GitLab API URL\n")
			return
		}
	}

	// 根据事件类型发送不同的通知
	switch body.ObjectKind {
	case "merge_request":
		// Merge Request 事件：发送 GitLab 评论
		err := notificationService.SendEnhancedNotification(
			apiURL,
			gitlabToken,
			body.Project.ID,
			body.ObjectAttributes.IID,
			comments,
			data.CommentType,
			data.Diff,
		)
		if err != nil {
			fmt.Printf("发送GitLab评论失败: %v\n", err)
		}

		// 推送webhook通知
		err = notificationService.SendWebhookNotification(
			data.GitlabInfo.Config.WebhookURL,
			int(data.GitlabInfo.Config.WebhookStatus),
			body.Project.PathWithNamespace,
			body.ObjectAttributes.MergeURL,
			comments,
			aiMessageID,
			data.MergeRequest,
		)
		if err != nil {
			fmt.Printf("发送Webhook通知失败: %v\n", err)
		}
	case "push":
		// Push 事件：只发送 webhook 通知，不发送 GitLab 评论
		fmt.Printf("Push事件：跳过GitLab评论，直接发送webhook通知\n")

		// 构建 Push 事件的 webhook URL
		pushWebhookURL := fmt.Sprintf("%s/-/commits/%s", data.GitlabInfo.Config.WebhookURL, body.After)

		err := notificationService.SendWebhookNotification(
			data.GitlabInfo.Config.WebhookURL,
			int(data.GitlabInfo.Config.WebhookStatus),
			body.Project.PathWithNamespace,
			pushWebhookURL,
			comments,
			aiMessageID,
			data.MergeRequest,
		)
		if err != nil {
			fmt.Printf("发送Push事件Webhook通知失败: %v\n", err)
		}
	default:
		fmt.Printf("未知的事件类型: %s，跳过通知发送\n", body.ObjectKind)
	}
}

// RAGAnalysisContext RAG分析上下文
type RAGAnalysisContext struct {
	ctx    context.Context
	cancel context.CancelFunc
	start  time.Time
}

// NewRAGAnalysisContext 创建RAG分析上下文
// func NewRAGAnalysisContext(timeout time.Duration) *RAGAnalysisContext {
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	return &RAGAnalysisContext{
// 		ctx:    ctx,
// 		cancel: cancel,
// 		start:  time.Now(),
// 	}
// }

// Done 检查是否完成
func (c *RAGAnalysisContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err 获取错误
func (c *RAGAnalysisContext) Err() error {
	return c.ctx.Err()
}

// Duration 获取耗时
func (c *RAGAnalysisContext) Duration() time.Duration {
	return time.Since(c.start)
}

// Cancel 取消操作
func (c *RAGAnalysisContext) Cancel() {
	c.cancel()
}
