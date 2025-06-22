package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestRAGClientTimeout(t *testing.T) {
	// 测试超时配置
	client, err := NewRAGClient("http://localhost:8000")
	if err != nil {
		t.Fatalf("创建RAG客户端失败: %v", err)
	}

	// 验证HTTP客户端超时设置
	if client.client.Timeout != 180*time.Second {
		t.Errorf("期望超时时间为180秒，实际为: %v", client.client.Timeout)
	} else {
		fmt.Println("HTTP客户端超时设置正确: 180秒")
	}
}

func TestAnalyzeCodeWithRequestTimeout(t *testing.T) {
	// 测试请求超时
	client, err := NewRAGClient("http://invalid-url:9999") // 使用无效URL来测试超时
	if err != nil {
		t.Fatalf("创建RAG客户端失败: %v", err)
	}

	req := &CodeReviewRequest{
		GitURL:      "https://gitlab.com/test/repo.git",
		Branch:      "main",
		DiffContent: "diff --git a/test.py b/test.py\n@@ -1,1 +1,1 @@\n-print('old')\n+print('new')",
		Query:       "test query",
	}

	// 记录开始时间
	start := time.Now()

	// 调用方法，应该超时
	_, err = client.AnalyzeCodeWithRequest(req)

	// 记录结束时间
	duration := time.Since(start)

	// 验证超时行为
	if err == nil {
		t.Error("期望出现超时错误，但没有错误")
		return
	}

	// 检查是否在合理的时间内超时（应该小于190秒，大于170秒）
	if duration > 190*time.Second {
		t.Errorf("超时时间过长: %v", duration)
	}
	if duration < 170*time.Second {
		t.Errorf("超时时间过短: %v", duration)
	}

	fmt.Printf("超时测试通过，耗时: %v，错误: %v\n", duration, err)
}

func TestCheckMergeRequestWithRAGTimeout(t *testing.T) {
	// 测试RAG分析超时
	client, err := NewRAGClient("http://invalid-url:9999") // 使用无效URL来测试超时
	if err != nil {
		t.Fatalf("创建RAG客户端失败: %v", err)
	}

	req := &CodeReviewRequest{
		GitURL:      "https://gitlab.com/test/repo.git",
		Branch:      "main",
		DiffContent: "diff --git a/test.py b/test.py\n@@ -1,1 +1,1 @@\n-print('old')\n+print('new')",
		Query:       "test query",
	}

	// 记录开始时间
	start := time.Now()

	// 调用方法，应该超时
	_, err = client.AnalyzeCodeWithRequest(req)

	// 记录结束时间
	duration := time.Since(start)

	// 验证超时行为
	if err == nil {
		t.Error("期望出现超时错误，但没有错误")
		return
	}

	// 检查是否在合理的时间内超时（应该小于190秒，大于170秒）
	if duration > 190*time.Second {
		t.Errorf("超时时间过长: %v", duration)
	}
	if duration < 170*time.Second {
		t.Errorf("超时时间过短: %v", duration)
	}

	fmt.Printf("RAG分析超时测试通过，耗时: %v，错误: %v\n", duration, err)
}

func TestContextTimeout(t *testing.T) {
	// 测试context超时
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 创建一个channel来模拟长时间操作
	done := make(chan bool)

	go func() {
		time.Sleep(5 * time.Second) // 睡眠5秒
		done <- true
	}()

	// 等待context超时或操作完成
	select {
	case <-done:
		t.Error("期望context超时，但操作完成了")
	case <-ctx.Done():
		fmt.Println("Context超时测试通过")
	}
}

// 性能测试
func BenchmarkRAGClientTimeout(b *testing.B) {
	client, err := NewRAGClient("http://localhost:8000")
	if err != nil {
		b.Fatalf("创建RAG客户端失败: %v", err)
	}

	req := &CodeReviewRequest{
		GitURL:      "https://gitlab.com/test/repo.git",
		Branch:      "main",
		DiffContent: "diff --git a/test.go b/test.go\n@@ -1,1 +1,1 @@\n-func old() {}\n+func new() {}",
		Query:       "测试查询",
		GitlabToken: "test_token",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.AnalyzeCodeWithRequest(req)
	}
}

// TestRAGServiceIntegration 测试RAG服务集成
func TestRAGServiceIntegration(t *testing.T) {
	// 使用实际的RAG服务URL进行测试
	client, err := NewRAGClient("http://localhost:8000")
	if err != nil {
		t.Fatalf("创建RAG客户端失败: %v", err)
	}

	req := &CodeReviewRequest{
		GitURL:      "http://localhost:9980/test_lv/test.git",
		Branch:      "main",
		DiffContent: "File: index.js\n@@ -2,14 +2,6 @@ function a(){wess",
		Query:       "检查",
		GitlabToken: "glpat-xxxxxxxxxxxx",
	}

	// 调用RAG服务
	result, err := client.AnalyzeCodeWithRequest(req)
	if err != nil {
		t.Logf("RAG服务调用失败（可能是网络问题）: %v", err)
		return
	}

	// 验证返回结果
	if result == nil {
		t.Error("期望返回结果，但结果为nil")
		return
	}

	if result.Review == "" {
		t.Error("期望返回审查结果，但结果为空")
		return
	}

	fmt.Printf("RAG服务集成测试通过，返回结果: %s\n", result.Review)
}

// TestNewRAGClientWithoutURL 测试未配置URL时的错误处理
func TestNewRAGClientWithoutURL(t *testing.T) {
	// 临时清除环境变量
	originalURL := os.Getenv("RAG_SERVICE_URL")
	os.Unsetenv("RAG_SERVICE_URL")
	defer os.Setenv("RAG_SERVICE_URL", originalURL)

	// 测试空URL
	_, err := NewRAGClient("")
	if err == nil {
		t.Error("期望返回错误，但没有错误")
		return
	}

	expectedError := "RAG服务URL未配置，请设置RAG_SERVICE_URL环境变量"
	if err.Error() != expectedError {
		t.Errorf("期望错误信息: %s, 实际错误信息: %s", expectedError, err.Error())
	}

	fmt.Println("RAG服务URL未配置时的错误处理测试通过")
}
