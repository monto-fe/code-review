package dto

type GitlabToken struct {
	ID        uint   `json:"id"`
	Namespace string `json:"namespace"`
	Token     string `json:"token"`
}
