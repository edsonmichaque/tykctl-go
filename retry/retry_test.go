package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRetry_Success(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	
	err := Retry(ctx, nil, func() error {
		callCount++
		if callCount == 1 {
			return errors.New("temporary error")
		}
		return nil
	})
	
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestRetry_MaxRetries(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxRetries:     2,
		InitialDelay:   10 * time.Millisecond,
		MaxDelay:       50 * time.Millisecond,
		BackoffFactor:  2.0,
		MaxElapsedTime: 1 * time.Second,
	}
	
	callCount := 0
	
	err := Retry(ctx, config, func() error {
		callCount++
		return errors.New("permanent error")
	})
	
	if err == nil {
		t.Error("Expected error, got success")
	}
	
	if callCount != 3 { // Initial call + 2 retries
		t.Errorf("Expected 3 calls, got %d", callCount)
	}
}

func TestRetryWithResult_Success(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	
	result, err := RetryWithResult(ctx, nil, func() (string, error) {
		callCount++
		if callCount == 1 {
			return "", errors.New("temporary error")
		}
		return "success", nil
	})
	
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	
	if result != "success" {
		t.Errorf("Expected 'success', got '%s'", result)
	}
	
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestRetryWithTimeout(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxRetries:     10,
		InitialDelay:   100 * time.Millisecond,
		MaxDelay:       1 * time.Second,
		BackoffFactor:  2.0,
		MaxElapsedTime: 10 * time.Second,
	}
	
	start := time.Now()
	err := RetryWithTimeout(ctx, config, 500*time.Millisecond, func() error {
		return errors.New("permanent error")
	})
	
	elapsed := time.Since(start)
	
	if err == nil {
		t.Error("Expected error, got success")
	}
	
	if elapsed > 1*time.Second {
		t.Errorf("Expected timeout around 500ms, got %v", elapsed)
	}
}

func TestRetryWithTimeoutAndResult(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxRetries:     10,
		InitialDelay:   100 * time.Millisecond,
		MaxDelay:       1 * time.Second,
		BackoffFactor:  2.0,
		MaxElapsedTime: 10 * time.Second,
	}
	
	start := time.Now()
	result, err := RetryWithTimeoutAndResult(ctx, config, 500*time.Millisecond, func() (int, error) {
		return 0, errors.New("permanent error")
	})
	
	elapsed := time.Since(start)
	
	if err == nil {
		t.Error("Expected error, got success")
	}
	
	if result != 0 {
		t.Errorf("Expected result 0, got %d", result)
	}
	
	if elapsed > 1*time.Second {
		t.Errorf("Expected timeout around 500ms, got %v", elapsed)
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: false,
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: false,
		},
		{
			name:     "authentication failed",
			err:      errors.New("authentication failed"),
			expected: false,
		},
		{
			name:     "authorization denied",
			err:      errors.New("authorization denied"),
			expected: false,
		},
		{
			name:     "invalid credentials",
			err:      errors.New("invalid credentials"),
			expected: false,
		},
		{
			name:     "network timeout",
			err:      errors.New("network timeout"),
			expected: true,
		},
		{
			name:     "server error",
			err:      errors.New("server error"),
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryableError(tt.err)
			if result != tt.expected {
				t.Errorf("IsRetryableError(%v) = %v, expected %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestRetryIfRetryable_NonRetryableError(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	
	err := RetryIfRetryable(ctx, nil, func() error {
		callCount++
		return errors.New("authentication failed")
	})
	
	if err == nil {
		t.Error("Expected error, got success")
	}
	
	if callCount != 1 {
		t.Errorf("Expected 1 call for non-retryable error, got %d", callCount)
	}
}

func TestRetryIfRetryable_RetryableError(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxRetries:     2,
		InitialDelay:   10 * time.Millisecond,
		MaxDelay:       50 * time.Millisecond,
		BackoffFactor:  2.0,
		MaxElapsedTime: 1 * time.Second,
	}
	
	callCount := 0
	
	err := RetryIfRetryable(ctx, config, func() error {
		callCount++
		return errors.New("network timeout")
	})
	
	if err == nil {
		t.Error("Expected error, got success")
	}
	
	if callCount != 3 { // Initial call + 2 retries
		t.Errorf("Expected 3 calls for retryable error, got %d", callCount)
	}
}

func TestRetryIfRetryableWithResult_NonRetryableError(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	
	result, err := RetryIfRetryableWithResult(ctx, nil, func() (string, error) {
		callCount++
		return "", errors.New("authentication failed")
	})
	
	if err == nil {
		t.Error("Expected error, got success")
	}
	
	if result != "" {
		t.Errorf("Expected empty result, got '%s'", result)
	}
	
	if callCount != 1 {
		t.Errorf("Expected 1 call for non-retryable error, got %d", callCount)
	}
}

func TestRetryIfRetryableWithResult_RetryableError(t *testing.T) {
	ctx := context.Background()
	config := &Config{
		MaxRetries:     2,
		InitialDelay:   10 * time.Millisecond,
		MaxDelay:       50 * time.Millisecond,
		BackoffFactor:  2.0,
		MaxElapsedTime: 1 * time.Second,
	}
	
	callCount := 0
	
	result, err := RetryIfRetryableWithResult(ctx, config, func() (string, error) {
		callCount++
		return "", errors.New("network timeout")
	})
	
	if err == nil {
		t.Error("Expected error, got success")
	}
	
	if result != "" {
		t.Errorf("Expected empty result, got '%s'", result)
	}
	
	if callCount != 3 { // Initial call + 2 retries
		t.Errorf("Expected 3 calls for retryable error, got %d", callCount)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries 5, got %d", config.MaxRetries)
	}
	
	if config.InitialDelay != time.Second {
		t.Errorf("Expected InitialDelay 1s, got %v", config.InitialDelay)
	}
	
	if config.MaxDelay != 30*time.Second {
		t.Errorf("Expected MaxDelay 30s, got %v", config.MaxDelay)
	}
	
	if config.BackoffFactor != 2.0 {
		t.Errorf("Expected BackoffFactor 2.0, got %f", config.BackoffFactor)
	}
	
	if config.MaxElapsedTime != 5*time.Minute {
		t.Errorf("Expected MaxElapsedTime 5m, got %v", config.MaxElapsedTime)
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	callCount := 0
	
	// Cancel context after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	
	err := Retry(ctx, nil, func() error {
		callCount++
		return errors.New("temporary error")
	})
	
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
	
	if callCount != 1 {
		t.Errorf("Expected 1 call before cancellation, got %d", callCount)
	}
}