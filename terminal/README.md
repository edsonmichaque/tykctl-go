# Terminal Package

The `terminal` package provides terminal capability detection and management for tykctl applications, offering information about terminal dimensions, color support, and TTY detection.

## Features

- **TTY Detection**: Detect if output is going to a terminal
- **Terminal Dimensions**: Get terminal width and height
- **Color Support**: Detect color support and disable colors when needed
- **Environment Integration**: Respect environment variables for terminal settings
- **Cross-platform**: Works consistently across different operating systems
- **Force TTY**: Option to force TTY behavior for testing

## Usage

### Basic Terminal Detection

```go
package main

import (
    "fmt"
    
    "github.com/edsonmichaque/tykctl-go/terminal"
)

func main() {
    // Create a new terminal instance
    term := terminal.New()
    
    // Check if output is a TTY
    if term.IsTTY() {
        fmt.Println("Output is going to a terminal")
    } else {
        fmt.Println("Output is not going to a terminal")
    }
    
    // Get terminal dimensions
    fmt.Printf("Terminal size: %dx%d\n", term.Width, term.Height)
    
    // Check color support
    if term.Color {
        fmt.Println("Terminal supports colors")
    } else {
        fmt.Println("Terminal does not support colors")
    }
}
```

### Terminal-Aware Output

```go
func terminalAwareOutput() {
    term := terminal.New()
    
    if term.IsTTY() {
        // Terminal output - use colors and formatting
        if term.Color {
            fmt.Println("\033[32m✓ Success\033[0m")
            fmt.Println("\033[31m✗ Error\033[0m")
        } else {
            fmt.Println("✓ Success")
            fmt.Println("✗ Error")
        }
        
        // Use terminal width for formatting
        fmt.Printf("Terminal width: %d characters\n", term.Width)
    } else {
        // Non-terminal output - plain text
        fmt.Println("Success")
        fmt.Println("Error")
    }
}
```

### Color Management

```go
func colorManagement() {
    term := terminal.New()
    
    // Check color support
    if term.Color && !term.NoColor {
        // Use colors
        fmt.Println("\033[1;34mBold Blue Text\033[0m")
        fmt.Println("\033[32mGreen Text\033[0m")
        fmt.Println("\033[31mRed Text\033[0m")
    } else {
        // No colors
        fmt.Println("Bold Blue Text")
        fmt.Println("Green Text")
        fmt.Println("Red Text")
    }
}
```

## Advanced Usage

### Dynamic Terminal Detection

```go
func dynamicTerminalDetection() {
    term := terminal.New()
    
    // Check various terminal conditions
    if term.IsTTY() {
        fmt.Println("Running in terminal")
        
        if term.Width < 80 {
            fmt.Println("Terminal is narrow, using compact output")
        } else {
            fmt.Println("Terminal is wide, using full output")
        }
        
        if term.Height < 24 {
            fmt.Println("Terminal is short, limiting output")
        }
    } else {
        fmt.Println("Not running in terminal, using plain output")
    }
}
```

### Terminal Configuration

```go
func terminalConfiguration() {
    term := terminal.New()
    
    // Check environment variables
    fmt.Printf("TTY: %t\n", term.IsTTY())
    fmt.Printf("Width: %d\n", term.Width)
    fmt.Printf("Height: %d\n", term.Height)
    fmt.Printf("Color: %t\n", term.Color)
    fmt.Printf("No Color: %t\n", term.NoColor)
    fmt.Printf("Force TTY: %t\n", term.ForceTTY)
    
    // Use configuration for output
    if term.IsTTY() && term.Color && !term.NoColor {
        fmt.Println("Using colored terminal output")
    } else {
        fmt.Println("Using plain text output")
    }
}
```

### Responsive Layout

```go
func responsiveLayout() {
    term := terminal.New()
    
    // Adjust layout based on terminal size
    if term.Width >= 120 {
        // Wide terminal - use multi-column layout
        fmt.Println("Using wide layout")
        displayWideLayout()
    } else if term.Width >= 80 {
        // Medium terminal - use standard layout
        fmt.Println("Using standard layout")
        displayStandardLayout()
    } else {
        // Narrow terminal - use compact layout
        fmt.Println("Using compact layout")
        displayCompactLayout()
    }
}

func displayWideLayout() {
    fmt.Println("Column 1    Column 2    Column 3    Column 4")
    fmt.Println("Data 1      Data 2      Data 3      Data 4")
}

func displayStandardLayout() {
    fmt.Println("Column 1    Column 2    Column 3")
    fmt.Println("Data 1      Data 2      Data 3")
}

func displayCompactLayout() {
    fmt.Println("Column 1    Column 2")
    fmt.Println("Data 1      Data 2")
}
```

