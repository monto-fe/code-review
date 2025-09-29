package handlers

import (
	"fmt"
	"sync"

	"code-review-go/internal/dto"
)

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(body dto.WebhookBody) error
	GetEventType() string
}

// EventRouter 事件路由器
type EventRouter struct {
	handlers map[string]EventHandler
	mu       sync.RWMutex
}

// NewEventRouter 创建新的事件路由器
func NewEventRouter() *EventRouter {
	return &EventRouter{
		handlers: make(map[string]EventHandler),
	}
}

// RegisterHandler 注册事件处理器
func (r *EventRouter) RegisterHandler(handler EventHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	eventType := handler.GetEventType()
	r.handlers[eventType] = handler
	fmt.Printf("注册事件处理器: %s\n", eventType)
}

// Route 路由事件到对应的处理器
func (r *EventRouter) Route(body dto.WebhookBody) error {
	// 1. 检测事件类型
	eventType := detectEventType(body)

	// 2. 获取对应的处理器
	r.mu.RLock()
	handler, exists := r.handlers[eventType]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("不支持的事件类型: %s", eventType)
	}

	// 3. 异步处理事件
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("事件处理发生panic: %v\n", r)
			}
		}()

		err := handler.Handle(body)
		if err != nil {
			fmt.Printf("事件处理失败 [%s]: %v\n", eventType, err)
		} else {
			fmt.Printf("事件处理成功 [%s]\n", eventType)
		}
	}()

	return nil
}

// detectEventType 检测事件类型
func detectEventType(body dto.WebhookBody) string {
	// 1. 检查 object_kind 字段
	if body.ObjectKind == "merge_request" {
		return "merge_request"
	}
	if body.ObjectKind == "push" {
		return "push"
	}

	// 2. 检查是否有 Merge Request 相关字段
	if body.ObjectAttributes.IID > 0 && body.ObjectAttributes.Title != "" {
		return "merge_request"
	}

	// 3. 检查是否有 Push 相关字段
	if body.Ref != "" && body.After != "" {
		return "push"
	}

	// 4. 其他事件类型暂不处理，返回 unknown
	return "unknown"
}

// GetRegisteredHandlers 获取已注册的处理器列表
func (r *EventRouter) GetRegisteredHandlers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handlers := make([]string, 0, len(r.handlers))
	for eventType := range r.handlers {
		handlers = append(handlers, eventType)
	}
	return handlers
}
