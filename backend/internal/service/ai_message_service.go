package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	dto "code-review-go/internal/dto"
	"code-review-go/internal/model"
	"code-review-go/internal/pkg/utils"
)

// AIMessageService AI消息服务
type AIMessageService struct {
	db         *gorm.DB
	ragService RAGService
}

func NewAIMessageService(db *gorm.DB, ragService RAGService) *AIMessageService {
	return &AIMessageService{
		db:         db,
		ragService: ragService,
	}
}

// GetAIMessage 获取AI消息列表
func (s *AIMessageService) GetAIMessage(params map[string]interface{}) ([]dto.AIMessageItem, int64, error) {
	var messages []model.AImessage
	var total int64

	query := s.db.Model(&model.AImessage{})

	// 添加查询条件
	if projectID, ok := params["projectId"].(uint); ok && projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if id, ok := params["id"].(uint); ok && id > 0 {
		query = query.Where("id = ?", id)
	}
	if projectNamespace, ok := params["projectNamespace"].(string); ok && projectNamespace != "" {
		query = query.Where("project_namespace = ?", projectNamespace)
	}
	if projectName, ok := params["projectName"].(string); ok && projectName != "" {
		query = query.Where("project_name = ?", projectName)
	}
	if createTime, ok := params["createTime"].(int64); ok && createTime > 0 {
		query = query.Where("create_time >= ?", createTime)
	}
	if passed, ok := params["passed"].(int); ok {
		query = query.Where("passed = ?", passed)
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
func (s *AIMessageService) CreateAIMessage(data *model.AImessage) (uint, error) {
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

// AnalyzeCodeWithRAG 使用RAG服务分析代码
func (s *AIMessageService) AnalyzeCodeWithRAG(projectID uint, mergeURL string, diffFiles []string) (string, error) {
	// 从mergeURL中提取gitURL和branch
	gitURL, branch, err := extractGitInfo(mergeURL)
	if err != nil {
		return "", fmt.Errorf("解析Git信息失败: %v", err)
	}

	// 获取代码变更内容
	codeChanges, err := s.getCodeChanges(gitURL, branch, diffFiles)
	if err != nil {
		return "", fmt.Errorf("获取代码变更失败: %v", err)
	}

	// 调用RAG服务分析代码
	review, err := s.ragService.LegacyAnalyzeCode(gitURL, branch, diffFiles, codeChanges)
	if err != nil {
		return "", fmt.Errorf("RAG分析失败: %v", err)
	}

	return review, nil
}

// getCodeChanges 获取代码变更内容
func (s *AIMessageService) getCodeChanges(gitURL string, branch string, diffFiles []string) (string, error) {
	// TODO: 实现从GitLab API获取代码变更的逻辑
	// 这里需要调用GitLab API获取代码变更内容
	// 返回格式化的代码变更字符串
	return "", nil
}

// extractGitInfo 从mergeURL中提取gitURL和branch
func extractGitInfo(mergeURL string) (string, string, error) {
	// 示例: https://gitlab.com/group/project/-/merge_requests/123
	parts := strings.Split(mergeURL, "/-/merge_requests/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid merge URL format")
	}

	baseURL := parts[0]
	gitURL := fmt.Sprintf("%s.git", baseURL)

	// 获取当前分支信息
	// 这里需要根据实际情况实现，可能需要调用Gitlab API
	branch := "main" // 默认分支

	return gitURL, branch, nil
}

func (s *AIMessageService) GetCheckCount(req dto.CheckCountRequest) (int64, error) {
	query := s.db.Model(&model.AImessage{})

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
	query := s.db.Model(&model.AImessage{})

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
