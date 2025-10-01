# API Package

The `api` package provides a high-level HTTP client for making API requests with built-in retry logic, middleware support, and response handling. It's designed to be a general-purpose HTTP client library that can be used in any Go project.

## Features

- **HTTP Client**: Simple interface for GET, POST, PUT, DELETE, PATCH requests
- **Retry Logic**: Configurable retry strategies with exponential backoff using industry-standard libraries
- **Middleware Support**: Chainable middleware for logging, timeouts, authentication, and custom logic
- **Response Handling**: Rich response objects with status code checking and JSON unmarshaling
- **Pagination Support**: Built-in pagination handling for API responses
- **Error Handling**: Comprehensive error types and handling with retryable error detection
- **Context Support**: Full context.Context integration for cancellation and timeouts
- **Configurable**: Flexible configuration options for different use cases
- **Framework Agnostic**: No dependencies on specific frameworks or libraries
- **Extensible**: Easy to extend with custom middleware and retry conditions
- **Production Ready**: Built with production use cases in mind

## Use Cases

This API package is designed for general-purpose HTTP client needs:

- **REST API Clients**: Build clients for any REST API service
- **Microservices Communication**: Inter-service communication in microservices architectures
- **Third-party Integrations**: Integrate with external APIs and services
- **Data Fetching**: Robust data fetching with retry logic and error handling
- **Web Scraping**: HTTP client for web scraping applications
- **API Testing**: HTTP client for API testing and validation
- **CLI Tools**: HTTP client for command-line tools and utilities
- **Background Jobs**: HTTP client for background job processing

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/your-org/your-project/pkg/api"
)

func main() {
    // Create API client with functional options
    client := api.New(
        api.WithBaseURL("https://api.example.com/v1"),
        api.WithClientTimeout(30 * time.Second),
        api.WithClientHeader("Accept", "application/json"),
        api.WithClientHeader("Authorization", "Bearer your-api-token"),
    )
    ctx := context.Background()
    
    // Make a GET request with functional options
    resp, err := client.Get(ctx, "/users",
        api.WithClientHeader("Accept", "application/json"),
        api.WithQuery("page", "1"),
        api.WithQuery("limit", "10"),
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    if resp.IsSuccess() {
        fmt.Printf("Success: %s\n", resp.String())
    } else {
        fmt.Printf("Error: %d - %s\n", resp.StatusCode, resp.String())
    }
}
```

## Functional Options

The API client supports functional options for configuring individual requests. This provides a clean, flexible way to add headers, query parameters, and other request-specific options.

### Available Options

#### Headers
```go
// Single header
client.Get(ctx, "/users", api.WithHeader("Accept", "application/json"))

// Multiple headers
client.Get(ctx, "/users", 
    api.WithHeaders(map[string]string{
        "Accept": "application/json",
        "X-API-Version": "v1",
    }),
)
```

#### Query Parameters
```go
// Single query parameter
client.Get(ctx, "/users", api.WithQuery("page", "1"))

// Multiple query parameters
client.Get(ctx, "/users",
    api.WithQueries(map[string]string{
        "page": "1",
        "limit": "10",
        "sort": "name",
    }),
)
```

#### Request Body
```go
// JSON body from interface
client.Post(ctx, "/users", userData, api.WithJSONBody(userData))

// Raw body
client.Post(ctx, "/users", nil, api.WithBody([]byte("raw data")))
```

#### Request Options
```go
// Timeout
client.Get(ctx, "/users", api.WithTimeout(10*time.Second))

// Retries
client.Get(ctx, "/users", api.WithRetries(3))

// Retry delay
client.Get(ctx, "/users", api.WithRetryDelay(2*time.Second))
```

### Complete Example
```go
// Complex request with multiple options
resp, err := client.Post(ctx, "/users", userData,
    api.WithClientHeader("Content-Type", "application/json"),
    api.WithHeader("X-Request-ID", "req-123"),
    api.WithQuery("validate", "true"),
    api.WithTimeout(15*time.Second),
    api.WithRetries(3),
)
```

## Client Configuration

The API client uses a flexible configuration interface that allows for different configuration strategies.

### Client Configuration with Functional Options

```go
// Create a client with functional options
client := api.New(
    api.WithBaseURL("https://api.example.com/v1"),
    api.WithClientTimeout(30 * time.Second),
    api.WithClientHeader("Authorization", "Bearer your-token"),
    api.WithClientHeader("Accept", "application/json"),
    api.WithClientHeader("Content-Type", "application/json"),
    api.WithUserAgent("my-app/1.0.0"),
)
```

### Available Client Options

```go
// Base URL
api.WithBaseURL("https://api.example.com/v1")

// Timeout
api.WithClientTimeout(30 * time.Second)

// Single header
api.WithClientHeader("Authorization", "Bearer token")

// Multiple headers
api.WithClientHeaders(map[string]string{
    "Accept": "application/json",
    "Content-Type": "application/json",
})

// User agent
api.WithUserAgent("my-app/1.0.0")
```

## Making Requests

### GET Request
```go
resp, err := client.Get(ctx, "/users")
```

### POST Request
```go
user := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
}

