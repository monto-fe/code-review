package gitlab_service

import (
	"code-review-go/internal/cache"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bytes"

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
	if err := s.db.Order("create_time DESC").Find(&gitlabList).Error; err != nil {
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
			CommentType:      info.CommentType,
			CreateTime:       info.CreateTime,
			UpdateTime:       info.UpdateTime,
		}
	}
	return responseList, nil
}

// CreateGitlabToken 创建Gitlab Token
func (s *GitlabService) CreateGitlabToken(data model.GitlabInfoCreate) (*dto.GitlabCreateUpdateResponse, error) {
	if data.API == "" || data.Token == "" {
		return nil, ErrParamsInvalid
	}

	now := time.Now().Unix()

	// 设置 CommentType 默认值
	commentType := data.CommentType
	if commentType == 0 {
		commentType = 1 // 默认使用普通评论
	}

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
		CommentType:      commentType,
		CreateTime:       now,
		UpdateTime:       now,
	}

	if data.Expired != 0 {
		gitlabInfo.Expired = data.Expired
	}

	if err := s.db.Create(gitlabInfo).Error; err != nil {
		return nil, err
	}

	// 创建新记录后，直接添加到缓存中，避免全量刷新
	cache.UpdateGitlabCacheItem(gitlabInfo.ID, *gitlabInfo)

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

		// 项目列表更新后再次刷新缓存
		if err := cache.RefreshGitlabCache(); err != nil {
			fmt.Printf("异步更新项目列表后刷新缓存失败: %v\n", err)
		}
	}()

	// 返回不包含敏感信息的响应
	return &dto.GitlabCreateUpdateResponse{
		ID:               gitlabInfo.ID,
		Name:             gitlabInfo.Name,
		API:              utils.MaskString(gitlabInfo.API),
		WebhookURL:       utils.MaskString(gitlabInfo.WebhookURL),
		WebhookName:      gitlabInfo.WebhookName,
		Status:           gitlabInfo.Status,
		GitlabVersion:    gitlabInfo.GitlabVersion,
		Expired:          gitlabInfo.Expired,
		GitlabURL:        gitlabInfo.GitlabURL,
		SourceBranch:     gitlabInfo.SourceBranch,
		TargetBranch:     gitlabInfo.TargetBranch,
		Prompt:           gitlabInfo.Prompt,
		WebhookStatus:    gitlabInfo.WebhookStatus,
		ProjectIds:       gitlabInfo.ProjectIds,
		ProjectIdsSynced: gitlabInfo.ProjectIdsSynced,
		RuleCheckStatus:  gitlabInfo.RuleCheckStatus,
		CommentType:      gitlabInfo.CommentType,
		CreateTime:       gitlabInfo.CreateTime,
		UpdateTime:       gitlabInfo.UpdateTime,
	}, nil
}

// UpdateGitlabInfo 更新Gitlab信息
func (s *GitlabService) UpdateGitlabInfo(data model.GitlabInfoUpdate) (*dto.GitlabCreateUpdateResponse, error) {
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
	if data.CommentType != 0 {
		updates["comment_type"] = data.CommentType
	}

	if err := s.db.Model(&existing).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 获取更新后的记录用于缓存更新
	var updated model.GitlabInfo
	if err := s.db.First(&updated, data.ID).Error; err != nil {
		return nil, err
	}

	// 智能缓存刷新：根据更新的字段决定刷新策略
	needFullRefresh := false
	if data.Token != "" || data.API != "" || data.Status != 0 {
		// 关键字段更新，需要全量刷新缓存
		needFullRefresh = true
	} else {
		// 其他字段更新，只更新单个缓存项
		cache.UpdateGitlabCacheItem(data.ID, updated)
	}

	// 如果需要全量刷新，则刷新整个缓存
	if needFullRefresh {
		if err := cache.RefreshGitlabCache(); err != nil {
			// 记录错误但不影响主流程
			fmt.Printf("更新后刷新缓存失败: %v\n", err)
		}
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

			// 项目列表更新后再次刷新缓存
			if err := cache.RefreshGitlabCache(); err != nil {
				fmt.Printf("异步更新项目列表后刷新缓存失败: %v\n", err)
			}
		}()
	}

	// 返回不包含敏感信息的响应
	return &dto.GitlabCreateUpdateResponse{
		ID:               updated.ID,
		Name:             updated.Name,
		API:              utils.MaskString(updated.API),
		WebhookURL:       utils.MaskString(updated.WebhookURL),
		WebhookName:      updated.WebhookName,
		Status:           updated.Status,
		GitlabVersion:    updated.GitlabVersion,
		Expired:          updated.Expired,
		GitlabURL:        updated.GitlabURL,
		SourceBranch:     updated.SourceBranch,
		TargetBranch:     updated.TargetBranch,
		Prompt:           updated.Prompt,
		WebhookStatus:    updated.WebhookStatus,
		ProjectIds:       updated.ProjectIds,
		ProjectIdsSynced: updated.ProjectIdsSynced,
		RuleCheckStatus:  updated.RuleCheckStatus,
		CommentType:      updated.CommentType,
		CreateTime:       updated.CreateTime,
		UpdateTime:       updated.UpdateTime,
	}, nil
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
	if err := s.db.Delete(&existing).Error; err != nil {
		return err
	}

	// 删除记录后，直接从缓存中移除，避免全量刷新
	cache.DeleteGitlabCacheItem(id)

	return nil
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

