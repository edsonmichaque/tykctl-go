package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/edsonmichaque/tykctl-go/httpclient"
)

// Client represents an API client
type Client struct {
	httpClient *httpclient.Client
	BaseURL    string
	Timeout    time.Duration
	config     *ClientConfig
}

// ClientOption is a functional option for configuring the API client
type ClientOption func(*ClientConfig)

// ClientConfig holds the configuration for the API client
type ClientConfig struct {
	BaseURL   string
	Timeout   time.Duration
	Headers   map[string]string
	UserAgent string
}

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *ClientConfig) {
		c.BaseURL = baseURL
	}
}

// WithClientTimeout sets the timeout for the client
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(c *ClientConfig) {
		c.Timeout = timeout
	}
}

// WithClientHeader adds a header to the client
func WithClientHeader(key, value string) ClientOption {
	return func(c *ClientConfig) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		c.Headers[key] = value
	}
}

// WithClientHeaders adds multiple headers to the client
func WithClientHeaders(headers map[string]string) ClientOption {
	return func(c *ClientConfig) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

// WithUserAgent sets the user agent for the client
func WithUserAgent(userAgent string) ClientOption {
	return func(c *ClientConfig) {
		c.UserAgent = userAgent
	}
}

// New creates a new API client with functional options
func New(opts ...ClientOption) *Client {
	// Create default configuration
	config := &ClientConfig{
		Timeout:   30 * time.Second,
		Headers:   make(map[string]string),
		UserAgent: "tykctl-go/1.0.0",
	}

	// Apply functional options
	for _, opt := range opts {
		opt(config)
	}

	// Create HTTP client
	client := httpclient.New()

	if config.BaseURL != "" {
		client.SetBaseURL(config.BaseURL)
	}

	if config.Timeout > 0 {
		client.SetTimeout(config.Timeout)
	} else {
		client.SetTimeout(30 * time.Second)
	}

	if config.UserAgent != "" {
		client.SetUserAgent(config.UserAgent)
	}

	// Set custom headers
	if config.Headers != nil {
		for k, v := range config.Headers {
			client.SetHeader(k, v)
		}
	}

	return &Client{
		httpClient: client,
		BaseURL:    config.BaseURL,
		Timeout:    config.Timeout,
		config:     config,
	}
}

// SetHTTPClient allows setting a custom HTTP client
func (c *Client) SetHTTPClient(client *httpclient.Client) {
	c.httpClient = client
}

// GetHTTPClient returns the underlying HTTP client
func (c *Client) GetHTTPClient() *httpclient.Client {
	return c.httpClient
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *ClientConfig {
	return c.config
}

// SetBaseURL updates the base URL
func (c *Client) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
	c.httpClient.SetBaseURL(baseURL)
}

// SetTimeout updates the timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.Timeout = timeout
	c.httpClient.SetTimeout(timeout)
}

// SetHeader sets a header
func (c *Client) SetHeader(key, value string) {
	c.httpClient.SetHeader(key, value)
}

// SetAuthorization sets the authorization header
func (c *Client) SetAuthorization(token string) {
	c.httpClient.SetAuthorization(token)
}

// SetUserAgent sets the user agent header
func (c *Client) SetUserAgent(userAgent string) {
	c.httpClient.SetUserAgent(userAgent)
}

// Clone creates a copy of the client with the same configuration
func (c *Client) Clone() *Client {
	config := c.config
	return New(
		WithBaseURL(config.BaseURL),
		WithClientTimeout(config.Timeout),
		WithClientHeaders(config.Headers),
		WithUserAgent(config.UserAgent),
	)
}

// WithBaseURL creates a new client with a different base URL
func (c *Client) WithBaseURL(baseURL string) *Client {
	return New(
		WithBaseURL(baseURL),
		WithClientTimeout(c.config.Timeout),
		WithClientHeaders(c.config.Headers),
		WithUserAgent(c.config.UserAgent),
	)
}

// WithTimeout creates a new client with a different timeout
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	return New(
		WithBaseURL(c.config.BaseURL),
		WithClientTimeout(timeout),
		WithClientHeaders(c.config.Headers),
		WithUserAgent(c.config.UserAgent),
	)
}

// WithHeader creates a new client with an additional header
func (c *Client) WithHeader(key, value string) *Client {
	// Copy existing headers
	newHeaders := make(map[string]string)
	if c.config.Headers != nil {
		for k, v := range c.config.Headers {
			newHeaders[k] = v
		}
	}
	newHeaders[key] = value

	return New(
		WithBaseURL(c.config.BaseURL),
		WithClientTimeout(c.config.Timeout),
		WithClientHeaders(newHeaders),
		WithUserAgent(c.config.UserAgent),
	)
}

