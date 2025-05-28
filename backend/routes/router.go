package routes

import (
	"code-review-go/routes/ai_check"
	"code-review-go/routes/ai_config"

	"code-review-go/routes/gitlab"
	"code-review-go/routes/user"

	"github.com/gin-gonic/gin"
)

// InitializeRoutes 初始化所有路由
func InitializeRoutes(r *gin.Engine) {
	// API 版本分组
	v1 := r.Group("/v1")
	{
		// 注册用户路由
		user.RegisterRoutes(v1)
		// Webhook相关路由
		ai_check.RegisterRoutes(v1)
		// 规则相关路由
		ai_config.RegisterRoutes(v1)
		// Gitlab相关路由
		gitlab.RegisterRoutes(v1)
	}
}
