package business

import (
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
	"code-review-go/internal/service/ai"
	"code-review-go/internal/service/gitlab_service"
	"code-review-go/internal/service/webhook/processors"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ManualCheckService struct {
	db            *gorm.DB
	gitlabService *gitlab_service.GitlabService
	aiRuleService *ai.AIRuleService
}

func NewManualCheckService(db *gorm.DB, gitlabService *gitlab_service.GitlabService, aiRuleService *ai.AIRuleService) *ManualCheckService {
	return &ManualCheckService{
		db:            db,
		gitlabService: gitlabService,
		aiRuleService: aiRuleService,
	}
}

// CreateManualCheckTask 创建手动审核任务
func (s *ManualCheckService) CreateManualCheckTask(req dto.ManualCheckRequest, userID uint) (*model.ManualCheckTask, error) {
	// 解析合并请求URL，提取项目ID、合并请求ID和项目名称
	projectID, mergeID, projectName, err := s.parseMergeURL(req.MergeURL)
	if err != nil {
		return nil, fmt.Errorf("解析合并请求URL失败: %w", err)
	}

	// 序列化机器人配置列表
	// 确保BotConfigs不为nil，避免序列化为null
	botConfigs := req.BotConfigs
	if botConfigs == nil {
		botConfigs = []dto.BotConfig{}
	}

	botConfigsJSON, err := json.Marshal(botConfigs)
	if err != nil {
		return nil, fmt.Errorf("序列化机器人配置失败: %w", err)
	}

	// 通过req.AIModelID获取缓存中的model名称
	aiModel, err := s.getAIModelByID(req.AIModelID)
	if err != nil {
		return nil, fmt.Errorf("获取AI模型失败: %w", err)
	}

	fmt.Println("botConfigsJSON", string(botConfigsJSON))
	fmt.Println("req.BotConfigs", projectID, mergeID, projectName)

	// 创建任务记录
	now := time.Now().Unix()
	task := &model.ManualCheckTask{
		UserID:      userID,
		ProjectID:   projectID,
		ProjectName: projectName, // 使用解析出的项目名称
		MergeID:     mergeID,
		MergeURL:    req.MergeURL,
		Status:      1, // 进行中
		AIModel:     aiModel.Model,
		BotConfigs:  string(botConfigsJSON),
		CreateTime:  now,
		UpdateTime:  now,
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	// 异步获取并更新项目信息和AI模型信息
	go func() {
		s.updateTaskWithCachedInfo(task.ID, projectID, projectName, req.AIModelID)
	}()

	return task, nil
}

// ExecuteManualCheckTask 执行手动审核任务
func (s *ManualCheckService) ExecuteManualCheckTask(taskID uint) error {
	// 获取任务信息
	var task model.ManualCheckTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("任务不存在: %w", err)
	}

	// 更新任务状态为进行中
	s.db.Model(&task).Update("status", 1)

	// 构建webhook body
	webhookBody := dto.WebhookBody{
		ObjectKind: "merge_request",
		Project: dto.ProjectInfo{
			ID:                task.ProjectID,
			Name:              task.ProjectName,
			PathWithNamespace: task.ProjectName,
		},
		ObjectAttributes: dto.ObjectAttributes{
			IID:   task.MergeID,
			Title: task.MergeTitle,
			State: "opened",
		},
	}

	// 执行AI审核
	result, err := processors.CheckMergeRequestWithAI(webhookBody)

	// 更新任务状态和结果
	updateData := map[string]interface{}{
		"update_time": time.Now().Unix(),
	}

	if err != nil {
		updateData["status"] = 3 // 失败
		updateData["error_message"] = err.Error()
	} else {
		updateData["status"] = 2 // 完成
		updateData["result"] = result
	}

	if err := s.db.Model(&task).Updates(updateData).Error; err != nil {
		return fmt.Errorf("更新任务状态失败: %w", err)
	}

	return err
}

// GetManualCheckHistory 获取手动审核历史
func (s *ManualCheckService) GetManualCheckHistory(req dto.ManualCheckHistoryRequest, userID uint) (*dto.ManualCheckHistoryResponse, error) {
	query := s.db.Model(&model.ManualCheckTask{}).Where("user_id = ?", userID)

	// 状态筛选
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 时间范围筛选
	if req.StartDate != "" {
		startTime, _ := strconv.ParseInt(req.StartDate, 10, 64)
		query = query.Where("create_time >= ?", startTime)
	}
	if req.EndDate != "" {
		endTime, _ := strconv.ParseInt(req.EndDate, 10, 64)
		query = query.Where("create_time <= ?", endTime)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	offset := (req.Current - 1) * req.PageSize

	var tasks []model.ManualCheckTask
	if err := query.Offset(offset).Limit(req.PageSize).Order("create_time DESC").Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("查询任务列表失败: %w", err)
	}

	// 转换为响应格式
	var items []dto.ManualCheckHistoryItem
	for _, task := range tasks {
		statusText := "进行中"
		switch task.Status {
		case 2:
			statusText = "完成"
		case 3:
			statusText = "失败"
		}

		items = append(items, dto.ManualCheckHistoryItem{
			ID:          task.ID,
			ProjectID:   task.ProjectID,
			ProjectName: task.ProjectName,
			MergeID:     task.MergeID,
			MergeTitle:  task.MergeTitle,
			Status:      int(task.Status),
			StatusText:  statusText,
			CreateTime:  task.CreateTime,
			UpdateTime:  task.UpdateTime,
		})
	}

	return &dto.ManualCheckHistoryResponse{
		List:  items,
		Total: total,
	}, nil
}

// GetManualCheckResult 获取手动审核结果详情
func (s *ManualCheckService) GetManualCheckResult(taskID uint, userID uint) (*dto.ManualCheckResultResponse, error) {
	var task model.ManualCheckTask
	if err := s.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		return nil, fmt.Errorf("任务不存在或无权限访问: %w", err)
	}

	// 反序列化机器人配置
	var botConfigs []dto.BotConfig
	if task.BotConfigs != "" {
		if err := json.Unmarshal([]byte(task.BotConfigs), &botConfigs); err != nil {
			// 如果反序列化失败，返回空列表
			botConfigs = []dto.BotConfig{}
		}
	}

	return &dto.ManualCheckResultResponse{
		ID:          task.ID,
		ProjectID:   task.ProjectID,
		ProjectName: task.ProjectName,
		MergeID:     task.MergeID,
		MergeTitle:  task.MergeTitle,
		MergeURL:    task.MergeURL,
		Status:      int(task.Status),
		Result:      task.Result,
		AIModel:     task.AIModel,
		BotConfigs:  botConfigs,
		RuleName:    task.RuleName,
		CreateTime:  task.CreateTime,
		UpdateTime:  task.UpdateTime,
	}, nil
}

