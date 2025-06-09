package model

// AIConfig AI配置模型
type AIConfig struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	Name       string `gorm:"type:varchar(100);not null" json:"name"`
	APIURL     string `gorm:"column:api_url;type:varchar(200);not null" json:"api_url"`
	APIKey     string `gorm:"column:api_key;type:varchar(200);not null" json:"api_key"`
	Model      string `gorm:"type:varchar(50);not null" json:"model"`
	Type       string `gorm:"type:varchar(50);not null" json:"type"`
	IsActive   int8   `gorm:"column:is_active;not null;default:2" json:"is_active"` // 1: 启用, 2: 禁用
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (AIConfig) TableName() string {
	return TableAIConfig
}

// AIConfigCreate 创建AI配置请求
type AIConfigCreate struct {
	Name     string `json:"name"` // 预留字段，暂时不使用
	APIURL   string `json:"api_url" binding:"required"`
	APIKey   string `json:"api_key" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Model    string `json:"model" binding:"required"`
	IsActive int8   `json:"is_active"` // 1: 启用, 2: 禁用
}

// AIConfigUpdate 更新AI配置请求
type AIConfigUpdate struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name"`
	APIURL   string `json:"api_url"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	Type     string `json:"type"`      // UCloud, DeepSeek, OpenAI
	IsActive int8   `json:"is_active"` // 1: 启用, 2: 禁用
}

// AIConfigInfo AI配置信息
type AIConfigInfo struct {
	Name   string `json:"name"`
	APIKey string `json:"api_key"`
	APIURL string `json:"api_url"`
	Type   string `json:"type"`
	Model  string `json:"model"`
}

// AIConfigDelete 删除AI配置请求
type AIConfigDelete struct {
	ID uint `json:"id" binding:"required"`
}

// AIConfigResponse AI配置响应
type AIConfigResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	APIURL     string `json:"api_url"`
	Model      string `json:"model"`
	Type       string `json:"type"`
	IsActive   int8   `json:"is_active"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

// ToResponse 转换为响应结构体
func (c *AIConfig) ToResponse() AIConfigResponse {
	return AIConfigResponse{
		ID:         c.ID,
		Name:       c.Name,
		APIURL:     c.APIURL,
		Model:      c.Model,
		Type:       c.Type,
		IsActive:   c.IsActive,
		CreateTime: c.CreateTime,
		UpdateTime: c.UpdateTime,
	}
}
