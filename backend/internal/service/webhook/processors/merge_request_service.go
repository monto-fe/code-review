package processors

import (
	"fmt"
	"time"

	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/utils"
)

// MergeRequestService Merge Request 服务
type MergeRequestService struct {
	// 这里应该注入实际的依赖服务
	// 暂时使用 nil，实际使用时需要注入
}

// NewMergeRequestService 创建 Merge Request 服务
func NewMergeRequestService() *MergeRequestService {
	return &MergeRequestService{}
}

// ProcessMergeRequest 处理 Merge Request 事件
func (s *MergeRequestService) ProcessMergeRequest(body dto.WebhookBody) error {
	startTime := time.Now()
	fmt.Printf("开始处理Merge Request: ProjectID=%d, MergeRequestID=%d\n",
		body.Project.ID, body.ObjectAttributes.IID)

	// 直接调用现有的 AI 检查逻辑
	result, err := CheckMergeRequestWithAI(body)
	if err != nil {
		return fmt.Errorf("AI检查失败: %v", err)
	}

	fmt.Printf("Merge Request处理完成 (耗时: %v): %s\n", time.Since(startTime), result)
	return nil
}

// ProcessMergeRequestEvent 处理 Merge Request 事件（对外接口）
func ProcessMergeRequestEvent(body dto.WebhookBody) error {
	service := NewMergeRequestService()
	return service.ProcessMergeRequest(body)
}

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
