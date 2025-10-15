# Telemetry Package

The telemetry package provides anonymous usage analytics for tykctl-go, helping improve the CLI tool while respecting user privacy.

## Features

- **Anonymous Data Collection**: Collects usage statistics without personal information
- **Privacy-First Design**: No sensitive data (API keys, tokens, etc.) is collected
- **Opt-in/Opt-out**: Users can easily enable or disable telemetry
- **Configurable**: Multiple configuration options and environment variables
- **Robust**: Handles network failures gracefully with retry logic
- **Non-blocking**: Telemetry doesn't impact CLI performance
- **Batched Transmission**: Events are batched and sent periodically

## Quick Start

```go
package main

import (
    "log"
    "time"
    
    "go.uber.org/zap"
    "github.com/edsonmichaque/tykctl-go/telemetry"
)

func main() {
    // Create logger
    logger, _ := zap.NewDevelopment()
    defer logger.Sync()

    // Create telemetry manager
    manager, err := telemetry.NewManager(logger)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Close()

    // Track a command execution
    startTime := time.Now()
    // ... execute your command ...
    duration := time.Since(startTime)
    
    if err := manager.TrackCommand("tykctl config get", duration, true); err != nil {
        logger.Error("Failed to track command", zap.Error(err))
    }
}
```

## Configuration

### Configuration File

Telemetry configuration is stored in `~/.config/tykctl/telemetry.yaml`:

```yaml
enabled: true
endpoint: "https://telemetry.tyk.io/v1/events"
batch_size: 100
flush_interval: "5m"
retry_attempts: 3
retry_delay: "1s"
timeout: "30s"
user_agent: "tykctl-go-telemetry/1.0"
```

### Environment Variables

- `TYKCTL_TELEMETRY_ENABLED`: Enable/disable telemetry (true/false)
- `TYKCTL_TELEMETRY_ENDPOINT`: Custom telemetry endpoint
- `TYKCTL_NO_TELEMETRY`: Disable telemetry for current session

### CLI Commands

```bash
# Check telemetry status
tykctl telemetry status

# Enable telemetry
tykctl telemetry enable

# Disable telemetry
tykctl telemetry disable

# Show configuration
tykctl telemetry config

# Flush pending events
tykctl telemetry flush
```

## Data Collection

### What Data is Collected

- **Command Usage**: Which commands are executed and how frequently
- **Performance Metrics**: Command execution times and success rates
- **Feature Usage**: Which features are used most
- **Error Information**: Anonymous error types and frequencies
- **System Information**: Operating system, CLI version, architecture
- **Session Data**: Anonymous session and user identifiers

### What Data is NOT Collected

- **Personal Information**: No usernames, emails, or personal data
- **Sensitive Data**: No API keys, tokens, passwords, or credentials
- **Repository Data**: No repository names or content
- **Network Data**: No IP addresses or network information
- **File Content**: No file contents or sensitive configuration data

## Usage Examples

### Basic Command Tracking

```go
// Track command execution
manager.TrackCommand("tykctl gateway create", duration, success)
```

### Feature Usage Tracking

```go
// Track feature usage
manager.TrackFeature("gateway_creation", map[string]interface{}{
    "template_used": "basic",
    "plugins_count": 3,
})
```

### API Call Tracking

```go
// Track API calls
manager.TrackAPICall("/api/gateways", "GET", 200, duration)
```

### Error Tracking

```go
// Track errors
manager.TrackError("validation_error", "Invalid configuration")
```

### Performance Tracking

```go
// Track performance metrics
manager.TrackPerformance("config_load_time", 45, map[string]interface{}{
    "file_size": 1024,
    "format": "yaml",
})
```

### Extension Tracking

```go
// Track extension usage
manager.TrackExtensionUsage("tykctl-plugin-auth", success, duration)
```

## Integration Patterns

### Command Wrapper

```go
// Wrap command functions with telemetry
integration := telemetry.NewIntegration(manager, logger)
wrappedFunc := integration.CommandWrapper(yourCommandFunc)
err := wrappedFunc()
```

### API Client Wrapper

```go
// Wrap API clients with telemetry
wrapper := telemetry.NewAPIClientWrapper(apiClient, manager, logger, "/api/gateways")
wrapper.TrackRequest("GET", startTime, err)
```

### Extension Wrapper

```go
// Wrap extension execution with telemetry
wrapper := telemetry.NewExtensionWrapper(manager, logger)
err := wrapper.ExecuteExtension("my-extension", extensionFunc)
```

## Event Builder

Use the fluent event builder for complex events:

```go
// Build custom events
event := telemetry.NewEventBuilder(telemetry.EventTypeCommand).
    Command("tykctl gateway create").
    Duration(1*time.Second).
    Success(true).
    Property("gateway_name", "my-gateway").
    Property("gateway_type", "http").
    Build()

manager.client.Track(event)
```

## Storage

### File Storage (Default)

Events are stored in `~/.cache/tykctl/telemetry/` and sent in batches.

### Memory Storage

For testing or when file storage is unavailable, events are stored in memory.

## Transport

### HTTP Transport (Default)

Events are sent via HTTPS POST requests to the configured endpoint.

### No-Op Transport

For testing or when telemetry is disabled, events are not sent.

## Privacy and Security

### Data Sanitization

All events are automatically sanitized to remove sensitive information:

```go
// Sensitive keys are automatically removed
sensitiveKeys := []string{
    "token", "key", "secret", "password", "auth",
    "credential", "api_key", "access_token", "refresh_token",
}
```

### Anonymous Identifiers

- **Session ID**: Generated per CLI session
- **User ID**: Generated and hashed for anonymity

### Error Message Sanitization

Error messages are sanitized to remove sensitive patterns.

## Best Practices

### 1. Always Check if Telemetry is Enabled

```go
if !manager.IsEnabled() {
    return // Skip telemetry
}
```

### 2. Use Non-blocking Tracking

```go
// Don't block command execution for telemetry failures
if err := manager.TrackCommand(cmd, duration, success); err != nil {
    logger.Debug("Telemetry failed", zap.Error(err))
    // Continue with command execution
}
```

### 3. Batch Related Events

```go
// Group related events together
manager.TrackFeature("gateway_creation", map[string]interface{}{
    "step": "validation",
    "template": "basic",
})
```

### 4. Use Appropriate Event Types

- `EventTypeCommand`: For command executions
- `EventTypeError`: For errors and failures
- `EventTypeFeature`: For feature usage
- `EventTypePerformance`: For performance metrics

## Troubleshooting

### Check Telemetry Status

```bash
tykctl telemetry status
```

### View Configuration

```bash
tykctl telemetry config
```

### Flush Events Manually

```bash
tykctl telemetry flush
```

### Disable Telemetry

```bash
# Via CLI
tykctl telemetry disable

# Via environment variable
export TYKCTL_NO_TELEMETRY=1

# Via configuration
tykctl telemetry config
# Edit the configuration file to set enabled: false
```

### Debug Telemetry

```bash
# Enable debug logging
export TYKCTL_DEBUG=true
tykctl telemetry status
```

## Contributing

When adding new telemetry events:

1. **Follow Privacy Guidelines**: Never collect sensitive data
2. **Use Appropriate Event Types**: Choose the right event type
3. **Add Properties Carefully**: Only add necessary, non-sensitive properties
4. **Test Thoroughly**: Ensure telemetry doesn't break functionality
5. **Document Changes**: Update this README for new features

## License

This package is part of tykctl-go and follows the same license terms.