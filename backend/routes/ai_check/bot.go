package ai_check

import (
	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/constants"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service/infrastructure"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BotReview 机器人代码审查
// @Summary 机器人代码审查
// @Description 使用指定的机器人角色对代码进行审查，支持多个内置机器人角色
// @Tags 机器人审查
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param request body dto.BotReviewRequest true "审查请求参数" example({"bot_name":"security_reviewer","type":"OpenAI","model":"gpt-4","api_url":"https://api.openai.com/v1","code_content":"function test() { return 'hello'; }","additional_prompt":"请特别关注性能问题"})
// @Success 200 {object} response.Response{data=dto.BotReviewResult}
// @Router /v1/ai/check/bot [post]
func BotReview(c *gin.Context) {
	var req dto.BotReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	// 获取用户ID（用于后续扩展，如记录审查历史）
	userIDStr := c.GetHeader("userId")
	if userIDStr == "" {
		response.Error(c, nil, "用户ID不存在", int(constants.RetCodeOperatorMissing))
		return
	}

	_, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.Error(c, err, "用户ID格式错误", int(constants.RetCodeInvalidUserID))
		return
	}

	// 执行机器人审查
	botService := infrastructure.GetBotService()
	result, err := botService.ExecuteBotReview(req.BotName, req.Type, req.Model, req.APIURL, req.CodeContent, req.AdditionalPrompt)
	if err != nil {
		response.Error(c, err, "机器人审查失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, result, "机器人审查完成", int(constants.RetCodeSuccess))
}

// GetBotRoles 获取所有机器人角色
// @Summary 获取所有机器人角色
// @Description 获取系统中所有可用的机器人角色列表
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response{data=dto.BotRolesListResponse}
// @Router /v1/ai/check/bot/roles [get]
func GetBotRoles(c *gin.Context) {
	botService := infrastructure.GetBotService()
	roles := botService.GetAllBotRoles()

	// 转换为响应格式
	var roleResponses []dto.BotRoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, dto.BotRoleResponse{
			Name:        role.Name,
			Description: role.Description,
			Category:    role.Category,
		})
	}

	response.Success(c, dto.BotRolesListResponse{
		Roles: roleResponses,
	}, "获取机器人角色列表成功", int(constants.RetCodeSuccess))
}

// GetBotRolesByCategory 根据分类获取机器人角色
// @Summary 根据分类获取机器人角色
// @Description 根据指定分类获取机器人角色列表
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param category path string true "机器人分类" example("安全")
// @Success 200 {object} response.Response{data=dto.BotRolesListResponse}
// @Router /v1/ai/check/bot/roles/{category} [get]
func GetBotRolesByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		response.Error(c, nil, "分类参数不能为空", int(constants.RetCodeBadRequest))
		return
	}

	botService := infrastructure.GetBotService()
	roles := botService.GetBotRolesByCategory(category)

	// 转换为响应格式
	var roleResponses []dto.BotRoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, dto.BotRoleResponse{
			Name:        role.Name,
			Description: role.Description,
			Category:    role.Category,
		})
	}

	response.Success(c, dto.BotRolesListResponse{
		Roles: roleResponses,
	}, "获取机器人角色列表成功", int(constants.RetCodeSuccess))
}

// GetBotRoleDetail 获取机器人角色详情
// @Summary 获取机器人角色详情
// @Description 获取指定机器人角色的详细信息
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param bot_name path string true "机器人名称" example("安全审查机器人")
// @Success 200 {object} response.Response{data=infrastructure.BotRole}
// @Router /v1/ai/check/bot/roles/detail/{bot_name} [get]
func GetBotRoleDetail(c *gin.Context) {
	botName := c.Param("bot_name")
	if botName == "" {
		response.Error(c, nil, "机器人名称不能为空", int(constants.RetCodeBadRequest))
		return
	}

	botService := infrastructure.GetBotService()
	role, err := botService.GetBotRole(botName)
	if err != nil {
		response.Error(c, err, "机器人角色不存在", int(constants.RetCodeBadRequest))
		return
	}

	response.Success(c, role, "获取机器人角色详情成功", int(constants.RetCodeSuccess))
}
