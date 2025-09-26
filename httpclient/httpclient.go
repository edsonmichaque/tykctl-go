package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an HTTP client
type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// New creates a new HTTP client
func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

// NewWithBaseURL creates a new HTTP client with a base URL
func NewWithBaseURL(baseURL string) *Client {
	client := New()
	client.baseURL = baseURL
	return client
}

// NewWithTimeout creates a new HTTP client with a custom timeout
func NewWithTimeout(timeout time.Duration) *Client {
	client := New()
	client.httpClient.Timeout = timeout
	return client
}

// SetBaseURL sets the base URL
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// GetBaseURL returns the base URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// SetTimeout sets the client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// GetTimeout returns the client timeout
func (c *Client) GetTimeout() time.Duration {
	return c.httpClient.Timeout
}

// SetHeader sets a header
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// GetHeader returns a header value
func (c *Client) GetHeader(key string) string {
	return c.headers[key]
}

// SetHeaders sets multiple headers
func (c *Client) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		c.headers[k] = v
	}
}

// GetHeaders returns all headers
func (c *Client) GetHeaders() map[string]string {
	return c.headers
}

// SetAuthorization sets the Authorization header
func (c *Client) SetAuthorization(auth string) {
	c.SetHeader("Authorization", auth)
}

// GetAuthorization returns the Authorization header value
func (c *Client) GetAuthorization() string {
	return c.GetHeader("Authorization")
}

// Get makes a GET request
func (c *Client) Get(path string) ([]byte, error) {
	return c.GetWithContext(context.Background(), path)
}

