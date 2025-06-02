package providers

import (
	"code-review-go/internal/model"
	"fmt"
)

// AIProvider AI提供者接口
type AIProvider interface {
	CallAI(prompt string) (string, error)
}

// AIProviderFactory AI提供者工厂
type AIProviderFactory struct {
	providers map[string]func(*model.AIConfig) AIProvider
}

// NewAIProviderFactory 创建AI提供者工厂
func NewAIProviderFactory() *AIProviderFactory {
	factory := &AIProviderFactory{
		providers: make(map[string]func(*model.AIConfig) AIProvider),
	}
	// 注册OpenAI提供者
	factory.RegisterProvider("DeepSeek", func(config *model.AIConfig) AIProvider {
		return NewDeepSeekProvider(config)
	})
	// 注册UCloud提供者
	factory.RegisterProvider("UCloud", func(config *model.AIConfig) AIProvider {
		return NewUCAIProvider(config)
	})
	return factory
}

// RegisterProvider 注册AI提供者
func (f *AIProviderFactory) RegisterProvider(name string, provider func(*model.AIConfig) AIProvider) {
	f.providers[name] = provider
}

// GetProvider 获取AI提供者
func (f *AIProviderFactory) GetProvider(config *model.AIConfig) (AIProvider, error) {
	provider, ok := f.providers[config.Type]
	if !ok {
		return nil, fmt.Errorf("不支持的AI提供者: %s", config.Type)
	}
	return provider(config), nil
}
