package service

// RAGService 定义RAG服务的接口
type RAGService interface {
	LegacyAnalyzeCode(gitURL string, branch string, diffFiles []string, codeChanges string) (string, error)
}
