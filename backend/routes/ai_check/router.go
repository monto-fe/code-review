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

		// 机器人审查相关路由
		check.POST("/check/bot", BotReview)                              // 机器人代码审查
		check.GET("/check/bot/roles", GetBotRoles)                       // 获取所有机器人角色
		check.GET("/check/bot/roles/:category", GetBotRolesByCategory)   // 根据分类获取机器人角色
		check.GET("/check/bot/roles/detail/:bot_name", GetBotRoleDetail) // 获取机器人角色详情
	}
}
