package ai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
)

// AIMessageService AI消息服务
type AIMessageService struct {
	db *gorm.DB
}

func NewAIMessageService(db *gorm.DB) *AIMessageService {
	return &AIMessageService{
		db: db,
	}
}

// GetAIMessage 获取AI消息列表
func (s *AIMessageService) GetAIMessage(params map[string]interface{}) ([]dto.AIMessageItem, int64, error) {
	var messages []model.AIMessage
	var total int64

	query := s.db.Model(&model.AIMessage{})

	// 基础筛选条件
	if id, ok := params["id"].(uint); ok && id > 0 {
		query = query.Where("id = ?", id)
	}
	if passed, ok := params["passed"].(int); ok {
		query = query.Where("passed = ?", passed)
	}

	// 项目ID筛选
	if projectID, ok := params["project_id"].(uint); ok && projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}

	// 项目命名空间筛选
	if projectNamespace, ok := params["project_namespace"].(string); ok && projectNamespace != "" {
		query = query.Where("project_namespace = ?", projectNamespace)
	}

	// 时间范围筛选（时间戳格式）
	if startDateVal, exists := params["start_date"]; exists && startDateVal != nil {
		if startDate, ok := startDateVal.(int64); ok && startDate > 0 {
			query = query.Where("create_time >= ?", startDate)
		}
	}
	if endDateVal, exists := params["end_date"]; exists && endDateVal != nil {
		if endDate, ok := endDateVal.(int64); ok && endDate > 0 {
			query = query.Where("create_time <= ?", endDate)
		}
	}

	// 项目ID多选筛选
	if projectIDsVal, exists := params["project_ids"]; exists && projectIDsVal != nil {
		if projectIDs, ok := projectIDsVal.([]uint); ok && len(projectIDs) > 0 {
			query = query.Where("project_id IN ?", projectIDs)
		}
	}

	// 人工评分多选筛选
	if humanRatingsVal, exists := params["human_ratings"]; exists && humanRatingsVal != nil {
		if humanRatings, ok := humanRatingsVal.([]int8); ok && len(humanRatings) > 0 {
			query = query.Where("human_rating IN ?", humanRatings)
		}
	}

	// 项目命名空间多选筛选
	if projectNamespacesVal, exists := params["project_namespaces"]; exists && projectNamespacesVal != nil {
		if projectNamespaces, ok := projectNamespacesVal.([]string); ok && len(projectNamespaces) > 0 {
			query = query.Where("project_namespace IN ?", projectNamespaces)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if offset, ok := params["offset"].(int); ok {
		query = query.Offset(offset)
	}
	if limit, ok := params["limit"].(int); ok {
		query = query.Limit(limit)
	}

	// 按创建时间倒序排序
	query = query.Order("create_time DESC")

	if err := query.Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应模型
	responseList := make([]dto.AIMessageItem, len(messages))
	for i, msg := range messages {
		responseList[i] = dto.AIMessageItem{
			ID:               msg.ID,
			ProjectID:        msg.ProjectID,
			ProjectName:      msg.ProjectName,
			ProjectNamespace: msg.ProjectNamespace,
			MergeDescription: msg.MergeDescription,
			MergeURL:         msg.MergeURL,
			MergeID:          uint(parseUint(msg.MergeID)),
			AIModel:          msg.AIModel,
			Rule:             int8(msg.Rule),
			RuleID:           msg.RuleID,
			Result:           msg.Result,
			HumanRating:      int8(msg.HumanRating),
			Remark:           msg.Remark,
			Passed:           msg.Passed,
			CheckedBy:        msg.CheckedBy,
			Status:           int8(msg.HumanRating),
			CreateTime:       msg.CreateTime,
			UpdateTime:       msg.CreateTime,
		}
	}

	return responseList, total, nil
}

// UpdateHumanRatingAndRemark 更新人工评分和备注
func (s *AIMessageService) UpdateHumanRatingAndRemark(data model.AIMessageUpdate) error {
	now := time.Now().Unix()
	return s.db.Model(&model.AIMessage{}).
		Where("id = ?", data.ID).
		Updates(map[string]interface{}{
			"human_rating": data.HumanRating,
			"remark":       data.Remark,
			"update_time":  now,
		}).Error
}

// parseUint 将字符串转换为 uint64，如果转换失败返回 0
func parseUint(s string) uint64 {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return u
}

// CreateAIMessage 创建AI消息
func (s *AIMessageService) CreateAIMessage(data *model.AIMessage) (uint, error) {
	if err := s.db.Create(data).Error; err != nil {
		return 0, err
	}
	return data.ID, nil
}

// SendMarkdownToWechatBot 向企业微信机器人发送 Markdown 消息
func SendMarkdownToWechatBot(webhookURL, markdownContent string) error {
	// 构造请求体
	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": markdownContent,
		},
	}

	// 发送 POST 请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(utils.MustMarshal(payload)))
	if err != nil {
		return fmt.Errorf("发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送失败: %s", string(body))
	}

	return nil
}

