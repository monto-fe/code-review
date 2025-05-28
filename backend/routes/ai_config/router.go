package ai_config

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup) {
	aiConfig := r.Group("/ai")
	{
		aiConfig.GET("/config", GetAIConfig)
		aiConfig.POST("/config", CreateAIConfig)
		aiConfig.PUT("/config", UpdateAIConfig)
		aiConfig.DELETE("/config", DeleteAIConfig)
		aiConfig.GET("/message", GetAIMessage)
		aiConfig.POST("/message", CreateAIMessage)
		aiConfig.PUT("/message", UpdateAIMessage)
	}
}
