package telemetry

import (
	"testing"
	"time"
)

func TestBasicFunctionality(t *testing.T) {
	// Test creating a no-op client
	client := NewNoOpClient()
	if client.IsEnabled() {
		t.Error("NoOp client should report as disabled")
	}
	
	// Test creating an event
	event := NewEventBuilder(EventTypeCommand).
		Command("test").
		Success(true).
		Build()
	
	if event.Command != "test" {
		t.Error("Event command should be 'test'")
	}
	
	if !event.Success {
		t.Error("Event should be successful")
	}
	
	// Test tracking with no-op client (should not fail)
	if err := client.Track(event); err != nil {
		t.Errorf("Tracking with no-op client should not fail: %v", err)
	}
	
	// Test closing no-op client
	if err := client.Close(); err != nil {
		t.Errorf("Closing no-op client should not fail: %v", err)
	}
}

func TestMockTransport(t *testing.T) {
	// Create mock transport
	mockTransport := NewMockTransport()
	
	// Create mock storage
	mockStorage := NewMockStorage()
	
	// Create config
	config := DefaultConfig()
	
	// Create client
	client := NewClient(config, mockTransport, mockStorage)
	defer client.Close()
	
	// Track an event
	event := NewEventBuilder(EventTypeCommand).
		Command("test").
		Success(true).
		Build()
	
	if err := client.Track(event); err != nil {
		t.Errorf("Failed to track event: %v", err)
	}
	
	// Flush events
	if err := client.Flush(); err != nil {
		t.Errorf("Failed to flush events: %v", err)
	}
	
	// Check what was sent
	sentEvents := mockTransport.GetEvents()
	if len(sentEvents) != 1 {
		t.Errorf("Expected 1 batch, got %d", len(sentEvents))
	}
	
	if len(sentEvents[0]) != 1 {
		t.Errorf("Expected 1 event in batch, got %d", len(sentEvents[0]))
	}
	
	if sentEvents[0][0].Command != "test" {
		t.Errorf("Expected command 'test', got '%s'", sentEvents[0][0].Command)
	}
}

func TestEventBuilder(t *testing.T) {
	// Test building a command event
	event := NewEventBuilder(EventTypeCommand).
		Command("test-command").
		Duration(100 * time.Millisecond).
		Success(true).
		Property("key1", "value1").
		Property("key2", 42).
		Build()
	
	if event.EventType != EventTypeCommand {
		t.Errorf("Expected EventTypeCommand, got %s", event.EventType)
	}
	
	if event.Command != "test-command" {
		t.Errorf("Expected command 'test-command', got '%s'", event.Command)
	}
	
	if event.Duration != 100 {
		t.Errorf("Expected duration 100ms, got %d", event.Duration)
	}
	
	if !event.Success {
		t.Error("Expected success to be true")
	}
	
	if event.Properties["key1"] != "value1" {
		t.Errorf("Expected property key1='value1', got '%v'", event.Properties["key1"])
	}
	
	if event.Properties["key2"] != 42 {
		t.Errorf("Expected property key2=42, got '%v'", event.Properties["key2"])
	}
}

func TestErrorEvent(t *testing.T) {
	// Test building an error event
	event := NewEventBuilder(EventTypeError).
		Error("test_error", "Something went wrong").
		Build()
	
	if event.EventType != EventTypeError {
		t.Errorf("Expected EventTypeError, got %s", event.EventType)
	}
	
	if event.ErrorType != "test_error" {
		t.Errorf("Expected error type 'test_error', got '%s'", event.ErrorType)
	}
	
	if event.ErrorMessage != "Something went wrong" {
		t.Errorf("Expected error message 'Something went wrong', got '%s'", event.ErrorMessage)
	}
	
	if event.Success {
		t.Error("Error event should not be successful")
	}
}

func TestFeatureEvent(t *testing.T) {
	// Test building a feature event
	event := NewEventBuilder(EventTypeFeature).
		Feature("auto_scaling").
		Success(true).
		Properties(map[string]interface{}{
			"min_replicas": 2,
			"max_replicas": 10,
		}).
		Build()
	
	if event.EventType != EventTypeFeature {
		t.Errorf("Expected EventTypeFeature, got %s", event.EventType)
	}
	
	if event.Feature != "auto_scaling" {
		t.Errorf("Expected feature 'auto_scaling', got '%s'", event.Feature)
	}
	
	if !event.Success {
		t.Error("Feature event should be successful")
	}
	
	if event.Properties["min_replicas"] != 2 {
		t.Errorf("Expected min_replicas=2, got %v", event.Properties["min_replicas"])
	}
	
	if event.Properties["max_replicas"] != 10 {
		t.Errorf("Expected max_replicas=10, got %v", event.Properties["max_replicas"])
	}
}

func TestSanitizeEvent(t *testing.T) {
	// Test sanitizing an event with sensitive data
	event := &Event{
		EventType: EventTypeCommand,
		Command:   "test",
		Success:   true,
		Properties: map[string]interface{}{
			"api_key":    "secret123",
			"password":   "mypass",
			"safe_data":  "this is safe",
		},
		ErrorMessage: "Error with token abc123",
	}
	
	sanitized := SanitizeEvent(event)
	
	// Sensitive properties should be removed
	if sanitized.Properties["api_key"] != nil {
		t.Error("api_key should be removed from properties")
	}
	
	if sanitized.Properties["password"] != nil {
		t.Error("password should be removed from properties")
	}
	
	// Safe data should remain
	if sanitized.Properties["safe_data"] != "this is safe" {
		t.Error("safe_data should remain in properties")
	}
	
	// Error message should be sanitized (though our current implementation is basic)
	if sanitized.ErrorMessage == "" {
		t.Error("Error message should not be empty after sanitization")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if !config.Enabled {
		t.Error("Default config should have telemetry enabled")
	}
	
	if config.Endpoint == "" {
		t.Error("Default config should have an endpoint")
	}
	
	if config.BatchSize <= 0 {
		t.Error("Default config should have a positive batch size")
	}
	
	if config.FlushInterval <= 0 {
		t.Error("Default config should have a positive flush interval")
	}
	
	if config.RetryAttempts < 0 {
		t.Error("Default config should have non-negative retry attempts")
	}
	
	if config.RetryDelay <= 0 {
		t.Error("Default config should have a positive retry delay")
	}
	
	if config.Timeout <= 0 {
		t.Error("Default config should have a positive timeout")
	}
	
	if config.UserAgent == "" {
		t.Error("Default config should have a user agent")
	}
}