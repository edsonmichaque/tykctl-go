# Example Package

The `example` package contains a comprehensive demonstration of the tykctl-go extension system, showcasing how to search, install, manage, and run extensions.

## Overview

This example demonstrates the complete extension lifecycle:

1. **Extension Discovery**: Search for available extensions
2. **Extension Management**: List installed and available extensions
3. **Extension Installation**: Install extensions from repositories
4. **Extension Execution**: Run installed extensions

## Features Demonstrated

- Extension searching with GitHub integration
- Extension installation and management
- Extension runner functionality
- Error handling and logging
- Context-aware operations

## Usage

Run the example:

```bash
go run main.go
```

## Code Examples

### Basic Extension Installer

```go
// Create extension installer
configDir := "/tmp/tykctl-config"
installer := extension.NewInstaller(configDir)

// Search for extensions
extensions, err := installer.SearchExtensions(ctx, "tyk", 10)
if err != nil {
    // Handle error
}

// List installed extensions
installed, err := installer.ListInstalledExtensions(ctx)
```

### Extension Runner

```go
// Create extension runner
runner := extension.NewRunner(configDir)

// List available extensions
available, err := runner.ListAvailableExtensions()

// Check if extension is available
if runner.IsExtensionAvailable("example") {
    // Run extension
    err := runner.RunExtension(ctx, "example", []string{"--help"})
}
```

### Advanced Configuration

```go
// Installer with GitHub token
installer := extension.NewInstaller(
    configDir,
    extension.WithGitHubToken("your-token-here"),
)

// Installer with custom logger
customLogger := zap.NewExample()
installer := extension.NewInstaller(
    configDir,
    extension.WithLogger(customLogger),
)
```

## Configuration

The example uses a temporary configuration directory (`/tmp/tykctl-config`). In production, this would typically be:

- `~/.config/tykctl/` on Linux/macOS
- `%APPDATA%/tykctl/` on Windows

## Dependencies

- `github.com/edsonmichaque/tykctl-go/extension`: Extension management
- `go.uber.org/zap`: Structured logging

## Error Handling

The example demonstrates proper error handling:

```go
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

## Safety Features

- Installation and execution examples are commented out to prevent side effects
- Proper error handling and logging
- Context-aware operations with cancellation support

## Extension Lifecycle

1. **Discovery**: Search for extensions using keywords
2. **Installation**: Install extensions from GitHub repositories
3. **Management**: List and manage installed extensions
4. **Execution**: Run extensions with custom arguments

## GitHub Integration

The example shows how to integrate with GitHub for extension discovery:

- Search repositories by keyword
- Filter by stars and other criteria
- Handle GitHub API rate limits
- Support for private repositories with tokens

## Running Extensions

Extensions can be run with custom arguments:

```go
// Run with arguments
err := runner.RunExtension(ctx, "example", []string{"--help", "--verbose"})

// Run with no arguments
err := runner.RunExtension(ctx, "example", []string{})
```

## Logging

The example demonstrates different logging approaches:

- Default logging
- Custom logger configuration
- Structured logging with zap

## Best Practices

- Always use context for cancellation
- Handle errors appropriately
- Use functional options for configuration
- Implement proper logging
- Test extension availability before execution

## Troubleshooting

Common issues and solutions:

1. **Extension not found**: Check if the extension is installed
2. **Permission errors**: Ensure proper file permissions
3. **Network issues**: Check internet connectivity for GitHub access
4. **Configuration errors**: Verify configuration directory exists

## License

This example is part of the tykctl-go project and follows the same license terms.
