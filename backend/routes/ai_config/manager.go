package ai_config

import (
	"code-review-go/internal/database"
	"code-review-go/internal/pkg/response"
	"code-review-go/internal/service"

	"github.com/gin-gonic/gin"
)

// GetAIManagerList 获取AI管理器列表
// @Summary 获取AI管理器列表
// @Description 获取所有AI管理器列表，服务的集成能力，不允许删除
// @Tags AI配置
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response{data=[]model.AIManager}
// @Router /v1/ai/manager [get]
func GetAIManagerList(c *gin.Context) {
	managerService := service.NewAIManagerService(database.DB)
	managers, err := managerService.GetManagerList()
	if err != nil {
		response.Error(c, err, "获取AI管理器列表失败", 500)
		return
	}
	response.Success(c, managers, "获取成功", 0)
}
