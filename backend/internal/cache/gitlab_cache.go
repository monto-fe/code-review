package cache

import (
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"strings"
	"sync"
)

var (
	gitlabCache     map[uint]dto.GitlabCacheItem
	gitlabCacheLock sync.RWMutex
)

// InitGitlabCache 初始化缓存
func InitGitlabCache() error {
	var infos []model.GitlabInfo
	if err := database.DB.Where("status = ?", 1).Find(&infos).Error; err != nil {
		return err
	}

	cache := make(map[uint]dto.GitlabCacheItem)
	for _, info := range infos {
		var projectIDs []string
		if info.ProjectIds == "" {
			// 先写入空的，异步去加载
			go func(info model.GitlabInfo) {
				ids, err := utils.FetchProjectIDs(info.Token, info.API)
				if err != nil || len(ids) == 0 {
					// 获取不到就跳过
					return
				}
				projectIDsStr := strings.Join(ids, ",")
				// 更新数据库
				_ = database.DB.Model(&model.GitlabInfo{}).Where("id = ?", info.ID).
					Update("project_ids", projectIDsStr).Error
				// 更新缓存（加锁）
				gitlabCacheLock.Lock()
				item := gitlabCache[info.ID]
				item.ProjectIDs = ids
				item.Config.ProjectIds = projectIDsStr
				gitlabCache[info.ID] = item
				gitlabCacheLock.Unlock()
			}(info)
			projectIDs = []string{}
		} else {
			projectIDs = strings.Split(info.ProjectIds, ",")
		}

		cache[info.ID] = dto.GitlabCacheItem{
			Token:           info.Token,
			Config:          info,
			ProjectIDs:      projectIDs,
			Prompt:          info.Prompt,
			WebhookURL:      info.WebhookURL,
			WebhookStatus:   info.WebhookStatus,
			RuleCheckStatus: info.RuleCheckStatus,
		}
	}

	gitlabCacheLock.Lock()
	gitlabCache = cache
	gitlabCacheLock.Unlock()
	return nil
}

// GetGitlabCache 获取缓存
func GetGitlabCache() map[uint]dto.GitlabCacheItem {
	gitlabCacheLock.RLock()
	defer gitlabCacheLock.RUnlock()
	return gitlabCache
}

// FindTokenByProjectID 根据项目ID查找对应的Token和配置
func FindTokenByProjectID(projectID string, gitlabCache map[uint]dto.GitlabCacheItem) (string, dto.GitlabCacheItem, bool) {
	gitlabCacheLock.RLock()
	defer gitlabCacheLock.RUnlock()

	for _, item := range gitlabCache {
		for _, pid := range item.ProjectIDs {
			if pid == projectID {
				return item.Token, item, true
			}
		}
	}
	return "", dto.GitlabCacheItem{}, false
}
