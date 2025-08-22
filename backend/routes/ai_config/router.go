package ai_config

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup) {
	aiConfig := r.Group("/ai")
	{
		aiConfig.GET("/config", GetAIConfig)
		aiConfig.POST("/config", CreateAIConfig)
		aiConfig.PUT("/config", UpdateAIConfig)
		aiConfig.DELETE("/config", DeleteAIConfig)
		aiConfig.POST("/message/list", GetAIMessage) // 获取AI消息列表
		aiConfig.POST("/message", CreateAIMessage)   // 创建AI消息
		aiConfig.PUT("/message", UpdateAIMessage)
		aiConfig.GET("/manager", GetAIManagerList)
		aiConfig.GET("/merge/check-count", GetCheckCount)
		aiConfig.GET("/merge/problem-chart", GetProblemChart)
		aiConfig.GET("/project-namespaces", GetProjectNamespaceList) // 新增：获取项目命名空间列表
	}
}
