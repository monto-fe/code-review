package service

import (
	"fmt"
	"time"
)

// RAGClientType RAG客户端类型
type RAGClientType string

const (
	RAGClientTypeBasic     RAGClientType = "basic"
	RAGClientTypeOptimized RAGClientType = "optimized"
)

// CreateRAGClient 创建RAG客户端工厂函数
func CreateRAGClient(clientType RAGClientType, baseURL string) (RAGService, error) {
	switch clientType {
	case RAGClientTypeBasic:
		return NewRAGClient(baseURL), nil
	case RAGClientTypeOptimized:
		config := &RAGClientConfig{
			BaseURL:           baseURL,
			PoolSize:          10,
			Timeout:           180 * time.Second,
			MaxRetries:        3,
			RetryDelay:        1 * time.Second,
			CacheTTL:          5 * time.Minute,
			MaxCacheSize:      1000,
			EnableCompression: true,
		}
		return NewOptimizedRAGClient(config)
	default:
		return nil, fmt.Errorf("不支持的RAG客户端类型: %s", clientType)
	}
}

// CreateDefaultRAGClient 创建默认的RAG客户端（优化版本）
func CreateDefaultRAGClient(baseURL string) (RAGService, error) {
	return CreateRAGClient(RAGClientTypeOptimized, baseURL)
}