// 获取diff refs
func GetMergeRequestInfoRefs(gitlabAPI, projectID, mergeRequestIID, gitlabToken string) (*model.MergeRequestInfo, error) {
	url := fmt.Sprintf("%s/v4/projects/%s/merge_requests/%s", gitlabAPI, projectID, mergeRequestIID)

	body, err := utils.CommonGetRequest("GET", url, gitlabToken, nil)
	if err != nil {
		return nil, fmt.Errorf("获取合并请求失败: %v", err)
	}

	var mergeRequest model.MergeRequestInfo
	if err := json.Unmarshal(body, &mergeRequest); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &mergeRequest, nil
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

	var changes []model.Change
	for _, item := range diff.Changes {
		//过滤一下
		if strings.HasPrefix(item.Diff, "Binary files") {
			fmt.Printf("过滤掉二进制文件:%v\n", item.NewPath)
			continue
		}
		if item.Diff == "" {
			fmt.Printf("过滤掉空文件变更:%v\n", item.NewPath)
			continue
		}

		if !isCodeFile(item.NewPath) {
			fmt.Printf("过滤掉非代码文件或忽略文件:%v\n", item.NewPath)
			continue
		}
		changes = append(changes, item)
	}
	return changes, nil
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

	fmt.Printf("GitLab评论请求结果成功: %+v\n", result)

	return &result, nil
}

// CreateDiscussionForPush 为 Push 事件创建 Discussion
func CreateDiscussionForPush(gitlabAPI string, projectID int, gitlabToken, title, body string) (*CommentResponse, error) {
	url := fmt.Sprintf("%s/v4/projects/%d/discussions", gitlabAPI, projectID)

	requestBody := map[string]string{
		"body": body,
	}

	responseBody, err := utils.CommonGetRequest("POST", url, gitlabToken, requestBody)
	if err != nil {
		return nil, fmt.Errorf("创建Discussion失败: %v", err)
	}

	var result CommentResponse
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// LineCommentRequest 行级评论请求
type LineCommentRequest struct {
	Body     string `json:"body"`
	Position struct {
		BaseSHA      string `json:"base_sha"`
		StartSHA     string `json:"start_sha"`
		HeadSHA      string `json:"head_sha"`
		NewPath      string `json:"new_path"`
		NewLine      int    `json:"new_line"`
		PositionType string `json:"position_type"`
	} `json:"position"`
}

// PostLineCommentToGitLab 向GitLab发送行级评论
func PostLineCommentToGitLab(gitlabAPI string, projectID, mergeRequestID int, gitlabToken, comment string, position LineCommentRequest) (*CommentResponse, error) {
	url := fmt.Sprintf("%s/v4/projects/%d/merge_requests/%d/discussions", gitlabAPI, projectID, mergeRequestID)
	fmt.Printf("GitLab行级评论请求：\nURL: %s\nToken: %s...%s\nComment: %s\nPosition: %+v\n",
		url,
		gitlabToken[:4], gitlabToken[len(gitlabToken)-4:],
		comment,
		position,
	)
	requestBody := map[string]interface{}{
		"body": comment,
		"position": map[string]interface{}{
			"base_sha":      position.Position.BaseSHA,
			"start_sha":     position.Position.StartSHA,
			"head_sha":      position.Position.HeadSHA,
			"new_path":      position.Position.NewPath,
			"new_line":      position.Position.NewLine,
			"position_type": position.Position.PositionType,
		},
	}
	fmt.Printf("请求体: %+v\n", requestBody)

	// 直接使用http包发送请求，以便获取详细错误信息
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("PRIVATE-TOKEN", gitlabToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("GitLab API错误响应，状态码: %d，响应内容: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("GitLab API返回错误状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	var result CommentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("解析响应失败，响应内容: %s\n", string(body))
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	fmt.Printf("行级评论发送成功，响应: %+v\n", result)
	return &result, nil
}

// CommentInfo 评论信息
type CommentInfo struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // error, warning, info
	Type     string `json:"type"`     // "line" 或 "general"
}

// 返回新文件的最大行号
func getNewFileLineCount(diffText string) int {
	lines := strings.Split(diffText, "\n")
	maxLine := 0
	curLine := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			// 解析 @@ -a,b +c,d @@
			// 例：@@ -1,10 +1,9 @@
			parts := strings.Split(line, " ")
			for _, part := range parts {
				if strings.HasPrefix(part, "+") {
					plus := strings.TrimPrefix(part, "+")
					plus = strings.Split(plus, ",")[0]
					if n, err := strconv.Atoi(plus); err == nil {
						curLine = n - 1 // diff 行号起始
					}
				}
			}
		} else if len(line) > 0 && (line[0] == ' ' || line[0] == '+') {
			curLine++
			if curLine > maxLine {
				maxLine = curLine
			}
		}
	}
	return maxLine
}

