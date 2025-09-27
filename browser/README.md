# Browser Package

The `browser` package provides cross-platform browser functionality for opening URLs in the default system browser.

## Features

- **Cross-platform Support**: Works on Windows, macOS, and Linux
- **Default Browser Detection**: Automatically detects and uses the system's default browser
- **Multiple Browser Fallback**: On Linux, tries multiple common browsers if the default isn't available
- **Simple Interface**: Clean, easy-to-use API for opening URLs

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/browser"
)

func main() {
    // Create a new browser instance
    b := browser.New()
    
    // Open a URL in the default browser
    err := b.Open("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("URL opened successfully!")
}
```

### Interface Usage

```go
// Use the Browser interface for dependency injection
func openURL(b browser.Browser, url string) error {
    return b.Open(url)
}

// Create browser instance
b := browser.New()

// Use with dependency injection
err := openURL(b, "https://github.com")
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

The package returns descriptive errors for common failure scenarios:

```go
err := browser.Open("https://example.com")
if err != nil {
    switch {
    case strings.Contains(err.Error(), "no suitable browser found"):
        fmt.Println("No browser found on this system")
    case strings.Contains(err.Error(), "unsupported platform"):
        fmt.Println("Platform not supported")
    default:
        fmt.Printf("Failed to open browser: %v\n", err)
    }
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

- **Error Handling**: Always handle errors gracefully, as browser availability can vary
- **User Feedback**: Provide clear feedback when opening URLs
- **Fallback Options**: Consider providing manual instructions if browser opening fails
- **Testing**: Test on different platforms to ensure compatibility
- **Dependency Injection**: Use the `Browser` interface for better testability

## Dependencies

- No external dependencies
- Uses only Go standard library (`os/exec`, `runtime`, `fmt`)

## Thread Safety

The `DefaultBrowser` implementation is thread-safe and can be used concurrently from multiple goroutines.