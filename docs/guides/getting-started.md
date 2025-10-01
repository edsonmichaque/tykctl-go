# Getting Started Guide

This guide will help you get up and running with TykCtl-Go quickly.

## 🚀 Quick Start

### Installation

```bash
go get github.com/edsonmichaque/tykctl-go
```

### Basic Usage

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
    
    // Create plugin manager for your extension
    manager := plugin.NewManager("my-extension", cfg)
    
    // Discover available plugins
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d plugins\n", len(plugins))
    
    // Execute a plugin if available
    if len(plugins) > 0 {
        err = manager.Execute(ctx, plugins[0].Path, []string{"--help"})
        if err != nil {
            log.Fatal(err)
        }
    }
}
```

## 🔧 Configuration

### Environment Variables

Set up your environment:

```bash
# Global configuration
export TYKCTL_DEBUG=true
export TYKCTL_VERBOSE=true

# Plugin configuration
export TYKCTL_PLUGIN_TIMEOUT="5m"

# Extension-specific configuration
export TYKCTL_MY_EXTENSION_PLUGIN_DIR="/path/to/plugins"
export TYKCTL_MY_EXTENSION_PLUGIN_TIMEOUT="10m"
```

### Configuration File

Create a `config.yaml` file:

```yaml
debug: true
verbose: true

plugins:
  execution:
    timeout: "5m"
  
  discovery:
    system_paths:
      - "/usr/local/lib/tykctl/my-extension/plugins"
    user_paths:
      - "~/.local/share/tykctl/my-extension/plugins"
```

## 🔌 Plugin System

### Creating a Plugin

Create a simple plugin:

```bash
#!/bin/bash
# my-plugin.sh

echo "Hello from my plugin!"
echo "Arguments: $@"
echo "Environment:"
env | grep TYKCTL
```

Make it executable:

```bash
chmod +x my-plugin.sh
```

### Installing a Plugin

```go
// Install plugin from file
err := manager.InstallFromFile(ctx, "/path/to/my-plugin.sh", "/plugin/dir", "my-plugin")

// Install plugin from directory
err := manager.InstallFromDirectory(ctx, "/path/to/plugin/dir", "/plugin/dir", "")

// Create plugin template
err := manager.CreateTemplate("new-plugin", "/plugin/dir")
```

### Executing Plugins

```go
// Execute with default timeout
err := manager.Execute(ctx, "/plugin/dir/tykctl-my-extension-my-plugin", []string{"arg1", "arg2"})

// Execute with custom timeout
err := manager.ExecuteWithTimeout(ctx, "/plugin/dir/tykctl-my-extension-my-plugin", []string{"arg1"}, 30*time.Second)
```

## 🏗️ Extension Development

### Basic Extension Structure

```go
package main

import (
    "context"
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
    }
    
    // Add plugin commands
    rootCmd.AddCommand(NewPluginCommand())
    
    // Register plugins as subcommands
    RegisterPlugins(ctx, rootCmd)
    
    // Execute
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func NewPluginCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "plugin",
        Short: "Plugin management",
        // ... plugin subcommands
    }
}

func RegisterPlugins(ctx context.Context, rootCmd *cobra.Command) error {
    cfg, _ := config.NewConfigManager()
    manager := plugin.NewManager("my-extension", cfg)
    
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        return err
    }
    
    // Register plugins as subcommands
    for _, plugin := range plugins {
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

## 📝 Common Patterns

### Error Handling

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to execute plugin %s: %w", pluginName, err)
}

// Check for specific errors
if errors.Is(err, plugin.ErrPluginNotFound) {
    // Handle plugin not found
}
```

### Configuration Access

```go
// Get configuration
cfg, err := config.NewConfigManager()
if err != nil {
    return err
}

// Access specific values
debug := cfg.GetBool("debug")
timeout := cfg.GetString("plugins.execution.timeout")
```

### Plugin Discovery

```go
// Discover plugins
plugins, err := manager.DiscoverPlugins(ctx)
if err != nil {
    return err
}

// Filter plugins
var filteredPlugins []plugin.Plugin
for _, p := range plugins {
    if strings.Contains(p.Name, "deploy") {
        filteredPlugins = append(filteredPlugins, p)
    }
}
```

## 🧪 Testing

### Unit Testing

```go
func TestPluginExecution(t *testing.T) {
    // Create test plugin
    testPlugin := createTestPlugin(t)
    defer os.Remove(testPlugin)
    
    // Test execution
    manager := plugin.NewManager("test", mockConfigProvider{})
    err := manager.Execute(context.Background(), testPlugin, []string{"test"})
    
    assert.NoError(t, err)
}
```

### Integration Testing

```go
func TestPluginDiscovery(t *testing.T) {
    // Set up test environment
    testDir := createTestPluginDir(t)
    defer os.RemoveAll(testDir)
    
    // Test discovery
    manager := plugin.NewManager("test", testConfigProvider{testDir})
    plugins, err := manager.DiscoverPlugins(context.Background())
    
    assert.NoError(t, err)
    assert.Len(t, plugins, 1)
}
```

## 🚀 Next Steps

1. **Explore the API**: Check out the [API Documentation](../api/README.md)
2. **Plugin Development**: Learn more in the [Plugin Guide](../plugin/README.md)
3. **Configuration**: Deep dive into [Configuration Management](../config/README.md)
4. **Extension Framework**: Build extensions with the [Extension Guide](../extension/README.md)

## 🆘 Troubleshooting

### Common Issues

**Plugin not found:**
- Check plugin naming convention (`tykctl-<extension>-<name>`)
- Verify plugin is executable
- Check plugin discovery paths

**Configuration not loading:**
- Verify environment variables
- Check config file syntax
- Ensure proper file permissions

**Timeout issues:**
- Check timeout configuration
- Verify plugin execution time
- Review timeout error messages

### Getting Help

- **Documentation**: Browse the docs directory
- **Issues**: [GitHub Issues](https://github.com/edsonmichaque/tykctl-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/edsonmichaque/tykctl-go/discussions)