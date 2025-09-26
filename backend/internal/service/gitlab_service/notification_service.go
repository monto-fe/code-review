package gitlab_service

import (
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service/ai"
	"code-review-go/internal/service/common"
	"fmt"
	"strings"
	"sync"
)

// NotificationServiceImpl 通知服务实现
type NotificationServiceImpl struct {
	mu          sync.RWMutex
	initialized bool
}

// NewNotificationService 创建通知服务实例
func NewNotificationService() common.NotificationService {
	return &NotificationServiceImpl{}
}

// Initialize 初始化通知服务
func (n *NotificationServiceImpl) Initialize() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.initialized {
		return nil
	}

	// 通知服务不需要特殊初始化
	n.initialized = true
	return nil
}

// SendGitLabComment 发送GitLab评论
func (n *NotificationServiceImpl) SendGitLabComment(api, token string, projectID, mergeRequestID int, comment string) error {
	_, err := PostCommentToGitLab(api, projectID, mergeRequestID, token, comment)
	return err
}

// SendLineComments 发送行级评论
func (n *NotificationServiceImpl) SendLineComments(api, token string, projectID, mergeRequestID int, comments []dto.LineComment, diff []model.Change) ([]string, error) {
	return PostLineComments(api, projectID, mergeRequestID, token, comments, diff)
}

// SendWebhookNotification 发送Webhook通知
func (n *NotificationServiceImpl) SendWebhookNotification(webhookURL string, status int, projectNamespace, mergeURL, comments string, aiMessageID uint, mergeRequest *model.MergeRequestInfo) error {
	pushWebhookIfNeeded(webhookURL, int8(status), projectNamespace, mergeURL, comments, aiMessageID, mergeRequest)
	return nil
}

// pushWebhookIfNeeded 推送webhook通知
func pushWebhookIfNeeded(webhookURL string, webhookStatus int8, pathWithNamespace, mergeURL, comments string, aiMessageId uint, mergeRequest *model.MergeRequestInfo) {
	if webhookURL != "" && webhookStatus == 1 {
		webhookContent := pushWeChatInfo(pathWithNamespace, mergeURL, comments, aiMessageId)
		_ = ai.SendMarkdownToWechatBot(webhookURL, webhookContent)
		fmt.Println("推送webhook成功", mergeRequest)
	}
}

// pushWeChatInfo 构建企业微信推送信息
func pushWeChatInfo(pathWithNamespace, mergeURL, comments string, aiMessageId uint) string {
	return fmt.Sprintf("项目: %s\n合并请求: %s\nAI检查结果: %s\nAI消息ID: %d", pathWithNamespace, mergeURL, comments, aiMessageId)
}

// IsInitialized 检查是否已初始化
func (n *NotificationServiceImpl) IsInitialized() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.initialized
}

// addQualityFeedbackCheckboxes 添加评论质量反馈复选框
func addQualityFeedbackCheckboxes(comments string) string {
	return comments + "\n\n---\n**请评价此评论的质量：**\n" +
		"- [ ] 评论内容准确且有用\n" +
		"- [ ] 精准定位并提供建议\n" +
		"- [ ] 建议不够具体\n" +
		"- [ ] 完全误导性建议\n"
}

// SendEnhancedNotification 发送增强通知（包含质量反馈）
func (n *NotificationServiceImpl) SendEnhancedNotification(api, token string, projectID, mergeRequestID int, comments string, commentType int8, diff []model.Change) error {
	if commentType == utils.CommentTypeCommon {
		// 为普通评论添加质量反馈复选框
		commentWithFeedback := addQualityFeedbackCheckboxes(comments)
		return n.SendGitLabComment(api, token, projectID, mergeRequestID, commentWithFeedback)
	} else {
		// 解析所有评论内容
		allComments := ParseCommentsForLineComments(comments, diff)

		// 分类评论
		lineComments, generalComments := ClassifyComments(allComments)

		// 处理行级评论
		var failedLineComments []string
		if len(lineComments) > 0 {
			// 转换 CommentInfo 到 dto.LineComment
			dtoLineComments := make([]dto.LineComment, len(lineComments))
			for i, comment := range lineComments {
				dtoLineComments[i] = dto.LineComment{
					File:     comment.File,
					Line:     comment.Line,
					Message:  comment.Message,
					Severity: comment.Severity,
				}
			}
			var err error
			failedLineComments, err = n.SendLineComments(api, token, projectID, mergeRequestID, dtoLineComments, diff)
			if err != nil {
				fmt.Printf("发送行级评论失败: %v\n", err)
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
			return n.SendGitLabComment(api, token, projectID, mergeRequestID, generalCommentContent.String())
		}
	}

	return nil
}
