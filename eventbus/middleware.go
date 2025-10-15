// Package eventbus provides middleware implementations for event processing.
package eventbus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Middleware defines the interface for event processing middleware.
type Middleware interface {
	// Process processes an event before it reaches handlers.
	Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error
}

// MiddlewareFunc is a function type that implements Middleware.
type MiddlewareFunc func(ctx context.Context, event *Event, next func(context.Context, *Event) error) error

// Process implements the Middleware interface.
func (f MiddlewareFunc) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	return f(ctx, event, next)
}

// ChainMiddleware chains multiple middleware together.
func ChainMiddleware(middleware ...Middleware) Middleware {
	return MiddlewareFunc(func(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
		// Build the chain in reverse order
		chain := next
		for i := len(middleware) - 1; i >= 0; i-- {
			middleware := middleware[i]
			chain = func(mw Middleware, n func(context.Context, *Event) error) func(context.Context, *Event) error {
				return func(ctx context.Context, event *Event) error {
					return mw.Process(ctx, event, n)
				}
			}(middleware, chain)
		}
		return chain(ctx, event)
	})
}

// LoggingMiddleware logs all events and their processing.
type LoggingMiddleware struct {
	logger *zap.Logger
}

// NewLoggingMiddleware creates a new logging middleware.
func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

// Process logs the event processing.
func (m *LoggingMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	start := time.Now()
	
	m.logger.Info("Processing event",
		zap.String("id", event.ID),
		zap.String("type", string(event.Type)),
		zap.String("source", event.Source),
		zap.Time("timestamp", event.Timestamp))

	err := next(ctx, event)
	duration := time.Since(start)

	if err != nil {
		m.logger.Error("Event processing failed",
			zap.String("id", event.ID),
			zap.String("type", string(event.Type)),
			zap.Duration("duration", duration),
			zap.Error(err))
	} else {
		m.logger.Debug("Event processed successfully",
			zap.String("id", event.ID),
			zap.String("type", string(event.Type)),
			zap.Duration("duration", duration))
	}

	return err
}

// MetricsMiddleware collects metrics for events.
type MetricsMiddleware struct {
	metrics map[string]interface{}
	mu      sync.RWMutex
}

// NewMetricsMiddleware creates a new metrics middleware.
func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: make(map[string]interface{}),
	}
}

// Process collects metrics for the event.
func (m *MetricsMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	start := time.Now()
	err := next(ctx, event)
	duration := time.Since(start)

	m.mu.Lock()
	defer m.mu.Unlock()

	// Update event type counters
	eventTypeKey := fmt.Sprintf("events.%s.count", event.Type)
	if count, exists := m.metrics[eventTypeKey]; exists {
		m.metrics[eventTypeKey] = count.(int64) + 1
	} else {
		m.metrics[eventTypeKey] = int64(1)
	}

	// Update duration metrics
	durationKey := fmt.Sprintf("events.%s.duration", event.Type)
	if durations, exists := m.metrics[durationKey]; exists {
		durations.([]time.Duration) = append(durations.([]time.Duration), duration)
		m.metrics[durationKey] = durations
	} else {
		m.metrics[durationKey] = []time.Duration{duration}
	}

	// Update error counters
	if err != nil {
		errorKey := fmt.Sprintf("events.%s.errors", event.Type)
		if count, exists := m.metrics[errorKey]; exists {
			m.metrics[errorKey] = count.(int64) + 1
		} else {
			m.metrics[errorKey] = int64(1)
		}
	}

	return err
}

// GetMetrics returns the collected metrics.
func (m *MetricsMiddleware) GetMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{})
	for k, v := range m.metrics {
		result[k] = v
	}
	return result
}

// ValidationMiddleware validates events before processing.
type ValidationMiddleware struct {
	validators map[EventType]func(*Event) error
	mu         sync.RWMutex
}

// NewValidationMiddleware creates a new validation middleware.
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validators: make(map[EventType]func(*Event) error),
	}
}

// AddValidator adds a validator for an event type.
func (m *ValidationMiddleware) AddValidator(eventType EventType, validator func(*Event) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.validators[eventType] = validator
}

// Process validates the event before processing.
func (m *ValidationMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	m.mu.RLock()
	validator, exists := m.validators[event.Type]
	m.mu.RUnlock()

	if exists && validator != nil {
		if err := validator(event); err != nil {
			return fmt.Errorf("event validation failed: %w", err)
		}
	}

	return next(ctx, event)
}

// RateLimitMiddleware limits the rate of event processing.
type RateLimitMiddleware struct {
	limiters map[EventType]*RateLimiter
	mu       sync.RWMutex
}

// NewRateLimitMiddleware creates a new rate limit middleware.
func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiters: make(map[EventType]*RateLimiter),
	}
}

// SetRateLimit sets the rate limit for an event type.
func (m *RateLimitMiddleware) SetRateLimit(eventType EventType, rate int, per time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[eventType] = NewRateLimiter(rate, per)
}

