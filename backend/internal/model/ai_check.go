package model

// MergeRequestInfo 合并请求信息

type DiffRefs struct {
	BaseSHA  string `json:"base_sha"`
	HeadSHA  string `json:"head_sha"`
	StartSHA string `json:"start_sha"`
}
type MergeRequestInfo struct {
	ID           int    `json:"id"`
	IID          int    `json:"iid"`
	ProjectID    uint   `json:"project_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	State        string `json:"state"`
	WebURL       string `json:"web_url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
	Sha          string `json:"sha"`
	User         struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Project struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		Description       string `json:"description"`
		PathWithNamespace string `json:"path_with_namespace"`
	} `json:"project"`
	ObjectAttributes struct {
		Action       string `json:"action"`
		State        string `json:"state"`
		SourceBranch string `json:"source_branch"`
		TargetBranch string `json:"target_branch"`
	} `json:"object_attributes"`
	DiffRefs DiffRefs `json:"diff_refs"`
}

// Change 文件变更信息
type Change struct {
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	Diff        string `json:"diff"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
}

// AICheck AI检查记录
type AICheck struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	ProjectID  uint   `gorm:"index;not null" json:"project_id"`
	MergeID    uint   `gorm:"index;not null" json:"merge_id"`
	Result     string `gorm:"type:text;not null" json:"result"`
	Status     int8   `gorm:"default:1" json:"status"` // 1: 正常, -1: 删除
	CreateTime int64  `gorm:"not null" json:"create_time"`
	UpdateTime int64  `gorm:"not null" json:"update_time"`
}

// ManualCheckTask 手动审核任务
type ManualCheckTask struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	UserID       uint   `gorm:"index;not null" json:"user_id"`          // 触发用户ID
	ProjectID    int    `gorm:"index;not null" json:"project_id"`       // 项目ID
	MergeID      int    `gorm:"index;not null" json:"merge_id"`         // 合并请求ID
	ProjectName  string `gorm:"type:varchar(200)" json:"project_name"`  // 项目名称
	MergeTitle   string `gorm:"type:varchar(500)" json:"merge_title"`   // 合并请求标题
	MergeURL     string `gorm:"type:varchar(500)" json:"merge_url"`     // 合并请求URL
	Status       int8   `gorm:"default:1" json:"status"`                // 状态: 1-进行中, 2-完成, 3-失败
	Result       string `gorm:"type:text" json:"result"`                // 审核结果
	ErrorMessage string `gorm:"type:varchar(500)" json:"error_message"` // 错误信息
	AIModelID    uint   `gorm:"not null" json:"ai_model_id"`            // AI模型ID
	AIModel      string `gorm:"type:varchar(50)" json:"ai_model"`       // AI模型名称
	BotConfigs   string `gorm:"type:text" json:"bot_configs"`           // 机器人配置列表(JSON格式)
	RuleID       uint   `gorm:"default:0" json:"rule_id"`               // 规则ID
	RuleName     string `gorm:"type:varchar(100)" json:"rule_name"`     // 规则名称
	CreateTime   int64  `gorm:"not null" json:"create_time"`
	UpdateTime   int64  `gorm:"not null" json:"update_time"`
}

// TableName 指定表名
func (ManualCheckTask) TableName() string {
	return "t_manual_check_task"
}