// 获取所有可评论的 new_line 行号
func getCommentableLines(diffText string) map[int]bool {
	lines := strings.Split(diffText, "\n")
	commentable := make(map[int]bool)
	curLine := 0

	fmt.Printf("解析diff文本:\n%s\n", diffText)

	for i, line := range lines {
		fmt.Printf("第%d行: %s\n", i, line)

		if strings.HasPrefix(line, "@@") {
			// 解析 @@ -a,b +c,d @@
			parts := strings.Split(line, " ")
			for _, part := range parts {
				if strings.HasPrefix(part, "+") {
					plus := strings.TrimPrefix(part, "+")
					plus = strings.Split(plus, ",")[0]
					if n, err := strconv.Atoi(plus); err == nil {
						curLine = n - 1
						fmt.Printf("找到新文件起始行号: %d\n", curLine)
					}
				}
			}
		} else if len(line) > 0 && (line[0] == ' ' || line[0] == '+') {
			curLine++
			commentable[curLine] = true
			fmt.Printf("标记可评论行号: %d (内容: %s)\n", curLine, line)
		}
	}

	fmt.Printf("最终可评论行号: %+v\n", commentable)
	return commentable
}

// ParseCommentsForLineComments 解析评论内容，返回所有评论（包括行级和普通评论）
func ParseCommentsForLineComments(comments string, diff []model.Change) []CommentInfo {
	var allComments []CommentInfo

	// 统计每个文件的可评论行号
	fileCommentableLines := make(map[string]map[int]bool)
	for _, change := range diff {
		fileCommentableLines[change.NewPath] = getCommentableLines(change.Diff)
	}

	lines := strings.Split(comments, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		commentInfo := parseLineComment(line, fileCommentableLines)
		if commentInfo != nil {
			allComments = append(allComments, *commentInfo)
		}
	}

	return allComments
}

// 修改 parseLineComment，返回所有解析的评论
func parseLineComment(line string, fileCommentableLines map[string]map[int]bool) *CommentInfo {
	// 格式: "src/main.go:15: 问题"
	if strings.Contains(line, ":") {
		parts := strings.SplitN(line, ":", 3)
		if len(parts) >= 3 {
			filePath := parts[0]
			if lineNum, err := strconv.Atoi(parts[1]); err == nil {
				message := parts[2]
				commentInfo := &CommentInfo{
					File:     filePath,
					Line:     lineNum,
					Message:  strings.TrimSpace(message),
					Severity: determineSeverity(message),
				}

				// 检查是否可以进行行级评论
				commentable := fileCommentableLines[filePath]
				if commentable != nil && commentable[lineNum] {
					commentInfo.Type = "line"
				} else {
					commentInfo.Type = "general"
				}

				return commentInfo
			}
		}
	}
	return nil
}

// ClassifyComments 将评论分类为行级评论和普通评论
func ClassifyComments(comments []CommentInfo) ([]CommentInfo, []CommentInfo) {
	var lineComments []CommentInfo
	var generalComments []CommentInfo

	for _, comment := range comments {
		if comment.Type == "line" {
			lineComments = append(lineComments, comment)
		} else {
			generalComments = append(generalComments, comment)
		}
	}

	return lineComments, generalComments
}

