package dto

type AIManagerListResponse struct {
	ID         uint   `json:"id"`
	Type       string `json:"type"`
	Model      string `json:"model"`
	APIURL     string `json:"api_url"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
