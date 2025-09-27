# Editor Package

The `editor` package provides cross-platform file editing functionality, allowing applications to open files in external editors with proper timeout handling and context support.

## Features

- **Cross-platform Support**: Works on Windows, macOS, and Linux
- **Default Editor Detection**: Automatically detects the system's default editor
- **Custom Editor Support**: Allows specifying custom editors and arguments
- **Timeout Handling**: Configurable timeout for editor operations
- **Context Support**: Full context.Context integration for cancellation
- **Error Handling**: Comprehensive error handling and reporting

## Usage

### Basic File Editing

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/editor"
)

func main() {
    // Create a new editor instance
    e := editor.New()
    
    // Edit a file
    ctx := context.Background()
    err := e.EditFile(ctx, "config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("File edited successfully!")
}
```

### Custom Editor Configuration

```go
func main() {
    // Create editor with specific editor
    e := editor.NewWithEditor("vim")
    
    // Set custom arguments
    e.SetArgs([]string{"-c", "set number"})
    
    // Set timeout
    e.SetTimeout(10 * time.Minute)
    
    // Edit file
    ctx := context.Background()
    err := e.EditFile(ctx, "config.yaml")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Editor with Context and Timeout

```go
func main() {
    e := editor.New()
    
    // Set timeout
    e.SetTimeout(5 * time.Minute)
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
    defer cancel()
    
    // Edit file with context
    err := e.EditFile(ctx, "config.yaml")
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("Editor operation timed out")
        } else {
            fmt.Printf("Editor error: %v\n", err)
        }
        return
    }
    
    fmt.Println("File edited successfully!")
}
```

## Editor Detection

### Default Editor Priority

The package automatically detects the system's default editor in the following order:

1. **Environment Variables**:
   - `EDITOR` (primary)
   - `VISUAL` (fallback)

2. **Platform-specific Defaults**:
   - **Windows**: `notepad.exe`
   - **macOS**: `open -t` (opens in default text editor)
   - **Linux**: `nano`, `vim`, `vi` (in order of preference)

### Custom Editor Setup

```go
// Use specific editor
e := editor.NewWithEditor("code") // VS Code

// Use editor with arguments
e := editor.NewWithEditor("vim")
e.SetArgs([]string{"-c", "set number", "-c", "set syntax=yaml"})

// Use editor with specific mode
e := editor.NewWithEditor("emacs")
e.SetArgs([]string{"--no-window-system", "+10:10"})
```

## Advanced Usage

### Editor with Validation

```go
func editConfigFile(filename string) error {
    e := editor.New()
    e.SetTimeout(10 * time.Minute)
    
    ctx := context.Background()
    
    // Check if file exists
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return fmt.Errorf("file does not exist: %s", filename)
    }
    
    // Edit file
    err := e.EditFile(ctx, filename)
    if err != nil {
        return fmt.Errorf("failed to edit file: %w", err)
    }
    
    // Validate edited file
    return validateConfigFile(filename)
}
```

### Editor with Backup

```go
func editWithBackup(filename string) error {
    // Create backup
    backupFile := filename + ".backup"
    if err := copyFile(filename, backupFile); err != nil {
        return fmt.Errorf("failed to create backup: %w", err)
    }
    
    // Edit file
    e := editor.New()
    ctx := context.Background()
    
    err := e.EditFile(ctx, filename)
    if err != nil {
        // Restore backup on error
        if restoreErr := copyFile(backupFile, filename); restoreErr != nil {
            return fmt.Errorf("edit failed and backup restore failed: %w, %w", err, restoreErr)
        }
        return fmt.Errorf("edit failed, backup restored: %w", err)
    }
    
    // Remove backup on success
    os.Remove(backupFile)
    return nil
}
```

### Editor with Progress Indication

```go
import (
    "github.com/edsonmichaque/tykctl-go/progress"
)

func editWithProgress(filename string) error {
    spinner := progress.New()
    spinner.SetMessage("Opening editor...")
    
    e := editor.New()
    
    // Start spinner
    spinner.Start()
    defer spinner.Stop()
    
    ctx := context.Background()
    err := e.EditFile(ctx, filename)
    
    if err != nil {
        spinner.SetMessage("Edit failed")
        return err
    }
    
    spinner.SetMessage("Edit completed")
    return nil
}
```

## Integration Examples

### With CLI Commands

```go
import (
    "github.com/spf13/cobra"
    "github.com/edsonmichaque/tykctl-go/editor"
)

func createEditCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "edit [file]",
        Short: "Edit a configuration file",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            filename := args[0]
            
            e := editor.New()
            
            // Get editor from flag
            if editorFlag, _ := cmd.Flags().GetString("editor"); editorFlag != "" {
                e.SetEditor(editorFlag)
            }
            
            // Get timeout from flag
            if timeoutFlag, _ := cmd.Flags().GetDuration("timeout"); timeoutFlag > 0 {
                e.SetTimeout(timeoutFlag)
            }
            
            return e.EditFile(cmd.Context(), filename)
        },
    }
}

func init() {
    editCmd := createEditCommand()
    editCmd.Flags().StringP("editor", "e", "", "Editor to use")
    editCmd.Flags().DurationP("timeout", "t", 5*time.Minute, "Editor timeout")
    
    rootCmd.AddCommand(editCmd)
}
```

### With Configuration Management

```go
func editUserConfig() error {
    configDir := os.Getenv("TYKCTL_CONFIG_DIR")
    if configDir == "" {
        configDir = filepath.Join(os.Getenv("HOME"), ".config", "tykctl")
    }
    
    configFile := filepath.Join(configDir, "config.yaml")
    
    // Ensure config directory exists
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }
    
    // Create default config if it doesn't exist
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        if err := createDefaultConfig(configFile); err != nil {
            return fmt.Errorf("failed to create default config: %w", err)
        }
    }
    
    // Edit config file
    e := editor.New()
    return e.EditFile(context.Background(), configFile)
}
```

## Error Handling

### Common Error Scenarios

```go
err := e.EditFile(ctx, filename)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "no suitable editor found"):
        fmt.Println("No editor found. Please set EDITOR environment variable.")
    case strings.Contains(err.Error(), "context deadline exceeded"):
        fmt.Println("Editor operation timed out.")
    case strings.Contains(err.Error(), "file not found"):
        fmt.Println("File does not exist.")
    default:
        fmt.Printf("Editor error: %v\n", err)
    }
}
```

### Error Recovery

```go
func editWithRetry(filename string) error {
    e := editor.New()
    
    for attempt := 1; attempt <= 3; attempt++ {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
        
        err := e.EditFile(ctx, filename)
        cancel()
        
        if err == nil {
            return nil
        }
        
        if attempt < 3 {
            fmt.Printf("Edit attempt %d failed, retrying...\n", attempt)
            time.Sleep(time.Second)
        }
    }
    
    return fmt.Errorf("failed to edit file after 3 attempts")
}
```

## Best Practices

- **Timeout Configuration**: Set appropriate timeouts for editor operations
- **Context Usage**: Use context for cancellation and timeout handling
- **Error Handling**: Handle editor errors gracefully with user-friendly messages
- **File Validation**: Validate files before and after editing
- **Backup Strategy**: Consider creating backups for important files
- **Editor Detection**: Provide fallback options when default editor isn't available

## Environment Variables

- `EDITOR` - Primary editor to use
- `VISUAL` - Visual editor (fallback)
- `TYKCTL_EDITOR` - Tykctl-specific editor override

## Dependencies

- No external dependencies
- Uses only Go standard library (`os/exec`, `context`, `time`, `os`)