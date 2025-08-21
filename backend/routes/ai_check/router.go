package ai_check

import (
	"code-review-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/webhook/merge", AICheck)

	// 手动审核相关路由
	check := r.Group("/ai")
	check.Use(middleware.AuthenticateJWT())
	{
		check.POST("/check/manual", ManualCheck)             // 手动触发代码审核
		check.GET("/check/history", GetManualCheckHistory)   // 获取审核历史
		check.GET("/check/result/:id", GetManualCheckResult) // 获取审核结果详情
	}
}
