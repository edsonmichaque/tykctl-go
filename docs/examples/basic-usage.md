# Basic Usage Examples

This document provides practical examples of using TykCtl-Go in your applications.

## Quick Start Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/config"
    "github.com/edsonmichaque/tykctl-go/plugin"
)

func main() {
    ctx := context.Background()
    
    // Initialize configuration
    cfg, err := config.NewConfigManager()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create plugin manager
    manager := plugin.NewManager("my-extension", cfg)
    
    // Discover plugins
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d plugins\n", len(plugins))
    
    // List plugins
    for _, p := range plugins {
        fmt.Printf("- %s (%s)\n", p.Name, p.Path)
    }
}
```

## Configuration Management

### Loading Configuration

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/config"
)

func main() {
    // Load configuration
    cfg, err := config.NewConfigManager()
    if err != nil {
        log.Fatal(err)
    }
    
    // Access configuration values
    debug := cfg.GetBool("debug")
    timeout := cfg.GetString("client.timeout")
    output := cfg.GetString("output")
    
    fmt.Printf("Debug: %v\n", debug)
    fmt.Printf("Timeout: %s\n", timeout)
    fmt.Printf("Output: %s\n", output)
}
```

### Environment Variable Configuration

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/edsonmichaque/tykctl-go/config"
)

func main() {
    // Set environment variables
    os.Setenv("TYKCTL_DEBUG", "true")
    os.Setenv("TYKCTL_CLIENT_TIMEOUT", "30")
    os.Setenv("TYK_PORTAL_URL", "https://portal.example.com")
    
    // Load configuration
    cfg, err := config.NewConfigManager()
    if err != nil {
        panic(err)
    }
    
    // Configuration will include environment variables
    fmt.Printf("Debug: %v\n", cfg.GetBool("debug"))
    fmt.Printf("Portal URL: %s\n", os.Getenv("TYK_PORTAL_URL"))
}
```

## Plugin Management

### Installing Plugins

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/config"
    "github.com/edsonmichaque/tykctl-go/plugin"
)

func main() {
    ctx := context.Background()
    
    // Create plugin manager
    cfg, _ := config.NewConfigManager()
    manager := plugin.NewManager("portal", cfg)
    
    // Install plugin from file
    err := manager.InstallFromFile(ctx, "/path/to/plugin.sh", "/plugin/dir", "my-plugin")
    if err != nil {
        log.Fatal(err)
    }
    
    // Install plugin from directory
    err = manager.InstallFromDirectory(ctx, "/path/to/plugin/dir", "/plugin/dir", "")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create plugin template
    err = manager.CreateTemplate("new-plugin", "/plugin/dir")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Plugins installed successfully!")
}
```

### Executing Plugins

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/config"
    "github.com/edsonmichaque/tykctl-go/plugin"
)

func main() {
    ctx := context.Background()
    
    // Create plugin manager
    cfg, _ := config.NewConfigManager()
    manager := plugin.NewManager("portal", cfg)
    
    // Discover plugins
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    if len(plugins) == 0 {
        fmt.Println("No plugins found")
        return
    }
    
    // Execute plugin with default timeout
    err = manager.Execute(ctx, plugins[0].Path, []string{"--help"})
    if err != nil {
        log.Fatal(err)
    }
    
    // Execute plugin with custom timeout
    err = manager.ExecuteWithTimeout(ctx, plugins[0].Path, []string{"arg1", "arg2"}, 30*time.Second)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Plugin executed successfully!")
}
```

## Alias Management

### Creating and Managing Aliases

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/alias"
)

func main() {
    ctx := context.Background()
    
    // Create alias manager
    provider := alias.NewInMemoryConfigProvider()
    manager := alias.NewManager(provider, []string{"help", "version"})
    
    // Create simple aliases
    err := manager.SetAlias(ctx, "co", "products checkout")
    if err != nil {
        log.Fatal(err)
    }
    
    err = manager.SetAlias(ctx, "users", "users list")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create parameterized aliases
    err = manager.SetAlias(ctx, "getuser", "users get $1")
    if err != nil {
        log.Fatal(err)
    }
    
    err = manager.SetAlias(ctx, "deploy", "products deploy $1 --env $2")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create shell aliases
    err = manager.SetAlias(ctx, "cleanup", "!rm -rf /tmp/tykctl-*")
    if err != nil {
        log.Fatal(err)
    }
    
    err = manager.SetAlias(ctx, "logs", "!tail -f /var/log/tykctl.log")
    if err != nil {
        log.Fatal(err)
    }
    
    // List all aliases
    aliases := manager.ListAliases(ctx)
    fmt.Printf("Created %d aliases:\n", len(aliases))
    for name, expansion := range aliases {
        aliasType := manager.GetAliasType(expansion)
        fmt.Printf("- %s (%s): %s\n", name, aliasType, expansion)
    }
}
```

