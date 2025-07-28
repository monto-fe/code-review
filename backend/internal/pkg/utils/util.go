package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"code-review-go/config"
	"code-review-go/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// IsEmpty checks if a value is empty
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return v == 0
	case map[string]interface{}:
		return len(v) == 0
	case []interface{}:
		return len(v) == 0
	default:
		return false
	}
}

// PageSize default pagination size
var PageSize = struct {
	Offset int
	Limit  int
}{
	Offset: 0,
	Limit:  10,
}

// Token constants
const (
	TokenSecretKey      = "a8cL!9@#2$FzQwErTyUiOp1234567890zxcvbnmASDFGHJKLQWERTYUIOP!@#$%^&*()_+-="
	TokenExpired        = "8h"
	RefreshTokenExpired = "3d"
	RetCodeSuccess      = 0
	RetCodeNotLogIn     = 10000
)

// GenerateSignToken generates a JWT token
func GenerateSignToken(option interface{}, secretKey string, expiresIn string) (string, error) {
	if expiresIn == "" {
		expiresIn = TokenExpired
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": option,
		"exp":  time.Now().Add(parseDuration(expiresIn)).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

// CheckSignToken verifies a JWT token
func CheckSignToken(jwtToken string, secretKey string) (interface{}, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["data"], nil
	}

	return nil, fmt.Errorf("invalid token")
}

// IsTokenExpired checks if a token is expired
func IsTokenExpired(tokenExpiredAt int64) bool {
	now := time.Now().Unix()
	return tokenExpiredAt-now < 0
}

// SplitIntegerToObject splits an integer into an object with equal distribution
func SplitIntegerToObject(n int, j []string) map[string]int {
	result := make(map[string]int)
	length := len(j)

	if length == 0 {
		return result
	}

	quotient := n / length
	remainder := n % length

	for _, element := range j {
		result[element] = quotient
	}

	for i := 0; i < remainder; i++ {
		result[j[i]]++
	}

	return result
}

// Base64ToBlob converts a base64 string to a blob
// func Base64ToBlob(base64String string, mimeType string) ([]byte, error) {
// 	return base64.StdEncoding.DecodeString(base64String)
// }

// GenerateAudioFileKey generates a unique key for audio files
// func GenerateAudioFileKey(userID int, questionID int) string {
// 	now := time.Now()
// 	formattedDate := now.Format("2006_01_02")
// 	timestamp := now.Unix()
// 	return fmt.Sprintf("%d_%d_%s_%d.mp3", userID, timestamp, formattedDate, questionID)
// }

// PushWeChatInfo 生成企业微信群消息
func PushWeChatInfo(pathWithNamespace, mergeURL, result string, id int) string {
	return fmt.Sprintf(`🔍 您的「%s」合并请求「[%s](%s)」触发了AI检测，详情如下：
%s

📝 请您到 [系统](http://%s:%s/aicodecheck/commentList?id=%d) 中反馈Bug检测效果，帮助我们不断优化～`,
		pathWithNamespace, mergeURL, mergeURL, result,
		config.LoadConfig().IP, config.LoadConfig().Port, id)
}

// parseDuration parses a duration string into time.Duration
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		// Handle common cases
		switch s {
		case "1d":
			return 24 * time.Hour
		case "2d":
			return 48 * time.Hour
		case "3d":
			return 72 * time.Hour
		case "7d":
			return 168 * time.Hour
		case "30d":
			return 720 * time.Hour
		default:
			return 8 * time.Hour // Default to 8 hours
		}
	}
	return d
}

