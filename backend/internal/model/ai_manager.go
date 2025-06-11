package model

type AIManager struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Type       string `json:"type"`
	Model      string `json:"model"`
	APIURL     string `json:"api_url"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

// TableName 指定表名
func (AIManager) TableName() string {
	return TableAIManager
}
