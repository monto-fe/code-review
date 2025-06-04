package service

import (
	"time"

	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"

	"gorm.io/gorm"
)

// ResourceService 资源服务
type ResourceService struct {
	db *gorm.DB
}

// NewResourceService 创建资源服务实例
func NewResourceService(db *gorm.DB) *ResourceService {
	return &ResourceService{db: db}
}

// Create 创建资源
func (s *ResourceService) Create(req *model.ResourceReq) (*model.Resource, error) {
	resource := &model.Resource{
		Namespace:  req.Namespace,
		Resource:   req.Resource,
		Name:       req.Name,
		Category:   req.Category,
		Operator:   req.Operator,
		CreateTime: req.CreateTime,
		UpdateTime: req.UpdateTime,
	}

	if err := s.db.Create(resource).Error; err != nil {
		return nil, err
	}

	return resource, nil
}

// Update 更新资源
func (s *ResourceService) Update(resource *dto.UpdateResourceReq) error {
	resource.UpdateTime = time.Now().Unix()
	return s.db.Model(resource).Updates(resource).Error
}

// DeleteSelf 删除资源
func (s *ResourceService) DeleteSelf(id uint) error {
	return s.db.Delete(&model.Resource{}, id).Error
}

// FindWithAllChildren 获取项目组及其子项目组信息
func (s *ResourceService) FindWithAllChildren(query *model.ResourceQuery) ([]*model.Resource, int64, error) {
	var resources []*model.Resource
	var total int64

	db := s.db.Model(&model.Resource{}).Where("namespace = ?", query.Namespace)

	if query.Resource != "" {
		db = db.Where("resource LIKE ?", "%"+query.Resource+"%")
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if len(query.Category) > 0 {
		db = db.Where("category IN ?", query.Category)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := db.Offset(query.Offset).Limit(query.Limit).Order("id DESC").Find(&resources).Error; err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}

// FindAllById 通过ID查询资源
func (s *ResourceService) FindAllById(id uint) (*model.Resource, error) {
	var resource model.Resource
	if err := s.db.First(&resource, id).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

// FindAllByIds 通过ID列表查询资源
func (s *ResourceService) FindAllByIds(ids []uint) ([]*model.Resource, error) {
	var resources []*model.Resource
	if err := s.db.Where("id IN ?", ids).Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

// FindByNamespaceAndName 通过命名空间和名称查询资源
func (s *ResourceService) FindByNamespaceAndName(namespace, name string) (*model.Resource, error) {
	var resource model.Resource
	if err := s.db.Where("namespace = ? AND name = ?", namespace, name).First(&resource).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

// 获取资源列表
func (s *ResourceService) GetResourceList(query model.ResourceQuery) ([]model.Resource, int64, error) {
	db := s.db.Model(&model.Resource{})
	db = db.Where("namespace = ?", query.Namespace)
	if query.Resource != "" {
		db = db.Where("resource = ?", query.Resource)
	}
	if query.Name != "" {
		db = db.Where("name = ?", query.Name)
	}
	if len(query.Category) > 0 {
		db = db.Where("category IN ?", query.Category)
	}
	var total int64
	db.Count(&total)
	var resources []model.Resource
	if err := db.Limit(query.Limit).Offset(query.Offset).Find(&resources).Error; err != nil {
		return nil, 0, err
	}
	return resources, total, nil
}
