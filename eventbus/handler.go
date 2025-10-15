// Package eventbus provides event handler interfaces and implementations.
package eventbus

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Handler defines the interface for event handlers.
type Handler interface {
	// Handle processes an event.
	Handle(ctx context.Context, event *Event) error

	// CanHandle returns true if this handler can process the given event type.
	CanHandle(eventType EventType) bool

	// GetName returns the handler name.
	GetName() string

	// GetPriority returns the handler priority (higher numbers execute first).
	GetPriority() int

	// GetTimeout returns the handler timeout.
	GetTimeout() time.Duration
}

// HandlerFunc is a function type that implements Handler.
type HandlerFunc func(ctx context.Context, event *Event) error

// Handle implements the Handler interface.
func (f HandlerFunc) Handle(ctx context.Context, event *Event) error {
	return f(ctx, event)
}

// CanHandle always returns true for function handlers.
func (f HandlerFunc) CanHandle(eventType EventType) bool {
	return true
}

// GetName returns a default name for function handlers.
func (f HandlerFunc) GetName() string {
	return "function-handler"
}

// GetPriority returns the default priority.
func (f HandlerFunc) GetPriority() int {
	return 0
}

// GetTimeout returns the default timeout.
func (f HandlerFunc) GetTimeout() time.Duration {
	return 30 * time.Second
}

// BaseHandler provides a base implementation for event handlers.
type BaseHandler struct {
	name     string
	priority int
	timeout  time.Duration
	canHandle func(EventType) bool
	handle   func(context.Context, *Event) error
}

// NewBaseHandler creates a new base handler.
func NewBaseHandler(name string, priority int, timeout time.Duration) *BaseHandler {
	return &BaseHandler{
		name:     name,
		priority: priority,
		timeout:  timeout,
		canHandle: func(EventType) bool { return true },
		handle:   func(context.Context, *Event) error { return nil },
	}
}

// Handle processes an event.
func (h *BaseHandler) Handle(ctx context.Context, event *Event) error {
	return h.handle(ctx, event)
}

// CanHandle checks if the handler can process the event type.
func (h *BaseHandler) CanHandle(eventType EventType) bool {
	return h.canHandle(eventType)
}

// GetName returns the handler name.
func (h *BaseHandler) GetName() string {
	return h.name
}

// GetPriority returns the handler priority.
func (h *BaseHandler) GetPriority() int {
	return h.priority
}

// GetTimeout returns the handler timeout.
func (h *BaseHandler) GetTimeout() time.Duration {
	return h.timeout
}

// SetCanHandle sets the can handle function.
func (h *BaseHandler) SetCanHandle(fn func(EventType) bool) *BaseHandler {
	h.canHandle = fn
	return h
}

// SetHandle sets the handle function.
func (h *BaseHandler) SetHandle(fn func(context.Context, *Event) error) *BaseHandler {
	h.handle = fn
	return h
}

// FilteredHandler wraps another handler with filtering capabilities.
type FilteredHandler struct {
	handler Handler
	filter  EventFilter
}

// NewFilteredHandler creates a new filtered handler.
func NewFilteredHandler(handler Handler, filter EventFilter) *FilteredHandler {
	return &FilteredHandler{
		handler: handler,
		filter:  filter,
	}
}

// Handle processes an event if it matches the filter.
func (h *FilteredHandler) Handle(ctx context.Context, event *Event) error {
	if !h.filter.Matches(event) {
		return nil // Skip processing
	}
	return h.handler.Handle(ctx, event)
}

// CanHandle checks if the handler can process the event type.
func (h *FilteredHandler) CanHandle(eventType EventType) bool {
	return h.handler.CanHandle(eventType)
}

// GetName returns the handler name.
func (h *FilteredHandler) GetName() string {
	return h.handler.GetName()
}

// GetPriority returns the handler priority.
func (h *FilteredHandler) GetPriority() int {
	return h.handler.GetPriority()
}

// GetTimeout returns the handler timeout.
func (h *FilteredHandler) GetTimeout() time.Duration {
	return h.handler.GetTimeout()
}

// RetryHandler wraps another handler with retry capabilities.
type RetryHandler struct {
	handler     Handler
	maxRetries  int
	retryDelay  time.Duration
	backoffFunc func(int) time.Duration
}

// NewRetryHandler creates a new retry handler.
func NewRetryHandler(handler Handler, maxRetries int, retryDelay time.Duration) *RetryHandler {
	return &RetryHandler{
		handler:    handler,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		backoffFunc: func(attempt int) time.Duration {
			return retryDelay * time.Duration(attempt)
		},
	}
}

