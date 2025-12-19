package providers

import (
	"bytes"
	"code-review-go/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// QwenProvider 通义千问实现
type QwenProvider struct {
	config *model.AIConfig
}

// NewQwenProvider 创建通义千问提供者
func NewQwenProvider(config *model.AIConfig) *QwenProvider {
	return &QwenProvider{
		config: config,
	}
}

// CallAI 调用通义千问接口
func (p *QwenProvider) CallAI(prompt []map[string]string) (string, error) {
	fmt.Printf("通义千问配置: URL=%s, Model=%s, Type=%s\n", p.config.APIURL, p.config.Model, p.config.Type)

	/**
	[]map[string]string{
			{"role": "system", "content": "你是一个专业的代码审核机器人，请仔细检查代码并提供详细的改进建议。"},
			{"role": "user", "content": prompt},
		}
	*/
	requestBody := map[string]interface{}{
		"model":       p.config.Model,
		"messages":    prompt,
		"stream":      false,
		"temperature": 0.1,
		"top_p":       0.8,
	}

	fmt.Printf("请求体: %+v\n", requestBody)

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

	// 读取响应体用于错误调试
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应体: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("通义千问 API返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
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

	// 使用已读取的body进行解析
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v, 响应体: %s", err, string(body))
	}

	if result.Error != nil {
		return "", fmt.Errorf("通义千问 API错误: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("通义千问 API返回空结果")
	}

	return result.Choices[0].Message.Content, nil
}
