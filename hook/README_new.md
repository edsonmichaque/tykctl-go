# Hook Package

A flexible, Go-idiomatic hook system for tykctl extensions that allows custom event handling and execution using interfaces and concrete implementations.

## Overview

The hook package provides a clean, modular system for managing and executing hooks. Following Go best practices, it uses interfaces for extensibility and concrete implementations for specific functionality.

## Architecture

The package follows the same Go-idiomatic pattern as the template package:

- **Executor** (interface) - for executing hooks
- **BuiltinExecutor** (concrete) - for builtin Go hooks
- **ExternalExecutor** (concrete) - for external script hooks
- **RegoExecutor** (concrete) - for Rego policy hooks
- **Validator** (interface) - for validating hook data
- **HookValidator** (concrete) - for validating hooks
- **Processor** (concrete) - orchestrates all executors

## Key Features

- **Interface-Based Design**: Clean interfaces for extensibility
- **Multiple Hook Types**: Support for builtin, external, and Rego hooks
- **Context Support**: Full context.Context integration for cancellation and timeouts
- **Validation**: Built-in validation for hook data
- **Error Handling**: Comprehensive error handling with structured error types
- **Go Idioms**: Follows Go best practices and conventions

## Usage

### Basic Hook Processing

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/hook"
)

func main() {
    // Create a simple processor with only builtin hooks
    processor := hook.NewSimpleProcessor(nil)
    ctx := context.Background()
    
    // Register some builtin hooks
    err := processor.RegisterBuiltin(ctx, "before-install", func(ctx context.Context, data *hook.Data) error {
        fmt.Printf("Installing extension: %s\n", data.ExtensionName)
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Create hook data
    hookData := hook.NewData("before-install", "my-extension").
        WithMetadata("version", "1.0.0")
    
    // Execute hooks
    err = processor.Execute(ctx, "before-install", hookData)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Full Hook Processing

```go
func fullHookProcessing() {
    logger, _ := zap.NewDevelopment()
    
    // Create a processor with script and policy executors
    builtinExecutor := hook.NewBuiltinExecutor()
    validator := hook.NewValidator()
    scriptExecutor := hook.NewScriptExecutor(logger, "/tmp/hooks")
    regoExecutor := hook.NewRegoExecutor(logger, "/tmp/policies")
    
    processor := hook.NewProcessor(
        logger,
        validator,
        builtinExecutor,
        scriptExecutor,
        nil, // schemaExecutor
        regoExecutor,
    )
    ctx := context.Background()
    
    // Register builtin hooks
    processor.RegisterBuiltin(ctx, "before-install", func(ctx context.Context, data *hook.Data) error {
        logger.Info("Before install hook", zap.String("extension", data.ExtensionName))
        return nil
    })
    
    // Create hook data
    hookData := hook.NewData("before-install", "my-extension").
        WithMetadataMap(map[string]interface{}{
            "version": "1.0.0",
            "author":  "example",
        })
    
    // Execute hooks (builtin, external, and Rego)
    err := processor.Execute(ctx, "before-install", hookData)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Processor Types

The package provides focused processors, one for each executor type:

#### Builtin Processor
```go
// Only builtin Go hooks
processor := hook.NewBuiltinProcessor(logger)
```

#### Script Processor
```go
// Only external script hooks
processor := hook.NewScriptProcessor(logger, "/tmp/hooks")
```

#### Schema Processor
```go
// Only schema validation hooks
processor := hook.NewSchemaProcessor(logger, "/tmp/schemas")
```

#### Policy Processor
```go
// Only Rego policy hooks
processor := hook.NewPolicyProcessor(logger, "/tmp/policies")
```

#### Custom Processor
```go
// Create a custom processor with specific executors
builtinExecutor := hook.NewBuiltinExecutor()
validator := hook.NewValidator()
scriptExecutor := hook.NewScriptExecutor(logger, "/tmp/hooks")
schemaExecutor := hook.NewSchemaExecutor(logger, "/tmp/schemas")
regoExecutor := hook.NewRegoExecutor(logger, "/tmp/policies")

processor := hook.NewProcessor(
    logger,
    validator,
    builtinExecutor,
    scriptExecutor,
    schemaExecutor,
    regoExecutor,
)
```

### Custom Executors

```go
func customExecutors() {
    // Create individual executors
    builtinExecutor := hook.NewBuiltinExecutor()
    scriptExecutor := hook.NewScriptExecutor(nil, "/tmp/hooks")
    regoExecutor := hook.NewRegoExecutor(nil, "/tmp/policies")
    validator := hook.NewValidator()
    
    // Create processor with custom executors
    processor := hook.NewProcessor(
        nil, // logger
        validator,
        builtinExecutor,
        scriptExecutor,
        nil, // schema executor
        regoExecutor,
    )
    
    // Use processor...
}
```

## Hook Types

The package uses string literals for hook types. Common hook types include:

- `"before-install"` - Before extension installation
- `"after-install"` - After extension installation
- `"before-uninstall"` - Before extension uninstallation
- `"after-uninstall"` - After extension uninstallation
- `"before-run"` - Before extension execution
- `"after-run"` - After extension execution
- `"before-create"` - Before resource creation
- `"after-create"` - After resource creation
- `"before-update"` - Before resource update
- `"after-update"` - After resource update
- `"before-delete"` - Before resource deletion
- `"after-delete"` - After resource deletion

You can use any string as a hook type, making the system flexible for custom use cases.

## Hook Data

Hook data is structured and validated:

```go
type Data struct {
    Type          Type                   `json:"type"`
    ExtensionName string                 `json:"extension"`
    Error         error                  `json:"error,omitempty"`
    Metadata      map[string]interface{} `json:"metadata,omitempty"`
}
```

### Builder Pattern

Use the builder pattern for creating hook data:

```go
hookData := hook.NewData("before-install", "my-extension").
    WithExtensionPath("/path/to/extension").
    WithError(someError).
    WithMetadata("version", "1.0.0").
    WithMetadataMap(map[string]interface{}{
        "author": "example",
        "license": "MIT",
    })
```

## Executors

### BuiltinExecutor

Executes Go functions as hooks:

```go
executor := hook.NewBuiltinExecutor()

// Register a hook
executor.Register(ctx, "before-install", func(ctx context.Context, data *hook.Data) error {
    // Hook logic here
    return nil
})

// Execute hooks
err := executor.Execute(ctx, "before-install", hookData)
```

### ExternalExecutor

Executes external script hooks:

```go
executor := hook.NewExternalExecutor("/tmp/hooks", logger)

// Hooks are automatically discovered from the directory
// Files should be named like: before-install.sh, after-install.py, etc.
err := executor.Execute(ctx, "before-install", hookData)
```

### RegoExecutor

Executes Rego policy hooks:

```go
executor := hook.NewRegoExecutor(logger, "/tmp/policies")

// Policies are automatically discovered from the directory
// Files should be named like: before-install.rego, after-install.rego, etc.
err := executor.Execute(ctx, "before-install", hookData)
```

## Validation

The package includes built-in validation:

```go
validator := hook.NewValidator()

// Validate hook data
err := validator.Validate(hookData)
if err != nil {
    // Handle validation error
}
```

### Validation Rules

- Hook type is required
- Extension name is required and must be alphanumeric with hyphens/underscores
- Extension path must be valid if provided
- Metadata keys must be alphanumeric with hyphens/underscores
- Metadata values must be valid types (string, int, float64, bool, array, object)

## Error Handling

The package provides structured error handling:

```go
// Check error types
if hook.IsValidationError(err) {
    // Handle validation error
}

if hook.IsExecutorError(err) {
    // Handle executor error
}

// Wrap errors with context
err = hook.WrapError(err, "failed to execute hooks")
err = hook.WrapErrorf(err, "failed to execute hooks for %s", hookType)
```

### Error Types

- **ValidationError**: Validation-specific errors
- **ExecutorError**: Executor-specific errors
- **Standard Errors**: Common error variables

## Integration Examples

### With Extension Management

```go
func installExtension(processor *hook.Processor, extensionName string) error {
    ctx := context.Background()
    
    // Execute before install hooks
    hookData := hook.NewData("before-install", extensionName)
    if err := processor.Execute(ctx, "before-install", hookData); err != nil {
        return fmt.Errorf("before install hooks failed: %w", err)
    }
    
    // Perform installation
    // ... installation logic ...
    
    // Execute after install hooks
    hookData.Type = "after-install"
    if err := processor.Execute(ctx, "after-install", hookData); err != nil {
        return fmt.Errorf("after install hooks failed: %w", err)
    }
    
    return nil
}
```

### With CLI Commands

```go
func createInstallCommand(processor *hook.Processor) *cobra.Command {
    return &cobra.Command{
        Use:   "install <extension>",
        Short: "Install an extension",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            return installExtension(processor, args[0])
        },
    }
}
```

## Best Practices

- **Use Context**: Always pass context for cancellation and timeouts
- **Validate Data**: Use the built-in validator for hook data
- **Handle Errors**: Use structured error handling
- **Builder Pattern**: Use the builder pattern for creating hook data
- **Interface Segregation**: Use specific interfaces when possible
- **Dependency Injection**: Inject dependencies through constructors

## Migration Guide

The package has been redesigned with a modern, modular architecture:

```go
// Use specialized processors for specific hook types
builtinProcessor := hook.NewBuiltinProcessor(logger)
scriptProcessor := hook.NewScriptProcessor(logger, "/tmp/hooks")
policyProcessor := hook.NewPolicyProcessor(logger, "/tmp/policies")
schemaProcessor := hook.NewSchemaProcessor(logger, "/tmp/schemas")

// Or use the general processor for multiple hook types
processor := hook.NewProcessor(logger, validator, builtinExecutor, scriptExecutor, schemaExecutor, regoExecutor)
```

## Dependencies

- `go.uber.org/zap`: Structured logging
- `context`: Context support
- `os/exec`: External script execution

## Testing

The package includes comprehensive tests:

```bash
go test -v ./hook
```

## Examples

See the `example.go` and `rego_examples.go` files for complete examples.
