package service

import (
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"time"

	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: database.DB}
}

func (s *RoleService) GetRoleList(params dto.RoleQuery) ([]dto.Role, int64, error) {
	var roles []dto.Role
	var total int64
	if err := database.DB.Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (s *RoleService) CreateRole(role *dto.CreateRoleRequest) error {
	return s.db.Create(role).Error
}

func (s *RoleService) UpdateRole(role *dto.UpdateRoleReq) error {
	return s.db.Save(role).Error
}

func (s *RoleService) DeleteRole(req *dto.DeleteRoleRequest) error {
	return s.db.Delete(&model.Role{}, req.ID).Error
}

func (s *RoleService) CreateRoleWithResources(req *dto.CreateRoleRequest) error {
	// 1. 创建角色
	role := &model.Role{
		Name:       req.Name,
		Describe:   req.Desc,
		Operator:   req.Operator,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
		// 其他字段
	}
	if err := s.db.Create(role).Error; err != nil {
		return err
	}

	// 2. 创建角色-资源关联
	var roleResources []dto.RoleResource
	for _, rid := range req.ResourceIDs {
		roleResources = append(roleResources, dto.RoleResource{
			RoleID:     role.ID,
			ResourceID: rid,
		})
	}
	if len(roleResources) > 0 {
		if err := s.db.Create(&roleResources).Error; err != nil {
			return err
		}
	}
	return nil
}
