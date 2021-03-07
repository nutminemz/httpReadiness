package model

type HealthRequest struct {
	TotalWebsites int   `json:"total_websites"`
	Success       int   `json:"success"`
	Failure       int   `json:"failure"`
	TotalTime     int64 `json:"total_time"`
}

type HealthResponse struct {
	Message *string `json:"message"`
}
