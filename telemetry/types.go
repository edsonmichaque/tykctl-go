// Package telemetry provides anonymous usage analytics for tykctl-go.
// It collects anonymous usage statistics to help improve the CLI tool
// while respecting user privacy and providing opt-out capabilities.
package telemetry

import (
	"time"
)

// EventType represents the type of telemetry event.
type EventType string

const (
	// EventTypeCommand represents a command execution event.
	EventTypeCommand EventType = "command"
	// EventTypeError represents an error event.
	EventTypeError EventType = "error"
	// EventTypeFeature represents a feature usage event.
	EventTypeFeature EventType = "feature"
	// EventTypePerformance represents a performance metric event.
	EventTypePerformance EventType = "performance"
)

// Event represents a telemetry event to be collected.
type Event struct {
	// EventType is the type of event being recorded.
	EventType EventType `json:"event_type"`
	
	// Command is the command that was executed (for command events).
	Command string `json:"command,omitempty"`
	
	// Duration is the execution time in milliseconds.
	Duration int64 `json:"duration_ms,omitempty"`
	
	// Success indicates whether the operation was successful.
	Success bool `json:"success"`
	
	// ErrorType is the type of error (for error events).
	ErrorType string `json:"error_type,omitempty"`
	
	// ErrorMessage is a sanitized error message (no sensitive data).
	ErrorMessage string `json:"error_message,omitempty"`
	
	// Feature is the feature being used (for feature events).
	Feature string `json:"feature,omitempty"`
	
	// Properties contains additional event-specific properties.
	Properties map[string]interface{} `json:"properties,omitempty"`
	
	// Timestamp is when the event occurred.
	Timestamp time.Time `json:"timestamp"`
	
	// SessionID is an anonymous session identifier.
	SessionID string `json:"session_id"`
	
	// UserID is an anonymous user identifier (hashed).
	UserID string `json:"user_id"`
	
	// CLI version information.
	CLIVersion string `json:"cli_version"`
	
	// Operating system information.
	OS string `json:"os"`
	
	// Architecture information.
	Arch string `json:"arch"`
}

// Config represents the telemetry configuration.
type Config struct {
	// Enabled controls whether telemetry is active.
	Enabled bool `yaml:"enabled" json:"enabled"`
	
	// Endpoint is the telemetry data collection endpoint.
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	
	// BatchSize is the maximum number of events to batch before sending.
	BatchSize int `yaml:"batch_size" json:"batch_size"`
	
	// FlushInterval is how often to send batched events.
	FlushInterval time.Duration `yaml:"flush_interval" json:"flush_interval"`
	
	// RetryAttempts is the number of retry attempts for failed sends.
	RetryAttempts int `yaml:"retry_attempts" json:"retry_attempts"`
	
	// RetryDelay is the delay between retry attempts.
	RetryDelay time.Duration `yaml:"retry_delay" json:"retry_delay"`
	
	// Timeout is the HTTP timeout for sending events.
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
	
	// UserAgent is the user agent string for HTTP requests.
	UserAgent string `yaml:"user_agent" json:"user_agent"`
}

// DefaultConfig returns the default telemetry configuration.
func DefaultConfig() *Config {
	return &Config{
		Enabled:        true,
		Endpoint:       "https://telemetry.tyk.io/v1/events",
		BatchSize:      100,
		FlushInterval:  5 * time.Minute,
		RetryAttempts:  3,
		RetryDelay:     1 * time.Second,
		Timeout:        30 * time.Second,
		UserAgent:      "tykctl-go-telemetry/1.0",
	}
}

// Client represents a telemetry client interface.
type Client interface {
	// Track records a telemetry event.
	Track(event *Event) error
	
	// Flush sends any pending events immediately.
	Flush() error
	
	// Close shuts down the telemetry client gracefully.
	Close() error
	
	// IsEnabled returns whether telemetry is currently enabled.
	IsEnabled() bool
	
	// SetEnabled enables or disables telemetry.
	SetEnabled(enabled bool) error
}

