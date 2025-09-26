package providers

import (
	"bytes"
	"code-review-go/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UCAIProvider UCAI实现
type UCAIProvider struct {
	config *model.AIConfig
}

// NewUCAIProvider 创建UCAI提供者
func NewUCAIProvider(config *model.AIConfig) *UCAIProvider {
	return &UCAIProvider{
		config: config,
	}
}

// CallAI 调用UCAI接口
func (p *UCAIProvider) CallAI(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model": p.config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "你是一个专业的代码审核机器人，请仔细检查代码并提供详细的改进建议。"},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", p.config.APIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("UCAI API返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("UCAI API错误: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("UCAI API返回空结果")
	}

	return result.Choices[0].Message.Content, nil
}
