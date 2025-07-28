package service

import (
	"code-review-go/config"
	"code-review-go/internal/cache"
	"code-review-go/internal/database"
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
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
	ragService    RAGService
	config        *config.Config
	aiRuleService *AIRuleService
	db            *gorm.DB
	mu            sync.RWMutex
	initialized   bool
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

// CheckMergeRequestWithRAGOptimized 优化后的RAG检查函数
func CheckMergeRequestWithRAGOptimized(body dto.WebhookBody) (string, *AnalysisData, uint, error) {
	// 获取RAG服务管理器
	manager := GetRAGServiceManager()
	if err := manager.Initialize(); err != nil {
		return "", nil, 0, err
	}

	// 前置验证
	if err := manager.ValidateRequest(body); err != nil {
		return "", nil, 0, err
	}

	// 获取必要数据
	data, err := manager.PrepareData(body)
	if err != nil {
		return "", nil, 0, err
	}

	// 执行RAG分析
	ragResult, err := manager.performRAGAnalysis(data)
	if err != nil {
		return "", data, 0, err
	}

	// 执行AI增强分析
	prompt := manager.generateEnhancedPrompt(ragResult, data)
	finalResult, err := manager.PerformAIEnhancement(prompt, data)
	if err != nil {
		return "", data, 0, err
	}

	// 保存结果
	aiMessage, err := manager.SaveResult(body, finalResult, data)
	if err != nil {
		return "", data, 0, err
	}

	// 发送通知 - 修复空指针问题
	var aiMessageID uint
	if aiMessage != nil {
		aiMessageID = aiMessage.ID
	}
	manager.SendNotifications(body, finalResult, data, aiMessageID)

	return finalResult, data, aiMessageID, nil
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

	// 使用工厂函数创建RAG客户端
	ragService, err := CreateDefaultRAGClient(m.config.RAGServiceURL)
	if err != nil {
		return fmt.Errorf("创建RAG客户端失败: %v", err)
	}

	m.ragService = ragService

	// 初始化数据库服务
	m.db = database.GetDB()
	m.aiRuleService = NewAIRuleService(m.db)

	m.initialized = true
	return nil
}

// ValidateRequest 前置验证
func (m *RAGServiceManager) ValidateRequest(body dto.WebhookBody) error {
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

// PrepareData 准备分析所需的数据
func (m *RAGServiceManager) PrepareData(body dto.WebhookBody) (*AnalysisData, error) {
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
	analysis, err := m.ragService.AnalyzeCodeWithRequest(req)
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

// PerformAIEnhancement 执行AI增强分析
func (m *RAGServiceManager) PerformAIEnhancement(prompt string, data *AnalysisData) (string, error) {
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

	// 调用AI服务
	comments, err := provider.CallAI(prompt)
	if err != nil {
		return "", fmt.Errorf("调用AI服务失败: %v", err)
	}

	return comments, nil
}

// SaveResult 保存分析结果
func (m *RAGServiceManager) SaveResult(body dto.WebhookBody, comments string, data *AnalysisData) (*model.AImessage, error) {
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

// addQualityFeedbackCheckboxes 添加评论质量反馈复选框
func addQualityFeedbackCheckboxes(comments string) string {
	return comments + "\n\n---\n**请评价此评论的质量：**\n" +
		"- [ ] 评论内容准确且有用\n" +
		"- [ ] 精准定位并提供建议\n" +
		"- [ ] 建议不够具体\n" +
		"- [ ] 完全误导性建议\n"
}

// SendNotifications 发送通知
func (m *RAGServiceManager) SendNotifications(body dto.WebhookBody, comments string, data *AnalysisData, aiMessageID uint) {
	if data.CommentType == utils.CommentTypeCommon {
		// 为普通评论添加质量反馈复选框
		commentWithFeedback := addQualityFeedbackCheckboxes(comments)
		PostCommentToGitLab(data.GitlabInfo.Config.API, body.Project.ID, body.ObjectAttributes.IID, data.GitlabToken, commentWithFeedback)
	} else {
		// 解析所有评论内容
		allComments := ParseCommentsForLineComments(comments, data.Diff)
		fmt.Printf("解析的所有评论: %+v\n", allComments)

		// 分类评论
		lineComments, generalComments := ClassifyComments(allComments)
		fmt.Printf("行级评论: %+v\n", lineComments)
		fmt.Printf("普通评论: %+v\n", generalComments)

		// 处理行级评论
		var failedLineComments []string
		if len(lineComments) > 0 {
			var err error
			failedLineComments, err = PostLineComments(data.GitlabInfo.Config.API, body.Project.ID, body.ObjectAttributes.IID, data.GitlabToken, lineComments, data.Diff)
			if err != nil {
				fmt.Printf("发送行级评论失败: %v\n", err)
			}
			if len(failedLineComments) > 0 {
				fmt.Printf("失败的行级评论: %v\n", failedLineComments)
			}
		}

		// 构建普通评论内容
		var generalCommentContent strings.Builder

		// 1. 添加原有的普通评论
		for _, comment := range generalComments {
			generalCommentContent.WriteString(fmt.Sprintf("- %s:%d: %s\n", comment.File, comment.Line, comment.Message))
		}

		// 2. 添加失败的行级评论
		if len(failedLineComments) > 0 {
			if generalCommentContent.Len() > 0 {
				generalCommentContent.WriteString("\n")
			}
			generalCommentContent.WriteString("以下评论因行级评论失败，转为普通评论：\n")
			for _, failedComment := range failedLineComments {
				generalCommentContent.WriteString("- ")
				generalCommentContent.WriteString(failedComment)
				generalCommentContent.WriteString("\n")
			}
		}

		// 3. 添加评论质量反馈复选框
		if generalCommentContent.Len() > 0 {
			generalCommentContent.WriteString(addQualityFeedbackCheckboxes(""))
		}

		// 发送普通评论（如果有内容）
		if generalCommentContent.Len() > 0 {
			_, err := PostCommentToGitLab(data.GitlabInfo.Config.API, body.Project.ID, body.ObjectAttributes.IID, data.GitlabToken, generalCommentContent.String())
			if err != nil {
				fmt.Printf("普通评论失败: %v\n", err)
			}
		}
	}

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
}

// NewRAGClientPool 创建RAG客户端连接池
func NewRAGClientPool(baseURL string, poolSize int) (*RAGClientPool, error) {
	pool := &RAGClientPool{
		clients: make(chan *RAGClient, poolSize),
		baseURL: baseURL,
	}

	// 预创建客户端
	for i := 0; i < poolSize; i++ {
		client := NewRAGClient(baseURL)
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
