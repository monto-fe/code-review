package gitlab

import (
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service"

	"github.com/gin-gonic/gin"
)

// GetGitlabList 获取 Gitlab 列表
// @Summary 获取 Gitlab 列表
// @Description 获取所有 Gitlab 配置信息
// @Tags Gitlab
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response
// @Router /v1/gitlab [get]
func GetGitlabList(c *gin.Context) {
	gitlabService := service.NewGitlabService(database.DB)
	gitlabList, err := gitlabService.GetGitlabInfo()
	if err != nil {
		response.Error(c, err, "获取Gitlab列表失败", 500)
		return
	}
	response.Success(c, gin.H{
		"data": gitlabList,
	}, "获取成功", 0)
}

// CreateGitlabToken 创建 Gitlab Token
// @Summary 创建 Gitlab Token
// @Description 创建新的 Gitlab Token
// @Tags Gitlab
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body model.GitlabInfoCreate true "Gitlab 配置信息"
// @Success 200 {object} response.Response
// @Router /v1/gitlab [post]
func CreateGitlabToken(c *gin.Context) {
	var req model.GitlabInfoCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	gitlabService := service.NewGitlabService(database.DB)
	result, err := gitlabService.CreateGitlabToken(req)
	if err != nil {
		response.Error(c, err, "创建Gitlab Token失败", 500)
		return
	}

	response.Success(c, gin.H{
		"data": result,
	}, "创建成功", 0)
}

// UpdateGitlabToken 更新 Gitlab Token
// @Summary 更新 Gitlab Token
// @Description 更新现有的 Gitlab Token
// @Tags Gitlab
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body model.GitlabInfo true "Gitlab 配置信息"
// @Success 200 {object} response.Response
// @Router /v1/gitlab [put]
func UpdateGitlabToken(c *gin.Context) {
	var req model.GitlabInfoUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	if req.ID == 0 {
		response.Error(c, nil, "id is required", 400)
		return
	}

	gitlabService := service.NewGitlabService(database.DB)
	result, err := gitlabService.UpdateGitlabInfo(req)
	if err != nil {
		response.Error(c, err, "更新Gitlab Token失败", 500)
		return
	}

	if result == nil {
		response.Error(c, nil, "更新失败", 400)
		return
	}

	response.Success(c, gin.H{
		"data": result,
	}, "更新成功", 0)
}

// DeleteGitlabToken 删除 Gitlab Token
// @Summary 删除 Gitlab Token
// @Description 删除指定的 Gitlab Token
// @Tags Gitlab
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body dto.GitlabDeleteRequest true "Gitlab Token ID"
// @Success 200 {object} response.Response
// @Router /v1/gitlab [delete]
func DeleteGitlabToken(c *gin.Context) {
	var req dto.GitlabDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	gitlabService := service.NewGitlabService(database.DB)
	err := gitlabService.DeleteGitlabToken(req.ID)
	if err != nil {
		response.Error(c, err, "删除Gitlab Token失败", 500)
		return
	}

	response.Success(c, nil, "删除成功", 0)
}

// RefreshGitlabToken 刷新 Gitlab Token
// @Summary 刷新 Gitlab Token
// @Description 刷新所有 Gitlab Token
// @Tags Gitlab
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response
// @Router /v1/gitlab/token/refresh [post]
func RefreshGitlabToken(c *gin.Context) {
	// TODO: 实现刷新 Gitlab Token 逻辑
}
