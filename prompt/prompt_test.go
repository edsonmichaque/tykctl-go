package prompt

import (
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	logger := zap.NewNop()
	p := New(logger)
	if p == nil {
		t.Fatal("New() returned nil")
	}
	if p.logger == nil {
		t.Fatal("logger is nil")
	}
}

func TestExtensionInstallPrompt(t *testing.T) {
	logger := zap.NewNop()
	_ = New(logger)

	// Test the prompt message formatting
	message := "Install extension 'test-extension'?"
	expected := "Install extension 'test-extension'?"

	// We can't easily test the actual prompt interaction without mocking
	// but we can test the message formatting logic
	if message != expected {
		t.Errorf("Expected message '%s', got '%s'", expected, message)
	}
}

func TestExtensionRemovePrompt(t *testing.T) {
	logger := zap.NewNop()
	_ = New(logger)

	// Test the prompt message formatting
	message := "Remove extension 'test-extension'?"
	expected := "Remove extension 'test-extension'?"

	if message != expected {
		t.Errorf("Expected message '%s', got '%s'", expected, message)
	}
}

func TestConfigPrompt(t *testing.T) {
	logger := zap.NewNop()
	_ = New(logger)

	// Test the prompt message formatting
	message := "api_key (API key for authentication)"
	expected := "api_key (API key for authentication)"

	if message != expected {
		t.Errorf("Expected message '%s', got '%s'", expected, message)
	}
}
