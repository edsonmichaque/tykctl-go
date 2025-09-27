# Extension Package

The `extension` package provides comprehensive extension management for tykctl, including installation, execution, and lifecycle management of extensions from GitHub repositories.

## Features

- **GitHub Integration**: Install extensions directly from GitHub repositories
- **Hook System**: Integrated hook system for extension lifecycle events
- **Version Management**: Support for specific version installation and management
- **Extension Discovery**: Search and discover extensions on GitHub
- **Execution Engine**: Run installed extensions with proper context and environment
- **Configuration Management**: XDG-based configuration directory management
- **Functional Options**: Clean configuration using functional options pattern

## Usage

### Basic Extension Installation

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/extension"
)

func main() {
    // Create extension installer
    configDir := "/tmp/tykctl-config"
    installer := extension.NewInstaller(configDir)
    
    ctx := context.Background()
    
    // Install an extension
    err := installer.InstallExtension(ctx, "owner", "repo-name")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Extension installed successfully!")
}
```

### Extension Installation with Options

```go
func main() {
    // Create installer with functional options
    installer := extension.NewInstaller(
        "/tmp/tykctl-config",
        extension.WithGitHubToken("your-github-token"),
        extension.WithLogger(logger),
        extension.WithHooks(hookManager),
    )
    
    ctx := context.Background()
    
    // Install with specific version
    err := installer.InstallExtensionWithVersion(ctx, "owner", "repo", "v1.2.3")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Extension Discovery

```go
func searchExtensions() error {
    installer := extension.NewInstaller("/tmp/tykctl-config")
    ctx := context.Background()
    
    // Search for extensions
    extensions, err := installer.SearchExtensions(ctx, "tykctl", 10)
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d extensions:\n", len(extensions))
    for _, ext := range extensions {
        fmt.Printf("- %s: %s (%d stars)\n", 
            ext.Name, ext.Description, ext.Stars)
    }
    
    return nil
}
```

### Extension Execution

```go
func runExtension() error {
    runner := extension.NewRunner("/tmp/tykctl-config")
    ctx := context.Background()
    
    // List available extensions
    available, err := runner.ListAvailableExtensions()
    if err != nil {
        return err
    }
    
    fmt.Printf("Available extensions: %v\n", available)
    
    // Run an extension
    err = runner.RunExtension(ctx, "my-extension", []string{"--help"})
    if err != nil {
        return err
    }
    
    return nil
}
```

## Advanced Usage

### Extension with Hooks

```go
import (
    "github.com/edsonmichaque/tykctl-go/extension"
    "github.com/edsonmichaque/tykctl-go/hook"
)

func installWithHooks() error {
    // Create hook manager
    hookManager := hook.New(
        hook.WithExternalHookDir("/tmp/tykctl-hooks"),
        hook.WithLogger(logger),
    )
    
    // Register hooks
    ctx := context.Background()
    hookManager.RegisterBuiltin(ctx, extension.HookTypeBeforeInstall, func(ctx context.Context, data interface{}) error {
        hookData := data.(*hook.HookData)
        fmt.Printf("About to install: %s\n", hookData.ExtensionName)
        return nil
    })
    
    // Create installer with hooks
    installer := extension.NewInstaller(
        "/tmp/tykctl-config",
        extension.WithHooks(hookManager),
    )
    
    // Install extension (hooks will be executed)
    return installer.InstallExtension(ctx, "owner", "repo")
}
```

### Extension Management

```go
func manageExtensions() error {
    installer := extension.NewInstaller("/tmp/tykctl-config")
    ctx := context.Background()
    
    // List installed extensions
    installed, err := installer.ListInstalledExtensions(ctx)
    if err != nil {
        return err
    }
    
    fmt.Printf("Installed extensions (%d):\n", len(installed))
    for _, ext := range installed {
        fmt.Printf("- %s v%s (%s)\n", 
            ext.Name, ext.Version, ext.Repository)
    }
    
    // Check if extension is installed
    if installer.IsExtensionInstalled(ctx, "my-extension") {
        fmt.Println("Extension is installed")
        
        // Uninstall extension
        err = installer.UninstallExtension(ctx, "my-extension")
        if err != nil {
            return err
        }
        fmt.Println("Extension uninstalled")
    }
    
    return nil
}
```

### Extension with Custom Configuration

```go
func installWithCustomConfig() error {
    // Create installer with custom options
    installer := extension.NewInstaller(
        "/custom/config/dir",
        extension.WithGitHubToken(os.Getenv("GITHUB_TOKEN")),
        extension.WithLogger(customLogger),
        extension.WithTimeout(5*time.Minute),
    )
    
    ctx := context.Background()
    
    // Install with custom options
    err := installer.InstallExtensionWithOptions(ctx, "owner", "repo", &extension.InstallOptions{
        Version: "v1.0.0",
        Force:   true,
        Update:  false,
    })
    
    return err
}
```

## Extension Lifecycle

### Hook Types

The extension package defines several hook types for lifecycle management:

```go
const (
    HookTypeBeforeInstall   hook.HookType = "extension-before-install"
    HookTypeAfterInstall    hook.HookType = "extension-after-install"
    HookTypeBeforeUninstall hook.HookType = "extension-before-uninstall"
    HookTypeAfterUninstall  hook.HookType = "extension-after-uninstall"
    HookTypeBeforeRun       hook.HookType = "extension-before-run"
)
```

### Hook Integration

```go
func setupExtensionHooks(hookManager *hook.Manager) {
    ctx := context.Background()
    
    // Before install hook
    hookManager.RegisterBuiltin(ctx, extension.HookTypeBeforeInstall, func(ctx context.Context, data interface{}) error {
        hookData := data.(*hook.HookData)
        
        // Validate extension
        if hookData.ExtensionName == "" {
            return fmt.Errorf("extension name is required")
        }
        
        fmt.Printf("Validating extension: %s\n", hookData.ExtensionName)
        return nil
    })
    
    // After install hook
    hookManager.RegisterBuiltin(ctx, extension.HookTypeAfterInstall, func(ctx context.Context, data interface{}) error {
        hookData := data.(*hook.HookData)
        
        // Log installation
        fmt.Printf("Extension installed: %s\n", hookData.ExtensionName)
        
        // Update extension registry
        return updateExtensionRegistry(hookData.ExtensionName)
    })
}
```

## Configuration

### XDG-based Configuration

Extensions are managed using XDG-based configuration directories:

- **Config Directory**: `~/.config/tykctl/` (or `TYKCTL_CONFIG_DIR`)
- **Extensions Directory**: `~/.config/tykctl/extensions/`
- **Cache Directory**: `~/.cache/tykctl/`

### Environment Variables

- `TYKCTL_CONFIG_DIR` - Custom configuration directory
- `TYKCTL_EXTENSIONS_DIR` - Custom extensions directory
- `GITHUB_TOKEN` - GitHub token for API access
- `TYKCTL_GITHUB_TOKEN` - Tykctl-specific GitHub token

## Extension Structure

### Installed Extension

```go
type Installed struct {
    Name        string    `yaml:"name"`
    Version     string    `yaml:"version"`
    Repository  string    `yaml:"repository"`
    InstalledAt time.Time `yaml:"installed_at"`
    Path        string    `yaml:"path"`
}
```

### Extension Info

```go
type Info struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Stars       int       `json:"stargazers_count"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## Integration Examples

### With CLI Commands

```go
import (
    "github.com/spf13/cobra"
    "github.com/edsonmichaque/tykctl-go/extension"
)

func createInstallCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "install <owner> <repo>",
        Short: "Install an extension",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            owner, repo := args[0], args[1]
            
            installer := extension.NewInstaller("/tmp/tykctl-config")
            return installer.InstallExtension(cmd.Context(), owner, repo)
        },
    }
}

func createListCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "list",
        Short: "List installed extensions",
        RunE: func(cmd *cobra.Command, args []string) error {
            installer := extension.NewInstaller("/tmp/tykctl-config")
            
            installed, err := installer.ListInstalledExtensions(cmd.Context())
            if err != nil {
                return err
            }
            
            for _, ext := range installed {
                fmt.Printf("%s v%s (%s)\n", ext.Name, ext.Version, ext.Repository)
            }
            
            return nil
        },
    }
}
```

### With Progress Indicators

```go
import (
    "github.com/edsonmichaque/tykctl-go/progress"
)

func installWithProgress(owner, repo string) error {
    spinner := progress.New()
    spinner.SetMessage("Installing extension...")
    
    installer := extension.NewInstaller("/tmp/tykctl-config")
    
    // Start spinner
    spinner.Start()
    defer spinner.Stop()
    
    ctx := context.Background()
    err := installer.InstallExtension(ctx, owner, repo)
    
    if err != nil {
        spinner.SetMessage("Installation failed")
        return err
    }
    
    spinner.SetMessage("Installation completed")
    return nil
}
```

## Best Practices

- **Error Handling**: Always handle installation and execution errors gracefully
- **Hook Usage**: Use hooks for extension lifecycle management
- **Version Management**: Pin to specific versions in production
- **Configuration**: Use XDG-based configuration directories
- **Logging**: Use the provided logger for consistent logging
- **Context**: Use context for cancellation and timeout handling

## Dependencies

- `github.com/adrg/xdg` - XDG directory support
- `github.com/google/go-github/v75/github` - GitHub API client
- `go.uber.org/zap` - Logging library
- `gopkg.in/yaml.v3` - YAML configuration support