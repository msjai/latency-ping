package entity

import "time"

// WebSite .-
type WebSite struct {
	Name    string         `json:"Name"`
	Address string         `json:"Address"`
	Latency *time.Duration `json:"Latency"`
}