// Handle processes an event with retry logic.
func (h *RetryHandler) Handle(ctx context.Context, event *Event) error {
	var lastErr error
	
	for attempt := 0; attempt <= h.maxRetries; attempt++ {
		if attempt > 0 {
			delay := h.backoffFunc(attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := h.handler.Handle(ctx, event)
		if err == nil {
			return nil
		}
		
		lastErr = err
	}

	return fmt.Errorf("handler failed after %d retries: %w", h.maxRetries, lastErr)
}

// CanHandle checks if the handler can process the event type.
func (h *RetryHandler) CanHandle(eventType EventType) bool {
	return h.handler.CanHandle(eventType)
}

// GetName returns the handler name.
func (h *RetryHandler) GetName() string {
	return h.handler.GetName()
}

// GetPriority returns the handler priority.
func (h *RetryHandler) GetPriority() int {
	return h.handler.GetPriority()
}

// GetTimeout returns the handler timeout.
func (h *RetryHandler) GetTimeout() time.Duration {
	return h.handler.GetTimeout()
}

// TimeoutHandler wraps another handler with timeout capabilities.
type TimeoutHandler struct {
	handler Handler
	timeout time.Duration
}

// NewTimeoutHandler creates a new timeout handler.
func NewTimeoutHandler(handler Handler, timeout time.Duration) *TimeoutHandler {
	return &TimeoutHandler{
		handler: handler,
		timeout: timeout,
	}
}

// Handle processes an event with timeout.
func (h *TimeoutHandler) Handle(ctx context.Context, event *Event) error {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- h.handler.Handle(ctx, event)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("handler timeout after %v: %w", h.timeout, ctx.Err())
	}
}

// CanHandle checks if the handler can process the event type.
func (h *TimeoutHandler) CanHandle(eventType EventType) bool {
	return h.handler.CanHandle(eventType)
}

// GetName returns the handler name.
func (h *TimeoutHandler) GetName() string {
	return h.handler.GetName()
}

// GetPriority returns the handler priority.
func (h *TimeoutHandler) GetPriority() int {
	return h.handler.GetPriority()
}

// GetTimeout returns the handler timeout.
func (h *TimeoutHandler) GetTimeout() time.Duration {
	return h.timeout
}

// BatchHandler processes multiple events in batches.
type BatchHandler struct {
	handler    Handler
	batchSize  int
	flushDelay time.Duration
	events     []*Event
	mu         sync.Mutex
	flushChan  chan struct{}
	stopChan   chan struct{}
}

// NewBatchHandler creates a new batch handler.
func NewBatchHandler(handler Handler, batchSize int, flushDelay time.Duration) *BatchHandler {
	bh := &BatchHandler{
		handler:    handler,
		batchSize:  batchSize,
		flushDelay: flushDelay,
		events:     make([]*Event, 0, batchSize),
		flushChan:  make(chan struct{}, 1),
		stopChan:   make(chan struct{}),
	}

	go bh.flushLoop()
	return bh
}

// Handle adds an event to the batch.
func (h *BatchHandler) Handle(ctx context.Context, event *Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = append(h.events, event)

	if len(h.events) >= h.batchSize {
		select {
		case h.flushChan <- struct{}{}:
		default:
		}
	}

	return nil
}

// CanHandle checks if the handler can process the event type.
func (h *BatchHandler) CanHandle(eventType EventType) bool {
	return h.handler.CanHandle(eventType)
}

// GetName returns the handler name.
func (h *BatchHandler) GetName() string {
	return h.handler.GetName()
}

// GetPriority returns the handler priority.
func (h *BatchHandler) GetPriority() int {
	return h.handler.GetPriority()
}

// GetTimeout returns the handler timeout.
func (h *BatchHandler) GetTimeout() time.Duration {
	return h.handler.GetTimeout()
}

// flushLoop periodically flushes the batch.
func (h *BatchHandler) flushLoop() {
	ticker := time.NewTicker(h.flushDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.flush()
		case <-h.flushChan:
			h.flush()
		case <-h.stopChan:
			h.flush() // Flush remaining events
			return
		}
	}
}

// flush processes the current batch of events.
func (h *BatchHandler) flush() {
	h.mu.Lock()
	if len(h.events) == 0 {
		h.mu.Unlock()
		return
	}

	events := make([]*Event, len(h.events))
	copy(events, h.events)
	h.events = h.events[:0]
	h.mu.Unlock()

	// Process each event in the batch
	for _, event := range events {
		ctx := context.Background()
		_ = h.handler.Handle(ctx, event)
	}
}

// Close stops the batch handler.
func (h *BatchHandler) Close() error {
	close(h.stopChan)
	return nil
}

// HandlerRegistry manages event handlers.
type HandlerRegistry struct {
	handlers map[EventType][]Handler
	mu       sync.RWMutex
}

// NewHandlerRegistry creates a new handler registry.
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[EventType][]Handler),
	}
}

// Register registers a handler for an event type.
func (r *HandlerRegistry) Register(eventType EventType, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	handlers := r.handlers[eventType]
	handlers = append(handlers, handler)
	
	// Sort by priority (higher priority first)
	for i := len(handlers) - 1; i > 0; i-- {
		if handlers[i].GetPriority() > handlers[i-1].GetPriority() {
			handlers[i], handlers[i-1] = handlers[i-1], handlers[i]
		}
	}

	r.handlers[eventType] = handlers
}

// Unregister removes a handler for an event type.
func (r *HandlerRegistry) Unregister(eventType EventType, handlerName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	handlers := r.handlers[eventType]
	for i, h := range handlers {
		if h.GetName() == handlerName {
			r.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// GetHandlers returns handlers for an event type.
func (r *HandlerRegistry) GetHandlers(eventType EventType) []Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handlers := r.handlers[eventType]
	result := make([]Handler, len(handlers))
	copy(result, handlers)
	return result
}

// Walk iterates over all registered handlers, calling fn for each event type and its handlers.
// If fn returns false, the iteration stops.
func (r *HandlerRegistry) Walk(fn func(eventType EventType, handlers []Handler) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for eventType, handlers := range r.handlers {
		// Create a copy of handlers to avoid holding the lock
		handlersCopy := make([]Handler, len(handlers))
		copy(handlersCopy, handlers)
		
		if !fn(eventType, handlersCopy) {
			break
		}
	}
}

// GetAllHandlers returns all registered handlers.
// Deprecated: Use Walk instead for better performance.
func (r *HandlerRegistry) GetAllHandlers() map[EventType][]Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[EventType][]Handler)
	for eventType, handlers := range r.handlers {
		result[eventType] = make([]Handler, len(handlers))
		copy(result[eventType], handlers)
	}
	return result
}

// Clear removes all handlers.
func (r *HandlerRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers = make(map[EventType][]Handler)
}