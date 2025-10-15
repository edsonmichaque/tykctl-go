# Event Bus

The Event Bus provides a comprehensive event-driven architecture for tykctl-go. It's a generic event system that allows implementers to define their own event types and handlers.

## Features

- **Synchronous and Asynchronous Processing**: Choose between sync and async event processing
- **Rich Event Model**: Events with metadata, correlation IDs, and timestamps
- **Middleware Support**: Built-in middleware for logging, metrics, validation, rate limiting, and more
- **Handler Management**: Flexible handler registration and management
- **Error Handling**: Retry logic, circuit breakers, and timeout protection
- **Performance**: High-performance async processing with configurable workers
- **Monitoring**: Built-in metrics and statistics
- **No External Dependencies**: Self-contained with no external dependencies
- **Generic Design**: Implementers define their own event types

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/eventbus"
    "go.uber.org/zap"
)

// Define your event types
const (
    EventTypeAPICreate eventbus.EventType = "api.create"
    EventTypeAPIUpdate eventbus.EventType = "api.update"
)

func main() {
    // Create logger
    logger, _ := zap.NewDevelopment()
    
    // Create event bus
    bus := eventbus.New(
        eventbus.WithLogger(logger),
        eventbus.WithAsyncWorkers(5),
        eventbus.WithAsyncQueueSize(1000),
    )
    
    // Subscribe to events
    subscription, err := bus.Subscribe(EventTypeAPICreate, eventbus.HandlerFunc(
        func(ctx context.Context, event *eventbus.Event) error {
            log.Printf("API created: %+v", event.Data)
            return nil
        },
    ))
    if err != nil {
        log.Fatal(err)
    }
    defer subscription.Unsubscribe()
    
    // Publish events
    event := eventbus.NewEvent(EventTypeAPICreate, map[string]interface{}{
        "api_id": "api-123",
        "name": "My API",
    }).WithSource("tykctl-gateway")
    
    // Synchronous publishing
    err = bus.Publish(event)
    if err != nil {
        log.Fatal(err)
    }
    
    // Asynchronous publishing
    err = bus.PublishAsync(event)
    if err != nil {
        log.Fatal(err)
    }
    
    // Close the bus
    bus.Close()
}
```

## Event Types

The event bus is generic and does not define any predefined event types. Implementers should define their own event types based on their needs.

### Defining Event Types

```go
// Define your own event types
const (
    EventTypeAPICreate      eventbus.EventType = "api.create"
    EventTypeAPIUpdate      eventbus.EventType = "api.update"
    EventTypeAPIDelete      eventbus.EventType = "api.delete"
    EventTypeCommandStart   eventbus.EventType = "command.start"
    EventTypeCommandComplete eventbus.EventType = "command.complete"
    EventTypeSessionStart   eventbus.EventType = "session.start"
    EventTypeConfigLoad     eventbus.EventType = "config.load"
    // ... define as many as needed
)
```

### Event Type Naming Conventions

- Use descriptive, hierarchical names
- Use dot notation for namespacing (e.g., `api.create`, `command.start`)
- Use lowercase with hyphens for multi-word types
- Group related events by functionality

## Event Structure

Events have a rich structure with the following fields:

```go
type Event struct {
    ID            string                 `json:"id"`             // Unique event ID
    Type          EventType              `json:"type"`           // Event type
    Data          interface{}            `json:"data"`           // Event payload
    Metadata      map[string]interface{} `json:"metadata"`       // Additional metadata
    Timestamp     time.Time              `json:"timestamp"`      // When event occurred
    Source        string                 `json:"source"`         // Event source
    Version       string                 `json:"version"`        // Event schema version
    CorrelationID string                 `json:"correlation_id"` // Links related events
    ParentID      string                 `json:"parent_id"`      // Links to parent event
}
```

## Handlers

### Basic Handler

```go
handler := eventbus.HandlerFunc(func(ctx context.Context, event *eventbus.Event) error {
    log.Printf("Processing event: %s", event.Type)
    return nil
})
```

### Custom Handler

```go
type MyHandler struct {
    name string
}

func (h *MyHandler) Handle(ctx context.Context, event *eventbus.Event) error {
    log.Printf("Handler %s processing event: %s", h.name, event.Type)
    return nil
}

func (h *MyHandler) CanHandle(eventType eventbus.EventType) bool {
    return eventType == eventbus.EventTypeAPICreate
}

func (h *MyHandler) GetName() string {
    return h.name
}

func (h *MyHandler) GetPriority() int {
    return 100
}