resp, err := client.Post(ctx, "/users", user)
```

### PUT Request
```go
user := map[string]interface{}{
    "id":    1,
    "name":  "John Doe Updated",
    "email": "john.updated@example.com",
}

resp, err := client.Put(ctx, "/users/1", user)
```

### DELETE Request
```go
resp, err := client.Delete(ctx, "/users/1")
```

### PATCH Request
```go
updates := map[string]interface{}{
    "name": "John Doe Patched",
}

resp, err := client.Patch(ctx, "/users/1", updates)
```

### Generic Request
```go
resp, err := client.Request(ctx, "OPTIONS", "/users", nil)
```

## Response Handling

```go
resp, err := client.Get(ctx, "/users")
if err != nil {
    return err
}

// Check response status
if resp.IsSuccess() {
    fmt.Println("Request successful")
} else if resp.IsClientError() {
    fmt.Println("Client error (4xx)")
} else if resp.IsServerError() {
    fmt.Println("Server error (5xx)")
}

// Get response data
fmt.Printf("Status: %d\n", resp.StatusCode)
fmt.Printf("Body: %s\n", resp.String())

// Get specific header
contentType := resp.GetHeader("Content-Type")

// Unmarshal JSON response
var users []User
if err := resp.UnmarshalJSON(&users); err != nil {
    return err
}
```

## Retry Logic

### Basic Retry
```go
retryConfig := api.NewExponentialBackoffConfig(
    3,                    // max retries
    1*time.Second,        // initial delay
    10*time.Second,       // max delay
    30*time.Second,       // max elapsed time
)

resp, err := api.WithRetry(ctx, retryConfig, func() (*api.Response, error) {
    return client.Get(ctx, "/users")
})
```

### Custom Retry Condition
```go
customRetry := &api.CustomRetryCondition{
    RetryableStatusCodes: []int{429, 500, 502, 503, 504},
}

retryConfig := api.RetryConfig{
    MaxRetries:     5,
    InitialDelay:   2 * time.Second,
    MaxDelay:       2 * time.Second, // Constant delay
    MaxElapsedTime: 30 * time.Second,
    Multiplier:     1.0, // No exponential growth
    Retryable:      customRetry,
}
```

### Retry Configurations

#### Exponential Backoff
```go
retryConfig := api.NewExponentialBackoffConfig(
    3,                    // max retries
    1*time.Second,        // initial delay
    30*time.Second,       // max delay
    2*time.Minute,        // max elapsed time
)
```

#### Constant Backoff
```go
retryConfig := api.NewConstantBackoffConfig(
    5,                    // max retries
    2*time.Second,        // delay
)
```

#### Custom Configuration
```go
retryConfig := api.RetryConfig{
    MaxRetries:     3,
    InitialDelay:   1 * time.Second,
    MaxDelay:       30 * time.Second,
    MaxElapsedTime: 2 * time.Minute,
    Multiplier:     2.0,
    Retryable:      &api.DefaultRetryCondition{},
}
```

## Middleware

### Logging Middleware
```go
middleware := api.LoggingMiddleware(logger)
```

### Timeout Middleware
```go
middleware := api.TimeoutMiddleware(10 * time.Second)
```

### Authentication Middleware
```go
authMiddleware := api.AuthMiddleware(func(req *api.Request) error {
    req.Headers["Authorization"] = "Bearer " + token
    return nil
})
```

### Header Middleware
```go
headerMiddleware := api.HeaderMiddleware(map[string]string{
    "User-Agent": "my-app/1.0.0",
    "Accept":     "application/json",
})
```

### Chaining Middleware
```go
middlewares := api.ChainMiddleware(
    api.LoggingMiddleware(logger),
    api.TimeoutMiddleware(10*time.Second),
    api.HeaderMiddleware(headers),
    authMiddleware,
)
```

## Pagination

```go
// Make paginated request
resp, err := client.Get(ctx, "/users?page=1&per_page=10")
if err != nil {
    return err
}

