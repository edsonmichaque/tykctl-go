// Package api provides a high-level HTTP client for making API requests with built-in retry logic, middleware support, and response handling.
//
// Features:
//   - HTTP Client: Simple interface for GET, POST, PUT, DELETE, PATCH requests
//   - Retry Logic: Configurable retry strategies with exponential backoff
//   - Middleware Support: Chainable middleware for logging, timeouts, authentication
//   - Response Handling: Rich response objects with status code checking and JSON unmarshaling
//   - Pagination Support: Built-in pagination handling for API responses
//   - Error Handling: Comprehensive error types and handling with retryable error detection
//   - Context Support: Full context.Context integration for cancellation and timeouts
//   - Configurable: Flexible configuration options for different use cases
//
// Example:
//   client := api.New(
//       api.WithBaseURL("https://api.example.com"),
//       api.WithClientTimeout(30*time.Second),
//   )
//
//   resp, err := client.Get(ctx, "/users",
//       api.WithHeader("X-API-Version", "v1"),
//       api.WithQuery("page", "1"),
//   )
package api