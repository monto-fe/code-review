package ai_check

import (
	"fmt"
	"time"

	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service"

	"github.com/gin-gonic/gin"
)

// AICheck 处理 AI 检查请求
// @Summary 触发 AI 代码审查
// @Description 处理 GitLab 合并请求的 webhook，自动触发 AI 代码审查与评论
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
		c.JSON(400, gin.H{
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("RAG检查请求: ProjectID=%d, MergeRequestID=%d\n",
		body.Project.ID, body.ObjectAttributes.IID)

	// 💡 防止重复执行
	if utils.IsDuplicateWebhook(body) {
		fmt.Println("重复 webhook 请求，跳过处理")
		return
	}

	// 2. 立即响应（不阻塞webhook）
	response.Success(c, gin.H{
		"projectId":      body.Project.ID,
		"mergeRequestId": body.ObjectAttributes.IID,
		"optimized":      true,
		"timestamp":      time.Now().Unix(),
	}, "AI检查已启动，请稍候查看结果", 0)

	// 3. 异步处理优化的RAG检查
	go handleOptimizedAICheck(body)
}

// handleOptimizedAICheck 处理优化的AI检查
func handleOptimizedAICheck(body dto.WebhookBody) {
	startTime := time.Now()
	fmt.Printf("开始AI检查流程: %s\n", startTime.Format("2006-01-02 15:04:05"))

	// 检查Merge Request状态
	if !service.ShouldProcessState(body) {
		fmt.Printf("跳过非opened状态的合并请求: %+v\n", body)
		return
	}

	// 使用优化的RAG检查服务
	result, data, aiMessageID, err := service.CheckMergeRequestWithRAGOptimized(body)
	fmt.Printf("RAG检查结果返回: %s\n", result)

	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("RAG检查失败 (耗时: %v): %v\n", duration, err)

		// 最后回退到AI检查
		aiResult, aiErr := service.CheckMergeRequestWithAI(body)
		if aiErr != nil {
			if data != nil {
				manager := service.GetRAGServiceManager()
				manager.SendNotifications(body, fmt.Sprintf("AI审核检查失败: %v", aiErr), data, aiMessageID)
			}
			fmt.Printf("所有检查方式都失败: %v\n", aiErr)
			return
		}
		fmt.Printf("通过AI检查成功 (耗时: %v): %s\n", duration, aiResult)
		return
	}

	fmt.Printf("RAG检查成功 (耗时: %v): %s\n", duration, result)
}
