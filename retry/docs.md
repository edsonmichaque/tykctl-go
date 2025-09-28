# Retry Package Documentation

The retry package provides robust retry logic with exponential backoff for Go applications using `github.com/cenkalti/backoff/v4` underneath.

## Overview

The retry package is designed to handle transient failures in distributed systems by automatically retrying operations with intelligent backoff strategies. It provides both simple retry functions and advanced configuration options for complex scenarios.

## Features

- **Exponential Backoff**: Configurable exponential backoff with jitter
- **Context Support**: Full context.Context integration for cancellation and timeouts
- **Configurable**: Customizable retry parameters (max retries, delays, etc.)
- **Type Safe**: Generic support for operations that return results
- **Retryable Error Detection**: Smart error handling to avoid retrying non-retryable errors
- **Timeout Support**: Built-in timeout support for operations

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/retry"
)

func main() {
    ctx := context.Background()
    
    // Simple operation retry
    err := retry.Retry(ctx, nil, func() error {
        // Your operation here
        return someOperation()
    })
    
    if err != nil {
        fmt.Printf("Operation failed: %v\n", err)
    }
}
```

### With Custom Configuration

```go
config := &retry.Config{
    MaxRetries:     10,
    InitialDelay:   500 * time.Millisecond,
    MaxDelay:       30 * time.Second,
    BackoffFactor:  2.0,
    MaxElapsedTime: 10 * time.Minute,
}

err := retry.Retry(ctx, config, func() error {
    return someOperation()
})
```

### With Result

```go
result, err := retry.RetryWithResult(ctx, nil, func() (string, error) {
    return someOperationThatReturnsString()
})

if err != nil {
    fmt.Printf("Operation failed: %v\n", err)
} else {
    fmt.Printf("Result: %s\n", result)
}
```

## API Reference

### Core Functions

#### `Retry(ctx, config, operation)`
Executes an operation with retry logic using exponential backoff.

**Parameters:**
- `ctx`: Context for cancellation and timeout
- `config`: Retry configuration (nil for defaults)
- `operation`: Function to retry

**Returns:**
- `error`: Error if operation fails after all retries

#### `RetryWithResult[T](ctx, config, operation)`
Executes an operation that returns a result with retry logic.

**Parameters:**
- `ctx`: Context for cancellation and timeout
- `config`: Retry configuration (nil for defaults)
- `operation`: Function that returns (T, error)

**Returns:**
- `T`: Result of the operation
- `error`: Error if operation fails after all retries

#### `RetryWithTimeout(ctx, config, timeout, operation)`
Executes an operation with a timeout context.

**Parameters:**
- `ctx`: Context for cancellation
- `config`: Retry configuration (nil for defaults)
- `timeout`: Maximum time for the entire operation
- `operation`: Function to retry

**Returns:**
- `error`: Error if operation fails after all retries or timeout

#### `RetryWithTimeoutAndResult[T](ctx, config, timeout, operation)`
Executes an operation with a timeout context and returns a result.

**Parameters:**
- `ctx`: Context for cancellation
- `config`: Retry configuration (nil for defaults)
- `timeout`: Maximum time for the entire operation
- `operation`: Function that returns (T, error)

**Returns:**
- `T`: Result of the operation
- `error`: Error if operation fails after all retries or timeout

### Smart Error Handling

#### `IsRetryableError(err)`
Checks if an error should trigger a retry.

**Parameters:**
- `err`: Error to check

**Returns:**
- `bool`: True if error is retryable

#### `RetryIfRetryable(ctx, config, operation)`
Executes an operation with retry logic only if the error is retryable.

**Parameters:**
- `ctx`: Context for cancellation and timeout
- `config`: Retry configuration (nil for defaults)
- `operation`: Function to retry

**Returns:**
- `error`: Error if operation fails after all retries

#### `RetryIfRetryableWithResult[T](ctx, config, operation)`
Executes an operation with retry logic only if the error is retryable.

**Parameters:**
- `ctx`: Context for cancellation and timeout
- `config`: Retry configuration (nil for defaults)
- `operation`: Function that returns (T, error)

**Returns:**
- `T`: Result of the operation
- `error`: Error if operation fails after all retries

### Configuration

#### `Config` Struct

```go
type Config struct {
    MaxRetries     int           // Maximum number of retry attempts
    InitialDelay   time.Duration // Initial delay between retries
    MaxDelay       time.Duration // Maximum delay between retries
    BackoffFactor  float64       // Multiplier for exponential backoff
    MaxElapsedTime time.Duration // Maximum total time for all retries
}
```

#### `DefaultConfig()`
Returns default retry configuration:

```go
config := retry.DefaultConfig()
// MaxRetries: 5
// InitialDelay: 1 second
// MaxDelay: 30 seconds
// BackoffFactor: 2.0
// MaxElapsedTime: 5 minutes
```

## Usage Examples

### HTTP Client Retry

```go
func makeHTTPRequestWithRetry(ctx context.Context, url string) (*http.Response, error) {
    return retry.RetryWithResult(ctx, nil, func() (*http.Response, error) {
        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
            return nil, err
        }
        
        client := &http.Client{Timeout: 30 * time.Second}
        resp, err := client.Do(req)
        if err != nil {
            return nil, err
        }
        
        if resp.StatusCode >= 500 {
            return resp, fmt.Errorf("server error: %d", resp.StatusCode)
        }
        
        return resp, nil
    })
}
```

### Database Operation Retry

```go
func saveDataWithRetry(ctx context.Context, data Data) error {
    config := &retry.Config{
        MaxRetries:     3,
        InitialDelay:   1 * time.Second,
        MaxDelay:       10 * time.Second,
        BackoffFactor:  2.0,
        MaxElapsedTime: 2 * time.Minute,
    }
    
    return retry.Retry(ctx, config, func() error {
        return db.Save(data)
    })
}
```

### API Client with Smart Retry

```go
func callAPIWithSmartRetry(ctx context.Context, endpoint string) (*APIResponse, error) {
    return retry.RetryIfRetryableWithResult(ctx, nil, func() (*APIResponse, error) {
        resp, err := httpClient.Get(endpoint)
        if err != nil {
            return nil, err
        }
        
        if resp.StatusCode == 401 {
            // Authentication error - not retryable
            return nil, fmt.Errorf("authentication failed")
        }
        
        if resp.StatusCode >= 500 {
            // Server error - retryable
            return nil, fmt.Errorf("server error: %d", resp.StatusCode)
        }
        
        return parseResponse(resp), nil
    })
}
```

### Timeout with Retry

```go
func operationWithTimeout(ctx context.Context) error {
    config := &retry.Config{
        MaxRetries:     5,
        InitialDelay:   100 * time.Millisecond,
        MaxDelay:       1 * time.Second,
        BackoffFactor:  2.0,
        MaxElapsedTime: 30 * time.Second,
    }
    
    return retry.RetryWithTimeout(ctx, config, 10*time.Second, func() error {
        return someOperation()
    })
}
```

## Error Handling

### Retryable vs Non-Retryable Errors

The package automatically detects retryable vs non-retryable errors:

**Non-Retryable Errors:**
- `context.Canceled`
- `context.DeadlineExceeded`
- Authentication/authorization errors
- Invalid credentials
- Malformed requests

**Retryable Errors:**
- Network timeouts
- Temporary service unavailability
- Rate limiting (with backoff)
- General server errors (5xx)

### Custom Error Detection

You can implement custom error detection by wrapping non-retryable errors with `backoff.Permanent()`:

```go
err := retry.Retry(ctx, nil, func() error {
    result, err := someOperation()
    if err != nil {
        if isNonRetryableError(err) {
            return backoff.Permanent(err)
        }
        return err
    }
    return nil
})
```

## Best Practices

### 1. Use Appropriate Timeouts

Always use context timeouts to prevent operations from running indefinitely:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

err := retry.Retry(ctx, config, operation)
```

