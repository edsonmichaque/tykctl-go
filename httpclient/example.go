package httpclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Example demonstrates the general-purpose HTTP client usage
func Example() {
	// Create a new HTTP client
	client := New()

	// Set base URL
	client.SetBaseURL("https://httpbin.org")

	// Set common headers
	client.SetUserAgent("tykctl-go/1.0.0")
	client.SetAccept("application/json")

	// Example 1: Basic GET request
	fmt.Println("=== Basic GET Request ===")
	data, err := client.Get("/get")
	if err != nil {
		log.Printf("GET request failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", string(data))
	}

	// Example 2: GET request with context and timeout
	fmt.Println("\n=== GET Request with Context ===")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err = client.GetWithContext(ctx, "/get")
	if err != nil {
		log.Printf("GET request with context failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", string(data))
	}

	// Example 3: POST request with JSON data
	fmt.Println("\n=== POST Request with JSON ===")
	jsonData := map[string]interface{}{
		"name":    "tykctl-go",
		"version": "1.0.0",
	}

	response, err := client.PostJSON("/post", jsonData)
	if err != nil {
		log.Printf("POST JSON request failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", string(response))
	}

	// Example 4: Generic request method
	fmt.Println("\n=== Generic Request Method ===")
	requestData := []byte(`{"test": "data"}`)
	data, err = client.Request("POST", "/post", requestData)
	if err != nil {
		log.Printf("Generic request failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", string(data))
	}

	// Example 5: Request with additional headers
	fmt.Println("\n=== Request with Additional Headers ===")
	extraHeaders := map[string]string{
		"X-Custom-Header": "custom-value",
		"X-Request-ID":    "12345",
	}

	data, err = client.RequestWithHeaders("GET", "/headers", nil, extraHeaders)
	if err != nil {
		log.Printf("Request with headers failed: %v", err)
	} else {
		fmt.Printf("Response: %s\n", string(data))
	}

	// Example 6: Full response handling
	fmt.Println("\n=== Full Response Handling ===")
	resp, err := client.GetResponse("/get")
	if err != nil {
		log.Printf("GET with response failed: %v", err)
	} else {
		fmt.Printf("Status Code: %d\n", resp.StatusCode)
		fmt.Printf("Is Success: %t\n", resp.IsSuccess())
		fmt.Printf("Content-Type: %s\n", resp.GetHeader("Content-Type"))
		fmt.Printf("Response Body: %s\n", string(resp.Body))
	}

	// Example 7: Different authentication methods
	fmt.Println("\n=== Authentication Examples ===")

	// Bearer token authentication
	client.SetAuthorization("Bearer your-token-here")
	fmt.Println("Set Bearer token authentication")

	// Basic authentication
	client.SetAuthorization("Basic dXNlcm5hbWU6cGFzc3dvcmQ=")
	fmt.Println("Set Basic authentication")

	// Custom authentication
	client.SetAuthorization("Custom auth-scheme your-custom-token")
	fmt.Println("Set custom authentication")

	// Example 8: Different HTTP methods
	fmt.Println("\n=== Different HTTP Methods ===")

	// HEAD request
	err = client.Head("/get")
	if err != nil {
		log.Printf("HEAD request failed: %v", err)
	} else {
		fmt.Println("HEAD request successful")
	}

	// OPTIONS request
	data, err = client.Options("/get")
	if err != nil {
		log.Printf("OPTIONS request failed: %v", err)
	} else {
		fmt.Printf("OPTIONS response: %s\n", string(data))
	}

	// PATCH request
	patchData := []byte(`{"status": "updated"}`)
	data, err = client.Patch("/patch", patchData)
	if err != nil {
		log.Printf("PATCH request failed: %v", err)
	} else {
		fmt.Printf("PATCH response: %s\n", string(data))
	}

	fmt.Println("\nHTTP client example completed successfully!")
}

// AuthenticationExample demonstrates different authentication patterns
func AuthenticationExample() {
	client := New()
	client.SetBaseURL("https://api.example.com")

	// Example 1: Bearer Token
	client.SetAuthorization("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")

	// Example 2: Basic Auth
	client.SetAuthorization("Basic dXNlcm5hbWU6cGFzc3dvcmQ=")

	// Example 3: API Key in header
	client.SetHeader("X-API-Key", "your-api-key-here")

	// Example 4: Custom authentication scheme
	client.SetAuthorization("Digest username=\"user\", realm=\"realm\", nonce=\"nonce\"")

	// Example 5: Multiple authentication headers
	client.SetHeader("X-Auth-Token", "token-value")
	client.SetHeader("X-Auth-Signature", "signature-value")

	fmt.Println("Authentication examples configured")
}

// AdvancedExample demonstrates advanced HTTP client features
func AdvancedExample() {
	// Create client with custom timeout
	client := NewWithTimeout(10 * time.Second)
	client.SetBaseURL("https://httpbin.org")

	// Set up custom headers
	client.SetHeaders(map[string]string{
		"User-Agent":      "tykctl-go/1.0.0",
		"Accept":          "application/json",
		"Accept-Encoding": "gzip, deflate",
	})

	// Create a custom HTTP client with additional configuration
	customClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
	}
	client.SetHTTPClient(customClient)

	// Make a request with the custom configuration
	data, err := client.Get("/get")
	if err != nil {
		log.Printf("Advanced request failed: %v", err)
	} else {
		fmt.Printf("Advanced response: %s\n", string(data))
	}
}
