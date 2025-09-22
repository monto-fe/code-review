package service

import (
	"code-review-go/internal/cache"
	"code-review-go/internal/model"
	"code-review-go/internal/service/providers"
	"fmt"
	"sync"
)

// AIServiceImpl AI服务实现
type AIServiceImpl struct {
	config      *model.AIConfig
	mu          sync.RWMutex
	initialized bool
}

// NewAIService 创建AI服务实例
func NewAIService() AIService {
	return &AIServiceImpl{}
}

// Initialize 初始化AI服务
func (a *AIServiceImpl) Initialize() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.initialized {
		return nil
	}

	// 获取AI配置
	aiConfig, ok := cache.GetAIConfigCache()
	if !ok {
		return fmt.Errorf("初始化AI配置管理器失败")
	}
	if aiConfig.APIURL == "" {
		return fmt.Errorf("AI服务URL未配置")
	}

	a.config = aiConfig
	a.initialized = true
	return nil
}

// CallAI 调用AI服务
func (a *AIServiceImpl) CallAI(prompt string) (string, error) {
	a.mu.RLock()
	if !a.initialized {
		a.mu.RUnlock()
		return "", fmt.Errorf("AI服务未初始化")
	}
	config := a.config
	a.mu.RUnlock()

	// 根据AI配置类型选择提供者
	var provider providers.AIProvider
	switch config.Type {
	case "UCloud":
		provider = providers.NewUCAIProvider(config)
	case "DeepSeek":
		provider = providers.NewDeepSeekProvider(config)
	case "Qwen":
		provider = providers.NewQwenProvider(config)
	default:
		return "", fmt.Errorf("不支持的AI模型类型: %s", config.Type)
	}

	// 调用AI服务
	return provider.CallAI(prompt)
}

// IsInitialized 检查是否已初始化
func (a *AIServiceImpl) IsInitialized() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.initialized
}