// GetManualCheckService 获取手动审核服务实例
func GetManualCheckService() *ManualCheckService {
	return NewManualCheckService(database.DB, gitlab_service.NewGitlabService(database.DB), ai.NewAIRuleService(database.DB))
}

// parseMergeURL 解析合并请求URL，提取项目ID、合并请求ID和项目名称
func (s *ManualCheckService) parseMergeURL(mergeURL string) (int, int, string, error) {
	// 支持的GitLab URL格式:
	// 1. https://gitlab.example.com/project/repo/-/merge_requests/123
	// 2. https://gitlab.example.com/group/project/-/merge_requests/456
	// 3. http://165.1.1.72:9980/usms/api/-/merge_requests/3
	// 4. https://gitlab.com/username/project/-/merge_requests/789

	// 使用正则表达式匹配合并请求URL格式
	// 匹配模式: /-/merge_requests/数字
	re := regexp.MustCompile(`/-/merge_requests/(\d+)(?:/.*)?$`)
	matches := re.FindStringSubmatch(mergeURL)
	if len(matches) < 2 {
		return 0, 0, "", fmt.Errorf("无效的合并请求URL格式，请确保URL包含 /-/merge_requests/数字")
	}

	// 提取合并请求ID
	mergeID, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, "", fmt.Errorf("无法解析合并请求ID: %v", err)
	}

	// 提取项目路径
	// 从URL中提取项目路径部分，用于获取真正的项目ID和项目名称
	projectPath := s.extractProjectPath(mergeURL)
	if projectPath == "" {
		return 0, 0, "", fmt.Errorf("无法从URL中提取项目路径")
	}

	// 尝试通过GitLab API获取真正的项目ID
	projectID, projectName, err := s.getRealProjectID(projectPath, mergeURL)
	if err != nil {
		// 如果获取失败，回退到使用哈希值生成项目ID
		fmt.Printf("获取真正项目ID失败，使用哈希值: %v\n", err)
		projectID = s.generateProjectID(projectPath)
		projectName = s.generateProjectName(projectPath)
	}

	return projectID, mergeID, projectName, nil
}

