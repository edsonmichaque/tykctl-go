package api

import (
	"context"
	"fmt"
	"time"
)

// Example demonstrates basic API client usage
func Example() {
	// Create API client with functional options
	client := New(
		WithBaseURL("https://api.example.com/v1"),
		WithClientTimeout(30*time.Second),
		WithClientHeader("Accept", "application/json"),
		WithClientHeader("Authorization", "Bearer your-api-token"),
	)
	ctx := context.Background()

	// Make a GET request with functional options
	resp, err := client.Get(ctx, "/users",
		WithHeader("Accept", "application/json"),
		WithQuery("page", "1"),
		WithQuery("limit", "10"),
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

// ExampleWithRetry demonstrates retry functionality
func ExampleWithRetry() {
	client := New(
		WithBaseURL("https://api.example.com/v1"),
		WithClientTimeout(30 * time.Second),
	)
	ctx := context.Background()

	// Create retry config with exponential backoff
	retryConfig := NewExponentialBackoffConfig(
		3,              // max retries
		1*time.Second,  // initial delay
		10*time.Second, // max delay
		30*time.Second, // max elapsed time
	)

	// Make request with retry
	resp, err := WithRetry(ctx, retryConfig, func() (*Response, error) {
		return client.Get(ctx, "/users")
	})

	if err != nil {
		fmt.Printf("Error after retries: %v\n", err)
		return
	}

	fmt.Printf("Success: %s\n", resp.String())
}

// ExampleWithMiddleware demonstrates middleware usage
func ExampleWithMiddleware() {
	client := New(
		WithBaseURL("https://api.example.com/v1"),
		WithClientTimeout(30 * time.Second),
	)
	ctx := context.Background()

	// Create middleware chain
	middlewares := ChainMiddleware(
		LoggingMiddleware(nil),
		TimeoutMiddleware(10*time.Second),
		HeaderMiddleware(map[string]string{
			"User-Agent": "my-app/1.0.0",
		}),
	)

	// This would be used with a custom request handler
	// For demonstration, we'll just show the middleware setup
	_ = middlewares

	// Make request
	resp, err := client.Get(ctx, "/users")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Success: %s\n", resp.String())
}

// ExamplePaginatedRequest demonstrates paginated requests
func ExamplePaginatedRequest() {
	client := New(
		WithBaseURL("https://api.example.com/v1"),
		WithClientTimeout(30 * time.Second),
	)
	ctx := context.Background()

	// Make paginated request
	resp, err := client.Get(ctx, "/users?page=1&per_page=10")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Parse paginated response
	var paginatedResp PaginatedResponse
	if err := resp.UnmarshalJSON(&paginatedResp); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	fmt.Printf("Total users: %d\n", paginatedResp.Pagination.Total)
	fmt.Printf("Current page: %d\n", paginatedResp.Pagination.Page)
	fmt.Printf("Has next page: %t\n", paginatedResp.Pagination.HasNext)
}

// ExampleCustomRetryCondition demonstrates custom retry logic
func ExampleCustomRetryCondition() {
	// Custom retry condition that retries on specific status codes
	customRetryCondition := &CustomRetryCondition{
		RetryableStatusCodes: []int{429, 500, 502, 503, 504},
	}

	retryConfig := RetryConfig{
		MaxRetries:     5,
		InitialDelay:   2 * time.Second,
		MaxDelay:       2 * time.Second, // Constant delay
		MaxElapsedTime: 30 * time.Second,
		Multiplier:     1.0, // No exponential growth
		Retryable:      customRetryCondition,
	}

	client := New(
		WithBaseURL("https://api.example.com/v1"),
		WithClientTimeout(30 * time.Second),
	)
	ctx := context.Background()

	// Make request with custom retry logic
	resp, err := WithRetry(ctx, retryConfig, func() (*Response, error) {
		return client.Get(ctx, "/users")
	})

	if err != nil {
		fmt.Printf("Error after retries: %v\n", err)
		return
	}

	fmt.Printf("Success: %s\n", resp.String())
}

// CustomRetryCondition implements custom retry logic
type CustomRetryCondition struct {
	RetryableStatusCodes []int
}

func (c *CustomRetryCondition) ShouldRetry(err error, response *Response) bool {
	if err != nil {
		return true // Retry on network errors
	}

	if response == nil {
		return false
	}

	// Check if status code is in retryable list
	for _, code := range c.RetryableStatusCodes {
		if response.StatusCode == code {
			return true
		}
	}

	return false
}

// ExampleClientFlexibility demonstrates client flexibility and general-purpose usage
func ExampleClientFlexibility() {
	// Create base client
	client := New(
		WithBaseURL("https://api.example.com/v1"),
		WithClientTimeout(30 * time.Second),
		WithClientHeaders(map[string]string{
			"Content-Type": "application/json",
		}),
	)
	ctx := context.Background()

	// Clone client for different use cases
	userClient := client.Clone()
	userClient.SetBaseURL("https://api.example.com/v1/users")

	// Create client with different timeout
	fastClient := client.WithTimeout(5 * time.Second)

	// Create client with additional headers
	authClient := client.WithHeader("Authorization", "Bearer token123")

	// Create client for different API version
	v2Client := client.WithBaseURL("https://api.example.com/v2")

	// Use different clients for different purposes
	fmt.Println("Base client:", client.GetConfig().BaseURL)
	fmt.Println("User client:", userClient.GetConfig().BaseURL)
	fmt.Println("Fast client timeout:", fastClient.GetConfig().Timeout)
	fmt.Println("Auth client headers:", authClient.GetConfig().Headers)
	fmt.Println("V2 client:", v2Client.GetConfig().BaseURL)

	// Make requests with different clients
	resp1, err := client.Get(ctx, "/status")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Status: %d\n", resp1.StatusCode)

	resp2, err := userClient.Get(ctx, "/profile")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("User profile: %d\n", resp2.StatusCode)
}

// ExampleFunctionalOptions demonstrates the new functional options for requests
func ExampleFunctionalOptions() {
	client := New(
		WithBaseURL("https://api.example.com/v1"),
		WithClientTimeout(30 * time.Second),
	)
	ctx := context.Background()

	// GET request with headers and query parameters
	resp1, err := client.Get(ctx, "/users",
		WithHeader("Accept", "application/json"),
		WithHeader("X-API-Version", "v1"),
		WithQuery("page", "1"),
		WithQuery("limit", "20"),
		WithQuery("sort", "name"),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("GET with options: %d\n", resp1.StatusCode)

	// POST request with JSON body and custom headers
	userData := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	resp2, err := client.Post(ctx, "/users", userData,
		WithHeader("Content-Type", "application/json"),
		WithHeader("X-Request-ID", "req-123"),
		WithQuery("validate", "true"),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("POST with options: %d\n", resp2.StatusCode)

	// PUT request with custom body and timeout
	resp3, err := client.Put(ctx, "/users/123", nil,
		WithJSONBody(map[string]string{"status": "active"}),
		WithHeader("If-Match", "etag-value"),
		WithTimeout(10*time.Second),
		WithRetries(3),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("PUT with options: %d\n", resp3.StatusCode)

	// DELETE request with query parameters
	resp4, err := client.Delete(ctx, "/users/123",
		WithQuery("force", "true"),
		WithHeader("X-Confirm-Delete", "yes"),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("DELETE with options: %d\n", resp4.StatusCode)

	// PATCH request with multiple headers
	resp5, err := client.Patch(ctx, "/users/123", map[string]string{"name": "Jane Doe"},
		WithHeaders(map[string]string{
			"Content-Type":    "application/json",
			"X-Update-Source": "admin-panel",
		}),
		WithQueries(map[string]string{
			"notify": "true",
			"audit":  "true",
		}),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("PATCH with options: %d\n", resp5.StatusCode)
}
