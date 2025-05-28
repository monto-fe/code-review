package service

import (
	"code-review-go/internal/cache"
	"code-review-go/internal/database"
	"code-review-go/internal/model"
	"code-review-go/internal/service/providers"
	"fmt"
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

// CheckMergeRequestWithAI 使用aiConfig中status=1的ai模型检查合并请求
func CheckMergeRequestWithAI(mergeRequest *model.MergeRequestInfo, diff []model.Change, gitlabAPI, gitlabToken string) (string, uint, error) {
	// 获取AI配置
	aiConfig, ok := cache.GetAIConfigCache()
	if !ok {
		return "", 0, fmt.Errorf("初始化AI配置管理器失败: %v", ok)
	}

	if aiConfig.APIURL == "" {
		return "apiUrl未配置", 0, nil
	}

	// 初始化服务
	db := database.GetDB()
	aiRuleService := NewAIRuleService(db)

	// 获取项目规则
	customRule, err := aiRuleService.GetCustomRuleByProjectID(mergeRequest.ProjectID)
	if err != nil {
		return "", 0, fmt.Errorf("获取自定义规则失败: %v", err)
	}

	var currentRule string
	if customRule != nil && customRule.Rule != "" {
		currentRule = customRule.Rule
	} else {
		// 获取项目主要语言
		language, err := GetDominantLanguage(gitlabAPI, fmt.Sprintf("%d", mergeRequest.ProjectID), gitlabToken)
		if err != nil {
			return "", 0, fmt.Errorf("获取项目语言失败: %v", err)
		}

		// 获取通用规则
		commonRule, err := aiRuleService.GetCommonRule(language)
		if err != nil {
			return "", 0, fmt.Errorf("获取通用规则失败: %v", err)
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

	// 生成提示词
	prompt := generatePrompt(finalRule, mergeRequest, diff)

	fmt.Println("prompt", prompt)

	// 根据aiConfig.Type选择AIProvider
	var provider providers.AIProvider
	switch aiConfig.Type {
	case "ucai":
		provider = providers.NewUCAIProvider(aiConfig)
	case "deepseek":
		provider = providers.NewDeepSeekProvider(aiConfig)
	default:
		return "", 0, fmt.Errorf("不支持的AI模型类型: %s", aiConfig.Type)
	}

	comments, err := provider.CallAI(prompt)
	if err != nil {
		return "", 0, fmt.Errorf("调用AI接口失败: %v", err)
	}
	fmt.Println("AI检查结果:", comments)

	// 保存AI消息
	aiMessage := &model.AImessage{
		ProjectID:  mergeRequest.ProjectID,
		MergeURL:   mergeRequest.WebURL,
		MergeID:    fmt.Sprintf("%d", mergeRequest.IID),
		AIModel:    aiConfig.Model,
		Rule:       model.RuleType(1), // 默认规则类型
		RuleID:     1,                 // 默认规则ID
		Result:     comments,
		CreateTime: time.Now().Unix(),
	}

	if err := db.Create(aiMessage).Error; err != nil {
		return "", 0, fmt.Errorf("保存AI消息失败: %v", err)
	}

	return comments, aiMessage.ID, nil
}

// generatePrompt 生成AI提示词
func generatePrompt(rule string, mergeRequest *model.MergeRequestInfo, diff []model.Change) string {
	var diffContent string
	for _, change := range diff {
		diffContent += fmt.Sprintf("File: %s\n%s\n\n", change.NewPath, change.Diff)
	}

	return fmt.Sprintf(`请根据以下规则审查代码：
规则：%s

合并请求信息：
标题：%s
描述：%s

代码变更：
%s`, rule, mergeRequest.Title, mergeRequest.Description, diffContent)
}

func PushWeChatInfo(pathWithNamespace, mergeURL, comments string, aiMessageId uint) string {
	return fmt.Sprintf("项目: %s\n合并请求: %s\nAI检查结果: %s\nAI消息ID: %d", pathWithNamespace, mergeURL, comments, aiMessageId)
}
