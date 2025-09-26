# tykctl-go Documentation

## Overview

`tykctl-go` is a comprehensive Go library for creating Tyk CLI extensions with best practices, templates, and utilities. It provides a framework for building consistent, maintainable, and feature-rich CLI extensions.

## Core Components

### 1. Extension Framework

The `Extension` struct is the core of the framework:

```go
type Extension struct {
    name        string
    version     string
    description string
    rootCmd     *cobra.Command
    logger      *zap.Logger
    config      *Config
}
```

#### Creating an Extension

```go
ext := NewExtension("my-extension", "1.0.0")
ext.SetDescription("My awesome extension")
```

#### Adding Commands

```go
ext.AddCommand(NewCommand("hello", "Say hello", helloHandler))
ext.AddCommands(cmd1, cmd2, cmd3)
```

#### Executing an Extension

```go
if err := ext.Execute(context.Background()); err != nil {
    log.Fatal(err)
}
```

### 2. Command Creation

#### Basic Command

```go
cmd := NewCommand("hello", "Say hello", func(cmd *cobra.Command, args []string) error {
    cmd.Println("Hello from my extension!")
    return nil
})
```

#### Root Command

```go
rootCmd := NewRootCommand("my-extension", "My Extension Description")
```

#### Version Command

```go
versionCmd := NewVersionCommand("my-extension", "1.0.0")
```

### 3. Template System

The template system provides pre-built templates for common extension patterns:

#### Available Templates

- `basic` - Basic extension template
- `api` - API management extension
- `webhook` - Webhook extension
- `monitoring` - Monitoring extension
- `auth` - Authentication extension
- `data` - Data processing extension

#### Using Templates

```go
// Generate basic extension
gen := NewGenerator("my-extension", "./output")
err := gen.Generate()

// Generate with specific template
err := gen.GenerateWithTemplate("api")
```

#### Template Data

```go
type TemplateData struct {
    ExtensionName    string
    ExtensionVersion string
    ModuleName       string
    Description      string
    Author           string
    License          string
    GoVersion        string
}
```

### 4. Generator System

The generator system helps create new extensions:

#### Basic Generation

```go
gen := NewGenerator("my-extension", "./output")
gen.SetVersion("1.0.0")
gen.SetDescription("My extension")
gen.SetAuthor("My Name")
gen.SetLicense("MIT")
err := gen.Generate()
```

#### Template-Specific Generation

```go
gen := NewGenerator("my-api-extension", "./output")
err := gen.GenerateWithTemplate("api")
```

#### Pre-built Generators

```go
// Create specific extension types
err := CreateBasicExtension("my-basic", "./output")
err := CreateAPIExtension("my-api", "./output")
err := CreateWebhookExtension("my-webhook", "./output")
err := CreateMonitoringExtension("my-monitoring", "./output")
err := CreateAuthExtension("my-auth", "./output")
err := CreateDataExtension("my-data", "./output")
```

### 5. Utility Functions

#### String Utilities

```go
utils := NewStringUtils()

// Convert to different cases
kebab := utils.ToKebabCase("my-awesome-extension")
snake := utils.ToSnakeCase("my-awesome-extension")
camel := utils.ToCamelCase("my-awesome-extension")
pascal := utils.ToPascalCase("my-awesome-extension")
```

#### File Utilities

```go
fileUtils := NewFileUtils()

// Create directory
err := fileUtils.CreateDirectory("./my-dir")

// Write file
err := fileUtils.WriteFile("./file.txt", "content")

// Read file
content, err := fileUtils.ReadFile("./file.txt")

// Check existence
exists := fileUtils.FileExists("./file.txt")
isDir := fileUtils.DirectoryExists("./dir")
```

#### Validation Utilities

```go
validation := NewValidationUtils()

// Validate extension name
err := validation.ValidateExtensionName("my-extension")

// Validate version
err := validation.ValidateVersion("1.0.0")

// Validate path
err := validation.ValidatePath("./path/to/file")
```

### 6. CLI Tool

The `tykctl-go` CLI tool provides command-line interface for generating extensions:

#### Generate Extension

```bash
# Basic extension
tykctl-go generate my-extension

# With template
tykctl-go generate my-api-extension --template api

# With custom options
tykctl-go generate my-extension \
  --template basic \
  --output ./my-output \
  --version 1.0.0 \
  --description "My awesome extension" \
  --author "My Name" \
  --license MIT
```

#### List Templates

```bash
tykctl-go list
```

#### Show Version

```bash
tykctl-go version
```

### 7. Example Manager

The example manager helps create and manage example extensions:

