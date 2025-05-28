package service

import (
	"code-review-go/internal/model"

	"gorm.io/gorm"
)

// PermissionService 权限服务
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// GetRolePermissions 获取角色权限
func (s *PermissionService) GetRolePermissions(query *model.PermissionQuery) ([]*model.RolePermission, int64, error) {
	var permissions []*model.RolePermission
	var total int64

	db := s.db.Model(&model.RolePermission{}).Where("namespace = ? AND role_id = ?", query.Namespace, query.RoleID)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := db.Offset(query.Offset).Limit(query.Limit).Order("id DESC").Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// GetRoleResourceIds 获取角色资源ID列表
func (s *PermissionService) GetRoleResourceIds(namespace string, roleID uint) ([]*model.RolePermission, error) {
	var permissions []*model.RolePermission
	if err := s.db.Where("namespace = ? AND role_id = ?", namespace, roleID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// AssetRolePermission 分配单个角色权限
func (s *PermissionService) AssetRolePermission(req *model.RoleSinglePermissionReq) (*model.RolePermission, error) {
	permission := &model.RolePermission{
		Namespace:  req.Namespace,
		RoleID:     req.RoleID,
		ResourceID: req.ResourceID,
		Operator:   req.Operator,
		Describe:   req.Describe,
	}

	if err := s.db.Create(permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

// AssertBulkRolePermission 批量分配角色权限
func (s *PermissionService) AssertBulkRolePermission(namespace string, roleID uint, resourceIDs []uint, describe, operator string, createTime, updateTime int64) error {
	permissions := make([]*model.RolePermission, len(resourceIDs))
	for i, resourceID := range resourceIDs {
		permissions[i] = &model.RolePermission{
			Namespace:  namespace,
			RoleID:     roleID,
			ResourceID: resourceID,
			Describe:   describe,
			Operator:   operator,
			CreateTime: createTime,
			UpdateTime: updateTime,
		}
	}

	return s.db.Create(&permissions).Error
}

// AssertUserRole 批量分配用户角色
func (s *PermissionService) AssertUserRole(req *model.UserRoleReq) error {
	userRoles := make([]*model.UserRole, len(req.RoleIDs))
	for i, roleID := range req.RoleIDs {
		userRoles[i] = &model.UserRole{
			Namespace:  req.Namespace,
			User:       req.User,
			RoleID:     roleID,
			Status:     1,
			Operator:   req.Operator,
			CreateTime: req.CreateTime,
			UpdateTime: req.UpdateTime,
		}
	}

	return s.db.Create(&userRoles).Error
}

// CancelRolePermission 取消角色权限
func (s *PermissionService) CancelRolePermission(namespace string, roleID, resourceID uint) error {
	return s.db.Where("namespace = ? AND role_id = ? AND resource_id = ?", namespace, roleID, resourceID).Delete(&model.RolePermission{}).Error
}

// CancelUserRole 取消用户角色
func (s *PermissionService) CancelUserRole(namespace, user string, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return nil
	}
	return s.db.Where("namespace = ? AND user = ? AND role_id IN ?", namespace, user, roleIDs).Delete(&model.UserRole{}).Error
}

// DeleteSelf 删除自身权限
func (s *PermissionService) DeleteSelf(userID, jobID uint, rest map[string]interface{}) ([]*model.RolePermission, error) {
	var permissions []*model.RolePermission
	query := s.db.Model(&model.RolePermission{}).Where("user_id = ? AND job_id = ?", userID, jobID)

	for key, value := range rest {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetSelfPermissions 获取自身权限
func (s *PermissionService) GetSelfPermissions(query *model.SelfPermissionQuery) ([]*model.Resource, error) {
	// 1. 获取用户角色ID列表
	var userRoles []*model.UserRole
	if err := s.db.Where("user = ? AND namespace = ?", query.User, query.Namespace).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	roleIDs := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	// 2. 获取角色权限ID列表
	var rolePermissions []*model.RolePermission
	if err := s.db.Where("role_id IN ? AND namespace = ?", roleIDs, query.Namespace).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	resourceIDs := make([]uint, len(rolePermissions))
	for i, rp := range rolePermissions {
		resourceIDs[i] = rp.ResourceID
	}

	// 3. 获取资源列表
	var resources []*model.Resource
	db := s.db.Where("id IN ? AND namespace = ?", resourceIDs, query.Namespace)
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if err := db.Find(&resources).Error; err != nil {
		return nil, err
	}

	return resources, nil
}
