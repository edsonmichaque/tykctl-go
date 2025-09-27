# Logger Package

The `logger` package provides a structured logging solution built on top of Zap, offering configurable logging levels, colored output, and production-ready features for tykctl applications.

## Features

- **Zap Integration**: Built on the high-performance Zap logging library
- **Configurable Levels**: Support for Debug, Info, Warn, and Error levels
- **Colored Output**: Optional colored output for better readability
- **Production Ready**: Optimized for both development and production environments
- **Environment Aware**: Automatic configuration based on environment variables
- **Structured Logging**: JSON-formatted logs with structured fields

## Usage

### Basic Logging

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/logger"
)

func main() {
    // Create logger with default configuration
    config := logger.Config{
        Debug:   false,
        Verbose: false,
        NoColor: false,
    }
    
    zapLogger := logger.New(config)
    
    // Use the logger
    zapLogger.Info("Application started")
    zapLogger.Debug("Debug information")
    zapLogger.Warn("Warning message")
    zapLogger.Error("Error occurred")
    
    fmt.Println("Logging completed!")
}
```

### Development Logging

```go
func developmentLogging() {
    // Create logger for development
    config := logger.Config{
        Debug:   true,  // Enable debug level
        Verbose: true,  // Enable verbose output
        NoColor: false, // Enable colors
    }
    
    zapLogger := logger.New(config)
    
    // Log with structured fields
    zapLogger.Info("User logged in",
        zap.String("user_id", "123"),
        zap.String("email", "user@example.com"),
        zap.Duration("login_time", time.Second*2),
    )
    
    // Debug logging
    zapLogger.Debug("Processing request",
        zap.String("method", "POST"),
        zap.String("path", "/api/users"),
        zap.Int("status_code", 201),
    )
}
```

### Production Logging

```go
func productionLogging() {
    // Create logger for production
    config := logger.Config{
        Debug:   false, // Disable debug in production
        Verbose: false, // Disable verbose output
        NoColor: true,  // Disable colors for production
    }
    
    zapLogger := logger.New(config)
    
    // Structured logging for production
    zapLogger.Info("Request processed",
        zap.String("method", "GET"),
        zap.String("path", "/api/data"),
        zap.Int("status_code", 200),
        zap.Duration("duration", time.Millisecond*150),
        zap.String("user_agent", "MyApp/1.0.0"),
    )
    
    // Error logging with context
    zapLogger.Error("Database connection failed",
        zap.String("database", "postgres"),
        zap.String("host", "localhost:5432"),
        zap.Error(err),
    )
}
```

## Advanced Usage

### Custom Logger Configuration

```go
func customLoggerConfig() {
    // Check environment variables
    debug := os.Getenv("DEBUG") == "true"
    verbose := os.Getenv("VERBOSE") == "true"
    noColor := os.Getenv("NO_COLOR") != ""
    
    config := logger.Config{
        Debug:   debug,
        Verbose: verbose,
        NoColor: noColor,
    }
    
    zapLogger := logger.New(config)
    
    // Use logger
    zapLogger.Info("Logger configured from environment")
}
```

### Logger with Context

```go
func loggerWithContext() {
    config := logger.Config{Debug: true}
    zapLogger := logger.New(config)
    
    // Create context with logger
    ctx := context.WithValue(context.Background(), "logger", zapLogger)
    
    // Use logger in context
    if logger, ok := ctx.Value("logger").(*logger.Logger); ok {
        logger.Info("Logging from context")
    }
}
```

### Structured Logging Patterns

```go
func structuredLogging() {
    config := logger.Config{Debug: true}
    zapLogger := logger.New(config)
    
    // Request logging
    zapLogger.Info("HTTP request",
        zap.String("method", "POST"),
        zap.String("url", "/api/users"),
        zap.String("user_id", "123"),
        zap.Duration("duration", time.Millisecond*250),
        zap.Int("status_code", 201),
    )
    
    // Business logic logging
    zapLogger.Info("User created",
        zap.String("user_id", "123"),
        zap.String("email", "user@example.com"),
        zap.String("role", "admin"),
        zap.Time("created_at", time.Now()),
    )
    
    // Error logging with stack trace
    zapLogger.Error("Failed to process payment",
        zap.String("payment_id", "pay_123"),
        zap.String("user_id", "123"),
        zap.Float64("amount", 99.99),
        zap.Error(err),
    )
}
```

## Integration Examples

### With HTTP Client

```go
import (
    "github.com/edsonmichaque/tykctl-go/httpclient"
    "github.com/edsonmichaque/tykctl-go/logger"
)