// Parse paginated response
var paginatedResp api.PaginatedResponse
if err := resp.UnmarshalJSON(&paginatedResp); err != nil {
    return err
}

fmt.Printf("Total users: %d\n", paginatedResp.Pagination.Total)
fmt.Printf("Current page: %d\n", paginatedResp.Pagination.Page)
fmt.Printf("Has next page: %t\n", paginatedResp.Pagination.HasNext)
```

## Error Handling

```go
resp, err := client.Get(ctx, "/users")
if err != nil {
    // Handle network or request errors
    if retryErr, ok := err.(*api.RetryableError); ok {
        fmt.Printf("Retry failed after %d attempts: %v\n", 
            retryErr.MaxRetries, retryErr.Err)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

if !resp.IsSuccess() {
    // Handle HTTP error responses
    apiErr := api.NewError(
        resp.StatusCode,
        "API request failed",
        resp.String(),
        resp.Headers,
    )
    return apiErr
}
```

## Request Options

```go
options := api.DefaultRequestOptions().
    WithHeader("X-Custom-Header", "value").
    WithTimeout(15 * time.Second).
    WithRetries(5).
    WithRetryDelay(2 * time.Second)
```

## Examples

See `example.go` for comprehensive usage examples including:
- Basic API client usage
- Retry functionality
- Middleware usage
- Paginated requests
- Custom retry conditions

## Client Flexibility

The API client is designed to be flexible and adaptable to different use cases:

### Immutable Client Creation
```go
// Create base client with functional options
baseClient := api.New(
    api.WithBaseURL("https://api.example.com/v1"),
    api.WithClientTimeout(30 * time.Second),
)

// Clone for different endpoints
userClient := baseClient.WithBaseURL("https://api.example.com/v1/users")
adminClient := baseClient.WithBaseURL("https://api.example.com/v1/admin")

// Create clients with different timeouts
fastClient := baseClient.WithTimeout(5 * time.Second)
slowClient := baseClient.WithTimeout(60 * time.Second)

// Add authentication
authClient := baseClient.WithHeader("Authorization", "Bearer token123")

// Add multiple headers
customClient := baseClient.WithHeaders(map[string]string{
    "X-API-Key": "your-api-key",
    "X-Client-Version": "1.0.0",
})
```

### Dynamic Configuration
```go
// Modify client configuration at runtime
client.SetBaseURL("https://api.example.com/v2")
client.SetTimeout(45 * time.Second)
client.SetHeader("Authorization", "Bearer new-token")
client.SetUserAgent("MyApp/2.0.0")

// Access underlying HTTP client for advanced usage
httpClient := client.GetHTTPClient()
config := client.GetConfig()
```

## Integration with Other Packages

The API package integrates seamlessly with other Go packages:

```go
// Use with any logger package
logger := logrus.New()
client := api.New(
    api.WithBaseURL("https://api.example.com"),
    api.WithClientHeaders(map[string]string{
        "Authorization": "Bearer " + token,
    }),
)

// Use with progress indicators for long-running requests
spinner := progress.New()
err := spinner.WithContext(ctx, "Fetching data...", func() error {
    resp, err := client.Get(ctx, "/data")
    if err != nil {
        return err
    }
    // Process response...
    return nil
})

// Use with any HTTP client library
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}
client.SetHTTPClient(httpClient)
```
