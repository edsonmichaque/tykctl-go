package hook

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	manager := New()
	if manager == nil {
		t.Fatal("New() returned nil")
	}
	
	if manager.builtin == nil {
		t.Error("Builtin manager should not be nil")
	}
	
	// External and Rego managers are initialized by default
	if manager.external == nil {
		t.Error("External manager should not be nil")
	}
	
	if manager.rego == nil {
		t.Error("Rego manager should not be nil")
	}
}

func TestNewWithOptions(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	hookDir := "/tmp/test-hooks"
	
	manager := New(
		WithExternalHookDir(hookDir),
		WithLogger(logger),
	)
	
	if manager == nil {
		t.Fatal("New() returned nil")
	}
	
	if manager.externalHookDir != hookDir {
		t.Errorf("Expected external hook dir '%s', got '%s'", hookDir, manager.externalHookDir)
	}
	
	if manager.logger != logger {
		t.Error("Logger was not set correctly")
	}
	
	if manager.external == nil {
		t.Error("External manager should not be nil")
	}
	
	if manager.rego == nil {
		t.Error("Rego manager should not be nil")
	}
}

func TestWithExternalHookDir(t *testing.T) {
	hookDir := "/tmp/test-hooks"
	option := WithExternalHookDir(hookDir)
	
	manager := &Manager{}
	option(manager)
	
	if manager.externalHookDir != hookDir {
		t.Errorf("Expected external hook dir '%s', got '%s'", hookDir, manager.externalHookDir)
	}
	
	if manager.external == nil {
		t.Error("External manager should be initialized")
	}
}

func TestWithLogger(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	option := WithLogger(logger)
	
	manager := &Manager{}
	option(manager)
	
	if manager.logger != logger {
		t.Error("Logger was not set correctly")
	}
	
	if manager.external == nil {
		t.Error("External manager should be initialized")
	}
	
	if manager.rego == nil {
		t.Error("Rego manager should be initialized")
	}
}

func TestRegisterBuiltin(t *testing.T) {
	manager := New()
	ctx := context.Background()
	
	hookType := HookType("test-hook")
	hookFunc := func(ctx context.Context, data *HookData) error {
		return nil
	}
	
	manager.RegisterBuiltin(ctx, hookType, hookFunc)
	
	// Verify hook was registered
	hooks := manager.ListBuiltin(ctx, hookType)
	if len(hooks) != 1 {
		t.Errorf("Expected 1 hook, got %d", len(hooks))
	}
}

func TestExecuteBuiltin(t *testing.T) {
	manager := New()
	ctx := context.Background()
	
	hookType := HookType("test-hook")
	var executed bool
	
	hookFunc := func(ctx context.Context, data *HookData) error {
		executed = true
		return nil
	}
	
	manager.RegisterBuiltin(ctx, hookType, hookFunc)
	
	hookData := &HookData{
		ExtensionName: "test-extension",
		ExtensionPath: "/path/to/extension",
		Metadata: map[string]interface{}{
			"test": "data",
		},
	}
	
	err := manager.Execute(ctx, hookType, hookData)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
	
	if !executed {
		t.Error("Hook function was not executed")
	}
}

func TestExecuteMultipleBuiltin(t *testing.T) {
	manager := New()
	ctx := context.Background()
	
	hookType := HookType("test-hook")
	var executionCount int
	
	hookFunc1 := func(ctx context.Context, data *HookData) error {
		executionCount++
		return nil
	}
	
	hookFunc2 := func(ctx context.Context, data *HookData) error {
		executionCount++
		return nil
	}
	
	manager.RegisterBuiltin(ctx, hookType, hookFunc1)
	manager.RegisterBuiltin(ctx, hookType, hookFunc2)
	
	hookData := &HookData{
		ExtensionName: "test-extension",
	}
	
	err := manager.Execute(ctx, hookType, hookData)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
	
	if executionCount != 2 {
		t.Errorf("Expected 2 executions, got %d", executionCount)
	}
}

func TestExecuteHookError(t *testing.T) {
	manager := New()
	ctx := context.Background()
	
	hookType := HookType("test-hook")
	expectedError := fmt.Errorf("hook error")
	
	hookFunc := func(ctx context.Context, data *HookData) error {
		return expectedError
	}
	
	manager.RegisterBuiltin(ctx, hookType, hookFunc)
	
	hookData := &HookData{
		ExtensionName: "test-extension",
	}
	
	err := manager.Execute(ctx, hookType, hookData)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	
	// The error is wrapped, so we check if it contains the expected error
	if err == nil || !strings.Contains(err.Error(), expectedError.Error()) {
		t.Errorf("Expected error containing '%v', got '%v'", expectedError, err)
	}
}

