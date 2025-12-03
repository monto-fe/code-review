package ai_check

import (
	"fmt"
	"strings"
	"time"

	"code-review-go/internal/cache"
	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/constants"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service/webhook/handlers"

	"github.com/gin-gonic/gin"
)

// 全局事件路由器实例
var eventRouter *handlers.EventRouter

// init 初始化事件路由器
func init() {
	eventRouter = handlers.NewEventRouter()

	// 注册事件处理器
	eventRouter.RegisterHandler(handlers.NewMergeRequestHandler())
	eventRouter.RegisterHandler(handlers.NewPushHandler())

	fmt.Printf("事件路由器初始化完成，已注册处理器: %v\n", eventRouter.GetRegisteredHandlers())
}

// AICheck 处理 AI 检查请求（事件驱动架构）
// @Summary 触发 AI 代码审查
// @Description 处理 GitLab 合并请求和 Push 事件的 webhook，自动触发 AI 代码审查与评论
// @Tags Webhook
// @Accept json
// @Produce json
// @Param data body dto.WebhookBody true "Webhook 触发参数"
// @Success 200 {object} response.Response
// @Router /v1/webhook/merge [post]
func AICheck(c *gin.Context) {
	// 1. 解析请求体
	var body dto.WebhookBody
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Printf("解析请求体失败: %v\n", err)
		c.JSON(400, gin.H{
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 2. 基础验证
	if body.Project.ID == 0 {
		fmt.Printf("无效的项目ID: %d\n", body.Project.ID)
		c.JSON(400, gin.H{
			"msg":   "无效的项目ID",
			"error": "Project ID is required",
		})
		return
	}

	// 3. 检查项目是否在白名单中
	if cache.IsProjectWhitelisted(uint(body.Project.ID)) {
		fmt.Printf("项目在白名单中，跳过代码审查: ProjectID=%d, ProjectName=%s\n",
			body.Project.ID, body.Project.Name)
		response.Success(c, gin.H{
			"eventType": "whitelisted",
			"projectId": body.Project.ID,
			"timestamp": time.Now().Unix(),
			"skipped":   true,
			"reason":    "项目在白名单中，已跳过代码审查",
		}, "项目在白名单中，已跳过代码审查", int(constants.RetCodeSuccess))
		return
	}

	// 4. 检测事件类型
	eventType := detectEventType(body)

	// 5. 检查是否为支持的事件类型
	if eventType == "unknown" {
		fmt.Printf("不支持的事件类型，跳过处理: ProjectID=%d, ObjectKind=%s\n",
			body.Project.ID, body.ObjectKind)
		return
	}

	// 6. 立即响应（不阻塞webhook）
	responseData := gin.H{
		"eventType":  eventType,
		"timestamp":  time.Now().Unix(),
		"optimized":  true,
		"processing": true,
	}

	switch eventType {
	case "merge_request":
		responseData["projectId"] = body.Project.ID
		responseData["mergeRequestId"] = body.ObjectAttributes.IID
		responseData["targetBranch"] = body.ObjectAttributes.TargetBranch
		fmt.Printf("Merge Request事件: ProjectID=%d, MergeRequestID=%d, TargetBranch=%s\n",
			body.Project.ID, body.ObjectAttributes.IID, body.ObjectAttributes.TargetBranch)
	case "push":
		responseData["projectId"] = body.Project.ID
		responseData["branch"] = extractBranchFromRef(body.Ref)
		responseData["commits"] = body.TotalCommitsCount
		fmt.Printf("Push事件: ProjectID=%d, Branch=%s, Commits=%d\n",
			body.Project.ID, extractBranchFromRef(body.Ref), body.TotalCommitsCount)
	}

	response.Success(c, responseData, "AI检查已启动，请稍候查看结果", int(constants.RetCodeSuccess))

	// 7. 使用事件路由器处理事件
	err := eventRouter.Route(body)
	if err != nil {
		fmt.Printf("事件路由失败: %v\n", err)
	}
}

// detectEventType 检测事件类型
func detectEventType(body dto.WebhookBody) string {
	// 1. 检查 object_kind 字段
	if body.ObjectKind == "merge_request" {
		return "merge_request"
	}
	if body.ObjectKind == "push" {
		return "push"
	}

	// 2. 检查是否有 Merge Request 相关字段
	if body.ObjectAttributes.IID > 0 && body.ObjectAttributes.Title != "" {
		return "merge_request"
	}

	// 3. 检查是否有 Push 相关字段
	if body.Ref != "" && body.After != "" {
		return "push"
	}

	// 4. 其他事件类型暂不处理，返回 unknown
	return "unknown"
}

// extractBranchFromRef 从 ref 中提取分支名称
func extractBranchFromRef(ref string) string {
	// 从 refs/heads/feature-branch 提取 feature-branch
	if strings.HasPrefix(ref, "refs/heads/") {
		return strings.TrimPrefix(ref, "refs/heads/")
	}
	return ref
}
