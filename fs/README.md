# Filesystem Package

The `fs` package provides a filesystem abstraction layer for tykctl, offering context-aware filesystem operations with support for both real and in-memory filesystems.

## Features

- **Context Support**: All operations support context.Context for cancellation and timeouts
- **Filesystem Abstraction**: Uses afero for filesystem abstraction
- **In-Memory Support**: Support for in-memory filesystems for testing
- **Idempotent Operations**: Safe operations that can be called multiple times
- **Error Handling**: Comprehensive error handling with context awareness
- **Cross-platform**: Works consistently across different operating systems

## Usage

### Basic Filesystem Operations

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/edsonmichaque/tykctl-go/fs"
)

func main() {
    // Create a new filesystem instance
    filesystem := fs.New()
    
    ctx := context.Background()
    
    // Create directory
    err := filesystem.MkdirAll(ctx, "/tmp/myapp/config", 0755)
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if file exists
    exists, err := filesystem.Exists(ctx, "/tmp/myapp/config/app.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    if !exists {
        fmt.Println("Config file does not exist")
    }
    
    fmt.Println("Filesystem operations completed!")
}
```

### In-Memory Filesystem for Testing

```go
func testWithInMemoryFS() {
    // Create in-memory filesystem
    memFS := fs.NewMem()
    
    ctx := context.Background()
    
    // Create directory structure
    err := memFS.MkdirAll(ctx, "/app/config", 0755)
    if err != nil {
        log.Fatal(err)
    }
    
    // Write test data
    testData := []byte("test configuration")
    err = memFS.WriteFile(ctx, "/app/config/test.yaml", testData, 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    // Read data back
    data, err := memFS.ReadFile(ctx, "/app/config/test.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Read data: %s\n", string(data))
}
```

### Context-Aware Operations

```go
func contextAwareOperations() error {
    filesystem := fs.New()
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Perform operations with context
    err := filesystem.MkdirAll(ctx, "/tmp/large-operation", 0755)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("operation timed out")
        }
        return err
    }
    
    return nil
}
```

## Advanced Usage

### File Operations with Validation

```go
func safeFileOperations() error {
    filesystem := fs.New()
    ctx := context.Background()
    
    configPath := "/tmp/app/config.yaml"
    
    // Check if parent directory exists
    parentDir := filepath.Dir(configPath)
    exists, err := filesystem.Exists(ctx, parentDir)
    if err != nil {
        return err
    }
    
    if !exists {
        // Create parent directory
        err = filesystem.MkdirAll(ctx, parentDir, 0755)
        if err != nil {
            return fmt.Errorf("failed to create directory: %w", err)
        }
    }
    
    // Write configuration file
    configData := []byte("app: config\nversion: 1.0.0\n")
    err = filesystem.WriteFile(ctx, configPath, configData, 0644)
    if err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }
    
    // Verify file was written
    exists, err = filesystem.Exists(ctx, configPath)
    if err != nil {
        return err
    }
    
    if !exists {
        return fmt.Errorf("config file was not created")
    }
    
    return nil
}
```

### Directory Operations

```go
func directoryOperations() error {
    filesystem := fs.New()
    ctx := context.Background()
    
    baseDir := "/tmp/myapp"
    
    // Create directory structure
    dirs := []string{
        filepath.Join(baseDir, "config"),
        filepath.Join(baseDir, "logs"),
        filepath.Join(baseDir, "cache"),
        filepath.Join(baseDir, "data"),
    }
    
    for _, dir := range dirs {
        err := filesystem.MkdirAll(ctx, dir, 0755)
        if err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dir, err)
        }
    }
    
    // List directory contents
    entries, err := filesystem.ReadDir(ctx, baseDir)
    if err != nil {
        return err
    }
    
    fmt.Printf("Created %d directories:\n", len(entries))
    for _, entry := range entries {
        fmt.Printf("- %s\n", entry.Name())
    }
    
    return nil
}
```

### File Copying and Moving

```go
func fileOperations() error {
    filesystem := fs.New()
    ctx := context.Background()
    
    sourceFile := "/tmp/source.txt"
    destFile := "/tmp/destination.txt"
    
    // Create source file
    sourceData := []byte("Hello, World!")
    err := filesystem.WriteFile(ctx, sourceFile, sourceData, 0644)
    if err != nil {
        return err
    }
    
    // Copy file
    err = filesystem.CopyFile(ctx, sourceFile, destFile)
    if err != nil {
        return err
    }
    
    // Verify copy
    destData, err := filesystem.ReadFile(ctx, destFile)
    if err != nil {
        return err
    }
    
    if string(destData) != string(sourceData) {
        return fmt.Errorf("file copy verification failed")
    }
    
    // Move file
    movedFile := "/tmp/moved.txt"
    err = filesystem.MoveFile(ctx, destFile, movedFile)
    if err != nil {
        return err
    }
    
    // Verify move
    exists, err := filesystem.Exists(ctx, destFile)
    if err != nil {
        return err
    }
    
    if exists {
        return fmt.Errorf("source file still exists after move")
    }
    
    exists, err = filesystem.Exists(ctx, movedFile)
    if err != nil {
        return err
    }
    
    if !exists {
        return fmt.Errorf("destination file does not exist after move")
    }
    
    return nil
}
```

## Integration Examples

### With Configuration Management

```go
func loadConfiguration() (*Config, error) {
    filesystem := fs.New()
    ctx := context.Background()
    
    configPath := filepath.Join(os.Getenv("HOME"), ".config", "myapp", "config.yaml")
    
    // Check if config exists
    exists, err := filesystem.Exists(ctx, configPath)
    if err != nil {
        return nil, err
    }
    
    if !exists {
        // Create default config
        defaultConfig := &Config{
            Server: "localhost:8080",
            Debug:  false,
        }
        
        configData, err := yaml.Marshal(defaultConfig)
        if err != nil {
            return nil, err
        }
        
        // Ensure config directory exists
        configDir := filepath.Dir(configPath)
        err = filesystem.MkdirAll(ctx, configDir, 0755)
        if err != nil {
            return nil, err
        }
        
        // Write default config
        err = filesystem.WriteFile(ctx, configPath, configData, 0644)
        if err != nil {
            return nil, err
        }
    }
    
    // Read config
    configData, err := filesystem.ReadFile(ctx, configPath)
    if err != nil {
        return nil, err
    }
    
    var config Config
    err = yaml.Unmarshal(configData, &config)
    if err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### With Logging

```go
func setupLogging() error {
    filesystem := fs.New()
    ctx := context.Background()
    
    logDir := "/tmp/myapp/logs"
    
    // Create log directory
    err := filesystem.MkdirAll(ctx, logDir, 0755)
    if err != nil {
        return err
    }
    
    // Create log file
    logFile := filepath.Join(logDir, "app.log")
    
    // Check if log file exists and get its size
    exists, err := filesystem.Exists(ctx, logFile)
    if err != nil {
        return err
    }
    
    if exists {
        info, err := filesystem.Stat(ctx, logFile)
        if err != nil {
            return err
        }
        
        // Rotate log if it's too large (> 10MB)
        if info.Size() > 10*1024*1024 {
            rotatedFile := filepath.Join(logDir, "app.log.1")
            err = filesystem.MoveFile(ctx, logFile, rotatedFile)
            if err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

### With Testing

```go
func TestFileOperations(t *testing.T) {
    // Use in-memory filesystem for testing
    filesystem := fs.NewMem()
    ctx := context.Background()
    
    // Test directory creation
    err := filesystem.MkdirAll(ctx, "/test/dir", 0755)
    assert.NoError(t, err)
    
    // Test file writing
    testData := []byte("test data")
    err = filesystem.WriteFile(ctx, "/test/dir/file.txt", testData, 0644)
    assert.NoError(t, err)
    
    // Test file reading
    data, err := filesystem.ReadFile(ctx, "/test/dir/file.txt")
    assert.NoError(t, err)
    assert.Equal(t, testData, data)
    
    // Test file existence
    exists, err := filesystem.Exists(ctx, "/test/dir/file.txt")
    assert.NoError(t, err)
    assert.True(t, exists)
}
```

## Error Handling

### Context-Aware Error Handling

```go
func handleContextErrors() error {
    filesystem := fs.New()
    
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    err := filesystem.MkdirAll(ctx, "/tmp/slow-operation", 0755)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("operation timed out")
        }
        if ctx.Err() == context.Canceled {
            return fmt.Errorf("operation was canceled")
        }
        return fmt.Errorf("filesystem operation failed: %w", err)
    }
    
    return nil
}
```

### File Operation Error Handling

```go
func robustFileOperations() error {
    filesystem := fs.New()
    ctx := context.Background()
    
    configPath := "/tmp/config.yaml"
    
    // Try to read config with fallback
    configData, err := filesystem.ReadFile(ctx, configPath)
    if err != nil {
        if os.IsNotExist(err) {
            // Create default config
            defaultConfig := []byte("default: config\n")
            err = filesystem.WriteFile(ctx, configPath, defaultConfig, 0644)
            if err != nil {
                return fmt.Errorf("failed to create default config: %w", err)
            }
            configData = defaultConfig
        } else {
            return fmt.Errorf("failed to read config: %w", err)
        }
    }
    
    return nil
}
```

## Best Practices

- **Context Usage**: Always use context for cancellation and timeout handling
- **Error Handling**: Handle filesystem errors gracefully with appropriate fallbacks
- **Idempotent Operations**: Use idempotent operations like `MkdirAll` for safety
- **Testing**: Use in-memory filesystem for unit tests
- **Path Handling**: Use `filepath` package for cross-platform path handling
- **Permissions**: Set appropriate file permissions for security

## Dependencies

- `github.com/spf13/afero` - Filesystem abstraction library