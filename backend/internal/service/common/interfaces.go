package common

import (
	"code-review-go/internal/dto"
	"code-review-go/internal/model"

	"gorm.io/gorm"
)

// DatabaseService 数据库服务接口
type DatabaseService interface {
	// 初始化数据库服务
	Initialize() error
	// 获取数据库连接
	GetDB() *gorm.DB
	// 获取AI规则服务
	GetAIRuleService() interface{}
	// 检查是否已初始化
	IsInitialized() bool
}

// RAGServiceInterface RAG服务接口
type RAGServiceInterface interface {
	// 初始化RAG服务
	Initialize() error
	// 分析代码
	AnalyzeCodeWithRequest(req *dto.CodeReviewRequest) (*dto.CodeReviewResponse, error)
	// 检查是否已初始化
	IsInitialized() bool
}

// AIService AI服务接口
type AIService interface {
	// 初始化AI服务
	Initialize() error
	// 调用AI服务
	CallAI(prompt string) (string, error)
	// 检查是否已初始化
	IsInitialized() bool
}

// NotificationService 通知服务接口
type NotificationService interface {
	// 初始化通知服务
	Initialize() error
	// 发送GitLab评论
	SendGitLabComment(api, token string, projectID, mergeRequestID int, comment string) error
	// 发送行级评论
	SendLineComments(api, token string, projectID, mergeRequestID int, comments []dto.LineComment, diff []model.Change) ([]string, error)
	// 发送增强通知
	SendEnhancedNotification(api, token string, projectID, mergeRequestID int, comments string, commentType int8, diff []model.Change) error
	// 发送Webhook通知
	SendWebhookNotification(webhookURL string, status int, projectNamespace, mergeURL, comments string, aiMessageID uint, mergeRequest *model.MergeRequestInfo) error
	// 检查是否已初始化
	IsInitialized() bool
}

// ConfigService 配置服务接口
type ConfigService interface {
	// 初始化配置服务
	Initialize() error
	// 获取配置
	GetConfig() *Config
	// 检查是否已初始化
	IsInitialized() bool
}

// 通用配置结构
type Config struct {
	RAGServiceURL string
	AIServiceURL  string
	DatabaseURL   string
	// 其他配置项...
}

// RAGClientPool RAG客户端连接池
type RAGClientPool struct {
	// 连接池实现
}

// NewRAGClientPool 创建RAG客户端连接池
func NewRAGClientPool(baseURL string, poolSize int) (*RAGClientPool, error) {
	return &RAGClientPool{}, nil
}

// GetClient 获取客户端
func (p *RAGClientPool) GetClient() interface{} {
	return nil
}

// ReturnClient 归还客户端
func (p *RAGClientPool) ReturnClient(client interface{}) {
	// 实现归还逻辑
}
