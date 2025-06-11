package dto

type ResourceReq struct {
	Resource string `json:"resource" binding:"required"`
	Category string `json:"category" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type UpdateResourceReq struct {
	ID         uint   `json:"id" binding:"required"`
	Resource   string `json:"resource" binding:"required"`
	Category   string `json:"category" binding:"required"`
	Name       string `json:"name" binding:"required"`
	UpdateTime int64  `json:"update_time"`
}

type Resource struct {
	// 根据你的 Resource 字段定义
	ID         uint   `json:"id"`
	Namespace  string `json:"namespace"`
	Resource   string `json:"resource"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Operator   string `json:"operator"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type ResourceListResponse struct {
	Data  []Resource `json:"data"`
	Total int64      `json:"total"`
}
