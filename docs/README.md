# TykCtl-Go Documentation

Welcome to the TykCtl-Go documentation. This directory contains comprehensive documentation for the TykCtl Go framework and its components.

## 📚 Documentation Structure

### Core Components

- **[API Documentation](api/)** - API client utilities and helpers
- **[Configuration Management](config/)** - Configuration system, environment variables, and discovery
- **[Plugin System](plugin/)** - Cross-platform plugin management and execution
- **[Extension Framework](extension/)** - Extension development and management

### Features

- **[Hooks System](hooks/)** - Event-driven hooks and automation
- **[Templates](templates/)** - Template system for resource generation
- **[Progress Tracking](progress/)** - Progress indicators and status tracking

### Guides and Examples

- **[User Guides](guides/)** - Step-by-step guides for common tasks
- **[Examples](examples/)** - Code examples and sample configurations
- **[Development Guide](development.md)** - Contributing and development guidelines

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
    "github.com/edsonmichaque/tykctl-go/config"
    "github.com/edsonmichaque/tykctl-go/plugin"
)

func main() {
    ctx := context.Background()
    
    // Initialize configuration
    cfg, err := config.NewConfigManager()
    if err != nil {
        panic(err)
    }
    
    // Create plugin manager
    manager := plugin.NewManager("my-extension", cfg)
    
    // Discover plugins
    plugins, err := manager.DiscoverPlugins(ctx)
    if err != nil {
        panic(err)
    }
    
    // Execute plugin
    if len(plugins) > 0 {
        err = manager.Execute(ctx, plugins[0].Path, []string{"arg1", "arg2"})
        if err != nil {
            panic(err)
        }
    }
}
```

## 🔧 Configuration

TykCtl-Go supports multiple configuration methods:

1. **Environment Variables** - `TYKCTL_*` variables
2. **Configuration Files** - YAML/JSON config files
3. **Command Line Flags** - Runtime configuration
4. **Extension-Specific Settings** - Per-extension configuration

See [Configuration Guide](config/README.md) for detailed information.

## 🔌 Plugin System

The plugin system provides:

- **Cross-Platform Support** - Linux, macOS, Windows, Unix
- **Timeout Configuration** - Configurable execution timeouts
- **Environment Setup** - Rich context for plugins
- **Discovery** - Automatic plugin discovery
- **Management** - Install, remove, list plugins

See [Plugin Documentation](plugin/README.md) for complete details.

## 🏗️ Extension Development

Create TykCtl extensions with:

- **Cobra Integration** - CLI command framework
- **Configuration Management** - Built-in config system
- **Plugin Support** - Plugin system integration
- **Hooks** - Event-driven automation
- **Templates** - Resource generation

See [Extension Guide](extension/README.md) for development details.

## 📖 API Reference

- **[API Client](api/README.md)** - HTTP client utilities
- **[Configuration API](config/README.md)** - Configuration management
- **[Plugin API](plugin/README.md)** - Plugin system API
- **[Extension API](extension/README.md)** - Extension framework API

## 🤝 Contributing

See [Development Guide](development.md) for:

- Setting up development environment
- Code style guidelines
- Testing requirements
- Pull request process

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

## 🆘 Support

- **Documentation**: Browse this documentation
- **Issues**: [GitHub Issues](https://github.com/edsonmichaque/tykctl-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/edsonmichaque/tykctl-go/discussions)