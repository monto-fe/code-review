package service

import (
	"code-review-go/config"
	"code-review-go/internal/cache"
	"code-review-go/internal/database"
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/service/providers"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

// RAGServiceManager RAG服务管理器，单例模式
type RAGServiceManager struct {
	ragClient     *RAGClient
	config        *config.Config
	aiRuleService *AIRuleService
	db            *gorm.DB
	mu            sync.RWMutex
	initialized   bool
}

var (
	ragManager *RAGServiceManager
	ragOnce    sync.Once
)

// GetRAGServiceManager 获取RAG服务管理器单例
func GetRAGServiceManager() *RAGServiceManager {
	ragOnce.Do(func() {
		ragManager = &RAGServiceManager{}
	})
	return ragManager
}

// Initialize 初始化RAG服务管理器
func (m *RAGServiceManager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return nil
	}

	// 加载配置
	m.config = config.LoadConfig()
	if m.config.RAGServiceURL == "" {
		return fmt.Errorf("RAG服务URL未配置，请设置RAG_SERVICE_URL环境变量")
	}

	// 创建RAG客户端
	ragClient, err := NewRAGClient(m.config.RAGServiceURL)
	if err != nil {
		return fmt.Errorf("创建RAG客户端失败: %v", err)
	}
	m.ragClient = ragClient

	// 初始化数据库服务
	m.db = database.GetDB()
	m.aiRuleService = NewAIRuleService(m.db)

	m.initialized = true
	return nil
}

// CheckMergeRequestWithRAGOptimized 优化后的RAG检查函数
func CheckMergeRequestWithRAGOptimized(body dto.WebhookBody) (string, error) {
	// 获取RAG服务管理器
	manager := GetRAGServiceManager()
	if err := manager.Initialize(); err != nil {
		return "", err
	}

	// 前置验证
	if err := manager.validateRequest(body); err != nil {
		return "", err
	}

	// 获取必要数据
	data, err := manager.prepareData(body)
	if err != nil {
		return "", err
	}

	// 执行RAG分析
	ragResult, err := manager.performRAGAnalysis(data)
	if err != nil {
		return "", err
	}

	// 执行AI增强分析
	finalResult, err := manager.performAIEnhancement(ragResult, data)
	if err != nil {
		return "", err
	}

	// 保存结果
	aiMessage, err := manager.saveResult(body, finalResult, data)
	if err != nil {
		return "", err
	}

	// 发送通知
	manager.sendNotifications(body, finalResult, data, aiMessage.ID)

	return finalResult, nil
}

// validateRequest 前置验证
func (m *RAGServiceManager) validateRequest(body dto.WebhookBody) error {
	// 验证GitLab配置
	gitlabCache := cache.GetGitlabCache()
	_, gitlabInfo, ok := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)
	if !ok {
		return fmt.Errorf("请配置gitlab Token")
	}

	// 验证分支匹配
	if !branchMatch(gitlabInfo.Config.SourceBranch, body.ObjectAttributes.SourceBranch) {
		return fmt.Errorf("源分支不匹配")
	}
	if !branchMatch(gitlabInfo.Config.TargetBranch, body.ObjectAttributes.TargetBranch) {
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
}

// prepareData 准备分析所需的数据
func (m *RAGServiceManager) prepareData(body dto.WebhookBody) (*AnalysisData, error) {
	// 获取GitLab配置
	gitlabCache := cache.GetGitlabCache()
	gitlabToken, gitlabInfo, _ := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)

	// 获取合并请求信息
	mergeRequest, err := GetMergeRequestInfo(gitlabInfo.Config.API, strconv.Itoa(body.Project.ID), gitlabToken)
	if err != nil {
		return nil, fmt.Errorf("获取合并请求信息失败: %v", err)
	}

	// 获取差异信息
	diff := getMergeDiff(gitlabInfo.Config.API, body.Project.ID, body.ObjectAttributes.IID, gitlabToken)

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
		MergeRequest: mergeRequest,
		Diff:         diff,
		AIConfig:     aiConfig,
		FinalRule:    finalRule,
		GitURL:       gitURL,
		DiffStr:      diffStr,
	}, nil
}

