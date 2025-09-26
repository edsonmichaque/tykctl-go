# pkg/hook

A generic hook system for tykctl extensions that allows custom event handling and execution.

## Overview

The hook package provides a flexible, generic system for managing and executing hooks. Unlike traditional hook systems that define specific events, this package allows extensions to define their own events and implement custom hook logic.

## Key Features

- **Generic Event System**: Extensions define their own event types
- **Hook Registry**: Register custom event handlers
- **Context Passing**: Rich context with custom data for extensions
- **Error Handling**: Proper error propagation and logging
- **Flexible Execution**: Extensions control when and how hooks are executed

## Usage

### Basic Hook Management

```go
// Create a hook manager
hm := hook.NewHookManager(hookDir)

// Create a hook
h, err := hm.CreateHook("my-hook", "Description", "script.sh")
if err != nil {
    log.Fatal(err)
}

// Enable the hook
err = hm.EnableHook("my-hook")
if err != nil {
    log.Fatal(err)
}
```

### Custom Event Handling

```go
// Define custom events
const (
    EventBeforeCreate HookEvent = "before-create"
    EventAfterCreate  HookEvent = "after-create"
    EventBeforeDelete HookEvent = "before-delete"
)

// Create hook registry
registry := hook.NewHookRegistry()

// Register handlers
registry.RegisterHandler(EventBeforeCreate, func(ctx context.Context, hookCtx *hook.HookContext) error {
    // Custom logic before create
    return nil
})

// Execute hooks
hookCtx := &hook.HookContext{
    Event: EventBeforeCreate,
    Data:  map[string]interface{}{"resource": "user"},
}
err := registry.ExecuteHandlers(ctx, EventBeforeCreate, hookCtx)
```

### Hook Context

The `HookContext` provides rich information for hook execution:

```go
type HookContext struct {
    Event       HookEvent              // Custom event type
    Command     string                 // Command being executed
    Args        []string               // Command arguments
    Extension   string                 // Extension name
    WorkingDir  string                 // Working directory
    Environment map[string]string      // Environment variables
    Data        map[string]interface{} // Custom data for extensions
}
```

## Configuration

Hooks are stored in the XDG-based configuration directory:

- **Default Location**: `~/.config/tykctl/hooks/`
- **Configurable**: Can be set via `TYKCTL_HOOKS_DIRECTORY` environment variable
- **Integration**: Uses the main tykctl configuration system

## Environment Variables

- `TYKCTL_HOOKS_ENABLED` - Enable/disable hooks (default: true)
- `TYKCTL_HOOKS_DIRECTORY` - Custom hook directory
- `TYKCTL_HOOKS_TIMEOUT` - Hook execution timeout (default: 30s)
- `TYKCTL_HOOKS_MAX_RETRIES` - Maximum retry attempts (default: 3)

## CLI Commands

The hook system is accessible via the main tykctl CLI:

```bash
# List all hooks
tykctl hook list

# Create a new hook
tykctl hook create my-hook --description "My hook" --script "script.sh"

# Enable a hook
tykctl hook enable my-hook

# Disable a hook
tykctl hook disable my-hook

# Delete a hook
tykctl hook delete my-hook

# Test a hook
tykctl hook test my-hook --event "test-event"
```

## Extension Integration

Extensions can integrate with the hook system by:

1. **Defining Events**: Create custom event types for their domain
2. **Registering Handlers**: Use `HookRegistry` to register event handlers
3. **Executing Hooks**: Call `ExecuteHandlers` at appropriate points
4. **Providing Context**: Pass relevant data via `HookContext`

## Best Practices

- **Event Naming**: Use descriptive, namespaced event names (e.g., `user-before-create`)
- **Error Handling**: Always handle hook execution errors gracefully
- **Context Data**: Provide meaningful data in `HookContext.Data`
- **Idempotency**: Design hooks to be idempotent where possible
- **Logging**: Use the provided logger for hook execution logging

## Examples

See the extension examples for practical usage:
- `tykctl-portal` - Uses hooks for user management events
- `tykctl-cloud` - Uses hooks for resource lifecycle events
- `tykctl-dashboard` - Uses hooks for configuration changes
