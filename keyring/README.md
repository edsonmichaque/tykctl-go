# Keyring Package

A simple wrapper around the [zalando/go-keyring](https://github.com/zalando/go-keyring) library that provides secure storage and retrieval of secrets using the system keyring.

## Features

- **Cross-platform support**: Works on macOS, Linux/Unix, and Windows
- **Simple API**: Easy-to-use interface for storing and retrieving secrets
- **Secure**: Uses the system's native keyring services
- **Lightweight wrapper**: Minimal overhead over the underlying zalando/go-keyring library

## Supported Platforms

- **macOS**: Uses the macOS Keychain via the `security` command-line tool
- **Linux/Unix**: Uses the Secret Service API via D-Bus
- **Windows**: Uses the Windows Credential Manager
- **Fallback**: Returns errors for unsupported platforms

## Installation

```bash
go get github.com/edsonmichaque/tykctl-go/keyring
```

## Usage

### Basic Operations

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/keyring"
)

func main() {
    ctx := context.Background()
    
    // Store a secret
    err := keyring.Set(ctx, "my-service", "my-user", "my-secret")
    if err != nil {
        log.Fatal(err)
    }
    
    // Retrieve a secret
    secret, err := keyring.Get(ctx, "my-service", "my-user")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(secret)
    
    // Delete a secret
    err = keyring.Delete(ctx, "my-service", "my-user")
    if err != nil {
        log.Fatal(err)
    }
    
    // Purge all secrets for a service
    err = keyring.Purge(ctx, "my-service")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Context Support

The package supports context for cancellation and timeouts:

```go
package main

import (
    "context"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/keyring"
)

func main() {
    // Example with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    secret, err := keyring.Get(ctx, "service", "user")
    if err != nil {
        // Handle timeout or other errors
        log.Fatal(err)
    }
    
    // Example with cancellation
    ctx2, cancel2 := context.WithCancel(context.Background())
    cancel2() // Cancel immediately
    
    err = keyring.Set(ctx2, "service", "user", "password")
    if err == context.Canceled {
        fmt.Println("Operation was cancelled")
    }
}
```

### Error Handling

The package defines several error types:

- `ErrNotFound`: Secret not found in keyring
- `ErrSetDataTooBig`: Data too large for the platform
- `ErrUnsupportedPlatform`: Platform not supported

```go
secret, err := keyring.Get(ctx, "service", "user")
if err == keyring.ErrNotFound {
    fmt.Println("Secret not found")
} else if err != nil {
    log.Fatal(err)
}
```

## Platform-Specific Notes

### macOS

- Uses the macOS Keychain
- Automatically handles encoding for multi-line and non-ASCII passwords
- Limited to ~3000 bytes total (service + username + password)

### Linux/Unix

- Uses the Secret Service API (D-Bus)
- Requires a D-Bus session
- No theoretical size limit, but performance degrades with large values (>100KiB)

### Windows

- Uses Windows Credential Manager
- Password limited to 2560 bytes
- Service name limited to 32KiB (practical limit: 30KiB)

## Testing

Run the tests with:

```bash
go test ./keyring/...
```

The tests will automatically detect the platform and run appropriate tests. Note that some tests may require user interaction (e.g., unlocking keychains).

## API Reference

### Functions

- `Set(ctx context.Context, service, user, password string) error` - Store a secret
- `Get(ctx context.Context, service, user string) (string, error)` - Retrieve a secret
- `Delete(ctx context.Context, service, user string) error` - Delete a secret
- `Purge(ctx context.Context, service string) error` - Purge all secrets for a service (best-effort)

### Error Types

- `ErrNotFound` - Secret not found in keyring
- `ErrSetDataTooBig` - Data too large for the platform
- `ErrUnsupportedPlatform` - Platform not supported

### Important Notes

- **Purge**: Since the underlying zalando/go-keyring doesn't support bulk deletion, `Purge` attempts to delete common user patterns (`user`, `admin`, `root`, `default`, `test`, `demo`) for the given service. This is a best-effort operation and may not delete all secrets.

## Dependencies

This package depends on [github.com/zalando/go-keyring](https://github.com/zalando/go-keyring) which provides the actual keyring implementation.

## License

This package is part of the tykctl-go project and follows the same license terms.