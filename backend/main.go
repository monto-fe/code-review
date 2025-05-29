package main

import (
	"log"

	"code-review-go/config"
	"code-review-go/docs"
	"code-review-go/internal/cache"
	"code-review-go/internal/database"
	"code-review-go/internal/middleware"
	"code-review-go/internal/websocket"
	"code-review-go/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Code Review API
// @version 1.0
// @description Code Review API documentation
// @host localhost:9000
// @BasePath /v1
func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置 Gin 模式
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化数据库连接
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 执行数据库迁移
	// if err := database.AutoMigrate(); err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }

	// 初始化 Gitlab 缓存
	if err := cache.InitGitlabCache(); err != nil {
		log.Fatalf("Failed to initialize gitlab cache: %v", err)
	}

	// 初始化 AI 配置缓存
	if err := cache.InitAIConfigCache(); err != nil {
		log.Fatalf("Failed to initialize ai_config cache: %v", err)
	}

	// 初始化 Gin
	r := gin.Default()

	// 设置中间件
	middleware.SetupMiddleware(r)

	// 初始化 WebSocket Hub
	hub := websocket.NewHub()

	// 初始化 Swagger
	docs.SwaggerInfo.Title = "Code Review API"
	docs.SwaggerInfo.Description = "Code Review API documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.IP + ":" + cfg.Port
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.InfoInstanceName = "swagger"

	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 初始化路由
	routes.InitializeRoutes(r)
	// 初始化 WebSocket 路由
	routes.InitializeWebSocketRoutes(r, hub)

	// 启动服务器
	log.Printf("Server is running on http://%s:%s", cfg.IP, cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
