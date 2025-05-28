package dto

type AiConfig struct {
	ID        uint   `json:"id"`
	Namespace string `json:"namespace"`
	Config    string `json:"config"`
}
