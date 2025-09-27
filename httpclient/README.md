# HTTP Client Package

The `httpclient` package provides a simple, flexible HTTP client wrapper with support for base URLs, custom headers, timeouts, and JSON operations.

## Features

- **Base URL Support**: Set a base URL for all requests
- **Custom Headers**: Add custom headers to requests
- **Timeout Configuration**: Configurable request timeouts
- **JSON Support**: Built-in JSON marshaling and unmarshaling
- **Response Handling**: Rich response objects with status checking
- **Error Handling**: Comprehensive error handling and reporting
- **Context Support**: Full context.Context integration

## Usage

### Basic HTTP Client

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/httpclient"
)

func main() {
    // Create a new HTTP client
    client := httpclient.New()
    
    ctx := context.Background()
    
    // Make a GET request
    resp, err := client.Get(ctx, "https://api.github.com/users/octocat")
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.IsSuccess() {
        fmt.Printf("Response: %s\n", resp.String())
    } else {
        fmt.Printf("Error: %d - %s\n", resp.StatusCode, resp.String())
    }
}
```

### Client with Base URL

```go
func main() {
    // Create client with base URL
    client := httpclient.NewWithBaseURL("https://api.github.com")
    
    ctx := context.Background()
    
    // Make requests relative to base URL
    resp, err := client.Get(ctx, "/users/octocat")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User data: %s\n", resp.String())
}
```

### Client with Timeout

```go
func main() {
    // Create client with custom timeout
    client := httpclient.NewWithTimeout(10 * time.Second)
    
    ctx := context.Background()
    
    // Make request with timeout
    resp, err := client.Get(ctx, "https://slow-api.example.com/data")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Response: %s\n", resp.String())
}
```

### Client with Custom Headers

```go
func main() {
    client := httpclient.New()
    
    // Set custom headers
    client.SetHeader("Authorization", "Bearer your-token")
    client.SetHeader("User-Agent", "MyApp/1.0.0")
    client.SetHeader("Accept", "application/json")
    
    ctx := context.Background()
    
    resp, err := client.Get(ctx, "https://api.example.com/protected")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Response: %s\n", resp.String())
}
```

## Advanced Usage

### JSON Operations

```go
type User struct {
    ID       int    `json:"id"`
    Username string `json:"login"`
    Name     string `json:"name"`
    Email    string `json:"email"`
}

func jsonOperations() error {
    client := httpclient.NewWithBaseURL("https://api.github.com")
    client.SetHeader("Accept", "application/json")
    
    ctx := context.Background()
    
    // GET request with JSON unmarshaling
    resp, err := client.Get(ctx, "/users/octocat")
    if err != nil {
        return err
    }
    
    if !resp.IsSuccess() {
        return fmt.Errorf("API request failed: %d", resp.StatusCode)
    }
    
    var user User
    err = resp.UnmarshalJSON(&user)
    if err != nil {
        return err
    }
    
    fmt.Printf("User: %s (%s)\n", user.Name, user.Username)
    
    // POST request with JSON marshaling
    newUser := User{
        Username: "newuser",
        Name:     "New User",
        Email:    "newuser@example.com",
    }
    
    resp, err = client.PostJSON(ctx, "/users", newUser)
    if err != nil {
        return err
    }
    
    if resp.IsSuccess() {
        fmt.Println("User created successfully")
    }
    
    return nil
}
```

### Request Methods

```go
func allRequestMethods() error {
    client := httpclient.NewWithBaseURL("https://api.example.com")
    ctx := context.Background()
    
    // GET request
    resp, err := client.Get(ctx, "/users/1")
    if err != nil {
        return err
    }
    fmt.Printf("GET: %d\n", resp.StatusCode)
    
    // POST request
    data := map[string]string{"name": "John Doe"}
    resp, err = client.PostJSON(ctx, "/users", data)
    if err != nil {
        return err
    }
    fmt.Printf("POST: %d\n", resp.StatusCode)
    
    // PUT request
    updateData := map[string]string{"name": "Jane Doe"}
    resp, err = client.PutJSON(ctx, "/users/1", updateData)
    if err != nil {
        return err
    }
    fmt.Printf("PUT: %d\n", resp.StatusCode)
    
    // PATCH request
    patchData := map[string]string{"email": "jane@example.com"}
    resp, err = client.PatchJSON(ctx, "/users/1", patchData)
    if err != nil {
        return err
    }
    fmt.Printf("PATCH: %d\n", resp.StatusCode)
    
    // DELETE request
    resp, err = client.Delete(ctx, "/users/1")
    if err != nil {
        return err
    }
    fmt.Printf("DELETE: %d\n", resp.StatusCode)
    
    return nil
}
```

### Response Handling

```go
func responseHandling() error {
    client := httpclient.New()
    ctx := context.Background()
    
    resp, err := client.Get(ctx, "https://api.example.com/data")
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
    fmt.Printf("Content-Type: %s\n", contentType)
    
    // Get all headers
    headers := resp.GetHeaders()
    for name, values := range headers {
        fmt.Printf("%s: %v\n", name, values)
    }
    
    return nil
}
```

## Integration Examples

### With Configuration

```go
type APIConfig struct {
    BaseURL string
    Token   string
    Timeout time.Duration
}

