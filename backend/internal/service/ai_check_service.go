package service

import (
	"code-review-go/internal/cache"
	"code-review-go/internal/database"
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/service/providers"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AICheckService struct {
	db            *gorm.DB
	gitlabService *GitlabService
	aiRuleService *AIRuleService
}

func NewAICheckService(db *gorm.DB, gitlabService *GitlabService, aiRuleService *AIRuleService) *AICheckService {
	return &AICheckService{
		db:            db,
		gitlabService: gitlabService,
		aiRuleService: aiRuleService,
	}
}

// 使用大模型检查合并请求
func CheckMergeRequestWithAI(body dto.WebhookBody) (string, error) {
	// 获取 Gitlab 配置缓存
	gitlabCache := cache.GetGitlabCache()
	// 匹配对应的gitlabToken
	gitlabToken, gitlabInfoResult, ok := cache.FindTokenByProjectID(fmt.Sprintf("%d", body.Project.ID), gitlabCache)
	if !ok {
		return "", fmt.Errorf("请配置gitlab Token")
	}
	// TODO: 如果源分支或者目标分支为空，则需要记录这一条记录，并且需要通知用户更换token
	if !branchMatch(gitlabInfoResult.Config.SourceBranch, body.ObjectAttributes.SourceBranch) {
		return "", fmt.Errorf("源分支不匹配")
	}
	// TODO: 如果目标分支为空，则需要记录这一条记录，并且需要通知用户更换token
	if !branchMatch(gitlabInfoResult.Config.TargetBranch, body.ObjectAttributes.TargetBranch) {
		return "", fmt.Errorf("目标分支不匹配")
	}
	// 获取对应的mergeRequest
	mergeRequest, err := GetMergeRequestInfo(gitlabInfoResult.Config.API, strconv.Itoa(body.Project.ID), gitlabToken)
	if err != nil {
		return "", err
	}

	// 获取对应的diff
	diff := getMergeDiff(gitlabInfoResult.Config.API, body.Project.ID, body.ObjectAttributes.IID, gitlabToken)

	// 获取AI Config配置
	aiConfig, ok := cache.GetAIConfigCache()
	if !ok {
		return "", fmt.Errorf("初始化AI配置管理器失败: %v", ok)
	}
	if aiConfig.APIURL == "" {
		return "", fmt.Errorf("apiUrl未配置")
	}

	gitlabPrompt := gitlabInfoResult.Config.Prompt
	// 如果 gitlabPrompt 为空，使用默认的 gitlabPrompt, 2关闭自定义配置
	if gitlabPrompt == "" || gitlabInfoResult.RuleCheckStatus == 2 {
		gitlabPrompt = `输出要求：
    	仅输出以下两部分内容，如无内容可省略：
        1. 不符合检查原则的地方（使用 Markdown 表格格式）：
          - 不符合的代码文件和行号 | 违反的原则 | 修改建议
        2. 疑似 Bug 的地方（基于描述信息和逻辑判断分析，如无可省略）。`
	}
	// 生成提示词
	prompt := generatePrompt(gitlabPrompt, mergeRequest, diff)

	fmt.Println("prompt", prompt)

	// 根据aiConfig.Type选择AIProvider
	var provider providers.AIProvider
	switch aiConfig.Type {
	case "UCloud":
		provider = providers.NewUCAIProvider(aiConfig)
	case "DeepSeek":
		provider = providers.NewDeepSeekProvider(aiConfig)
	default:
		return "", fmt.Errorf("不支持的AI模型类型: %s", aiConfig.Type)
	}

	comments, err := provider.CallAI(prompt)
	if err != nil {
		return "", fmt.Errorf("调用AI接口失败: %v", err)
	}

	postComment(gitlabInfoResult.Config.API, body.Project.ID, body.ObjectAttributes.IID, gitlabToken, comments)

	// 保存AI消息
	passed := -1
	if strings.Contains(comments, "未发现Bug") {
		passed = 1
	}

	// 保存AI消息
	aiMessage := &model.AImessage{
		ProjectID:        uint(body.Project.ID),
		MergeURL:         body.ObjectAttributes.URL,
		ProjectName:      body.Project.Name,
		ProjectNamespace: body.Project.PathWithNamespace,
		MergeDescription: body.ObjectAttributes.Description,
		MergeID:          fmt.Sprintf("%d", body.ObjectAttributes.IID),
		AIModel:          aiConfig.Model,
		Rule:             model.RuleType(1), // 默认规则类型
		RuleID:           1,                 // 默认规则ID
		Result:           comments,
		Passed:           passed,
		CreateTime:       time.Now().Unix(),
	}

	// 初始化数据库服务
	db := database.GetDB()

	if err := db.Create(aiMessage).Error; err != nil {
		return "", fmt.Errorf("保存AI消息失败: %v", err)
	}

	pushWebhookIfNeeded(gitlabInfoResult.Config.WebhookURL, gitlabInfoResult.Config.WebhookStatus, body.Project.PathWithNamespace, body.ObjectAttributes.MergeURL, comments, aiMessage.ID, mergeRequest)

	return comments, nil
}

// 只处理opened状态的merge请求
func ShouldProcessState(state string) bool {
	return state == "opened"
}