// extractProjectPath 从GitLab URL中提取项目路径
func (s *ManualCheckService) extractProjectPath(mergeURL string) string {
	// 移除协议和域名部分
	// 例如: https://gitlab.example.com/group/project/-/merge_requests/123
	// 提取: group/project
	// 例如: http://127.0.0.1:9980/usms/api/-/merge_requests/3
	// 提取: usms/api

	// 使用正则表达式提取项目路径
	// 匹配模式: 协议://域名/项目路径/-/merge_requests/
	re := regexp.MustCompile(`^https?://[^/]+/(.+?)/-/merge_requests/`)
	matches := re.FindStringSubmatch(mergeURL)
	if len(matches) < 2 {
		return ""
	}

	projectPath := matches[1]

	// 确保项目路径不为空
	if projectPath == "" {
		return ""
	}

	return projectPath
}

// generateProjectID 根据项目路径生成项目ID
func (s *ManualCheckService) generateProjectID(projectPath string) int {
	// 使用项目路径的哈希值作为项目ID
	// 这样可以确保相同的项目路径总是生成相同的ID
	hash := fnv.New32a()
	hash.Write([]byte(projectPath))

	// 将哈希值转换为正整数
	projectID := int(hash.Sum32())
	if projectID < 0 {
		projectID = -projectID
	}

	// 确保ID不为0
	if projectID == 0 {
		projectID = 1
	}

	return projectID
}

// generateProjectName 根据项目路径生成项目名称
func (s *ManualCheckService) generateProjectName(projectPath string) string {
	// 将项目路径转换为更友好的项目名称
	// 例如: "usms/api" -> "usms/api"
	// 例如: "group/project" -> "group/project"
	// 例如: "username/project" -> "username/project"

	// 如果项目路径为空，返回默认名称
	if projectPath == "" {
		return "unknown-project"
	}

	// 直接使用项目路径作为项目名称
	// 可以根据需要进一步处理，比如替换特殊字符等
	return projectPath
}

// getRealProjectID 通过GitLab API获取真正的项目ID
func (s *ManualCheckService) getRealProjectID(projectPath, mergeURL string) (int, string, error) {
	// 从URL中提取GitLab实例信息
	gitlabAPI, err := s.extractGitlabAPI(mergeURL)
	if err != nil {
		return 0, "", fmt.Errorf("无法提取GitLab API地址: %v", err)
	}

	// 获取GitLab配置信息
	gitlabInfoList, err := s.gitlabService.GetGitlabInfo()
	if err != nil || len(gitlabInfoList) == 0 {
		return 0, "", fmt.Errorf("无法获取GitLab配置信息: %v", err)
	}

	// 查找匹配的GitLab配置
	var matchedConfig *dto.GitlabInfoResponse
	var matchedToken string
	for _, config := range gitlabInfoList {
		if config.Status == 1 && config.API != "" {
			// 检查API地址是否匹配（去掉协议前缀进行比较）
			configAPI := strings.TrimPrefix(config.API, "http://")
			configAPI = strings.TrimPrefix(configAPI, "https://")
			mergeAPI := strings.TrimPrefix(gitlabAPI, "http://")
			mergeAPI = strings.TrimPrefix(mergeAPI, "https://")

			if strings.Contains(mergeAPI, configAPI) || strings.Contains(configAPI, mergeAPI) {
				matchedConfig = &config
				// 从数据库获取完整的配置信息以获取Token
				var fullConfig model.GitlabInfo
				if err := s.db.First(&fullConfig, config.ID).Error; err == nil {
					matchedToken = fullConfig.Token
				}
				break
			}
		}
	}

	if matchedConfig == nil {
		return 0, "", fmt.Errorf("未找到匹配的GitLab配置")
	}

	if matchedToken == "" {
		return 0, "", fmt.Errorf("无法获取GitLab Token")
	}

	// 通过项目路径获取项目信息
	projectID, projectName, err := s.fetchProjectByPath(projectPath, matchedConfig.API, matchedToken)
	if err != nil {
		return 0, "", fmt.Errorf("获取项目信息失败: %v", err)
	}

	return projectID, projectName, nil
}

