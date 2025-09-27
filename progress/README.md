# Progress Package

The `progress` package provides visual progress indicators including spinners and progress bars for tykctl applications, built with the Bubble Tea framework for rich terminal UI experiences.

## Features

- **Spinner Support**: Animated spinners for indeterminate progress
- **Progress Bars**: Visual progress bars for determinate progress
- **Context Support**: Full context.Context integration for cancellation
- **Thread Safety**: Safe for concurrent use across goroutines
- **Customizable**: Configurable messages, characters, and styling
- **Terminal UI**: Rich terminal UI using Bubble Tea framework

## Usage

### Basic Spinner

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/progress"
)

func main() {
    // Create a new spinner
    spinner := progress.New()
    
    // Set message
    spinner.SetMessage("Loading...")
    
    // Start spinner
    spinner.Start()
    
    // Simulate work
    time.Sleep(3 * time.Second)
    
    // Stop spinner
    spinner.Stop()
    
    fmt.Println("Task completed!")
}
```

### Spinner with Context

```go
func spinnerWithContext() error {
    spinner := progress.New()
    spinner.SetMessage("Processing data...")
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Start spinner
    spinner.Start()
    defer spinner.Stop()
    
    // Simulate work with context
    select {
    case <-time.After(10 * time.Second):
        return fmt.Errorf("operation timed out")
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Progress Bar

```go
func progressBarExample() {
    // Create progress bar with total
    total := int64(100)
    bar := progress.NewBar(total)
    
    // Set message
    bar.SetMessage("Downloading files...")
    
    // Start progress bar
    bar.Start()
    defer bar.Stop()
    
    // Simulate progress
    for i := int64(0); i <= total; i++ {
        bar.SetProgress(i)
        time.Sleep(50 * time.Millisecond)
    }
    
    fmt.Println("Download completed!")
}
```

## Advanced Usage

### Custom Spinner Configuration

```go
func customSpinner() {
    spinner := progress.New()
    
    // Custom configuration
    spinner.SetMessage("Custom operation...")
    spinner.SetFrames([]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"})
    
    spinner.Start()
    defer spinner.Stop()
    
    // Simulate work
    time.Sleep(3 * time.Second)
}
```

### Custom Progress Bar

```go
func customProgressBar() {
    total := int64(50)
    bar := progress.NewBar(total)
    
    // Custom configuration
    bar.SetMessage("Processing items...")
    bar.SetWidth(30)
    bar.SetFillChar("█")
    bar.SetEmptyChar("░")
    
    bar.Start()
    defer bar.Stop()
    
    // Simulate progress
    for i := int64(0); i <= total; i++ {
        bar.SetProgress(i)
        time.Sleep(100 * time.Millisecond)
    }
}
```

### Multiple Progress Indicators

```go
func multipleProgress() {
    // Create multiple spinners
    spinner1 := progress.New()
    spinner1.SetMessage("Task 1...")
    
    spinner2 := progress.New()
    spinner2.SetMessage("Task 2...")
    
    // Start both spinners
    spinner1.Start()
    spinner2.Start()
    
    // Simulate concurrent work
    go func() {
        time.Sleep(2 * time.Second)
        spinner1.Stop()
        fmt.Println("Task 1 completed!")
    }()
    
    go func() {
        time.Sleep(3 * time.Second)
        spinner2.Stop()
        fmt.Println("Task 2 completed!")
    }()
    
    // Wait for both to complete
    time.Sleep(4 * time.Second)
}
```

## Integration Examples

### With HTTP Client

```go
import (
    "github.com/edsonmichaque/tykctl-go/httpclient"
    "github.com/edsonmichaque/tykctl-go/progress"
)

func httpClientWithProgress() error {
    spinner := progress.New()
    spinner.SetMessage("Fetching data...")
    
    spinner.Start()
    defer spinner.Stop()
    
    client := httpclient.NewWithBaseURL("https://api.example.com")
    ctx := context.Background()
    
    resp, err := client.Get(ctx, "/data")
    if err != nil {
        return err
    }
    
    if resp.IsSuccess() {
        fmt.Println("Data fetched successfully!")
    }
    
    return nil
}
```

### With File Operations

```go
func fileOperationsWithProgress() error {
    // Get file size for progress bar
    fileInfo, err := os.Stat("large-file.txt")
    if err != nil {
        return err
    }
    
    total := fileInfo.Size()
    bar := progress.NewBar(total)
    bar.SetMessage("Copying file...")
    
    bar.Start()
    defer bar.Stop()
    
    // Copy file with progress updates
    src, err := os.Open("large-file.txt")
    if err != nil {
        return err
    }
    defer src.Close()
    
    dst, err := os.Create("copy-of-large-file.txt")
    if err != nil {
        return err
    }
    defer dst.Close()
    
    // Copy with progress updates
    buffer := make([]byte, 32*1024) // 32KB buffer
    var copied int64
    
    for {
        n, err := src.Read(buffer)
        if err != nil && err != io.EOF {
            return err
        }
        
        if n == 0 {
            break
        }
        
        _, err = dst.Write(buffer[:n])
        if err != nil {
            return err
        }
        
        copied += int64(n)
        bar.SetProgress(copied)
    }
    
    return nil
}
```

### With Extension Installation

```go
func extensionInstallationWithProgress() error {
    spinner := progress.New()
    spinner.SetMessage("Installing extension...")
    
    spinner.Start()
    defer spinner.Stop()
    
    installer := extension.NewInstaller("/tmp/tykctl-config")
    ctx := context.Background()
    
    err := installer.InstallExtension(ctx, "owner", "repo")
    if err != nil {
        spinner.SetMessage("Installation failed")
        return err
    }
    
    spinner.SetMessage("Installation completed")
    return nil
}
```

## Context-Aware Progress

### With Timeout

```go
func progressWithTimeout() error {
    spinner := progress.New()
    spinner.SetMessage("Long running operation...")
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    spinner.Start()
    defer spinner.Stop()
    
    // Simulate long operation
    select {
    case <-time.After(15 * time.Second):
        return fmt.Errorf("operation timed out")
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### With Cancellation

```go
func progressWithCancellation() error {
    spinner := progress.New()
    spinner.SetMessage("Processing...")
    
    // Create cancellable context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    spinner.Start()
    defer spinner.Stop()
    
    // Simulate work that can be cancelled
    go func() {
        time.Sleep(5 * time.Second)
        cancel() // Cancel after 5 seconds
    }()
    
    // Wait for cancellation or completion
    <-ctx.Done()
    return ctx.Err()
}
```

## Customization Options

### Spinner Customization

```go
func customizeSpinner() {
    spinner := progress.New()
    
    // Custom frames
    spinner.SetFrames([]string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"})
    
    // Custom message
    spinner.SetMessage("Custom spinner message")
    
    spinner.Start()
    defer spinner.Stop()
    
    time.Sleep(3 * time.Second)
}
```

### Progress Bar Customization

```go
func customizeProgressBar() {
    total := int64(100)
    bar := progress.NewBar(total)
    
    // Custom width
    bar.SetWidth(40)
    
    // Custom characters
    bar.SetFillChar("▓")
    bar.SetEmptyChar("░")
    
    // Custom message
    bar.SetMessage("Custom progress bar")
    
    bar.Start()
    defer bar.Stop()
    
    // Simulate progress
    for i := int64(0); i <= total; i++ {
        bar.SetProgress(i)
        time.Sleep(50 * time.Millisecond)
    }
}
```

## Error Handling

### Progress with Error Recovery

```go
func progressWithErrorRecovery() error {
    spinner := progress.New()
    spinner.SetMessage("Attempting operation...")
    
    maxRetries := 3
    for attempt := 1; attempt <= maxRetries; attempt++ {
        spinner.SetMessage(fmt.Sprintf("Attempt %d/%d...", attempt, maxRetries))
        
        spinner.Start()
        
        // Simulate operation that might fail
        err := simulateOperation()
        
        spinner.Stop()
        
        if err == nil {
            fmt.Println("Operation succeeded!")
            return nil
        }
        
        if attempt < maxRetries {
            fmt.Printf("Attempt %d failed, retrying...\n", attempt)
            time.Sleep(time.Second)
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}

func simulateOperation() error {
    // Simulate random failure
    if time.Now().UnixNano()%3 == 0 {
        return fmt.Errorf("simulated error")
    }
    time.Sleep(time.Second)
    return nil
}
```

## Best Practices

- **Context Usage**: Always use context for cancellation and timeout handling
- **Resource Cleanup**: Always stop progress indicators when done
- **Thread Safety**: Progress indicators are thread-safe and can be used concurrently
- **User Feedback**: Provide meaningful messages for progress indicators
- **Error Handling**: Handle errors gracefully and update progress messages accordingly
- **Performance**: Avoid updating progress too frequently to prevent UI flickering

## Dependencies

- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- `github.com/charmbracelet/lipgloss` - Styling library