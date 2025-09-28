package httpclient

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
	
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.httpClient.Timeout)
	}
	
	if client.headers == nil {
		t.Error("Headers map should not be nil")
	}
	
	if len(client.headers) != 0 {
		t.Error("Headers map should be empty initially")
	}
}

func TestNewWithBaseURL(t *testing.T) {
	baseURL := "https://api.example.com"
	client := NewWithBaseURL(baseURL)
	
	if client == nil {
		t.Fatal("NewWithBaseURL() returned nil")
	}
	
	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", baseURL, client.baseURL)
	}
}

func TestNewWithTimeout(t *testing.T) {
	timeout := 10 * time.Second
	client := NewWithTimeout(timeout)
	
	if client == nil {
		t.Fatal("NewWithTimeout() returned nil")
	}
	
	if client.httpClient.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.httpClient.Timeout)
	}
}

func TestSetBaseURL(t *testing.T) {
	client := New()
	baseURL := "https://api.example.com"
	
	client.SetBaseURL(baseURL)
	
	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", baseURL, client.baseURL)
	}
}

func TestGetBaseURL(t *testing.T) {
	client := New()
	baseURL := "https://api.example.com"
	
	client.SetBaseURL(baseURL)
	result := client.GetBaseURL()
	
	if result != baseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", baseURL, result)
	}
}

func TestSetTimeout(t *testing.T) {
	client := New()
	timeout := 15 * time.Second
	
	client.SetTimeout(timeout)
	
	if client.httpClient.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.httpClient.Timeout)
	}
}

func TestGetTimeout(t *testing.T) {
	client := New()
	timeout := 20 * time.Second
	
	client.SetTimeout(timeout)
	result := client.GetTimeout()
	
	if result != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, result)
	}
}

func TestSetHeader(t *testing.T) {
	client := New()
	key := "Authorization"
	value := "Bearer token123"
	
	client.SetHeader(key, value)
	
	if client.headers[key] != value {
		t.Errorf("Expected header '%s'='%s', got '%s'", key, value, client.headers[key])
	}
}

func TestGetHeader(t *testing.T) {
	client := New()
	key := "Authorization"
	value := "Bearer token123"
	
	client.SetHeader(key, value)
	result := client.GetHeader(key)
	
	if result != value {
		t.Errorf("Expected header '%s'='%s', got '%s'", key, value, result)
	}
	
	// Test non-existent header
	result2 := client.GetHeader("NonExistent")
	if result2 != "" {
		t.Errorf("Expected empty string for non-existent header, got '%s'", result2)
	}
}

func TestGetHeaders(t *testing.T) {
	client := New()
	
	// Initially should be empty
	headers := client.GetHeaders()
	if len(headers) != 0 {
		t.Error("Headers should be empty initially")
	}
	
	// Add headers
	client.SetHeader("Authorization", "Bearer token123")
	client.SetHeader("Content-Type", "application/json")
	
	headers = client.GetHeaders()
	if len(headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(headers))
	}
	
	if headers["Authorization"] != "Bearer token123" {
		t.Error("Authorization header not found")
	}
	
	if headers["Content-Type"] != "application/json" {
		t.Error("Content-Type header not found")
	}
}

func TestSetHeaders(t *testing.T) {
	client := New()
	headers := map[string]string{
		"Authorization": "Bearer token123",
		"Content-Type":  "application/json",
		"User-Agent":    "MyApp/1.0",
	}
	
	client.SetHeaders(headers)
	
	result := client.GetHeaders()
	if len(result) != len(headers) {
		t.Errorf("Expected %d headers, got %d", len(headers), len(result))
	}
	
	for key, value := range headers {
		if result[key] != value {
			t.Errorf("Expected header '%s'='%s', got '%s'", key, value, result[key])
		}
	}
}

func TestSetHTTPClient(t *testing.T) {
	client := New()
	newHTTPClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	
	client.SetHTTPClient(newHTTPClient)
	
	if client.httpClient != newHTTPClient {
		t.Error("HTTP client was not set correctly")
	}
}

func TestGetHTTPClient(t *testing.T) {
	client := New()
	newHTTPClient := &http.Client{
		Timeout: 60 * time.Second,
	}
	
	client.SetHTTPClient(newHTTPClient)
	result := client.GetHTTPClient()
	
	if result != newHTTPClient {
		t.Error("GetHTTPClient did not return the set client")
	}
}

