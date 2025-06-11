package model

// MergeRequestInfo 合并请求信息
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
