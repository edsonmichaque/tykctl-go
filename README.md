# Tykctl-Go

A comprehensive Go library for creating Tyk CLI extensions with best practices, templates, and utilities. This library provides a framework for building consistent, maintainable, and feature-rich CLI extensions for the Tyk ecosystem.

## Features

- **Extension Management**: Install, manage, and run Tyk CLI extensions
- **API Client**: HTTP client with retry logic, middleware support, and functional options
- **Command Framework**: Structured command handling with Cobra integration
- **File System Utilities**: File watching, operations, and management
- **Interactive Components**: Prompts, progress indicators, and terminal UI
- **Hook System**: Event-driven hooks with Rego policy support
- **Logging**: Structured logging with Zap integration
- **Browser Integration**: Web browser automation and interaction
- **Script Execution**: Dynamic script running capabilities
- **Table Rendering**: Beautiful table formatting and display
- **JSON Processing**: Pure Go JQ integration for JSON manipulation
- **JSON Schema**: JSON Schema validation using external gojsonschema library
- **HTTP Client**: Simple HTTP client for API interactions
- **Editor Integration**: File editing with external editors
- **Version Management**: Semantic versioning utilities

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/TykTechnologies/tykctl-go/extension"
)

func main() {
    // Create extension installer
    configDir := "/tmp/tykctl-config"
    installer := extension.NewInstaller(configDir)
    
    ctx := context.Background()
    
    // Search for extensions
    extensions, err := installer.SearchExtensions(ctx, "tyk", 10)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Found %d extensions:\n", len(extensions))
    for _, ext := range extensions {
        fmt.Printf("- %s: %s (%d stars)\n", ext.Name, ext.Description, ext.Stars)
    }
}
```

## Installation

```bash
go get github.com/TykTechnologies/tykctl-go
```

## Packages

### Extension Management (`extension/`)
Core extension management functionality for installing, managing, and running Tyk CLI extensions.

```go
import "github.com/TykTechnologies/tykctl-go/extension"

installer := extension.NewInstaller("/config/dir")
extensions, err := installer.SearchExtensions(ctx, "tyk", 10)
```

### API Client (`api/`)
High-level HTTP client with functional options, retry logic, and middleware support.

```go
import "github.com/TykTechnologies/tykctl-go/api"

client := api.New(
    api.WithBaseURL("https://api.example.com"),
    api.WithClientTimeout(30*time.Second),
)
```

### HTTP Client (`httpclient/`)
Simple HTTP client for making API requests with context support.

```go
import "github.com/TykTechnologies/tykctl-go/httpclient"

client := httpclient.NewWithBaseURL("https://api.example.com")
data, err := client.Get("/users")
```

### Command Framework (`command/`)
Structured command handling with Cobra integration.

```go
import "github.com/TykTechnologies/tykctl-go/command"

cmd := command.New("myapp", "Short description", handler)
cmd.SetLogger(logger)
```

### File System (`fs/`)
File operations, watching, and management utilities.

```go
import "github.com/TykTechnologies/tykctl-go/fs"

watcher := fs.NewWatcher()
watcher.Watch("/path/to/watch")
```

### Interactive Components (`prompt/`, `progress/`, `terminal/`)
User interaction components for CLI applications.

```go
import "github.com/TykTechnologies/tykctl-go/prompt"

answer, err := prompt.Ask("What's your name?")
```

### Hook System (`hook/`)
Event-driven hooks with Rego policy support.

```go
import "github.com/TykTechnologies/tykctl-go/hook"

hm := hook.New()
hm.Register(ctx, "before-save", func(ctx context.Context, data interface{}) error {
    // Custom logic
    return nil
})
```

### Logging (`logger/`)
Structured logging with Zap integration.

```go
import "github.com/TykTechnologies/tykctl-go/logger"

log := logger.New()
log.Info("Application started")
```

### Browser Integration (`browser/`)
Web browser automation and interaction.

```go
import "github.com/TykTechnologies/tykctl-go/browser"

browser := browser.New()
browser.Open("https://example.com")
```

### Script Execution (`script/`)
Dynamic script running capabilities.

```go
import "github.com/TykTechnologies/tykctl-go/script"

sm := script.NewScriptManager("/scripts")
s, err := sm.CreateScript("my-script", "Description", "#!/bin/bash\necho 'Hello'")
```

### Table Rendering (`table/`)
Beautiful table formatting and display.

```go
import "github.com/TykTechnologies/tykctl-go/table"

tbl := table.New()
tbl.SetHeaders([]string{"Name", "Age", "City"})
tbl.AddRow([]string{"John", "30", "NYC"})
tbl.Render()
```

### JSON Processing (`jq/`)
Pure Go JQ integration for JSON manipulation and processing using gojq.

```go
import "github.com/TykTechnologies/tykctl-go/jq"

// Process JSON string
result, err := jq.ProcessString(jsonData, ".users[0].name")

// Process JSON bytes
data, err := jq.Process(jsonBytes, ".field")

// Process Go object
obj, err := jq.ProcessObject(myObject, ".property")

// Complex queries
activeUsers, err := jq.ProcessString(jsonData, ".users[] | select(.active) | .name")
```

### JSON Schema (`jsonschema/`)
JSON Schema validation using the gojsonschema library with context support.

```go
import (
    "context"
    "time"
    "github.com/TykTechnologies/tykctl-go/jsonschema"
)

// Create validator from schema string
validator, err := jsonschema.New(schemaString)

// Validate JSON data with context
ctx := context.Background()
result, err := validator.ValidateString(ctx, jsonData)

// Validate with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
result, err := validator.ValidateString(ctx, jsonData)

// Convenience functions
isValid, err := jsonschema.IsValidString(ctx, jsonData, schemaString)

// Validate files
result, err := jsonschema.ValidateFile(ctx, "data.json", "schema.json")

// Validate directories
results, err := jsonschema.ValidateDirectory(ctx, "./data", schemaString)
```

### Editor Integration (`editor/`)
File editing with external editors.

```go
import "github.com/TykTechnologies/tykctl-go/editor"

ed := editor.New()
content, err := ed.EditString(ctx, "Initial content")
```

### Version Management (`version/`)
Semantic versioning utilities.

```go
import "github.com/TykTechnologies/tykctl-go/version"

v := version.New("1.2.3")
fmt.Println(v.String()) // "1.2.3"
```

## Examples

See the `example/` directory for comprehensive usage examples demonstrating extension management, API client usage, and other package functionalities.

## Documentation

Detailed documentation is available in the `docs.md` file, which provides comprehensive information about:
- Extension framework and command creation
- Template system for generating extensions
- Hook system with Rego policy support
- Best practices for CLI extension development

## Key Features

- **Extension Discovery**: Search and discover Tyk CLI extensions from GitHub
- **Extension Management**: Install, uninstall, and manage extensions
- **Hook System**: Event-driven hooks with support for builtin, external, and Rego hooks
- **Template System**: Pre-built templates for common extension patterns
- **Rich CLI Components**: Progress indicators, prompts, tables, and terminal UI
- **JSON Processing**: Pure Go JQ integration for advanced JSON manipulation
- **JSON Schema**: Comprehensive JSON Schema validation with context support and detailed error reporting
- **HTTP Clients**: Both high-level API client and simple HTTP client
- **File Operations**: File watching, editing, and management utilities

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