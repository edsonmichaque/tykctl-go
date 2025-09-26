# Hook Package

A flexible hook system for tykctl extensions that allows custom event handling and execution using functional options.

## Overview

The hook package provides a flexible, generic system for managing and executing hooks. Unlike traditional hook systems that define specific events, this package allows extensions to define their own event types and implement custom hook logic using functional options for configuration.

## Key Features

- **Functional Options**: Clean configuration using functional options pattern
- **Generic Event System**: Extensions define their own event types
- **Multiple Hook Types**: Support for builtin, external, and Rego hooks
- **Context Passing**: Rich context with custom data for extensions
- **Error Handling**: Proper error propagation and logging
- **Flexible Execution**: Extensions control when and how hooks are executed

## Usage

### Basic Hook Manager Creation

```go
import "github.com/edsonmichaque/tykctl-go/hook"

// Create a hook manager with default settings
hm := hook.New()

// Create a hook manager with custom external hook directory
hm := hook.New(
    hook.WithExternalHookDir("/path/to/hooks"),
)

// Create a hook manager with logger and custom directory
hm := hook.New(
    hook.WithExternalHookDir("/path/to/hooks"),
    hook.WithLogger(logger),
)
```

### Custom Event Handling

```go
// Define custom events for your extension
const (
    EventBeforeCreate hook.HookType = "before-create"
    EventAfterCreate  hook.HookType = "after-create"
    EventBeforeDelete hook.HookType = "before-delete"
)

// Create hook manager
hm := hook.New(
    hook.WithExternalHookDir("/path/to/hooks"),
    hook.WithLogger(logger),
)

// Register hooks for specific events
hm.Register(ctx, EventBeforeCreate, func(ctx context.Context, data interface{}) error {
    hookData, ok := data.(*hook.HookData)
    if !ok {
        return fmt.Errorf("invalid hook data type")
    }
    
    // Custom logic before create
    fmt.Printf("Before creating: %s\n", hookData.ExtensionName)
    return nil
})

// Execute hooks for an event
hookData := &hook.HookData{
    ExtensionName: "my-extension",
    ExtensionPath: "/path/to/extension",
    Metadata: map[string]interface{}{
        "resource": "user",
        "action":   "create",
    },
}

err := hm.Execute(ctx, EventBeforeCreate, hookData)
if err != nil {
    log.Printf("Hook execution failed: %v", err)
}
```

### Hook Data Structure

The `HookData` provides rich information for hook execution:

```go
type HookData struct {
    ExtensionName string                 // Name of the extension
    ExtensionPath string                 // Path to the extension
    Error         error                  // Any error that occurred
    Metadata      map[string]interface{} // Custom data for extensions
}
```

### Builtin Hooks

Builtin hooks are Go functions that are registered programmatically:

```go
// Register a builtin hook
hm.RegisterBuiltin(ctx, EventBeforeCreate, func(ctx context.Context, data interface{}) error {
    // Builtin hook logic
    return nil
})

// List builtin hooks for an event
builtinHooks := hm.ListBuiltin(ctx, EventBeforeCreate)

// Count builtin hooks
count := hm.CountBuiltin(ctx, EventBeforeCreate)
```

### External Hooks

External hooks are scripts or executables stored in the hook directory:

```go
// Create an external hook
externalHook, err := hm.CreateExternalHook(ctx, "my-hook", "#!/bin/bash\necho 'Hello from hook'")
if err != nil {
    log.Fatal(err)
}

// List external hooks
externalHooks, err := hm.ListExternal(ctx)
if err != nil {
    log.Fatal(err)
}

// Enable/disable external hooks
err = hm.EnableExternalHook(ctx, "my-hook")
err = hm.DisableExternalHook(ctx, "my-hook")

// Delete external hook
err = hm.DeleteExternalHook(ctx, "my-hook")
```

### Rego Hooks

Rego hooks use Open Policy Agent (OPA) for policy-based hook execution:

