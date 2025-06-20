package service

import (
	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrParamsInvalid = errors.New("invalid parameters")
)

// ProjectIdsSynced 枚举
const (
	ProjectIdsSyncedFailed  int8 = 1 // 缓存失败
	ProjectIdsSyncedPending int8 = 2 // 缓存中
	ProjectIdsSyncedSuccess int8 = 3 // 缓存成功
)

// GitlabService Gitlab服务
type GitlabService struct {
	db *gorm.DB
}

// NewGitlabService 创建Gitlab服务实例
func NewGitlabService(db *gorm.DB) *GitlabService {
	return &GitlabService{db: db}
}

// GetGitlabInfo 获取Gitlab信息
func (s *GitlabService) GetGitlabInfo() ([]dto.GitlabInfoResponse, error) {
	var gitlabList []model.GitlabInfo
	if err := s.db.Find(&gitlabList).Error; err != nil {
		return nil, err
	}

	// 转换为响应模型，过滤敏感信息
	responseList := make([]dto.GitlabInfoResponse, len(gitlabList))
	for i, info := range gitlabList {
		responseList[i] = dto.GitlabInfoResponse{
			ID:               info.ID,
			Name:             info.Name,
			API:              utils.MaskString(info.API),
			WebhookURL:       utils.MaskString(info.WebhookURL),
			WebhookName:      info.WebhookName,
			Status:           info.Status,
			GitlabVersion:    info.GitlabVersion,
			Expired:          info.Expired,
			GitlabURL:        info.GitlabURL,
			SourceBranch:     info.SourceBranch,
			TargetBranch:     info.TargetBranch,
			Prompt:           info.Prompt,
			WebhookStatus:    info.WebhookStatus,
			ProjectIdsSynced: info.ProjectIdsSynced,
			RuleCheckStatus:  info.RuleCheckStatus,
			CreateTime:       info.CreateTime,
			UpdateTime:       info.UpdateTime,
		}
	}
	return responseList, nil
}

// CreateGitlabToken 创建Gitlab Token
func (s *GitlabService) CreateGitlabToken(data model.GitlabInfoCreate) (*model.GitlabInfo, error) {
	if data.API == "" || data.Token == "" {
		return nil, ErrParamsInvalid
	}

	now := time.Now().Unix()
	gitlabInfo := &model.GitlabInfo{
		API:              data.API,
		Name:             data.Name,
		Token:            data.Token,
		WebhookName:      data.WebhookName,
		WebhookURL:       data.WebhookURL,
		Status:           data.Status,
		GitlabVersion:    data.GitlabVersion,
		GitlabURL:        data.GitlabURL,
		SourceBranch:     data.SourceBranch,
		TargetBranch:     data.TargetBranch,
		Prompt:           data.Prompt,
		WebhookStatus:    data.WebhookStatus,
		RuleCheckStatus:  data.RuleCheckStatus,
		ProjectIdsSynced: ProjectIdsSyncedPending,
		CreateTime:       now,
		UpdateTime:       now,
	}

	if data.Expired != 0 {
		gitlabInfo.Expired = data.Expired
	}

	if err := s.db.Create(gitlabInfo).Error; err != nil {
		return nil, err
	}

	// 异步获取项目列表
	go func() {
		var projectIDsStr string
		var updateStatus int8

		if data.Token == "" {
			updateStatus = ProjectIdsSyncedFailed
		} else {
			projectIDs, err := utils.FetchProjectIDs(data.Token, data.API)
			if err != nil || len(projectIDs) == 0 {
				updateStatus = ProjectIdsSyncedFailed
			} else {
				projectIDsStr = strings.Join(projectIDs, ",")
				updateStatus = ProjectIdsSyncedSuccess
			}
		}

		updateMap := map[string]interface{}{
			"project_ids":        projectIDsStr,
			"project_ids_synced": updateStatus,
		}
		if dbErr := s.db.Model(&model.GitlabInfo{}).Where("id = ?", gitlabInfo.ID).Updates(updateMap).Error; dbErr != nil {
			s.db.Model(&model.GitlabInfo{}).Where("id = ?", gitlabInfo.ID).Update("project_ids_synced", ProjectIdsSyncedFailed)
		}
	}()

	return gitlabInfo, nil
}

