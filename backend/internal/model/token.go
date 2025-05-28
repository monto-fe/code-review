package model

// Token 登录令牌模型
type Token struct {
	ID         uint   `gorm:"primarykey;autoIncrement" json:"id"`
	Token      string `gorm:"type:varchar(255);not null" json:"token"`
	User       string `gorm:"type:varchar(50);not null" json:"user"`
	ExpiredAt  int64  `gorm:"not null" json:"expired_at"`
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (Token) TableName() string {
	return TableToken
}
