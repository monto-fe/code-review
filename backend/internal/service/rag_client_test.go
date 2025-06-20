package service

import (
	"code-review-go/internal/model"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRAGClientTimeout(t *testing.T) {
	// 测试超时配置
	client := NewRAGClient("http://localhost:8000")

	// 验证HTTP客户端超时设置
	if client.client.Timeout != 60*time.Second {
		t.Errorf("期望超时时间为60秒，实际为: %v", client.client.Timeout)
	}

	fmt.Println("HTTP客户端超时设置正确: 60秒")
}

func TestAnalyzeCodeWithRequestTimeout(t *testing.T) {
	// 创建RAG客户端
	client := NewRAGClient("http://invalid-url:9999") // 使用无效URL来测试超时

	// 准备测试请求
	req := &CodeReviewRequest{
		GitURL:      "https://gitlab.com/test/repo.git",
		Branch:      "main",
		DiffContent: "diff --git a/test.go b/test.go\n@@ -1,1 +1,1 @@\n-func old() {}\n+func new() {}",
		Query:       "测试查询",
		GitlabToken: "test_token",
	}

	// 记录开始时间
	start := time.Now()

	// 调用方法，应该超时
	_, err := client.AnalyzeCodeWithRequest(req)

	// 记录结束时间
	duration := time.Since(start)

	// 验证超时行为
	if err == nil {
		t.Error("期望出现超时错误，但没有错误")
		return
	}

	// 检查是否在合理的时间内超时（应该小于70秒，大于50秒）
	if duration > 70*time.Second {
		t.Errorf("超时时间过长: %v", duration)
	}

	if duration < 50*time.Second {
		t.Errorf("超时时间过短: %v", duration)
	}

	fmt.Printf("超时测试通过，耗时: %v，错误: %v\n", duration, err)
}

func TestCheckMergeRequestWithRAGTimeout(t *testing.T) {
	// 创建测试数据
	mergeRequest := &model.MergeRequestInfo{
		ProjectID:    1,
		IID:          1,
		WebURL:       "https://gitlab.com/test/repo/-/merge_requests/1",
		SourceBranch: "feature/test",
		Title:        "测试MR",
		Description:  "测试描述",
	}

	diff := []model.Change{
		{
			OldPath: "test.go",
			NewPath: "test.go",
			Diff:    "@@ -1,1 +1,1 @@\n-func old() {}\n+func new() {}",
		},
	}

	// 记录开始时间
	start := time.Now()

	// 调用方法，应该超时
	_, err := CheckMergeRequestWithRAG(mergeRequest, diff, "http://invalid-api", "invalid-token")

	// 记录结束时间
	duration := time.Since(start)

	// 验证超时行为
	if err == nil {
		t.Error("期望出现超时错误，但没有错误")
		return
	}

	// 检查是否在合理的时间内超时（应该小于70秒，大于50秒）
	if duration > 70*time.Second {
		t.Errorf("超时时间过长: %v", duration)
	}

	if duration < 50*time.Second {
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
	client := NewRAGClient("http://localhost:8000")

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
