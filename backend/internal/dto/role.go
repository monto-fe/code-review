package dto

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Desc        string `json:"desc"`
	Operator    string `json:"operator"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	ResourceIDs []uint `json:"resource_ids"`
}

type UpdateRoleReq struct {
	ID         uint   `json:"id" binding:"required"`
	Role       string `json:"role" binding:"required"`
	UpdateTime int64  `json:"update_time"`
}

type RoleQuery struct {
	Name   string `form:"name"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type RoleResource struct {
	RoleID     uint `json:"role_id" binding:"required"`
	ResourceID uint `json:"resource_id" binding:"required"`
}

type DeleteRoleRequest struct {
	ID uint `json:"id" binding:"required"`
}
