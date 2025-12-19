package gitlab_service

import (
	"bytes"
	dto "code-review-go/internal/dto"
	"code-review-go/internal/service/common"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// RAGClient RAG服务客户端
type RAGClient struct {
	baseURL    string
	client     *http.Client
	clientPool *common.RAGClientPool
	cache      *RAGResponseCache
	config     *RAGClientConfig
	metrics    *RAGMetrics
}

// RAGClientConfig RAG客户端配置
type RAGClientConfig struct {
	BaseURL            string
	PoolSize           int
	Timeout            time.Duration
	MaxRetries         int
	RetryDelay         time.Duration
	CacheTTL           time.Duration
	MaxCacheSize       int
	EnableCompression  bool
	EnableOptimization bool
}

// DefaultRAGClientConfig 默认配置
func DefaultRAGClientConfig() *RAGClientConfig {
	return &RAGClientConfig{
		BaseURL:            "http://localhost:8000/api/code-analysis",
		PoolSize:           10,
		Timeout:            180 * time.Second,
		MaxRetries:         2,
		RetryDelay:         1 * time.Second,
		CacheTTL:           5 * time.Minute,
		MaxCacheSize:       1000,
		EnableCompression:  true,
		EnableOptimization: false, // 默认关闭优化功能
	}
}

// NewRAGClient 创建RAG服务客户端
func NewRAGClient(baseURL string) *RAGClient {
	config := DefaultRAGClientConfig()
	config.BaseURL = baseURL

	client := &RAGClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		config:  config,
		metrics: NewRAGMetrics(),
	}

	// 如果启用优化功能，初始化连接池和缓存
	if config.EnableOptimization {
		clientPool, err := common.NewRAGClientPool(baseURL, config.PoolSize)
		if err == nil {
			client.clientPool = clientPool
		}
		client.cache = NewRAGResponseCache(config.CacheTTL, config.MaxCacheSize)
	}

	return client
}

// NewOptimizedRAGClient 创建优化的RAG客户端
func NewOptimizedRAGClient(config *RAGClientConfig) (*RAGClient, error) {
	if config == nil {
		config = DefaultRAGClientConfig()
	}

	config.EnableOptimization = true

	// 创建连接池
	clientPool, err := common.NewRAGClientPool(config.BaseURL, config.PoolSize)
	if err != nil {
		return nil, fmt.Errorf("创建连接池失败: %v", err)
	}

	// 创建缓存
	cache := NewRAGResponseCache(config.CacheTTL, config.MaxCacheSize)

	return &RAGClient{
		baseURL:    config.BaseURL,
		client:     &http.Client{Timeout: config.Timeout},
		clientPool: clientPool,
		cache:      cache,
		config:     config,
		metrics:    NewRAGMetrics(),
	}, nil
}

// Initialize 初始化RAG客户端（实现RAGServiceInterface）
func (c *RAGClient) Initialize() error {
	// RAG客户端不需要特殊初始化
	return nil
}

// IsInitialized 检查是否已初始化（实现RAGServiceInterface）
func (c *RAGClient) IsInitialized() bool {
	return true
}

// AnalyzeCodeWithRequest 使用请求对象进行分析
func (c *RAGClient) AnalyzeCodeWithRequest(req *dto.CodeReviewRequest) (*dto.CodeAnalysisResponse, error) {
	startTime := time.Now()

	// 如果启用优化功能，使用缓存和连接池
	if c.config != nil && c.config.EnableOptimization {
		return c.analyzeWithOptimization(req, startTime)
	}

	// 使用基础实现
	return c.analyzeBasic(req, startTime)
}

// analyzeWithOptimization 使用优化功能进行分析
func (c *RAGClient) analyzeWithOptimization(req *dto.CodeReviewRequest, startTime time.Time) (*dto.CodeAnalysisResponse, error) {
	// 生成缓存键
	cacheKey := c.generateCacheKey(req)

	// 尝试从缓存获取
	if cached, found := c.cache.Get(cacheKey); found {
		c.metrics.RecordCacheHit()
		return cached, nil
	}

	c.metrics.RecordCacheMiss()

	// 执行带重试的请求
	var result *dto.CodeAnalysisResponse
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
		c.metrics.RecordRequest(false, time.Since(startTime))
		return nil, err
	}

	// 缓存结果
	c.cache.Set(cacheKey, result)
	c.metrics.RecordRequest(true, time.Since(startTime))

	return result, nil
}

// analyzeBasic 基础分析实现
func (c *RAGClient) analyzeBasic(req *dto.CodeReviewRequest, startTime time.Time) (*dto.CodeAnalysisResponse, error) {
	// 创建超时的上下文
	timeout := 180 * time.Second
	if c.config != nil {
		timeout = c.config.Timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建带超时的请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("RAG服务请求超时(%v): %v", timeout, err)
		}
		return nil, fmt.Errorf("请求RAG服务失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RAG服务返回错误状态码: %d", resp.StatusCode)
	}

	var result dto.CodeAnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析RAG服务响应失败: %v", err)
	}

	if c.metrics != nil {
		c.metrics.RecordRequest(true, time.Since(startTime))
	}

	return &result, nil
}

// performRequest 执行单个请求（优化版本）
func (c *RAGClient) performRequest(req *dto.CodeReviewRequest) (*dto.CodeAnalysisResponse, error) {
	// 从连接池获取客户端（如果可用）
	var httpClient *http.Client
	if c.clientPool != nil {
		poolClient := c.clientPool.GetClient()
		defer c.clientPool.ReturnClient(poolClient)
		// 暂时使用默认客户端，连接池功能待完善
		httpClient = c.client
	} else {
		httpClient = c.client
	}

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
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("RAG服务请求超时(%v): %v", c.config.Timeout, err)
		}
		return nil, fmt.Errorf("请求RAG服务失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		all, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("RAG服务返回错误 code: %d, info: %s", resp.StatusCode, string(all))
	}

	// 解析响应
	var result dto.CodeAnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析RAG服务响应失败: %v", err)
	}

	return &result, nil
}

// generateCacheKey 生成缓存键
func (c *RAGClient) generateCacheKey(req *dto.CodeReviewRequest) string {
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

// GenerateReview 根据分析结果生成代码审查
func (c *RAGClient) GenerateReview(req *dto.CodeReviewRequest) (string, error) {
	analysis, err := c.AnalyzeCodeWithRequest(req)
	if err != nil {
		return "", err
	}

	// 直接返回RAG服务的审查结果
	return analysis.Review, nil
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
	response  *dto.CodeAnalysisResponse
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
func (c *RAGResponseCache) Get(key string) (*dto.CodeAnalysisResponse, bool) {
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
func (c *RAGResponseCache) Set(key string, response *dto.CodeAnalysisResponse) {
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
