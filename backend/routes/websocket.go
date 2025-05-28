package routes

import (
	"code-review-go/internal/websocket"

	"github.com/gin-gonic/gin"
)

// InitializeWebSocketRoutes 初始化WebSocket路由
func InitializeWebSocketRoutes(r *gin.Engine, hub *websocket.Hub) {
	// 启动Hub
	go hub.Run()

	// WebSocket路由
	r.GET("/ws", websocket.HandleWebSocket(hub))
}
