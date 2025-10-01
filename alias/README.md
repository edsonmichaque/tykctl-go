# TykCtl-Go Alias Package

A comprehensive alias management system for TykCtl extensions that provides command shortcuts, shell integration, and parameter expansion.

## Overview

The `alias` package provides a unified interface for managing command aliases across TykCtl extensions. It handles alias creation, execution, validation, and registration with Cobra commands.

## Features

- **Command Aliases**: Create shortcuts for commonly used commands
- **Shell Integration**: Execute shell commands from aliases (prefixed with `!`)
- **Parameter Expansion**: Support for `$1`, `$2`, `$*`, `$@` parameter substitution
- **Validation**: Comprehensive alias name and expansion validation
- **Cobra Integration**: Seamless integration with Cobra command framework
- **Config Providers**: Flexible configuration storage (in-memory, extension-specific)
- **Conflict Detection**: Automatic detection of command name conflicts

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "github.com/edsonmichaque/tykctl-go/alias"
    "github.com/spf13/cobra"
)

func main() {
    // Create config provider
    configProvider := alias.NewInMemoryConfigProvider()
    
    // Create alias manager
    manager := alias.NewManager(configProvider, []string{"help", "version"})
    
    // Create command builder
    builder := alias.NewCommandBuilder(manager)
    
    // Build alias command
    aliasCmd := builder.BuildAliasCommand()
    
    // Add to root command
    rootCmd := &cobra.Command{Use: "myapp"}
    rootCmd.AddCommand(aliasCmd)
    
    // Register aliases as subcommands
    registrar := alias.NewRegistrar(manager)
    ctx := context.Background()
    registrar.RegisterAliases(ctx, rootCmd)
}
```

### Extension Integration

```go
// Create extension-specific config provider
configProvider := alias.NewExtensionConfigProvider(
    setAliasFunc,    // Your extension's SetAlias function
    getAliasFunc,    // Your extension's GetAlias function
    deleteAliasFunc, // Your extension's DeleteAlias function
    listAliasesFunc, // Your extension's ListAliases function
)

// Create manager with reserved command names
reservedNames := []string{
    "products", "users", "organizations", "apps", 
    "auth", "context", "configure", "version",
}

manager := alias.NewManager(configProvider, reservedNames)
```

## Alias Types

### Command Aliases

Simple command shortcuts:

```bash
# Set alias
tykctl alias set co "products checkout"
tykctl alias set users "users list"

# Use alias
tykctl co my-product
tykctl users
```

### Shell Aliases

Execute shell commands (prefixed with `!`):

```bash
# Set shell alias
tykctl alias set cleanup "!rm -rf /tmp/tykctl-*"
tykctl alias set logs "!tail -f /var/log/tykctl.log"

# Use shell alias
tykctl cleanup
tykctl logs
```

### Parameterized Aliases

Use parameters with `$1`, `$2`, etc.:

```bash
# Set parameterized alias
tykctl alias set getuser "users get $1"
tykctl alias set deploy "products deploy $1 --env $2"

# Use with parameters
tykctl getuser john@example.com
tykctl deploy my-api production
```

### Complex Aliases

Multi-command workflows:

```bash
# Set complex alias
tykctl alias set setup "products create $1 && products checkout $1 && products publish $1"

# Use complex alias
tykctl setup my-new-api
```

## API Reference

### Manager

The `Manager` struct handles core alias operations:

```go
type Manager struct {
    configProvider ConfigProvider
    reservedNames  []string
}

// SetAlias sets an alias
func (m *Manager) SetAlias(ctx context.Context, name, expansion string) error

// GetAlias retrieves an alias
func (m *Manager) GetAlias(ctx context.Context, name string) (string, bool)

// DeleteAlias deletes an alias
func (m *Manager) DeleteAlias(ctx context.Context, name string) error

// ListAliases returns all aliases
func (m *Manager) ListAliases(ctx context.Context) map[string]string

