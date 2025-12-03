package ai_check

import (
	"code-review-go/internal/database"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/constants"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service/business"

	"github.com/gin-gonic/gin"
)

// CreateReviewWhitelist 创建白名单记录
// @Summary 创建代码审查白名单
// @Description 创建新的代码审查白名单记录
// @Tags 代码审查白名单
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body model.ReviewWhitelistCreate true "白名单信息"
// @Success 200 {object} response.Response{data=model.ReviewWhitelist}
// @Router /v1/webhook/whitelist [post]
func CreateReviewWhitelist(c *gin.Context) {
	var req model.ReviewWhitelistCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	// 从请求头获取操作人
	operator := c.GetHeader("remoteUser")
	if operator == "" {
		response.Error(c, nil, "操作人信息缺失", int(constants.RetCodeOperatorMissing))
		return
	}
	req.Operator = operator

	// 获取服务实例并创建白名单
	service := business.NewReviewWhitelistService(database.DB)
	whitelist, err := service.Create(&req)
	if err != nil {
		response.Error(c, err, "创建白名单失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, whitelist, "创建成功", int(constants.RetCodeSuccess))
}

// UpdateReviewWhitelist 更新白名单记录
// @Summary 更新代码审查白名单
// @Description 更新现有的代码审查白名单记录
// @Tags 代码审查白名单
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body model.ReviewWhitelistUpdate true "白名单信息"
// @Success 200 {object} response.Response
// @Router /v1/webhook/whitelist [put]
func UpdateReviewWhitelist(c *gin.Context) {
	var req model.ReviewWhitelistUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	// 从请求头获取操作人
	operator := c.GetHeader("remoteUser")
	if operator != "" {
		req.Operator = operator
	}

	// 获取服务实例并更新白名单
	service := business.NewReviewWhitelistService(database.DB)
	err := service.Update(&req)
	if err != nil {
		response.Error(c, err, "更新白名单失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, nil, "更新成功", int(constants.RetCodeSuccess))
}

// DeleteReviewWhitelist 删除白名单记录
// @Summary 删除代码审查白名单
// @Description 删除现有的代码审查白名单记录
// @Tags 代码审查白名单
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param data body model.ReviewWhitelistDelete true "白名单ID"
// @Success 200 {object} response.Response
// @Router /v1/webhook/whitelist [delete]
func DeleteReviewWhitelist(c *gin.Context) {
	var req model.ReviewWhitelistDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", int(constants.RetCodeBadRequest))
		return
	}

	// 获取服务实例并删除白名单
	service := business.NewReviewWhitelistService(database.DB)
	err := service.Delete(req.ID)
	if err != nil {
		response.Error(c, err, "删除白名单失败", int(constants.RetCodeInternalError))
		return
	}

	response.Success(c, nil, "删除成功", int(constants.RetCodeSuccess))
}

// GetReviewWhitelistList 获取白名单列表
// @Summary 获取代码审查白名单列表
// @Description 分页获取代码审查白名单列表
// @Tags 代码审查白名单
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param current query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param project_name query string false "项目名称（模糊查询）"
// @Param project_id query int false "项目ID"
// @Success 200 {object} response.Response{data=object}
// @Router /v1/webhook/whitelist/list [post]
func GetReviewWhitelistList(c *gin.Context) {
	var query model.ReviewWhitelistQuery

	// 从查询参数获取分页信息
	query.Current = utils.GetQueryInt(c, "current", 1)
	query.PageSize = utils.GetQueryInt(c, "pageSize", 10)
	query.ProjectName = c.Query("project_name")
	query.ProjectID = uint(utils.GetQueryInt(c, "project_id", 0))

	// 获取服务实例并查询数据
	service := business.NewReviewWhitelistService(database.DB)
	whitelists, total, err := service.GetList(&query)
	if err != nil {
		response.Error(c, err, "获取白名单列表失败", int(constants.RetCodeInternalError))
		return
	}

	// 转换为响应结构体
	responseList := make([]model.ReviewWhitelistResponse, len(whitelists))
	for i, wl := range whitelists {
		responseList[i] = wl.ToResponse()
	}

	response.Success(c, map[string]interface{}{
		"data":  responseList,
		"total": total,
	}, "获取成功", int(constants.RetCodeSuccess))
}
