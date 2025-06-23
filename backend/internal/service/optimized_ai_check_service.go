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

// OptimizedAICheckService 优化的AI检查服务
type OptimizedAICheckService struct {
	ragClient     *OptimizedRAGClient
	aiProvider    providers.AIProvider
	aiRuleService *AIRuleService
	db            *gorm.DB
	config        *config.Config
	metrics       *RAGMetrics
}

// CheckContext 检查上下文，包含所有必要的数据
type CheckContext struct {
	mu           sync.RWMutex
	GitlabInfo   *dto.GitlabCacheItem
	GitlabToken  string
	MergeRequest *model.MergeRequestInfo
	Diff         []model.Change
	AIConfig     *model.AIConfig
	FinalRule    string
	GitURL       string
	DiffStr      string
	RAGResult    *CodeAnalysisResponse
	AIResult     string
	Errors       []error
}

// NewOptimizedAICheckService 创建优化的AI检查服务
func NewOptimizedAICheckService() (*OptimizedAICheckService, error) {
	// 加载配置
	cfg := config.LoadConfig()
	if cfg.RAGServiceURL == "" {
		return nil, fmt.Errorf("RAG服务URL未配置")
	}

	// 创建优化的RAG客户端
	ragConfig := &RAGClientConfig{
		BaseURL:           cfg.RAGServiceURL,
		PoolSize:          10,
		Timeout:           180 * time.Second,
		MaxRetries:        3,
		RetryDelay:        1 * time.Second,
		CacheTTL:          5 * time.Minute,
		MaxCacheSize:      1000,
		EnableCompression: true,
	}
	ragClient, err := NewOptimizedRAGClient(ragConfig)
	if err != nil {
		return nil, fmt.Errorf("创建RAG客户端失败: %v", err)
	}

	// 获取AI配置
	aiConfig, ok := cache.GetAIConfigCache()
	if !ok {
		return nil, fmt.Errorf("AI配置未找到")
	}

	// 创建AI提供者
	var aiProvider providers.AIProvider
	switch aiConfig.Type {
	case "UCloud":
		aiProvider = providers.NewUCAIProvider(aiConfig)
	case "DeepSeek":
		aiProvider = providers.NewDeepSeekProvider(aiConfig)
	default:
		return nil, fmt.Errorf("不支持的AI模型类型: %s", aiConfig.Type)
	}

	// 初始化数据库服务
	db := database.GetDB()
	aiRuleService := NewAIRuleService(db)

	// 创建指标收集器
	metrics := NewRAGMetrics()

	return &OptimizedAICheckService{
		ragClient:     ragClient,
		aiProvider:    aiProvider,
		aiRuleService: aiRuleService,
		db:            db,
		config:        cfg,
		metrics:       metrics,
	}, nil
}

// CheckMergeRequestOptimized 优化的合并请求检查主函数
func CheckMergeRequestOptimized(body dto.WebhookBody) (string, error) {
	// 检查Merge Request状态
	// if !ShouldProcessState(body.ObjectAttributes.State) {
	// 	return "", fmt.Errorf("跳过非opened状态的合并请求")
	// }

	// 创建优化的AI检查服务
	service, err := NewOptimizedAICheckService()
	if err != nil {
		return "", fmt.Errorf("创建AI检查服务失败: %v", err)
	}

	// 执行优化的检查流程
	return service.ExecuteOptimizedCheck(body)
}

// ExecuteOptimizedCheck 执行优化的检查流程
func (s *OptimizedAICheckService) ExecuteOptimizedCheck(body dto.WebhookBody) (string, error) {
	// 1. 快速验证阶段
	ctx, err := s.prepareContext(body)
	if err != nil {
		return "", err
	}

	// 2. 并发数据获取阶段
	if err := s.concurrentDataFetch(ctx, body); err != nil {
		return "", err
	}

	// 3. 并发AI分析阶段
	result, err := s.concurrentAIAnalysis(ctx)
	if err != nil {
		// 回退到传统方式
		return s.fallbackToTraditional(body)
	}

	// 4. 异步处理阶段
	go s.asyncProcessing(body, result, ctx)

	return result, nil
}

