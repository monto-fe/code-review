package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// OptimizedRAGClient 优化后的RAG客户端
type OptimizedRAGClient struct {
	baseURL    string
	clientPool *RAGClientPool
	cache      *RAGResponseCache
	config     *RAGClientConfig
}

// RAGClientConfig RAG客户端配置
type RAGClientConfig struct {
	BaseURL           string
	PoolSize          int
	Timeout           time.Duration
	MaxRetries        int
	RetryDelay        time.Duration
	CacheTTL          time.Duration
	MaxCacheSize      int
	EnableCompression bool
}

// DefaultRAGClientConfig 默认配置
func DefaultRAGClientConfig() *RAGClientConfig {
	return &RAGClientConfig{
		BaseURL:           "http://localhost:8000/api/code-analysis",
		PoolSize:          10,
		Timeout:           180 * time.Second,
		MaxRetries:        2,
		RetryDelay:        1 * time.Second,
		CacheTTL:          5 * time.Minute,
		MaxCacheSize:      1000,
		EnableCompression: true,
	}
}

// NewOptimizedRAGClient 创建优化后的RAG客户端
func NewOptimizedRAGClient(config *RAGClientConfig) (*OptimizedRAGClient, error) {
	if config == nil {
		config = DefaultRAGClientConfig()
	}

	// 创建连接池
	clientPool, err := NewRAGClientPool(config.BaseURL, config.PoolSize)
	if err != nil {
		return nil, fmt.Errorf("创建连接池失败: %v", err)
	}

	// 创建缓存
	cache := NewRAGResponseCache(config.CacheTTL, config.MaxCacheSize)

	return &OptimizedRAGClient{
		baseURL:    config.BaseURL,
		clientPool: clientPool,
		cache:      cache,
		config:     config,
	}, nil
}

// AnalyzeCodeWithRequestOptimized 优化后的代码分析请求
func (c *OptimizedRAGClient) AnalyzeCodeWithRequestOptimized(req *CodeReviewRequest) (*CodeAnalysisResponse, error) {
	// 生成缓存键
	cacheKey := c.generateCacheKey(req)

	// 尝试从缓存获取
	if cached, found := c.cache.Get(cacheKey); found {
		return cached, nil
	}

	// 执行带重试的请求
	var result *CodeAnalysisResponse
	var err error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		result, err = c.performRequest(req)
		if err == nil {
			break
		}

		// 如果不是最后一次尝试，等待后重试
		if attempt < c.config.MaxRetries {
			time.Sleep(c.config.RetryDelay * time.Duration(attempt+1))
		}
	}

	if err != nil {
		return nil, err
	}

	// 缓存结果
	c.cache.Set(cacheKey, result)

	return result, nil
}

// performRequest 执行单个请求
func (c *OptimizedRAGClient) performRequest(req *CodeReviewRequest) (*CodeAnalysisResponse, error) {
	// 从连接池获取客户端
	client := c.clientPool.GetClient()
	defer c.clientPool.ReturnClient(client)

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	// 序列化请求
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	if c.config.EnableCompression {
		httpReq.Header.Set("Accept-Encoding", "gzip, deflate")
	}

	// 执行请求
	resp, err := client.client.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("RAG服务请求超时(%v): %v", c.config.Timeout, err)
		}
		return nil, fmt.Errorf("请求RAG服务失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RAG服务返回错误状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var result CodeAnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析RAG服务响应失败: %v", err)
	}

	return &result, nil
}

// generateCacheKey 生成缓存键
func (c *OptimizedRAGClient) generateCacheKey(req *CodeReviewRequest) string {
	// 使用请求的关键信息生成缓存键
	keyData := map[string]string{
		"git_url":      req.GitURL,
		"branch":       req.Branch,
		"diff_content": req.DiffContent,
		"query":        req.Query,
	}

	keyBytes, _ := json.Marshal(keyData)
	return string(keyBytes)
}

// RAGResponseCache RAG响应缓存
type RAGResponseCache struct {
	cache         map[string]*cacheEntry
	ttl           time.Duration
	maxSize       int
	mu            sync.RWMutex
	cleanupTicker *time.Ticker
	stopChan      chan struct{}
}

