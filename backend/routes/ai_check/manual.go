package ai_check

import (
	"strconv"

	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/constants"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service/business"

	"github.com/gin-gonic/gin"
)

// ManualCheck 手动触发代码审核
// @Summary 手动触发代码审核
// @Description 用户输入合并请求链接和AI模型ID，手动触发代码审核。支持配置多个机器人，每个机器人都有对应的提示词。系统会异步从缓存中获取项目信息和AI模型信息并更新到任务中。
// @Tags 手动代码审核
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param request body dto.ManualCheckRequest true "审核请求参数" example({"merge_url":"https://gitlab.example.com/project/repo/-/merge_requests/123","ai_model_id":1,"bot_configs":[{"bot_name":"安全审查机器人","bot_prompt":"请重点检查代码中的安全漏洞，包括SQL注入、XSS攻击、权限绕过等问题"},{"bot_name":"性能优化机器人","bot_prompt":"请分析代码性能问题，关注算法复杂度、内存使用、数据库查询优化等"}]})
// @Success 200 {object} response.Response{data=dto.ManualCheckResponse}
// @Router /v1/ai/check/manual [post]
func ManualCheck(c *gin.Context) {
	var req dto.ManualCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	// 获取用户ID
	userIDStr := c.GetHeader("userId")
	if userIDStr == "" {
		response.Error(c, nil, "用户ID不存在", int(constants.RetCodeOperatorMissing))
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.Error(c, err, "用户ID格式错误", int(constants.RetCodeInvalidUserID))
		return
	}

	// 创建手动审核任务
	manualCheckService := business.GetManualCheckService()
	task, err := manualCheckService.CreateManualCheckTask(req, uint(userID))
	if err != nil {
		response.Error(c, err, "创建审核任务失败", int(constants.RetCodeInternalError))
		return
	}

	// 异步执行审核任务
	go func() {
		if err := manualCheckService.ExecuteManualCheckTask(task.ID); err != nil {
			// 记录错误日志，但不影响响应
			// log.Printf("执行审核任务失败: %v", err)
		}
	}()

	response.Success(c, dto.ManualCheckResponse{
		TaskID:        task.ID,
		Status:        "processing",
		EstimatedTime: 30, // 预估30秒完成
	}, "审核任务已创建", int(constants.RetCodeSuccess))
}

// GetManualCheckHistory 获取手动审核历史
// @Summary 获取手动审核历史
// @Description 查询手动触发的审核历史记录
// @Tags 手动代码审核
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param current query int false "当前页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query int false "状态筛选：0-全部，1-进行中，2-完成，3-失败" default(0)
// @Param start_date query string false "开始日期" example(1703001600)
// @Param end_date query string false "结束日期" example(1705680000)
// @Success 200 {object} response.Response{data=dto.ManualCheckHistoryResponse}
// @Router /v1/ai/check/history [get]
func GetManualCheckHistory(c *gin.Context) {
	// 获取查询参数
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	req := dto.ManualCheckHistoryRequest{
		Current:   current,
		PageSize:  pageSize,
		Status:    status,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 获取用户ID
	userIDStr := c.GetHeader("userId")
	if userIDStr == "" {
		response.Error(c, nil, "用户ID不存在", int(constants.RetCodeOperatorMissing))
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.Error(c, err, "用户ID格式错误", int(constants.RetCodeInvalidUserID))
		return
	}

	// 查询审核历史
	manualCheckService := business.GetManualCheckService()
	result, err := manualCheckService.GetManualCheckHistory(req, uint(userID))
	if err != nil {
		response.Error(c, err, "查询审核历史失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, result, "查询成功", int(constants.RetCodeSuccess))
}

// GetManualCheckResult 获取手动审核结果详情
// @Summary 获取手动审核结果详情
// @Description 获取特定审核任务的详细结果
// @Tags 手动代码审核
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response{data=dto.ManualCheckResultResponse}
// @Router /v1/ai/check/result/{id} [get]
func GetManualCheckResult(c *gin.Context) {
	// 获取任务ID
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		response.Error(c, err, "任务ID格式错误", int(constants.RetCodeBadRequest))
		return
	}

	// 获取用户ID
	userIDStr := c.GetHeader("userId")
	if userIDStr == "" {
		response.Error(c, nil, "用户ID不存在", int(constants.RetCodeOperatorMissing))
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.Error(c, err, "用户ID格式错误", int(constants.RetCodeInvalidUserID))
		return
	}

	// 获取审核结果
	manualCheckService := business.GetManualCheckService()
	result, err := manualCheckService.GetManualCheckResult(uint(taskID), uint(userID))
	if err != nil {
		response.Error(c, err, "获取审核结果失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, result, "查询成功", int(constants.RetCodeSuccess))
}