// WithHeaders creates a new client with additional headers
func (c *Client) WithHeaders(headers map[string]string) *Client {
	// Copy existing headers
	newHeaders := make(map[string]string)
	if c.config.Headers != nil {
		for k, v := range c.config.Headers {
			newHeaders[k] = v
		}
	}

	// Add new headers
	for k, v := range headers {
		newHeaders[k] = v
	}

	return New(
		WithBaseURL(c.config.BaseURL),
		WithClientTimeout(c.config.Timeout),
		WithClientHeaders(newHeaders),
		WithUserAgent(c.config.UserAgent),
	)
}

// buildQueryString builds a query string from a map of parameters
func (c *Client) buildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	var parts []string
	for k, v := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, "&")
}

// Get makes a GET request
func (c *Client) Get(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	req := &Request{
		Method: "GET",
		Path:   path,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	// Build query string if present
	fullPath := path
	if len(req.Query) > 0 {
		queryStr := c.buildQueryString(req.Query)
		if queryStr != "" {
			fullPath = path + "?" + queryStr
		}
	}

	// Set headers on HTTP client temporarily
	originalHeaders := make(map[string]string)
	for k, v := range req.Headers {
		originalHeaders[k] = v
		c.httpClient.SetHeader(k, v)
	}

	httpResp, err := c.httpClient.GetResponseWithContext(ctx, fullPath)

	// Restore original headers
	for k := range req.Headers {
		c.httpClient.SetHeader(k, "")
	}
	for k, v := range originalHeaders {
		c.httpClient.SetHeader(k, v)
	}

	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Headers,
		Body:       httpResp.Body,
		Duration:   0, // Will be set by middleware
	}, nil
}

// Post makes a POST request
func (c *Client) Post(ctx context.Context, path string, data interface{}, opts ...RequestOption) (*Response, error) {
	req := &Request{
		Method: "POST",
		Path:   path,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	// Handle body - either from data parameter or from options
	var body []byte
	var err error
	if data != nil && len(req.Body) == 0 {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
	} else {
		body = req.Body
	}

	// Build query string if present
	fullPath := path
	if len(req.Query) > 0 {
		queryStr := c.buildQueryString(req.Query)
		if queryStr != "" {
			fullPath = path + "?" + queryStr
		}
	}

	// Set headers on HTTP client temporarily
	originalHeaders := make(map[string]string)
	for k, v := range req.Headers {
		originalHeaders[k] = v
		c.httpClient.SetHeader(k, v)
	}

	httpResp, err := c.httpClient.PostResponseWithContext(ctx, fullPath, body)

	// Restore original headers
	for k := range req.Headers {
		c.httpClient.SetHeader(k, "")
	}
	for k, v := range originalHeaders {
		c.httpClient.SetHeader(k, v)
	}

	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Headers,
		Body:       httpResp.Body,
		Duration:   0, // Will be set by middleware
	}, nil
}

// Put makes a PUT request
func (c *Client) Put(ctx context.Context, path string, data interface{}, opts ...RequestOption) (*Response, error) {
	req := &Request{
		Method: "PUT",
		Path:   path,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	// Handle body - either from data parameter or from options
	var body []byte
	var err error
	if data != nil && len(req.Body) == 0 {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
	} else {
		body = req.Body
	}

	// Build query string if present
	fullPath := path
	if len(req.Query) > 0 {
		queryStr := c.buildQueryString(req.Query)
		if queryStr != "" {
			fullPath = path + "?" + queryStr
		}
	}

	// Set headers on HTTP client temporarily
	originalHeaders := make(map[string]string)
	for k, v := range req.Headers {
		originalHeaders[k] = v
		c.httpClient.SetHeader(k, v)
	}

	respBody, err := c.httpClient.PutWithContext(ctx, fullPath, body)

	// Restore original headers
	for k := range req.Headers {
		c.httpClient.SetHeader(k, "")
	}
	for k, v := range originalHeaders {
		c.httpClient.SetHeader(k, v)
	}

	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: 200, // Assume success if no error
		Headers:    make(map[string]string),
		Body:       respBody,
		Duration:   0,
	}, nil
}

// Delete makes a DELETE request
func (c *Client) Delete(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	req := &Request{
		Method: "DELETE",
		Path:   path,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	// Build query string if present
	fullPath := path
	if len(req.Query) > 0 {
		queryStr := c.buildQueryString(req.Query)
		if queryStr != "" {
			fullPath = path + "?" + queryStr
		}
	}

	// Set headers on HTTP client temporarily
	originalHeaders := make(map[string]string)
	for k, v := range req.Headers {
		originalHeaders[k] = v
		c.httpClient.SetHeader(k, v)
	}

	body, err := c.httpClient.DeleteWithContext(ctx, fullPath)

	// Restore original headers
	for k := range req.Headers {
		c.httpClient.SetHeader(k, "")
	}
	for k, v := range originalHeaders {
		c.httpClient.SetHeader(k, v)
	}

	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: 200, // Assume success if no error
		Headers:    make(map[string]string),
		Body:       body,
		Duration:   0,
	}, nil
}