// UpdateGitlabInfo 更新Gitlab信息
func (s *GitlabService) UpdateGitlabInfo(data model.GitlabInfoUpdate) (*model.GitlabInfo, error) {
	if data.ID == 0 {
		return nil, ErrParamsInvalid
	}

	// 检查记录是否存在
	var existing model.GitlabInfo
	if err := s.db.First(&existing, data.ID).Error; err != nil {
		return nil, err
	}

	// 更新字段
	updates := map[string]interface{}{
		"update_time": time.Now().Unix(),
	}

	if data.API != "" {
		updates["api"] = data.API
	}
	if data.Name != "" {
		updates["name"] = data.Name
	}
	if data.Token != "" {
		updates["token"] = data.Token
		updates["project_ids_synced"] = ProjectIdsSyncedPending // 如果更新了 token，重置同步状态
	}
	if data.WebhookName != "" {
		updates["webhook_name"] = data.WebhookName
	}
	if data.WebhookURL != "" {
		updates["webhook_url"] = data.WebhookURL
	}
	if data.Status != 0 {
		updates["status"] = data.Status
	}
	if data.GitlabVersion != "" {
		updates["gitlab_version"] = data.GitlabVersion
	}
	if data.GitlabURL != "" {
		updates["gitlab_url"] = data.GitlabURL
	}
	if data.Expired != 0 {
		updates["expired"] = data.Expired
	}
	if data.SourceBranch != "" {
		updates["source_branch"] = data.SourceBranch
	}
	if data.TargetBranch != "" {
		updates["target_branch"] = data.TargetBranch
	}
	if data.Prompt != "" {
		updates["prompt"] = data.Prompt
	}
	if data.WebhookStatus != 0 {
		updates["webhook_status"] = data.WebhookStatus
	}
	if data.RuleCheckStatus != 0 {
		updates["rule_check_status"] = data.RuleCheckStatus
	}

	if err := s.db.Model(&existing).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 如果更新了 token，异步获取新的项目列表
	if data.Token != "" {
		go func() {
			projectIDs, err := utils.FetchProjectIDs(data.Token, data.API)
			projectIDsStr := ""
			if err == nil {
				projectIDsStr = strings.Join(projectIDs, ",")
			}
			updateMap := map[string]interface{}{
				"project_ids":        projectIDsStr,
				"project_ids_synced": ProjectIdsSyncedSuccess,
			}
			if err != nil || projectIDsStr == "" {
				updateMap["project_ids_synced"] = ProjectIdsSyncedFailed
			}
			if dbErr := s.db.Model(&model.GitlabInfo{}).Where("id = ?", data.ID).Updates(updateMap).Error; dbErr != nil {
				s.db.Model(&model.GitlabInfo{}).Where("id = ?", data.ID).Update("project_ids_synced", ProjectIdsSyncedFailed)
			}
		}()
	} else {
		s.db.Model(&model.GitlabInfo{}).Where("id = ?", data.ID).Update("project_ids_synced", ProjectIdsSyncedFailed)
	}

	// 获取更新后的记录
	var updated model.GitlabInfo
	if err := s.db.First(&updated, data.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

// DeleteGitlabToken 删除Gitlab Token
func (s *GitlabService) DeleteGitlabToken(id uint) error {
	if id == 0 {
		return ErrParamsInvalid
	}

	// 检查记录是否存在
	var existing model.GitlabInfo
	if err := s.db.First(&existing, id).Error; err != nil {
		return err
	}

	// 删除记录
	return s.db.Delete(&existing).Error
}

// GetProjectLanguages 获取项目语言
func GetProjectLanguages(gitlabAPI, projectID, gitlabToken string) (map[string]float64, error) {
	url := fmt.Sprintf("%s/v4/projects/%s/languages", gitlabAPI, projectID)

	body, err := utils.CommonGetRequest("GET", url, gitlabToken, nil)
	if err != nil {
		return nil, fmt.Errorf("获取项目语言失败: %v", err)
	}

	var languages map[string]float64
	if err := json.Unmarshal(body, &languages); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return languages, nil
}

// GetDominantLanguage 获取主要语言
func GetDominantLanguage(gitlabAPI, projectID, gitlabToken string) (string, error) {
	languages, err := GetProjectLanguages(gitlabAPI, projectID, gitlabToken)
	if err != nil {
		return "", err
	}

	// 计算总行数
	var totalLines float64
	for _, lines := range languages {
		totalLines += lines
	}

	// 如果没有语言统计数据，返回空字符串
	if totalLines == 0 {
		return "", nil
	}

	// 计算每种语言的占比，并找到占比最大的语言
	var dominantLanguage string
	var maxPercentage float64

	for language, lines := range languages {
		percentage := (lines / totalLines) * 100
		if percentage > maxPercentage {
			maxPercentage = percentage
			dominantLanguage = language
		}
	}

	return dominantLanguage, nil
}

// GetMergeRequestInfo 获取合并请求信息
func GetMergeRequestInfo(gitlabAPI, projectID, gitlabToken string) (*model.MergeRequestInfo, error) {
	url := fmt.Sprintf("%s/v4/projects/%s/merge_requests", gitlabAPI, projectID)

	body, err := utils.CommonGetRequest("GET", url, gitlabToken, nil)
	if err != nil {
		return nil, fmt.Errorf("获取合并请求失败: %v", err)
	}

	var mergeRequests []model.MergeRequestInfo
	if err := json.Unmarshal(body, &mergeRequests); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(mergeRequests) == 0 {
		return nil, fmt.Errorf("no merge requests found")
	}

	return &mergeRequests[0], nil
}

// GetMergeRequestDiff 获取合并请求差异
func GetMergeRequestDiff(gitlabAPI, projectID, mergeRequestID, gitlabToken string) ([]model.Change, error) {
	url := fmt.Sprintf("%s/v4/projects/%s/merge_requests/%s/changes", gitlabAPI, projectID, mergeRequestID)

	body, err := utils.CommonGetRequest("GET", url, gitlabToken, nil)
	if err != nil {
		return nil, fmt.Errorf("获取合并请求差异失败: %v", err)
	}

	var diff struct {
		Changes []model.Change `json:"changes"`
	}
	if err := json.Unmarshal(body, &diff); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return diff.Changes, nil
}

// CommentResponse GitLab评论响应
type CommentResponse struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	Author struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PostCommentToGitLab 向GitLab发送评论
func PostCommentToGitLab(gitlabAPI string, projectID, mergeRequestID int, gitlabToken, comment string) (*CommentResponse, error) {
	url := fmt.Sprintf("%s/v4/projects/%d/merge_requests/%d/notes", gitlabAPI, projectID, mergeRequestID)

	requestBody := map[string]string{
		"body": comment,
	}

	body, err := utils.CommonGetRequest("POST", url, gitlabToken, requestBody)
	if err != nil {
		return nil, fmt.Errorf("发送评论失败: %v", err)
	}

	var result CommentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// 通过id获取Gitlab Token 详情，不返回token
func (s *GitlabService) GetGitlabTokenDetail(id uint) (*model.GitlabInfo, error) {
	var gitlabInfo model.GitlabInfo
	if err := s.db.First(&gitlabInfo, id).Error; err != nil {
		return nil, err
	}
	gitlabInfo.Token = ""
	return &gitlabInfo, nil
}
