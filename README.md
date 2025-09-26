# Tykctl-Go

A comprehensive Go SDK for building CLI applications with Tyk. This library provides a rich set of tools and utilities for creating powerful command-line interfaces.

## Features

- **API Client**: HTTP client with retry logic, middleware support, and functional options
- **Command Framework**: Structured command handling with Cobra integration
- **File System Utilities**: File watching, operations, and management
- **Interactive Components**: Prompts, progress indicators, and terminal UI
- **Extension System**: Plugin architecture for extending functionality
- **Hook System**: Event-driven hooks with Rego policy support
- **Logging**: Structured logging with Zap integration
- **Browser Integration**: Web browser automation and interaction
- **Script Execution**: Dynamic script running capabilities
- **Table Rendering**: Beautiful table formatting and display
- **Version Management**: Semantic versioning utilities

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/api"
)

func main() {
    // Create API client with functional options
    client := api.New(
        api.WithBaseURL("https://api.example.com/v1"),
        api.WithClientTimeout(30*time.Second),
        api.WithClientHeader("Accept", "application/json"),
    )
    
    ctx := context.Background()
    
    // Make a request
    resp, err := client.Get(ctx, "/users",
        api.WithHeader("X-API-Version", "v1"),
        api.WithQuery("page", "1"),
    )
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Status: %d\n", resp.StatusCode)
    fmt.Printf("Response: %s\n", resp.String())
}
```

## Installation

```bash
go get github.com/edsonmichaque/tykctl-go
```

## Packages

### API Client (`api/`)
High-level HTTP client with functional options, retry logic, and middleware support.

```go
import "github.com/edsonmichaque/tykctl-go/api"

client := api.New(
    api.WithBaseURL("https://api.example.com"),
    api.WithClientTimeout(30*time.Second),
)
```

### Command Framework (`command/`)
Structured command handling with Cobra integration.

```go
import "github.com/edsonmichaque/tykctl-go/command"

cmd := command.New("myapp")
cmd.AddCommand(subCmd)
```

### File System (`fs/`)
File operations, watching, and management utilities.

```go
import "github.com/edsonmichaque/tykctl-go/fs"

watcher := fs.NewWatcher()
watcher.Watch("/path/to/watch")
```

### Interactive Components (`prompt/`, `progress/`, `terminal/`)
User interaction components for CLI applications.

```go
import "github.com/edsonmichaque/tykctl-go/prompt"

answer, err := prompt.Ask("What's your name?")
```

### Extension System (`extension/`)
Plugin architecture for extending functionality.

```go
import "github.com/edsonmichaque/tykctl-go/extension"

runner := extension.NewRunner()
runner.RegisterPlugin("my-plugin", pluginFunc)
```

### Hook System (`hook/`)
Event-driven hooks with Rego policy support.

```go
import "github.com/edsonmichaque/tykctl-go/hook"

hook.Register("before-save", func(ctx *hook.Context) error {
    // Custom logic
    return nil
})
```

### Logging (`logger/`)
Structured logging with Zap integration.

```go
import "github.com/edsonmichaque/tykctl-go/logger"

log := logger.New()
log.Info("Application started")
```

### Browser Integration (`browser/`)
Web browser automation and interaction.

```go
import "github.com/edsonmichaque/tykctl-go/browser"

browser := browser.New()
browser.Open("https://example.com")
```

### Script Execution (`script/`)
Dynamic script running capabilities.

```go
import "github.com/edsonmichaque/tykctl-go/script"

runner := script.NewRunner()
result, err := runner.Execute("console.log('Hello World')")
```

### Table Rendering (`table/`)
Beautiful table formatting and display.

```go
import "github.com/edsonmichaque/tykctl-go/table"

tbl := table.New()
tbl.AddRow("Name", "Age", "City")
tbl.Render()
```

### Version Management (`version/`)
Semantic versioning utilities.

```go
import "github.com/edsonmichaque/tykctl-go/version"

v := version.New("1.2.3")
fmt.Println(v.String()) // "1.2.3"
```

## Examples

See the `example/` directory for comprehensive usage examples.

## Documentation

Detailed documentation is available in the `docs.md` file.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions, please open an issue on GitHub.