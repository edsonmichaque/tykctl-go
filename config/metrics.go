package config

import (
	"time"
)

// MetricsOptions provides configuration for metrics
type MetricsOptions struct {
	Extension string
}

// NewMetrics creates a new metrics instance
func NewMetrics(opts MetricsOptions) (Metrics, error) {
	return &noopMetrics{}, nil
}

// noopMetrics implements a no-op metrics collector
type noopMetrics struct{}

func (m *noopMetrics) RecordLoadDuration(duration time.Duration) {
	// No-op implementation
}

func (m *noopMetrics) RecordDiscoveryDuration(resourceType string, duration time.Duration) {
	// No-op implementation
}

func (m *noopMetrics) Close() error {
	return nil
}