// determineSeverity 确定评论的严重程度
func determineSeverity(message string) string {
	messageLower := strings.ToLower(message)

	// 错误级别关键词
	errorKeywords := []string{"error", "错误", "严重", "critical", "致命", "exception", "异常", "crash", "崩溃"}
	for _, keyword := range errorKeywords {
		if strings.Contains(messageLower, keyword) {
			return "error"
		}
	}

	// 警告级别关键词
	warningKeywords := []string{"warning", "警告", "注意", "attention", "潜在", "potential", "可能", "maybe"}
	for _, keyword := range warningKeywords {
		if strings.Contains(messageLower, keyword) {
			return "warning"
		}
	}

	// 默认返回信息级别
	return "info"
}

// PostLineComments 发送行级评论，返回失败的行级评论列表
func PostLineComments(gitlabAPI string, projectID, mergeRequestID int, gitlabToken string, comments []dto.LineComment, diff []model.Change) ([]string, error) {
	// 获取合并请求的SHA信息
	mergeRequest, err := GetMergeRequestInfoRefs(gitlabAPI, strconv.Itoa(projectID), strconv.Itoa(mergeRequestID), gitlabToken)
	if err != nil {
		return nil, fmt.Errorf("获取合并请求信息失败: %v", err)
	}
	baseSHA := mergeRequest.DiffRefs.BaseSHA
	startSHA := mergeRequest.DiffRefs.StartSHA
	headSHA := mergeRequest.DiffRefs.HeadSHA

	failedComments := make([]string, 0)

	for _, comment := range comments {
		// 查找对应的文件变更
		var targetChange *model.Change
		for _, change := range diff {
			if change.NewPath == comment.File || change.OldPath == comment.File {
				targetChange = &change
				break
			}
		}
		if targetChange == nil {
			// 文件变更都找不到，API一定会失败
			failedComments = append(failedComments, fmt.Sprintf("%s:%d: %s", comment.File, comment.Line, comment.Message))
			continue
		}

		// 构建位置信息
		position := LineCommentRequest{}
		// 为行级评论添加质量反馈复选框
		position.Body = fmt.Sprintf("[%s] %s\n\n---\n**请评价此评论的质量：**\n- [ ] 精准定位并提供建议\n- [ ] 部分有效但需要改进\n- [ ] 完全误导性建议\n- [ ] 建议不相关\n\n**评价说明：**\n- 精准定位：评论准确指出了代码问题，并提供了具体可行的改进建议\n- 部分有效：评论有一定道理，但建议不够具体或存在小问题\n- 完全误导：评论对代码理解有误，建议会引入新的问题\n- 不相关：评论与当前代码变更无关\n\n**您的反馈将帮助我们改进AI代码审查质量！**",
			strings.ToUpper(comment.Severity), comment.Message)
		position.Position.BaseSHA = baseSHA
		position.Position.StartSHA = startSHA
		position.Position.HeadSHA = headSHA
		position.Position.NewPath = comment.File
		position.Position.NewLine = comment.Line
		position.Position.PositionType = "text"

		// 只以API返回为准
		_, err := PostLineCommentToGitLab(gitlabAPI, projectID, mergeRequestID, gitlabToken, position.Body, position)
		if err != nil {
			failedComments = append(failedComments, fmt.Sprintf("%s:%d: %s", comment.File, comment.Line, comment.Message))
		}
	}

	return failedComments, nil
}

// GetGitlabTokenDetail 通过id获取Gitlab Token 详情，不返回token
func (s *GitlabService) GetGitlabTokenDetail(id uint) (*model.GitlabInfo, error) {
	var gitlabInfo model.GitlabInfo
	if err := s.db.First(&gitlabInfo, id).Error; err != nil {
		return nil, err
	}
	gitlabInfo.Token = ""
	return &gitlabInfo, nil
}

// GetGitlabTokenProjects 获取Gitlab Token项目信息列表
// 只返回ID、Name和ProjectIds字段，全量拉取不分页
func (s *GitlabService) GetGitlabTokenProjects() ([]dto.GitlabTokenProjectResponse, error) {
	var gitlabList []model.GitlabInfo

	// 只查询需要的字段，提高性能
	if err := s.db.Select("id, name, project_ids").Find(&gitlabList).Error; err != nil {
		return nil, err
	}

	// 转换为响应模型
	responseList := make([]dto.GitlabTokenProjectResponse, len(gitlabList))
	for i, info := range gitlabList {
		responseList[i] = dto.GitlabTokenProjectResponse{
			ID:         info.ID,
			Name:       info.Name,
			ProjectIds: info.ProjectIds,
		}
	}

	return responseList, nil
}

