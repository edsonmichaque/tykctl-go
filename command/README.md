# Command Package

The `command` package provides enhanced command functionality built on top of Cobra, adding logging and context support for tykctl extensions.

## Features

- **Cobra Integration**: Built on top of the popular Cobra CLI framework
- **Logging Support**: Integrated logging with zap logger
- **Context Support**: Full context.Context integration for cancellation and timeouts
- **Fluent API**: Method chaining for clean command configuration
- **Extension Ready**: Designed specifically for tykctl extensions

## Usage

### Basic Command Creation

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/edsonmichaque/tykctl-go/command"
    "go.uber.org/zap"
)

func main() {
    // Create a new command
    cmd := command.New("hello", "Say hello", func(cmd *cobra.Command, args []string) error {
        fmt.Println("Hello, World!")
        return nil
    })
    
    // Execute the command
    cmd.Execute()
}
```

### Command with Long Description

```go
cmd := command.NewWithLong(
    "install",
    "Install an extension",
    "Install a tykctl extension from a repository",
    func(cmd *cobra.Command, args []string) error {
        // Installation logic
        return nil
    },
)
```

### Command with Logger

```go
import (
    "go.uber.org/zap"
    "github.com/edsonmichaque/tykctl-go/command"
)

func main() {
    // Create logger
    logger, _ := zap.NewDevelopment()
    
    // Create command with logger
    cmd := command.New("test", "Test command", func(cmd *cobra.Command, args []string) error {
        // Access logger from command
        if cmdLogger := cmd.Context().Value("logger"); cmdLogger != nil {
            if zapLogger, ok := cmdLogger.(*zap.Logger); ok {
                zapLogger.Info("Command executed")
            }
        }
        return nil
    }).SetLogger(logger)
    
    cmd.Execute()
}
```

### Command with Context

```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    cmd := command.New("long-running", "Long running command", func(cmd *cobra.Command, args []string) error {
        // Access context from command
        cmdCtx := cmd.Context()
        
        // Use context for cancellation
        select {
        case <-cmdCtx.Done():
            return cmdCtx.Err()
        case <-time.After(5 * time.Second):
            fmt.Println("Command completed")
            return nil
        }
    }).SetContext(ctx)
    
    cmd.Execute()
}
```

### Fluent API Usage

```go
cmd := command.New("complex", "Complex command", handler).
    SetLogger(logger).
    SetContext(ctx).
    AddFlag("verbose", "v", false, "Enable verbose output").
    AddStringFlag("config", "c", "", "Configuration file")
```

## Advanced Usage

### Custom Command with Flags

```go
func createInstallCommand() *command.Command {
    cmd := command.NewWithLong(
        "install",
        "Install extension",
        "Install a tykctl extension from GitHub repository",
        installHandler,
    )
    
    // Add flags
    cmd.Flags().StringP("repo", "r", "", "Repository URL")
    cmd.Flags().StringP("version", "v", "latest", "Extension version")
    cmd.Flags().BoolP("force", "f", false, "Force installation")
    
    return cmd
}

func installHandler(cmd *cobra.Command, args []string) error {
    repo, _ := cmd.Flags().GetString("repo")
    version, _ := cmd.Flags().GetString("version")
    force, _ := cmd.Flags().GetBool("force")
    
    fmt.Printf("Installing %s@%s (force: %v)\n", repo, version, force)
    return nil
}
```

### Command with Subcommands

```go
func createRootCommand() *command.Command {
    rootCmd := command.New("tykctl", "Tykctl CLI tool", nil)
    
    // Add subcommands
    rootCmd.AddCommand(createInstallCommand())
    rootCmd.AddCommand(createListCommand())
    rootCmd.AddCommand(createRemoveCommand())
    
    return rootCmd
}
```

### Command with Validation

```go
func createValidatedCommand() *command.Command {
    cmd := command.New("validate", "Validate configuration", func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            return fmt.Errorf("configuration file required")
        }
        
        configFile := args[0]
        if _, err := os.Stat(configFile); os.IsNotExist(err) {
            return fmt.Errorf("configuration file not found: %s", configFile)
        }
        
        fmt.Printf("Validating %s...\n", configFile)
        return nil
    })
    
    return cmd
}
```

## Integration with Other Packages

### With Logger Package

```go
import (
    "github.com/edsonmichaque/tykctl-go/command"
    "github.com/edsonmichaque/tykctl-go/logger"
)

func main() {
    // Create logger
    logConfig := logger.Config{
        Debug:   true,
        Verbose: false,
        NoColor: false,
    }
    zapLogger := logger.New(logConfig)
    
    // Create command with logger
    cmd := command.New("test", "Test command", handler).
        SetLogger(zapLogger.Logger)
    
    cmd.Execute()
}
```

### With Extension Package

```go
func createExtensionCommand() *command.Command {
    cmd := command.New("extension", "Manage extensions", nil)
    
    // Install subcommand
    installCmd := command.New("install", "Install extension", func(cmd *cobra.Command, args []string) error {
        if len(args) < 2 {
            return fmt.Errorf("usage: extension install <owner> <repo>")
        }
        
        installer := extension.NewInstaller("/tmp/tykctl-config")
        return installer.InstallExtension(cmd.Context(), args[0], args[1])
    })
    
    cmd.AddCommand(installCmd)
    return cmd
}
```

## Command Structure

### Command Type

```go
type Command struct {
    *cobra.Command
    logger  *zap.Logger
    context context.Context
}
```

### Available Methods

- `SetLogger(logger *zap.Logger) *Command` - Set the logger
- `SetContext(ctx context.Context) *Command` - Set the context
- `GetLogger() *zap.Logger` - Get the logger
- `GetContext() context.Context` - Get the context

## Best Practices

- **Error Handling**: Always return errors from command handlers
- **Context Usage**: Use context for cancellation and timeouts
- **Logging**: Use the provided logger for consistent logging
- **Validation**: Validate arguments and flags early in command handlers
- **Help Text**: Provide clear help text and examples
- **Subcommands**: Use subcommands for complex functionality

## Examples

### Complete CLI Application

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/command"
    "github.com/edsonmichaque/tykctl-go/logger"
    "go.uber.org/zap"
)

func main() {
    // Create logger
    logConfig := logger.Config{Debug: true}
    zapLogger := logger.New(logConfig)
    
    // Create root command
    rootCmd := command.New("myapp", "My CLI Application", nil).
        SetLogger(zapLogger.Logger)
    
    // Add subcommands
    rootCmd.AddCommand(createHelloCommand(zapLogger.Logger))
    rootCmd.AddCommand(createVersionCommand())
    
    // Execute
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func createHelloCommand(logger *zap.Logger) *command.Command {
    return command.New("hello", "Say hello", func(cmd *cobra.Command, args []string) error {
        logger.Info("Hello command executed")
        fmt.Println("Hello, World!")
        return nil
    }).SetLogger(logger)
}

func createVersionCommand() *command.Command {
    return command.New("version", "Show version", func(cmd *cobra.Command, args []string) error {
        fmt.Println("v1.0.0")
        return nil
    })
}
```

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `go.uber.org/zap` - Logging library