// MustMarshal 将对象转换为 JSON 字节
func MustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// CommonGetRequest 执行Get API请求
func CommonGetRequest(method, url, gitlabToken string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func MaskString(s string) string {
	if len(s) <= 8 {
		return s // 长度不足8直接返回原文
	}
	return s[:4] + "***" + s[len(s)-4:]
}

// GetQueryInt 获取 query 参数并转换为 int，带默认值
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	valStr := c.DefaultQuery(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

// FetchProjectIDs 获取项目ID列表
func FetchProjectIDs(token, gitlabAPI string) ([]string, error) {
	var projectIDs []string
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf("%s/v4/projects?per_page=%d&page=%d", gitlabAPI, perPage, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("PRIVATE-TOKEN", token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var projects []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
			return nil, err
		}

		if len(projects) == 0 {
			break
		}

		for _, proj := range projects {
			if id, ok := proj["id"].(float64); ok {
				projectIDs = append(projectIDs, fmt.Sprintf("%.0f", id))
			}
		}

		if len(projects) < perPage {
			break
		}
		page++
	}

	return projectIDs, nil
}

const CommentTypeCommon = 1 // 普通评论
const CommentTypeInline = 2 // 行级评论

const CodeReviewPrompt = `找出潜在问题：可能存在的bug、逻辑错误、安全隐患、资源泄漏等。`

// GeneratePrompt 生成常规AI提示词
func GeneratePrompt(rule string, mergeRequest *model.MergeRequestInfo, diff []model.Change) string {
	var diffContent string
	for _, change := range diff {
		diffContent += fmt.Sprintf("File: %s\n%s\n\n", change.NewPath, change.Diff)
	}
	return GenerateAICheckCommonPrompt(rule, mergeRequest.Title, mergeRequest.Description, diffContent)
}

// GenerateRAGEnhancedInlinePrompt RAG行级评论提示词
func GenerateRAGEnhancedInlinePrompt(ragResult, rule, title, description, diffContent string) string {
	return fmt.Sprintf(`
你是一位资深的代码审查专家。基于RAG服务的初步分析结果，请进行进一步的代码审查。

### RAG服务分析结果
%s

### 审查规则
%s

### 代码信息
**标题**: %s
**描述**: %s

### 代码差异
%s

请基于RAG分析结果和上述规则进行深入审查，找出潜在问题：可能存在的bug、逻辑错误、安全隐患、资源泄漏等，用中文输出。

**输出要求：**
1. 每个问题单独一行，格式如下（严格按照此格式输出，便于自动解析）：
   文件路径:行号:问题描述:修改建议
2. 示例：
   src/main.go:15:可能存在空指针:建议增加非空判断
   src/utils.go:42:未处理的错误:建议添加错误处理
3. 如果没有发现问题，请输出：✅ 代码审查通过 - 未发现明显问题，代码质量良好

**重要注意事项：**
- 行号必须是diff中实际存在的可评论行（即diff视图中有"+"或" "的行）
- 对于文件级问题（如"文件末尾缺少换行符"、"文件编码问题"等），输出普通评论
- 确保指定的行号在diff patch中确实存在，避免无效行号导致评论失败

请不要输出多余的解释或表格，只输出上述格式内容。
`, ragResult, rule, title, description, diffContent)
}

// GenerateRAGEnhancedCommonPrompt RAG普通评论提示词
func GenerateRAGEnhancedCommonPrompt(ragResult, rule, title, description, diffContent string) string {
	return fmt.Sprintf(`
你是一位资深的代码审查专家。基于RAG服务的初步分析结果，请进行进一步的代码审查。

### RAG服务分析结果
%s

### 审查规则
%s

### 代码信息
**标题**: %s
**描述**: %s

### 代码差异
%s

请基于RAG分析结果和上述规则进行深入审查，找出潜在问题：可能存在的bug、逻辑错误、安全隐患、资源泄漏等，用中文输出。
如果有问题用Markdown表格格式输出：
      | 代码文件路径(行号) | 疑似Bug | 修改建议 |
如果没有发现问题，请输出：✅ 代码审查通过 - 未发现明显问题，代码质量良好
`, ragResult, rule, title, description, diffContent)
}

// GenerateAICheckInlinePrompt 生成AI检查行级评论提示词
func GenerateAICheckInlinePrompt(rule, title, description, diffContent string) string {
	return fmt.Sprintf(`请检查以下代码差异（diff），确保其符合以下要求：
规则：%s

代码信息：
标题：%s
描述：%s

代码差异：
%s

请基于Diff内容和上述规则进行深入审查，找出潜在问题：可能存在的bug、逻辑错误、安全隐患、资源泄漏等，用中文输出。
如果有问题用Markdown表格格式输出：
      | 代码文件路径(行号) | 疑似Bug | 修改建议 |

`, rule, title, description, diffContent)
}

// GenerateAICheckCommonPrompt 生成AI检查普通评论提示词
func GenerateAICheckCommonPrompt(rule, title, description, diffContent string) string {
	return fmt.Sprintf(`请检查以下代码差异（diff），确保其符合以下要求：
规则：%s

代码信息：
标题：%s
描述：%s

代码差异：
%s

请基于Diff内容和上述规则进行深入审查，找出潜在问题：可能存在的bug、逻辑错误、安全隐患、资源泄漏等，用中文输出。
如果有问题用Markdown表格格式输出：
      | 代码文件路径(行号) | 疑似Bug | 修改建议 |
如果没有发现问题，请输出：✅ 代码审查通过 - 未发现明显问题，代码质量良好`,
		rule, title, description, diffContent)
}
