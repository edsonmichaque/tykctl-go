package script

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewScriptManager(t *testing.T) {
	scriptDir := "/tmp/test-scripts"
	manager := NewScriptManager(scriptDir)
	
	if manager == nil {
		t.Fatal("NewScriptManager() returned nil")
	}
	
	if manager.scriptDir != scriptDir {
		t.Errorf("Expected script dir '%s', got '%s'", scriptDir, manager.scriptDir)
	}
}

func TestNewScriptManagerWithLogger(t *testing.T) {
	scriptDir := "/tmp/test-scripts"
	logger, _ := zap.NewDevelopment()
	manager := NewScriptManagerWithLogger(scriptDir, logger)
	
	if manager == nil {
		t.Fatal("NewScriptManagerWithLogger() returned nil")
	}
	
	if manager.scriptDir != scriptDir {
		t.Errorf("Expected script dir '%s', got '%s'", scriptDir, manager.scriptDir)
	}
	
	if manager.logger != logger {
		t.Error("Logger was not set correctly")
	}
}

func TestNewScriptRegistry(t *testing.T) {
	registry := NewScriptRegistry()
	
	if registry == nil {
		t.Fatal("NewScriptRegistry() returned nil")
	}
	
	if registry.handlers == nil {
		t.Error("Handlers map should be initialized")
	}
	
	if len(registry.handlers) != 0 {
		t.Error("Handlers map should be empty initially")
	}
}

func TestRegisterHandler(t *testing.T) {
	registry := NewScriptRegistry()
	event := ScriptEvent("test-event")
	
	handler := func(ctx context.Context, scriptCtx *ScriptContext) error {
		return nil
	}
	
	registry.RegisterHandler(event, handler)
	
	if len(registry.handlers) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(registry.handlers))
	}
	
	if registry.handlers[event] == nil {
		t.Error("Handler was not registered")
	}
}

func TestExecuteHandlers(t *testing.T) {
	registry := NewScriptRegistry()
	ctx := context.Background()
	event := ScriptEvent("test-event")
	
	var executed bool
	handler := func(ctx context.Context, scriptCtx *ScriptContext) error {
		executed = true
		return nil
	}
	
	registry.RegisterHandler(event, handler)
	
	scriptCtx := &ScriptContext{
		Event: event,
		Data:  map[string]interface{}{"test": "data"},
	}
	
	err := registry.ExecuteHandlers(ctx, event, scriptCtx)
	if err != nil {
		t.Errorf("ExecuteHandlers failed: %v", err)
	}
	
	if !executed {
		t.Error("Handler was not executed")
	}
}

func TestExecuteHandlersMultiple(t *testing.T) {
	registry := NewScriptRegistry()
	ctx := context.Background()
	event := ScriptEvent("test-event")
	
	var executionCount int
	handler1 := func(ctx context.Context, scriptCtx *ScriptContext) error {
		executionCount++
		return nil
	}
	
	handler2 := func(ctx context.Context, scriptCtx *ScriptContext) error {
		executionCount++
		return nil
	}
	
	registry.RegisterHandler(event, handler1)
	registry.RegisterHandler(event, handler2)
	
	scriptCtx := &ScriptContext{
		Event: event,
		Data:  map[string]interface{}{"test": "data"},
	}
	
	err := registry.ExecuteHandlers(ctx, event, scriptCtx)
	if err != nil {
		t.Errorf("ExecuteHandlers failed: %v", err)
	}
	
	if executionCount != 2 {
		t.Errorf("Expected 2 executions, got %d", executionCount)
	}
}

func TestExecuteHandlersError(t *testing.T) {
	registry := NewScriptRegistry()
	ctx := context.Background()
	event := ScriptEvent("test-event")
	
	expectedError := fmt.Errorf("handler error")
	handler := func(ctx context.Context, scriptCtx *ScriptContext) error {
		return expectedError
	}
	
	registry.RegisterHandler(event, handler)
	
	scriptCtx := &ScriptContext{
		Event: event,
		Data:  map[string]interface{}{"test": "data"},
	}
	
	err := registry.ExecuteHandlers(ctx, event, scriptCtx)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	// The error is wrapped, so we check if it contains the expected error
	if err == nil || !strings.Contains(err.Error(), expectedError.Error()) {
		t.Errorf("Expected error containing '%v', got '%v'", expectedError, err)
	}
}

func TestExecuteHandlersNoHandlers(t *testing.T) {
	registry := NewScriptRegistry()
	ctx := context.Background()
	event := ScriptEvent("test-event")
	
	scriptCtx := &ScriptContext{
		Event: event,
		Data:  map[string]interface{}{"test": "data"},
	}
	
	err := registry.ExecuteHandlers(ctx, event, scriptCtx)
	if err != nil {
		t.Errorf("ExecuteHandlers should not fail with no handlers: %v", err)
	}
}