// ExecuteAlias executes an alias
func (m *Manager) ExecuteAlias(ctx context.Context, aliasName string, args []string) error
```

### CommandBuilder

The `CommandBuilder` creates Cobra commands:

```go
type CommandBuilder struct {
    manager *Manager
}

// BuildAliasCommand creates the main alias command
func (cb *CommandBuilder) BuildAliasCommand() *cobra.Command

// BuildSetCommand creates the alias set command
func (cb *CommandBuilder) BuildSetCommand() *cobra.Command

// BuildListCommand creates the alias list command
func (cb *CommandBuilder) BuildListCommand() *cobra.Command

// BuildDeleteCommand creates the alias delete command
func (cb *CommandBuilder) BuildDeleteCommand() *cobra.Command
```

### Registrar

The `Registrar` handles alias registration:

```go
type Registrar struct {
    manager *Manager
}

// RegisterAliases registers aliases as subcommands
func (r *Registrar) RegisterAliases(ctx context.Context, rootCmd *cobra.Command) error

// RegisterAliasesWithValidation registers with conflict detection
func (r *Registrar) RegisterAliasesWithValidation(ctx context.Context, rootCmd *cobra.Command, reservedNames []string) error

// ValidateAliases validates all configured aliases
func (r *Registrar) ValidateAliases(ctx context.Context, reservedNames []string) []ValidationError
```

## Configuration Providers

### InMemoryConfigProvider

Stores aliases in memory (useful for testing):

```go
provider := alias.NewInMemoryConfigProvider()
```

### ExtensionConfigProvider

Integrates with extension configuration:

```go
provider := alias.NewExtensionConfigProvider(
    setAliasFunc,    // func(ctx context.Context, name, expansion string) error
    getAliasFunc,    // func(ctx context.Context, name string) (string, bool)
    deleteAliasFunc, // func(ctx context.Context, name string) error
    listAliasesFunc, // func(ctx context.Context) map[string]string
)
```

### ConfigProviderBuilder

Builds config providers with a fluent interface:

```go
provider := alias.NewConfigProviderBuilder().
    WithInMemory().
    Build()

// Or with custom provider
provider := alias.NewConfigProviderBuilder().
    WithCustom(myCustomProvider).
    Build()
```

## Parameter Expansion

The alias system supports several parameter expansion patterns:

- **`$1`, `$2`, etc.** - Individual arguments
- **`$*`** - All arguments as a single string
- **`$@`** - All arguments as separate words

### Examples

```bash
# Alias definition
tykctl alias set deploy "products deploy $1 --env $2 --confirm"

# Usage
tykctl deploy my-api production

# Expands to:
# products deploy my-api --env production --confirm
```

## Validation

### Alias Name Validation

- Cannot be empty
- Cannot contain whitespace
- Cannot contain shell metacharacters (`&|;()<>`)
- Cannot conflict with reserved command names

### Expansion Validation

- Cannot be empty
- Cannot start with whitespace

### Reserved Names

Common reserved command names:

```go
reservedNames := []string{
    "help", "version", "config", "alias",
    "products", "users", "organizations", "apps",
    "auth", "context", "configure",
}
```

## Error Handling

The alias system provides comprehensive error handling:

```go
// Check if alias exists
if _, exists := manager.GetAlias(ctx, "myalias"); !exists {
    return fmt.Errorf("alias 'myalias' not found")
}

// Validate before setting
if err := manager.SetAlias(ctx, "invalid name", "command"); err != nil {
    return fmt.Errorf("invalid alias: %w", err)
}

