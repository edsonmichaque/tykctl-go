package template

import (
	"context"
	"net/http"
	"time"

	"github.com/edsonmichaque/tykctl-go/retry"
)

// Options represents template loading options.
type Options struct {
	// Template source (one of these should be set)
	Name string // Built-in template name
	URL  string // URL to load template from
	File string // File path to load template from

	// HTTP client configuration for URL loading
	HTTPTimeout time.Duration
	HTTPHeaders map[string]string
	RetryConfig *retry.Config

	// File loading configuration
	FileTimeout time.Duration

	// Validation options
	Insecure bool
}

// Option is a functional option for configuring template loading
type Option func(*Options)

// WithName sets the built-in template name
func WithName(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

// WithURL sets the URL to load template from
func WithURL(url string) Option {
	return func(opts *Options) {
		opts.URL = url
	}
}

// WithFile sets the file path to load template from
func WithFile(file string) Option {
	return func(opts *Options) {
		opts.File = file
	}
}

// WithHTTPTimeout sets the HTTP timeout for URL loading
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.HTTPTimeout = timeout
	}
}

// WithRetryConfig sets the retry configuration
func WithRetryConfig(config *retry.Config) Option {
	return func(opts *Options) {
		opts.RetryConfig = config
	}
}

// WithHTTPHeaders sets custom HTTP headers
func WithHTTPHeaders(headers map[string]string) Option {
	return func(opts *Options) {
		opts.HTTPHeaders = headers
	}
}

// WithFileTimeout sets the file loading timeout
func WithFileTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.FileTimeout = timeout
	}
}

// WithInsecure sets the insecure option (disables SSL validation)
func WithInsecure(insecure bool) Option {
	return func(opts *Options) {
		opts.Insecure = insecure
	}
}

// NewOptions creates a new Options with default values
func NewOptions() *Options {
	return &Options{
		HTTPTimeout: 30 * time.Second,
		HTTPHeaders: make(map[string]string),
		RetryConfig: retry.DefaultConfig(),
		FileTimeout: 10 * time.Second,
		Insecure:    false,
	}
}

// Config represents loader configuration.
type Config struct {
	Client interface{} // *http.Client
}

// Resolve resolves a template from various sources (built-in, file, or URL)
func Resolve(ctx context.Context, opts ...Option) (*Template, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Apply default options and user-provided options
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	// Determine template source and load accordingly
	if options.URL != "" {
		// Load from URL
		config := &Config{
			Client: &http.Client{Timeout: options.HTTPTimeout},
		}
		urlLoader := NewURLLoader(config, options.URL, options)
		return urlLoader.Load(ctx)
	}

	if options.File != "" {
		// Load from file
		fileLoader := NewFileLoader(options.File)
		return fileLoader.Load(ctx)
	}

	if options.Name != "" {
		// Load built-in template
		return resolveBuiltinTemplate(options.Name)
	}

	// No source specified
	return nil, NewTemplateError(ErrorTypeNoSource, "", ErrNoSourceSpecified)
}


