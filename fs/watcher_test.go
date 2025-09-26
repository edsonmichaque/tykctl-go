package fs

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func TestWatcher_New(t *testing.T) {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer watcher.Stop()

	if watcher == nil {
		t.Fatal("Watcher should not be nil")
	}
}

func TestWatcher_NewWithLogger(t *testing.T) {
	logger := zap.NewNop()
	watcher, err := NewWatcherWithLogger(logger)
	if err != nil {
		t.Fatalf("NewWatcherWithLogger failed: %v", err)
	}
	defer watcher.Stop()

	if watcher.logger != logger {
		t.Fatal("Logger should be set correctly")
	}
}

func TestWatcher_Watch(t *testing.T) {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer watcher.Stop()

	// Create a temporary directory
	dir := filepath.Join(os.TempDir(), "tykctl-watcher-test")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Watch the directory
	err = watcher.Watch(dir)
	if err != nil {
		t.Fatalf("Watch failed: %v", err)
	}
}

func TestWatcher_AddHandler(t *testing.T) {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer watcher.Stop()

	// Add a handler
	handler := WatchHandlerFunc(func(ctx context.Context, event WatchEvent) error {
		return nil
	})

	watcher.AddHandler("*.txt", handler)
	watcher.AddHandlerFunc("*.log", func(ctx context.Context, event WatchEvent) error {
		return nil
	})

	// Check that handlers were added
	if len(watcher.handlers) != 2 {
		t.Fatalf("Expected 2 handlers, got %d", len(watcher.handlers))
	}
}

func TestWatcher_StartStop(t *testing.T) {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}

	// Start the watcher
	watcher.Start()

	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)

	// Stop the watcher
	watcher.Stop()
}

func TestWatcher_ContextCancellation(t *testing.T) {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}

	// Start the watcher
	watcher.Start()

	// Cancel the context
	watcher.cancel()

	// Give it a moment to stop
	time.Sleep(10 * time.Millisecond)

	// Stop should be safe to call multiple times
	watcher.Stop()
}

func TestWatcher_GlobalWatcher(t *testing.T) {
	// Test that global watcher can be created
	watcher1 := GetGlobalWatcher()
	watcher2 := GetGlobalWatcher()

	// Should be the same instance
	if watcher1 != watcher2 {
		t.Fatal("Global watcher should be singleton")
	}

	// Start and stop global watcher
	StartGlobalWatcher()
	time.Sleep(10 * time.Millisecond)
	StopGlobalWatcher()
}

func TestWatcher_WatchConfigFile(t *testing.T) {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer watcher.Stop()

	// Create a temporary config file
	configDir := filepath.Join(os.TempDir(), "tykctl-config-test")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	defer os.RemoveAll(configDir)

	configFile := filepath.Join(configDir, "config.yaml")

	// Track reload calls
	var reloadCount int
	var mu sync.Mutex
	reloadFunc := func() error {
		mu.Lock()
		reloadCount++
		mu.Unlock()
		return nil
	}

	// Watch the config file
	err = watcher.WatchConfigFile(configFile, reloadFunc)
	if err != nil {
		t.Fatalf("WatchConfigFile failed: %v", err)
	}

	// Start the watcher
	watcher.Start()
	defer watcher.Stop()

	// Create the config file
	err = os.WriteFile(configFile, []byte("test: config"), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Wait for the event to be processed
	time.Sleep(100 * time.Millisecond)

	// Check if reload was called
	mu.Lock()
	count := reloadCount
	mu.Unlock()

	if count == 0 {
		t.Fatal("Reload function should have been called")
	}
}

func TestWatcher_WatchExtensions(t *testing.T) {
	watcher, err := NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher failed: %v", err)
	}
	defer watcher.Stop()

	// Create a temporary extensions directory
	extDir := filepath.Join(os.TempDir(), "tykctl-extensions-test")
	err = os.MkdirAll(extDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create extensions directory: %v", err)
	}
	defer os.RemoveAll(extDir)

	// Track change calls
	var changeCount int
	var mu sync.Mutex
	changeFunc := func() error {
		mu.Lock()
		changeCount++
		mu.Unlock()
		return nil
	}

	// Watch the extensions directory
	err = watcher.WatchExtensions(extDir, changeFunc)
	if err != nil {
		t.Fatalf("WatchExtensions failed: %v", err)
	}

	// Start the watcher
	watcher.Start()
	defer watcher.Stop()

	// Create an extension file
	extFile := filepath.Join(extDir, "tykctl-test")
	err = os.WriteFile(extFile, []byte("#!/bin/bash\necho test"), 0755)
	if err != nil {
		t.Fatalf("Failed to create extension file: %v", err)
	}

	// Wait for the event to be processed
	time.Sleep(100 * time.Millisecond)

	// Check if change was called
	mu.Lock()
	count := changeCount
	mu.Unlock()

	if count == 0 {
		t.Fatal("Change function should have been called")
	}
}

func TestWatchEvent(t *testing.T) {
	event := WatchEvent{
		Name:      "/test/file.txt",
		Op:        fsnotify.Write,
		Timestamp: time.Now(),
	}

	if event.Name != "/test/file.txt" {
		t.Fatal("Event name should be set correctly")
	}

	if event.Op != fsnotify.Write {
		t.Fatal("Event operation should be set correctly")
	}

	if event.Timestamp.IsZero() {
		t.Fatal("Event timestamp should be set")
	}
}

func TestWatchHandlerFunc(t *testing.T) {
	var called bool
	handler := WatchHandlerFunc(func(ctx context.Context, event WatchEvent) error {
		called = true
		return nil
	})

	ctx := context.Background()
	event := WatchEvent{Name: "test", Op: fsnotify.Write, Timestamp: time.Now()}

	err := handler.HandleEvent(ctx, event)
	if err != nil {
		t.Fatalf("HandleEvent failed: %v", err)
	}

	if !called {
		t.Fatal("Handler function should have been called")
	}
}
