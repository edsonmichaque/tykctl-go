# Hook Package

A flexible, Go-idiomatic hook system for tykctl extensions that allows custom event handling and execution using interfaces and concrete implementations.

## Overview

The hook package provides a clean, modular system for managing and executing hooks. Following Go best practices, it uses interfaces for extensibility and concrete implementations for specific functionality.

## Architecture

The package follows Go best practices with a clean, modular design:

### Executors
- **BuiltinExecutor** - executes builtin Go hooks
- **ScriptExecutor** - executes external script hooks with discovery
- **RegoExecutor** - executes Rego policy hooks with discovery
- **SchemaExecutor** - executes JSON Schema validation hooks with discovery

### Processors
- **BuiltinProcessor** - handles only builtin hooks
- **ScriptProcessor** - handles only script hooks with discovery
- **PolicyProcessor** - handles only Rego policy hooks with discovery
- **SchemaProcessor** - handles only JSON Schema validation hooks with discovery
- **Processor** - orchestrates multiple executors

### Validators
- **BuiltinValidator** - validates builtin hook data
- **ScriptValidator** - validates script hook data
- **PolicyValidator** - validates Rego policy hook data
- **SchemaValidator** - validates JSON Schema hook data
- **HookValidator** - general-purpose hook validator

## Key Features

- **Modular Design**: Specialized processors for different hook types
- **Automatic Discovery**: Script, policy, and schema files are automatically discovered
- **Variadic Parameters**: Flexible discovery methods that can find all or specific hook types
- **Context Support**: Full context.Context integration for cancellation and timeouts
- **Smart Caching**: Efficient caching of discovered files for better performance
- **Comprehensive Validation**: Built-in validation for all hook types
- **Structured Logging**: Detailed logging with zap for debugging and monitoring
- **Go Idioms**: Follows Go best practices and conventions

## Usage

### Constructor Patterns

All constructors follow consistent patterns with logger as the first parameter:

```go
// Executors
builtinExecutor := hook.NewBuiltinExecutor(logger)
scriptExecutor := hook.NewScriptExecutor(logger, "/tmp/hooks")
regoExecutor := hook.NewRegoExecutor(logger, "/tmp/policies")
schemaExecutor := hook.NewSchemaExecutor(logger, "/tmp/schemas")

// Processors
builtinProcessor := hook.NewBuiltinProcessor(logger)
scriptProcessor := hook.NewScriptProcessor(logger, "/tmp/hooks")
policyProcessor := hook.NewPolicyProcessor(logger, "/tmp/policies")
schemaProcessor := hook.NewSchemaProcessor(logger, "/tmp/schemas")

// Validators
builtinValidator := hook.NewBuiltinValidator()
scriptValidator := hook.NewScriptValidator("/tmp/hooks")
policyValidator := hook.NewPolicyValidator("/tmp/policies")
schemaValidator := hook.NewSchemaValidator("/tmp/schemas")
```

### Discovery API

All discovery methods support variadic parameters for flexible usage:

```go
// Discover all items
allScripts, err := scriptExecutor.discoverScripts(ctx)
allPolicies, err := regoExecutor.discoverPolicies(ctx)
allSchemas, err := schemaExecutor.discoverSchemas(ctx)

// Discover specific hook types
scripts, err := scriptExecutor.discoverScripts(ctx, "before-install")
policies, err := regoExecutor.discoverPolicies(ctx, "before-install", "after-install")
schemas, err := schemaExecutor.discoverSchemas(ctx, "before-install")
```

### File Discovery Patterns

The system supports flexible file discovery with consistent patterns:

#### Directory Structure
```
<hookDir>/
├── <hookType>/          # Directory containing multiple files
│   ├── file1
│   ├── file2
│   └── ...
└── <hookType>           # Single file named after hook type
```

#### Script Discovery
- **Naming**: Files must be named exactly like the hook type (with dashes)
- **Extension**: No file extension required
- **Permissions**: Must be executable
- **Order**: Files in directories are processed in lexicographic order

#### Policy Discovery (Rego)
- **Naming**: Files must be named exactly like the hook type (with dashes)
- **Extension**: Must have `.rego` extension
- **Order**: Files in directories are processed in lexicographic order

#### Schema Discovery (JSON Schema)
- **Naming**: Files must be named exactly like the hook type (with dashes)
- **Extension**: Must have `.json` extension
- **Validation**: Must be valid JSON Schema format
- **Order**: Files in directories are processed in lexicographic order

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
    // Create a builtin processor
    processor := hook.NewBuiltinProcessor(logger)
    ctx := context.Background()
    
    // Register some builtin hooks
    err := processor.Register("before-install", func(ctx context.Context, data *hook.Data) error {
        fmt.Printf("Installing extension: %s\n", data.Extension)
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
    builtinExecutor := hook.NewBuiltinExecutor(logger)
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
    processor.Register("before-install", func(ctx context.Context, data *hook.Data) error {
        logger.Info("Before install hook", zap.String("extension", data.Extension))
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

### Specialized Processors

The package provides focused processors, one for each executor type, with discovery capabilities:

#### Builtin Processor
```go
processor := hook.NewBuiltinProcessor(logger)

// Register builtin hooks
processor.Register("before-install", func(ctx context.Context, data *hook.Data) error {
    // Builtin hook logic
    return nil
})

// Execute hooks
err := processor.Execute(ctx, "before-install", hookData)
```

#### Script Processor
```go
processor := hook.NewScriptProcessor(logger, "/tmp/hooks")

// Discover scripts
scripts, err := processor.ListScripts(ctx, "before-install")
allScripts, err := processor.DiscoverAllScripts()

// Execute hooks (scripts are automatically discovered)
err := processor.Execute(ctx, "before-install", hookData)
```

#### Policy Processor
```go
processor := hook.NewPolicyProcessor(logger, "/tmp/policies")

// Discover policies
policies, err := processor.ListPolicies(ctx, "before-install")
allPolicies, err := processor.DiscoverAllPolicies()

// Execute hooks (policies are automatically discovered)
err := processor.Execute(ctx, "before-install", hookData)
```

#### Schema Processor
```go
processor := hook.NewSchemaProcessor(logger, "/tmp/schemas")

// Discover schemas
schemas, err := processor.ListJSONSchemas(ctx, "before-install")
allSchemas, err := processor.DiscoverAllJSONSchemas()

// Execute hooks (schemas are automatically discovered)
err := processor.Execute(ctx, "before-install", hookData)
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
executor := hook.NewRegoExecutor("/tmp/policies", logger)

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
