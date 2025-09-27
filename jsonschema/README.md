# JSON Schema Package

The `jsonschema` package provides comprehensive JSON Schema validation functionality, allowing you to validate JSON data against schemas from various sources including strings, files, and URLs.

## Features

- **Schema Validation**: Validate JSON data against JSON Schema specifications
- **Multiple Sources**: Support for schemas from strings, files, and URLs
- **Detailed Errors**: Comprehensive validation error reporting with field-level details
- **Context Support**: Full context.Context integration for cancellation and timeouts
- **Flexible Input**: Support for various input formats (strings, bytes, readers)
- **Error Recovery**: Detailed error information for debugging validation issues

## Usage

### Basic Schema Validation

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/jsonschema"
)

func main() {
    // Define a JSON schema
    schema := `{
        "type": "object",
        "properties": {
            "name": {"type": "string"},
            "age": {"type": "number", "minimum": 0},
            "email": {"type": "string", "format": "email"}
        },
        "required": ["name", "age"]
    }`
    
    // Create validator
    validator, err := jsonschema.New(schema)
    if err != nil {
        log.Fatal(err)
    }
    
    // Sample JSON data
    jsonData := `{
        "name": "John Doe",
        "age": 30,
        "email": "john@example.com"
    }`
    
    // Validate data
    ctx := context.Background()
    result, err := validator.ValidateString(ctx, jsonData)
    if err != nil {
        log.Fatal(err)
    }
    
    if result.Valid {
        fmt.Println("Data is valid!")
    } else {
        fmt.Printf("Validation failed: %d errors\n", len(result.Errors))
        for _, err := range result.Errors {
            fmt.Printf("- %s: %s\n", err.Field, err.Description)
        }
    }
}
```

### Schema from File

```go
func validateFromFile() error {
    // Create validator from schema file
    validator, err := jsonschema.NewFromFile("schema.json")
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    
    // Validate data from file
    result, err := validator.ValidateFile(ctx, "data.json")
    if err != nil {
        return err
    }
    
    if !result.Valid {
        fmt.Printf("Validation failed with %d errors:\n", len(result.Errors))
        for _, err := range result.Errors {
            fmt.Printf("- %s: %s\n", err.Field, err.Description)
        }
    }
    
    return nil
}
```

### Schema from URL

```go
func validateFromURL() error {
    // Create validator from remote schema
    validator, err := jsonschema.NewFromURL("https://json-schema.org/draft/2020-12/schema")
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    
    // Validate local data against remote schema
    jsonData := `{"$schema": "https://json-schema.org/draft/2020-12/schema"}`
    result, err := validator.ValidateString(ctx, jsonData)
    if err != nil {
        return err
    }
    
    if result.Valid {
        fmt.Println("Schema is valid!")
    }
    
    return nil
}
```

## Advanced Usage

### Complex Schema Validation

```go
func complexValidation() error {
    schema := `{
        "type": "object",
        "properties": {
            "user": {
                "type": "object",
                "properties": {
                    "id": {"type": "integer", "minimum": 1},
                    "name": {"type": "string", "minLength": 1},
                    "email": {"type": "string", "format": "email"},
                    "profile": {
                        "type": "object",
                        "properties": {
                            "age": {"type": "integer", "minimum": 0, "maximum": 150},
                            "country": {"type": "string", "enum": ["US", "CA", "UK", "DE"]}
                        },
                        "required": ["age"]
                    }
                },
                "required": ["id", "name", "email"]
            },
            "settings": {
                "type": "object",
                "properties": {
                    "notifications": {"type": "boolean"},
                    "theme": {"type": "string", "enum": ["light", "dark"]}
                }
            }
        },
        "required": ["user"]
    }`
    
    validator, err := jsonschema.New(schema)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    
    // Valid data
    validData := `{
        "user": {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "profile": {
                "age": 30,
                "country": "US"
            }
        },
        "settings": {
            "notifications": true,
            "theme": "dark"
        }
    }`
    
    result, err := validator.ValidateString(ctx, validData)
    if err != nil {
        return err
    }
    
    if result.Valid {
        fmt.Println("Complex validation passed!")
    }
    
    return nil
}
```

### Array Validation

```go
func arrayValidation() error {
    schema := `{
        "type": "array",
        "items": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "name": {"type": "string"},
                "tags": {
                    "type": "array",
                    "items": {"type": "string"},
                    "minItems": 1,
                    "uniqueItems": true
                }
            },
            "required": ["id", "name"]
        },
        "minItems": 1,
        "maxItems": 10
    }`
    
    validator, err := jsonschema.New(schema)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    
    jsonData := `[
        {"id": 1, "name": "Item 1", "tags": ["tag1", "tag2"]},
        {"id": 2, "name": "Item 2", "tags": ["tag3"]}
    ]`
    
    result, err := validator.ValidateString(ctx, jsonData)
    if err != nil {
        return err
    }
    
    if result.Valid {
        fmt.Println("Array validation passed!")
    } else {
        for _, err := range result.Errors {
            fmt.Printf("Error: %s - %s\n", err.Field, err.Description)
        }
    }
    
    return nil
}
```

### Conditional Validation

```go
func conditionalValidation() error {
    schema := `{
        "type": "object",
        "properties": {
            "type": {"type": "string", "enum": ["user", "admin"]},
            "permissions": {"type": "array", "items": {"type": "string"}}
        },
        "required": ["type"],
        "if": {
            "properties": {"type": {"const": "admin"}}
        },
        "then": {
            "properties": {
                "permissions": {
                    "type": "array",
                    "minItems": 1,
                    "items": {"enum": ["read", "write", "delete"]}
                }
            },
            "required": ["permissions"]
        }
    }`
    
    validator, err := jsonschema.New(schema)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    
    // Admin user (requires permissions)
    adminData := `{
        "type": "admin",
        "permissions": ["read", "write"]
    }`
    
    result, err := validator.ValidateString(ctx, adminData)
    if err != nil {
        return err
    }
    
    if result.Valid {
        fmt.Println("Admin validation passed!")
    }
    
    return nil
}
```

## Integration Examples

### With Configuration Validation

```go
func validateConfig() error {
    configSchema := `{
        "type": "object",
        "properties": {
            "server": {
                "type": "object",
                "properties": {
                    "host": {"type": "string", "minLength": 1},
                    "port": {"type": "integer", "minimum": 1, "maximum": 65535}
                },
                "required": ["host", "port"]
            },
            "database": {
                "type": "object",
                "properties": {
                    "url": {"type": "string", "format": "uri"},
                    "max_connections": {"type": "integer", "minimum": 1, "maximum": 100}
                },
                "required": ["url"]
            }
        },
        "required": ["server", "database"]
    }`
    
    validator, err := jsonschema.New(configSchema)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    
    // Validate configuration file
    result, err := validator.ValidateFile(ctx, "config.json")
    if err != nil {
        return err
    }
    
    if !result.Valid {
        return fmt.Errorf("configuration validation failed: %v", result.Errors)
    }
    
    return nil
}
```

### With API Request Validation

```go
func validateAPIRequest() error {
    requestSchema := `{
        "type": "object",
        "properties": {
            "method": {"type": "string", "enum": ["GET", "POST", "PUT", "DELETE"]},
            "path": {"type": "string", "pattern": "^/api/v[0-9]+/.*$"},
            "headers": {
                "type": "object",
                "properties": {
                    "Authorization": {"type": "string", "pattern": "^Bearer .+$"},
                    "Content-Type": {"type": "string", "enum": ["application/json"]}
                },
                "required": ["Authorization"]
            },
            "body": {"type": "object"}
        },
        "required": ["method", "path"]
    }`
    
    validator, err := jsonschema.New(requestSchema)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    
    requestData := `{
        "method": "POST",
        "path": "/api/v1/users",
        "headers": {
            "Authorization": "Bearer token123",
            "Content-Type": "application/json"
        },
        "body": {
            "name": "John Doe",
            "email": "john@example.com"
        }
    }`
    
    result, err := validator.ValidateString(ctx, requestData)
    if err != nil {
        return err
    }
    
    if result.Valid {
        fmt.Println("API request is valid!")
    }
    
    return nil
}
```

### With Data Processing Pipeline

```go
func processDataWithValidation() error {
    dataSchema := `{
        "type": "array",
        "items": {
            "type": "object",
            "properties": {
                "id": {"type": "string", "pattern": "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"},
                "name": {"type": "string", "minLength": 1, "maxLength": 100},
                "value": {"type": "number", "minimum": 0},
                "timestamp": {"type": "string", "format": "date-time"}
            },
            "required": ["id", "name", "value", "timestamp"]
        }
    }`
    
    validator, err := jsonschema.New(dataSchema)
    if err != nil {
        return err
    }
    
    ctx := context.Background()
    
    // Process data file
    result, err := validator.ValidateFile(ctx, "data.json")
    if err != nil {
        return err
    }
    
    if !result.Valid {
        fmt.Printf("Data validation failed with %d errors:\n", len(result.Errors))
        for _, err := range result.Errors {
            fmt.Printf("- Field '%s': %s\n", err.Field, err.Description)
        }
        return fmt.Errorf("data validation failed")
    }
    
    fmt.Println("Data validation passed!")
    return nil
}
```

## Error Handling

### Detailed Error Reporting

```go
func detailedErrorHandling() error {
    validator, err := jsonschema.New(schema)
    if err != nil {
        return fmt.Errorf("failed to create validator: %w", err)
    }
    
    ctx := context.Background()
    result, err := validator.ValidateString(ctx, jsonData)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    if !result.Valid {
        fmt.Printf("Validation failed with %d errors:\n", len(result.Errors))
        
        for i, err := range result.Errors {
            fmt.Printf("Error %d:\n", i+1)
            fmt.Printf("  Field: %s\n", err.Field)
            fmt.Printf("  Description: %s\n", err.Description)
            if err.Context != "" {
                fmt.Printf("  Context: %s\n", err.Context)
            }
            fmt.Println()
        }
        
        return fmt.Errorf("validation failed")
    }
    
    return nil
}
```

### Context-Aware Validation

```go
func contextAwareValidation() error {
    validator, err := jsonschema.New(schema)
    if err != nil {
        return err
    }
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    result, err := validator.ValidateString(ctx, jsonData)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("validation timed out")
        }
        return fmt.Errorf("validation failed: %w", err)
    }
    
    return nil
}
```

## Best Practices

- **Schema Design**: Design schemas to be clear, specific, and maintainable
- **Error Handling**: Always handle validation errors gracefully
- **Performance**: Cache validators for repeated use
- **Context Usage**: Use context for cancellation and timeout handling
- **Documentation**: Document schema requirements and constraints
- **Testing**: Test schemas with various valid and invalid data samples

## Dependencies

- `github.com/xeipuuv/gojsonschema` - JSON Schema validation library