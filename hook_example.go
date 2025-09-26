package main

import (
	"context"
	"fmt"
	"log"

	"github.com/edsonmichaque/tykctl-go/extension"
	"github.com/edsonmichaque/tykctl-go/hook"
	"go.uber.org/zap"
)

// HookExample demonstrates how to use both builtin and external hooks
func main() {
	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a hook manager with custom external hook directory
	hookDir := "/tmp/tykctl-hooks"
	manager := hook.NewWithLogger(hookDir, logger)

	// Setup builtin hooks
	setupBuiltinHooks(manager)

	// Setup external hooks
	setupExternalHooks(manager)

	// Create extension installer with hooks
	configDir := "/tmp/tykctl-config"
	installer := extension.NewInstaller(configDir, extension.WithHooks(manager))

	ctx := context.Background()

	// Example: Install an extension (both builtin and external hooks will be executed)
	fmt.Println("Installing extension...")
	err := installer.InstallExtension(ctx, "example", "my-extension")
	if err != nil {
		log.Printf("Installation failed: %v", err)
		return
	}

	// List hooks
	listHooks(manager)

	fmt.Println("Hook example completed successfully!")
}

// setupBuiltinHooks configures builtin Go hooks
func setupBuiltinHooks(manager *hook.Manager) {
	ctx := context.Background()

	// Register builtin hooks directly
	manager.RegisterBuiltin(ctx, hook.HookTypeBeforeInstall, func(ctx context.Context, data *hook.HookData) error {
		log.Printf("🔧 Builtin Hook: Before installing %s", data.ExtensionName)
		return nil
	})

	manager.RegisterBuiltin(ctx, hook.HookTypeAfterInstall, func(ctx context.Context, data *hook.HookData) error {
		log.Printf("✅ Builtin Hook: After installing %s", data.ExtensionName)
		return nil
	})

	// Register predefined builtin hooks
	manager.RegisterBuiltin(ctx, hook.HookTypeBeforeInstall, hook.LoggingHook(hook.HookTypeBeforeInstall))
	manager.RegisterBuiltin(ctx, hook.HookTypeAfterInstall, hook.TimingHook(hook.HookTypeAfterInstall))

	// Register validation hook
	manager.RegisterBuiltin(ctx, hook.HookTypeBeforeInstall, hook.ValidationHook(func(data *hook.HookData) error {
		if data.ExtensionName == "" {
			return fmt.Errorf("extension name cannot be empty")
		}
		log.Printf("🔍 Builtin Hook: Validating extension %s", data.ExtensionName)
		return nil
	}))

	// Register metrics hook
	manager.RegisterBuiltin(ctx, hook.HookTypeAfterInstall, hook.MetricsHook(func(operation string, metrics map[string]interface{}) error {
		log.Printf("📊 Builtin Hook: Metrics for %s - %+v", operation, metrics)
		return nil
	}))
}

// setupExternalHooks configures external file-based hooks
func setupExternalHooks(manager *hook.Manager) {
	// Create external hooks (similar to Git hooks)
	createExternalHooks(manager)
}

// createExternalHooks creates external hook files
func createExternalHooks(manager *hook.Manager) {
	ctx := context.Background()

	// Create before-install hook
	beforeInstallScript := `#!/bin/bash
echo "🔧 External Hook: Before installing $TYKCTL_HOOK_EXTENSION"
echo "Event: $TYKCTL_HOOK_EVENT"
echo "Extension: $TYKCTL_HOOK_EXTENSION"
echo "Path: $TYKCTL_HOOK_PATH"
echo "Working Dir: $TYKCTL_HOOK_WORKING_DIR"
echo "External hook executed successfully!"
`

	_, err := manager.CreateExternalHook(ctx, "before-install", beforeInstallScript)
	if err != nil {
		log.Printf("Failed to create before-install hook: %v", err)
	}

	// Create after-install hook
	afterInstallScript := `#!/bin/bash
echo "✅ External Hook: After installing $TYKCTL_HOOK_EXTENSION"
echo "Event: $TYKCTL_HOOK_EVENT"
echo "Extension: $TYKCTL_HOOK_EXTENSION"
echo "Path: $TYKCTL_HOOK_PATH"
echo "External hook executed successfully!"
`

	_, err = manager.CreateExternalHook(ctx, "after-install", afterInstallScript)
	if err != nil {
		log.Printf("Failed to create after-install hook: %v", err)
	}

	// Create a Python hook
	pythonHook := `#!/usr/bin/env python3
import os
import sys

print("🐍 Python External Hook: Processing extension")
print(f"Event: {os.environ.get('TYKCTL_HOOK_EVENT', 'N/A')}")
print(f"Extension: {os.environ.get('TYKCTL_HOOK_EXTENSION', 'N/A')}")
print(f"Path: {os.environ.get('TYKCTL_HOOK_PATH', 'N/A')}")
print("Python external hook executed successfully!")
`

	_, err = manager.CreateExternalHook(ctx, "python-hook", pythonHook)
	if err != nil {
		log.Printf("Failed to create Python hook: %v", err)
	}

	// Create a validation hook
	validationScript := `#!/bin/bash
echo "🔍 External Hook: Validating installation"

# Check if extension name is provided
if [ -z "$TYKCTL_HOOK_EXTENSION" ]; then
    echo "ERROR: Extension name is required"
    exit 1
fi

# Check if working directory exists
if [ ! -d "$TYKCTL_HOOK_WORKING_DIR" ]; then
    echo "ERROR: Working directory does not exist: $TYKCTL_HOOK_WORKING_DIR"
    exit 1
fi

echo "External validation passed for extension: $TYKCTL_HOOK_EXTENSION"
`

	_, err = manager.CreateExternalHook(ctx, "validate-install", validationScript)
	if err != nil {
		log.Printf("Failed to create validation hook: %v", err)
	}
}

