// Package eventbus provides configuration for the event bus.
package eventbus

import (
	"time"

	"go.uber.org/zap"
)

// Config contains configuration for the event bus.
type Config struct {
	// AsyncWorkers is the number of async workers.
	AsyncWorkers int

	// AsyncQueueSize is the size of the async queue.
	AsyncQueueSize int

	// Logger is the logger instance.
	Logger *zap.Logger

	// DefaultTimeout is the default timeout for event processing.
	DefaultTimeout time.Duration

	// MaxRetries is the maximum number of retries for failed events.
	MaxRetries int

	// RetryDelay is the delay between retries.
	RetryDelay time.Duration

	// EnableMetrics enables metrics collection.
	EnableMetrics bool

	// EnableLogging enables event logging.
	EnableLogging bool

	// EnableValidation enables event validation.
	EnableValidation bool

	// EnableRateLimit enables rate limiting.
	EnableRateLimit bool

	// EnableCircuitBreaker enables circuit breaker.
	EnableCircuitBreaker bool

	// EnableSanitization enables data sanitization.
	EnableSanitization bool

	// Middleware contains middleware configuration.
	Middleware MiddlewareConfig
}

// MiddlewareConfig contains middleware-specific configuration.
type MiddlewareConfig struct {
	// Logging configures logging middleware.
	Logging LoggingConfig

	// Metrics configures metrics middleware.
	Metrics MetricsConfig

	// Validation configures validation middleware.
	Validation ValidationConfig

	// RateLimit configures rate limiting middleware.
	RateLimit RateLimitConfig

	// CircuitBreaker configures circuit breaker middleware.
	CircuitBreaker CircuitBreakerConfig

	// Retry configures retry middleware.
	Retry RetryConfig

	// Timeout configures timeout middleware.
	Timeout TimeoutConfig

	// Sanitization configures sanitization middleware.
	Sanitization SanitizationConfig
}

// LoggingConfig contains logging middleware configuration.
type LoggingConfig struct {
	// Enabled enables logging middleware.
	Enabled bool

	// Level is the log level for events.
	Level zap.Level

	// IncludeData includes event data in logs.
	IncludeData bool

	// IncludeMetadata includes event metadata in logs.
	IncludeMetadata bool
}

// MetricsConfig contains metrics middleware configuration.
type MetricsConfig struct {
	// Enabled enables metrics middleware.
	Enabled bool

	// CollectDurations collects duration metrics.
	CollectDurations bool

	// CollectCounts collects count metrics.
	CollectCounts bool

	// CollectErrors collects error metrics.
	CollectErrors bool

	// FlushInterval is the interval for flushing metrics.
	FlushInterval time.Duration
}

// ValidationConfig contains validation middleware configuration.
type ValidationConfig struct {
	// Enabled enables validation middleware.
	Enabled bool

	// StrictMode enables strict validation.
	StrictMode bool

	// RequiredFields are fields that must be present.
	RequiredFields map[EventType][]string

	// CustomValidators are custom validation functions.
	CustomValidators map[EventType]func(*Event) error
}

// RateLimitConfig contains rate limiting middleware configuration.
type RateLimitConfig struct {
	// Enabled enables rate limiting middleware.
	Enabled bool

	// DefaultRate is the default rate limit.
	DefaultRate int

	// DefaultPer is the default time period.
	DefaultPer time.Duration

	// PerEventType allows different rates per event type.
	PerEventType map[EventType]RateLimit
}

// RateLimit represents a rate limit configuration.
type RateLimit struct {
	Rate int
	Per  time.Duration
}

// CircuitBreakerConfig contains circuit breaker middleware configuration.
type CircuitBreakerConfig struct {
	// Enabled enables circuit breaker middleware.
	Enabled bool

	// DefaultFailureThreshold is the default failure threshold.
	DefaultFailureThreshold int

	// DefaultTimeout is the default timeout.
	DefaultTimeout time.Duration

	// PerEventType allows different settings per event type.
	PerEventType map[EventType]CircuitBreakerSettings
}

// CircuitBreakerSettings represents circuit breaker settings.
type CircuitBreakerSettings struct {
	FailureThreshold int
	Timeout          time.Duration
}

// RetryConfig contains retry middleware configuration.
type RetryConfig struct {
	// Enabled enables retry middleware.
	Enabled bool

	// MaxRetries is the maximum number of retries.
	MaxRetries int

	// RetryDelay is the delay between retries.
	RetryDelay time.Duration

	// BackoffMultiplier is the backoff multiplier.
	BackoffMultiplier float64

	// MaxRetryDelay is the maximum retry delay.
	MaxRetryDelay time.Duration
}

// TimeoutConfig contains timeout middleware configuration.
type TimeoutConfig struct {
	// Enabled enables timeout middleware.
	Enabled bool

	// DefaultTimeout is the default timeout.
	DefaultTimeout time.Duration

	// PerEventType allows different timeouts per event type.
	PerEventType map[EventType]time.Duration
}

