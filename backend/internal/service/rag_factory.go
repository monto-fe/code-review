package service

import (
	"time"
)

// CreateRAGClient 创建RAG客户端工厂函数
func CreateRAGClient(baseURL string) (RAGServiceInterface, error) {
	config := &RAGClientConfig{
		BaseURL:            baseURL,
		PoolSize:           10,
		Timeout:            180 * time.Second,
		MaxRetries:         3,
		RetryDelay:         1 * time.Second,
		CacheTTL:           5 * time.Minute,
		MaxCacheSize:       1000,
		EnableCompression:  true,
		EnableOptimization: true,
	}
	return NewOptimizedRAGClient(config)
}

// CreateDefaultRAGClient 创建默认的RAG客户端（优化版本）
func CreateDefaultRAGClient(baseURL string) (RAGServiceInterface, error) {
	return CreateRAGClient(baseURL)
}
