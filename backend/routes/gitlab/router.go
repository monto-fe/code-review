package gitlab

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup) {
	gitlab := r.Group("/gitlab")
	{
		gitlab.GET("", GetGitlabList)
		gitlab.POST("", CreateGitlabToken)
		gitlab.PUT("", UpdateGitlabToken)
		gitlab.DELETE("", DeleteGitlabToken)
		gitlab.GET("/token", GetGitlabTokenDetail)
		gitlab.POST("/token/refresh", RefreshGitlabToken)
		gitlab.GET("/token/projects", GetGitlabTokenProjects) // 新增：获取Token项目信息列表
	}
}
