# TykCtl Plugin Package

A reusable, cross-platform plugin system for TykCtl extensions that provides comprehensive plugin management, execution, and environment setup across Linux, macOS, Windows, and Unix-like systems.

## Overview

The `plugin` package provides a unified interface for managing plugins across TykCtl extensions. It handles plugin installation, execution, discovery, and environment variable setup with extension-specific naming conventions and cross-platform compatibility.

## Features

- **Cross-Platform Support**: Works on Linux, macOS, Windows, and Unix-like systems
- **Plugin Installation**: From files, directories, or templates
- **Naming Convention**: Automatic detection of `tykctl-<extension>-*` plugins
- **Environment Setup**: Rich context and configuration access for plugins
- **Plugin Discovery**: Automatic discovery across multiple paths
- **Template Generation**: Create new plugin templates (bash/batch)
- **Wrapper Scripts**: Handle multiple executables in directories (bash/batch)
- **Executable Detection**: Platform-aware executable file detection

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "github.com/edsonmichaque/tykctl-go/plugin"
)

// Implement ConfigProvider interface
type MyConfigProvider struct {
    // Your config implementation
}

func (c *MyConfigProvider) GetConfigDir() string {
    return "/path/to/config"
}

func (c *MyConfigProvider) GetPluginDir(ctx context.Context) string {
    return "/path/to/plugins"
}

func (c *MyConfigProvider) GetPluginDiscoveryPaths(ctx context.Context) []string {
    return []string{"/path/to/plugins", "/usr/local/lib/tykctl/my-extension/plugins"}
}

// Create plugin manager
func main() {
    configProvider := &MyConfigProvider{}
    manager := plugin.NewManager("my-extension", configProvider)
    
    ctx := context.Background()
    
    // Install plugin from directory
    err := manager.InstallFromDirectory(ctx, "/source/dir", "/plugin/dir", "")
    if err != nil {
        log.Fatal(err)
    }
    
    // Execute plugin
    err = manager.Execute(ctx, "/plugin/dir/tykctl-my-extension-my-plugin", []string{"arg1", "arg2"})
    if err != nil {
        log.Fatal(err)
    }
}
```

### Plugin Installation

#### From Directory
```go
// Install from directory (detects tykctl-my-extension-* executables)
err := manager.InstallFromDirectory(ctx, "/source/dir", "/plugin/dir", "")

// Install with custom name
err := manager.InstallFromDirectory(ctx, "/source/dir", "/plugin/dir", "custom-name")
```

#### From File
```go
// Install single executable file
err := manager.InstallFromFile(ctx, "/path/to/script.sh", "/plugin/dir", "")

// Install with custom name
err := manager.InstallFromFile(ctx, "/path/to/script.sh", "/plugin/dir", "my-plugin")
```

#### Create Template
```go
// Create new plugin template
err := manager.CreateTemplate("new-plugin", "/plugin/dir")
```

### Plugin Execution

```go
// Execute plugin with arguments
err := manager.Execute(ctx, "/plugin/dir/tykctl-my-extension-my-plugin", []string{"arg1", "arg2"})
```

### Plugin Discovery

```go
// Discover all plugins
plugins, err := manager.DiscoverPlugins(ctx)
for _, plugin := range plugins {
    fmt.Printf("Plugin: %s at %s\n", plugin.Name, plugin.Path)
}
```

### Plugin Removal

```go
// Remove plugin
err := manager.RemovePlugin(ctx, "my-plugin", "/plugin/dir")
```

## Cross-Platform Support

The plugin system automatically adapts to different operating systems:

### **Executable Detection:**
- **Linux/macOS/Unix**: Checks file permissions (`mode&0111 != 0`)
- **Windows**: Checks file permissions and common executable extensions (`.exe`, `.bat`, `.cmd`, `.ps1`)

### **Script Generation:**
- **Unix-like systems**: Generates bash scripts with shebang (`#!/bin/bash`)
- **Windows**: Generates batch scripts (`.bat`) with Windows batch syntax

### **File Extensions:**
- **Windows**: Automatically adds `.exe` extension for executables, `.bat` for scripts
- **Unix-like**: No extensions added (relies on shebang and permissions)

### **Path Handling:**
- Uses `filepath.Join()` for cross-platform path construction
- Handles Windows backslashes and Unix forward slashes automatically

## Plugin Naming Convention

Plugins must follow the naming convention: `tykctl-<extension>-<name>`

Examples:
- `tykctl-portal-deploy`
- `tykctl-dashboard-backup`
- `tykctl-gateway-monitor`

## Environment Variables

Plugins receive comprehensive environment variables:

### Plugin Identification
- `TYKCTL_PLUGIN_NAME`: Plugin name (without prefix)
- `TYKCTL_PLUGIN_PATH`: Full path to plugin executable
- `TYKCTL_PLUGIN_EXTENSION`: Extension name

