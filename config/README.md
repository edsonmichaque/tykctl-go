# TykCtl Configuration Package

A production-ready configuration package for TykCtl extensions with enhanced features, multi-context support, validation, and comprehensive resource management.

## Features

- **Multi-Context Support**: Manage multiple configuration contexts (dev, staging, prod)
- **Enhanced Discovery**: Rich metadata and filtering for hooks, plugins, templates, and cache configs
- **Validation Framework**: Built-in validation with custom validators
- **Caching**: In-memory caching with TTL support
- **Structured Logging**: Comprehensive logging with different levels
- **Metrics**: Built-in metrics collection (no-op implementation)
- **Context Management**: Easy context switching and isolation
- **Resource Discovery**: Automatic discovery of hooks, plugins, templates, and cache configurations

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/edsonmichaque/tykctl/examples/extensions/tykctl-ai-studiov2/pkg/config"
    "github.com/edsonmichaque/tykctl/examples/extensions/tykctl-ai-studiov2/pkg/config/types"
)

type MyConfig struct {
    URL        string        `mapstructure:"url" validate:"required,url"`
    Token      string        `mapstructure:"token" validate:"required,min_length=10"`
    MaxTokens  int           `mapstructure:"max_tokens" validate:"range=1,10000"`
    Temperature float64     `mapstructure:"temperature" validate:"range=0,2"`
    CacheEnabled bool       `mapstructure:"cache_enabled"`
    CacheTTL    time.Duration `mapstructure:"cache_ttl"`
}

// Implement DefaultSetter interface
func (c *MyConfig) SetDefaults() {
    if c.MaxTokens == 0 {
        c.MaxTokens = 512
    }
    if c.Temperature == 0 {
        c.Temperature = 0.7
    }
    if c.CacheTTL == 0 {
        c.CacheTTL = time.Hour
    }
}

// Implement Validator interface
func (c *MyConfig) Validate() error {
    if c.URL == "" {
        return fmt.Errorf("url is required")
    }
    if c.Token == "" {
        return fmt.Errorf("token is required")
    }
    return nil
}

