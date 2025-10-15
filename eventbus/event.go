// Package eventbus provides a comprehensive event-driven architecture for tykctl-go.
// It replaces the hook system with a more powerful, scalable event bus.
package eventbus

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// EventType represents the type of event in the system.
// Implementers should define their own event type constants.
type EventType string

// Event represents a single event in the system.
type Event struct {
	// ID is a unique identifier for the event.
	ID string `json:"id"`

	// Type is the type of event.
	Type EventType `json:"type"`

	// Data contains the event payload.
	Data interface{} `json:"data"`

	// Metadata contains additional event metadata.
	Metadata map[string]interface{} `json:"metadata"`

	// Timestamp is when the event occurred.
	Timestamp time.Time `json:"timestamp"`

	// Source identifies where the event originated.
	Source string `json:"source"`

	// Version is the event schema version.
	Version string `json:"version"`

	// CorrelationID links related events.
	CorrelationID string `json:"correlation_id,omitempty"`

	// ParentID links to a parent event.
	ParentID string `json:"parent_id,omitempty"`
}

// NewEvent creates a new event with the given type and data.
func NewEvent(eventType EventType, data interface{}) *Event {
	return &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Data:      data,
		Metadata:  make(map[string]interface{}),
		Timestamp: time.Now(),
		Version:   "1.0",
	}
}

// generateEventID generates a simple event ID without external dependencies.
func generateEventID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// WithMetadata adds metadata to the event.
func (e *Event) WithMetadata(key string, value interface{}) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// WithMetadataMap adds multiple metadata entries to the event.
func (e *Event) WithMetadataMap(metadata map[string]interface{}) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	for k, v := range metadata {
		e.Metadata[k] = v
	}
	return e
}

// WithSource sets the event source.
func (e *Event) WithSource(source string) *Event {
	e.Source = source
	return e
}

// WithCorrelationID sets the correlation ID for linking related events.
func (e *Event) WithCorrelationID(correlationID string) *Event {
	e.CorrelationID = correlationID
	return e
}

// WithParentID sets the parent event ID.
func (e *Event) WithParentID(parentID string) *Event {
	e.ParentID = parentID
	return e
}

// Clone creates a deep copy of the event.
func (e *Event) Clone() *Event {
	clone := &Event{
		ID:            e.ID,
		Type:          e.Type,
		Data:          e.Data, // Shallow copy - caller responsible for deep copy if needed
		Timestamp:     e.Timestamp,
		Source:        e.Source,
		Version:       e.Version,
		CorrelationID: e.CorrelationID,
		ParentID:      e.ParentID,
	}

	if e.Metadata != nil {
		clone.Metadata = make(map[string]interface{})
		for k, v := range e.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}

// String returns a string representation of the event.
func (e *Event) String() string {
	return fmt.Sprintf("Event{ID: %s, Type: %s, Source: %s, Timestamp: %s}", 
		e.ID, e.Type, e.Source, e.Timestamp.Format(time.RFC3339))
}

// MarshalJSON customizes JSON marshaling for the event.
func (e *Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	return json.Marshal(&struct {
		*Alias
		Timestamp string `json:"timestamp"`
	}{
		Alias:     (*Alias)(e),
		Timestamp: e.Timestamp.Format(time.RFC3339Nano),
	})
}

// UnmarshalJSON customizes JSON unmarshaling for the event.
func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		*Alias
		Timestamp string `json:"timestamp"`
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	e.Timestamp, err = time.Parse(time.RFC3339Nano, aux.Timestamp)
	if err != nil {
		// Fallback to RFC3339
		e.Timestamp, err = time.Parse(time.RFC3339, aux.Timestamp)
		if err != nil {
			return fmt.Errorf("invalid timestamp format: %w", err)
		}
	}

	return nil
}

// EventFilter represents a filter for querying events.
type EventFilter struct {
	Types       []EventType           `json:"types,omitempty"`
	Sources     []string              `json:"sources,omitempty"`
	StartTime   *time.Time            `json:"start_time,omitempty"`
	EndTime     *time.Time            `json:"end_time,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CorrelationID string              `json:"correlation_id,omitempty"`
	ParentID    string                `json:"parent_id,omitempty"`
}

// Matches checks if an event matches the filter criteria.
func (f *EventFilter) Matches(event *Event) bool {
	// Check event types
	if len(f.Types) > 0 {
		found := false
		for _, eventType := range f.Types {
			if event.Type == eventType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check sources
	if len(f.Sources) > 0 {
		found := false
		for _, source := range f.Sources {
			if event.Source == source {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check time range
	if f.StartTime != nil && event.Timestamp.Before(*f.StartTime) {
		return false
	}
	if f.EndTime != nil && event.Timestamp.After(*f.EndTime) {
		return false
	}

	// Check correlation ID
	if f.CorrelationID != "" && event.CorrelationID != f.CorrelationID {
		return false
	}

	// Check parent ID
	if f.ParentID != "" && event.ParentID != f.ParentID {
		return false
	}

	// Check metadata
	if len(f.Metadata) > 0 {
		for key, value := range f.Metadata {
			if event.Metadata == nil {
				return false
			}
			if event.Metadata[key] != value {
				return false
			}
		}
	}

	return true
}

// EventContext provides context for event processing.
type EventContext struct {
	// Event is the event being processed.
	Event *Event

	// Context is the Go context for cancellation and timeouts.
	Context interface{} // Using interface{} to avoid import cycle

	// Metadata contains processing metadata.
	Metadata map[string]interface{}

	// CancelFunc can be used to cancel event processing.
	CancelFunc func()
}

// NewEventContext creates a new event context.
func NewEventContext(event *Event, ctx interface{}) *EventContext {
	return &EventContext{
		Event:    event,
		Context:  ctx,
		Metadata: make(map[string]interface{}),
	}
}

// WithMetadata adds metadata to the event context.
func (ec *EventContext) WithMetadata(key string, value interface{}) *EventContext {
	if ec.Metadata == nil {
		ec.Metadata = make(map[string]interface{})
	}
	ec.Metadata[key] = value
	return ec
}