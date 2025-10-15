// Package eventbus provides rate limiting and circuit breaker implementations.
package eventbus

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter.
type RateLimiter struct {
	rate     int
	per      time.Duration
	tokens   int
	lastTime time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		per:      per,
		tokens:   rate,
		lastTime: time.Now(),
	}
}

// Allow checks if a request is allowed under the rate limit.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTime)
	
	// Add tokens based on elapsed time
	tokensToAdd := int(elapsed.Nanoseconds() / rl.per.Nanoseconds())
	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.rate {
			rl.tokens = rl.rate
		}
		rl.lastTime = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// CircuitBreaker implements a circuit breaker pattern.
type CircuitBreaker struct {
	failureThreshold int
	timeout          time.Duration
	failureCount     int
	lastFailureTime  time.Time
	state            CircuitState
	mu               sync.Mutex
}

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	// CircuitClosed means the circuit is closed and requests are allowed.
	CircuitClosed CircuitState = iota
	// CircuitOpen means the circuit is open and requests are blocked.
	CircuitOpen
	// CircuitHalfOpen means the circuit is half-open and testing requests.
	CircuitHalfOpen
)

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(failureThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		timeout:          timeout,
		state:            CircuitClosed,
	}
}

// Allow checks if a request is allowed through the circuit breaker.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if now.Sub(cb.lastFailureTime) >= cb.timeout {
			cb.state = CircuitHalfOpen
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

// RecordResult records the result of a request.
func (cb *CircuitBreaker) RecordResult(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	if success {
		cb.failureCount = 0
		if cb.state == CircuitHalfOpen {
			cb.state = CircuitClosed
		}
	} else {
		cb.failureCount++
		cb.lastFailureTime = now

		if cb.state == CircuitClosed && cb.failureCount >= cb.failureThreshold {
			cb.state = CircuitOpen
		} else if cb.state == CircuitHalfOpen {
			cb.state = CircuitOpen
		}
	}
}

// GetState returns the current state of the circuit breaker.
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// GetFailureCount returns the current failure count.
func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.failureCount
}

// Reset resets the circuit breaker to its initial state.
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.failureCount = 0
	cb.state = CircuitClosed
	cb.lastFailureTime = time.Time{}
}