### Extension-Specific Directories
- `TYKCTL_<EXTENSION>_CONFIG_DIR`: Extension configuration directory
- `TYKCTL_<EXTENSION>_PLUGIN_DIR`: Extension plugin directory
- `TYKCTL_<EXTENSION>_GLOBAL_CONFIG_DIR`: Global TYKCTL configuration directory

### API Configuration
- `TYK_<EXTENSION>_URL`: Extension API URL (if set)
- `TYK_<EXTENSION>_TOKEN`: Extension API token (if set)

### Context and Debug
- `TYKCTL_<EXTENSION>_CONTEXT`: Current context (if set)
- `TYKCTL_<EXTENSION>_DEBUG`: Debug flag (if set)
- `TYKCTL_<EXTENSION>_VERBOSE`: Verbose flag (if set)

### Plugin Discovery
- `TYKCTL_<EXTENSION>_PLUGIN_DISCOVERY_PATHS`: Colon-separated discovery paths

## Directory Installation Logic

When installing from a directory:

1. **Plugin-Named Executables**: If `tykctl-<extension>-*` files are found, install them directly
2. **Single Executable**: If one regular executable is found, copy it directly
3. **Multiple Executables**: If multiple regular executables are found, create a wrapper script
4. **Mixed Directory**: Plugin-named executables take priority over regular ones

## Wrapper Scripts

For directories with multiple executables, a bash wrapper script is generated:

```bash
#!/bin/bash
# Plugin wrapper
set -euo pipefail

PLUGIN_DIR="/source/dir"
COMMAND="${1:-help}"

case "$COMMAND" in
    "help")
        echo "Available commands:"
        echo "  - deploy"
        echo "  - test"
        ;;
    "deploy")
        exec "$PLUGIN_DIR/deploy" "$@"
        ;;
    "test")
        exec "$PLUGIN_DIR/test" "$@"
        ;;
    *)
        echo "Unknown command: $COMMAND"
        exit 1
        ;;
esac
```

## Plugin Templates

Generated plugin templates include:

```bash
#!/bin/bash
# Plugin Template
set -euo pipefail

PLUGIN_NAME="my-plugin"
PLUGIN_VERSION="1.0.0"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$PLUGIN_NAME] $*" >&2
}

case "${1:-}" in
    "version")
        echo "$PLUGIN_VERSION"
        ;;
    "info")
        echo "Name: $PLUGIN_NAME"
        echo "Version: $PLUGIN_VERSION"
        echo "Description: My-plugin plugin for tykctl-my-extension"
        ;;
    *)
        echo "Usage: $0 {version|info}"
        exit 1
        ;;
esac
```

## Integration with Extensions

### Cobra Command Integration

```go
func NewPluginCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "plugin",
        Short: "Manage plugins",
    }
    
    // Create plugin manager
    manager := plugin.NewManager("my-extension", configProvider)
    
    cmd.AddCommand(&cobra.Command{
        Use:   "install <plugin-name-or-path>",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            pluginNameOrPath := args[0]
            
            // Check if it's a file or directory
            if stat, err := os.Stat(pluginNameOrPath); err == nil {
                if stat.IsDir() {
                    return manager.InstallFromDirectory(ctx, pluginNameOrPath, pluginDir, "")
                } else {
                    return manager.InstallFromFile(ctx, pluginNameOrPath, pluginDir, "")
                }
            }
            
            // Create template
            return manager.CreateTemplate(pluginNameOrPath, pluginDir)
        },
    })
    
    return cmd
}
```

### Direct Plugin Execution

```go
func RegisterPlugins(ctx context.Context, rootCmd *cobra.Command) error {
    manager := plugin.NewManager("my-extension", configProvider)
    
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        return err
    }
    
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

## Benefits

1. **Reusable**: Single package for all TykCtl extensions
2. **Consistent**: Unified plugin management across extensions
3. **Rich Context**: Comprehensive environment variable setup
4. **Flexible**: Supports multiple installation methods
5. **Developer Friendly**: Easy integration with existing extensions
6. **Production Ready**: Robust error handling and validation

## Extension-Specific Environment Variables

The package automatically generates extension-specific environment variables:

- Extension name is converted to uppercase
- Variables follow pattern: `TYKCTL_<EXTENSION>_*`
- API variables follow pattern: `TYK_<EXTENSION>_*`

Examples for "portal" extension:
- `TYKCTL_PORTAL_CONFIG_DIR`
- `TYKCTL_PORTAL_PLUGIN_DIR`
- `TYK_PORTAL_URL`
- `TYK_PORTAL_TOKEN`
- `TYKCTL_PORTAL_CONTEXT`
- `TYKCTL_PORTAL_DEBUG`