### Executing Aliases

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/alias"
)

func main() {
    ctx := context.Background()
    
    // Create alias manager
    provider := alias.NewInMemoryConfigProvider()
    manager := alias.NewManager(provider, []string{"help", "version"})
    
    // Set up aliases
    manager.SetAlias(ctx, "deploy", "products deploy $1 --env $2")
    manager.SetAlias(ctx, "getuser", "users get $1")
    manager.SetAlias(ctx, "cleanup", "!rm -rf /tmp/tykctl-*")
    
    // Preview alias expansion
    preview := manager.ExpandAliasPreview("deploy", []string{"my-api", "production"})
    fmt.Printf("Deploy alias would execute: %s\n", preview)
    
    preview = manager.ExpandAliasPreview("getuser", []string{"john@example.com"})
    fmt.Printf("Getuser alias would execute: %s\n", preview)
    
    // Execute aliases (in a real application, these would call the actual commands)
    fmt.Println("Aliases are ready for execution!")
}
```

### Cobra Integration with Aliases

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/spf13/cobra"
    "github.com/edsonmichaque/tykctl-go/alias"
)

func main() {
    ctx := context.Background()
    
    // Create root command
    rootCmd := &cobra.Command{
        Use:   "my-extension",
        Short: "My TykCtl Extension with Aliases",
        Long:  "A comprehensive extension with alias support",
    }
    
    // Create alias manager
    provider := alias.NewInMemoryConfigProvider()
    manager := alias.NewManager(provider, []string{"help", "version", "my-extension"})
    
    // Add alias command
    builder := alias.NewCommandBuilder(manager)
    rootCmd.AddCommand(builder.BuildAliasCommand())
    
    // Register aliases as subcommands
    registrar := alias.NewRegistrar(manager)
    err := registrar.RegisterAliasesWithValidation(ctx, rootCmd, []string{"help", "version"})
    if err != nil {
        log.Fatal(err)
    }
    
    // Set up some default aliases
    manager.SetAlias(ctx, "co", "products checkout")
    manager.SetAlias(ctx, "deploy", "products deploy $1 --env $2")
    manager.SetAlias(ctx, "users", "users list")
    
    // Execute command
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}
```

## Extension Development

### Basic Extension Structure

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/spf13/cobra"
    "github.com/edsonmichaque/tykctl-go/config"
    "github.com/edsonmichaque/tykctl-go/plugin"
)

func main() {
    ctx := context.Background()
    
    // Create root command
    rootCmd := &cobra.Command{
        Use:   "my-extension",
        Short: "My TykCtl Extension",
        Long:  "A comprehensive extension for managing my resources",
    }
    
    // Add global flags
    rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging")
    rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
    
    // Add commands
    rootCmd.AddCommand(NewResourceCommand())
    rootCmd.AddCommand(NewPluginCommand())
    
    // Register plugins as subcommands
    RegisterPlugins(ctx, rootCmd)
    
    // Execute
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func NewResourceCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "resource",
        Short: "Manage resources",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Println("Resource management")
            return nil
        },
    }
}

func NewPluginCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "plugin",
        Short: "Plugin management",
    }
    
    cmd.AddCommand(&cobra.Command{
        Use:   "list",
        Short: "List plugins",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            cfg, _ := config.NewConfigManager()
            manager := plugin.NewManager("my-extension", cfg)
            
            plugins, err := manager.DiscoverPlugins(ctx)
            if err != nil {
                return err
            }
            
            fmt.Printf("Found %d plugins:\n", len(plugins))
            for _, p := range plugins {
                fmt.Printf("- %s (%s)\n", p.Name, p.Path)
            }
            
            return nil
        },
    })
    
    return cmd
}

func RegisterPlugins(ctx context.Context, rootCmd *cobra.Command) error {
    cfg, _ := config.NewConfigManager()
    manager := plugin.NewManager("my-extension", cfg)
    
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        return err
    }
    
    // Get existing command names to avoid conflicts
    existingCommands := make(map[string]bool)
    for _, cmd := range rootCmd.Commands() {
        existingCommands[cmd.Name()] = true
        for _, alias := range cmd.Aliases {
            existingCommands[alias] = true
        }
    }
    
    // Register plugins as subcommands
    for _, plugin := range plugins {
        if existingCommands[plugin.Name] {
            continue // Skip if there's a command conflict
        }
        
        pluginCmd := &cobra.Command{
            Use:   plugin.Name,
            Short: fmt.Sprintf("Execute plugin: %s", plugin.Name),
            RunE: func(cmd *cobra.Command, args []string) error {
                return manager.Execute(ctx, plugin.Path, args)
            },
        }
        
        rootCmd.AddCommand(pluginCmd)
    }
    
    return nil
}
```

## Error Handling

### Proper Error Handling

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/config"
    "github.com/edsonmichaque/tykctl-go/plugin"
)

func main() {
    ctx := context.Background()
    
    // Initialize with proper error handling
    cfg, err := config.NewConfigManager()
    if err != nil {
        log.Fatalf("Failed to initialize configuration: %v", err)
    }
    
    manager := plugin.NewManager("portal", cfg)
    
    // Discover plugins with error handling
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        log.Fatalf("Failed to discover plugins: %v", err)
    }
    
    if len(plugins) == 0 {
        fmt.Println("No plugins found")
        return
    }
    
    // Execute plugin with error handling
    err = manager.Execute(ctx, plugins[0].Path, []string{"--help"})
    if err != nil {
        // Check for specific error types
        if errors.Is(err, plugin.ErrPluginNotFound) {
            log.Fatal("Plugin not found")
        }
        if errors.Is(err, plugin.ErrPluginTimeout) {
            log.Fatal("Plugin execution timed out")
        }
        log.Fatalf("Plugin execution failed: %v", err)
    }
    
    fmt.Println("Plugin executed successfully!")
}
```

## Testing Examples

### Unit Testing

```go
package main

import (
    "context"
    "os"
    "testing"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/config"
    "github.com/edsonmichaque/tykctl-go/plugin"
    "github.com/stretchr/testify/assert"
)

func TestPluginExecution(t *testing.T) {
    // Create test plugin
    testPlugin := createTestPlugin(t)
    defer os.Remove(testPlugin)
    
    // Set up test environment
    os.Setenv("TYK_PORTAL_URL", "https://test.example.com")
    os.Setenv("TYK_PORTAL_TOKEN", "test-token")
    defer os.Unsetenv("TYK_PORTAL_URL")
    defer os.Unsetenv("TYK_PORTAL_TOKEN")
    
    // Create plugin manager
    cfg, err := config.NewConfigManager()
    assert.NoError(t, err)
    
    manager := plugin.NewManager("portal", cfg)
    
    // Test plugin execution
    err = manager.Execute(context.Background(), testPlugin, []string{"--help"})
    assert.NoError(t, err)
}

func TestPluginTimeout(t *testing.T) {
    // Create long-running test plugin
    testPlugin := createLongRunningPlugin(t)
    defer os.Remove(testPlugin)
    
    cfg, _ := config.NewConfigManager()
    manager := plugin.NewManager("portal", cfg)
    
    // Test timeout
    err := manager.ExecuteWithTimeout(context.Background(), testPlugin, []string{}, 1*time.Second)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "timed out")
}

func createTestPlugin(t *testing.T) string {
    // Create a simple test plugin
    pluginContent := `#!/bin/bash
