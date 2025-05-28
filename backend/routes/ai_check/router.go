package ai_check

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/webhook/merge", AICheck)
}