func createClientFromConfig(config APIConfig) *httpclient.Client {
    client := httpclient.NewWithBaseURL(config.BaseURL)
    client.SetTimeout(config.Timeout)
    client.SetHeader("Authorization", "Bearer "+config.Token)
    client.SetHeader("Accept", "application/json")
    client.SetHeader("User-Agent", "MyApp/1.0.0")
    
    return client
}
```

### With Error Handling

```go
func robustAPIRequest() error {
    client := httpclient.NewWithBaseURL("https://api.example.com")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    resp, err := client.Get(ctx, "/data")
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("request timed out")
        }
        return fmt.Errorf("request failed: %w", err)
    }
    
    if !resp.IsSuccess() {
        return fmt.Errorf("API error: %d - %s", resp.StatusCode, resp.String())
    }
    
    return nil
}
```

### With Retry Logic

```go
func requestWithRetry() error {
    client := httpclient.New()
    
    maxRetries := 3
    for attempt := 1; attempt <= maxRetries; attempt++ {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        
        resp, err := client.Get(ctx, "https://api.example.com/data")
        cancel()
        
        if err == nil && resp.IsSuccess() {
            return nil
        }
        
        if attempt < maxRetries {
            fmt.Printf("Attempt %d failed, retrying...\n", attempt)
            time.Sleep(time.Duration(attempt) * time.Second)
        }
    }
    
    return fmt.Errorf("request failed after %d attempts", maxRetries)
}
```

## Client Configuration

### Dynamic Configuration

```go
func configureClient() *httpclient.Client {
    client := httpclient.New()
    
    // Set base URL
    client.SetBaseURL("https://api.example.com/v1")
    
    // Set timeout
    client.SetTimeout(30 * time.Second)
    
    // Set headers
    client.SetHeader("Authorization", "Bearer token123")
    client.SetHeader("Accept", "application/json")
    client.SetHeader("Content-Type", "application/json")
    client.SetHeader("User-Agent", "MyApp/1.0.0")
    
    return client
}
```

### Client Information

```go
func clientInfo(client *httpclient.Client) {
    fmt.Printf("Base URL: %s\n", client.GetBaseURL())
    fmt.Printf("Timeout: %v\n", client.GetTimeout())
    
    headers := client.GetHeaders()
    fmt.Printf("Headers: %v\n", headers)
}
```

## Error Handling

### Common Error Scenarios

```go
func handleErrors() error {
    client := httpclient.New()
    ctx := context.Background()
    
    resp, err := client.Get(ctx, "https://api.example.com/data")
    if err != nil {
        // Handle network errors
        if strings.Contains(err.Error(), "timeout") {
            return fmt.Errorf("request timed out")
        }
        if strings.Contains(err.Error(), "connection refused") {
            return fmt.Errorf("server is not available")
        }
        return fmt.Errorf("network error: %w", err)
    }
    
    if !resp.IsSuccess() {
        // Handle HTTP errors
        switch resp.StatusCode {
        case 400:
            return fmt.Errorf("bad request: %s", resp.String())
        case 401:
            return fmt.Errorf("unauthorized: check your credentials")
        case 403:
            return fmt.Errorf("forbidden: insufficient permissions")
        case 404:
            return fmt.Errorf("not found: resource does not exist")
        case 429:
            return fmt.Errorf("rate limited: too many requests")
        case 500:
            return fmt.Errorf("server error: %s", resp.String())
        default:
            return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.String())
        }
    }
    
    return nil
}
```

## Best Practices

- **Context Usage**: Always use context for cancellation and timeout handling
- **Error Handling**: Handle both network and HTTP errors appropriately
- **Headers**: Set appropriate headers for your API (Accept, User-Agent, etc.)
- **Timeouts**: Set reasonable timeouts for your use case
- **JSON**: Use the built-in JSON methods for type safety
- **Base URL**: Use base URLs for consistent API endpoints

## Dependencies

- No external dependencies
- Uses only Go standard library (`net/http`, `context`, `time`, `encoding/json`)