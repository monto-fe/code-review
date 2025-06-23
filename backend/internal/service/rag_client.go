package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"code-review-go/internal/pkg/utils"
)

// RAGClient RAG服务客户端
type RAGClient struct {
	baseURL string
	client  *http.Client
}

// NewRAGClient 创建RAG服务客户端
func NewRAGClient(baseURL string) *RAGClient {
	return &RAGClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 180 * time.Second, // 设置3分钟超时
		},
	}
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

// generateEnhancedPrompt 生成增强的提示词
func (m *RAGServiceManager) generateEnhancedPrompt(ragResult string, data *AnalysisData) string {
	return utils.GenerateRAGEnhancedPrompt(
		ragResult,
		data.FinalRule,
		data.MergeRequest.Title,
		data.MergeRequest.Description,
		data.DiffStr,
	)
}
