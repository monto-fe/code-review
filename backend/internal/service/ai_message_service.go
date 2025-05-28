package service

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
	return &AIMessageService{db: db}
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
			ID:          msg.ID,
			ProjectID:   msg.ProjectID,
			MergeURL:    msg.MergeURL,
			MergeID:     uint(parseUint(msg.MergeID)),
			AIModel:     msg.AIModel,
			Rule:        int8(msg.Rule),
			RuleID:      msg.RuleID,
			Result:      msg.Result,
			HumanRating: int8(msg.HumanRating),
			Remark:      msg.Remark,
			Passed:      msg.Passed,
			CheckedBy:   msg.CheckedBy,
			Status:      int8(msg.HumanRating),
			CreateTime:  msg.CreateTime,
			UpdateTime:  msg.CreateTime,
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
