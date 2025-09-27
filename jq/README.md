# JQ Package

The `jq` package provides JSON querying functionality using the JQ language, allowing you to filter, transform, and manipulate JSON data with powerful query expressions.

## Features

- **JQ Language Support**: Full support for JQ query language syntax
- **JSON Processing**: Filter, transform, and manipulate JSON data
- **Go Integration**: Native Go implementation using gojq
- **Error Handling**: Comprehensive error handling for malformed queries and data
- **Type Safety**: Proper handling of different JSON data types
- **Performance**: Efficient JSON processing with minimal memory overhead

## Usage

### Basic JSON Querying

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/jq"
)

func main() {
    // Sample JSON data
    jsonData := []byte(`{
        "users": [
            {"id": 1, "name": "Alice", "age": 30},
            {"id": 2, "name": "Bob", "age": 25},
            {"id": 3, "name": "Charlie", "age": 35}
        ]
    }`)
    
    // Query to get all user names
    query := ".users[].name"
    
    result, err := jq.Process(jsonData, query)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User names: %s\n", string(result))
    // Output: ["Alice","Bob","Charlie"]
}
```

### Filtering Data

```go
func filterUsers() error {
    jsonData := []byte(`{
        "users": [
            {"id": 1, "name": "Alice", "age": 30, "active": true},
            {"id": 2, "name": "Bob", "age": 25, "active": false},
            {"id": 3, "name": "Charlie", "age": 35, "active": true}
        ]
    }`)
    
    // Filter active users
    query := ".users[] | select(.active == true)"
    
    result, err := jq.Process(jsonData, query)
    if err != nil {
        return err
    }
    
    fmt.Printf("Active users: %s\n", string(result))
    // Output: [{"id":1,"name":"Alice","age":30,"active":true},{"id":3,"name":"Charlie","age":35,"active":true}]
    
    return nil
}
```

### Data Transformation

```go
func transformData() error {
    jsonData := []byte(`{
        "products": [
            {"name": "Laptop", "price": 999.99, "category": "Electronics"},
            {"name": "Book", "price": 19.99, "category": "Education"},
            {"name": "Phone", "price": 699.99, "category": "Electronics"}
        ]
    }`)
    
    // Transform to simplified format
    query := ".products[] | {name: .name, price: .price, category: .category}"
    
    result, err := jq.Process(jsonData, query)
    if err != nil {
        return err
    }
    
    fmt.Printf("Transformed data: %s\n", string(result))
    
    return nil
}
```

## Advanced Usage

### Complex Queries

```go
func complexQueries() error {
    jsonData := []byte(`{
        "orders": [
            {"id": 1, "customer": "Alice", "items": [{"name": "Laptop", "price": 999.99}]},
            {"id": 2, "customer": "Bob", "items": [{"name": "Book", "price": 19.99}]},
            {"id": 3, "customer": "Alice", "items": [{"name": "Phone", "price": 699.99}]}
        ]
    }`)
    
    // Get total order value for Alice
    query := `.orders[] | select(.customer == "Alice") | .items[].price | add`
    
    result, err := jq.Process(jsonData, query)
    if err != nil {
        return err
    }
    
    fmt.Printf("Alice's total: %s\n", string(result))
    
    return nil
}
```

### Array Operations

```go
func arrayOperations() error {
    jsonData := []byte(`{
        "numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    }`)
    
    // Get even numbers
    query := ".numbers[] | select(. % 2 == 0)"
    
    result, err := jq.Process(jsonData, query)
    if err != nil {
        return err
    }
    
    fmt.Printf("Even numbers: %s\n", string(result))
    // Output: [2,4,6,8,10]
    
    // Get sum of all numbers
    sumQuery := ".numbers | add"
    sumResult, err := jq.Process(jsonData, sumQuery)
    if err != nil {
        return err
    }
    
    fmt.Printf("Sum: %s\n", string(sumResult))
    // Output: 55
    
    return nil
}
```

### String Operations

```go
func stringOperations() error {
    jsonData := []byte(`{
        "messages": [
            {"text": "Hello World", "priority": "high"},
            {"text": "Good morning", "priority": "low"},
            {"text": "Important notice", "priority": "high"}
        ]
    }`)
    
    // Get uppercase messages with high priority
    query := `.messages[] | select(.priority == "high") | .text | ascii_upcase`
    
    result, err := jq.Process(jsonData, query)
    if err != nil {
        return err
    }
    
    fmt.Printf("High priority messages: %s\n", string(result))
    // Output: ["HELLO WORLD","IMPORTANT NOTICE"]
    
    return nil
}
```

## Integration Examples

### With HTTP Client

```go
import (
    "github.com/edsonmichaque/tykctl-go/httpclient"
    "github.com/edsonmichaque/tykctl-go/jq"
)

