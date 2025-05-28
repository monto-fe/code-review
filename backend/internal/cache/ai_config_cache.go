package cache

import (
	"code-review-go/internal/database"
	"code-review-go/internal/model"
	"sync"
)

var (
	aiConfigCache     *model.AIConfig
	aiConfigCacheLock sync.RWMutex
)

// InitAIConfigCache 从数据库加载所有 ai_config 信息到缓存
func InitAIConfigCache() error {
	var config model.AIConfig
	if err := database.DB.Where("is_active = ?", true).First(&config).Error; err != nil {
		aiConfigCacheLock.Lock()
		aiConfigCache = nil
		aiConfigCacheLock.Unlock()
		return nil
	}

	aiConfigCacheLock.Lock()
	aiConfigCache = &config
	aiConfigCacheLock.Unlock()
	return nil
}

// GetAIConfigCache 获取AI配置缓存
func GetAIConfigCache() (*model.AIConfig, bool) {
	aiConfigCacheLock.RLock()
	defer aiConfigCacheLock.RUnlock()
	return aiConfigCache, aiConfigCache != nil
}

// SetAIConfigCache 设置AI配置缓存
func SetAIConfigCache(config *model.AIConfig) {
	aiConfigCacheLock.Lock()
	defer aiConfigCacheLock.Unlock()
	aiConfigCache = config
}
