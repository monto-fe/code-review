package service

import (
	"code-review-go/internal/model"

	"gorm.io/gorm"
)

type AIManagerService struct {
	db *gorm.DB
}

func NewAIManagerService(db *gorm.DB) *AIManagerService {
	return &AIManagerService{db: db}
}

func (s *AIManagerService) GetManagerList() ([]model.AIManager, error) {
	var managers []model.AIManager
	if err := s.db.Find(&managers).Error; err != nil {
		return nil, err
	}
	return managers, nil
}
