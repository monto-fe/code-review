package service

import (
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
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

// 只处理opened状态的merge请求
func ShouldProcessState(body dto.WebhookBody) bool {
	if body.ObjectKind == "merge_request" {
		return body.ObjectAttributes.State == "opened"
	}
	if body.ObjectKind == "note" && body.MergeRequest.State == "opened" {
		// note 事件通常没有 state 字段，如有需要可补充获取 MR 状态的逻辑
		return true
	}
	return false
}

func BranchMatch(cfg, actual string) bool {
	if cfg == "" {
		return true
	}
	return cfg == actual
}

func GetMergeDiff(api string, projectID, iid int, token string) []model.Change {
	if api != "" && projectID != 0 && iid != 0 {
		diff, _ := GetMergeRequestDiff(api, strconv.Itoa(projectID), strconv.Itoa(iid), token)
		return diff
	}
	return nil
}

// func postComment(api string, projectID, iid int, token, comments string) {
// 	if comments != "" && iid != 0 && projectID != 0 && token != "" && api != "" {
// 		// 获取差异信息用于行级评论
// 		diff := GetMergeDiff(api, projectID, iid, token)

// 		// 解析所有评论内容
// 		allComments := ParseCommentsForLineComments(comments, diff)

// 		// 分类评论
// 		lineComments, generalComments := ClassifyComments(allComments)

// 		// 处理行级评论
// 		var failedLineComments []string
// 		if len(lineComments) > 0 {
// 			// 转换 CommentInfo 到 dto.LineComment
// 			dtoLineComments := make([]dto.LineComment, len(lineComments))
// 			for i, comment := range lineComments {
// 				dtoLineComments[i] = dto.LineComment{
// 					File:     comment.File,
// 					Line:     comment.Line,
// 					Message:  comment.Message,
// 					Severity: comment.Severity,
// 				}
// 			}
// 			failedLineComments, err := PostLineComments(api, projectID, iid, token, dtoLineComments, diff)
// 			if err != nil {
// 				fmt.Printf("发送行级评论失败: %v\n", err)
// 			}
// 			if len(failedLineComments) > 0 {
// 				fmt.Printf("失败的行级评论: %v\n", failedLineComments)
// 			}
// 		}

// 		// 构建普通评论内容
// 		var generalCommentContent strings.Builder

// 		// 1. 添加原有的普通评论
// 		for _, comment := range generalComments {
// 			generalCommentContent.WriteString(fmt.Sprintf("- %s:%d: %s\n", comment.File, comment.Line, comment.Message))
// 		}

// 		// 2. 添加失败的行级评论
// 		if len(failedLineComments) > 0 {
// 			if generalCommentContent.Len() > 0 {
// 				generalCommentContent.WriteString("\n")
// 			}
// 			generalCommentContent.WriteString("以下评论因行级评论失败，转为普通评论：\n")
// 			for _, failedComment := range failedLineComments {
// 				generalCommentContent.WriteString("- ")
// 				generalCommentContent.WriteString(failedComment)
// 				generalCommentContent.WriteString("\n")
// 			}
// 		}

// 		// 3. 添加评论质量反馈复选框
// 		if generalCommentContent.Len() > 0 {
// 			generalCommentContent.WriteString("\n---\n")
// 			generalCommentContent.WriteString("**请评价此评论的质量：**\n")
// 			generalCommentContent.WriteString("- [ ] 精准定位并提供建议\n")
// 			generalCommentContent.WriteString("- [ ] 部分有效但需要改进\n")
// 			generalCommentContent.WriteString("- [ ] 完全误导性建议\n")
// 			generalCommentContent.WriteString("- [ ] 建议不相关\n")
// 		}

// 		// 发送普通评论（如果有内容）
// 		if generalCommentContent.Len() > 0 {
// 			_, err := PostCommentToGitLab(api, projectID, iid, token, generalCommentContent.String())
// 			if err != nil {
// 				fmt.Println("普通评论失败:", err)
// 			}
// 		}
// 	} else {
// 		fmt.Println("评论提交失败:", comments, iid, projectID, token, api)
// 	}
// }

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

// 实现gitlab推送评论到对应的行

// TestLineCommentParsing 测试行级评论解析功能
func TestLineCommentParsing() {
	// 模拟AI输出
	testComments := `代码检查结果：

src/main.go:15: 这里有一个潜在的空指针问题
src/utils.go:42: 建议添加错误处理
文件: src/config.go, 行: 28, 问题: 变量命名不符合规范
src/handler.go:67: 性能警告：这里可能存在内存泄漏

总体评价：代码质量良好，但需要注意上述问题。`

	// 模拟差异数据
	testDiff := []model.Change{
		{NewPath: "src/main.go", OldPath: "src/main.go"},
		{NewPath: "src/utils.go", OldPath: "src/utils.go"},
		{NewPath: "src/config.go", OldPath: "src/config.go"},
		{NewPath: "src/handler.go", OldPath: "src/handler.go"},
	}

	// 解析行级评论
	lineComments := ParseCommentsForLineComments(testComments, testDiff)

	fmt.Println("解析到的行级评论:")
	for i, comment := range lineComments {
		fmt.Printf("%d. 文件: %s, 行: %d, 级别: %s, 消息: %s\n",
			i+1, comment.File, comment.Line, comment.Severity, comment.Message)
	}
}
