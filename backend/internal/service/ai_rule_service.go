package service

import (
	"code-review-go/internal/model"

	"gorm.io/gorm"
)

// AIRuleService AI规则服务
type AIRuleService struct {
	db *gorm.DB
}

// NewAIRuleService 创建AI规则服务实例
func NewAIRuleService(db *gorm.DB) *AIRuleService {
	return &AIRuleService{db: db}
}

// GetCustomRuleByProjectID 获取项目的自定义规则
func (s *AIRuleService) GetCustomRuleByProjectID(projectID uint) (*model.CustomRule, error) {
	var rule model.CustomRule
	err := s.db.Where("project_id = ?", projectID).First(&rule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

// GetCommonRule 获取通用规则
func (s *AIRuleService) GetCommonRule(language string) ([]model.CommonRule, error) {
	var rules []model.CommonRule
	err := s.db.Where("language = ?", language).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}
