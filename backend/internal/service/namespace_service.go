package service

import (
	"time"

	"code-review-go/internal/model"

	"gorm.io/gorm"
)

// NamespaceService 命名空间服务
type NamespaceService struct {
	db *gorm.DB
}

// NewNamespaceService 创建命名空间服务实例
func NewNamespaceService(db *gorm.DB) *NamespaceService {
	return &NamespaceService{db: db}
}

// Create 创建命名空间
func (s *NamespaceService) Create(req *model.NamespaceReq) (*model.Namespace, error) {
	namespace := &model.Namespace{
		Namespace:  req.Namespace,
		Parent:     req.Parent,
		Name:       req.Name,
		Describe:   req.Describe,
		Operator:   req.Operator,
		CreateTime: req.CreateTime,
		UpdateTime: req.UpdateTime,
	}

	if err := s.db.Create(namespace).Error; err != nil {
		return nil, err
	}

	return namespace, nil
}

// Update 更新命名空间
func (s *NamespaceService) Update(namespace *model.Namespace) error {
	// 如果只有ID和Describe字段，则只更新Describe
	if namespace.Describe != "" {
		return s.db.Model(namespace).Update("describe", namespace.Describe).Error
	}

	// 更新所有字段
	namespace.UpdateTime = time.Now().Unix()
	return s.db.Model(namespace).Updates(namespace).Error
}

// DeleteSelf 删除命名空间
func (s *NamespaceService) DeleteSelf(id uint) error {
	return s.db.Delete(&model.Namespace{}, id).Error
}

// FindWithAllChildren 获取项目组及其子项目组信息
func (s *NamespaceService) FindWithAllChildren(query *model.NamespaceQuery) ([]*model.Namespace, int64, error) {
	var namespaces []*model.Namespace
	var total int64

	db := s.db.Model(&model.Namespace{})
	if query.Namespace != "" {
		db = db.Where("namespace = ?", query.Namespace)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := db.Offset(query.Offset).Limit(query.Limit).Order("id DESC").Find(&namespaces).Error; err != nil {
		return nil, 0, err
	}

	return namespaces, total, nil
}

// CheckNamespace 校验项目组是否存在
func (s *NamespaceService) CheckNamespace(check *model.NamespaceCheck) (bool, error) {
	var count int64
	if err := s.db.Model(&model.Namespace{}).Where("namespace = ? OR name = ?", check.Namespace, check.Name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