// prepareContext 快速验证和准备上下文
func (s *OptimizedAICheckService) prepareContext(body dto.WebhookBody) (*CheckContext, error) {
	// 验证GitLab Token
	gitlabCache := cache.GetGitlabCache()
	gitlabToken, gitlabInfo, ok := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)
	if !ok {
		return nil, fmt.Errorf("请配置GitLab Token")
	}

	// 检查分支匹配
	if !branchMatch(gitlabInfo.Config.SourceBranch, body.ObjectAttributes.SourceBranch) {
		return nil, fmt.Errorf("源分支不匹配")
	}
	if !branchMatch(gitlabInfo.Config.TargetBranch, body.ObjectAttributes.TargetBranch) {
		return nil, fmt.Errorf("目标分支不匹配")
	}

	// 验证AI配置
	aiConfig, ok := cache.GetAIConfigCache()
	if !ok {
		return nil, fmt.Errorf("AI配置失败")
	}

	return &CheckContext{
		GitlabInfo:  &gitlabInfo,
		GitlabToken: gitlabToken,
		AIConfig:    aiConfig,
	}, nil
}

// concurrentDataFetch 并发数据获取
func (s *OptimizedAICheckService) concurrentDataFetch(ctx *CheckContext, body dto.WebhookBody) error {
	var wg sync.WaitGroup
	var errors []error
	var mu sync.Mutex

	// goroutine 1: 获取Merge Request信息
	wg.Add(1)
	go func() {
		defer wg.Done()
		mergeRequest, err := GetMergeRequestInfo(
			ctx.GitlabInfo.Config.API,
			strconv.Itoa(body.Project.ID),
			ctx.GitlabToken,
		)
		if err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("获取Merge Request失败: %v", err))
			mu.Unlock()
			return
		}
		ctx.mu.Lock()
		ctx.MergeRequest = mergeRequest
		ctx.GitURL = s.buildGitURL(mergeRequest.WebURL)
		ctx.mu.Unlock()
	}()

	// goroutine 2: 获取代码Diff
	wg.Add(1)
	go func() {
		defer wg.Done()
		diff := getMergeDiff(
			ctx.GitlabInfo.Config.API,
			body.Project.ID,
			body.ObjectAttributes.IID,
			ctx.GitlabToken,
		)
		ctx.mu.Lock()
		ctx.Diff = diff
		ctx.DiffStr = s.buildDiffString(diff)
		ctx.mu.Unlock()
	}()

	// 等待所有goroutine完成
	wg.Wait()

	// 检查关键错误
	if ctx.MergeRequest == nil {
		// MergeRequest 是后续步骤的必需数据
		return fmt.Errorf("获取Merge Request信息失败，无法继续")
	}

	// 在获取Merge Request信息之后，再获取规则
	finalRule, err := s.buildFinalRule(ctx)
	if err != nil {
		return fmt.Errorf("获取规则失败: %v", err)
	}
	ctx.mu.Lock()
	ctx.FinalRule = finalRule
	ctx.mu.Unlock()

	// 检查并发过程中的其他错误
	if len(errors) > 0 {
		// 返回第一个遇到的错误
		return errors[0]
	}

	return nil
}

