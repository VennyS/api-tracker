package models

import "time"

type APIRequestLog struct {
	Timestamp   time.Time `json:"timestamp"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	StatusCode  int       `json:"status_code"`
	LatencyMs   float64   `json:"latency_ms"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	ServiceName string    `json:"service_name"`
}
