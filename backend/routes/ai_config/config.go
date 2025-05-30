package ai_config

import (
	"code-review-go/internal/database"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service"

	"github.com/gin-gonic/gin"
)

// GetAIConfig 获取AI配置列表
// @Summary 获取 AI 配置
// @Description 获取 AI 配置信息
// @Tags AI配置
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=object}
// @Router /v1/ai/config [get]
func GetAIConfig(c *gin.Context) {
	// 分页参数
	page := utils.GetQueryInt(c, "page", 1)
	pageSize := utils.GetQueryInt(c, "page_size", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 获取服务实例并查询数据
	aiConfigManager, err := service.NewAIConfigManager(database.DB)
	if err != nil {
		response.Error(c, err, "初始化AI配置管理器失败", 500)
		return
	}

	configs, total, err := aiConfigManager.GetConfigListPaged(page, pageSize)
	if err != nil {
		response.Error(c, err, "获取AI配置列表失败", 500)
		return
	}

	// 转换为不含 api_key 的响应结构体
	responseList := make([]model.AIConfigResponse, len(configs))
	for i, config := range configs {
		responseList[i] = config.ToResponse()
	}

	response.Success(c, map[string]interface{}{
		"data":  responseList,
		"total": total,
	}, "获取成功", 0)
}

// CreateAIConfig 创建 AI 配置
// @Summary 创建 AI 配置
// @Description 创建新的 AI 配置
// @Tags AI配置
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body model.AIConfigCreate true "AI配置信息"
// @Success 200 {object} response.Response{data=model.AIConfig}
// @Router /v1/ai/config [post]
func CreateAIConfig(c *gin.Context) {
	var req model.AIConfigCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	// 获取服务实例并创建配置
	aiConfigManager, err := service.NewAIConfigManager(database.DB)
	if err != nil {
		response.Error(c, err, "初始化AI配置管理器失败", 500)
		return
	}

	config, err := aiConfigManager.AddAIConfig(req)
	if err != nil {
		response.Error(c, err, "创建AI配置失败", 500)
		return
	}

	response.Success(c, config, "创建成功", 0)
}

// UpdateAIConfig 更新 AI 配置
// @Summary 更新 AI 配置
// @Description 更新现有的 AI 配置
// @Tags AI配置
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body model.AIConfigUpdate true "AI配置信息"
// @Success 200 {object} response.Response
// @Router /v1/ai/config [put]
func UpdateAIConfig(c *gin.Context) {
	var req model.AIConfigUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	// 获取服务实例并更新配置
	aiConfigManager, err := service.NewAIConfigManager(database.DB)
	if err != nil {
		response.Error(c, err, "初始化AI配置管理器失败", 500)
		return
	}

	if err := aiConfigManager.UpdateAIConfig(req); err != nil {
		response.Error(c, err, "更新AI配置失败", 500)
		return
	}

	response.Success(c, nil, "更新成功", 0)
}

// DeleteAIConfig 删除AI配置
// @Summary 删除AI配置
// @Description 删除现有的 AI 配置
// @Tags AI配置
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body model.AIConfigDelete true "配置ID"
// @Success 200 {object} response.Response
// @Router /v1/ai/config [delete]
func DeleteAIConfig(c *gin.Context) {
	var req model.AIConfigDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	aiConfigManager, err := service.NewAIConfigManager(database.DB)
	if err != nil {
		response.Error(c, err, "初始化AI配置管理器失败", 500)
		return
	}

	if err := aiConfigManager.DeleteAIConfig(req.ID); err != nil {
		response.Error(c, err, "删除AI配置失败", 500)
		return
	}

	response.Success(c, nil, "删除成功", 0)
}
