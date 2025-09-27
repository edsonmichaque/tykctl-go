// Package httpclient provides a simple HTTP client for making API requests with context support.
//
// Features:
//   - HTTP Methods: Support for GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
//   - Context Support: Full context.Context integration for cancellation and timeouts
//   - Header Management: Easy header setting and management
//   - JSON Support: Built-in JSON marshaling and unmarshaling
//   - Response Handling: Rich response objects with status code checking
//   - Error Handling: Comprehensive error handling with HTTP status codes
//   - Base URL Support: Configurable base URLs for API endpoints
//
// Example:
//   client := httpclient.NewWithBaseURL("https://api.example.com")
//   data, err := client.Get("/users")
//   err = client.PostJSON("/users", userData)
package httpclient