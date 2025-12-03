package cache

import (
	"code-review-go/internal/database"
	"code-review-go/internal/model"
	"sync"
)

var (
	reviewWhitelistCache     map[uint]bool
	reviewWhitelistCacheLock sync.RWMutex
)

// InitReviewWhitelistCache 从数据库加载所有白名单项目ID到缓存
func InitReviewWhitelistCache() error {
	var whitelists []model.ReviewWhitelist
	if err := database.DB.Find(&whitelists).Error; err != nil {
		reviewWhitelistCacheLock.Lock()
		reviewWhitelistCache = make(map[uint]bool)
		reviewWhitelistCacheLock.Unlock()
		return err
	}

	reviewWhitelistCacheLock.Lock()
	reviewWhitelistCache = make(map[uint]bool)
	for _, wl := range whitelists {
		reviewWhitelistCache[wl.ProjectID] = true
	}
	reviewWhitelistCacheLock.Unlock()
	return nil
}

// GetReviewWhitelistCache 获取白名单缓存
func GetReviewWhitelistCache() map[uint]bool {
	reviewWhitelistCacheLock.RLock()
	defer reviewWhitelistCacheLock.RUnlock()
	return reviewWhitelistCache
}

// IsProjectWhitelisted 检查项目是否在白名单中
func IsProjectWhitelisted(projectID uint) bool {
	reviewWhitelistCacheLock.RLock()
	defer reviewWhitelistCacheLock.RUnlock()
	return reviewWhitelistCache[projectID]
}

// RefreshReviewWhitelistCache 刷新白名单缓存
func RefreshReviewWhitelistCache() error {
	return InitReviewWhitelistCache()
}
