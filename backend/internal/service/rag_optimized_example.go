package service

import (
	"code-review-go/config"
	dto "code-review-go/internal/dto"
	"fmt"
	"time"
)

// ExampleUsage 展示如何使用优化后的RAG服务
func ExampleUsage() {
	// 1. 创建优化后的RAG客户端
	config := &RAGClientConfig{
		BaseURL:           "http://localhost:8000/api/code-analysis",
		PoolSize:          10,
		Timeout:           180 * time.Second,
		MaxRetries:        3,
		RetryDelay:        1 * time.Second,
		CacheTTL:          5 * time.Minute,
		MaxCacheSize:      1000,
		EnableCompression: true,
	}

	optimizedClient, err := NewOptimizedRAGClient(config)
	if err != nil {
		fmt.Printf("创建优化RAG客户端失败: %v\n", err)
		return
	}

	// 2. 准备请求
	req := &CodeReviewRequest{
		GitURL:      "https://gitlab.com/user/repo.git",
		Branch:      "main",
		DiffContent: "diff --git a/file.py b/file.py\n@@ -1,3 +1,3 @@\n-def old_function():\n+def new_function():\n     return True",
		Query:       "检查代码质量和安全性",
		GitlabToken: "your-gitlab-token",
	}

	// 3. 执行分析
	start := time.Now()
	result, err := optimizedClient.AnalyzeCodeWithRequestOptimized(req)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("RAG分析失败: %v\n", err)
		return
	}

	fmt.Printf("RAG分析成功，耗时: %v\n", duration)
	fmt.Printf("分析结果: %s\n", result.Review)
}

// ExampleWithManager 展示如何使用RAG服务管理器
func ExampleWithManager() {
	// 1. 获取RAG服务管理器
	manager := GetRAGServiceManager()

	// 2. 初始化管理器
	if err := manager.Initialize(); err != nil {
		fmt.Printf("初始化RAG服务管理器失败: %v\n", err)
		return
	}

	// 3. 准备webhook数据（示例）
	body := dto.WebhookBody{
		Project: dto.ProjectInfo{
			ID: 123,
		},
		ObjectAttributes: dto.ObjectAttributes{
			IID:          456,
			SourceBranch: "feature-branch",
			TargetBranch: "main",
			State:        "opened",
			URL:          "https://gitlab.com/user/repo/-/merge_requests/456",
			Description:  "这是一个测试合并请求",
		},
	}

	// 4. 执行优化后的RAG检查
	start := time.Now()
	result, err := CheckMergeRequestWithRAGOptimized(body)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("RAG检查失败: %v\n", err)
		return
	}

	fmt.Printf("RAG检查成功，耗时: %v\n", duration)
	fmt.Printf("检查结果: %s\n", result)
}

// ExamplePerformanceComparison 性能对比示例
func ExamplePerformanceComparison() {
	fmt.Println("=== RAG服务性能对比 ===")

	// 测试原始RAG客户端
	fmt.Println("1. 测试原始RAG客户端...")
	originalClient, _ := NewRAGClient("http://localhost:8000/api/code-analysis")

	req := &CodeReviewRequest{
		GitURL:      "https://gitlab.com/user/repo.git",
		Branch:      "main",
		DiffContent: "diff --git a/test.py b/test.py\n@@ -1,1 +1,1 @@\n-print('old')\n+print('new')",
		Query:       "测试查询",
		GitlabToken: "test-token",
	}

	// 测试原始客户端性能
	start := time.Now()
	_, err1 := originalClient.AnalyzeCodeWithRequest(req)
	originalDuration := time.Since(start)

	if err1 != nil {
		fmt.Printf("原始客户端测试失败: %v\n", err1)
	} else {
		fmt.Printf("原始客户端耗时: %v\n", originalDuration)
	}

	// 测试优化后的RAG客户端
	fmt.Println("2. 测试优化后的RAG客户端...")
	optimizedClient, _ := NewOptimizedRAGClient(nil)

	start = time.Now()
	_, err2 := optimizedClient.AnalyzeCodeWithRequestOptimized(req)
	optimizedDuration := time.Since(start)

	if err2 != nil {
		fmt.Printf("优化客户端测试失败: %v\n", err2)
	} else {
		fmt.Printf("优化客户端耗时: %v\n", optimizedDuration)
	}

	// 性能对比
	if err1 == nil && err2 == nil {
		improvement := float64(originalDuration-optimizedDuration) / float64(originalDuration) * 100
		fmt.Printf("性能提升: %.2f%%\n", improvement)
	}
}