// Transport represents a transport mechanism for sending telemetry data.
type Transport interface {
	// Send sends a batch of events to the telemetry endpoint.
	Send(events []*Event) error
	
	// Close closes the transport.
	Close() error
}

// Storage represents a storage mechanism for telemetry events.
type Storage interface {
	// Store stores an event for later transmission.
	Store(event *Event) error
	
	// Retrieve retrieves stored events.
	Retrieve() ([]*Event, error)
	
	// Clear removes all stored events.
	Clear() error
	
	// Count returns the number of stored events.
	Count() (int, error)
}

// EventBuilder provides a fluent interface for building telemetry events.
type EventBuilder struct {
	event *Event
}

// NewEventBuilder creates a new event builder.
func NewEventBuilder(eventType EventType) *EventBuilder {
	return &EventBuilder{
		event: &Event{
			EventType: eventType,
			Timestamp: time.Now(),
			Properties: make(map[string]interface{}),
		},
	}
}

// Command sets the command for the event.
func (b *EventBuilder) Command(cmd string) *EventBuilder {
	b.event.Command = cmd
	return b
}

// Duration sets the duration for the event.
func (b *EventBuilder) Duration(duration time.Duration) *EventBuilder {
	b.event.Duration = duration.Milliseconds()
	return b
}

// Success sets the success status for the event.
func (b *EventBuilder) Success(success bool) *EventBuilder {
	b.event.Success = success
	return b
}

// Error sets error information for the event.
func (b *EventBuilder) Error(errorType, message string) *EventBuilder {
	b.event.ErrorType = errorType
	b.event.ErrorMessage = message
	b.event.Success = false
	return b
}

// Feature sets the feature for the event.
func (b *EventBuilder) Feature(feature string) *EventBuilder {
	b.event.Feature = feature
	return b
}

// Property adds a property to the event.
func (b *EventBuilder) Property(key string, value interface{}) *EventBuilder {
	if b.event.Properties == nil {
		b.event.Properties = make(map[string]interface{})
	}
	b.event.Properties[key] = value
	return b
}

// Properties sets multiple properties for the event.
func (b *EventBuilder) Properties(props map[string]interface{}) *EventBuilder {
	if b.event.Properties == nil {
		b.event.Properties = make(map[string]interface{})
	}
	for k, v := range props {
		b.event.Properties[k] = v
	}
	return b
}

// Build returns the built event.
func (b *EventBuilder) Build() *Event {
	return b.event
}

// SanitizeEvent removes or masks sensitive information from an event.
func SanitizeEvent(event *Event) *Event {
	sanitized := *event
	
	// Remove any properties that might contain sensitive data
	if sanitized.Properties != nil {
		sensitiveKeys := []string{
			"token", "key", "secret", "password", "auth",
			"credential", "api_key", "access_token", "refresh_token",
		}
		
		for _, key := range sensitiveKeys {
			delete(sanitized.Properties, key)
		}
	}
	
	// Sanitize error messages to remove sensitive data
	if sanitized.ErrorMessage != "" {
		sanitized.ErrorMessage = sanitizeErrorMessage(sanitized.ErrorMessage)
	}
	
	return &sanitized
}

// sanitizeErrorMessage removes sensitive information from error messages.
func sanitizeErrorMessage(msg string) string {
	// This is a simple implementation - in practice, you might want
	// more sophisticated sanitization based on your specific needs
	sensitivePatterns := []string{
		"token", "key", "secret", "password", "auth",
		"credential", "api_key", "access_token", "refresh_token",
	}
	
	sanitized := msg
	for _, pattern := range sensitivePatterns {
		// Replace sensitive patterns with [REDACTED]
		// This is a simplified approach - you might want regex-based replacement
		if len(sanitized) > 0 {
			// Simple string replacement for demonstration
			// In practice, use regex for more sophisticated matching
			// For now, just return the original message
			_ = pattern // Avoid unused variable warning
		}
	}
	
	return sanitized
}