func TestScriptContext(t *testing.T) {
	scriptCtx := &ScriptContext{
		Event:       ScriptEvent("test-event"),
		Command:     "test-command",
		Args:        []string{"arg1", "arg2"},
		Extension:   "test-extension",
		WorkingDir:  "/tmp",
		Environment: map[string]string{"KEY": "value"},
		Data:        map[string]interface{}{"test": "data"},
	}
	
	if scriptCtx.Event != ScriptEvent("test-event") {
		t.Errorf("Expected Event 'test-event', got '%s'", scriptCtx.Event)
	}
	
	if scriptCtx.Command != "test-command" {
		t.Errorf("Expected Command 'test-command', got '%s'", scriptCtx.Command)
	}
	
	if len(scriptCtx.Args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(scriptCtx.Args))
	}
	
	if scriptCtx.Extension != "test-extension" {
		t.Errorf("Expected Extension 'test-extension', got '%s'", scriptCtx.Extension)
	}
	
	if scriptCtx.WorkingDir != "/tmp" {
		t.Errorf("Expected WorkingDir '/tmp', got '%s'", scriptCtx.WorkingDir)
	}
	
	if scriptCtx.Environment["KEY"] != "value" {
		t.Error("Environment KEY not set correctly")
	}
	
	if scriptCtx.Data["test"] != "data" {
		t.Error("Data test not set correctly")
	}
}

func TestScriptEvent(t *testing.T) {
	event := ScriptEvent("test-event")
	
	if string(event) != "test-event" {
		t.Errorf("Expected event 'test-event', got '%s'", string(event))
	}
}

func TestScriptStruct(t *testing.T) {
	script := &Script{
		Name:        "test-script",
		Description: "Test script",
		Script:      "echo 'Hello World'",
		Enabled:     true,
		Timeout:     30 * time.Second,
		Environment: map[string]string{"KEY": "value"},
	}
	
	if script.Name != "test-script" {
		t.Errorf("Expected Name 'test-script', got '%s'", script.Name)
	}
	
	if script.Description != "Test script" {
		t.Errorf("Expected Description 'Test script', got '%s'", script.Description)
	}
	
	if script.Script != "echo 'Hello World'" {
		t.Errorf("Expected Script 'echo 'Hello World'', got '%s'", script.Script)
	}
	
	if !script.Enabled {
		t.Error("Script should be enabled")
	}
	
	if script.Timeout != 30*time.Second {
		t.Errorf("Expected Timeout 30s, got %v", script.Timeout)
	}
	
	if script.Environment["KEY"] != "value" {
		t.Error("Environment KEY not set correctly")
	}
}

func TestContextCancellation(t *testing.T) {
	registry := NewScriptRegistry()
	ctx := context.Background()
	
	event := ScriptEvent("test-event")
	handler := func(ctx context.Context, scriptCtx *ScriptContext) error {
		// Simple handler that respects context
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
	
	registry.RegisterHandler(event, handler)
	
	scriptCtx := &ScriptContext{
		Event: event,
		Data:  map[string]interface{}{"test": "data"},
	}
	
	err := registry.ExecuteHandlers(ctx, event, scriptCtx)
	if err != nil {
		t.Errorf("ExecuteHandlers failed: %v", err)
	}
}

func TestConcurrentHandlers(t *testing.T) {
	registry := NewScriptRegistry()
	ctx := context.Background()
	event := ScriptEvent("test-event")
	
	var executionCount int
	handler := func(ctx context.Context, scriptCtx *ScriptContext) error {
		executionCount++
		return nil
	}
	
	registry.RegisterHandler(event, handler)
	
	scriptCtx := &ScriptContext{
		Event: event,
		Data:  map[string]interface{}{"test": "data"},
	}
	
	// Execute handlers concurrently
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func() {
			err := registry.ExecuteHandlers(ctx, event, scriptCtx)
			done <- err
		}()
	}
	
	// Wait for all to complete
	for i := 0; i < 10; i++ {
		err := <-done
		if err != nil {
			t.Errorf("Concurrent execution failed: %v", err)
		}
	}
	
	if executionCount != 10 {
		t.Errorf("Expected 10 executions, got %d", executionCount)
	}
}

// Benchmark tests
func BenchmarkNewScriptRegistry(b *testing.B) {
	for i := 0; i < b.N; i++ {
		registry := NewScriptRegistry()
		_ = registry
	}
}

func BenchmarkRegisterHandler(b *testing.B) {
	registry := NewScriptRegistry()
	event := ScriptEvent("test-event")
	handler := func(ctx context.Context, scriptCtx *ScriptContext) error {
		return nil
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.RegisterHandler(event, handler)
	}
}

func BenchmarkExecuteHandlers(b *testing.B) {
	registry := NewScriptRegistry()
	ctx := context.Background()
	event := ScriptEvent("test-event")
	handler := func(ctx context.Context, scriptCtx *ScriptContext) error {
		return nil
	}
	
	registry.RegisterHandler(event, handler)
	scriptCtx := &ScriptContext{
		Event: event,
		Data:  map[string]interface{}{"test": "data"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.ExecuteHandlers(ctx, event, scriptCtx)
	}
}

func BenchmarkExecuteHandlersConcurrent(b *testing.B) {
	registry := NewScriptRegistry()
	ctx := context.Background()
	event := ScriptEvent("test-event")
	handler := func(ctx context.Context, scriptCtx *ScriptContext) error {
		return nil
	}
	
	registry.RegisterHandler(event, handler)
	scriptCtx := &ScriptContext{
		Event: event,
		Data:  map[string]interface{}{"test": "data"},
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			registry.ExecuteHandlers(ctx, event, scriptCtx)
		}
	})
}