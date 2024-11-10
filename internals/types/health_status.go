package types

import "time"

type HealthStatus struct {
	IsHealthy   bool      `json:"is_healthy"`
	LastChecked time.Time `json:"last_checked"`
}