// ExampleConfiguration 配置示例
func ExampleConfiguration() {
	fmt.Println("=== RAG服务配置示例 ===")

	// 1. 默认配置
	fmt.Println("1. 默认配置:")
	defaultConfig := DefaultRAGClientConfig()
	fmt.Printf("   BaseURL: %s\n", defaultConfig.BaseURL)
	fmt.Printf("   PoolSize: %d\n", defaultConfig.PoolSize)
	fmt.Printf("   Timeout: %v\n", defaultConfig.Timeout)
	fmt.Printf("   MaxRetries: %d\n", defaultConfig.MaxRetries)
	fmt.Printf("   CacheTTL: %v\n", defaultConfig.CacheTTL)

	// 2. 自定义配置
	fmt.Println("2. 自定义配置:")
	customConfig := &RAGClientConfig{
		BaseURL:           "http://rag-service:8000/api/code-analysis",
		PoolSize:          20,
		Timeout:           300 * time.Second,
		MaxRetries:        5,
		RetryDelay:        2 * time.Second,
		CacheTTL:          10 * time.Minute,
		MaxCacheSize:      2000,
		EnableCompression: true,
	}
	fmt.Printf("   BaseURL: %s\n", customConfig.BaseURL)
	fmt.Printf("   PoolSize: %d\n", customConfig.PoolSize)
	fmt.Printf("   Timeout: %v\n", customConfig.Timeout)
	fmt.Printf("   MaxRetries: %d\n", customConfig.MaxRetries)
	fmt.Printf("   CacheTTL: %v\n", customConfig.CacheTTL)

	// 3. 从环境变量加载配置
	fmt.Println("3. 从环境变量加载配置:")
	envConfig := config.LoadConfig()
	fmt.Printf("   RAGServiceURL: %s\n", envConfig.RAGServiceURL)
}

// ExampleErrorHandling 错误处理示例
func ExampleErrorHandling() {
	fmt.Println("=== RAG服务错误处理示例 ===")

	// 1. 配置错误
	fmt.Println("1. 配置错误处理:")
	invalidConfig := &RAGClientConfig{
		BaseURL: "", // 空URL
	}
	_, err := NewOptimizedRAGClient(invalidConfig)
	if err != nil {
		fmt.Printf("   配置错误: %v\n", err)
	}

	// 2. 网络错误
	fmt.Println("2. 网络错误处理:")
	networkConfig := &RAGClientConfig{
		BaseURL:    "http://invalid-url:9999",
		MaxRetries: 2,
		Timeout:    5 * time.Second,
	}
	client, _ := NewOptimizedRAGClient(networkConfig)

	req := &CodeReviewRequest{
		GitURL:      "https://gitlab.com/user/repo.git",
		Branch:      "main",
		DiffContent: "test diff",
		Query:       "test query",
	}

	_, err = client.AnalyzeCodeWithRequestOptimized(req)
	if err != nil {
		fmt.Printf("   网络错误: %v\n", err)
	}

	// 3. 超时错误
	fmt.Println("3. 超时错误处理:")
	timeoutConfig := &RAGClientConfig{
		BaseURL:    "http://localhost:8000/api/code-analysis",
		Timeout:    1 * time.Second, // 很短的超时时间
		MaxRetries: 1,
	}
	timeoutClient, _ := NewOptimizedRAGClient(timeoutConfig)

	_, err = timeoutClient.AnalyzeCodeWithRequestOptimized(req)
	if err != nil {
		fmt.Printf("   超时错误: %v\n", err)
	}
}

// ExampleMetrics 指标监控示例
func ExampleMetrics() {
	fmt.Println("=== RAG服务指标监控示例 ===")

	// 创建指标收集器
	metrics := NewRAGMetrics()

	// 模拟一些请求
	metrics.RecordRequest(true, 2*time.Second)
	metrics.RecordRequest(true, 1*time.Second)
	metrics.RecordRequest(false, 3*time.Second)
	metrics.RecordCacheHit()
	metrics.RecordCacheMiss()
	metrics.RecordCacheHit()

	// 获取统计信息
	stats := metrics.GetStats()
	fmt.Printf("总请求数: %v\n", stats["total_requests"])
	fmt.Printf("成功请求数: %v\n", stats["successful_requests"])
	fmt.Printf("失败请求数: %v\n", stats["failed_requests"])
	fmt.Printf("成功率: %.2f%%\n", stats["success_rate"])
	fmt.Printf("缓存命中数: %v\n", stats["cache_hits"])
	fmt.Printf("缓存未命中数: %v\n", stats["cache_misses"])
	fmt.Printf("缓存命中率: %.2f%%\n", stats["cache_hit_rate"])
	fmt.Printf("平均响应时间: %v\n", stats["average_response_time"])
}
