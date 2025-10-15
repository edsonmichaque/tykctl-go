// Package telemetry provides anonymous usage analytics for tykctl-go.
package telemetry

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// client implements the Client interface.
type client struct {
	config    *Config
	transport Transport
	storage   Storage
	enabled   bool
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// NewClient creates a new telemetry client.
func NewClient(config *Config, transport Transport, storage Storage) Client {
	ctx, cancel := context.WithCancel(context.Background())
	
	c := &client{
		config:    config,
		transport: transport,
		storage:   storage,
		enabled:   config.Enabled,
		ctx:       ctx,
		cancel:    cancel,
	}
	
	// Start the background worker for batching and sending events
	c.wg.Add(1)
	go c.worker()
	
	return c
}

// Track records a telemetry event.
func (c *client) Track(event *Event) error {
	c.mu.RLock()
	enabled := c.enabled
	c.mu.RUnlock()
	
	if !enabled {
		return nil
	}
	
	// Sanitize the event to remove sensitive data
	sanitized := SanitizeEvent(event)
	
	// Store the event for later transmission
	return c.storage.Store(sanitized)
}

// Flush sends any pending events immediately.
func (c *client) Flush() error {
	c.mu.RLock()
	enabled := c.enabled
	c.mu.RUnlock()
	
	if !enabled {
		return nil
	}
	
	// Retrieve all stored events
	events, err := c.storage.Retrieve()
	if err != nil {
		return fmt.Errorf("failed to retrieve events: %w", err)
	}
	
	if len(events) == 0 {
		return nil
	}
	
	// Send events in batches
	return c.sendBatches(events)
}

// Close shuts down the telemetry client gracefully.
func (c *client) Close() error {
	// Cancel the context to stop the worker
	c.cancel()
	
	// Wait for the worker to finish
	c.wg.Wait()
	
	// Flush any remaining events
	if err := c.Flush(); err != nil {
		return fmt.Errorf("failed to flush events on close: %w", err)
	}
	
	// Close the transport
	if err := c.transport.Close(); err != nil {
		return fmt.Errorf("failed to close transport: %w", err)
	}
	
	return nil
}

// IsEnabled returns whether telemetry is currently enabled.
func (c *client) IsEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.enabled
}

// SetEnabled enables or disables telemetry.
func (c *client) SetEnabled(enabled bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.enabled = enabled
	c.config.Enabled = enabled
	
	// If disabling, flush any pending events
	if !enabled {
		go func() {
			if err := c.Flush(); err != nil {
				// Log error but don't fail the operation
				fmt.Printf("Warning: failed to flush telemetry events: %v\n", err)
			}
		}()
	}
	
	return nil
}

// worker runs in the background to batch and send events periodically.
func (c *client) worker() {
	defer c.wg.Done()
	
	ticker := time.NewTicker(c.config.FlushInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.Flush(); err != nil {
				// Log error but continue running
				fmt.Printf("Warning: failed to flush telemetry events: %v\n", err)
			}
		}
	}
}

// sendBatches sends events in batches according to the configured batch size.
func (c *client) sendBatches(events []*Event) error {
	batchSize := c.config.BatchSize
	if batchSize <= 0 {
		batchSize = 100 // Default batch size
	}
	
	for i := 0; i < len(events); i += batchSize {
		end := i + batchSize
		if end > len(events) {
			end = len(events)
		}
		
		batch := events[i:end]
		
		// Send the batch with retry logic
		if err := c.sendWithRetry(batch); err != nil {
			return fmt.Errorf("failed to send batch %d-%d: %w", i, end-1, err)
		}
	}
	
	// Clear the storage after successful send
	return c.storage.Clear()
}

// sendWithRetry sends a batch with retry logic.
func (c *client) sendWithRetry(batch []*Event) error {
	var lastErr error
	
	for attempt := 0; attempt <= c.config.RetryAttempts; attempt++ {
		if attempt > 0 {
			// Wait before retrying
			select {
			case <-c.ctx.Done():
				return c.ctx.Err()
			case <-time.After(c.config.RetryDelay * time.Duration(attempt)):
			}
		}
		
		if err := c.transport.Send(batch); err != nil {
			lastErr = err
			continue
		}
		
		// Success
		return nil
	}
	
	return fmt.Errorf("failed to send batch after %d attempts: %w", c.config.RetryAttempts+1, lastErr)
}

// NoOpClient is a no-operation client that discards all telemetry events.
type NoOpClient struct{}

// NewNoOpClient creates a new no-operation telemetry client.
func NewNoOpClient() Client {
	return &NoOpClient{}
}

// Track discards the event.
func (c *NoOpClient) Track(event *Event) error {
	return nil
}

// Flush does nothing.
func (c *NoOpClient) Flush() error {
	return nil
}

// Close does nothing.
func (c *NoOpClient) Close() error {
	return nil
}

// IsEnabled returns false.
func (c *NoOpClient) IsEnabled() bool {
	return false
}

// SetEnabled does nothing.
func (c *NoOpClient) SetEnabled(enabled bool) error {
	return nil
}