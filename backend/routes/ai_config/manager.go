package ai_config

import (
	"code-review-go/internal/database"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetAIManagerList 获取AI管理器列表
// @Summary 获取AI管理器列表
// @Description 获取所有AI管理器列表
// @Tags AI配置
// @Accept json
// @Produce json
// @Param jwt_token header string true "JWT认证Token"
// @Success 200 {object} response.Response{data=[]model.AIManager}
// @Router /v1/ai/manager [get]
func GetAIManagerList(c *gin.Context) {
	// 查询所有数据
	var managers []model.AIManager
	if err := database.DB.Find(&managers).Error; err != nil {
		response.Error(c, err, "获取AI管理器列表失败", 500)
		return
	}

	response.Success(c, managers, "获取成功", 0)
}
