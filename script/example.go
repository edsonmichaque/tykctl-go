package script

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
)

// Example demonstrates how to use the script system
func Example() {
	// Create a logger
	logger, _ := zap.NewDevelopment()

	// Create a script manager
	scriptDir := "/tmp/tykctl-scripts"
	sm := NewScriptManagerWithLogger(scriptDir, logger)

	// Define custom events
	const (
		EventBeforeInstall ScriptEvent = "before-install"
		EventAfterInstall  ScriptEvent = "after-install"
		EventBeforeDelete  ScriptEvent = "before-delete"
		EventAfterDelete   ScriptEvent = "after-delete"
	)

	// Create a script registry
	registry := NewScriptRegistry()

	// Register handlers
	registry.RegisterHandler(EventBeforeInstall, func(ctx context.Context, scriptCtx *ScriptContext) error {
		log.Printf("Before install handler: %s", scriptCtx.Extension)
		return sm.ExecuteScriptsForEvent(ctx, EventBeforeInstall, scriptCtx)
	})

	registry.RegisterHandler(EventAfterInstall, func(ctx context.Context, scriptCtx *ScriptContext) error {
		log.Printf("After install handler: %s", scriptCtx.Extension)
		return sm.ExecuteScriptsForEvent(ctx, EventAfterInstall, scriptCtx)
	})

	// Create some example scripts
	createExampleScripts(sm)

	// Create script context
	scriptCtx := &ScriptContext{
		Event:       EventBeforeInstall,
		Command:     "install",
		Args:        []string{"my-extension"},
		Extension:   "my-extension",
		WorkingDir:  "/tmp",
		Environment: map[string]string{"TYKCTL_DEBUG": "true"},
		Data: map[string]interface{}{
			"version": "1.0.0",
			"author":  "example",
		},
	}

	// Execute handlers
	ctx := context.Background()
	err := registry.ExecuteHandlers(ctx, EventBeforeInstall, scriptCtx)
	if err != nil {
		log.Printf("Handler execution failed: %v", err)
	}

	// List available scripts
	scripts, err := sm.ListScripts()
	if err != nil {
		log.Printf("Failed to list scripts: %v", err)
		return
	}

	fmt.Printf("Available scripts: %d\n", len(scripts))
	for _, script := range scripts {
		fmt.Printf("- %s: %s (enabled: %t)\n", script.Name, script.Description, script.Enabled)
	}
}

// createExampleScripts creates some example scripts for demonstration
func createExampleScripts(sm *ScriptManager) {
	// Create a simple bash script
	bashScript := `#!/bin/bash
echo "=== Script: $0 ==="
echo "Event: $TYKCTL_SCRIPT_EVENT"
echo "Command: $TYKCTL_SCRIPT_COMMAND"
echo "Args: $TYKCTL_SCRIPT_ARGS"
echo "Extension: $TYKCTL_SCRIPT_EXTENSION"
echo "Working Dir: $TYKCTL_SCRIPT_WORKING_DIR"
echo "Debug: $TYKCTL_DEBUG"
echo "Script executed successfully!"
`

	_, err := sm.CreateScript("example-bash", "Example bash script", bashScript)
	if err != nil {
		log.Printf("Failed to create bash script: %v", err)
	}

	// Create a Python script
	pythonScript := `#!/usr/bin/env python3
import os
import sys

print("=== Python Script ===")
print(f"Event: {os.environ.get('TYKCTL_SCRIPT_EVENT', 'N/A')}")
print(f"Command: {os.environ.get('TYKCTL_SCRIPT_COMMAND', 'N/A')}")
print(f"Args: {os.environ.get('TYKCTL_SCRIPT_ARGS', 'N/A')}")
print(f"Extension: {os.environ.get('TYKCTL_SCRIPT_EXTENSION', 'N/A')}")
print(f"Working Dir: {os.environ.get('TYKCTL_SCRIPT_WORKING_DIR', 'N/A')}")
print(f"Debug: {os.environ.get('TYKCTL_DEBUG', 'N/A')}")
print("Python script executed successfully!")
`

	_, err = sm.CreateScript("example-python", "Example Python script", pythonScript)
	if err != nil {
		log.Printf("Failed to create Python script: %v", err)
	}

	// Create a validation script
	validationScript := `#!/bin/bash
echo "Validating extension installation..."

# Check if extension name is provided
if [ -z "$TYKCTL_SCRIPT_EXTENSION" ]; then
    echo "ERROR: Extension name is required"
    exit 1
fi

# Check if working directory exists
if [ ! -d "$TYKCTL_SCRIPT_WORKING_DIR" ]; then
    echo "ERROR: Working directory does not exist: $TYKCTL_SCRIPT_WORKING_DIR"
    exit 1
fi

echo "Validation passed for extension: $TYKCTL_SCRIPT_EXTENSION"
`

	_, err = sm.CreateScript("validate-install", "Installation validation script", validationScript)
	if err != nil {
		log.Printf("Failed to create validation script: %v", err)
	}
}

