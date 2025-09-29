# Template Package

The `template` package provides a unified template system for loading, processing, and validating templates from various sources. It supports both file-based and URL-based template loading with comprehensive validation and processing capabilities.

## Features

- **Multi-source Loading**: Load templates from local files, directories, or remote URLs
- **Template Processing**: Process templates with variable substitution using Go's text/template
- **Validation**: Comprehensive template and variable validation
- **Context Support**: Full context-aware operations with cancellation support
- **Retry Logic**: Built-in retry mechanism for URL loading
- **Custom Headers**: Support for authentication and custom HTTP headers
- **Directory Templates**: Load templates from directories with flexible resolution
- **Custom Resolvers**: Define custom template resolution logic
- **Type Safety**: Strongly typed template structures and validation

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "github.com/edsonmichaque/tykctl-go/template"
)

func main() {
    // Load template from file
    loader := template.NewFileLoader("template.yaml")
    
    // Load template
    ctx := context.Background()
    tmpl, err := loader.Load(ctx)
    if err != nil {
        panic(err)
    }
    
    // Validate template
    validator := template.NewValidator()
    if err := validator.Validate(tmpl); err != nil {
        panic(err)
    }
    
    // Process template
    processor := template.NewProcessor(loader, validator)
    variables := map[string]interface{}{
        "Name":        "My Product",
        "Description": "A great product",
    }
    
    result, err := processor.Process(ctx, variables)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Processed template: %s\n", string(result))
}
```

## Template Structure

Templates are defined in YAML format with the following structure:

```yaml
name: "Product Template"
description: "Template for creating products"
resource_type: "product"
version: "1.0.0"
author: "tykctl-portal"
tags: ["api", "product", "standard"]
variables:
  - name: "Name"
    type: "string"
    required: true
    description: "Product name"
    validation:
      min_length: 1
      max_length: 100
  - name: "Description"
    type: "string"
    required: false
    description: "Product description"
content:
  display_name: "{{.Name}}"
  description: "{{.Description}}"
  visibility: "public"
  policies:
    - name: "Rate Limiting"
      config:
        rate_limit: 1000
        per: "hour"
```

## Loading Templates

### File Loading

```go
// Load from local file
loader := template.NewFileLoader("template.yaml")
template, err := loader.Load(ctx)
```

### Directory Loading

```go
// Load from directory by name
loader := template.NewDirLoader("/path/to/templates", "product", nil)
template, err := loader.Load(ctx)

// List available templates
templates, err := loader.ListTemplates()

// Get template path
path, err := loader.GetTemplatePath("product")
```

### URL Loading

```go
// Load from URL with custom options
options := &template.Options{
    Headers: map[string]string{
        "Authorization": "Bearer token",
    },
    Timeout:    30,
    Retries:    3,
    RetryDelay: 1,
}

config := &template.Config{
    Client: &http.Client{},
}

loader := template.NewURLLoader(config, "https://api.example.com/template.yaml", options)
template, err := loader.Load(ctx)
```

### Custom Resolvers

```go
// Create custom resolver
resolver := template.NewCustomResolver(func(dir, name string) (string, error) {
    // Custom resolution logic
    path := filepath.Join(dir, name+".custom")
    if _, err := os.Stat(path); err != nil {
        return "", err
    }
    return path, nil
})

// Use custom resolver
loader := template.NewDirLoader("/templates", "product", resolver)
```

## Processing Templates

The processor loads, validates, and processes templates in one step:

```go
validator := template.NewValidator()
processor := template.NewProcessor(loader, validator)

variables := map[string]interface{}{
    "Name":        "My Product",
    "Description": "A great product",
}

// This will load, validate, and process the template
result, err := processor.Process(ctx, variables)
```

## Validating Templates

```go
validator := template.NewValidator()

// Validate template structure
err := validator.Validate(template)
if err != nil {
    // Handle validation errors
}
```

## Configuration

### Options

```go
options := &template.Options{
    Headers:      map[string]string{}, // Custom HTTP headers
    Timeout:      30,                  // Request timeout in seconds
    Retries:      3,                   // Number of retry attempts
    RetryDelay:   1,                   // Delay between retries in seconds
    ValidateSSL:  true,                // SSL certificate validation
    AllowedHosts: []string{},          // Allowed hosts for URL loading
}
```

### Config

```go
config := &template.Config{
    Client: &http.Client{
        Timeout: 30 * time.Second,
    },
}
```

## Error Handling

The package provides a comprehensive error system with specific error types and structured error information:

### Standard Errors

```go
import "github.com/edsonmichaque/tykctl-go/template"

// Template errors
if errors.Is(err, template.ErrTemplateNotFound) {
    // Template not found
}

if errors.Is(err, template.ErrTemplateInvalid) {
    // Invalid template format
}

// Network errors
if errors.Is(err, template.ErrNetworkTimeout) {
    // Network timeout
}

if errors.Is(err, template.ErrHostBlocked) {
    // Host not in allowlist
}
```

### Structured Error Types

```go
// Check error types
if template.IsValidationError(err) {
    var validationErr *template.ValidationError
    if errors.As(err, &validationErr) {
        fmt.Printf("Field: %s, Rule: %s, Value: %v\n", 
            validationErr.Field, validationErr.Rule, validationErr.Value)
    }
}

if template.IsLoaderError(err) {
    var loaderErr *template.LoaderError
    if errors.As(err, &loaderErr) {
        fmt.Printf("Source: %s, Type: %s\n", 
            loaderErr.Source, loaderErr.Type)
    }
}

if template.IsProcessingError(err) {
    var processingErr *template.ProcessingError
    if errors.As(err, &processingErr) {
        fmt.Printf("Template: %s, Variable: %s\n", 
            processingErr.Template, processingErr.Variable)
    }
}
```

### Error Wrapping

```go
// Wrap errors with context
err = template.WrapError(err, "failed to process template")
err = template.WrapErrorf(err, "failed to process template %s", templateName)

// Add specific context
err = template.ErrWithPath(err, "/path/to/template.yaml")
err = template.ErrWithTemplate(err, "product-template")
err = template.ErrWithVariable(err, "name")
```

## Variable Types

Supported variable types:

- `string`: Text values
- `integer`: Whole numbers
- `number`: Decimal numbers
- `boolean`: True/false values
- `array`: Lists of values
- `object`: Key-value pairs

## Validation Rules

Variables support various validation rules:

```yaml
variables:
  - name: "Name"
    type: "string"
    validation:
      min_length: 1
      max_length: 100
      pattern: "^[a-zA-Z0-9\\s]+$"
      enum: ["option1", "option2", "option3"]
  - name: "Count"
    type: "integer"
    validation:
      min_value: 0
      max_value: 1000
```

## Examples

See the `example_test.go` file for comprehensive usage examples including:

- Basic template loading and processing
- File-based template loading
- URL-based template loading with authentication
- Template validation
- Error handling

## Testing

Run the tests with:

```bash
go test -v ./template
```

## Dependencies

- `gopkg.in/yaml.v3`: YAML parsing and marshaling
- `text/template`: Go template processing
- Standard library packages for HTTP, context, and file operations

## License

This package is part of the tykctl-go project and follows the same license terms.
