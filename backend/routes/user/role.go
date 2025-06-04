package user

import (
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service"

	"github.com/gin-gonic/gin"
)

func GetRoleList(c *gin.Context) {
	var req dto.RoleQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	roleService := service.NewRoleService(database.DB)
	roles, total, err := roleService.GetRoleList(req)
	if err != nil {
		response.Error(c, err, "获取失败", 500)
		return
	}
	response.Success(c, dto.RoleListResponse{
		Data:  roles,
		Count: int64(total),
	}, "获取成功", 0)
}

func CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	roleService := service.NewRoleService(database.DB)
	if err := roleService.CreateRoleWithResources(&req); err != nil {
		response.Error(c, err, "创建角色失败", 500)
		return
	}
	response.Success(c, nil, "创建成功", 0)
}

func UpdateRole(c *gin.Context) {
	var req dto.UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	roleService := service.NewRoleService(database.DB)
	if err := roleService.UpdateRole(&req); err != nil {
		response.Error(c, err, "更新角色失败", 500)
		return
	}
	response.Success(c, nil, "更新成功", 0)
}

func DeleteRole(c *gin.Context) {
	var req dto.DeleteRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "参数错误", 400)
		return
	}
	roleService := service.NewRoleService(database.DB)
	if err := roleService.DeleteRole(&req); err != nil {
		response.Error(c, err, "删除角色失败", 500)
	}
}