func processAPIResponse() error {
    client := httpclient.NewWithBaseURL("https://api.github.com")
    ctx := context.Background()
    
    // Get GitHub user data
    resp, err := client.Get(ctx, "/users/octocat")
    if err != nil {
        return err
    }
    
    if !resp.IsSuccess() {
        return fmt.Errorf("API request failed: %d", resp.StatusCode)
    }
    
    // Process response with JQ
    query := "{name: .name, login: .login, followers: .followers, public_repos: .public_repos}"
    result, err := jq.Process(resp.Body, query)
    if err != nil {
        return err
    }
    
    fmt.Printf("User info: %s\n", string(result))
    return nil
}
```

### With Configuration Processing

```go
func processConfig() error {
    configData := []byte(`{
        "database": {
            "host": "localhost",
            "port": 5432,
            "name": "myapp"
        },
        "redis": {
            "host": "localhost",
            "port": 6379
        },
        "features": {
            "auth": true,
            "logging": false,
            "metrics": true
        }
    }`)
    
    // Extract database configuration
    dbQuery := ".database | {host: .host, port: .port, name: .name}"
    dbResult, err := jq.Process(configData, dbQuery)
    if err != nil {
        return err
    }
    
    fmt.Printf("Database config: %s\n", string(dbResult))
    
    // Extract enabled features
    featuresQuery := ".features | to_entries | map(select(.value == true)) | map(.key)"
    featuresResult, err := jq.Process(configData, featuresQuery)
    if err != nil {
        return err
    }
    
    fmt.Printf("Enabled features: %s\n", string(featuresResult))
    
    return nil
}
```

### With Log Processing

```go
func processLogs() error {
    logData := []byte(`{
        "logs": [
            {"level": "ERROR", "message": "Database connection failed", "timestamp": "2023-01-01T10:00:00Z"},
            {"level": "INFO", "message": "User logged in", "timestamp": "2023-01-01T10:01:00Z"},
            {"level": "ERROR", "message": "File not found", "timestamp": "2023-01-01T10:02:00Z"},
            {"level": "WARN", "message": "High memory usage", "timestamp": "2023-01-01T10:03:00Z"}
        ]
    }`)
    
    // Get error messages
    errorQuery := `.logs[] | select(.level == "ERROR") | {message: .message, timestamp: .timestamp}`
    errorResult, err := jq.Process(logData, errorQuery)
    if err != nil {
        return err
    }
    
    fmt.Printf("Error logs: %s\n", string(errorResult))
    
    // Count logs by level
    countQuery := `.logs | group_by(.level) | map({level: .[0].level, count: length})`
    countResult, err := jq.Process(logData, countQuery)
    if err != nil {
        return err
    }
    
    fmt.Printf("Log counts: %s\n", string(countResult))
    
    return nil
}
```

## Error Handling

### Query Validation

```go
func validateQuery(query string) error {
    // Test query with empty data
    testData := []byte("{}")
    
    _, err := jq.Process(testData, query)
    if err != nil {
        return fmt.Errorf("invalid JQ query: %w", err)
    }
    
    return nil
}
```

### Data Validation

```go
func validateJSONData(data []byte) error {
    // Test with simple query
    _, err := jq.Process(data, ".")
    if err != nil {
        return fmt.Errorf("invalid JSON data: %w", err)
    }
    
    return nil
}
```

## Common JQ Patterns

### Data Extraction

```go
// Extract specific fields
query := ".users[] | {id: .id, name: .name}"

// Extract nested values
query := ".data.items[].metadata.tags[]"

// Extract with conditions
query := ".items[] | select(.price > 100) | .name"
```

### Data Aggregation

```go
// Count items
query := ".items | length"

// Sum values
query := ".items[].price | add"

// Average
query := ".items[].price | add / length"

// Min/Max
query := ".items[].price | min"
query := ".items[].price | max"
```

### Data Transformation

```go
// Rename fields
query := ".users[] | {user_id: .id, full_name: .name}"

// Add computed fields
query := ".users[] | {name: .name, age: .age, category: (if .age > 30 then "adult" else "young" end)}"

// Filter and transform
query := ".users[] | select(.active) | {name: .name, email: .email}"
```

## Best Practices

- **Query Testing**: Test JQ queries with sample data before using in production
- **Error Handling**: Always handle JQ processing errors gracefully
- **Performance**: Use specific queries rather than processing large datasets
- **Documentation**: Document complex JQ queries for maintainability
- **Validation**: Validate JSON data before processing with JQ

## Dependencies

- `github.com/itchyny/gojq` - Pure Go implementation of JQ