### 2. Configure Retry Parameters Wisely

- **MaxRetries**: Balance between resilience and performance
- **InitialDelay**: Start with reasonable delays (100ms-1s)
- **MaxDelay**: Cap delays to prevent excessive waiting
- **BackoffFactor**: Use 2.0 for exponential backoff
- **MaxElapsedTime**: Set reasonable total time limits

### 3. Handle Different Error Types

```go
err := retry.RetryIfRetryable(ctx, config, func() error {
    err := operation()
    if err != nil {
        // Log the error for debugging
        log.Printf("Operation failed: %v", err)
    }
    return err
})
```

### 4. Use Result Functions for Better Type Safety

```go
// Instead of using global variables
var result string
err := retry.Retry(ctx, nil, func() error {
    var err error
    result, err = someOperation()
    return err
})

// Use the result function
result, err := retry.RetryWithResult(ctx, nil, func() (string, error) {
    return someOperation()
})
```

### 5. Monitor and Log Retry Attempts

```go
config := &retry.Config{
    MaxRetries:     5,
    InitialDelay:   time.Second,
    MaxDelay:       30 * time.Second,
    BackoffFactor:  2.0,
    MaxElapsedTime: 5 * time.Minute,
}

err := retry.Retry(ctx, config, func() error {
    err := operation()
    if err != nil {
        log.Printf("Retry attempt failed: %v", err)
    }
    return err
})
```

## Integration with tykctl-go

The retry package is designed to work seamlessly with other tykctl-go packages:

### With Client Package

```go
import (
    "github.com/edsonmichaque/tykctl-go/retry"
    "github.com/edsonmichaque/tykctl/examples/extensions/tykctl-portal/pkg/client"
)

// Create client with retry configuration
retryConfig := &retry.Config{
    MaxRetries:     5,
    InitialDelay:   time.Second,
    MaxDelay:       30 * time.Second,
    BackoffFactor:  2.0,
    MaxElapsedTime: 5 * time.Minute,
}

portalClient := client.New("http://localhost:3001", "token",
    client.WithRetryConfig(retryConfig),
)

// Bootstrap with retry
req := client.BootstrapRequest{
    Username:  "admin@company.com",
    Password:  "password123",
    FirstName: "John",
    LastName:  "Doe",
}

result, err := portalClient.Bootstrap(ctx, req)
```

## Testing

The package includes comprehensive tests covering:

- Basic retry functionality
- Configuration options
- Error handling
- Context cancellation
- Timeout scenarios
- Smart error detection

Run tests with:

```bash
go test ./retry
```

## License

This package is part of the tykctl-go project and follows the same license terms.