package api

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// RetryConfig represents retry configuration
type RetryConfig struct {
	MaxRetries     int
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	MaxElapsedTime time.Duration
	Multiplier     float64
	Retryable      RetryCondition
}

// RetryCondition defines when to retry a request
type RetryCondition interface {
	ShouldRetry(err error, response *Response) bool
}

// DefaultRetryCondition implements default retry logic
type DefaultRetryCondition struct{}

func (c *DefaultRetryCondition) ShouldRetry(err error, response *Response) bool {
	if err != nil {
		return true // Retry on network errors
	}

	if response == nil {
		return false
	}

	// Retry on server errors (5xx) and rate limiting (429)
	return response.IsServerError() || response.StatusCode == 429
}

// RetryableError represents an error that can be retried
type RetryableError struct {
	Err        error
	Attempts   int
	MaxRetries int
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable error (attempts %d/%d): %v", e.Attempts, e.MaxRetries, e.Err)
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// WithRetry executes a function with retry logic using external backoff library
func WithRetry(ctx context.Context, config RetryConfig, fn func() (*Response, error)) (*Response, error) {
	var lastResponse *Response
	var attempts int

	// Create exponential backoff
	expBackoff := backoff.NewExponentialBackOff()

	// Configure backoff parameters
	if config.InitialDelay > 0 {
		expBackoff.InitialInterval = config.InitialDelay
	}
	if config.MaxDelay > 0 {
		expBackoff.MaxInterval = config.MaxDelay
	}
	if config.MaxElapsedTime > 0 {
		expBackoff.MaxElapsedTime = config.MaxElapsedTime
	}
	if config.Multiplier > 0 {
		expBackoff.Multiplier = config.Multiplier
	}

	// Create retry operation
	operation := func() error {
		attempts++

		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return backoff.Permanent(ctx.Err())
		default:
		}

		// Execute the function
		response, err := fn()
		lastResponse = response

		// Check if we should retry
		shouldRetry := false
		if config.Retryable != nil {
			shouldRetry = config.Retryable.ShouldRetry(err, response)
		} else {
			// Default retry condition
			condition := &DefaultRetryCondition{}
			shouldRetry = condition.ShouldRetry(err, response)
		}

		// If no retry needed, return success
		if !shouldRetry {
			return nil
		}

		// Check max retries
		if config.MaxRetries > 0 && attempts >= config.MaxRetries {
			return backoff.Permanent(&RetryableError{
				Err:        err,
				Attempts:   attempts,
				MaxRetries: config.MaxRetries,
			})
		}

		// Return error to trigger retry
		return err
	}

	// Execute with backoff
	err := backoff.Retry(operation, expBackoff)
	if err != nil {
		return lastResponse, err
	}

	return lastResponse, nil
}

// WithRetryNotify executes a function with retry logic and notification callback
func WithRetryNotify(ctx context.Context, config RetryConfig, fn func() (*Response, error), notify func(error, time.Duration)) (*Response, error) {
	var lastResponse *Response
	var attempts int

	// Create exponential backoff
	expBackoff := backoff.NewExponentialBackOff()

	// Configure backoff parameters
	if config.InitialDelay > 0 {
		expBackoff.InitialInterval = config.InitialDelay
	}
	if config.MaxDelay > 0 {
		expBackoff.MaxInterval = config.MaxDelay
	}
	if config.MaxElapsedTime > 0 {
		expBackoff.MaxElapsedTime = config.MaxElapsedTime
	}
	if config.Multiplier > 0 {
		expBackoff.Multiplier = config.Multiplier
	}

	// Create retry operation
	operation := func() error {
		attempts++

		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return backoff.Permanent(ctx.Err())
		default:
		}

		// Execute the function
		response, err := fn()
		lastResponse = response

		// Check if we should retry
		shouldRetry := false
		if config.Retryable != nil {
			shouldRetry = config.Retryable.ShouldRetry(err, response)
		} else {
			// Default retry condition
			condition := &DefaultRetryCondition{}
			shouldRetry = condition.ShouldRetry(err, response)
		}

		// If no retry needed, return success
		if !shouldRetry {
			return nil
		}

		// Check max retries
		if config.MaxRetries > 0 && attempts >= config.MaxRetries {
			return backoff.Permanent(&RetryableError{
				Err:        err,
				Attempts:   attempts,
				MaxRetries: config.MaxRetries,
			})
		}

		// Return error to trigger retry
		return err
	}

	// Create notification function
	notifyFunc := func(err error, duration time.Duration) {
		if notify != nil {
			notify(err, duration)
		}
	}

	// Execute with backoff and notification
	err := backoff.RetryNotify(operation, expBackoff, notifyFunc)
	if err != nil {
		return lastResponse, err
	}

	return lastResponse, nil
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialDelay:   1 * time.Second,
		MaxDelay:       30 * time.Second,
		MaxElapsedTime: 2 * time.Minute,
		Multiplier:     2.0,
		Retryable:      &DefaultRetryCondition{},
	}
}

// NewExponentialBackoffConfig creates a retry config with exponential backoff
func NewExponentialBackoffConfig(maxRetries int, initialDelay, maxDelay, maxElapsedTime time.Duration) RetryConfig {
	return RetryConfig{
		MaxRetries:     maxRetries,
		InitialDelay:   initialDelay,
		MaxDelay:       maxDelay,
		MaxElapsedTime: maxElapsedTime,
		Multiplier:     2.0,
		Retryable:      &DefaultRetryCondition{},
	}
}

// NewConstantBackoffConfig creates a retry config with constant backoff
func NewConstantBackoffConfig(maxRetries int, delay time.Duration) RetryConfig {
	return RetryConfig{
		MaxRetries:     maxRetries,
		InitialDelay:   delay,
		MaxDelay:       delay,
		MaxElapsedTime: time.Duration(maxRetries) * delay,
		Multiplier:     1.0,
		Retryable:      &DefaultRetryCondition{},
	}
}
