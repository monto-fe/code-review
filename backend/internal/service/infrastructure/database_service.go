package infrastructure

import (
	"code-review-go/internal/database"
	"code-review-go/internal/service/ai"
	"code-review-go/internal/service/common"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

// DatabaseServiceImpl 数据库服务实现
type DatabaseServiceImpl struct {
	db            *gorm.DB
	aiRuleService *ai.AIRuleService
	mu            sync.RWMutex
	initialized   bool
}

// NewDatabaseService 创建数据库服务实例
func NewDatabaseService() common.DatabaseService {
	return &DatabaseServiceImpl{}
}

// Initialize 初始化数据库服务
func (d *DatabaseServiceImpl) Initialize() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.initialized {
		return nil
	}

	// 获取数据库连接
	d.db = database.GetDB()
	if d.db == nil {
		return fmt.Errorf("获取数据库连接失败")
	}

	// 初始化AI规则服务
	d.aiRuleService = ai.NewAIRuleService(d.db)

	d.initialized = true
	return nil
}

// GetDB 获取数据库连接
func (d *DatabaseServiceImpl) GetDB() *gorm.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db
}

// GetAIRuleService 获取AI规则服务
func (d *DatabaseServiceImpl) GetAIRuleService() interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.aiRuleService
}

// IsInitialized 检查是否已初始化
func (d *DatabaseServiceImpl) IsInitialized() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.initialized
}
