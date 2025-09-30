package hook

import (
	"context"
	"testing"
)

func TestNewBuiltinDispatcher(t *testing.T) {
	dispatcher := NewBuiltinDispatcher(nil)
	if dispatcher == nil {
		t.Fatal("NewBuiltinDispatcher() returned nil")
	}

	// Test internal fields directly
	if dispatcher.builtinExecutor == nil {
		t.Error("Builtin executor should not be nil")
	}

	if dispatcher.validator == nil {
		t.Error("Validator should not be nil")
	}
}

func TestDispatcherExecute(t *testing.T) {
	dispatcher := NewBuiltinDispatcher(nil)
	ctx := context.Background()

	// Register a test hook
	hookExecuted := false
	err := dispatcher.Register("before-install", func(ctx context.Context, data *Data) error {
		hookExecuted = true
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to register hook: %v", err)
	}

	// Create test data
	hookData := NewData("before-install", "test-extension")

	// Execute hooks
	err = dispatcher.Execute(ctx, "before-install", hookData)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !hookExecuted {
		t.Error("Hook was not executed")
	}
}

func TestDispatcherExecuteWithValidation(t *testing.T) {
	dispatcher := NewBuiltinDispatcher(nil)
	ctx := context.Background()

	// Test with invalid data (empty extension name)
	hookData := NewData("before-install", "")

	// Execute should fail due to validation
	err := dispatcher.Execute(ctx, "before-install", hookData)
	if err == nil {
		t.Error("Expected validation error, got nil")
	}

	// Test with valid data
	hookData = NewData("before-install", "test-extension")

	err = dispatcher.Execute(ctx, "before-install", hookData)
	if err != nil {
		t.Fatalf("Execute failed with valid data: %v", err)
	}
}

func TestBuiltinExecutor(t *testing.T) {
	executor := NewBuiltinExecutor()
	ctx := context.Background()

	// Register a hook
	hookExecuted := false
	executor.Register(ctx, "before-install", func(ctx context.Context, data *Data) error {
		hookExecuted = true
		return nil
	})

	// Execute hook
	hookData := NewData("before-install", "test-extension")
	err := executor.Execute(ctx, "before-install", hookData)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !hookExecuted {
		t.Error("Hook was not executed")
	}
}

func TestHookValidator(t *testing.T) {
	validator := NewValidator()

	// Test valid data
	validData := NewData("before-install", "test-extension")
	err := validator.Validate(validData)
	if err != nil {
		t.Errorf("Valid data should pass validation: %v", err)
	}

	// Test invalid data - empty hook type
	invalidData := &Data{
		Type:      "",
		Extension: "test-extension",
	}
	err = validator.Validate(invalidData)
	if err == nil {
		t.Error("Invalid data should fail validation")
	}

	// Test invalid data - empty extension name
	invalidData = &Data{
		Type:      "before-install",
		Extension: "",
	}
	err = validator.Validate(invalidData)
	if err == nil {
		t.Error("Invalid data should fail validation")
	}

	// Test invalid data - invalid extension name format
	invalidData = &Data{
		Type:      "before-install",
		Extension: "invalid@name!",
	}
	err = validator.Validate(invalidData)
	if err == nil {
		t.Error("Invalid data should fail validation")
	}
}

func TestHookDataBuilder(t *testing.T) {
	// Test basic creation
	hookData := NewData("before-install", "test-extension")
	if hookData.Type != "before-install" {
		t.Errorf("Expected hook type 'before-install', got %s", hookData.Type)
	}
	if hookData.Extension != "test-extension" {
		t.Errorf("Expected extension name 'test-extension', got '%s'", hookData.Extension)
	}

	// Test with metadata
	hookData = NewData("before-install", "test-extension").
		WithMetadata("version", "1.0.0").
		WithMetadata("author", "test")
	if hookData.Metadata["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%v'", hookData.Metadata["version"])
	}
	if hookData.Metadata["author"] != "test" {
		t.Errorf("Expected author 'test', got '%v'", hookData.Metadata["author"])
	}

	// Test with metadata map
	metadata := map[string]interface{}{
		"version": "2.0.0",
		"license": "MIT",
	}
	hookData = NewData("before-install", "test-extension").
		WithMetadataMap(metadata)
	if hookData.Metadata["version"] != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%v'", hookData.Metadata["version"])
	}
	if hookData.Metadata["license"] != "MIT" {
		t.Errorf("Expected license 'MIT', got '%v'", hookData.Metadata["license"])
	}
}

func TestDispatcherBuiltinHookManagement(t *testing.T) {
	dispatcher := NewBuiltinDispatcher(nil)

	// Test registering hooks
	hook1 := func(ctx context.Context, data *Data) error { return nil }
	hook2 := func(ctx context.Context, data *Data) error { return nil }

	err := dispatcher.Register("before-install", hook1)
	if err != nil {
		t.Fatalf("Failed to register first hook: %v", err)
	}

	err = dispatcher.Register("before-install", hook2)
	if err != nil {
		t.Fatalf("Failed to register second hook: %v", err)
	}

	// Test counting hooks
	count := dispatcher.CountBuiltinHooks("before-install")
	if count != 2 {
		t.Errorf("Expected 2 hooks, got %d", count)
	}

	// Test listing hook types
	types := dispatcher.ListBuiltinHooks()
	if len(types) != 1 || types[0] != "before-install" {
		t.Errorf("Expected [before-install], got %v", types)
	}

	// Test unregistering a hook
	err = dispatcher.Unregister("before-install", hook1)
	if err != nil {
		t.Fatalf("Failed to unregister hook: %v", err)
	}

	// Test counting after unregister
	count = dispatcher.CountBuiltinHooks("before-install")
	if count != 1 {
		t.Errorf("Expected 1 hook after unregister, got %d", count)
	}
}
