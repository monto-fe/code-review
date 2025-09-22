package dto

// CodeReviewRequest 代码审查请求
type CodeReviewRequest struct {
	GitURL      string `json:"git_url"`
	Branch      string `json:"branch"`
	DiffContent string `json:"diff_content"`
	Query       string `json:"query,omitempty"`
	GitlabToken string `json:"gitlab_token,omitempty"`
}

// CodeAnalysisResponse 代码分析响应
type CodeAnalysisResponse struct {
	Review string `json:"review"`
}

// CodeReviewResponse 代码审查响应（别名）
type CodeReviewResponse = CodeAnalysisResponse

// LineComment 行级评论
type LineComment struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// RAGService 定义RAG服务接口
type RAGService interface {
	AnalyzeCodeWithRequest(req *CodeReviewRequest) (*CodeAnalysisResponse, error)
	GenerateReview(req *CodeReviewRequest) (string, error)
}