## Integration Examples

### With Progress Indicators

```go
import (
    "github.com/edsonmichaque/tykctl-go/progress"
    "github.com/edsonmichaque/tykctl-go/terminal"
)

func progressWithTerminal() {
    term := terminal.New()
    
    if term.IsTTY() {
        // Use animated progress in terminal
        spinner := progress.New()
        spinner.SetMessage("Processing...")
        spinner.Start()
        defer spinner.Stop()
        
        // Simulate work
        time.Sleep(3 * time.Second)
    } else {
        // Use simple progress for non-terminal
        fmt.Println("Processing...")
        time.Sleep(3 * time.Second)
        fmt.Println("Done")
    }
}
```

### With Table Output

```go
func tableWithTerminal() {
    term := terminal.New()
    
    // Create table
    t := table.New()
    t.SetHeaders([]string{"Name", "Status", "Time"})
    t.AddRow([]string{"Task 1", "Completed", "10:30"})
    t.AddRow([]string{"Task 2", "Running", "10:31"})
    
    if term.IsTTY() {
        // Use full table formatting in terminal
        t.Render()
    } else {
        // Use simple format for non-terminal
        fmt.Println("Name    Status     Time")
        fmt.Println("Task 1  Completed  10:30")
        fmt.Println("Task 2  Running    10:31")
    }
}
```

### With Logging

```go
func loggingWithTerminal() {
    term := terminal.New()
    
    // Configure logger based on terminal
    config := logger.Config{
        Debug:   true,
        Verbose: true,
        NoColor: term.NoColor || !term.Color,
    }
    
    zapLogger := logger.New(config)
    
    // Log with terminal awareness
    if term.IsTTY() {
        zapLogger.Info("Running in terminal with colors")
    } else {
        zapLogger.Info("Running in non-terminal environment")
    }
}
```

## Environment Variables

### Supported Environment Variables

- `COLUMNS` - Terminal width
- `LINES` - Terminal height
- `NO_COLOR` - Disable colors
- `FORCE_TTY` - Force TTY behavior

### Environment Integration

```go
func environmentIntegration() {
    term := terminal.New()
    
    // Check environment variables
    if os.Getenv("NO_COLOR") != "" {
        fmt.Println("Colors disabled by NO_COLOR environment variable")
    }
    
    if os.Getenv("COLUMNS") != "" {
        fmt.Printf("Terminal width set by COLUMNS: %d\n", term.Width)
    }
    
    if os.Getenv("FORCE_TTY") != "" {
        fmt.Println("TTY behavior forced by FORCE_TTY environment variable")
    }
}
```

## Testing and Development

### Force TTY for Testing

```go
func testingWithForceTTY() {
    // Set environment variable for testing
    os.Setenv("FORCE_TTY", "1")
    
    term := terminal.New()
    
    if term.ForceTTY {
        fmt.Println("TTY behavior forced for testing")
    }
    
    // Test terminal-aware code
    if term.IsTTY() {
        fmt.Println("Testing terminal output")
    }
}
```

### Terminal Simulation

```go
func simulateTerminal() {
    // Simulate different terminal conditions
    testCases := []struct {
        name   string
        width  int
        height int
        color  bool
        noColor bool
    }{
        {"Wide Color Terminal", 120, 40, true, false},
        {"Narrow Terminal", 60, 20, true, false},
        {"No Color Terminal", 80, 24, false, true},
        {"Small Terminal", 40, 10, true, false},
    }
    
    for _, tc := range testCases {
        fmt.Printf("Testing %s:\n", tc.name)
        
        // Simulate terminal conditions
        os.Setenv("COLUMNS", fmt.Sprintf("%d", tc.width))
        os.Setenv("LINES", fmt.Sprintf("%d", tc.height))
        if tc.noColor {
            os.Setenv("NO_COLOR", "1")
        } else {
            os.Unsetenv("NO_COLOR")
        }
        
        term := terminal.New()
        fmt.Printf("  Width: %d, Height: %d, Color: %t\n", 
            term.Width, term.Height, term.Color)
    }
}
```

## Best Practices

- **TTY Detection**: Always check if output is going to a terminal
- **Color Management**: Respect NO_COLOR environment variable
- **Responsive Design**: Adjust output based on terminal dimensions
- **Fallback Behavior**: Provide plain text fallbacks for non-terminal output
- **Environment Respect**: Honor terminal-related environment variables
- **Testing**: Test with different terminal conditions

## Dependencies

- No external dependencies
- Uses only Go standard library (`os`, `strconv`, `strings`)