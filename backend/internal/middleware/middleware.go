package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupMiddleware 设置所有中间件
func SetupMiddleware(r *gin.Engine) {
	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 配置日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 配置请求体大小限制
	r.MaxMultipartMemory = 50 << 20 // 50 MiB
}
