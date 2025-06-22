package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// RAGClient RAG服务客户端
type RAGClient struct {
	baseURL string
	client  *http.Client
}

// NewRAGClient 创建RAG服务客户端
func NewRAGClient(baseURL string) (*RAGClient, error) {
	if baseURL == "" {
		baseURL = os.Getenv("RAG_SERVICE_URL")
		if baseURL == "" {
			return nil, fmt.Errorf("RAG服务URL未配置，请设置RAG_SERVICE_URL环境变量")
		}
	}
	return &RAGClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 180 * time.Second, // 设置3分钟超时
		},
	}, nil
}

// CodeReviewRequest 代码审查请求
type CodeReviewRequest struct {
	GitURL      string `json:"git_url"`
	Branch      string `json:"branch"`
	DiffContent string `json:"diff_content"`
	Query       string `json:"query,omitempty"`
	GitlabToken string `json:"gitlab_token,omitempty"`
}

// CodeAnalysisResponse 代码分析响应
type CodeAnalysisResponse struct {
	Review string `json:"review"`
}

// CodeReviewPrompt 代码审查提示词模板
const CodeReviewPrompt = `你是一个专业的代码审查助手。请基于以下信息提供全面的代码审查建议：

1. 代码变更：
%s

2. 相关上下文：
%s

请从以下几个方面提供详细的审查建议：
1. 代码质量：代码的可读性、可维护性、性能等方面
2. 潜在问题：可能存在的bug、安全隐患、资源泄漏等
3. 改进建议：具体的优化方案和最佳实践
4. 示例代码：如果可能，提供改进后的代码示例

请以结构化的方式输出审查结果。`

// AnalyzeCode 实现RAGService接口的方法
func (c *RAGClient) AnalyzeCode(gitURL string, branch string, diffFiles []string, codeChanges string) (string, error) {
	req := &CodeReviewRequest{
		GitURL:      gitURL,
		Branch:      branch,
		DiffContent: codeChanges,
		Query:       "代码审查和优化建议",
	}

	analysis, err := c.AnalyzeCodeWithRequest(req)
	if err != nil {
		return "", err
	}

	return analysis.Review, nil
}

// AnalyzeCodeWithRequest 使用请求对象进行分析
func (c *RAGClient) AnalyzeCodeWithRequest(req *CodeReviewRequest) (*CodeAnalysisResponse, error) {
	// 创建3分钟超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建带超时的请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		c.baseURL,
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("RAG服务请求超时(3分钟): %v", err)
		}
		return nil, fmt.Errorf("请求RAG服务失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RAG服务返回错误状态码: %d", resp.StatusCode)
	}

	var result CodeAnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析RAG服务响应失败: %v", err)
	}

	return &result, nil
}

// GenerateReview 根据分析结果生成代码审查
func (c *RAGClient) GenerateReview(req *CodeReviewRequest) (string, error) {
	analysis, err := c.AnalyzeCodeWithRequest(req)
	if err != nil {
		return "", err
	}

	// 直接返回RAG服务的审查结果
	return analysis.Review, nil
}

// LegacyAnalyzeCode 实现RAGService接口的方法
func (c *RAGClient) LegacyAnalyzeCode(gitURL string, branch string, diffFiles []string, codeChanges string) (string, error) {
	req := &CodeReviewRequest{
		GitURL:      gitURL,
		Branch:      branch,
		DiffContent: codeChanges,
		Query:       "代码审查和优化建议",
	}

	analysis, err := c.AnalyzeCodeWithRequest(req)
	if err != nil {
		return "", err
	}

	return analysis.Review, nil
}