```go
em := NewExampleManager("./examples")

// Create example
err := em.CreateExample("my-example", "basic")

// List examples
examples, err := em.ListExamples()

// Get example path
path := em.GetExamplePath("my-example")

// Create all examples
err := em.CreateAllExamples()
```

### 8. Configuration

Extensions support configuration through the `Config` struct:

```go
type Config struct {
    Debug   bool
    Verbose bool
    NoColor bool
}
```

#### Using Configuration

```go
ext := NewExtension("my-extension", "1.0.0")
config := ext.GetConfig()

if config.Debug {
    // Debug mode
}

logger := ext.GetLogger()
logger.Info("Extension started")
```

### 9. Best Practices

#### 1. Use the Framework

Always use the provided framework instead of building from scratch:

```go
// Good
ext := NewExtension("my-extension", "1.0.0")
ext.AddCommand(NewCommand("hello", "Say hello", helloHandler))

// Avoid
cmd := &cobra.Command{...}
```

#### 2. Handle Errors Properly

Always handle errors and provide meaningful messages:

```go
func helloHandler(cmd *cobra.Command, args []string) error {
    if err := doSomething(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}
```

#### 3. Use Structured Logging

Use the provided logger for consistent logging:

```go
logger := ext.GetLogger()
logger.Info("Command executed", zap.String("command", "hello"))
logger.Debug("Debug information", zap.String("arg", arg))
```

#### 4. Provide Helpful Help Text

Always provide clear help text for commands:

```go
cmd := NewCommand("hello", "Say hello to the user", helloHandler)
cmd.Long = "This command says hello to the user with a friendly message."
```

#### 5. Use Templates

Use the provided templates for consistent structure:

```go
// Generate extension with template
gen := NewGenerator("my-extension", "./output")
err := gen.GenerateWithTemplate("api")
```

#### 6. Test Your Extensions

Always test your extensions:

```go
func TestMyExtension(t *testing.T) {
    ext := NewExtension("test-extension", "1.0.0")
    // Test your extension
}
```

#### 7. Follow Go Conventions

Follow Go naming conventions and idioms:

```go
// Good
func NewMyCommand() *cobra.Command

// Avoid
func new_my_command() *cobra.Command
```

### 10. Integration with Existing Extensions

The library integrates seamlessly with existing Tyk CLI extensions:

```go
// In your extension's main.go
package main

import (
    "context"
    "fmt"
    "os"
    "github.com/edsonmichaque/tykctl/pkg/tykctl-go"
)

func main() {
    ext := tykctl.NewExtension("my-extension", "1.0.0")
    
    // Add your commands
    ext.AddCommand(tykctl.NewCommand("hello", "Say hello", helloHandler))
    
    // Execute
    if err := ext.Execute(context.Background()); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### 11. Advanced Usage

#### Custom Templates

You can create custom templates by extending the template system:

```go
// Custom template data
data := TemplateData{
    ExtensionName:    "my-custom-extension",
    ExtensionVersion: "1.0.0",
    ModuleName:       "my-custom-extension",
    Description:      "My custom extension",
    Author:           "My Name",
    License:          "MIT",
    GoVersion:        "1.25",
}

// Generate with custom template
err := GenerateFromTemplate("custom", data, "./output")
```

#### Custom Generators

You can create custom generators for specific use cases:

```go
func CreateCustomExtension(name, outputDir string) error {
    gen := NewGenerator(name, outputDir)
    gen.SetDescription("Custom extension")
    gen.SetAuthor("Custom Author")
    gen.SetLicense("Custom License")
    return gen.GenerateWithTemplate("custom")
}
```

### 12. Troubleshooting

#### Common Issues

1. **Template not found**: Make sure the template name is correct and exists
2. **Permission denied**: Check file permissions for output directory
3. **Invalid extension name**: Use only alphanumeric characters, hyphens, and underscores
4. **Missing dependencies**: Run `go mod tidy` to install dependencies

#### Debug Mode

Enable debug mode for detailed logging:

```go
ext := NewExtension("my-extension", "1.0.0")
ext.GetConfig().Debug = true
```

#### Validation

Always validate inputs:

```go
if err := ValidateExtensionName(name); err != nil {
    return fmt.Errorf("invalid extension name: %w", err)
}
```

## Conclusion

The `tykctl-go` library provides a comprehensive framework for creating Tyk CLI extensions in Go. It includes templates, utilities, validation, and best practices to help developers create consistent, maintainable, and feature-rich CLI extensions.

For more examples and advanced usage, see the `example_usage.go` file and the test suite.