func (h *MyHandler) GetTimeout() time.Duration {
    return 30 * time.Second
}
```

### Handler Wrappers

The event bus provides several handler wrappers for common patterns:

#### Filtered Handler
```go
filter := eventbus.EventFilter{
    Sources: []string{"tykctl-gateway"},
    Types:   []eventbus.EventType{eventbus.EventTypeAPICreate},
}
filteredHandler := eventbus.NewFilteredHandler(handler, filter)
```

#### Retry Handler
```go
retryHandler := eventbus.NewRetryHandler(handler, 3, 1*time.Second)
```

#### Timeout Handler
```go
timeoutHandler := eventbus.NewTimeoutHandler(handler, 30*time.Second)
```

#### Batch Handler
```go
batchHandler := eventbus.NewBatchHandler(handler, 10, 5*time.Second)
defer batchHandler.Close()
```

## Middleware

The event bus supports middleware for cross-cutting concerns:

### Logging Middleware
```go
logger, _ := zap.NewDevelopment()
loggingMiddleware := eventbus.NewLoggingMiddleware(logger)
bus.SetMiddleware(loggingMiddleware)
```

### Metrics Middleware
```go
metricsMiddleware := eventbus.NewMetricsMiddleware()
bus.SetMiddleware(metricsMiddleware)
```

### Validation Middleware
```go
validationMiddleware := eventbus.NewValidationMiddleware()
validationMiddleware.AddValidator(eventbus.EventTypeAPICreate, func(event *eventbus.Event) error {
    if event.Data == nil {
        return fmt.Errorf("event data is required")
    }
    return nil
})
bus.SetMiddleware(validationMiddleware)
```

### Rate Limiting Middleware
```go
rateLimitMiddleware := eventbus.NewRateLimitMiddleware()
rateLimitMiddleware.SetRateLimit(eventbus.EventTypeAPICreate, 100, 1*time.Minute)
bus.SetMiddleware(rateLimitMiddleware)
```

### Circuit Breaker Middleware
```go
circuitBreakerMiddleware := eventbus.NewCircuitBreakerMiddleware()
circuitBreakerMiddleware.SetCircuitBreaker(eventbus.EventTypeAPICreate, 5, 30*time.Second)
bus.SetMiddleware(circuitBreakerMiddleware)
```

### Retry Middleware
```go
retryMiddleware := eventbus.NewRetryMiddleware(3, 1*time.Second)
bus.SetMiddleware(retryMiddleware)
```

### Timeout Middleware
```go
timeoutMiddleware := eventbus.NewTimeoutMiddleware(30*time.Second)
bus.SetMiddleware(timeoutMiddleware)
```

### Sanitization Middleware
```go
sanitizationMiddleware := eventbus.NewSanitizationMiddleware()
sanitizationMiddleware.AddSanitizer("password", func(v interface{}) interface{} {
    return "[REDACTED]"
})
bus.SetMiddleware(sanitizationMiddleware)
```

## Configuration

The event bus can be configured with various options:

```go
bus := eventbus.New(
    eventbus.WithAsyncWorkers(10),
    eventbus.WithAsyncQueueSize(1000),
    eventbus.WithLogger(logger),
    eventbus.WithDefaultTimeout(30*time.Second),
    eventbus.WithMaxRetries(3),
    eventbus.WithRetryDelay(1*time.Second),
    eventbus.WithMetrics(true),
    eventbus.WithLogging(true),
    eventbus.WithValidation(true),
    eventbus.WithRateLimit(true),
    eventbus.WithCircuitBreaker(true),
    eventbus.WithSanitization(true),
)
```

## Statistics

The event bus provides statistics about its operation:

```go
stats := bus.GetStats()
fmt.Printf("Events published: %d\n", stats.EventsPublished)
fmt.Printf("Events processed: %d\n", stats.EventsProcessed)
fmt.Printf("Events failed: %d\n", stats.EventsFailed)
fmt.Printf("Active subscriptions: %d\n", stats.ActiveSubscriptions)
```

## Custom Event Types

Since the event bus is generic, you need to define your own event types. Here's how to create a complete event system:

### 1. Define Event Types
```go
// Define your event types
const (
    EventTypeAPICreate      eventbus.EventType = "api.create"
    EventTypeAPIUpdate      eventbus.EventType = "api.update"
    EventTypeAPIDelete      eventbus.EventType = "api.delete"
    EventTypeCommandStart   eventbus.EventType = "command.start"
    EventTypeCommandComplete eventbus.EventType = "command.complete"
)
```

### 2. Create Event Bus
```go
// Create event bus
eventBus := eventbus.New(
    eventbus.WithLogger(logger),
    eventbus.WithAsyncWorkers(5),
)
```

### 3. Subscribe to Events
```go
// Subscribe to your custom events
subscription, err := eventBus.Subscribe(EventTypeAPICreate, eventbus.HandlerFunc(
    func(ctx context.Context, event *eventbus.Event) error {
        log.Printf("API created: %+v", event.Data)
        return nil
    },
))
```

### 4. Publish Events
```go
// Publish events using your custom types
event := eventbus.NewEvent(EventTypeAPICreate, apiData)
eventBus.Publish(event)
```

## Best Practices

1. **Use appropriate event types**: Choose specific event types for better filtering and handling
2. **Handle errors gracefully**: Always handle errors in your event handlers
3. **Use correlation IDs**: Link related events with correlation IDs
4. **Set appropriate timeouts**: Configure timeouts based on your handler requirements
5. **Monitor performance**: Use metrics middleware to monitor event processing performance
6. **Sanitize sensitive data**: Use sanitization middleware to protect sensitive information
7. **Use async processing**: Use async publishing for non-critical events to improve performance
8. **Implement proper logging**: Use logging middleware for debugging and monitoring

## Examples

See the `examples/` directory for complete examples of using the event bus in different scenarios.