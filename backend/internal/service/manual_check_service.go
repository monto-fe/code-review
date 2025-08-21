package service

import (
	"code-review-go/internal/database"
	"code-review-go/internal/dto"
	"code-review-go/internal/model"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ManualCheckService struct {
	db            *gorm.DB
	gitlabService *GitlabService
	aiRuleService *AIRuleService
}

func NewManualCheckService(db *gorm.DB, gitlabService *GitlabService, aiRuleService *AIRuleService) *ManualCheckService {
	return &ManualCheckService{
		db:            db,
		gitlabService: gitlabService,
		aiRuleService: aiRuleService,
	}
}

// CreateManualCheckTask 创建手动审核任务
func (s *ManualCheckService) CreateManualCheckTask(req dto.ManualCheckRequest, userID uint) (*model.ManualCheckTask, error) {
	// 获取Gitlab配置信息
	gitlabInfoList, err := s.gitlabService.GetGitlabInfo()
	if err != nil || len(gitlabInfoList) == 0 {
		return nil, fmt.Errorf("未找到Gitlab配置信息")
	}
	gitlabInfo := gitlabInfoList[0] // 使用第一个配置

	// 验证项目是否存在 - 这里简化处理，实际应该调用Gitlab API
	projectName := fmt.Sprintf("project-%d", req.ProjectID) // 临时项目名称

	// 验证合并请求是否存在 - 这里简化处理，实际应该调用Gitlab API
	mergeTitle := fmt.Sprintf("Merge Request #%d", req.MergeID) // 临时标题
	mergeURL := fmt.Sprintf("%s/project/%d/merge_requests/%d", gitlabInfo.GitlabURL, req.ProjectID, req.MergeID)

	// 获取AI模型和规则信息
	aiModel := req.AIModel
	if aiModel == "" {
		aiModel = "gpt-4" // 默认模型
	}

	ruleName := ""
	if req.RuleID > 0 {
		// 这里应该调用规则服务获取规则名称，暂时简化处理
		ruleName = fmt.Sprintf("Rule-%d", req.RuleID)
	}

	// 创建任务记录
	now := time.Now().Unix()
	task := &model.ManualCheckTask{
		UserID:      userID,
		ProjectID:   req.ProjectID,
		MergeID:     req.MergeID,
		ProjectName: projectName,
		MergeTitle:  mergeTitle,
		MergeURL:    mergeURL,
		Status:      1, // 进行中
		AIModel:     aiModel,
		RuleID:      req.RuleID,
		RuleName:    ruleName,
		CreateTime:  now,
		UpdateTime:  now,
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

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
	result, err := CheckMergeRequestWithAI(webhookBody)

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
		RuleName:    task.RuleName,
		CreateTime:  task.CreateTime,
		UpdateTime:  task.UpdateTime,
	}, nil
}

// GetManualCheckService 获取手动审核服务实例
func GetManualCheckService() *ManualCheckService {
	return NewManualCheckService(database.DB, NewGitlabService(database.DB), NewAIRuleService(database.DB))
}
