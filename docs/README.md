# TykCtl-Go Documentation

Welcome to the TykCtl-Go documentation. This directory contains comprehensive documentation for the TykCtl Go framework and its components.

## 📚 Documentation Structure

### Core Components

- **[API Documentation](api/)** - API client utilities and helpers
- **[Configuration Management](config/)** - Configuration system, environment variables, and discovery
- **[Plugin System](plugin/)** - Cross-platform plugin management and execution
- **[Extension Framework](extension/)** - Extension development and management
- **[Alias System](alias/)** - Command alias management and execution

### Features

- **[Hooks System](hooks/)** - Event-driven hooks and automation
- **[Templates](templates/)** - Template system for resource generation
- **[Progress Tracking](progress/)** - Progress indicators and status tracking
- **[Alias Management](alias/)** - Command shortcuts and automation

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
    "github.com/edsonmichaque/tykctl-go/alias"
    "github.com/spf13/cobra"
)

func main() {
    ctx := context.Background()
    
    // Initialize configuration
    cfg, err := config.NewConfigManager()
    if err != nil {
        panic(err)
    }
    
    // Create plugin manager
    pluginManager := plugin.NewManager("my-extension", cfg)
    
    // Create alias manager
    aliasProvider := alias.NewExtensionConfigProvider(
        setAliasFunc, getAliasFunc, deleteAliasFunc, listAliasesFunc,
    )
    aliasManager := alias.NewManager(aliasProvider, []string{"help", "version"})
    
    // Create root command
    rootCmd := &cobra.Command{Use: "myapp"}
    
    // Add alias command
    aliasBuilder := alias.NewCommandBuilder(aliasManager)
    rootCmd.AddCommand(aliasBuilder.BuildAliasCommand())
    
    // Register aliases as subcommands
    aliasRegistrar := alias.NewRegistrar(aliasManager)
    aliasRegistrar.RegisterAliases(ctx, rootCmd)
    
    // Execute command
    rootCmd.Execute()
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

## 🔗 Alias System

The alias system provides:

- **Command Shortcuts** - Create shortcuts for commonly used commands
- **Shell Integration** - Execute shell commands from aliases
- **Parameter Expansion** - Support for `$1`, `$2`, `$*`, `$@` parameters
- **Validation** - Comprehensive alias name and expansion validation
- **Cobra Integration** - Seamless integration with Cobra commands

See [Alias Documentation](alias/README.md) for complete details.

## 🏗️ Extension Development

Create TykCtl extensions with:

- **Cobra Integration** - CLI command framework
- **Configuration Management** - Built-in config system
- **Plugin Support** - Plugin system integration
- **Alias System** - Command alias management
- **Hooks** - Event-driven automation
- **Templates** - Resource generation

See [Extension Guide](extension/README.md) for development details.

## 📖 API Reference

- **[API Client](api/README.md)** - HTTP client utilities
- **[Configuration API](config/README.md)** - Configuration management
- **[Plugin API](plugin/README.md)** - Plugin system API
- **[Alias API](alias/README.md)** - Alias system API
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