```go
// Register a Rego hook
regoHook := &hook.RegoHook{
    Name:    "validation-hook",
    Policy:  "package hooks\n\ndefault allow = false\n\nallow {\n    input.action == \"create\"\n}",
    Enabled: true,
}

err := hm.RegisterRegoHook(ctx, regoHook)
if err != nil {
    log.Fatal(err)
}

// Execute Rego hook
input := map[string]interface{}{
    "action": "create",
    "user":   "john",
}

result, err := hm.ExecuteRegoHook(ctx, "validation-hook", input)
if err != nil {
    log.Fatal(err)
}

if result.Allow {
    fmt.Println("Hook allows the action")
} else {
    fmt.Println("Hook denies the action")
}
```

## Configuration Options

### Functional Options

- `WithExternalHookDir(dir string)` - Set the external hook directory
- `WithLogger(logger *zap.Logger)` - Set the logger for the hook manager

### Default Configuration

- **External Hook Directory**: Uses XDG-based configuration directory (`~/.config/tykctl/hooks/`)
- **Logger**: Default no-op logger if not specified
- **Builtin Hooks**: Always available
- **External Hooks**: Available if directory is specified
- **Rego Hooks**: Available if logger is provided

## Environment Variables

- `TYKCTL_HOOKS_ENABLED` - Enable/disable hooks (default: true)
- `TYKCTL_HOOKS_DIRECTORY` - Custom hook directory
- `TYKCTL_HOOKS_TIMEOUT` - Hook execution timeout (default: 30s)
- `TYKCTL_HOOKS_MAX_RETRIES` - Maximum retry attempts (default: 3)

## Extension Integration

Extensions can integrate with the hook system by:

1. **Defining Events**: Create custom `HookType` constants for their domain
2. **Registering Hooks**: Use `Register()` to register hook functions
3. **Executing Hooks**: Call `Execute()` at appropriate points
4. **Providing Data**: Pass relevant data via `HookData`

### Example Extension Integration

```go
package myextension

import (
    "context"
    "fmt"
    
    "github.com/edsonmichaque/tykctl-go/hook"
)

// Define extension-specific events
const (
    EventUserCreate hook.HookType = "user-create"
    EventUserDelete hook.HookType = "user-delete"
)

type MyExtension struct {
    hookManager *hook.Manager
}

func NewMyExtension() *MyExtension {
    hm := hook.New(
        hook.WithExternalHookDir("/path/to/my/hooks"),
        hook.WithLogger(logger),
    )
    
    return &MyExtension{
        hookManager: hm,
    }
}

func (e *MyExtension) CreateUser(ctx context.Context, userData map[string]interface{}) error {
    // Execute before-create hooks
    hookData := &hook.HookData{
        ExtensionName: "my-extension",
        ExtensionPath: "/path/to/my/extension",
        Metadata: map[string]interface{}{
            "action": "create",
            "user":   userData,
        },
    }
    
    if err := e.hookManager.Execute(ctx, EventUserCreate, hookData); err != nil {
        return fmt.Errorf("pre-create hook failed: %w", err)
    }
    
    // Perform actual user creation
    // ...
    
    return nil
}
```

## Best Practices

- **Event Naming**: Use descriptive, namespaced event names (e.g., `user-before-create`)
- **Error Handling**: Always handle hook execution errors gracefully
- **Context Data**: Provide meaningful data in `HookData.Metadata`
- **Idempotency**: Design hooks to be idempotent where possible
- **Logging**: Use the provided logger for hook execution logging
- **Functional Options**: Use functional options for clean configuration
- **Type Safety**: Always validate `HookData` type in hook functions

## Examples

See the extension examples for practical usage:
- `tykctl-portal` - Uses hooks for user management events
- `tykctl-cloud` - Uses hooks for resource lifecycle events
- `tykctl-dashboard` - Uses hooks for configuration changes