// listHooks demonstrates hook management
func listHooks(manager *hook.Manager) {
	ctx := context.Background()
	fmt.Println("\n=== Hook Management ===")

	// List builtin hooks
	fmt.Println("\nBuiltin Hooks:")
	for _, hookType := range manager.HookTypes(ctx) {
		count := manager.CountBuiltin(ctx, hookType)
		fmt.Printf("- %s: %d hooks\n", hookType, count)
	}

	// List external hooks
	fmt.Println("\nExternal Hooks:")
	externalHooks, err := manager.ListExternal(ctx)
	if err != nil {
		log.Printf("Failed to list external hooks: %v", err)
		return
	}

	for _, hook := range externalHooks {
		status := "disabled"
		if hook.Enabled {
			status = "enabled"
		}
		fmt.Printf("- %s: %s (%s)\n", hook.Name, hook.Path, status)
	}

	// Count total hooks
	builtinCount := 0
	for _, hookType := range manager.HookTypes(ctx) {
		builtinCount += manager.CountBuiltin(ctx, hookType)
	}

	externalCount, _ := manager.CountExternal(ctx)
	fmt.Printf("\nTotal: %d builtin hooks, %d external hooks\n", builtinCount, externalCount)
}

// BuiltinHookExample demonstrates builtin hook usage
func BuiltinHookExample() {
	manager := hook.New()
	ctx := context.Background()

	// Register a custom builtin hook
	manager.RegisterBuiltin(ctx, hook.HookTypeBeforeInstall, func(ctx context.Context, data *hook.HookData) error {
		fmt.Printf("Custom builtin hook: Installing %s\n", data.ExtensionName)
		return nil
	})

	// Register predefined hooks
	manager.RegisterBuiltin(ctx, hook.HookTypeBeforeInstall, hook.LoggingHook(hook.HookTypeBeforeInstall))
	manager.RegisterBuiltin(ctx, hook.HookTypeAfterInstall, hook.TimingHook(hook.HookTypeAfterInstall))

	// Execute hooks
	hookData := &hook.HookData{
		ExtensionName: "test-extension",
		ExtensionPath: "/path/to/extension",
		Metadata: map[string]interface{}{
			"version": "1.0.0",
		},
	}

	err := manager.Execute(ctx, hook.HookTypeBeforeInstall, hookData)
	if err != nil {
		log.Printf("Builtin hooks failed: %v", err)
		return
	}

	fmt.Println("Builtin hook example completed")
}

// ExternalHookExample demonstrates external hook usage
func ExternalHookExample() {
	logger, _ := zap.NewDevelopment()
	manager := hook.NewWithLogger("/tmp/tykctl-hooks", logger)
	ctx := context.Background()

	// Create external hooks
	script := `#!/bin/bash
echo "External hook executed for: $TYKCTL_HOOK_EXTENSION"
`

	_, err := manager.CreateExternalHook(ctx, "test-hook", script)
	if err != nil {
		log.Printf("Failed to create external hook: %v", err)
		return
	}

	// Enable the hook
	err = manager.EnableExternalHook(ctx, "test-hook")
	if err != nil {
		log.Printf("Failed to enable external hook: %v", err)
		return
	}

	// Execute hooks
	hookData := &hook.HookData{
		ExtensionName: "test-extension",
		ExtensionPath: "/path/to/extension",
	}

	err = manager.Execute(ctx, hook.HookTypeBeforeInstall, hookData)
	if err != nil {
		log.Printf("External hooks failed: %v", err)
		return
	}

	// Clean up
	manager.DeleteExternalHook(ctx, "test-hook")
	fmt.Println("External hook example completed")
}

// MixedHookExample demonstrates using both builtin and external hooks together
func MixedHookExample() {
	logger, _ := zap.NewDevelopment()
	manager := hook.NewWithLogger("/tmp/tykctl-hooks", logger)
	ctx := context.Background()

	// Register builtin hooks
	manager.RegisterBuiltin(ctx, hook.HookTypeBeforeInstall, func(ctx context.Context, data *hook.HookData) error {
		fmt.Println("Builtin: Before install")
		return nil
	})

	manager.RegisterBuiltin(ctx, hook.HookTypeAfterInstall, func(ctx context.Context, data *hook.HookData) error {
		fmt.Println("Builtin: After install")
		return nil
	})

	// Create external hooks
	beforeScript := `#!/bin/bash
echo "External: Before install"
`
	afterScript := `#!/bin/bash
echo "External: After install"
`

	manager.CreateExternalHook(ctx, "before-install", beforeScript)
	manager.CreateExternalHook(ctx, "after-install", afterScript)

	// Execute hooks
	hookData := &hook.HookData{
		ExtensionName: "mixed-extension",
		ExtensionPath: "/path/to/extension",
	}

	// Execute before install hooks (both builtin and external)
	err := manager.Execute(ctx, hook.HookTypeBeforeInstall, hookData)
	if err != nil {
		log.Printf("Before install hooks failed: %v", err)
		return
	}

	// Execute after install hooks (both builtin and external)
	err = manager.Execute(ctx, hook.HookTypeAfterInstall, hookData)
	if err != nil {
		log.Printf("After install hooks failed: %v", err)
		return
	}

	fmt.Println("Mixed hook example completed")
}