// GetWithContext makes a GET request with context
func (c *Client) GetWithContext(ctx context.Context, path string) ([]byte, error) {
	req, err := c.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// Post makes a POST request
func (c *Client) Post(path string, data []byte) ([]byte, error) {
	return c.PostWithContext(context.Background(), path, data)
}

// PostWithContext makes a POST request with context
func (c *Client) PostWithContext(ctx context.Context, path string, data []byte) ([]byte, error) {
	req, err := c.newRequest(ctx, "POST", path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// Put makes a PUT request
func (c *Client) Put(path string, data []byte) ([]byte, error) {
	return c.PutWithContext(context.Background(), path, data)
}

// PutWithContext makes a PUT request with context
func (c *Client) PutWithContext(ctx context.Context, path string, data []byte) ([]byte, error) {
	req, err := c.newRequest(ctx, "PUT", path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// Delete makes a DELETE request
func (c *Client) Delete(path string) ([]byte, error) {
	return c.DeleteWithContext(context.Background(), path)
}

// DeleteWithContext makes a DELETE request with context
func (c *Client) DeleteWithContext(ctx context.Context, path string) ([]byte, error) {
	req, err := c.newRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// Patch makes a PATCH request
func (c *Client) Patch(path string, data []byte) ([]byte, error) {
	return c.PatchWithContext(context.Background(), path, data)
}

// PatchWithContext makes a PATCH request with context
func (c *Client) PatchWithContext(ctx context.Context, path string, data []byte) ([]byte, error) {
	req, err := c.newRequest(ctx, "PATCH", path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// newRequest creates a new HTTP request
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Set headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Set content type for POST/PUT/PATCH requests
	if body != nil && (method == "POST" || method == "PUT" || method == "PATCH") {
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	return req, nil
}

// doRequest executes an HTTP request
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetJSON makes a GET request and unmarshals the response
func (c *Client) GetJSON(path string, v interface{}) error {
	data, err := c.Get(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// PostJSON makes a POST request with JSON data
func (c *Client) PostJSON(path string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return c.Post(path, jsonData)
}

// PostJSONResponse makes a POST request with JSON data and unmarshals the response
func (c *Client) PostJSONResponse(path string, data interface{}, response interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := c.Post(path, jsonData)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, response)
}

// PutJSON makes a PUT request with JSON data
func (c *Client) PutJSON(path string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return c.Put(path, jsonData)
}

// PatchJSON makes a PATCH request with JSON data
func (c *Client) PatchJSON(path string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return c.Patch(path, jsonData)
}

// SetUserAgent sets the User-Agent header
func (c *Client) SetUserAgent(userAgent string) {
	c.SetHeader("User-Agent", userAgent)
}

// SetContentType sets the Content-Type header
func (c *Client) SetContentType(contentType string) {
	c.SetHeader("Content-Type", contentType)
}

// SetAccept sets the Accept header
func (c *Client) SetAccept(accept string) {
	c.SetHeader("Accept", accept)
}

// SetCustomHeader sets a custom header
func (c *Client) SetCustomHeader(key, value string) {
	c.SetHeader(key, value)
}

// GetCustomHeader returns a custom header value
func (c *Client) GetCustomHeader(key string) string {
	return c.GetHeader(key)
}

// RemoveHeader removes a header
func (c *Client) RemoveHeader(key string) {
	delete(c.headers, key)
}

// ClearHeaders clears all headers
func (c *Client) ClearHeaders() {
	c.headers = make(map[string]string)
}

// GetHTTPClient returns the underlying HTTP client
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// SetHTTPClient sets the underlying HTTP client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// Request makes a generic HTTP request
func (c *Client) Request(method, path string, body []byte) ([]byte, error) {
	return c.RequestWithContext(context.Background(), method, path, body)
}

// RequestWithContext makes a generic HTTP request with context
func (c *Client) RequestWithContext(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := c.newRequest(ctx, method, path, reader)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// RequestWithHeaders makes a generic HTTP request with additional headers
func (c *Client) RequestWithHeaders(method, path string, body []byte, headers map[string]string) ([]byte, error) {
	return c.RequestHeadersWithContext(context.Background(), method, path, body, headers)
}

// RequestHeadersWithContext makes a generic HTTP request with additional headers and context
func (c *Client) RequestHeadersWithContext(ctx context.Context, method, path string, body []byte, headers map[string]string) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := c.newRequest(ctx, method, path, reader)
	if err != nil {
		return nil, err
	}

	// Add additional headers for this request
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.doRequest(req)
}

// Head makes a HEAD request
func (c *Client) Head(path string) error {
	return c.HeadWithContext(context.Background(), path)
}

// HeadWithContext makes a HEAD request with context
func (c *Client) HeadWithContext(ctx context.Context, path string) error {
	req, err := c.newRequest(ctx, "HEAD", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

// Options makes an OPTIONS request
func (c *Client) Options(path string) ([]byte, error) {
	return c.OptionsWithContext(context.Background(), path)
}

// OptionsWithContext makes an OPTIONS request with context
func (c *Client) OptionsWithContext(ctx context.Context, path string) ([]byte, error) {
	req, err := c.newRequest(ctx, "OPTIONS", path, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req)
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// GetResponse makes a GET request and returns the full response
func (c *Client) GetResponse(path string) (*Response, error) {
	return c.GetResponseWithContext(context.Background(), path)
}

// GetResponseWithContext makes a GET request with context and returns the full response
func (c *Client) GetResponseWithContext(ctx context.Context, path string) (*Response, error) {
	req, err := c.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequestWithResponse(req)
}

// PostResponse makes a POST request and returns the full response
func (c *Client) PostResponse(path string, data []byte) (*Response, error) {
	return c.PostResponseWithContext(context.Background(), path, data)
}

// PostResponseWithContext makes a POST request with context and returns the full response
func (c *Client) PostResponseWithContext(ctx context.Context, path string, data []byte) (*Response, error) {
	req, err := c.newRequest(ctx, "POST", path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return c.doRequestWithResponse(req)
}

// doRequestWithResponse executes an HTTP request and returns the full response
func (c *Client) doRequestWithResponse(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert headers to map
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       body,
	}, nil
}

// IsSuccess checks if the status code indicates success
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError checks if the status code indicates a client error
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError checks if the status code indicates a server error
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// GetHeader returns a header value
func (r *Response) GetHeader(key string) string {
	return r.Headers[key]
}