// ScriptManagementExample demonstrates script management operations
func ScriptManagementExample() {
	logger, _ := zap.NewDevelopment()
	sm := NewScriptManagerWithLogger("/tmp/tykctl-scripts", logger)

	// Create a script
	_, err := sm.CreateScript("test-script", "Test script", `#!/bin/bash
echo "Test script executed"
`)
	if err != nil {
		log.Printf("Failed to create script: %v", err)
		return
	}

	// Get the script
	retrievedScript, err := sm.GetScript("test-script")
	if err != nil {
		log.Printf("Failed to get script: %v", err)
		return
	}

	fmt.Printf("Retrieved script: %s\n", retrievedScript.Name)

	// Disable the script
	err = sm.DisableScript("test-script")
	if err != nil {
		log.Printf("Failed to disable script: %v", err)
		return
	}

	// Enable the script
	err = sm.EnableScript("test-script")
	if err != nil {
		log.Printf("Failed to enable script: %v", err)
		return
	}

	// List scripts
	scripts, err := sm.ListScripts()
	if err != nil {
		log.Printf("Failed to list scripts: %v", err)
		return
	}

	fmt.Printf("Total scripts: %d\n", len(scripts))

	// Delete the script
	err = sm.DeleteScript("test-script")
	if err != nil {
		log.Printf("Failed to delete script: %v", err)
		return
	}

	fmt.Println("Script management example completed")
}

// ScriptExecutionExample demonstrates script execution with different contexts
func ScriptExecutionExample() {
	logger, _ := zap.NewDevelopment()
	sm := NewScriptManagerWithLogger("/tmp/tykctl-scripts", logger)

	// Create a test script
	scriptContent := `#!/bin/bash
echo "Script execution started"
echo "Event: $TYKCTL_SCRIPT_EVENT"
echo "Extension: $TYKCTL_SCRIPT_EXTENSION"
echo "Command: $TYKCTL_SCRIPT_COMMAND"
echo "Args: $TYKCTL_SCRIPT_ARGS"
echo "Working Dir: $TYKCTL_SCRIPT_WORKING_DIR"
echo "Custom Data: $TYKCTL_SCRIPT_CUSTOM_DATA"
echo "Script execution completed"
`

	script, err := sm.CreateScript("execution-test", "Execution test script", scriptContent)
	if err != nil {
		log.Printf("Failed to create script: %v", err)
		return
	}

	// Create script context
	scriptCtx := &ScriptContext{
		Event:      "test-event",
		Command:    "test",
		Args:       []string{"arg1", "arg2"},
		Extension:  "test-extension",
		WorkingDir: "/tmp",
		Environment: map[string]string{
			"TYKCTL_SCRIPT_CUSTOM_DATA": "custom-value",
		},
		Data: map[string]interface{}{
			"version": "1.0.0",
		},
	}

	// Execute the script
	ctx := context.Background()
	err = sm.ExecuteScript(ctx, script, scriptCtx)
	if err != nil {
		log.Printf("Script execution failed: %v", err)
		return
	}

	fmt.Println("Script execution example completed")
}
