package service

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"code-review-go/internal/cache"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
)

const (
	AIConfigActive   int8 = 1 // 启用
	AIConfigInactive int8 = 2 // 禁用
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
	if err := m.db.Where("is_active = ?", AIConfigActive).
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

// RefreshCache 刷新缓存
func (m *AIConfigManager) RefreshCache() error {
	var config model.AIConfig
	if err := m.db.Where("is_active = ?", AIConfigActive).First(&config).Error; err != nil {
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
// 添加AI配置时，默认不启用，不需要更新缓存
func (m *AIConfigManager) AddAIConfig(data model.AIConfigCreate) (*model.AIConfig, error) {
	now := time.Now().Unix()

	// 如果本次要激活该模型，先关闭其他已激活配置
	if data.IsActive == AIConfigActive {
		if err := m.db.Model(&model.AIConfig{}).Where("is_active = ?", AIConfigActive).Update("is_active", AIConfigInactive).Error; err != nil {
			return nil, err
		}
	}

	config := &model.AIConfig{
		Name:       data.Name,
		APIURL:     data.APIURL,
		APIKey:     data.APIKey,
		Model:      data.Model,
		Type:       data.Type,
		IsActive:   data.IsActive,
		CreateTime: now,
		UpdateTime: now,
	}

	if err := m.db.Create(config).Error; err != nil {
		return nil, err
	}

	// 如果新加的是激活状态，重新加载缓存
	if data.IsActive == AIConfigActive {
		if err := m.loadConfig(); err != nil {
			fmt.Printf("Failed to reload config after adding: %v\n", err)
		}
	}

	return config, nil
}

// UpdateAIConfig 更新AI配置
// 更新AI配置时，如果启用了当前model，则需要更新缓存
func (m *AIConfigManager) UpdateAIConfig(data model.AIConfigUpdate) error {
	now := time.Now().Unix()

	// 查找当前记录
	var current model.AIConfig
	if err := m.db.First(&current, data.ID).Error; err != nil {
		return err
	}

	// 如果本次要激活该模型
	if data.IsActive == AIConfigActive {
		// 查找当前已激活的模型
		var active model.AIConfig
		err := m.db.Where("is_active = ?", AIConfigActive).First(&active).Error
		if err == nil && active.ID != data.ID {
			// 取消当前激活
			if err := m.db.Model(&model.AIConfig{}).Where("id = ?", active.ID).Update("is_active", AIConfigInactive).Error; err != nil {
				return err
			}
		}
	}

	// 构造更新字段
	updateMap := map[string]interface{}{
		"name":        chooseString(data.Name, current.Name),
		"api_url":     chooseString(data.APIURL, current.APIURL),
		"api_key":     chooseString(data.APIKey, current.APIKey),
		"model":       chooseString(data.Model, current.Model),
		"type":        chooseString(data.Type, current.Type),
		"is_active":   chooseInt8(data.IsActive, current.IsActive),
		"update_time": now,
	}

	if err := m.db.Model(&model.AIConfig{}).
		Where("id = ?", data.ID).
		Updates(updateMap).Error; err != nil {
		return err
	}

	// 重新加载配置到内存
	if err := m.loadConfig(); err != nil {
		fmt.Printf("Failed to reload config after updating: %v\n", err)
	}

	return nil
}

// chooseString returns a if a is not empty, otherwise b
func chooseString(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// chooseInt8 returns a if a != 0, otherwise b
func chooseInt8(a, b int8) int8 {
	if a != 0 {
		return a
	}
	return b
}

// GetConfigList 获取配置列表
func (m *AIConfigManager) GetConfigListPaged(page, pageSize int) ([]model.AIConfig, int64, error) {
	var list []model.AIConfig
	var total int64
	offset := (page - 1) * pageSize
	db := m.db.Model(&model.AIConfig{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	for i := range list {
		list[i].APIKey = utils.MaskString(list[i].APIKey)
	}
	return list, total, nil
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
	var config model.AIConfig
	if err := m.db.First(&config, id).Error; err != nil {
		return err
	}
	if config.IsActive == AIConfigActive {
		return fmt.Errorf("请先禁用该配置后再删除")
	}

	if err := m.db.Delete(&model.AIConfig{}, id).Error; err != nil {
		return err
	}

	// 重新加载配置到内存
	if err := m.loadConfig(); err != nil {
		fmt.Printf("Failed to reload config after deleting: %v\n", err)
	}

	return nil
}
