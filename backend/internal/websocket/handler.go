package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// TODO: 后续优化，在GitlabToken加载时同步前端状态

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 表示一个WebSocket客户端
type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	mu   sync.Mutex
}

// Hub 维护所有活跃的客户端
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// NewHub 创建一个新的Hub实例
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 启动Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		client := &Client{
			conn: conn,
			send: make(chan []byte, 256),
			hub:  hub,
		}

		client.hub.register <- client

		// 启动goroutine处理读写
		go client.writePump()
		go client.readPump()
	}
}

// readPump 从WebSocket连接读取消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
		c.hub.broadcast <- message
	}
}

// writePump 向WebSocket连接写入消息
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.mu.Lock()
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()
			if err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}
		}
	}
}