// GetCommitCodeContent 获取指定commit的代码内容
func GetCommitCodeContent(gitlabAPI string, projectID int, commitSHA, gitlabToken string) ([]model.Change, error) {
	// 构建GitLab API URL来获取commit的文件树
	// 注意：GitLab 12.9.2 可能对某些参数有限制
	url := fmt.Sprintf("%s/v4/projects/%d/repository/tree?ref=%s&recursive=true&per_page=100", gitlabAPI, projectID, commitSHA)

	fmt.Printf("🔍 获取commit文件树 (GitLab 12.9.2兼容): URL=%s, ProjectID=%d, CommitSHA=%s\n", url, projectID, commitSHA)

	body, err := utils.CommonGetRequest("GET", url, gitlabToken, nil)
	if err != nil {
		fmt.Printf("❌ 获取commit文件树失败: %v\n", err)
		return nil, fmt.Errorf("获取commit文件树失败: %v", err)
	}

	var treeResponse []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Path string `json:"path"`
		Mode string `json:"mode"`
	}

	if err := json.Unmarshal(body, &treeResponse); err != nil {
		return nil, fmt.Errorf("解析文件树响应失败: %v", err)
	}

	// 过滤出代码文件（排除目录和二进制文件）
	var codeFiles []string
	for _, item := range treeResponse {
		if item.Type == "blob" && isCodeFile(item.Name) {
			if strings.HasPrefix(item.Name, ".") || strings.HasPrefix(item.Path, ".") {
				fmt.Printf("❌ 文件/目录以.开头忽略 %s: \n", item.Path)
				continue
			}
			codeFiles = append(codeFiles, item.Path)
		}
	}

	// 获取每个代码文件的内容
	var changes []model.Change
	for _, filePath := range codeFiles {
		content, err := getFileContent(gitlabAPI, projectID, commitSHA, filePath, gitlabToken)
		if err != nil {
			fmt.Printf("❌ 获取文件内容失败 %s: %v\n", filePath, err)
			continue
		}

		change := model.Change{
			OldPath: filePath,
			NewPath: filePath,
			Diff:    content, // 直接使用文件内容，不添加 +++ 前缀
		}
		changes = append(changes, change)
		fmt.Printf("✅ 成功获取文件内容: %s (大小: %d 字节)\n", filePath, len(content))
	}

	fmt.Printf("获取commit代码内容成功: %s, 代码文件数: %d\n", commitSHA, len(changes))
	return changes, nil
}

// GetBranchCodeContent 获取指定分支的代码内容（GitLab 12.9.2 兼容性方法）
func GetBranchCodeContent(gitlabAPI string, projectID int, branchName, gitlabToken string) ([]model.Change, error) {
	// 构建GitLab API URL来获取分支的文件树
	url := fmt.Sprintf("%s/v4/projects/%d/repository/tree?ref=%s&recursive=true&per_page=100", gitlabAPI, projectID, branchName)

	fmt.Printf("🔍 获取分支文件树 (GitLab 12.9.2兼容): URL=%s, ProjectID=%d, Branch=%s\n", url, projectID, branchName)

	body, err := utils.CommonGetRequest("GET", url, gitlabToken, nil)
	if err != nil {
		fmt.Printf("❌ 获取分支文件树失败: %v\n", err)
		return nil, fmt.Errorf("获取分支文件树失败: %v", err)
	}

	var treeResponse []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Path string `json:"path"`
		Mode string `json:"mode"`
	}

	if err := json.Unmarshal(body, &treeResponse); err != nil {
		fmt.Printf("❌ 解析分支文件树响应失败: %v, 响应内容: %s\n", err, string(body))
		return nil, fmt.Errorf("解析分支文件树响应失败: %v", err)
	}

	fmt.Printf("📋 分支文件树项目总数: %d\n", len(treeResponse))

	// 过滤出代码文件（排除目录和二进制文件）
	var codeFiles []string
	for _, item := range treeResponse {
		if strings.HasPrefix(item.Name, ".") || strings.HasPrefix(item.Path, ".") {
			fmt.Printf("❌ 文件/目录以.开头忽略 %s: \n", item.Path)
			continue
		}
		if item.Type == "blob" && isCodeFile(item.Name) {
			codeFiles = append(codeFiles, item.Path)
			fmt.Printf("✅ 识别为代码文件: %s\n", item.Path)
		} else {
			fmt.Printf("❌ 跳过文件: %s (Type=%s, isCodeFile=%v)\n", item.Path, item.Type, isCodeFile(item.Name))
		}
	}

	fmt.Printf("📝 过滤后的代码文件数: %d\n", len(codeFiles))

	// 获取每个代码文件的内容
	var changes []model.Change
	for _, filePath := range codeFiles {
		content, err := getFileContent(gitlabAPI, projectID, branchName, filePath, gitlabToken)
		if err != nil {
			fmt.Printf("❌ 获取文件内容失败 %s: %v\n", filePath, err)
			continue
		}

		change := model.Change{
			OldPath: filePath,
			NewPath: filePath,
			Diff:    content,
		}
		changes = append(changes, change)
		fmt.Printf("✅ 成功获取文件内容: %s (大小: %d 字节)\n", filePath, len(content))
	}

	fmt.Printf("🎯 获取分支代码内容成功: %s, 代码文件数: %d\n", branchName, len(changes))
	return changes, nil
}