func httpClientWithLogging() error {
    // Create logger
    config := logger.Config{Debug: true}
    zapLogger := logger.New(config)
    
    // Create HTTP client
    client := httpclient.NewWithBaseURL("https://api.example.com")
    client.SetHeader("User-Agent", "MyApp/1.0.0")
    
    ctx := context.Background()
    
    // Log request
    zapLogger.Info("Making API request",
        zap.String("method", "GET"),
        zap.String("url", "/users"),
    )
    
    // Make request
    resp, err := client.Get(ctx, "/users")
    if err != nil {
        zapLogger.Error("API request failed",
            zap.String("method", "GET"),
            zap.String("url", "/users"),
            zap.Error(err),
        )
        return err
    }
    
    // Log response
    zapLogger.Info("API request completed",
        zap.String("method", "GET"),
        zap.String("url", "/users"),
        zap.Int("status_code", resp.StatusCode),
        zap.Int("response_size", len(resp.Body)),
    )
    
    return nil
}
```

### With Command Package

```go
import (
    "github.com/edsonmichaque/tykctl-go/command"
    "github.com/edsonmichaque/tykctl-go/logger"
)

func commandWithLogging() {
    // Create logger
    config := logger.Config{Debug: true}
    zapLogger := logger.New(config)
    
    // Create command with logger
    cmd := command.New("test", "Test command", func(cmd *cobra.Command, args []string) error {
        // Access logger from command
        if cmdLogger := cmd.Context().Value("logger"); cmdLogger != nil {
            if zapLogger, ok := cmdLogger.(*logger.Logger); ok {
                zapLogger.Info("Command executed",
                    zap.String("command", "test"),
                    zap.Strings("args", args),
                )
            }
        }
        return nil
    }).SetLogger(zapLogger.Logger)
    
    cmd.Execute()
}
```

### With Extension Package

```go
func extensionWithLogging() error {
    // Create logger
    config := logger.Config{Debug: true}
    zapLogger := logger.New(config)
    
    // Create extension installer with logger
    installer := extension.NewInstaller(
        "/tmp/tykctl-config",
        extension.WithLogger(zapLogger.Logger),
    )
    
    ctx := context.Background()
    
    // Log installation start
    zapLogger.Info("Starting extension installation",
        zap.String("owner", "owner"),
        zap.String("repo", "repo"),
    )
    
    // Install extension
    err := installer.InstallExtension(ctx, "owner", "repo")
    if err != nil {
        zapLogger.Error("Extension installation failed",
            zap.String("owner", "owner"),
            zap.String("repo", "repo"),
            zap.Error(err),
        )
        return err
    }
    
    // Log installation success
    zapLogger.Info("Extension installed successfully",
        zap.String("owner", "owner"),
        zap.String("repo", "repo"),
    )
    
    return nil
}
```

## Configuration Options

### Config Structure

```go
type Config struct {
    Debug   bool // Enable debug level logging
    Verbose bool // Enable verbose output
    NoColor bool // Disable colored output
}
```

### Environment Variables

- `DEBUG=true` - Enable debug level logging
- `VERBOSE=true` - Enable verbose output
- `NO_COLOR=1` - Disable colored output

### Log Levels

- **Debug**: Detailed information for debugging
- **Info**: General information about application flow
- **Warn**: Warning messages for potential issues
- **Error**: Error messages for failures

## Performance Considerations

### Production Optimization

```go
func productionOptimization() {
    // Production configuration
    config := logger.Config{
        Debug:   false, // Disable debug for performance
        Verbose: false, // Disable verbose output
        NoColor: true,  // Disable colors for performance
    }
    
    zapLogger := logger.New(config)
    
    // Use structured logging efficiently
    zapLogger.Info("High-performance logging",
        zap.String("service", "api"),
        zap.Int("request_count", 1000),
        zap.Duration("avg_response_time", time.Millisecond*50),
    )
}
```

### Conditional Logging

```go
func conditionalLogging() {
    config := logger.Config{Debug: true}
    zapLogger := logger.New(config)
    
    // Only log expensive operations in debug mode
    if zapLogger.Core().Enabled(zap.DebugLevel) {
        zapLogger.Debug("Expensive operation completed",
            zap.Duration("duration", time.Second*5),
            zap.Int("items_processed", 10000),
        )
    }
}
```

## Best Practices

- **Structured Logging**: Use structured fields instead of string formatting
- **Appropriate Levels**: Use the right log level for each message
- **Context Information**: Include relevant context in log messages
- **Performance**: Avoid expensive operations in log statements
- **Production Ready**: Configure appropriately for production environments
- **Error Handling**: Always log errors with context and stack traces

## Dependencies

- `go.uber.org/zap` - High-performance structured logging library
- `go.uber.org/zap/zapcore` - Zap core functionality