// SanitizationConfig contains sanitization middleware configuration.
type SanitizationConfig struct {
	// Enabled enables sanitization middleware.
	Enabled bool

	// SensitiveFields are fields that should be sanitized.
	SensitiveFields []string

	// SanitizationFunction is the function to use for sanitization.
	SanitizationFunction func(string) string

	// CustomSanitizers are custom sanitizers for specific fields.
	CustomSanitizers map[string]func(interface{}) interface{}
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		AsyncWorkers:   10,
		AsyncQueueSize: 1000,
		Logger:         zap.NewNop(),
		DefaultTimeout: 30 * time.Second,
		MaxRetries:     3,
		RetryDelay:     1 * time.Second,
		EnableMetrics:  true,
		EnableLogging:  true,
		EnableValidation: false,
		EnableRateLimit: false,
		EnableCircuitBreaker: false,
		EnableSanitization: true,
		Middleware: MiddlewareConfig{
			Logging: LoggingConfig{
				Enabled:         true,
				Level:           zap.InfoLevel,
				IncludeData:     false,
				IncludeMetadata: true,
			},
			Metrics: MetricsConfig{
				Enabled:         true,
				CollectDurations: true,
				CollectCounts:   true,
				CollectErrors:   true,
				FlushInterval:   1 * time.Minute,
			},
			Validation: ValidationConfig{
				Enabled:         false,
				StrictMode:      false,
				RequiredFields:  make(map[EventType][]string),
				CustomValidators: make(map[EventType]func(*Event) error),
			},
			RateLimit: RateLimitConfig{
				Enabled:     false,
				DefaultRate: 100,
				DefaultPer:  1 * time.Minute,
				PerEventType: make(map[EventType]RateLimit),
			},
			CircuitBreaker: CircuitBreakerConfig{
				Enabled:                false,
				DefaultFailureThreshold: 5,
				DefaultTimeout:         30 * time.Second,
				PerEventType:           make(map[EventType]CircuitBreakerSettings),
			},
			Retry: RetryConfig{
				Enabled:           false,
				MaxRetries:        3,
				RetryDelay:        1 * time.Second,
				BackoffMultiplier: 2.0,
				MaxRetryDelay:     30 * time.Second,
			},
			Timeout: TimeoutConfig{
				Enabled:        false,
				DefaultTimeout: 30 * time.Second,
				PerEventType:   make(map[EventType]time.Duration),
			},
			Sanitization: SanitizationConfig{
				Enabled:              true,
				SensitiveFields:      []string{"password", "token", "secret", "api_key"},
				SanitizationFunction: func(s string) string { return "[REDACTED]" },
				CustomSanitizers:     DefaultSanitizers(),
			},
		},
	}
}

// Option is a function that configures the event bus.
type Option func(*Config)

// WithAsyncWorkers sets the number of async workers.
func WithAsyncWorkers(count int) Option {
	return func(c *Config) {
		c.AsyncWorkers = count
	}
}

// WithAsyncQueueSize sets the async queue size.
func WithAsyncQueueSize(size int) Option {
	return func(c *Config) {
		c.AsyncQueueSize = size
	}
}

// WithLogger sets the logger.
func WithLogger(logger *zap.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

// WithDefaultTimeout sets the default timeout.
func WithDefaultTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.DefaultTimeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(maxRetries int) Option {
	return func(c *Config) {
		c.MaxRetries = maxRetries
	}
}

// WithRetryDelay sets the retry delay.
func WithRetryDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.RetryDelay = delay
	}
}

// WithMetrics enables or disables metrics.
func WithMetrics(enabled bool) Option {
	return func(c *Config) {
		c.EnableMetrics = enabled
	}
}

// WithLogging enables or disables logging.
func WithLogging(enabled bool) Option {
	return func(c *Config) {
		c.EnableLogging = enabled
	}
}

// WithValidation enables or disables validation.
func WithValidation(enabled bool) Option {
	return func(c *Config) {
		c.EnableValidation = enabled
	}
}

// WithRateLimit enables or disables rate limiting.
func WithRateLimit(enabled bool) Option {
	return func(c *Config) {
		c.EnableRateLimit = enabled
	}
}

// WithCircuitBreaker enables or disables circuit breaker.
func WithCircuitBreaker(enabled bool) Option {
	return func(c *Config) {
		c.EnableCircuitBreaker = enabled
	}
}

// WithSanitization enables or disables sanitization.
func WithSanitization(enabled bool) Option {
	return func(c *Config) {
		c.EnableSanitization = enabled
	}
}

// WithLoggingConfig sets the logging configuration.
func WithLoggingConfig(config LoggingConfig) Option {
	return func(c *Config) {
		c.Middleware.Logging = config
	}
}

// WithMetricsConfig sets the metrics configuration.
func WithMetricsConfig(config MetricsConfig) Option {
	return func(c *Config) {
		c.Middleware.Metrics = config
	}
}

// WithValidationConfig sets the validation configuration.
func WithValidationConfig(config ValidationConfig) Option {
	return func(c *Config) {
		c.Middleware.Validation = config
	}
}

// WithRateLimitConfig sets the rate limit configuration.
func WithRateLimitConfig(config RateLimitConfig) Option {
	return func(c *Config) {
		c.Middleware.RateLimit = config
	}
}

// WithCircuitBreakerConfig sets the circuit breaker configuration.
func WithCircuitBreakerConfig(config CircuitBreakerConfig) Option {
	return func(c *Config) {
		c.Middleware.CircuitBreaker = config
	}
}

// WithRetryConfig sets the retry configuration.
func WithRetryConfig(config RetryConfig) Option {
	return func(c *Config) {
		c.Middleware.Retry = config
	}
}

// WithTimeoutConfig sets the timeout configuration.
func WithTimeoutConfig(config TimeoutConfig) Option {
	return func(c *Config) {
		c.Middleware.Timeout = config
	}
}

// WithSanitizationConfig sets the sanitization configuration.
func WithSanitizationConfig(config SanitizationConfig) Option {
	return func(c *Config) {
		c.Middleware.Sanitization = config
	}
}