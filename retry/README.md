# Retry Package

A robust retry package for Go that provides exponential backoff retry logic using `github.com/cenkalti/backoff/v4` underneath.

## Features

- **Exponential Backoff**: Configurable exponential backoff with jitter
- **Context Support**: Full context.Context integration for cancellation and timeouts
- **Configurable**: Customizable retry parameters (max retries, delays, etc.)
- **Type Safe**: Generic support for operations that return results
- **Retryable Error Detection**: Smart error handling to avoid retrying non-retryable errors
- **Timeout Support**: Built-in timeout support for operations

## Installation

```bash
go get github.com/edsonmichaque/tykctl-go
```

Then import the retry package:

```go
import "github.com/edsonmichaque/tykctl-go/retry"
```

## Usage

### Basic Retry

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

### Retry with Custom Configuration

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

### Retry with Result

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

### Retry with Timeout

```go
// Retry with 5-minute timeout
err := retry.RetryWithTimeout(ctx, nil, 5*time.Minute, func() error {
    return someOperation()
})
```

### Retry with Timeout and Result

```go
result, err := retry.RetryWithTimeoutAndResult(ctx, nil, 5*time.Minute, func() (int, error) {
    return someOperationThatReturnsInt()
})
```

### Smart Retry (Only Retry Retryable Errors)

```go
// Only retry if the error is retryable (not auth errors, etc.)
err := retry.RetryIfRetryable(ctx, nil, func() error {
    return someOperation()
})
```

### Retry with Result and Smart Error Handling

```go
result, err := retry.RetryIfRetryableWithResult(ctx, nil, func() (Data, error) {
    return someOperationThatReturnsData()
})
```

## Configuration

The `Config` struct allows you to customize retry behavior:

```go
type Config struct {
    MaxRetries     int           // Maximum number of retry attempts
    InitialDelay   time.Duration // Initial delay between retries
    MaxDelay       time.Duration // Maximum delay between retries
    BackoffFactor  float64       // Multiplier for exponential backoff
    MaxElapsedTime time.Duration // Maximum total time for all retries
}
```

### Default Configuration

```go
config := retry.DefaultConfig()
// MaxRetries: 5
// InitialDelay: 1 second
// MaxDelay: 30 seconds
// BackoffFactor: 2.0
// MaxElapsedTime: 5 minutes
```

## Error Handling

The package provides smart error handling:

- **Context Cancellation**: Respects context cancellation and deadlines
- **Non-Retryable Errors**: Automatically detects and stops retrying for certain error types
- **Permanent Errors**: Uses `backoff.Permanent()` to wrap non-retryable errors

### Retryable vs Non-Retryable Errors

**Non-Retryable Errors:**
- `context.Canceled`
- `context.DeadlineExceeded`
- Authentication/authorization errors
- Invalid credentials

**Retryable Errors:**
- Network timeouts
- Temporary service unavailability
- Rate limiting (with backoff)
- General server errors

## Examples

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

## License

This package is part of the tykctl-go project and follows the same license terms.