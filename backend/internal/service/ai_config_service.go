package service

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"code-review-go/internal/cache"
	"code-review-go/internal/model"
)

// AIConfigManager AI配置管理器
type AIConfigManager struct {
	config *model.AIConfigInfo
	db     *gorm.DB
	mu     sync.RWMutex
}

var (
	instance *AIConfigManager
	once     sync.Once
)

// NewAIConfigManager 创建AI配置管理器实例
func NewAIConfigManager(db *gorm.DB) (*AIConfigManager, error) {
	var initErr error
	once.Do(func() {
		instance = &AIConfigManager{
			db: db,
		}
		if err := instance.loadConfig(); err != nil {
			initErr = err
			return
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}

// loadConfig 加载配置
func (m *AIConfigManager) loadConfig() error {
	var config model.AIConfig
	if err := m.db.Where("is_active = ?", true).
		Order("update_time DESC").
		First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("No active AI config found, using fallback or skipping initialization.")
			return nil
		}
		return err
	}

	m.mu.Lock()
	m.config = &model.AIConfigInfo{
		Name:   config.Name,
		APIKey: config.APIKey,
		APIURL: config.APIURL,
		Model:  config.Model,
	}
	m.mu.Unlock()

	return nil
}

// startAutoRefresh 启动自动刷新
func (m *AIConfigManager) startAutoRefresh() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.loadConfig(); err != nil {
			fmt.Printf("Failed to refresh AI config: %v\n", err)
		}
	}
}

// RefreshCache 刷新缓存
func (m *AIConfigManager) RefreshCache() error {
	var config model.AIConfig
	if err := m.db.Where("is_active = ?", true).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cache.SetAIConfigCache(nil)
			return nil
		}
		return err
	}

	// 更新缓存
	cache.SetAIConfigCache(&config)
	return nil
}

// AddAIConfig 添加AI配置
func (m *AIConfigManager) AddAIConfig(data model.AIConfigCreate) (*model.AIConfig, error) {
	now := time.Now().Unix()
	config := &model.AIConfig{
		Name:       data.Name,
		APIURL:     data.APIURL,
		APIKey:     data.APIKey,
		Model:      data.Model,
		IsActive:   data.IsActive,
		CreateTime: now,
		UpdateTime: now,
	}

	if err := m.db.Create(config).Error; err != nil {
		return nil, err
	}

	// 重新加载配置到内存
	if err := m.loadConfig(); err != nil {
		fmt.Printf("Failed to reload config after adding: %v\n", err)
	}

	return config, nil
}

// UpdateAIConfig 更新AI配置
func (m *AIConfigManager) UpdateAIConfig(data model.AIConfigUpdate) error {
	now := time.Now().Unix()
	if err := m.db.Model(&model.AIConfig{}).
		Where("id = ?", data.ID).
		Updates(map[string]interface{}{
			"name":        data.Name,
			"api_url":     data.APIURL,
			"api_key":     data.APIKey,
			"model":       data.Model,
			"is_active":   data.IsActive,
			"update_time": now,
		}).Error; err != nil {
		return err
	}

	// 重新加载配置到内存
	if err := m.loadConfig(); err != nil {
		fmt.Printf("Failed to reload config after updating: %v\n", err)
	}

	return nil
}

// GetConfigList 获取配置列表
func (m *AIConfigManager) GetConfigList() ([]model.AIConfig, error) {
	var configs []model.AIConfig
	if err := m.db.Order("id DESC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetConfig 获取当前配置
func (m *AIConfigManager) GetConfig() (*model.AIConfigInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.config == nil {
		return nil, fmt.Errorf("AI config not loaded yet")
	}
	return m.config, nil
}

// DeleteAIConfig 删除AI配置
func (m *AIConfigManager) DeleteAIConfig(id uint) error {
	if err := m.db.Delete(&model.AIConfig{}, id).Error; err != nil {
		return err
	}

	// 重新加载配置到内存
	if err := m.loadConfig(); err != nil {
		fmt.Printf("Failed to reload config after deleting: %v\n", err)
	}

	return nil
}
