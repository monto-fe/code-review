package ai_check

import (
	"fmt"
	"time"

	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/constants"
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

	fmt.Printf("webhook请求信息: ProjectID=%d, MergeRequestID=%d\n",
		body.Project.ID, body.ObjectAttributes.IID)

	// 2. 立即响应（不阻塞webhook）
	response.Success(c, gin.H{
		"projectId":      body.Project.ID,
		"mergeRequestId": body.ObjectAttributes.IID,
		"optimized":      true,
		"timestamp":      time.Now().Unix(),
	}, "AI检查已启动，请稍候查看结果", int(constants.RetCodeSuccess))

	// 💡 防止重复执行
	if utils.IsDuplicateWebhook(body) {
		fmt.Println("重复webhook请求，跳过处理")
		return
	}
	// 检查Merge Request状态
	if !service.ShouldProcessState(body) {
		fmt.Printf("跳过非opened状态的合并请求: %+v\n", body)
		return
	}

	// 3. 异步处理优化的RAG检查
	go handleOptimizedAICheck(body)
}

// handleOptimizedAICheck 处理优化的AI检查
func handleOptimizedAICheck(body dto.WebhookBody) {
	startTime := time.Now()
	fmt.Printf("开始代码评审检查: %s\n", startTime.Format("2006-01-02 15:04:05"))

	// 1. 使用RAG服务检查
	result, err := service.CheckMergeRequestWithAI(body)
	if err != nil {
		fmt.Printf("检查失败 (耗时: %v): %v\n", time.Since(startTime), err)
		return
	}
	fmt.Printf("检查成功 (耗时: %v): %s\n", time.Since(startTime), result)
}
