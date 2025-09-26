package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Middleware represents a middleware function
type Middleware func(next func(context.Context, *Request) (*Response, error)) func(context.Context, *Request) (*Response, error)

// Request represents an API request
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Query   map[string]string
	Body    []byte
	Options *RequestOptions
}

// RequestOption is a functional option for configuring requests
type RequestOption func(*Request)

// WithHeader adds a header to the request
func WithHeader(key, value string) RequestOption {
	return func(req *Request) {
		if req.Headers == nil {
			req.Headers = make(map[string]string)
		}
		req.Headers[key] = value
	}
}

// WithHeaders adds multiple headers to the request
func WithHeaders(headers map[string]string) RequestOption {
	return func(req *Request) {
		if req.Headers == nil {
			req.Headers = make(map[string]string)
		}
		for k, v := range headers {
			req.Headers[k] = v
		}
	}
}

// WithQuery adds a query parameter to the request
func WithQuery(key, value string) RequestOption {
	return func(req *Request) {
		if req.Query == nil {
			req.Query = make(map[string]string)
		}
		req.Query[key] = value
	}
}

// WithQueries adds multiple query parameters to the request
func WithQueries(queries map[string]string) RequestOption {
	return func(req *Request) {
		if req.Query == nil {
			req.Query = make(map[string]string)
		}
		for k, v := range queries {
			req.Query[k] = v
		}
	}
}

// WithBody sets the request body
func WithBody(body []byte) RequestOption {
	return func(req *Request) {
		req.Body = body
	}
}

// WithJSONBody sets the request body as JSON
func WithJSONBody(v interface{}) RequestOption {
	return func(req *Request) {
		if v != nil {
			jsonData, err := json.Marshal(v)
			if err != nil {
				// Store error in request for later handling
				req.Headers["_json_error"] = err.Error()
				return
			}
			req.Body = jsonData
			req.Headers["Content-Type"] = "application/json"
		}
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) RequestOption {
	return func(req *Request) {
		if req.Options == nil {
			req.Options = &RequestOptions{}
		}
		req.Options.Timeout = timeout
	}
}

// WithRetries sets the number of retries
func WithRetries(retries int) RequestOption {
	return func(req *Request) {
		if req.Options == nil {
			req.Options = &RequestOptions{}
		}
		req.Options.Retries = retries
	}
}

// WithRetryDelay sets the retry delay
func WithRetryDelay(delay time.Duration) RequestOption {
	return func(req *Request) {
		if req.Options == nil {
			req.Options = &RequestOptions{}
		}
		req.Options.RetryDelay = delay
	}
}

// LoggingMiddleware creates a middleware that logs requests and responses
func LoggingMiddleware(logger interface{}) Middleware {
	return func(next func(context.Context, *Request) (*Response, error)) func(context.Context, *Request) (*Response, error) {
		return func(ctx context.Context, req *Request) (*Response, error) {
			start := time.Now()

			// Log request
			if logger != nil {
				// This would use the logger interface - simplified for now
				fmt.Printf("API Request: %s %s\n", req.Method, req.Path)
			}

			// Execute the request
			resp, err := next(ctx, req)

			// Log response
			if logger != nil {
				duration := time.Since(start)
				if err != nil {
					fmt.Printf("API Error: %s %s - %v (took %v)\n", req.Method, req.Path, err, duration)
				} else {
					fmt.Printf("API Response: %s %s - %d (took %v)\n", req.Method, req.Path, resp.StatusCode, duration)
				}
			}

			return resp, err
		}
	}
}

// TimeoutMiddleware creates a middleware that adds timeout to requests
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next func(context.Context, *Request) (*Response, error)) func(context.Context, *Request) (*Response, error) {
		return func(ctx context.Context, req *Request) (*Response, error) {
			// Use request timeout if specified, otherwise use middleware timeout
			requestTimeout := timeout
			if req.Options != nil && req.Options.Timeout > 0 {
				requestTimeout = req.Options.Timeout
			}

			// Create timeout context
			timeoutCtx, cancel := context.WithTimeout(ctx, requestTimeout)
			defer cancel()

			// Execute with timeout context
			return next(timeoutCtx, req)
		}
	}
}

// RetryMiddleware creates a middleware that adds retry logic to requests
func RetryMiddleware(config RetryConfig) Middleware {
	return func(next func(context.Context, *Request) (*Response, error)) func(context.Context, *Request) (*Response, error) {
		return func(ctx context.Context, req *Request) (*Response, error) {
			// Use request retry config if specified, otherwise use middleware config
			retryConfig := config
			if req.Options != nil && req.Options.Retries > 0 {
				retryConfig.MaxRetries = req.Options.Retries
				retryConfig.InitialDelay = req.Options.RetryDelay
				retryConfig.MaxDelay = req.Options.RetryDelay
			}

			// Execute with retry logic
			return WithRetry(ctx, retryConfig, func() (*Response, error) {
				return next(ctx, req)
			})
		}
	}
}

// AuthMiddleware creates a middleware that adds authentication
func AuthMiddleware(authFunc func(*Request) error) Middleware {
	return func(next func(context.Context, *Request) (*Response, error)) func(context.Context, *Request) (*Response, error) {
		return func(ctx context.Context, req *Request) (*Response, error) {
			// Apply authentication
			if err := authFunc(req); err != nil {
				return nil, fmt.Errorf("authentication failed: %w", err)
			}

			// Execute the request
			return next(ctx, req)
		}
	}
}

// HeaderMiddleware creates a middleware that adds headers to requests
func HeaderMiddleware(headers map[string]string) Middleware {
	return func(next func(context.Context, *Request) (*Response, error)) func(context.Context, *Request) (*Response, error) {
		return func(ctx context.Context, req *Request) (*Response, error) {
			// Add headers to request
			if req.Headers == nil {
				req.Headers = make(map[string]string)
			}

			for k, v := range headers {
				req.Headers[k] = v
			}

			// Execute the request
			return next(ctx, req)
		}
	}
}

// ChainMiddleware chains multiple middlewares together
func ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(next func(context.Context, *Request) (*Response, error)) func(context.Context, *Request) (*Response, error) {
		// Apply middlewares in reverse order (last middleware wraps the next function)
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
