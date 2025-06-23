package service

import (
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"fmt"
	"strconv"

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

	gitlabPrompt := data.GitlabInfo.Config.Prompt
	// 如果 gitlabPrompt 为空，使用默认的 gitlabPrompt, 2关闭自定义配置
	if gitlabPrompt == "" || data.GitlabInfo.RuleCheckStatus == 2 {
		gitlabPrompt = utils.CodeReviewPrompt
	}
	// 生成提示词
	prompt := utils.GeneratePrompt(gitlabPrompt, data.MergeRequest, data.Diff)

	// 执行AI分析
	comments, err := manager.performAIEnhancement(prompt, data)
	if err != nil {
		return "", err
	}

	// 保存结果
	aiMessage, err := manager.saveResult(body, comments, data)
	if err != nil {
		return "", err
	}

	// 发送通知 - 修复空指针问题
	var aiMessageID uint
	if aiMessage != nil {
		aiMessageID = aiMessage.ID
	}
	manager.sendNotifications(body, comments, data, aiMessageID)

	return comments, nil
}

// 只处理opened状态的merge请求
func ShouldProcessState(state string) bool {
	return state == "opened"
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
	} else {
		fmt.Println("评论提交失败:", comments, iid, projectID, token, api)
	}
}

func PushWeChatInfo(pathWithNamespace, mergeURL, comments string, aiMessageId uint) string {
	return fmt.Sprintf("项目: %s\n合并请求: %s\nAI检查结果: %s\nAI消息ID: %d", pathWithNamespace, mergeURL, comments, aiMessageId)
}

func pushWebhookIfNeeded(webhookURL string, webhookStatus int8, pathWithNamespace, mergeURL, comments string, aiMessageId uint, mergeRequest *model.MergeRequestInfo) {
	if webhookURL != "" && webhookStatus == 1 {
		webhookContent := PushWeChatInfo(pathWithNamespace, mergeURL, comments, aiMessageId)
		_ = SendMarkdownToWechatBot(webhookURL, webhookContent)
		fmt.Println("推送webhook成功", mergeRequest)
	}
}
