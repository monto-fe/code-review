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
	Title             string `json:"title"`
	Description       string `json:"description"`
	Action            string `json:"action"`
	ProjectID         int    `json:"project_id"`
	MergeURL          string `json:"merge_url"`
	State             string `json:"state"`
	SourceBranch      string `json:"source_branch"`
	TargetBranch      string `json:"target_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
	Note              string `json:"note,omitempty"`
}

type WebhookBody struct {
	ObjectKind       string           `json:"object_kind,omitempty"`
	Project          ProjectInfo      `json:"project"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
	MergeRequest     MergeRequest     `json:"merge_request,omitempty"`
}

type MergeRequest struct {
	State string `json:"state"`
}
