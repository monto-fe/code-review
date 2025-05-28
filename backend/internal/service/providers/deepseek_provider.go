package providers

import (
	"bytes"
	"code-review-go/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DeepSeekProvider DeepSeek实现
type DeepSeekProvider struct {
	config *model.AIConfig
}

// NewDeepSeekProvider 创建DeepSeek提供者
func NewDeepSeekProvider(config *model.AIConfig) *DeepSeekProvider {
	return &DeepSeekProvider{
		config: config,
	}
}

// CallAI 调用DeepSeek接口
func (p *DeepSeekProvider) CallAI(prompt string) (string, error) {
	fmt.Println("p.config.APIURL", p.config)
	requestBody := map[string]interface{}{
		"model": p.config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "你是一个专业的代码审核机器人，请仔细检查代码并提供详细的改进建议。"},
			{"role": "user", "content": prompt},
		},
		"stream": false,
	}

	fmt.Println("requestBody", requestBody)

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
		return "", fmt.Errorf("DeepSeek API返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
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
		return "", fmt.Errorf("DeepSeek API错误: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("DeepSeek API返回空结果")
	}

	return result.Choices[0].Message.Content, nil
}
