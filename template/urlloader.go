package template

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/edsonmichaque/tykctl-go/retry"
	"gopkg.in/yaml.v3"
)

// URLLoader loads templates from URLs.
type URLLoader struct {
	url     string
	options *Options
	client  *http.Client
}

// NewURLLoader creates a new URL loader.
func NewURLLoader(config *Config, url string, options *Options) *URLLoader {
	return &URLLoader{
		url:     url,
		options: options,
		client:  config.Client.(*http.Client),
	}
}

// Load loads a template from a URL.
func (u *URLLoader) Load(ctx context.Context) (*Template, error) {
	// Use retry package for robust retry logic
	return retry.RetryWithResult(ctx, u.options.RetryConfig, func() (*Template, error) {
		// Create request with timeout
		reqCtx := ctx
		if u.options.HTTPTimeout > 0 {
			var cancel context.CancelFunc
			reqCtx, cancel = context.WithTimeout(ctx, u.options.HTTPTimeout)
			defer cancel()
		}

		req, err := http.NewRequestWithContext(reqCtx, "GET", u.url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add custom headers
		for key, value := range u.options.HTTPHeaders {
			req.Header.Set(key, value)
		}

		// Make request
		resp, err := u.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch template: %w", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("template not found: %s", resp.Status)
		}

		// Read response
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		// Parse template
		var template Template
		if err := yaml.Unmarshal(data, &template); err != nil {
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}

		return &template, nil
	})
}
