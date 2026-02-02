package model

type TelemetryPayload struct {
	Timestamp string  `json:"timestamp"`
	Route     string  `json:"route"`
	UserAgent string  `json:"userAgent"`
	IP        *string `json:"ip,omitempty"`
	Error     *string `json:"error,omitempty"`
}
