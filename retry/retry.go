package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Config represents retry configuration
type Config struct {
	MaxRetries    int           `json:"max_retries"`
	InitialDelay  time.Duration `json:"initial_delay"`
	MaxDelay      time.Duration `json:"max_delay"`
	BackoffFactor float64       `json:"backoff_factor"`
	MaxElapsedTime time.Duration `json:"max_elapsed_time"`
}

// DefaultConfig returns default retry configuration
func DefaultConfig() *Config {
	return &Config{
		MaxRetries:     5,
		InitialDelay:   time.Second,
		MaxDelay:       30 * time.Second,
		BackoffFactor:  2.0,
		MaxElapsedTime: 5 * time.Minute,
	}
}

// Operation represents a function that can be retried
type Operation func() error

// OperationWithResult represents a function that returns a result and can be retried
type OperationWithResult[T any] func() (T, error)

// Retry executes an operation with retry logic using exponential backoff
func Retry(ctx context.Context, config *Config, operation Operation) error {
	if config == nil {
		config = DefaultConfig()
	}

	// Create backoff configuration
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = config.InitialDelay
	backoffConfig.MaxInterval = config.MaxDelay
	backoffConfig.Multiplier = config.BackoffFactor
	backoffConfig.MaxElapsedTime = config.MaxElapsedTime

	// Create backoff with context
	backoffWithContext := backoff.WithContext(backoffConfig, ctx)

	// Execute operation with backoff
	err := backoff.Retry(backoff.Operation(operation), backoffWithContext)
	if err != nil {
		return fmt.Errorf("operation failed after retries: %w", err)
	}

	return nil
}

// RetryWithResult executes an operation that returns a result with retry logic
func RetryWithResult[T any](ctx context.Context, config *Config, operation OperationWithResult[T]) (T, error) {
	var result T
	var lastErr error

	if config == nil {
		config = DefaultConfig()
	}

	// Create backoff configuration
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = config.InitialDelay
	backoffConfig.MaxInterval = config.MaxDelay
	backoffConfig.Multiplier = config.BackoffFactor
	backoffConfig.MaxElapsedTime = config.MaxElapsedTime

	// Create backoff with context
	backoffWithContext := backoff.WithContext(backoffConfig, ctx)

	// Execute operation with backoff
	err := backoff.Retry(func() error {
		var err error
		result, err = operation()
		lastErr = err
		return err
	}, backoffWithContext)

	if err != nil {
		return result, fmt.Errorf("operation failed after retries: %w", lastErr)
	}

	return result, nil
}

// RetryWithTimeout executes an operation with a timeout context
func RetryWithTimeout(ctx context.Context, config *Config, timeout time.Duration, operation Operation) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return Retry(ctx, config, operation)
}

// RetryWithTimeoutAndResult executes an operation with a timeout context and returns a result
func RetryWithTimeoutAndResult[T any](ctx context.Context, config *Config, timeout time.Duration, operation OperationWithResult[T]) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return RetryWithResult(ctx, config, operation)
}

// IsRetryableError checks if an error should trigger a retry
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for context cancellation
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}

	// Check for specific error types that should not be retried
	// Add more specific error checks as needed
	switch err.Error() {
	case "authentication failed", "authorization denied", "invalid credentials":
		return false
	}

	return true
}

// RetryIfRetryable executes an operation with retry logic only if the error is retryable
func RetryIfRetryable(ctx context.Context, config *Config, operation Operation) error {
	return Retry(ctx, config, func() error {
		err := operation()
		if err != nil && !IsRetryableError(err) {
			// Return a non-retryable error wrapped to stop retries
			return backoff.Permanent(err)
		}
		return err
	})
}

// RetryIfRetryableWithResult executes an operation with retry logic only if the error is retryable
func RetryIfRetryableWithResult[T any](ctx context.Context, config *Config, operation OperationWithResult[T]) (T, error) {
	var result T
	var lastErr error

	if config == nil {
		config = DefaultConfig()
	}

	// Create backoff configuration
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = config.InitialDelay
	backoffConfig.MaxInterval = config.MaxDelay
	backoffConfig.Multiplier = config.BackoffFactor
	backoffConfig.MaxElapsedTime = config.MaxElapsedTime

	// Create backoff with context
	backoffWithContext := backoff.WithContext(backoffConfig, ctx)

	// Execute operation with backoff
	err := backoff.Retry(func() error {
		var err error
		result, err = operation()
		lastErr = err
		if err != nil && !IsRetryableError(err) {
			// Return a non-retryable error wrapped to stop retries
			return backoff.Permanent(err)
		}
		return err
	}, backoffWithContext)

	if err != nil {
		return result, fmt.Errorf("operation failed after retries: %w", lastErr)
	}

	return result, nil
}