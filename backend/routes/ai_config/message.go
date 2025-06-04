package ai_config

import (
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getAIMessageService() *service.AIMessageService {
	// 或 database.DB，确保不是 nil
	ragClient := service.NewRAGClient("")
	return service.NewAIMessageService(database.DB, ragClient)
}

// GetAIMessage 获取AI Code Review列表
// @Summary 获取AI Code Review列表
// @Description 获取AI代码审查记录列表
// @Tags AI Code Review
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param project_id query int false "项目ID"
// @Param id query int false "审查记录ID"
// @Param current query int false "当前页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=dto.AIMessageListResponse}
// @Router /v1/ai/message [get]
func GetAIMessage(c *gin.Context) {
	var req dto.AIMessageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	// 转换参数
	params := make(map[string]interface{})
	if req.ProjectID > 0 {
		params["projectId"] = req.ProjectID
	}
	if req.ID > 0 {
		params["id"] = req.ID
	}
	if req.ProjectNamespace != "" {
		params["projectNamespace"] = req.ProjectNamespace
	}
	if req.ProjectName != "" {
		params["projectName"] = req.ProjectName
	}
	if req.CreateTime > 0 {
		params["createTime"] = req.CreateTime
	}
	if req.Passed != 0 {
		params["passed"] = req.Passed
	}

	// 计算分页参数
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	params["offset"] = (req.Current - 1) * req.PageSize
	params["limit"] = req.PageSize

	// 获取服务实例并查询数据
	rows, total, err := getAIMessageService().GetAIMessage(params)
	if err != nil {
		response.Error(c, err, "获取AI消息列表失败", 500)
		return
	}

	response.Success(c, dto.AIMessageListResponse{
		Data:  rows,
		Count: total,
	}, "获取成功", 0)
}

// CreateAIMessage 创建AI Code Review记录
// @Summary 创建AI Code Review记录
// @Description 创建新的AI代码审查记录
// @Tags AI Code Review
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body dto.AIMessageCreateRequest true "AI代码审查信息"
// @Success 200 {object} response.Response{data=dto.AIMessageCreateResponse}
// @Router /v1/ai/message [post]
func CreateAIMessage(c *gin.Context) {
	var req dto.AIMessageCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	// 创建AI消息
	aiMessage := &model.AImessage{
		ProjectID:  req.ProjectID,
		MergeURL:   req.MergeURL,
		MergeID:    strconv.FormatUint(uint64(req.MergeID), 10),
		AIModel:    req.AIModel,
		Rule:       model.RuleType(req.Rule),
		RuleID:     req.RuleID,
		Result:     req.Result,
		CreateTime: time.Now().Unix(),
	}

	// 获取服务实例并创建数据
	id, err := getAIMessageService().CreateAIMessage(aiMessage)
	if err != nil {
		response.Error(c, err, "创建AI消息失败", 500)
		return
	}

	response.Success(c, dto.AIMessageCreateResponse{
		ID:       id,
		Comments: req.Result,
	}, "创建成功", 0)
}

// UpdateAIMessage 更新AI Code Review记录
// @Summary 更新AI Code Review记录
// @Description 更新AI代码审查记录的人工评分和备注
// @Tags AI Code Review
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body dto.AIMessageUpdateRequest true "更新信息"
// @Success 200 {object} response.Response
// @Router /v1/ai/message [put]
func UpdateAIMessage(c *gin.Context) {
	var req dto.AIMessageUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}

	// 获取服务实例并更新数据
	err := getAIMessageService().UpdateHumanRatingAndRemark(model.AIMessageUpdate{
		ID:          req.ID,
		HumanRating: req.HumanRating,
		Remark:      req.Remark,
	})
	if err != nil {
		response.Error(c, err, "更新AI消息失败", 500)
		return
	}

	response.Success(c, nil, "更新成功", 0)
}

// GetCheckCount 获取检查数量
// @Summary 获取检查数量
// @Description 获取检查数量
// @Tags AI Code Review
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Param project_id query int false "项目ID"
// @Param project_namespace query string false "项目命名空间"
// @Param project_name query string false "项目名称"
// @Param passed query int false "-1未通过 1通过"
// @Success 200 {object} response.Response{data=dto.CheckCountResponse}
// @Router /v1/ai/merge/check-count [get]
func GetCheckCount(c *gin.Context) {
	// 检查t_message表中，passed为false的记录，统计每个项目的数量，查询最近30天的数据
	var req dto.CheckCountRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	// 检查开始时间和结束时间是否存在，查询对应的范围
	if req.StartTime == 0 {
		req.StartTime = time.Now().AddDate(0, 0, -30).Unix()
	}
	if req.EndTime == 0 {
		req.EndTime = time.Now().Unix()
	}

	total, err := getAIMessageService().GetCheckCount(req)
	if err != nil {
		response.Error(c, err, "获取失败", 500)
		return
	}

	response.Success(c, dto.CheckCountResponse{
		Count: int(total),
	}, "获取成功", 0)
}

// GetProblemChart 获取问题图表
// @Summary 获取问题图表
// @Description 获取问题图表
// @Tags AI Code Review
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Param project_id query int false "项目ID"
// @Param project_namespace query string false "项目命名空间"
// @Param project_name query string false "项目名称"
// @Param passed query int false "-1未通过 1通过"
// @Success 200 {object} response.Response{data=dto.AIProblemCountResponse}
// @Router /v1/ai/merge/problem-chart [get]
func GetProblemChart(c *gin.Context) {
	// 检查t_message表中，passed为false的记录，统计每个项目的数量，查询最近30天的数据
	var req dto.AIProblemCountRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	// 检查开始时间和结束时间是否存在，查询对应的范围
	if req.StartTime == 0 {
		req.StartTime = time.Now().AddDate(0, 0, -30).Unix()
	}
	if req.EndTime == 0 {
		req.EndTime = time.Now().Unix()
	}
	// 拉取最近30天的数据
	rows, err := getAIMessageService().GetProblemChart(req)
	if err != nil {
		response.Error(c, err, "获取失败", 500)
		return
	}

	response.Success(c, rows, "获取成功", 0)
}