// Patch makes a PATCH request
func (c *Client) Patch(ctx context.Context, path string, data interface{}, opts ...RequestOption) (*Response, error) {
	req := &Request{
		Method: "PATCH",
		Path:   path,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	// Handle body - either from data parameter or from options
	var body []byte
	var err error
	if data != nil && len(req.Body) == 0 {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
	} else {
		body = req.Body
	}

	// Build query string if present
	fullPath := path
	if len(req.Query) > 0 {
		queryStr := c.buildQueryString(req.Query)
		if queryStr != "" {
			fullPath = path + "?" + queryStr
		}
	}

	// Set headers on HTTP client temporarily
	originalHeaders := make(map[string]string)
	for k, v := range req.Headers {
		originalHeaders[k] = v
		c.httpClient.SetHeader(k, v)
	}

	respBody, err := c.httpClient.PatchWithContext(ctx, fullPath, body)

	// Restore original headers
	for k := range req.Headers {
		c.httpClient.SetHeader(k, "")
	}
	for k, v := range originalHeaders {
		c.httpClient.SetHeader(k, v)
	}

	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: 200, // Assume success if no error
		Headers:    make(map[string]string),
		Body:       respBody,
		Duration:   0,
	}, nil
}

// Request makes a generic request
func (c *Client) Request(ctx context.Context, method, path string, data interface{}) (*Response, error) {
	var body []byte
	var err error

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
	}

	respBody, err := c.httpClient.RequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	// Create a simple response for generic requests
	return &Response{
		StatusCode: 200, // Assume success if no error
		Headers:    make(map[string]string),
		Body:       respBody,
		Duration:   0,
	}, nil
}

// SetAuth sets the authorization header
func (c *Client) SetAuth(token string) {
	c.httpClient.SetAuthorization(token)
}

// GetBaseURL returns the base URL
func (c *Client) GetBaseURL() string {
	return c.BaseURL
}

// Response represents an API response
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	Duration   time.Duration
}

// IsSuccess checks if the response indicates success
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError checks if the response indicates a client error
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError checks if the response indicates a server error
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// GetHeader returns a header value
func (r *Response) GetHeader(key string) string {
	return r.Headers[key]
}

// UnmarshalJSON unmarshals the response body into the provided interface
func (r *Response) UnmarshalJSON(v interface{}) error {
	if len(r.Body) == 0 {
		return fmt.Errorf("response body is empty")
	}

	return json.Unmarshal(r.Body, v)
}

// String returns the response body as a string
func (r *Response) String() string {
	return string(r.Body)
}

// Error represents an API error
type Error struct {
	StatusCode int
	Message    string
	Body       string
	Headers    map[string]string
}

func (e *Error) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// NewError creates a new API error
func NewError(statusCode int, message, body string, headers map[string]string) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
		Body:       body,
		Headers:    headers,
	}
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page,omitempty"`
	PerPage    int  `json:"per_page,omitempty"`
	Total      int  `json:"total,omitempty"`
	TotalPages int  `json:"total_pages,omitempty"`
	HasNext    bool `json:"has_next,omitempty"`
	HasPrev    bool `json:"has_prev,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// RequestOptions represents options for API requests
type RequestOptions struct {
	Headers    map[string]string
	Timeout    time.Duration
	Retries    int
	RetryDelay time.Duration
}

// DefaultRequestOptions returns default request options
func DefaultRequestOptions() *RequestOptions {
	return &RequestOptions{
		Headers:    make(map[string]string),
		Timeout:    30 * time.Second,
		Retries:    3,
		RetryDelay: 1 * time.Second,
	}
}

// WithHeader adds a header to the request options
func (o *RequestOptions) WithHeader(key, value string) *RequestOptions {
	o.Headers[key] = value
	return o
}

// WithTimeout sets the timeout for the request
func (o *RequestOptions) WithTimeout(timeout time.Duration) *RequestOptions {
	o.Timeout = timeout
	return o
}

// WithRetries sets the number of retries
func (o *RequestOptions) WithRetries(retries int) *RequestOptions {
	o.Retries = retries
	return o
}

// WithRetryDelay sets the delay between retries
func (o *RequestOptions) WithRetryDelay(delay time.Duration) *RequestOptions {
	o.RetryDelay = delay
	return o
}
