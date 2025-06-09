package gitlab

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup) {
	gitlab := r.Group("/gitlab")
	{
		gitlab.GET("", GetGitlabList)
		gitlab.POST("", CreateGitlabToken)
		gitlab.PUT("", UpdateGitlabToken)
		gitlab.DELETE("", DeleteGitlabToken)
		gitlab.GET("/token/:id", GetGitlabTokenDetail)
	}
}
