package v2ray

import (
	"github.com/v2rayA/v2rayA/db/configure"
)

// LatencyProgressStatus represents the status of a latency test
type LatencyProgressStatus string

const (
	LatencyStatusPending LatencyProgressStatus = "pending"
	LatencyStatusRunning LatencyProgressStatus = "running"
	LatencyStatusDone    LatencyProgressStatus = "done"
	LatencyStatusError   LatencyProgressStatus = "error"
)

// LatencyProgress represents a latency test progress update
type LatencyProgress struct {
	TestType string                `json:"testType"` // "http" or "ping"
	Which    *configure.Which      `json:"which"`
	Status   LatencyProgressStatus `json:"status"`
	Latency  string                `json:"latency,omitempty"`
}

// BroadcastLatencyProgress sends a latency test progress update to all WebSocket clients
func BroadcastLatencyProgress(testType string, which *configure.Which, status LatencyProgressStatus, latency string) {
	if ApiFeed == nil {
		return
	}
	progress := LatencyProgress{
		TestType: testType,
		Which:    which,
		Status:   status,
		Latency:  latency,
	}
	ApiFeed.ProductMessage("latencyProgress", progress)
}
