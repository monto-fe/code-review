package helpers

import (
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/service/ai"
	"code-review-go/internal/service/gitlab_service"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type AICheckService struct {
	db            *gorm.DB
	gitlabService *gitlab_service.GitlabService
	aiRuleService *ai.AIRuleService
}

func NewAICheckService(db *gorm.DB, gitlabService *gitlab_service.GitlabService, aiRuleService *ai.AIRuleService) *AICheckService {
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
		diff, _ := gitlab_service.GetMergeRequestDiff(api, strconv.Itoa(projectID), strconv.Itoa(iid), token)
		return diff
	}
	return nil
}

func PushWeChatInfo(pathWithNamespace, mergeURL, comments string, aiMessageId uint) string {
	return fmt.Sprintf("项目: %s\n合并请求: %s\nAI检查结果: %s\nAI消息ID: %d", pathWithNamespace, mergeURL, comments, aiMessageId)
}

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
	lineComments := gitlab_service.ParseCommentsForLineComments(testComments, testDiff)

	fmt.Println("解析到的行级评论:")
	for i, comment := range lineComments {
		fmt.Printf("%d. 文件: %s, 行: %d, 级别: %s, 消息: %s\n",
			i+1, comment.File, comment.Line, comment.Severity, comment.Message)
	}
}
