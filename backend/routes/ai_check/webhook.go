package ai_check

import (
	"fmt"
	"strconv"

	"code-review-go/internal/cache"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service"

	"github.com/gin-gonic/gin"
)

type ProjectInfo struct {
	ID                int    `json:"id"`
	PathWithNamespace string `json:"path_with_namespace"`
}

type ObjectAttributes struct {
	IID          int    `json:"iid"`
	URL          string `json:"url"`
	Action       string `json:"action"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
}

type WebhookBody struct {
	Project          ProjectInfo      `json:"project"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
}

// AICheck 处理 AI 检查请求
// @Summary 触发 AI 代码审查
// @Description 处理 GitLab 合并请求的 webhook，自动触发 AI 代码审查与评论
// @Tags Webhook
// @Accept json
// @Produce json
// @Param data body ai_check.WebhookBody true "Webhook 触发参数"
// @Success 200 {object} response.Response
// @Router /v1/webhook/merge [post]
func AICheck(c *gin.Context) {

	// 1. 解析请求体
	var body WebhookBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误", "error": err.Error()})
		return
	}

	projectId := body.Project.ID
	iid := body.ObjectAttributes.IID
	mergeURL := body.ObjectAttributes.URL
	action := body.ObjectAttributes.Action
	sourceBranch := body.ObjectAttributes.SourceBranch
	targetBranch := body.ObjectAttributes.TargetBranch
	pathWithNamespace := body.Project.PathWithNamespace

	fmt.Println("body", mergeURL, pathWithNamespace)

	// 2. 立即响应
	response.Success(c, gin.H{
		"projectId":      projectId,
		"mergeRequestId": iid,
	}, "Webhook处理成功，等待AI检测", 0)

	// 3. 只处理 open/update
	if action != "open" && action != "update" {
		return
	}

	// 4. 获取gitlab缓存信息（伪代码，需结合你的实现）
	gitlabCache := cache.GetGitlabCache()
	gitlabToken, gitlabInfoResult, ok := cache.FindTokenByProjectID(strconv.Itoa(projectId), gitlabCache)
	if !ok {
		fmt.Println("请配置gitlab Token")
		return
	}

	gitlabAPI := gitlabInfoResult.Config.API
	webhookURL := gitlabInfoResult.Config.WebhookURL
	webhookStatus := gitlabInfoResult.Config.WebhookStatus

	sourceBranchCfg := gitlabInfoResult.Config.SourceBranch
	targetBranchCfg := gitlabInfoResult.Config.TargetBranch

	if sourceBranchCfg != "" && sourceBranchCfg != sourceBranch {
		fmt.Println("源分支不匹配")
		return
	}
	if targetBranchCfg != "" && targetBranchCfg != targetBranch {
		fmt.Println("目标分支不匹配")
		return
	}

	// 5. 获取merge信息
	mergeRequest, err := service.GetMergeRequestInfo(gitlabAPI, strconv.Itoa(projectId), gitlabToken)
	if err != nil {
		fmt.Println("获取merge信息失败:", err)
		return
	}

	// 6. 获取diff
	var diff []model.Change
	if gitlabAPI != "" && projectId != 0 && iid != 0 {
		diff, _ = service.GetMergeRequestDiff(gitlabAPI, strconv.Itoa(projectId), strconv.Itoa(iid), gitlabToken)
	}

	// 7. AI 检查
	comments, aiMessageId, err := service.CheckMergeRequestWithAI(mergeRequest, diff, gitlabAPI, gitlabToken)
	if err != nil {
		fmt.Println("AI检查失败:", err)
		return
	}
	fmt.Println("AI检查结果:", comments, aiMessageId)

	// 8. 写入GitLab评论
	if comments != "" && iid != 0 && projectId != 0 && gitlabToken != "" && gitlabAPI != "" {
		_, err := service.PostCommentToGitLab(gitlabAPI, projectId, iid, gitlabToken, comments)
		if err != nil {
			fmt.Println("评论提交失败:", err)
		}
	}

	// 9. 推送webhook
	if webhookURL != "" && webhookStatus == 1 {
		webhookContent := service.PushWeChatInfo(pathWithNamespace, mergeURL, comments, aiMessageId)
		_ = service.SendMarkdownToWechatBot(webhookURL, webhookContent)
		fmt.Println("推送webhook成功", mergeRequest)
	}

}
