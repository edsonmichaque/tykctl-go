package logger

import (
	"testing"
)

func TestNew(t *testing.T) {
	config := Config{
		Debug:   false,
		Verbose: false,
		NoColor: false,
	}

	logger := New(config)
	if logger == nil {
		t.Fatal("New() returned nil")
	}
	if logger.Logger == nil {
		t.Fatal("Logger is nil")
	}
}

func TestNewWithDebug(t *testing.T) {
	config := Config{
		Debug:   true,
		Verbose: false,
		NoColor: false,
	}

	logger := New(config)
	if logger == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNewWithVerbose(t *testing.T) {
	config := Config{
		Debug:   false,
		Verbose: true,
		NoColor: false,
	}

	logger := New(config)
	if logger == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNewWithNoColor(t *testing.T) {
	config := Config{
		Debug:   false,
		Verbose: false,
		NoColor: true,
	}

	logger := New(config)
	if logger == nil {
		t.Fatal("New() returned nil")
	}
}

func TestGlobalLogger(t *testing.T) {
	// Test that global logger can be initialized
	config := Config{
		Debug:   false,
		Verbose: false,
		NoColor: false,
	}

	InitGlobal(config)
	global := GetGlobal()
	if global == nil {
		t.Fatal("GetGlobal() returned nil")
	}

	// Test sync
	SyncGlobal()
}

func TestSync(t *testing.T) {
	config := Config{
		Debug:   false,
		Verbose: false,
		NoColor: false,
	}

	logger := New(config)
	if logger == nil {
		t.Fatal("New() returned nil")
	}

	// Test sync method
	logger.Sync()
}
