package gitlab

import (
	"code-review-go/internal/cache"
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/constants"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service"
	"strconv"

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
	// 获取Gitlab列表
	gitlabList, err := gitlabService.GetGitlabInfo()
	if err != nil {
		response.Error(c, err, "获取Gitlab列表失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, gin.H{
		"data": gitlabList,
	}, "获取成功", int(constants.RetCodeSuccess))
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
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	gitlabService := service.NewGitlabService(database.DB)
	// 创建Gitlab Token
	result, err := gitlabService.CreateGitlabToken(req)
	if err != nil {
		response.Error(c, err, "创建Gitlab Token失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, gin.H{
		"data": result,
	}, "创建成功", int(constants.RetCodeSuccess))
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
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	if req.ID == 0 {
		response.Error(c, nil, "id is required", int(constants.RetCodeBadRequest))
		return
	}

	gitlabService := service.NewGitlabService(database.DB)
	// 更新Gitlab Token
	result, err := gitlabService.UpdateGitlabInfo(req)
	if err != nil {
		response.Error(c, err, "更新Gitlab Token失败", int(constants.RetCodeInternalError))
		return
	}

	if result == nil {
		response.Error(c, nil, "更新失败", int(constants.RetCodeBadRequest))
		return
	}

	response.Success(c, gin.H{
		"data": result,
	}, "更新成功", int(constants.RetCodeSuccess))
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
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	gitlabService := service.NewGitlabService(database.DB)
	// 删除Gitlab Token
	err := gitlabService.DeleteGitlabToken(req.ID)
	if err != nil {
		response.Error(c, err, "删除Gitlab Token失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, nil, "删除成功", int(constants.RetCodeSuccess))
}

// RefreshGitlabToken 刷新 Gitlab Token
// @Summary 刷新 Gitlab Token
// @Description 刷新所有 Gitlab Token 缓存
// @Tags Gitlab
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response
// @Router /v1/gitlab/token/refresh [post]
func RefreshGitlabToken(c *gin.Context) {
	// 刷新缓存
	err := cache.RefreshGitlabCache()
	if err != nil {
		response.Error(c, err, "刷新缓存失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, nil, "缓存刷新成功", int(constants.RetCodeSuccess))
}

// GetGitlabTokenDetail 获取 Gitlab Token 详情
// @Summary 获取 Gitlab Token 详情
// @Description 获取指定的 Gitlab Token 详情
// @Tags Gitlab
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param id query int true "Gitlab Token ID"
// @Success 200 {object} response.Response
// @Router /v1/gitlab/token [get]
func GetGitlabTokenDetail(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		response.Error(c, nil, "id参数不能为空", int(constants.RetCodeBadRequest))
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, err, "id参数格式错误", int(constants.RetCodeBadRequest))
		return
	}
	gitlabService := service.NewGitlabService(database.DB)
	// 获取Gitlab Token详情
	gitlabInfo, err := gitlabService.GetGitlabTokenDetail(uint(idInt))
	if err != nil {
		response.Error(c, err, "获取Gitlab Token 详情失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, gin.H{
		"data": gitlabInfo,
	}, "获取Gitlab Token 详情成功", int(constants.RetCodeSuccess))
}

// GetGitlabTokenProjects 获取Gitlab Token项目信息列表
// @Summary 获取Gitlab Token项目信息列表
// @Description 获取所有Gitlab Token的名称和对应的ProjectIds，全量拉取不分页
// @Tags Gitlab管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response{data=dto.GitlabTokenProjectListResponse}
// @Router /v1/gitlab/token/projects [get]
func GetGitlabTokenProjects(c *gin.Context) {
	gitlabService := service.NewGitlabService(database.DB)

	// 获取Token项目信息列表
	tokenProjects, err := gitlabService.GetGitlabTokenProjects()
	if err != nil {
		response.Error(c, err, "获取Gitlab Token项目信息失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, dto.GitlabTokenProjectListResponse{
		Data: tokenProjects,
	}, "获取成功", int(constants.RetCodeSuccess))
}