// CheckMergeRequestWithAI 使用aiConfig中status=1的ai模型检查合并请求
func CheckMergeRequestWithRAG(mergeRequest *model.MergeRequestInfo, diff []model.Change, gitlabAPI, gitlabToken string) (string, error) {
	// 获取AI配置
	aiConfig, ok := cache.GetAIConfigCache()
	if !ok {
		return "", fmt.Errorf("初始化AI配置管理器失败: %v", ok)
	}

	if aiConfig.APIURL == "" {
		return "apiUrl未配置", nil
	}

	// 初始化服务
	db := database.GetDB()
	aiRuleService := NewAIRuleService(db)

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
		language, err := GetDominantLanguage(gitlabAPI, fmt.Sprintf("%d", mergeRequest.ProjectID), gitlabToken)
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

	// 获取 Gitlab 配置中的 Prompt
	gitlabCache := cache.GetGitlabCache()
	_, gitlabInfo, _ := cache.FindTokenByProjectID(fmt.Sprintf("%d", mergeRequest.ProjectID), gitlabCache)
	gitlabPrompt := gitlabInfo.Config.Prompt

	// 如果 gitlabPrompt 为空，使用默认的 gitlabPrompt
	if gitlabPrompt == "" {
		gitlabPrompt = "请检查代码是否符合以下要求：\n1. 代码风格是否规范\n2. 是否有潜在的安全问题\n3. 是否有性能优化空间\n4. 是否有重复代码\n5. 是否有未使用的代码\n6. 是否有更好的实现方式"
	}

	// 拼接规则
	var finalRule string
	if currentRule != "" {
		finalRule = fmt.Sprintf("%s\n%s", currentRule, gitlabPrompt)
	} else {
		finalRule = gitlabPrompt
	}

	// 创建RAG客户端
	ragClient := NewRAGClient("http://localhost:8000")

	// 将diff转换为字符串
	diffStr := ""
	for _, change := range diff {
		diffStr += fmt.Sprintf("diff --git a/%s b/%s\n", change.OldPath, change.NewPath)
		diffStr += change.Diff + "\n"
	}

	// 准备请求
	req := &CodeReviewRequest{
		GitURL:      mergeRequest.WebURL, // 使用WebURL作为git_url
		Branch:      mergeRequest.SourceBranch,
		DiffContent: diffStr,
		Query:       finalRule, // 直接使用规则作为查询
		GitlabToken: gitlabInfo.Token,
	}

	// 调用RAG服务获取代码分析结果
	_, err = ragClient.AnalyzeCodeWithRequest(req)
	if err != nil {
		return "", fmt.Errorf("调用RAG服务失败: %v", err)
	}

	// 使用大模型生成审查结果
	comments, err := ragClient.GenerateReview(req)
	if err != nil {
		return "", fmt.Errorf("生成审查结果失败: %v", err)
	}

	fmt.Println("代码审查结果:", comments)

	// 保存AI消息
	// passed := -1
	// if strings.Contains(comments, "未发现Bug") {
	// 	passed = 1
	// }

	// aiMessage := &model.AImessage{
	// 	ProjectID:  mergeRequest.ProjectID,
	// 	MergeURL:   mergeRequest.WebURL,
	// 	MergeID:    fmt.Sprintf("%d", mergeRequest.IID),
	// 	AIModel:    aiConfig.Model,    // 使用配置的模型
	// 	Rule:       model.RuleType(1), // 默认规则类型
	// 	RuleID:     1,                 // 默认规则ID
	// 	Result:     comments,
	// 	Passed:     passed,
	// 	CreateTime: time.Now().Unix(),
	// }

	// if err := db.Create(aiMessage).Error; err != nil {
	// 	return "", 0, fmt.Errorf("保存AI消息失败: %v", err)
	// }

	return comments, nil
}

// generatePrompt 生成AI提示词
func generatePrompt(rule string, mergeRequest *model.MergeRequestInfo, diff []model.Change) string {
	var diffContent string
	for _, change := range diff {
		diffContent += fmt.Sprintf("File: %s\n%s\n\n", change.NewPath, change.Diff)
	}

	return fmt.Sprintf(`请检查以下代码差异（diff），确保其符合以要求：
规则：%s

代码信息：
标题：%s
描述：%s

代码差异：
%s，请使用中文回答,如果没有BUG输出：未发现Bug`, rule, mergeRequest.Title, mergeRequest.Description, diffContent)
}

func PushWeChatInfo(pathWithNamespace, mergeURL, comments string, aiMessageId uint) string {
	return fmt.Sprintf("项目: %s\n合并请求: %s\nAI检查结果: %s\nAI消息ID: %d", pathWithNamespace, mergeURL, comments, aiMessageId)
}

func branchMatch(cfg, actual string) bool {
	if cfg == "" {
		return true
	}
	return cfg == actual
}

func getMergeDiff(api string, projectID, iid int, token string) []model.Change {
	if api != "" && projectID != 0 && iid != 0 {
		diff, _ := GetMergeRequestDiff(api, strconv.Itoa(projectID), strconv.Itoa(iid), token)
		return diff
	}
	return nil
}

func postComment(api string, projectID, iid int, token, comments string) {
	if comments != "" && iid != 0 && projectID != 0 && token != "" && api != "" {
		_, err := PostCommentToGitLab(api, projectID, iid, token, comments)
		if err != nil {
			fmt.Println("评论提交失败:", err)
		}
	}
}

func pushWebhookIfNeeded(webhookURL string, webhookStatus int8, pathWithNamespace, mergeURL, comments string, aiMessageId uint, mergeRequest *model.MergeRequestInfo) {
	if webhookURL != "" && webhookStatus == 1 {
		webhookContent := PushWeChatInfo(pathWithNamespace, mergeURL, comments, aiMessageId)
		_ = SendMarkdownToWechatBot(webhookURL, webhookContent)
		fmt.Println("推送webhook成功", mergeRequest)
	}
}
