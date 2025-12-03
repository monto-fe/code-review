package model

// ReviewWhitelist 代码审查白名单模型
type ReviewWhitelist struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	ProjectName string `gorm:"column:project_name;type:varchar(200);not null" json:"project_name"`
	ProjectID   uint   `gorm:"column:project_id;not null" json:"project_id"`
	Remark      string `gorm:"type:varchar(500)" json:"remark"`
	Operator    string `gorm:"type:varchar(50);not null" json:"operator"`
	CreateTime  int64  `gorm:"not null" json:"create_time"`
	UpdateTime  int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (ReviewWhitelist) TableName() string {
	return TableReviewWhitelist
}

// ReviewWhitelistCreate 创建白名单请求
type ReviewWhitelistCreate struct {
	ProjectName string `json:"project_name" binding:"required"`
	ProjectID   uint   `json:"project_id" binding:"required"`
	Remark      string `json:"remark"`
	Operator    string `json:"operator" binding:"required"`
}

// ReviewWhitelistUpdate 更新白名单请求
type ReviewWhitelistUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	ProjectName string `json:"project_name"`
	ProjectID   uint   `json:"project_id"`
	Remark      string `json:"remark"`
	Operator    string `json:"operator"`
}

// ReviewWhitelistDelete 删除白名单请求
type ReviewWhitelistDelete struct {
	ID uint `json:"id" binding:"required"`
}

// ReviewWhitelistQuery 查询白名单请求
type ReviewWhitelistQuery struct {
	Current     int    `json:"current" form:"current"`
	PageSize    int    `json:"pageSize" form:"pageSize"`
	ProjectName string `json:"project_name" form:"project_name"`
	ProjectID   uint   `json:"project_id" form:"project_id"`
}

// ReviewWhitelistResponse 白名单响应
type ReviewWhitelistResponse struct {
	ID          uint   `json:"id"`
	ProjectName string `json:"project_name"`
	ProjectID   uint   `json:"project_id"`
	Remark      string `json:"remark"`
	Operator    string `json:"operator"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
}

// ToResponse 转换为响应结构体
func (r *ReviewWhitelist) ToResponse() ReviewWhitelistResponse {
	return ReviewWhitelistResponse{
		ID:          r.ID,
		ProjectName: r.ProjectName,
		ProjectID:   r.ProjectID,
		Remark:      r.Remark,
		Operator:    r.Operator,
		CreateTime:  r.CreateTime,
		UpdateTime:  r.UpdateTime,
	}
}
