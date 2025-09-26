package hook

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example demonstrates how to use the hook system
func Example() {
	// Create a hook manager
	manager := New()
	ctx := context.Background()

	// Define custom hook types for this example
	const (
		HookTypeBeforeInstall HookType = "before-install"
		HookTypeAfterInstall  HookType = "after-install"
		HookTypeBeforeRun     HookType = "before-run"
		HookTypeAfterRun      HookType = "after-run"
	)

	// Register some predefined hooks
	manager.RegisterBuiltin(ctx, HookTypeBeforeInstall, LoggingHook(HookTypeBeforeInstall))
	manager.RegisterBuiltin(ctx, HookTypeAfterInstall, LoggingHook(HookTypeAfterInstall))
	manager.RegisterBuiltin(ctx, HookTypeBeforeRun, TimingHook(HookTypeBeforeRun))
	manager.RegisterBuiltin(ctx, HookTypeAfterRun, TimingHook(HookTypeAfterRun))

	// Register a custom hook
	manager.RegisterBuiltin(ctx, HookTypeBeforeInstall, func(ctx context.Context, data *HookData) error {
		log.Printf("Custom hook: Installing extension %s", data.ExtensionName)
		return nil
	})

	// Register a validation hook
	manager.RegisterBuiltin(ctx, HookTypeBeforeInstall, ValidationHook(func(data *HookData) error {
		if data.ExtensionName == "" {
			return fmt.Errorf("extension name cannot be empty")
		}
		return nil
	}))

	// Register a metrics hook
	manager.RegisterBuiltin(ctx, HookTypeAfterInstall, MetricsHook(func(operation string, metrics map[string]interface{}) error {
		log.Printf("Metrics: %s - %+v", operation, metrics)
		return nil
	}))

	// Execute hooks
	hookData := &HookData{
		ExtensionName: "my-extension",
		ExtensionPath: "/path/to/extension",
		Metadata: map[string]interface{}{
			"version": "1.0.0",
			"author":  "example",
		},
	}

	// Execute before install hooks
	if err := manager.Execute(ctx, HookTypeBeforeInstall, hookData); err != nil {
		log.Printf("Before install hooks failed: %v", err)
		return
	}

	// Simulate installation work
	time.Sleep(100 * time.Millisecond)

	// Execute after install hooks
	if err := manager.Execute(ctx, HookTypeAfterInstall, hookData); err != nil {
		log.Printf("After install hooks failed: %v", err)
		return
	}

	// Execute hooks asynchronously
	errChan := manager.ExecuteAsync(ctx, HookTypeBeforeRun, hookData)
	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("Async hooks failed: %v", err)
		}
	case <-time.After(1 * time.Second):
		log.Printf("Async hooks timed out")
	}

	// List registered hooks
	fmt.Printf("Registered hook types: %v\n", manager.HookTypes(ctx))
	fmt.Printf("Before install hooks count: %d\n", manager.Count(ctx, HookTypeBeforeInstall))
	fmt.Printf("Has before install hooks: %t\n", manager.HasHooks(ctx, HookTypeBeforeInstall))
}
