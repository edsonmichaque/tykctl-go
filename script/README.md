# pkg/tykctl-go/script

A generic script-based hook system for tykctl extensions that allows custom event handling and script execution.

## Overview

The script package provides a flexible, generic system for managing and executing scripts. Unlike traditional hook systems that define specific events, this package allows extensions to define their own events and implement custom script logic.

## Key Features

- **Generic Event System**: Extensions define their own event types
- **Script Registry**: Register custom event handlers
- **Context Passing**: Rich context with custom data for extensions
- **Error Handling**: Proper error propagation and logging
- **Flexible Execution**: Extensions control when and how scripts are executed
- **File-based Scripts**: Execute actual script files with proper environment setup

## Usage

### Basic Script Management

```go
import "github.com/edsonmichaque/tykctl/pkg/tykctl-go/script"

// Create a script manager
sm := script.NewScriptManager(scriptDir)

// Create a script
s, err := sm.CreateScript("my-script", "Description", "#!/bin/bash\necho 'Hello World'")
if err != nil {
    log.Fatal(err)
}

// Enable the script
err = sm.EnableScript("my-script")
if err != nil {
    log.Fatal(err)
}
```

### Custom Event Handling

```go
// Define custom events
const (
    EventBeforeCreate ScriptEvent = "before-create"
    EventAfterCreate  ScriptEvent = "after-create"
    EventBeforeDelete ScriptEvent = "before-delete"
)

// Create script registry
registry := script.NewScriptRegistry()

// Register handlers
registry.RegisterHandler(EventBeforeCreate, func(ctx context.Context, scriptCtx *script.ScriptContext) error {
    // Custom logic before create
    return nil
})

// Execute scripts
scriptCtx := &script.ScriptContext{
    Event: EventBeforeCreate,
    Data:  map[string]interface{}{"resource": "user"},
}
err := registry.ExecuteHandlers(ctx, EventBeforeCreate, scriptCtx)
```

### Script Context

The `ScriptContext` provides rich information for script execution:

```go
type ScriptContext struct {
    Event       ScriptEvent              // Custom event type
    Command     string                   // Command being executed
    Args        []string                 // Command arguments
    Extension   string                   // Extension name
    WorkingDir  string                   // Working directory
    Environment map[string]string        // Environment variables
    Data        map[string]interface{}   // Custom data for extensions
}
```

## Configuration

Scripts are stored in the XDG-based configuration directory:

- **Default Location**: `~/.config/tykctl/scripts/`
- **Configurable**: Can be set via `TYKCTL_SCRIPTS_DIRECTORY` environment variable
- **Integration**: Uses the main tykctl configuration system

## Environment Variables

Scripts receive the following environment variables:

- `TYKCTL_SCRIPT_EVENT` - The event that triggered the script
- `TYKCTL_SCRIPT_COMMAND` - The command being executed
- `TYKCTL_SCRIPT_ARGS` - Space-separated command arguments
- `TYKCTL_SCRIPT_EXTENSION` - The extension name
- `TYKCTL_SCRIPT_WORKING_DIR` - The working directory

## Script Execution

Scripts are executed with proper context and timeout handling:

```go
// Execute a single script
scriptCtx := &script.ScriptContext{
    Event:      "before-install",
    Command:    "install",
    Args:       []string{"my-extension"},
    Extension:  "my-extension",
    WorkingDir: "/tmp",
    Data:       map[string]interface{}{"version": "1.0.0"},
}

err := sm.ExecuteScript(ctx, script, scriptCtx)
```

## Extension Integration

Extensions can integrate with the script system by:

1. **Defining Events**: Create custom event types for their domain
2. **Registering Handlers**: Use `ScriptRegistry` to register event handlers
3. **Executing Scripts**: Call `ExecuteHandlers` at appropriate points
4. **Providing Context**: Pass relevant data via `ScriptContext`

## Best Practices

- **Event Naming**: Use descriptive, namespaced event names (e.g., `user-before-create`)
- **Error Handling**: Always handle script execution errors gracefully
- **Context Data**: Provide meaningful data in `ScriptContext.Data`
- **Idempotency**: Design scripts to be idempotent where possible
- **Logging**: Use the provided logger for script execution logging
- **Script Permissions**: Ensure scripts are executable (chmod +x)

## Examples

### Basic Script Creation

```go
// Create a script manager
sm := script.NewScriptManagerWithLogger(scriptDir, logger)

// Create a script
scriptContent := `#!/bin/bash
echo "Installing extension: $TYKCTL_SCRIPT_EXTENSION"
echo "Event: $TYKCTL_SCRIPT_EVENT"
echo "Command: $TYKCTL_SCRIPT_COMMAND"
`

script, err := sm.CreateScript("install-hook", "Installation hook", scriptContent)
if err != nil {
    log.Fatal(err)
}
```

### Event-based Script Execution

```go
// Define events
const (
    EventBeforeInstall ScriptEvent = "before-install"
    EventAfterInstall  ScriptEvent = "after-install"
)

// Create registry
registry := script.NewScriptRegistry()

// Register handlers
registry.RegisterHandler(EventBeforeInstall, func(ctx context.Context, scriptCtx *script.ScriptContext) error {
    // Execute all scripts for this event
    return sm.ExecuteScriptsForEvent(ctx, EventBeforeInstall, scriptCtx)
})

// Execute handlers
scriptCtx := &script.ScriptContext{
    Event:     EventBeforeInstall,
    Extension: "my-extension",
    Data:      map[string]interface{}{"version": "1.0.0"},
}

err := registry.ExecuteHandlers(ctx, EventBeforeInstall, scriptCtx)
```

### Script Validation

```go
// Create validator
validator := script.NewScriptValidator(logger)

// Validate script
err := validator.ValidateScript(script)
if err != nil {
    log.Printf("Script validation failed: %v", err)
}
```


## Integration with tykctl-go

The script package integrates seamlessly with the main tykctl-go SDK:

```go
import (
    "github.com/edsonmichaque/tykctl/pkg/tykctl-go/extension"
    "github.com/edsonmichaque/tykctl/pkg/tykctl-go/script"
)

// Create script manager
sm := script.NewScriptManager(scriptDir)

// Create extension installer with script integration
installer := extension.NewInstaller(configDir)

// Execute scripts before installation
scriptCtx := &script.ScriptContext{
    Event:     "before-install",
    Extension: "my-extension",
}

err := sm.ExecuteScriptsForEvent(ctx, "before-install", scriptCtx)
if err != nil {
    return err
}

// Proceed with installation
err = installer.InstallExtension(ctx, "owner", "repo")
```
