# Version Package

The `version` package provides version information management for tykctl applications, offering various version string formats and build information for extensions and applications.

## Features

- **Version Management**: Centralized version information
- **Build Information**: Git commit, build date, and Go version tracking
- **Multiple Formats**: Short, full, and custom version string formats
- **Extension Support**: Version information for specific extensions
- **Runtime Information**: Access to Go runtime version information

## Usage

### Basic Version Information

```go
package main

import (
    "fmt"
    
    "github.com/edsonmichaque/tykctl-go/version"
)

func main() {
    // Get version string
    fmt.Printf("Version: %s\n", version.String())
    
    // Get short version
    fmt.Printf("Short version: %s\n", version.Short())
    
    // Get full version information
    fmt.Printf("Full version:\n%s\n", version.Full())
}
```

### Extension Version Information

```go
func extensionVersionInfo() {
    extensionName := "tykctl-portal"
    
    // Get extension version info
    fmt.Printf("Extension: %s\n", version.Info(extensionName))
    
    // Get full extension version info
    fmt.Printf("Full extension info:\n%s\n", version.InfoFull(extensionName))
}
```

### Version Constants

```go
func versionConstants() {
    // Access version constants directly
    fmt.Printf("Version: %s\n", version.Version)
    fmt.Printf("Git Commit: %s\n", version.GitCommit)
    fmt.Printf("Build Date: %s\n", version.BuildDate)
    fmt.Printf("Go Version: %s\n", version.GoVersion)
}
```

## Advanced Usage

### Custom Version Display

```go
func customVersionDisplay() {
    // Create custom version display
    fmt.Println("=== Application Information ===")
    fmt.Printf("Application: tykctl\n")
    fmt.Printf("Version: %s\n", version.Version)
    fmt.Printf("Build: %s\n", version.GitCommit)
    fmt.Printf("Built: %s\n", version.BuildDate)
    fmt.Printf("Go: %s\n", version.GoVersion)
    fmt.Println("================================")
}
```

### Version Comparison

```go
func versionComparison() {
    currentVersion := version.Version
    
    // Simple version comparison (semantic versioning)
    if currentVersion >= "1.0.0" {
        fmt.Println("Running a stable version")
    } else {
        fmt.Println("Running a pre-release version")
    }
    
    // Check for specific version
    if currentVersion == "1.0.0" {
        fmt.Println("Running version 1.0.0")
    }
}
```

### Build Information Display

```go
func buildInformation() {
    fmt.Println("=== Build Information ===")
    fmt.Printf("Version: %s\n", version.Version)
    fmt.Printf("Git Commit: %s\n", version.GitCommit)
    fmt.Printf("Build Date: %s\n", version.BuildDate)
    fmt.Printf("Go Version: %s\n", version.GoVersion)
    
    // Additional build context
    fmt.Printf("Build Time: %s\n", time.Now().Format(time.RFC3339))
    fmt.Printf("Architecture: %s\n", runtime.GOARCH)
    fmt.Printf("Operating System: %s\n", runtime.GOOS)
    fmt.Println("==========================")
}
```

## Integration Examples

### With CLI Commands

```go
import (
    "github.com/spf13/cobra"
    "github.com/edsonmichaque/tykctl-go/version"
)

func createVersionCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Show version information",
        Run: func(cmd *cobra.Command, args []string) {
            // Check for verbose flag
            verbose, _ := cmd.Flags().GetBool("verbose")
            
            if verbose {
                fmt.Println(version.Full())
            } else {
                fmt.Println(version.String())
            }
        },
    }
}

func init() {
    versionCmd := createVersionCommand()
    versionCmd.Flags().BoolP("verbose", "v", false, "Show detailed version information")
    rootCmd.AddCommand(versionCmd)
}
```

### With Extension Management