func TestGet(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	resp, err := client.Get("/test")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	
	if string(resp) != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got '%s'", string(resp))
	}
}

func TestGetWithContext(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	ctx := context.Background()
	resp, err := client.GetWithContext(ctx, "/test")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	
	if string(resp) != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got '%s'", string(resp))
	}
}

func TestPost(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		
		// Read body
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		
		if string(body) != `{"message":"Hello"}` {
			t.Errorf("Expected body '{\"message\":\"Hello\"}', got '%s'", string(body))
		}
		
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	data := []byte(`{"message":"Hello"}`)
	resp, err := client.Post("/test", data)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	
	if string(resp) != "Created" {
		t.Errorf("Expected body 'Created', got '%s'", string(resp))
	}
}

func TestPostWithContext(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	data := []byte(`{"message":"Hello"}`)
	ctx := context.Background()
	resp, err := client.PostWithContext(ctx, "/test", data)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	
	if string(resp) != "Created" {
		t.Errorf("Expected body 'Created', got '%s'", string(resp))
	}
}

func TestPut(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Updated"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	data := []byte(`{"message":"Updated"}`)
	resp, err := client.Put("/test", data)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	
	if string(resp) != "Updated" {
		t.Errorf("Expected body 'Updated', got '%s'", string(resp))
	}
}

func TestDelete(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	resp, err := client.Delete("/test")
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	
	// DELETE with 204 should return empty body
	if len(resp) != 0 {
		t.Errorf("Expected empty body, got '%s'", string(resp))
	}
}

func TestPatch(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Patched"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	data := []byte(`{"message":"Patched"}`)
	resp, err := client.Patch("/test", data)
	if err != nil {
		t.Fatalf("PATCH request failed: %v", err)
	}
	
	if string(resp) != "Patched" {
		t.Errorf("Expected body 'Patched', got '%s'", string(resp))
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
	
	client := New()
	client.SetBaseURL(server.URL)
	
	resp, err := client.Request("OPTIONS", "/test", nil)
	if err != nil {
		t.Fatalf("OPTIONS request failed: %v", err)
	}
	
	if string(resp) != "Options" {
		t.Errorf("Expected body 'Options', got '%s'", string(resp))
	}
}

func TestPostJSON(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("JSON Posted"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	data := map[string]string{"message": "Hello"}
	resp, err := client.PostJSON("/test", data)
	if err != nil {
		t.Fatalf("PostJSON request failed: %v", err)
	}
	
	if string(resp) != "JSON Posted" {
		t.Errorf("Expected body 'JSON Posted', got '%s'", string(resp))
	}
}

func TestPutJSON(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("JSON Put"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	data := map[string]string{"message": "Updated"}
	resp, err := client.PutJSON("/test", data)
	if err != nil {
		t.Fatalf("PutJSON request failed: %v", err)
	}
	
	if string(resp) != "JSON Put" {
		t.Errorf("Expected body 'JSON Put', got '%s'", string(resp))
	}
}

func TestPatchJSON(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("JSON Patched"))
	}))
	defer server.Close()
	
	client := New()
	client.SetBaseURL(server.URL)
	
	data := map[string]string{"message": "Patched"}
	resp, err := client.PatchJSON("/test", data)
	if err != nil {
		t.Fatalf("PatchJSON request failed: %v", err)
	}
	
	if string(resp) != "JSON Patched" {
		t.Errorf("Expected body 'JSON Patched', got '%s'", string(resp))
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
	
	client := New()
	client.SetBaseURL(server.URL)
	client.SetTimeout(50 * time.Millisecond) // Shorter timeout
	
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	
	_, err := client.GetWithContext(ctx, "/test")
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

func BenchmarkNewWithBaseURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client := NewWithBaseURL("https://api.example.com")
		_ = client
	}
}

func BenchmarkSetHeader(b *testing.B) {
	client := New()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.SetHeader("Authorization", "Bearer token123")
	}
}

func BenchmarkGetHeaders(b *testing.B) {
	client := New()
	client.SetHeader("Authorization", "Bearer token123")
	client.SetHeader("Content-Type", "application/json")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.GetHeaders()
	}
}