// Handle execution errors
if err := manager.ExecuteAlias(ctx, "myalias", []string{"arg1"}); err != nil {
    return fmt.Errorf("alias execution failed: %w", err)
}
```

## Testing

### Unit Testing

```go
func TestAliasManager(t *testing.T) {
    // Create in-memory provider for testing
    provider := alias.NewInMemoryConfigProvider()
    manager := alias.NewManager(provider, []string{"help"})
    
    ctx := context.Background()
    
    // Test setting alias
    err := manager.SetAlias(ctx, "test", "echo hello")
    assert.NoError(t, err)
    
    // Test getting alias
    expansion, exists := manager.GetAlias(ctx, "test")
    assert.True(t, exists)
    assert.Equal(t, "echo hello", expansion)
    
    // Test listing aliases
    aliases := manager.ListAliases(ctx)
    assert.Len(t, aliases, 1)
    assert.Equal(t, "echo hello", aliases["test"])
}
```

### Integration Testing

```go
func TestAliasExecution(t *testing.T) {
    provider := alias.NewInMemoryConfigProvider()
    manager := alias.NewManager(provider, []string{})
    
    ctx := context.Background()
    
    // Set alias
    err := manager.SetAlias(ctx, "test", "echo $1")
    assert.NoError(t, err)
    
    // Test parameter expansion
    preview := manager.ExpandAliasPreview("echo $1", []string{"hello"})
    assert.Equal(t, "echo hello", preview)
}
```

## Best Practices

### 1. Reserved Names

Always define reserved command names to prevent conflicts:

```go
reservedNames := []string{
    "help", "version", "config", "alias",
    // Add your extension's command names
    "products", "users", "organizations",
}
```

### 2. Validation

Validate aliases before registration:

```go
registrar := alias.NewRegistrar(manager)
errors := registrar.ValidateAliases(ctx, reservedNames)
if len(errors) > 0 {
    for _, err := range errors {
        log.Printf("Alias validation error: %s", err)
    }
}
```

### 3. Error Handling

Provide clear error messages for alias operations:

```go
if err := manager.SetAlias(ctx, name, expansion); err != nil {
    return fmt.Errorf("failed to set alias '%s': %w", name, err)
}
```

### 4. Documentation

Document complex aliases:

```bash
# Complex deployment alias
tykctl alias set deploy "products create $1 && products checkout $1 && products publish $1"
```

## Examples

### Complete Extension Integration

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/alias"
    "github.com/spf13/cobra"
)

func main() {
    // Create config provider (integrate with your extension's config)
    configProvider := createExtensionConfigProvider()
    
    // Define reserved command names
    reservedNames := []string{
        "help", "version", "config", "alias",
        "products", "users", "organizations", "apps",
        "auth", "context", "configure",
    }
    
    // Create alias manager
    manager := alias.NewManager(configProvider, reservedNames)
    
    // Create root command
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "My TykCtl Extension",
    }
    
    // Add alias command
    builder := alias.NewCommandBuilder(manager)
    rootCmd.AddCommand(builder.BuildAliasCommand())
    
    // Register aliases as subcommands
    registrar := alias.NewRegistrar(manager)
    ctx := context.Background()
    
    if err := registrar.RegisterAliasesWithValidation(ctx, rootCmd, reservedNames); err != nil {
        log.Fatalf("Failed to register aliases: %v", err)
    }
    
    // Execute command
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func createExtensionConfigProvider() alias.ConfigProvider {
    // This would integrate with your extension's configuration system
    return alias.NewExtensionConfigProvider(
        setAliasFunc,
        getAliasFunc,
        deleteAliasFunc,
        listAliasesFunc,
    )
}

// These functions would be implemented by your extension
func setAliasFunc(ctx context.Context, name, expansion string) error {
    // Implementation depends on your config system
    return nil
}

func getAliasFunc(ctx context.Context, name string) (string, bool) {
    // Implementation depends on your config system
    return "", false
}

func deleteAliasFunc(ctx context.Context, name string) error {
    // Implementation depends on your config system
    return nil
}

func listAliasesFunc(ctx context.Context) map[string]string {
    // Implementation depends on your config system
    return make(map[string]string)
}
```

## Resources

- [Configuration Guide](../config/README.md)
- [Extension Framework](../extension/README.md)
- [Getting Started Guide](../guides/getting-started.md)