```go
func extensionVersionManagement() {
    // List extensions with version info
    extensions := []string{"tykctl-portal", "tykctl-cloud", "tykctl-dashboard"}
    
    fmt.Println("=== Extension Versions ===")
    for _, ext := range extensions {
        fmt.Printf("%s: %s\n", ext, version.Info(ext))
    }
    fmt.Println("==========================")
}
```

### With API Responses

```go
func apiVersionResponse() map[string]interface{} {
    return map[string]interface{}{
        "application": "tykctl",
        "version": version.Version,
        "build": map[string]interface{}{
            "commit": version.GitCommit,
            "date":   version.BuildDate,
            "go":     version.GoVersion,
        },
        "runtime": map[string]interface{}{
            "arch": runtime.GOARCH,
            "os":   runtime.GOOS,
        },
    }
}
```

### With Logging

```go
func loggingWithVersion() {
    config := logger.Config{Debug: true}
    zapLogger := logger.New(config)
    
    // Log version information
    zapLogger.Info("Application started",
        zap.String("version", version.Version),
        zap.String("commit", version.GitCommit),
        zap.String("build_date", version.BuildDate),
        zap.String("go_version", version.GoVersion),
    )
}
```

## Version Information Structure

### Available Variables

```go
var (
    Version   = "1.0.0"        // Application version
    GitCommit = "unknown"      // Git commit hash
    BuildDate = "unknown"      // Build timestamp
    GoVersion = runtime.Version() // Go runtime version
)
```

### Version Functions

- `String()` - Returns version string (e.g., "v1.0.0")
- `Short()` - Returns short version (e.g., "1.0.0")
- `Full()` - Returns full version information
- `Info(extensionName)` - Returns version info for extension
- `InfoFull(extensionName)` - Returns full version info for extension

## Build-Time Configuration

### Setting Version Information

Version information is typically set at build time using ldflags:

```bash
# Build with version information
go build -ldflags "-X github.com/edsonmichaque/tykctl-go/version.Version=1.0.0 -X github.com/edsonmichaque/tykctl-go/version.GitCommit=$(git rev-parse HEAD) -X github.com/edsonmichaque/tykctl-go/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o tykctl
```

### Makefile Example

```makefile
VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

build:
	go build -ldflags "-X github.com/edsonmichaque/tykctl-go/version.Version=$(VERSION) -X github.com/edsonmichaque/tykctl-go/version.GitCommit=$(COMMIT) -X github.com/edsonmichaque/tykctl-go/version.BuildDate=$(BUILD_DATE)" -o tykctl
```

## Testing

### Version Testing

```go
func TestVersion(t *testing.T) {
    // Test version string format
    versionStr := version.String()
    assert.True(t, strings.HasPrefix(versionStr, "v"))
    
    // Test short version
    shortVersion := version.Short()
    assert.NotEmpty(t, shortVersion)
    
    // Test full version contains all components
    fullVersion := version.Full()
    assert.Contains(t, fullVersion, "version")
    assert.Contains(t, fullVersion, "Git commit")
    assert.Contains(t, fullVersion, "Build date")
    assert.Contains(t, fullVersion, "Go version")
}
```

### Extension Version Testing

```go
func TestExtensionVersion(t *testing.T) {
    extensionName := "test-extension"
    
    // Test extension version info
    info := version.Info(extensionName)
    assert.Contains(t, info, extensionName)
    assert.Contains(t, info, version.Version)
    
    // Test full extension version info
    fullInfo := version.InfoFull(extensionName)
    assert.Contains(t, fullInfo, extensionName)
    assert.Contains(t, fullInfo, version.Version)
    assert.Contains(t, fullInfo, version.GitCommit)
}
```

## Best Practices

- **Build-Time Setting**: Set version information at build time using ldflags
- **Consistent Formatting**: Use consistent version string formatting
- **Extension Naming**: Use clear extension names in version info
- **Runtime Information**: Include relevant runtime information
- **Documentation**: Document version information in help text
- **Testing**: Test version information display and formatting

## Dependencies

- No external dependencies
- Uses only Go standard library (`fmt`, `runtime`)