// buildFinalRule 构建最终规则
func (m *RAGServiceManager) buildFinalRule(mergeRequest *model.MergeRequestInfo, gitlabInfo dto.GitlabCacheItem, gitlabToken string) (string, error) {
	// 获取项目规则
	customRule, err := m.aiRuleService.GetCustomRuleByProjectID(mergeRequest.ProjectID)
	if err != nil {
		return "", fmt.Errorf("获取自定义规则失败: %v", err)
	}

	var currentRule string
	if customRule != nil && customRule.Rule != "" {
		currentRule = customRule.Rule
	} else {
		// 获取项目主要语言
		language, err := GetDominantLanguage(gitlabInfo.Config.API, fmt.Sprintf("%d", mergeRequest.ProjectID), gitlabToken)
		if err != nil {
			return "", fmt.Errorf("获取项目语言失败: %v", err)
		}

		// 获取通用规则
		commonRule, err := m.aiRuleService.GetCommonRule(language)
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
func (m *RAGServiceManager) buildGitURL(webURL string) string {
	if strings.Contains(webURL, "/-/merge_requests/") {
		parts := strings.Split(webURL, "/-/merge_requests/")
		if len(parts) > 0 {
			return parts[0] + ".git"
		}
	}
	return webURL
}

// buildDiffString 构建差异字符串
func (m *RAGServiceManager) buildDiffString(diff []model.Change) string {
	var diffStr strings.Builder
	for _, change := range diff {
		diffStr.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", change.OldPath, change.NewPath))
		diffStr.WriteString(change.Diff)
		diffStr.WriteString("\n")
	}
	return diffStr.String()
}

// performRAGAnalysis 执行RAG分析
func (m *RAGServiceManager) performRAGAnalysis(data *AnalysisData) (string, error) {
	// 准备RAG服务请求
	req := &CodeReviewRequest{
		GitURL:      data.GitURL,
		Branch:      data.MergeRequest.SourceBranch,
		DiffContent: data.DiffStr,
		Query:       data.FinalRule,
		GitlabToken: data.GitlabToken,
	}

	// 调用RAG服务
	analysis, err := m.ragClient.AnalyzeCodeWithRequest(req)
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

// performAIEnhancement 执行AI增强分析
func (m *RAGServiceManager) performAIEnhancement(ragResult string, data *AnalysisData) (string, error) {
	// 根据AI配置类型选择提供者
	var provider providers.AIProvider
	switch data.AIConfig.Type {
	case "UCloud":
		provider = providers.NewUCAIProvider(data.AIConfig)
	case "DeepSeek":
		provider = providers.NewDeepSeekProvider(data.AIConfig)
	default:
		return "", fmt.Errorf("不支持的AI模型类型: %s", data.AIConfig.Type)
	}

	// 生成包含RAG上下文的提示词
	prompt := m.generateEnhancedPrompt(ragResult, data)

	// 调用AI服务
	comments, err := provider.CallAI(prompt)
	if err != nil {
		return "", fmt.Errorf("调用AI服务失败: %v", err)
	}

	return comments, nil
}

// generateEnhancedPrompt 生成增强的提示词
func (m *RAGServiceManager) generateEnhancedPrompt(ragResult string, data *AnalysisData) string {
	diffContent := data.DiffStr

	return fmt.Sprintf(`
你是一位资深的代码审查专家。基于RAG服务的初步分析结果，请进行进一步的代码审查。

### RAG服务分析结果
%s

### 审查规则
%s

### 代码信息
**标题**: %s
**描述**: %s

### 代码差异
%s

请基于RAG分析结果和上述规则进行深入审查，找出疑似Bug的地方，用中文输出。
如果有问题用Markdown表格格式输出：
      | 不符合的代码行号 | 疑似Bug | 修改建议 |
如果没有发现问题，请输出：'未发现Bug',不需要更多冗余信息。
`, ragResult, data.FinalRule, data.MergeRequest.Title, data.MergeRequest.Description, diffContent)
}

// saveResult 保存分析结果
func (m *RAGServiceManager) saveResult(body dto.WebhookBody, comments string, data *AnalysisData) (*model.AImessage, error) {
	// 判断是否通过
	passed := -1
	if strings.Contains(comments, "未发现Bug") || strings.Contains(comments, "通过") {
		passed = 1
	}

	// 创建AI消息记录
	aiMessage := &model.AImessage{
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
	if err := m.db.Create(aiMessage).Error; err != nil {
		return nil, fmt.Errorf("保存AI消息失败: %v", err)
	}

	return aiMessage, nil
}

// sendNotifications 发送通知
func (m *RAGServiceManager) sendNotifications(body dto.WebhookBody, comments string, data *AnalysisData, aiMessageID uint) {
	// 发送评论到GitLab
	postComment(data.GitlabInfo.Config.API, body.Project.ID, body.ObjectAttributes.IID, data.GitlabToken, comments)

	// 推送webhook通知
	pushWebhookIfNeeded(
		data.GitlabInfo.Config.WebhookURL,
		data.GitlabInfo.Config.WebhookStatus,
		body.Project.PathWithNamespace,
		body.ObjectAttributes.MergeURL,
		comments,
		aiMessageID,
		data.MergeRequest,
	)
}

// RAGClientPool RAG客户端连接池
type RAGClientPool struct {
	clients chan *RAGClient
	baseURL string
	mu      sync.Mutex
}

// NewRAGClientPool 创建RAG客户端连接池
func NewRAGClientPool(baseURL string, poolSize int) (*RAGClientPool, error) {
	pool := &RAGClientPool{
		clients: make(chan *RAGClient, poolSize),
		baseURL: baseURL,
	}

	// 预创建客户端
	for i := 0; i < poolSize; i++ {
		client, err := NewRAGClient(baseURL)
		if err != nil {
			return nil, err
		}
		pool.clients <- client
	}

	return pool, nil
}

// GetClient 获取客户端
func (p *RAGClientPool) GetClient() *RAGClient {
	return <-p.clients
}

// ReturnClient 归还客户端
func (p *RAGClientPool) ReturnClient(client *RAGClient) {
	select {
	case p.clients <- client:
	default:
		// 池已满，丢弃客户端
	}
}

// RAGAnalysisContext RAG分析上下文
type RAGAnalysisContext struct {
	ctx    context.Context
	cancel context.CancelFunc
	start  time.Time
}

// NewRAGAnalysisContext 创建RAG分析上下文
func NewRAGAnalysisContext(timeout time.Duration) *RAGAnalysisContext {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &RAGAnalysisContext{
		ctx:    ctx,
		cancel: cancel,
		start:  time.Now(),
	}
}

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
