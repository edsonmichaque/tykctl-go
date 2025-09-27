# Browser Package

The `browser` package provides comprehensive cross-platform browser functionality with advanced features for opening URLs, managing browser configurations, and handling various browser scenarios.

## Features

- **Cross-platform Support**: Works on Windows, macOS, and Linux
- **Default Browser Detection**: Automatically detects and uses the system's default browser
- **Multiple Browser Support**: Support for Chrome, Firefox, Safari, Edge, and more
- **Advanced Configuration**: Custom timeouts, browser arguments, and behavior options
- **Context Support**: Full context.Context integration for cancellation and timeouts
- **Background Operations**: Open URLs in background without blocking
- **New Window/Tab Support**: Control how URLs are opened (new window, new tab)
- **URL Validation**: Built-in URL format validation
- **Browser Discovery**: Find available browsers on the system
- **Comprehensive Error Handling**: Detailed error types and handling

## Quick Start

### Simplest Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/browser"
)

func main() {
    // Simplest way to open a URL
    err := browser.OpenURL("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("URL opened successfully!")
}
```

### Using the Browser Interface

```go
func main() {
    // Create a new browser instance
    b := browser.New()
    
    // Check if browser is available
    if b.IsAvailable() {
        fmt.Printf("Using browser: %s\n", b.GetName())
        
        // Open a URL
        err := b.Open("https://example.com")
        if err != nil {
            log.Fatal(err)
        }
    } else {
        fmt.Println("No browser available")
    }
}
```

## Advanced Usage

### Convenience Functions

```go
// Open with timeout
err := browser.OpenURLWithTimeout("https://example.com", 10*time.Second)

// Open in background (non-blocking)
err := browser.OpenURLInBackground("https://example.com")

// Open in new window
err := browser.OpenURLInNewWindow("https://example.com")

// Open in new tab
err := browser.OpenURLInNewTab("https://example.com")

// Open with specific browser
err := browser.OpenWithBrowser("firefox", "https://example.com")

// Open with specific browser and timeout
err := browser.OpenWithBrowserAndTimeout("chrome", "https://example.com", 15*time.Second)
```

### Custom Configuration

```go
// Create browser with custom configuration
config := browser.Config{
    BrowserName: "chrome",
    Timeout:     15 * time.Second,
    NewTab:      true,
    Args:        []string{"--incognito", "--disable-web-security"},
}

b := browser.NewWithConfig(config)
err := b.Open("https://example.com")
```

### Context Support

```go
// Open with context for cancellation and timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

b := browser.New()
err := b.OpenWithContext(ctx, "https://example.com")
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("Browser operation timed out")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Browser Discovery

```go
// Get default browser
defaultBrowser := browser.GetDefaultBrowser()
fmt.Printf("Default browser: %s\n", defaultBrowser)

// Get all available browsers
availableBrowsers := browser.GetAvailableBrowsers()
fmt.Printf("Available browsers: %v\n", availableBrowsers)

// Get detailed browser information
browserInfo := browser.GetBrowserInfo()
for _, info := range browserInfo {
    defaultMark := ""
    if info.Default {
        defaultMark = " (default)"
    }
    fmt.Printf("%s: available=%t%s\n", info.Name, info.Available, defaultMark)
}
```

### URL Validation

```go
// Validate URL before opening
url := "https://example.com"
err := browser.ValidateURL(url)
if err != nil {
    fmt.Printf("Invalid URL: %v\n", err)
    return
}

// Open validated URL
err = browser.OpenURL(url)
if err != nil {
    fmt.Printf("Error opening URL: %v\n", err)
}
```

## Platform Support

### Windows
- Uses `rundll32 url.dll,FileProtocolHandler` to open URLs
- Works with any default browser set in Windows

### macOS
- Uses the `open` command to open URLs
- Respects the system's default browser setting

### Linux
- Tries multiple browsers in order of preference:
  1. `xdg-open` (most common, respects desktop environment settings)
  2. `firefox`
  3. `google-chrome`
  4. `chromium-browser`
- Returns an error if no suitable browser is found

## Error Handling

The package provides comprehensive error handling with specific error types:

```go
err := browser.OpenURL("https://example.com")
if err != nil {
    // Handle different error types
    switch {
    case strings.Contains(err.Error(), "URL cannot be empty"):
        fmt.Println("Empty URL provided")
    case strings.Contains(err.Error(), "invalid URL format"):
        fmt.Println("Invalid URL format")
    case strings.Contains(err.Error(), "no suitable browser found"):
        fmt.Println("No browser available on this system")
    case strings.Contains(err.Error(), "browser") && strings.Contains(err.Error(), "is not available"):
        fmt.Println("Specified browser is not available")
    case strings.Contains(err.Error(), "unsupported platform"):
        fmt.Println("Platform not supported")
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

### Error Types

```go
// BrowserError - for browser-specific errors
type BrowserError struct {
    Browser string
    URL     string
    Err     error
}

// URLValidationError - for URL validation errors
type URLValidationError struct {
    URL string
    Err error
}
```

### Graceful Error Handling

```go
func openURLSafely(url string) error {
    // Validate URL first
    if err := browser.ValidateURL(url); err != nil {
        return fmt.Errorf("invalid URL: %w", err)
    }
    
    // Try to open with fallback
    err := browser.OpenURL(url)
    if err != nil {
        // Provide fallback instructions
        fmt.Printf("Failed to open browser: %v\n", err)
        fmt.Printf("Please manually open: %s\n", url)
        return nil // Don't fail the operation
    }
    
    return nil
}
```

## Use Cases

- **CLI Tools**: Open documentation, help pages, or web interfaces from command-line tools
- **Desktop Applications**: Open external links from desktop applications
- **Development Tools**: Open project URLs, documentation, or issue trackers
- **Web Applications**: Open external links from web-based tools
- **Automation Scripts**: Open URLs as part of automated workflows

## Integration Examples

### With CLI Commands

```go
import (
    "github.com/spf13/cobra"
    "github.com/edsonmichaque/tykctl-go/browser"
)

func openDocsCommand(cmd *cobra.Command, args []string) error {
    b := browser.New()
    return b.Open("https://docs.example.com")
}

var docsCmd = &cobra.Command{
    Use:   "docs",
    Short: "Open documentation in browser",
    RunE:  openDocsCommand,
}
```

### With User Confirmation

```go
func openWithConfirmation(b browser.Browser, url string) error {
    fmt.Printf("Opening %s in your browser...\n", url)
    
    err := b.Open(url)
    if err != nil {
        return fmt.Errorf("failed to open browser: %w", err)
    }
    
    fmt.Println("Browser opened successfully!")
    return nil
}
```

### With Error Recovery

```go
func openWithFallback(b browser.Browser, url string) error {
    err := b.Open(url)
    if err != nil {
        fmt.Printf("Failed to open browser: %v\n", err)
        fmt.Printf("Please manually open: %s\n", url)
        return nil // Don't fail the operation
    }
    return nil
}
```

## Best Practices

### Error Handling
- **Always validate URLs** before opening them using `ValidateURL()`
- **Handle browser unavailability** gracefully with fallback options
- **Use context for timeouts** to prevent hanging operations
- **Provide user feedback** when browser operations fail

### Performance
- **Use background opening** for non-blocking operations with `OpenInBackground()`
- **Set appropriate timeouts** to prevent long waits
- **Cache browser instances** when making multiple calls

### User Experience
- **Check browser availability** before attempting to open URLs
- **Provide fallback instructions** when browser opening fails
- **Use appropriate opening modes** (new window vs new tab)

### Testing
- **Test on different platforms** to ensure compatibility
- **Mock the Browser interface** for unit tests
- **Test error scenarios** including browser unavailability

## Examples

See `example.go` for comprehensive usage examples including:
- Basic and advanced usage patterns
- Error handling strategies
- Browser discovery and configuration
- Context-aware operations
- Batch URL operations

## Dependencies

- No external dependencies
- Uses only Go standard library (`os/exec`, `runtime`, `fmt`, `context`, `strings`, `time`)

## Thread Safety

The `DefaultBrowser` implementation is thread-safe and can be used concurrently from multiple goroutines.