echo "Test plugin executed"
echo "Arguments: $@"
`
    
    tmpFile, err := os.CreateTemp("", "test-plugin-*")
    assert.NoError(t, err)
    
    _, err = tmpFile.WriteString(pluginContent)
    assert.NoError(t, err)
    
    err = tmpFile.Close()
    assert.NoError(t, err)
    
    err = os.Chmod(tmpFile.Name(), 0755)
    assert.NoError(t, err)
    
    return tmpFile.Name()
}

func createLongRunningPlugin(t *testing.T) string {
    // Create a plugin that sleeps for 10 seconds
    pluginContent := `#!/bin/bash
sleep 10
echo "Long running plugin completed"
`
    
    tmpFile, err := os.CreateTemp("", "long-plugin-*")
    assert.NoError(t, err)
    
    _, err = tmpFile.WriteString(pluginContent)
    assert.NoError(t, err)
    
    err = tmpFile.Close()
    assert.NoError(t, err)
    
    err = os.Chmod(tmpFile.Name(), 0755)
    assert.NoError(t, err)
    
    return tmpFile.Name()
}
```

## Configuration Examples

### Custom Configuration

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/config"
)

func main() {
    // Create custom configuration
    customConfig := map[string]interface{}{
        "debug":   true,
        "verbose": true,
        "client": map[string]interface{}{
            "timeout": 30,
            "retries": 3,
        },
        "plugins": map[string]interface{}{
            "execution": map[string]interface{}{
                "timeout": "5m",
            },
        },
    }
    
    // Load configuration with custom values
    cfg, err := config.NewConfigManagerWithDefaults(customConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    // Access configuration
    fmt.Printf("Debug: %v\n", cfg.GetBool("debug"))
    fmt.Printf("Client timeout: %v\n", cfg.GetInt("client.timeout"))
    fmt.Printf("Plugin timeout: %s\n", cfg.GetString("plugins.execution.timeout"))
}
```

## Advanced Usage

### Plugin Management with Custom Configuration

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/config"
    "github.com/edsonmichaque/tykctl-go/plugin"
)

func main() {
    ctx := context.Background()
    
    // Set custom plugin timeout
    os.Setenv("TYKCTL_PORTAL_PLUGIN_TIMEOUT", "10m")
    
    cfg, err := config.NewConfigManager()
    if err != nil {
        log.Fatal(err)
    }
    
    manager := plugin.NewManager("portal", cfg)
    
    // Get configured timeout
    timeout := manager.GetConfiguredTimeout()
    fmt.Printf("Plugin timeout: %v\n", timeout)
    
    // Discover plugins
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Execute plugins with configured timeout
    for _, p := range plugins {
        fmt.Printf("Executing plugin: %s\n", p.Name)
        
        err = manager.Execute(ctx, p.Path, []string{"--version"})
        if err != nil {
            log.Printf("Plugin %s failed: %v", p.Name, err)
            continue
        }
        
        fmt.Printf("Plugin %s completed successfully\n", p.Name)
    }
}
```

## Resources

- [Getting Started Guide](../guides/getting-started.md)
- [Plugin Development Guide](../guides/plugin-development.md)
- [Alias System Documentation](../alias/README.md)
- [Configuration Documentation](../config/README.md)
- [Plugin System Documentation](../plugin/README.md)