// cacheEntry 缓存条目
type cacheEntry struct {
	response  *CodeAnalysisResponse
	timestamp time.Time
}

// NewRAGResponseCache 创建RAG响应缓存
func NewRAGResponseCache(ttl time.Duration, maxSize int) *RAGResponseCache {
	cache := &RAGResponseCache{
		cache:         make(map[string]*cacheEntry),
		ttl:           ttl,
		maxSize:       maxSize,
		cleanupTicker: time.NewTicker(ttl / 2), // 每TTL/2时间清理一次
		stopChan:      make(chan struct{}),
	}

	// 启动清理协程
	go cache.cleanupRoutine()

	return cache
}

// Get 获取缓存值
func (c *RAGResponseCache) Get(key string) (*CodeAnalysisResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// 检查是否过期
	if time.Since(entry.timestamp) > c.ttl {
		// 标记为过期，下次清理时删除
		return nil, false
	}

	return entry.response, true
}

// Set 设置缓存值
func (c *RAGResponseCache) Set(key string, response *CodeAnalysisResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 检查缓存大小
	if len(c.cache) >= c.maxSize {
		c.evictOldest()
	}

	c.cache[key] = &cacheEntry{
		response:  response,
		timestamp: time.Now(),
	}
}

// evictOldest 淘汰最旧的条目
func (c *RAGResponseCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range c.cache {
		if oldestKey == "" || entry.timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.timestamp
		}
	}

	if oldestKey != "" {
		delete(c.cache, oldestKey)
	}
}

// cleanupRoutine 清理协程
func (c *RAGResponseCache) cleanupRoutine() {
	for {
		select {
		case <-c.cleanupTicker.C:
			c.cleanup()
		case <-c.stopChan:
			c.cleanupTicker.Stop()
			return
		}
	}
}

// cleanup 清理过期条目
func (c *RAGResponseCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.cache {
		if now.Sub(entry.timestamp) > c.ttl {
			delete(c.cache, key)
		}
	}
}

// Close 关闭缓存
func (c *RAGResponseCache) Close() {
	close(c.stopChan)
}

// RAGMetrics RAG服务指标
type RAGMetrics struct {
	TotalRequests       int64
	SuccessfulRequests  int64
	FailedRequests      int64
	CacheHits           int64
	CacheMisses         int64
	AverageResponseTime time.Duration
	mu                  sync.RWMutex
}

// NewRAGMetrics 创建RAG指标
func NewRAGMetrics() *RAGMetrics {
	return &RAGMetrics{}
}

// RecordRequest 记录请求
func (m *RAGMetrics) RecordRequest(success bool, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalRequests++
	if success {
		m.SuccessfulRequests++
	} else {
		m.FailedRequests++
	}

	// 更新平均响应时间
	if m.SuccessfulRequests > 0 {
		totalTime := m.AverageResponseTime * time.Duration(m.SuccessfulRequests-1)
		m.AverageResponseTime = (totalTime + duration) / time.Duration(m.SuccessfulRequests)
	}
}

// RecordCacheHit 记录缓存命中
func (m *RAGMetrics) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheHits++
}

// RecordCacheMiss 记录缓存未命中
func (m *RAGMetrics) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CacheMisses++
}

// GetStats 获取统计信息
func (m *RAGMetrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	successRate := float64(0)
	if m.TotalRequests > 0 {
		successRate = float64(m.SuccessfulRequests) / float64(m.TotalRequests) * 100
	}

	cacheHitRate := float64(0)
	totalCacheRequests := m.CacheHits + m.CacheMisses
	if totalCacheRequests > 0 {
		cacheHitRate = float64(m.CacheHits) / float64(totalCacheRequests) * 100
	}

	return map[string]interface{}{
		"total_requests":        m.TotalRequests,
		"successful_requests":   m.SuccessfulRequests,
		"failed_requests":       m.FailedRequests,
		"success_rate":          successRate,
		"cache_hits":            m.CacheHits,
		"cache_misses":          m.CacheMisses,
		"cache_hit_rate":        cacheHitRate,
		"average_response_time": m.AverageResponseTime.String(),
	}
}
