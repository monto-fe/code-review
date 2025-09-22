package service

import (
	dto "code-review-go/internal/dto"
	"code-review-go/internal/pkg/utils"
	"fmt"
)

// CheckMergeRequestWithAI 使用指定的管理器执行RAG检查
func CheckMergeRequestWithAI(body dto.WebhookBody) (string, error) {
	// 获取主服务管理器
	manager := GetMainServiceManager()
	if err := manager.Initialize(); err != nil {
		fmt.Printf("服务初始化失败: %v\n", err)
		return "", err
	}

	// 前置验证
	if err := manager.ValidateRequest(body); err != nil {
		sendFailureNotification(manager, body, fmt.Sprintf("请求验证失败: %v", err))
		return "", err
	}

	// 获取必要数据
	data, err := manager.PrepareData(body)
	if err != nil {
		sendFailureNotification(manager, body, fmt.Sprintf("数据准备失败: %v", err))
		return "", err
	}

	// 1. 执行RAG分析
	prompt := ""
	ragResult, err := manager.PerformRAGAnalysis(data)
	if err == nil {
		prompt = manager.GenerateEnhancedPrompt(ragResult, data)
		fmt.Printf("RAG分析成功，提示词: %s\n", prompt)
		// 2. 如果RAG分析失败，则使用AI检查
	} else {
		gitlabPrompt := data.GitlabInfo.Config.Prompt
		// 如果 gitlabPrompt 为空，使用默认的 gitlabPrompt, 2关闭自定义配置
		if gitlabPrompt == "" || data.GitlabInfo.RuleCheckStatus == 2 {
			gitlabPrompt = utils.CodeReviewPrompt
		}
		// 生成提示词
		// 读取gitlab中的评论类型，然后选择不同的提示词
		if data.GitlabInfo.CommentType == utils.CommentTypeCommon {
			prompt = utils.GenerateAICheckCommonPrompt(gitlabPrompt, data.MergeRequest.Title, data.MergeRequest.Description, data.DiffStr)
		} else {
			prompt = utils.GenerateAICheckInlinePrompt(gitlabPrompt, data.MergeRequest.Title, data.MergeRequest.Description, data.DiffStr)
		}
		fmt.Printf("AI CodeReview分析的提示词: %s\n", prompt)
	}
	// AI检查
	result, err := manager.PerformAIEnhancement(prompt, data)
	if err != nil {
		sendFailureNotification(manager, body, fmt.Sprintf("AI CodeReview分析失败: %v", err))
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
