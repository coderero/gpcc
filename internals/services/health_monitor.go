package services

import (
	"connector/internals/types"
	"context"
	"sync"
	"time"
)

type HealthMonitor interface {
	// The `Start(context.Context) error` method in the `healthMonitor` struct is implementing the
	// `HealthMonitor` interface. This method starts a health monitoring process that runs periodically
	// based on the specified interval. It uses a ticker to trigger the health check function at regular
	// intervals. If the context is canceled (ctx.Done()), the method will return nil, indicating that the
	// health monitoring should stop.
	Start(context.Context) error
}

type healthMonitor struct {
	healthChan chan types.HealthStatus
	checkFunc  func() bool
	interval   time.Duration
	mutex      sync.RWMutex
	lastStatus types.HealthStatus
}

// The NewHealthMonitor function creates a new HealthMonitor instance with a specified check function
// and interval.
func NewHealthMonitor(checkFunc func() bool, interval time.Duration) HealthMonitor {
	return &healthMonitor{
		healthChan: make(chan types.HealthStatus),
		checkFunc:  checkFunc,
		interval:   interval,
	}
}

func (h *healthMonitor) Start(ctx context.Context) error {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			h.performHealthCheck()
		}
	}
}

func (h *healthMonitor) performHealthCheck() {
	status := types.HealthStatus{
		IsHealthy:   h.checkFunc(),
		LastChecked: time.Now(),
	}

	h.mutex.Lock()
	h.lastStatus = status
	h.mutex.Unlock()

	select {
	case h.healthChan <- status:
	default:
		// If the channel is full, drop the status
		<-h.healthChan
		h.healthChan <- status
	}
}
