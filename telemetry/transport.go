// Package telemetry provides anonymous usage analytics for tykctl-go.
package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// HTTPTransport implements the Transport interface using HTTP.
type HTTPTransport struct {
	client    *http.Client
	endpoint  string
	userAgent string
}

// NewHTTPTransport creates a new HTTP transport.
func NewHTTPTransport(endpoint, userAgent string, timeout time.Duration) Transport {
	return &HTTPTransport{
		client: &http.Client{
			Timeout: timeout,
		},
		endpoint:  endpoint,
		userAgent: userAgent,
	}
}

// Send sends a batch of events to the telemetry endpoint.
func (t *HTTPTransport) Send(events []*Event) error {
	if len(events) == 0 {
		return nil
	}
	
	// Marshal events to JSON
	payload, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		t.endpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept", "application/json")
	
	// Send request
	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telemetry endpoint returned status %d", resp.StatusCode)
	}
	
	return nil
}

// Close closes the transport.
func (t *HTTPTransport) Close() error {
	// HTTP client doesn't need explicit closing
	return nil
}

// MockTransport is a mock transport for testing.
type MockTransport struct {
	events [][]*Event
	errors []error
}

// NewMockTransport creates a new mock transport.
func NewMockTransport() *MockTransport {
	return &MockTransport{
		events: make([][]*Event, 0),
		errors: make([]error, 0),
	}
}

// Send records the events for later inspection.
func (t *MockTransport) Send(events []*Event) error {
	// Create a copy of the events to avoid reference issues
	eventCopy := make([]*Event, len(events))
	for i, event := range events {
		eventCopy[i] = event
	}
	
	t.events = append(t.events, eventCopy)
	
	// Return any configured error
	if len(t.errors) > 0 {
		err := t.errors[0]
		t.errors = t.errors[1:]
		return err
	}
	
	return nil
}

// Close does nothing for the mock transport.
func (t *MockTransport) Close() error {
	return nil
}

// GetEvents returns all events that were sent.
func (t *MockTransport) GetEvents() [][]*Event {
	return t.events
}

// SetError configures the next Send call to return an error.
func (t *MockTransport) SetError(err error) {
	t.errors = append(t.errors, err)
}

// Clear clears all recorded events and errors.
func (t *MockTransport) Clear() {
	t.events = make([][]*Event, 0)
	t.errors = make([]error, 0)
}

// FileTransport implements the Transport interface by writing events to a file.
type FileTransport struct {
	filename string
}

// NewFileTransport creates a new file transport.
func NewFileTransport(filename string) Transport {
	return &FileTransport{
		filename: filename,
	}
}

// Send writes events to a file in JSON format.
func (t *FileTransport) Send(events []*Event) error {
	if len(events) == 0 {
		return nil
	}
	
	// Marshal events to JSON
	payload, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}
	
	// Append to file
	file, err := os.OpenFile(t.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	// Write timestamp and events
	timestamp := time.Now().Format(time.RFC3339)
	_, err = fmt.Fprintf(file, "=== %s ===\n%s\n\n", timestamp, string(payload))
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	
	return nil
}

// Close closes the file transport.
func (t *FileTransport) Close() error {
	return nil
}