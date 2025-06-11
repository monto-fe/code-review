package user

import (
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service"

	"github.com/gin-gonic/gin"
)

// GetResourceList 获取资源列表
// @Summary 获取资源列表
// @Description 获取资源列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response{data=dto.ResourceListResponse}
// @Router /v1/user/resource [get]
func GetResourceList(c *gin.Context) {
	var req model.ResourceQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	resourceService := service.NewResourceService(database.GetDB())
	resources, total, err := resourceService.GetResourceList(req)
	if err != nil {
		response.Error(c, err, "获取资源列表失败", 500)
		return
	}
	response.Success(c, map[string]interface{}{
		"data":  resources,
		"total": total,
	}, "获取成功", 0)
}

// CreateResource 创建资源
// @Summary 创建资源
// @Description 创建资源
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param resource body dto.Resource true "资源"
// @Success 200 {object} response.Response{data=dto.Resource}
// @Router /v1/user/resource [post]
func CreateResource(c *gin.Context) {
	var req model.ResourceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	resourceService := service.NewResourceService(database.GetDB())
	resource, err := resourceService.Create(&req)
	if err != nil {
		response.Error(c, err, "创建资源失败", 500)
		return
	}
	response.Success(c, resource, "创建成功", 0)
}

// UpdateResource 更新资源
// @Summary 更新资源
// @Description 更新资源
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Param resource body dto.Resource true "资源"
// @Success 200 {object} response.Response{data=dto.Resource}
// @Router /v1/user/resource [put]
func UpdateResource(c *gin.Context) {
	var resource dto.UpdateResourceReq
	if err := c.ShouldBindJSON(&resource); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	resourceService := service.NewResourceService(database.GetDB())
	if err := resourceService.Update(&resource); err != nil {
		response.Error(c, err, "更新资源失败", 500)
		return
	}
	response.Success(c, resource, "更新成功", 0)
}

// DeleteResource 删除资源
// @Summary 删除资源
// @Description 删除资源
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response{data=dto.Resource}
// @Router /v1/user/resource [delete]
func DeleteResource(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		response.Error(c, err, "参数错误", 400)
		return
	}
	resourceService := service.NewResourceService(database.GetDB())
	if err := resourceService.DeleteSelf(req.ID); err != nil {
		response.Error(c, err, "删除资源失败", 500)
		return
	}
	response.Success(c, nil, "删除成功", 0)
}
