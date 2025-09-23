package service

import (
	"sync"
)

// DatabaseServiceManager 数据库服务管理器
type DatabaseServiceManager struct {
	service     DatabaseService
	mu          sync.RWMutex
	initialized bool
}

// NewDatabaseServiceManager 创建数据库服务管理器
func NewDatabaseServiceManager() *DatabaseServiceManager {
	return &DatabaseServiceManager{
		service: NewDatabaseService(),
	}
}

// Initialize 初始化数据库服务
func (m *DatabaseServiceManager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return nil
	}

	if err := m.service.Initialize(); err != nil {
		return err
	}

	m.initialized = true
	return nil
}

// GetService 获取数据库服务
func (m *DatabaseServiceManager) GetService() DatabaseService {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.service
}

// IsInitialized 检查是否已初始化
func (m *DatabaseServiceManager) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.initialized
}

// RAGServiceManagerImpl RAG服务管理器实现
type RAGServiceManagerImpl struct {
	service     RAGServiceInterface
	mu          sync.RWMutex
	initialized bool
}

// NewRAGServiceManagerImpl 创建RAG服务管理器实现
func NewRAGServiceManagerImpl() *RAGServiceManagerImpl {
	return &RAGServiceManagerImpl{
		service: NewRAGService(),
	}
}

// Initialize 初始化RAG服务
func (m *RAGServiceManagerImpl) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return nil
	}

	if err := m.service.Initialize(); err != nil {
		return err
	}

	m.initialized = true
	return nil
}

// GetService 获取RAG服务
func (m *RAGServiceManagerImpl) GetService() RAGServiceInterface {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.service
}

// IsInitialized 检查是否已初始化
func (m *RAGServiceManagerImpl) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.initialized
}

// AIServiceManager AI服务管理器
type AIServiceManager struct {
	service     AIService
	mu          sync.RWMutex
	initialized bool
}

// NewAIServiceManager 创建AI服务管理器
func NewAIServiceManager() *AIServiceManager {
	return &AIServiceManager{
		service: NewAIService(),
	}
}

// Initialize 初始化AI服务
func (m *AIServiceManager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return nil
	}

	if err := m.service.Initialize(); err != nil {
		return err
	}

	m.initialized = true
	return nil
}

// GetService 获取AI服务
func (m *AIServiceManager) GetService() AIService {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.service
}

// IsInitialized 检查是否已初始化
func (m *AIServiceManager) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.initialized
}

// NotificationServiceManager 通知服务管理器
type NotificationServiceManager struct {
	service     NotificationService
	mu          sync.RWMutex
	initialized bool
}

// NewNotificationServiceManager 创建通知服务管理器
func NewNotificationServiceManager() *NotificationServiceManager {
	return &NotificationServiceManager{
		service: NewNotificationService(),
	}
}

// Initialize 初始化通知服务
func (m *NotificationServiceManager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return nil
	}

	if err := m.service.Initialize(); err != nil {
		return err
	}

	m.initialized = true
	return nil
}

// GetService 获取通知服务
func (m *NotificationServiceManager) GetService() NotificationService {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.service
}

// IsInitialized 检查是否已初始化
func (m *NotificationServiceManager) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.initialized
}