func TestListBuiltin(t *testing.T) {
	manager := New()
	ctx := context.Background()
	
	hookType := HookType("test-hook")
	
	// Initially should be empty
	hooks := manager.ListBuiltin(ctx, hookType)
	if len(hooks) != 0 {
		t.Errorf("Expected 0 hooks, got %d", len(hooks))
	}
	
	// Register a hook
	hookFunc := func(ctx context.Context, data *HookData) error {
		return nil
	}
	
	manager.RegisterBuiltin(ctx, hookType, hookFunc)
	
	// Should now have 1 hook
	hooks = manager.ListBuiltin(ctx, hookType)
	if len(hooks) != 1 {
		t.Errorf("Expected 1 hook, got %d", len(hooks))
	}
}

func TestCountBuiltin(t *testing.T) {
	manager := New()
	ctx := context.Background()
	
	hookType := HookType("test-hook")
	
	// Initially should be 0
	count := manager.CountBuiltin(ctx, hookType)
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
	
	// Register hooks
	hookFunc1 := func(ctx context.Context, data *HookData) error {
		return nil
	}
	
	hookFunc2 := func(ctx context.Context, data *HookData) error {
		return nil
	}
	
	manager.RegisterBuiltin(ctx, hookType, hookFunc1)
	manager.RegisterBuiltin(ctx, hookType, hookFunc2)
	
	// Should now have 2 hooks
	count = manager.CountBuiltin(ctx, hookType)
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestHookTypes(t *testing.T) {
	manager := New()
	ctx := context.Background()
	
	// Initially should be empty
	types := manager.HookTypes(ctx)
	if len(types) != 0 {
		t.Errorf("Expected 0 hook types, got %d", len(types))
	}
	
	// Register hooks of different types
	hookType1 := HookType("type1")
	hookType2 := HookType("type2")
	
	hookFunc := func(ctx context.Context, data *HookData) error {
		return nil
	}
	
	manager.RegisterBuiltin(ctx, hookType1, hookFunc)
	manager.RegisterBuiltin(ctx, hookType2, hookFunc)
	
	// Should now have 2 types
	types = manager.HookTypes(ctx)
	if len(types) != 2 {
		t.Errorf("Expected 2 hook types, got %d", len(types))
	}
}

func TestHookData(t *testing.T) {
	hookData := &HookData{
		ExtensionName: "test-extension",
		ExtensionPath: "/path/to/extension",
		Error:         fmt.Errorf("test error"),
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		},
	}
	
	if hookData.ExtensionName != "test-extension" {
		t.Errorf("Expected ExtensionName 'test-extension', got '%s'", hookData.ExtensionName)
	}
	
	if hookData.ExtensionPath != "/path/to/extension" {
		t.Errorf("Expected ExtensionPath '/path/to/extension', got '%s'", hookData.ExtensionPath)
	}
	
	if hookData.Error == nil {
		t.Error("Expected error to be set")
	}
	
	if hookData.Metadata["key1"] != "value1" {
		t.Error("Metadata key1 not set correctly")
	}
	
	if hookData.Metadata["key2"] != 123 {
		t.Error("Metadata key2 not set correctly")
	}
}

func TestHookType(t *testing.T) {
	hookType := HookType("test-hook")
	
	if string(hookType) != "test-hook" {
		t.Errorf("Expected hook type 'test-hook', got '%s'", string(hookType))
	}
}

func TestContextCancellation(t *testing.T) {
	manager := New()
	ctx := context.Background()
	
	hookType := HookType("test-hook")
	hookFunc := func(ctx context.Context, data *HookData) error {
		// Simple hook that respects context
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}
	
	manager.RegisterBuiltin(ctx, hookType, hookFunc)
	
	hookData := &HookData{
		ExtensionName: "test-extension",
	}
	
	err := manager.Execute(ctx, hookType, hookData)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
}

// Benchmark tests
func BenchmarkNewManager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		manager := New()
		_ = manager
	}
}

func BenchmarkRegisterBuiltin(b *testing.B) {
	manager := New()
	ctx := context.Background()
	hookType := HookType("test-hook")
	hookFunc := func(ctx context.Context, data *HookData) error {
		return nil
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.RegisterBuiltin(ctx, hookType, hookFunc)
	}
}

func BenchmarkExecuteBuiltin(b *testing.B) {
	manager := New()
	ctx := context.Background()
	hookType := HookType("test-hook")
	hookFunc := func(ctx context.Context, data *HookData) error {
		return nil
	}
	
	manager.RegisterBuiltin(ctx, hookType, hookFunc)
	hookData := &HookData{
		ExtensionName: "test-extension",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Execute(ctx, hookType, hookData)
	}
}

func BenchmarkListBuiltin(b *testing.B) {
	manager := New()
	ctx := context.Background()
	hookType := HookType("test-hook")
	hookFunc := func(ctx context.Context, data *HookData) error {
		return nil
	}
	
	// Register some hooks
	for i := 0; i < 10; i++ {
		manager.RegisterBuiltin(ctx, hookType, hookFunc)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.ListBuiltin(ctx, hookType)
	}
}