// extractGitlabAPI 从合并请求URL中提取GitLab API地址
func (s *ManualCheckService) extractGitlabAPI(mergeURL string) (string, error) {
	// 使用正则表达式提取协议和域名部分
	re := regexp.MustCompile(`^(https?://[^/]+)`)
	matches := re.FindStringSubmatch(mergeURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法从URL中提取GitLab API地址")
	}
	return matches[1], nil
}

// fetchProjectByPath 通过项目路径获取项目ID和名称
func (s *ManualCheckService) fetchProjectByPath(projectPath, gitlabAPI, token string) (int, string, error) {
	// 使用GitLab API通过项目路径获取项目信息
	// URL编码项目路径
	encodedPath := strings.ReplaceAll(projectPath, "/", "%2F")
	url := fmt.Sprintf("%s/v4/projects/%s", gitlabAPI, encodedPath)

	// 发送请求
	body, err := utils.CommonGetRequest("GET", url, token, nil)
	if err != nil {
		return 0, "", fmt.Errorf("请求GitLab API失败: %v", err)
	}

	// 解析响应
	var projectInfo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Path string `json:"path"`
	}

	if err := json.Unmarshal(body, &projectInfo); err != nil {
		return 0, "", fmt.Errorf("解析项目信息失败: %v", err)
	}

	// 使用项目名称，如果没有则使用路径
	projectName := projectInfo.Name
	if projectName == "" {
		projectName = projectInfo.Path
	}
	if projectName == "" {
		projectName = projectPath
	}

	return projectInfo.ID, projectName, nil
}

// getAIModelByID 根据ID获取AI模型信息
func (s *ManualCheckService) getAIModelByID(modelID uint) (*model.AIConfig, error) {
	var aiConfig model.AIConfig
	if err := s.db.First(&aiConfig, modelID).Error; err != nil {
		return nil, fmt.Errorf("AI模型不存在: %w", err)
	}

	if aiConfig.IsActive != 1 {
		return nil, fmt.Errorf("AI模型未启用")
	}

	return &aiConfig, nil
}

// updateTaskWithCachedInfo 异步更新任务中的项目信息和AI模型信息
func (s *ManualCheckService) updateTaskWithCachedInfo(taskID uint, projectID int, projectName string, aiModelID uint) {
	updateData := map[string]interface{}{
		"update_time": time.Now().Unix(),
	}

	// 获取AI模型信息
	aiConfig, err := s.getAIModelByID(aiModelID)
	if err == nil {
		updateData["ai_model"] = aiConfig.Model
	} else {
		fmt.Printf("获取AI模型信息失败: %v\n", err)
	}

	// 获取Gitlab配置信息
	gitlabInfoList, err := s.gitlabService.GetGitlabInfo()
	if err == nil && len(gitlabInfoList) > 0 {
		// 这里可以根据项目ID匹配对应的Gitlab配置
		// 简化处理：使用第一个配置
		_ = gitlabInfoList[0] // 暂时不使用，后续可以根据需要扩展

		// 使用解析出的项目名称，如果没有则生成默认名称
		if projectName == "" {
			projectName = fmt.Sprintf("project-%d", projectID)
		}

		// 生成合并请求标题
		mergeTitle := fmt.Sprintf("Merge Request #%d", projectID) // 使用projectID作为mergeID的简化处理

		updateData["project_name"] = projectName
		updateData["merge_title"] = mergeTitle
	} else {
		fmt.Printf("获取Gitlab配置信息失败: %v\n", err)
		// 即使获取Gitlab配置失败，也使用解析出的项目名称
		if projectName != "" {
			updateData["project_name"] = projectName
		}
	}

	// 更新任务信息
	if err := s.db.Model(&model.ManualCheckTask{}).Where("id = ?", taskID).Updates(updateData).Error; err != nil {
		fmt.Printf("更新任务信息失败: %v\n", err)
	}
}
