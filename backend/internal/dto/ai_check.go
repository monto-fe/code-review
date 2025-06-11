package dto

type ProjectInfo struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PathWithNamespace string `json:"path_with_namespace"`
}

type ObjectAttributes struct {
	IID               int    `json:"iid"`
	URL               string `json:"url"`
	Action            string `json:"action"`
	ProjectID         int    `json:"project_id"`
	MergeURL          string `json:"merge_url"`
	State             string `json:"state"`
	SourceBranch      string `json:"source_branch"`
	TargetBranch      string `json:"target_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
}

type WebhookBody struct {
	Project          ProjectInfo      `json:"project"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
}