// concurrentAIAnalysis 并发AI分析
func (s *OptimizedAICheckService) concurrentAIAnalysis(ctx *CheckContext) (string, error) {
	// 准备RAG请求数据
	ragReq := &CodeReviewRequest{
		GitURL:      ctx.GitURL,
		Branch:      ctx.MergeRequest.SourceBranch,
		DiffContent: ctx.DiffStr,
		Query:       ctx.FinalRule,
		GitlabToken: ctx.GitlabToken,
	}

	// 创建180秒超时上下文
	analysisCtx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var ragResult *CodeAnalysisResponse
	var aiResult string
	var ragErr, aiErr error

	// goroutine 1: 调用RAG服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		result, err := s.ragClient.AnalyzeCodeWithRequestOptimized(ragReq)
		duration := time.Since(start)

		s.metrics.RecordRequest(err == nil, duration)

		if err != nil {
			ragErr = err
			return
		}
		ragResult = result
	}()

	// goroutine 2: 调用AI服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()

		// 生成AI提示词
		prompt := s.generateAIPrompt(ctx)
		result, err := s.aiProvider.CallAI(prompt)
		duration := time.Since(start)

		s.metrics.RecordRequest(err == nil, duration)

		if err != nil {
			aiErr = err
			return
		}
		aiResult = result
	}()

	// 等待分析完成或超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 分析完成
	case <-analysisCtx.Done():
		// 超时
		return "", fmt.Errorf("AI分析超时(180秒)")
	}

	// 结果选择策略
	return s.selectResult(ragResult, aiResult, ragErr, aiErr)
}

// selectResult 结果选择策略
func (s *OptimizedAICheckService) selectResult(
	ragResult *CodeAnalysisResponse,
	aiResult string,
	ragErr, aiErr error,
) (string, error) {
	// RAG成功
	if ragErr == nil && ragResult != nil && ragResult.Review != "" {
		return ragResult.Review, nil
	}

	// RAG失败，AI成功
	if ragErr != nil && aiErr == nil && aiResult != "" {
		return aiResult, nil
	}

	// 都失败
	if ragErr != nil && aiErr != nil {
		return "", fmt.Errorf("RAG和AI服务都失败: RAG错误=%v, AI错误=%v", ragErr, aiErr)
	}

	// 其他情况
	if ragResult != nil && ragResult.Review != "" {
		return ragResult.Review, nil
	}
	if aiResult != "" {
		return aiResult, nil
	}

	return "", fmt.Errorf("未获得有效的分析结果")
}

// asyncProcessing 异步处理
func (s *OptimizedAICheckService) asyncProcessing(body dto.WebhookBody, result string, ctx *CheckContext) {
	var wg sync.WaitGroup

	// goroutine 1: 保存AI消息到数据库
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.saveAIMessage(body, result)
	}()

	// goroutine 2: 发送评论到GitLab
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.sendGitLabComment(body, result, ctx)
	}()

	// goroutine 3: 推送Webhook通知
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.pushWebhookNotification(body, result, ctx)
	}()

	// 等待所有异步操作完成
	wg.Wait()
}

// saveAIMessage 保存AI消息
func (s *OptimizedAICheckService) saveAIMessage(body dto.WebhookBody, result string) {
	passed := -1
	if strings.Contains(result, "未发现Bug") || strings.Contains(result, "通过") {
		passed = 1
	}

	aiMessage := &model.AImessage{
		ProjectID:        uint(body.Project.ID),
		MergeURL:         body.ObjectAttributes.URL,
		ProjectName:      body.Project.Name,
		ProjectNamespace: body.Project.PathWithNamespace,
		MergeDescription: body.ObjectAttributes.Description,
		MergeID:          fmt.Sprintf("%d", body.ObjectAttributes.IID),
		AIModel:          "OptimizedRAG",
		Rule:             model.RuleType(1),
		RuleID:           1,
		Result:           result,
		Passed:           passed,
		CreateTime:       time.Now().Unix(),
	}

	if err := s.db.Create(aiMessage).Error; err != nil {
		fmt.Printf("保存AI消息失败: %v\n", err)
	}
}

// sendGitLabComment 发送GitLab评论
func (s *OptimizedAICheckService) sendGitLabComment(body dto.WebhookBody, result string, ctx *CheckContext) {
	postComment(
		ctx.GitlabInfo.Config.API,
		body.Project.ID,
		body.ObjectAttributes.IID,
		ctx.GitlabToken,
		result,
	)
}