// Process applies rate limiting to the event.
func (m *RateLimitMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	m.mu.RLock()
	limiter, exists := m.limiters[event.Type]
	m.mu.RUnlock()

	if exists && limiter != nil {
		if !limiter.Allow() {
			return fmt.Errorf("rate limit exceeded for event type %s", event.Type)
		}
	}

	return next(ctx, event)
}

// CircuitBreakerMiddleware implements circuit breaker pattern for event processing.
type CircuitBreakerMiddleware struct {
	breakers map[EventType]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerMiddleware creates a new circuit breaker middleware.
func NewCircuitBreakerMiddleware() *CircuitBreakerMiddleware {
	return &CircuitBreakerMiddleware{
		breakers: make(map[EventType]*CircuitBreaker),
	}
}

// SetCircuitBreaker sets the circuit breaker for an event type.
func (m *CircuitBreakerMiddleware) SetCircuitBreaker(eventType EventType, failureThreshold int, timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.breakers[eventType] = NewCircuitBreaker(failureThreshold, timeout)
}

// Process applies circuit breaker logic to the event.
func (m *CircuitBreakerMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	m.mu.RLock()
	breaker, exists := m.breakers[event.Type]
	m.mu.RUnlock()

	if exists && breaker != nil {
		if !breaker.Allow() {
			return fmt.Errorf("circuit breaker open for event type %s", event.Type)
		}

		err := next(ctx, event)
		breaker.RecordResult(err == nil)
		return err
	}

	return next(ctx, event)
}

// RetryMiddleware retries failed event processing.
type RetryMiddleware struct {
	maxRetries int
	retryDelay time.Duration
	backoffFunc func(int) time.Duration
}

// NewRetryMiddleware creates a new retry middleware.
func NewRetryMiddleware(maxRetries int, retryDelay time.Duration) *RetryMiddleware {
	return &RetryMiddleware{
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		backoffFunc: func(attempt int) time.Duration {
			return retryDelay * time.Duration(attempt)
		},
	}
}

// Process retries failed event processing.
func (m *RetryMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	var lastErr error

	for attempt := 0; attempt <= m.maxRetries; attempt++ {
		if attempt > 0 {
			delay := m.backoffFunc(attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := next(ctx, event)
		if err == nil {
			return nil
		}

		lastErr = err
	}

	return fmt.Errorf("event processing failed after %d retries: %w", m.maxRetries, lastErr)
}

// TimeoutMiddleware adds timeout to event processing.
type TimeoutMiddleware struct {
	timeout time.Duration
}

// NewTimeoutMiddleware creates a new timeout middleware.
func NewTimeoutMiddleware(timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{timeout: timeout}
}

// Process adds timeout to event processing.
func (m *TimeoutMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- next(ctx, event)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("event processing timeout after %v: %w", m.timeout, ctx.Err())
	}
}

// CorrelationMiddleware adds correlation tracking to events.
type CorrelationMiddleware struct {
	correlationID string
}

// NewCorrelationMiddleware creates a new correlation middleware.
func NewCorrelationMiddleware(correlationID string) *CorrelationMiddleware {
	return &CorrelationMiddleware{correlationID: correlationID}
}

// Process adds correlation ID to events.
func (m *CorrelationMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	if event.CorrelationID == "" {
		event.CorrelationID = m.correlationID
	}
	return next(ctx, event)
}

// SanitizationMiddleware sanitizes event data.
type SanitizationMiddleware struct {
	sanitizers map[string]func(interface{}) interface{}
	mu         sync.RWMutex
}

// NewSanitizationMiddleware creates a new sanitization middleware.
func NewSanitizationMiddleware() *SanitizationMiddleware {
	return &SanitizationMiddleware{
		sanitizers: make(map[string]func(interface{}) interface{}),
	}
}

// AddSanitizer adds a sanitizer for a specific field.
func (m *SanitizationMiddleware) AddSanitizer(field string, sanitizer func(interface{}) interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sanitizers[field] = sanitizer
}

// Process sanitizes event data.
func (m *SanitizationMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	m.mu.RLock()
	sanitizers := make(map[string]func(interface{}) interface{})
	for k, v := range m.sanitizers {
		sanitizers[k] = v
	}
	m.mu.RUnlock()

	// Sanitize metadata
	if event.Metadata != nil {
		for key, value := range event.Metadata {
			if sanitizer, exists := sanitizers[key]; exists {
				event.Metadata[key] = sanitizer(value)
			}
		}
	}

	return next(ctx, event)
}

// DefaultSanitizers returns common sanitizers for sensitive data.
func DefaultSanitizers() map[string]func(interface{}) interface{} {
	return map[string]func(interface{}) interface{}{
		"password":     func(v interface{}) interface{} { return "[REDACTED]" },
		"token":        func(v interface{}) interface{} { return "[REDACTED]" },
		"secret":       func(v interface{}) interface{} { return "[REDACTED]" },
		"api_key":      func(v interface{}) interface{} { return "[REDACTED]" },
		"access_token": func(v interface{}) interface{} { return "[REDACTED]" },
		"refresh_token": func(v interface{}) interface{} { return "[REDACTED]" },
	}
}