func (s *AIMessageService) GetCheckCount(req dto.CheckCountRequest) (int64, error) {
	query := s.db.Model(&model.AIMessage{})

	// 只有传入时才加条件
	if req.Passed != 0 {
		query = query.Where("passed = ?", req.Passed)
	}

	if req.ProjectID > 0 {
		query = query.Where("project_id = ?", req.ProjectID)
	}
	if req.ProjectNamespace != "" {
		query = query.Where("project_namespace = ?", req.ProjectNamespace)
	}
	if req.ProjectName != "" {
		query = query.Where("project_name = ?", req.ProjectName)
	}

	query = query.Where("create_time >= ?", req.StartTime)
	query = query.Where("create_time <= ?", req.EndTime)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (s *AIMessageService) GetProblemChart(req dto.AIProblemCountRequest) ([]dto.HumanRatingStat, error) {
	query := s.db.Model(&model.AIMessage{})

	query = query.Where("create_time >= ?", req.StartTime)
	query = query.Where("create_time <= ?", req.EndTime)

	// 统计各等级数量
	type result struct {
		Level int   `json:"level"`
		Count int64 `json:"count"`
	}
	var results []result

	err := query.Select("human_rating as level, COUNT(*) as count").
		Where("human_rating BETWEEN ? AND ?", 1, 5).
		Group("human_rating").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// 保证每个level都有（即使为0）
	stats := make(map[int]int64)
	for i := 1; i <= 5; i++ {
		stats[i] = 0
	}
	for _, r := range results {
		stats[r.Level] = r.Count
	}

	// 返回数组
	resp := make([]dto.HumanRatingStat, 0, 5)
	for i := 1; i <= 5; i++ {
		resp = append(resp, dto.HumanRatingStat{
			Level: i,
			Count: stats[i],
		})
	}

	return resp, nil
}

// GetProjectNamespaceList 获取项目命名空间列表（去重）
// 使用DISTINCT查询优化性能，避免全表扫描
// 支持时间段筛选，最大允许查询31天的数据
func (s *AIMessageService) GetProjectNamespaceList(startTime, endTime int64) ([]dto.ProjectNamespaceItem, error) {
	var results []struct {
		ProjectID        uint   `json:"project_id"`
		ProjectNamespace string `json:"project_namespace"`
	}

	query := s.db.Model(&model.AIMessage{})

	// 添加时间段筛选条件
	if startTime > 0 {
		query = query.Where("create_time >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("create_time <= ?", endTime)
	}

	// 使用DISTINCT查询，只选择需要的字段，提高查询效率
	err := query.Select("DISTINCT project_id, project_namespace").
		Order("project_id ASC, project_namespace ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应模型
	responseList := make([]dto.ProjectNamespaceItem, len(results))
	for i, result := range results {
		responseList[i] = dto.ProjectNamespaceItem{
			ProjectID:        result.ProjectID,
			ProjectNamespace: result.ProjectNamespace,
		}
	}

	return responseList, nil
}