// pushWebhookNotification 推送Webhook通知
func (s *OptimizedAICheckService) pushWebhookNotification(body dto.WebhookBody, result string, ctx *CheckContext) {
	pushWebhookIfNeeded(
		ctx.GitlabInfo.Config.WebhookURL,
		ctx.GitlabInfo.Config.WebhookStatus,
		body.Project.PathWithNamespace,
		body.ObjectAttributes.MergeURL,
		result,
		0, // TODO: 获取aiMessage.ID
		ctx.MergeRequest,
	)
}

// fallbackToTraditional 回退到传统方式
func (s *OptimizedAICheckService) fallbackToTraditional(body dto.WebhookBody) (string, error) {
	fmt.Println("回退到传统AI检查方式")
	return CheckMergeRequestWithAI(body)
}

// buildGitURL 构建Git URL
func (s *OptimizedAICheckService) buildGitURL(webURL string) string {
	if strings.Contains(webURL, "/-/merge_requests/") {
		parts := strings.Split(webURL, "/-/merge_requests/")
		if len(parts) > 0 {
			return parts[0] + ".git"
		}
	}
	return webURL
}

// buildDiffString 构建差异字符串
func (s *OptimizedAICheckService) buildDiffString(diff []model.Change) string {
	var diffStr strings.Builder
	for _, change := range diff {
		diffStr.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", change.OldPath, change.NewPath))
		diffStr.WriteString(change.Diff)
		diffStr.WriteString("\n")
	}
	return diffStr.String()
}

// buildFinalRule 构建最终规则
func (s *OptimizedAICheckService) buildFinalRule(ctx *CheckContext) (string, error) {
	// 获取项目规则
	customRule, err := s.aiRuleService.GetCustomRuleByProjectID(ctx.MergeRequest.ProjectID)
	if err != nil {
		return "", fmt.Errorf("获取自定义规则失败: %v", err)
	}

	var currentRule string
	if customRule != nil && customRule.Rule != "" {
		currentRule = customRule.Rule
	} else {
		// 获取项目主要语言
		language, err := GetDominantLanguage(
			ctx.GitlabInfo.Config.API,
			fmt.Sprintf("%d", ctx.MergeRequest.ProjectID),
			ctx.GitlabToken,
		)
		if err != nil {
			return "", fmt.Errorf("获取项目语言失败: %v", err)
		}

		// 获取通用规则
		commonRule, err := s.aiRuleService.GetCommonRule(language)
		if err != nil {
			return "", fmt.Errorf("获取通用规则失败: %v", err)
		}
		if len(commonRule) > 0 {
			currentRule = commonRule[0].Rule
		}
	}

	// 获取GitLab配置中的Prompt
	gitlabPrompt := ctx.GitlabInfo.Config.Prompt
	if gitlabPrompt == "" {
		gitlabPrompt = "请检查代码是否符合以下要求：\n1. 代码风格是否规范\n2. 是否有潜在的安全问题\n3. 是否有性能优化空间\n4. 是否有重复代码\n5. 是否有未使用的代码\n6. 是否有更好的实现方式"
	}

	// 拼接规则
	if currentRule != "" {
		return fmt.Sprintf("%s\n%s", currentRule, gitlabPrompt), nil
	}
	return gitlabPrompt, nil
}

// generateAIPrompt 生成AI提示词
func (s *OptimizedAICheckService) generateAIPrompt(ctx *CheckContext) string {
	var diffContent strings.Builder
	for _, change := range ctx.Diff {
		diffContent.WriteString(fmt.Sprintf("File: %s\n%s\n\n", change.NewPath, change.Diff))
	}

	return fmt.Sprintf(`请检查以下代码差异（diff），确保其符合以下要求：
规则：%s

代码信息：
标题：%s
描述：%s

代码差异：
%s

请使用中文回答，如果没有BUG输出：未发现Bug`,
		ctx.FinalRule,
		ctx.MergeRequest.Title,
		ctx.MergeRequest.Description,
		diffContent.String())
}

// GetMetrics 获取性能指标
func (s *OptimizedAICheckService) GetMetrics() map[string]interface{} {
	return s.metrics.GetStats()
}
