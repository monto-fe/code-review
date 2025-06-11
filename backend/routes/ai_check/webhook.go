package ai_check

import (
	"fmt"

	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/response"
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
		c.JSON(400, gin.H{"msg": "参数错误", "error": err.Error()})
		return
	}

	fmt.Println("body", body)

	// 2. 立即响应
	response.Success(c, gin.H{
		"projectId":      body.Project.ID,
		"mergeRequestId": body.ObjectAttributes.IID,
	}, "Webhook处理成功，等待AI检测", 0)

	// 3. 后续逻辑异步处理
	go handleAICheck(body)
}

func handleAICheck(body dto.WebhookBody) {
	if !service.ShouldProcessState(body.ObjectAttributes.State) {
		return
	}
	comments, err := service.CheckMergeRequestWithAI(body)
	if err != nil {
		fmt.Println("AI检查失败:", err)
		return
	}
	fmt.Println("AI检查结果:", comments)
}
