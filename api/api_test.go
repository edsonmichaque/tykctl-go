package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New()
	if client == nil {
		t.Fatal("New() returned nil")
	}
	
	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}
	
	if client.config == nil {
		t.Error("Config should not be nil")
	}
}

func TestNewWithOptions(t *testing.T) {
	baseURL := "https://api.example.com"
	timeout := 15 * time.Second
	
	client := New(
		WithBaseURL(baseURL),
		WithClientTimeout(timeout),
		WithClientHeader("Authorization", "Bearer token123"),
	)
	
	if client == nil {
		t.Fatal("New() returned nil")
	}
	
	if client.BaseURL != baseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", baseURL, client.BaseURL)
	}
	
	if client.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.Timeout)
	}
	
	if client.config.Headers["Authorization"] != "Bearer token123" {
		t.Error("Authorization header not set correctly")
	}
}

func TestWithBaseURL(t *testing.T) {
	baseURL := "https://api.example.com"
	option := WithBaseURL(baseURL)
	
	config := &ClientConfig{}
	option(config)
	
	if config.BaseURL != baseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", baseURL, config.BaseURL)
	}
}

func TestWithClientTimeout(t *testing.T) {
	timeout := 30 * time.Second
	option := WithClientTimeout(timeout)
	
	config := &ClientConfig{}
	option(config)
	
	if config.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, config.Timeout)
	}
}

func TestWithClientHeader(t *testing.T) {
	key := "Authorization"
	value := "Bearer token123"
	option := WithClientHeader(key, value)
	
	config := &ClientConfig{}
	option(config)
	
	if config.Headers == nil {
		t.Error("Headers map should be initialized")
	}
	
	if config.Headers[key] != value {
		t.Errorf("Expected header '%s'='%s', got '%s'", key, value, config.Headers[key])
	}
}

func TestWithClientHeaders(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token123",
		"Content-Type":  "application/json",
	}
	option := WithClientHeaders(headers)
	
	config := &ClientConfig{}
	option(config)
	
	if config.Headers == nil {
		t.Error("Headers map should be initialized")
	}
	
	for key, value := range headers {
		if config.Headers[key] != value {
			t.Errorf("Expected header '%s'='%s', got '%s'", key, value, config.Headers[key])
		}
	}
}

func TestWithUserAgent(t *testing.T) {
	userAgent := "MyApp/1.0.0"
	option := WithUserAgent(userAgent)
	
	config := &ClientConfig{}
	option(config)
	
	if config.UserAgent != userAgent {
		t.Errorf("Expected user agent '%s', got '%s'", userAgent, config.UserAgent)
	}
}

func TestGet(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello"}`))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	
	resp, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	if resp.String() != `{"message":"Hello"}` {
		t.Errorf("Expected body '{\"message\":\"Hello\"}', got '%s'", resp.String())
	}
}

func TestPost(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":123}`))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	
	data := map[string]string{"name": "Test"}
	resp, err := client.Post(context.Background(), "/test", data)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
	
	if resp.String() != `{"id":123}` {
		t.Errorf("Expected body '{\"id\":123}', got '%s'", resp.String())
	}
}

func TestPut(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated":true}`))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	
	data := map[string]string{"name": "Updated"}
	resp, err := client.Put(context.Background(), "/test", data)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	if resp.String() != `{"updated":true}` {
		t.Errorf("Expected body '{\"updated\":true}', got '%s'", resp.String())
	}
}

func TestDelete(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK) // Some APIs return 200 for DELETE
		w.Write([]byte("Deleted"))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	
	resp, err := client.Delete(context.Background(), "/test")
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	if resp.String() != "Deleted" {
		t.Errorf("Expected body 'Deleted', got '%s'", resp.String())
	}
}

func TestPatch(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"patched":true}`))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	
	data := map[string]string{"name": "Patched"}
	resp, err := client.Patch(context.Background(), "/test", data)
	if err != nil {
		t.Fatalf("PATCH request failed: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	if resp.String() != `{"patched":true}` {
		t.Errorf("Expected body '{\"patched\":true}', got '%s'", resp.String())
	}
}

func TestRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			t.Errorf("Expected OPTIONS method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Options"))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	
	resp, err := client.Request(context.Background(), "OPTIONS", "/test", nil)
	if err != nil {
		t.Fatalf("OPTIONS request failed: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	
	if resp.String() != "Options" {
		t.Errorf("Expected body 'Options', got '%s'", resp.String())
	}
}

func TestResponseMethods(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "custom-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello"}`))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	
	resp, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	
	// Test IsSuccess
	if !resp.IsSuccess() {
		t.Error("Response should be successful")
	}
	
	// Test IsClientError
	if resp.IsClientError() {
		t.Error("Response should not be client error")
	}
	
	// Test IsServerError
	if resp.IsServerError() {
		t.Error("Response should not be server error")
	}
	
	// Test GetHeader
	contentType := resp.GetHeader("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
	
	customHeader := resp.GetHeader("X-Custom-Header")
	if customHeader != "custom-value" {
		t.Errorf("Expected X-Custom-Header 'custom-value', got '%s'", customHeader)
	}
	
	// Test Headers field
	headers := resp.Headers
	if headers["Content-Type"] != "application/json" {
		t.Error("Content-Type header not found in headers")
	}
	
	// Test UnmarshalJSON
	var data map[string]string
	err = resp.UnmarshalJSON(&data)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	if data["message"] != "Hello" {
		t.Errorf("Expected message 'Hello', got '%s'", data["message"])
	}
}

func TestErrorResponse(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	
	resp, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	
	// Test error response methods
	if resp.IsSuccess() {
		t.Error("Response should not be successful")
	}
	
	if !resp.IsClientError() {
		t.Error("Response should be client error")
	}
	
	if resp.IsServerError() {
		t.Error("Response should not be server error")
	}
	
	if resp.String() != "Not Found" {
		t.Errorf("Expected body 'Not Found', got '%s'", resp.String())
	}
}

func TestContextCancellation(t *testing.T) {
	// Create test server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL), WithClientTimeout(50*time.Millisecond))
	
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	
	_, err := client.Get(ctx, "/test")
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

// Benchmark tests
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client := New()
		_ = client
	}
}

func BenchmarkNewWithOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client := New(
			WithBaseURL("https://api.example.com"),
			WithClientTimeout(30*time.Second),
			WithClientHeader("Authorization", "Bearer token123"),
		)
		_ = client
	}
}

func BenchmarkGetRequest(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello"}`))
	}))
	defer server.Close()
	
	client := New(WithBaseURL(server.URL))
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(ctx, "/test")
		if err != nil {
			b.Fatalf("GET request failed: %v", err)
		}
		_ = resp
	}
}