// getFileContent 获取单个文件的内容
func getFileContent(gitlabAPI string, projectID int, commitSHA, filePath, gitlabToken string) (string, error) {
	// 构建GitLab API URL来获取文件内容
	// 注意：GitLab 12.9.2 可能需要特殊的文件路径编码
	encodedPath := strings.ReplaceAll(filePath, "/", "%2F")
	url := fmt.Sprintf("%s/v4/projects/%d/repository/files/%s/raw?ref=%s",
		gitlabAPI, projectID, encodedPath, commitSHA)

	fmt.Printf("📄 获取文件内容 (GitLab 12.9.2兼容): URL=%s, FilePath=%s\n", url, filePath)

	body, err := utils.CommonGetRequest("GET", url, gitlabToken, nil)
	if err != nil {
		fmt.Printf("❌ 获取文件内容失败: %v\n", err)
		return "", fmt.Errorf("获取文件内容失败: %v", err)
	}

	content := string(body)
	fmt.Printf("✅ 文件内容获取成功: %s (大小: %d 字节)\n", filePath, len(content))
	return content, nil
}

// isCodeFile 判断是否为代码文件
func isCodeFile(filename string) bool {
	// 定义代码文件扩展名

	codeExtensions := map[string]bool{
		".go":         true,
		".js":         true,
		".ts":         true,
		".tsx":        true,
		".jsx":        true,
		".py":         true,
		".java":       true,
		".cpp":        true,
		".c":          true,
		".h":          true,
		".cs":         true,
		".php":        true,
		".rb":         true,
		".swift":      true,
		".kt":         true,
		".scala":      true,
		".rs":         true,
		".vue":        true,
		".html":       true,
		".css":        true,
		".scss":       true,
		".less":       true,
		".sql":        true,
		".sh":         true,
		".yaml":       true,
		".yml":        true,
		".json":       true,
		".xml":        true,
		".md":         true,
		".dockerfile": true,
		".makefile":   true,
	}

	value := os.Getenv("IGNORE_EXT")
	if value != "" {
		extArray := strings.Split(value, ",")
		for _, ext := range extArray {
			if codeExtensions[ext] == true {
				codeExtensions[ext] = false
			}
		}
	}

	// 检查文件扩展名
	for ext := range codeExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return codeExtensions[ext]
		}
	}

	// 检查特殊文件名
	specialFiles := map[string]bool{
		"dockerfile": true,
		"makefile":   true,
		"rakefile":   true,
		"gemfile":    true,
		"podfile":    true,
	}

	return specialFiles[strings.ToLower(filename)]
}

// CreateCommitComment 为 commit 创建评论
func CreateCommitComment(gitlabAPI string, projectID int, commitSHA, gitlabToken, title, comment string) (*CommentResponse, error) {
	url := fmt.Sprintf("%s/v4/projects/%d/repository/commits/%s/comments", gitlabAPI, projectID, commitSHA)

	// 构建评论内容
	commentBody := fmt.Sprintf("## %s\n\n%s", title, comment)

	requestBody := map[string]string{
		"note": commentBody,
	}

	body, err := utils.CommonGetRequest("POST", url, gitlabToken, requestBody)
	if err != nil {
		return nil, fmt.Errorf("创建Commit评论失败: %v", err)
	}

	var result CommentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}
