package business

import (
	"code-review-go/internal/cache"
	"code-review-go/internal/model"
	"time"

	"gorm.io/gorm"
)

// ReviewWhitelistService 代码审查白名单服务
type ReviewWhitelistService struct {
	db *gorm.DB
}

// NewReviewWhitelistService 创建白名单服务实例
func NewReviewWhitelistService(db *gorm.DB) *ReviewWhitelistService {
	return &ReviewWhitelistService{db: db}
}

// Create 创建白名单记录
func (s *ReviewWhitelistService) Create(req *model.ReviewWhitelistCreate) (*model.ReviewWhitelist, error) {
	now := time.Now().Unix()
	whitelist := &model.ReviewWhitelist{
		ProjectName: req.ProjectName,
		ProjectID:   req.ProjectID,
		Remark:      req.Remark,
		Operator:    req.Operator,
		CreateTime:  now,
		UpdateTime:  now,
	}

	if err := s.db.Create(whitelist).Error; err != nil {
		return nil, err
	}

	// 刷新缓存
	cache.RefreshReviewWhitelistCache()

	return whitelist, nil
}

// Update 更新白名单记录
func (s *ReviewWhitelistService) Update(req *model.ReviewWhitelistUpdate) error {
	now := time.Now().Unix()

	// 查找当前记录
	var current model.ReviewWhitelist
	if err := s.db.First(&current, req.ID).Error; err != nil {
		return err
	}

	// 构造更新字段
	updateMap := map[string]interface{}{
		"update_time": now,
	}

	if req.ProjectName != "" {
		updateMap["project_name"] = req.ProjectName
	}
	if req.ProjectID > 0 {
		updateMap["project_id"] = req.ProjectID
	}
	if req.Remark != "" {
		updateMap["remark"] = req.Remark
	}
	if req.Operator != "" {
		updateMap["operator"] = req.Operator
	}

	if err := s.db.Model(&model.ReviewWhitelist{}).
		Where("id = ?", req.ID).
		Updates(updateMap).Error; err != nil {
		return err
	}

	// 刷新缓存
	cache.RefreshReviewWhitelistCache()

	return nil
}

// Delete 删除白名单记录
func (s *ReviewWhitelistService) Delete(id uint) error {
	if err := s.db.Delete(&model.ReviewWhitelist{}, id).Error; err != nil {
		return err
	}

	// 刷新缓存
	cache.RefreshReviewWhitelistCache()

	return nil
}

// GetList 分页查询白名单列表
func (s *ReviewWhitelistService) GetList(query *model.ReviewWhitelistQuery) ([]*model.ReviewWhitelist, int64, error) {
	var whitelists []*model.ReviewWhitelist
	var total int64

	db := s.db.Model(&model.ReviewWhitelist{})

	// 条件查询
	if query.ProjectName != "" {
		db = db.Where("project_name LIKE ?", "%"+query.ProjectName+"%")
	}
	if query.ProjectID > 0 {
		db = db.Where("project_id = ?", query.ProjectID)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页参数
	page := query.Current
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 分页查询
	if err := db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&whitelists).Error; err != nil {
		return nil, 0, err
	}

	return whitelists, total, nil
}

// GetByProjectID 根据项目ID查询
func (s *ReviewWhitelistService) GetByProjectID(projectID uint) (*model.ReviewWhitelist, error) {
	var whitelist model.ReviewWhitelist
	if err := s.db.Where("project_id = ?", projectID).First(&whitelist).Error; err != nil {
		return nil, err
	}
	return &whitelist, nil
}

// GetAllProjectIDs 获取所有白名单项目ID列表（用于缓存）
func (s *ReviewWhitelistService) GetAllProjectIDs() ([]uint, error) {
	var projectIDs []uint
	if err := s.db.Model(&model.ReviewWhitelist{}).
		Pluck("project_id", &projectIDs).Error; err != nil {
		return nil, err
	}
	return projectIDs, nil
}
