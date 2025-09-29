package manager

import (
	"code-review-go/internal/service/common"
	"fmt"
	"sync"
)

// ServiceCoordinator 服务协调器，负责协调各个服务
type ServiceCoordinator struct {
	database     *DatabaseServiceManager
	rag          *RAGServiceManagerImpl
	ai           *AIServiceManager
	notification *NotificationServiceManager
	mu           sync.RWMutex
	initialized  bool
}

// NewServiceCoordinator 创建服务协调器
func NewServiceCoordinator() *ServiceCoordinator {
	return &ServiceCoordinator{
		database:     NewDatabaseServiceManager(),
		rag:          NewRAGServiceManagerImpl(),
		ai:           NewAIServiceManager(),
		notification: NewNotificationServiceManager(),
	}
}

// NewServiceCoordinatorWithManagers 使用指定的管理器创建协调器（用于测试）
func NewServiceCoordinatorWithManagers(
	database *DatabaseServiceManager,
	rag *RAGServiceManagerImpl,
	ai *AIServiceManager,
	notification *NotificationServiceManager,
) *ServiceCoordinator {
	return &ServiceCoordinator{
		database:     database,
		rag:          rag,
		ai:           ai,
		notification: notification,
	}
}

// Initialize 初始化所有服务
func (c *ServiceCoordinator) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return nil
	}

	// 初始化各个服务
	if err := c.database.Initialize(); err != nil {
		return fmt.Errorf("初始化数据库服务失败: %v", err)
	}

	if err := c.rag.Initialize(); err != nil {
		return fmt.Errorf("初始化RAG服务失败: %v", err)
	}

	if err := c.ai.Initialize(); err != nil {
		return fmt.Errorf("初始化AI服务失败: %v", err)
	}

	if err := c.notification.Initialize(); err != nil {
		return fmt.Errorf("初始化通知服务失败: %v", err)
	}

	c.initialized = true
	return nil
}

// GetDatabaseManager 获取数据库服务管理器
func (c *ServiceCoordinator) GetDatabaseManager() *DatabaseServiceManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.database
}

// GetRAGManager 获取RAG服务管理器
func (c *ServiceCoordinator) GetRAGManager() *RAGServiceManagerImpl {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rag
}

// GetAIManager 获取AI服务管理器
func (c *ServiceCoordinator) GetAIManager() *AIServiceManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ai
}

// GetNotificationManager 获取通知服务管理器
func (c *ServiceCoordinator) GetNotificationManager() *NotificationServiceManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.notification
}

// IsInitialized 检查是否已初始化
func (c *ServiceCoordinator) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}

// GetDatabaseService 获取数据库服务
func (c *ServiceCoordinator) GetDatabaseService() common.DatabaseService {
	return c.GetDatabaseManager().GetService()
}

// GetRAGService 获取RAG服务
func (c *ServiceCoordinator) GetRAGService() common.RAGServiceInterface {
	return c.GetRAGManager().GetService()
}

// GetAIService 获取AI服务
func (c *ServiceCoordinator) GetAIService() common.AIService {
	return c.GetAIManager().GetService()
}

// GetNotificationService 获取通知服务
func (c *ServiceCoordinator) GetNotificationService() common.NotificationService {
	return c.GetNotificationManager().GetService()
}