func main() {
    ctx := context.Background()
    
    // Create loader
    loader, err := config.NewLoader(ctx, config.LoaderOptions{
        Extension:      "my-app",
        CacheEnabled:   true,
        CacheTTL:       5 * time.Minute,
        LogLevel:       types.LogLevelInfo,
        MetricsEnabled: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer loader.Close()
    
    // Load configuration
    var cfg MyConfig
    if err := loader.Load(ctx, &cfg); err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("URL: %s\n", cfg.URL)
    fmt.Printf("Max Tokens: %d\n", cfg.MaxTokens)
}
```

### Multi-Context Usage

```go
func main() {
    ctx := context.Background()
    
    // Create context store
    ctxStore, err := config.NewContextStore(config.ContextStoreOptions{
        ConfigPath: "~/.tykctl/contexts",
        Logger:     logger,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Create development context
    devConfig := map[string]interface{}{
        "url":         "https://dev-api.example.com",
        "max_tokens":  256,
        "temperature": 0.8,
    }
    
    if err := ctxStore.CreateContext("dev", devConfig); err != nil {
        log.Printf("Failed to create dev context: %v", err)
    }
    
    // Create production context
    prodConfig := map[string]interface{}{
        "url":         "https://api.example.com",
        "max_tokens":  1024,
        "temperature": 0.3,
    }
    
    if err := ctxStore.CreateContext("prod", prodConfig); err != nil {
        log.Printf("Failed to create prod context: %v", err)
    }
    
    // Switch to development context
    if err := ctxStore.SwitchContext("dev"); err != nil {
        log.Printf("Failed to switch to dev context: %v", err)
    }
    
    // Load configuration for specific context
    loader, err := config.NewLoader(ctx, config.LoaderOptions{
        Extension: "my-app",
        Context:   "dev",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    var cfg MyConfig
    if err := loader.Load(ctx, &cfg); err != nil {
        log.Fatal(err)
    }
    
    // Configuration will be loaded from dev context
    fmt.Printf("URL: %s\n", cfg.URL) // https://dev-api.example.com
}
```

### Resource Discovery

```go
func main() {
    ctx := context.Background()
    
    loader, err := config.NewLoader(ctx, config.LoaderOptions{
        Extension: "my-app",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer loader.Close()
    
    // Discover hooks with filtering
    hooks, err := loader.DiscoverHooks(ctx, types.HookFilter{
        Events:     []string{"pre-command", "post-command"},
        Enabled:    &[]bool{true}[0],
        MinTimeout: 5 * time.Second,
        MaxTimeout: 60 * time.Second,
    })
    if err != nil {
        log.Printf("Failed to discover hooks: %v", err)
    } else {
        fmt.Printf("Found hooks:\n")
        for event, eventHooks := range hooks {
            fmt.Printf("  %s (%d hooks):\n", event, len(eventHooks))
            for _, hook := range eventHooks {
                fmt.Printf("    - %s: %s (priority: %d, timeout: %v)\n", 
                    hook.Name, hook.Path, hook.Priority, hook.Timeout)
            }
        }
    }
    
    // Discover plugins with filtering
    plugins, err := loader.DiscoverPlugins(ctx, types.PluginFilter{
        Commands: []string{"train", "evaluate"},
        Enabled:  &[]bool{true}[0],
    })
    if err != nil {
        log.Printf("Failed to discover plugins: %v", err)
    } else {
        fmt.Printf("Found %d plugins:\n", len(plugins))
        for _, plugin := range plugins {
            fmt.Printf("  - %s: %s (version: %s, commands: %v)\n", 
                plugin.Name, plugin.Path, plugin.Version, plugin.Commands)
        }
    }
    
    // Discover templates with filtering
    templates, err := loader.DiscoverTemplates(ctx, types.TemplateFilter{
        Types:   []string{"output", "input"},
        Formats: []string{"json", "yaml"},
        Enabled: &[]bool{true}[0],
    })
    if err != nil {
        log.Printf("Failed to discover templates: %v", err)
    } else {
        fmt.Printf("Found %d templates:\n", len(templates))
        for _, template := range templates {
            fmt.Printf("  - %s: %s (type: %s, format: %s)\n", 
                template.Name, template.Path, template.Type, template.Format)
        }
    }
}
```

## Configuration Files

### YAML Configuration

```yaml
# my-app.yaml
url: "https://api.example.com"
token: "your-token-here"
max_tokens: 512
temperature: 0.7
cache_enabled: true
cache_ttl: "1h"
```

### Environment Variables

```bash
# Environment variables automatically override config file values
export TYKCTL_MY_APP_URL="https://api.example.com"
export TYKCTL_MY_APP_MAX_TOKENS="1024"
export TYKCTL_MY_APP_TEMPERATURE="0.5"
export TYKCTL_MY_APP_CACHE_ENABLED="true"
export TYKCTL_MY_APP_CACHE_TTL="2h"
```

### Context Configuration

```yaml
# ~/.tykctl/contexts/dev.yaml
name: dev
description: Development environment
config:
  url: "https://dev-api.example.com"
  token: "dev-token-here"
  max_tokens: 256
  temperature: 0.8
  cache_enabled: true
  cache_ttl: "30m"
metadata:
  environment: development
  team: ai-team
  region: us-west-2
resources:
  hooks:
    pre-command:
      - name: "setup-dev-env"
        path: "/usr/local/bin/setup-dev-env"
        priority: 100
  plugins:
    - name: "dev-tools"
      path: "/usr/local/bin/tykctl-my-app-dev-tools"
      enabled: true
```

## Configuration Precedence

1. Environment variables (highest priority)
2. Context-specific configuration files
3. Extension-specific configuration files
4. Global configuration files
5. Default values (lowest priority)

## Resource Discovery Paths

### Hooks
- `/etc/tykctl/<extension>/hooks`
- `~/.local/share/tykctl/<extension>/hooks`
- `~/.tykctl/<extension>/hooks`
- `./.tykctl/<extension>/hooks`

### Plugins
- `/usr/local/bin`
- `/usr/bin`
- `~/.local/share/tykctl/<extension>/bin`
- `~/.tykctl/<extension>/bin`
- `./.tykctl/<extension>/bin`

### Templates
- `/etc/tykctl/<extension>/templates`
- `~/.local/share/tykctl/<extension>/templates`
- `~/.tykctl/<extension>/templates`
- `./.tykctl/<extension>/templates`

### Cache Configurations
- `/etc/tykctl/<extension>/cache`
- `~/.local/share/tykctl/<extension>/cache`
- `~/.tykctl/<extension>/cache`
- `./.tykctl/<extension>/cache`

## Validation

The package supports struct tag validation:

```go
type MyConfig struct {
    URL        string  `mapstructure:"url" validate:"required,url"`
    Token      string  `mapstructure:"token" validate:"required,min_length=10"`
    MaxTokens  int     `mapstructure:"max_tokens" validate:"range=1,10000"`
    Temperature float64 `mapstructure:"temperature" validate:"range=0,2"`
}
```

### Built-in Validators

- `required`: Field is required
- `min_length=N`: Minimum string length
- `max_length=N`: Maximum string length
- `range=min,max`: Numeric range validation
- `url`: URL format validation

## Interfaces

### DefaultSetter
```go
type DefaultSetter interface {
    SetDefaults()
}
```

### Validator
```go
type Validator interface {
    Validate() error
}
```

### Configurable
```go
type Configurable interface {
    Configure(cfg Config) error
}
```

## Dependencies

- `github.com/adrg/xdg` - XDG Base Directory Specification support
- `github.com/spf13/viper` - Configuration management
- `github.com/spf13/cobra` - CLI framework (for context commands)

## License

This package is part of the TykCtl project and follows the same license terms.