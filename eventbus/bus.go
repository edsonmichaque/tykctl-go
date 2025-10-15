// Package eventbus provides the core event bus implementation with both sync and async processing.
package eventbus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// EventBus defines the interface for the event bus.
type EventBus interface {
	// Publish publishes an event synchronously.
	Publish(event *Event) error

	// PublishAsync publishes an event asynchronously.
	PublishAsync(event *Event) error

	// Subscribe subscribes to events of a specific type.
	Subscribe(eventType EventType, handler Handler) (Subscription, error)

	// Unsubscribe removes a subscription.
	Unsubscribe(subscription Subscription) error

	// Close shuts down the event bus.
	Close() error

	// GetStats returns event bus statistics.
	GetStats() *Stats

	// SetMiddleware sets global middleware.
	SetMiddleware(middleware ...Middleware)
}

// Subscription represents an event subscription.
type Subscription interface {
	// ID returns the subscription ID.
	ID() string

	// EventType returns the subscribed event type.
	EventType() EventType

	// Handler returns the event handler.
	Handler() Handler

	// Unsubscribe removes the subscription.
	Unsubscribe() error
}

// Middleware defines middleware for event processing.
type Middleware interface {
	// Process processes an event before it reaches handlers.
	Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error
}

// Stats contains event bus statistics.
type Stats struct {
	EventsPublished    int64     `json:"events_published"`
	EventsProcessed    int64     `json:"events_processed"`
	EventsFailed       int64     `json:"events_failed"`
	ActiveSubscriptions int64    `json:"active_subscriptions"`
	StartTime          time.Time `json:"start_time"`
	LastEventTime      time.Time `json:"last_event_time"`
}

// eventBus implements the EventBus interface.
type eventBus struct {
	registry     *HandlerRegistry
	middleware   []Middleware
	stats        *Stats
	mu           sync.RWMutex
	logger       *zap.Logger
	asyncWorkers int
	asyncQueue   chan *Event
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// New creates a new event bus.
func New(options ...Option) EventBus {
	config := &Config{
		AsyncWorkers: 10,
		AsyncQueueSize: 1000,
		Logger: zap.NewNop(),
	}

	for _, option := range options {
		option(config)
	}

	eb := &eventBus{
		registry:     NewHandlerRegistry(),
		middleware:   make([]Middleware, 0),
		stats:        &Stats{StartTime: time.Now()},
		logger:       config.Logger,
		asyncWorkers: config.AsyncWorkers,
		asyncQueue:   make(chan *Event, config.AsyncQueueSize),
		stopChan:     make(chan struct{}),
	}

	// Start async workers
	for i := 0; i < eb.asyncWorkers; i++ {
		eb.wg.Add(1)
		go eb.asyncWorker(i)
	}

	return eb
}

// Publish publishes an event synchronously.
func (eb *eventBus) Publish(event *Event) error {
	eb.mu.Lock()
	eb.stats.EventsPublished++
	eb.stats.LastEventTime = time.Now()
	eb.mu.Unlock()

	handlers := eb.registry.GetHandlers(event.Type)
	if len(handlers) == 0 {
		eb.logger.Debug("No handlers for event type", zap.String("type", string(event.Type)))
		return nil
	}

	ctx := context.Background()
	return eb.processEvent(ctx, event, handlers)
}

// PublishAsync publishes an event asynchronously.
func (eb *eventBus) PublishAsync(event *Event) error {
	eb.mu.Lock()
	eb.stats.EventsPublished++
	eb.stats.LastEventTime = time.Now()
	eb.mu.Unlock()

	select {
	case eb.asyncQueue <- event:
		return nil
	default:
		return fmt.Errorf("async queue is full")
	}
}

// Subscribe subscribes to events of a specific type.
func (eb *eventBus) Subscribe(eventType EventType, handler Handler) (Subscription, error) {
	eb.registry.Register(eventType, handler)

	subscription := &subscription{
		id:        fmt.Sprintf("%s-%d", string(eventType), time.Now().UnixNano()),
		eventType: eventType,
		handler:   handler,
		bus:       eb,
	}

	eb.mu.Lock()
	eb.stats.ActiveSubscriptions++
	eb.mu.Unlock()

	eb.logger.Info("Subscribed to event type", 
		zap.String("type", string(eventType)),
		zap.String("handler", handler.GetName()))

	return subscription, nil
}

// Unsubscribe removes a subscription.
func (eb *eventBus) Unsubscribe(subscription Subscription) error {
	eb.registry.Unregister(subscription.EventType(), subscription.ID())

	eb.mu.Lock()
	eb.stats.ActiveSubscriptions--
	eb.mu.Unlock()

	eb.logger.Info("Unsubscribed from event type",
		zap.String("type", string(subscription.EventType())),
		zap.String("handler", subscription.Handler().GetName()))

	return nil
}

// Close shuts down the event bus.
func (eb *eventBus) Close() error {
	close(eb.stopChan)
	close(eb.asyncQueue)
	eb.wg.Wait()
	return nil
}

// GetStats returns event bus statistics.
func (eb *eventBus) GetStats() *Stats {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	stats := *eb.stats
	return &stats
}

// SetMiddleware sets global middleware.
func (eb *eventBus) SetMiddleware(middleware ...Middleware) {
	eb.middleware = middleware
}

// processEvent processes an event through handlers.
func (eb *eventBus) processEvent(ctx context.Context, event *Event, handlers []Handler) error {
	// Apply middleware
	next := func(ctx context.Context, event *Event) error {
		return eb.executeHandlers(ctx, event, handlers)
	}

	for i := len(eb.middleware) - 1; i >= 0; i-- {
		middleware := eb.middleware[i]
		next = func(mw Middleware, n func(context.Context, *Event) error) func(context.Context, *Event) error {
			return func(ctx context.Context, event *Event) error {
				return mw.Process(ctx, event, n)
			}
		}(middleware, next)
	}

	return next(ctx, event)
}

// executeHandlers executes all handlers for an event.
func (eb *eventBus) executeHandlers(ctx context.Context, event *Event, handlers []Handler) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	for _, handler := range handlers {
		if !handler.CanHandle(event.Type) {
			continue
		}

		wg.Add(1)
		go func(h Handler) {
			defer wg.Done()

			handlerCtx, cancel := context.WithTimeout(ctx, h.GetTimeout())
			defer cancel()

			err := h.Handle(handlerCtx, event)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("handler %s: %w", h.GetName(), err))
				mu.Unlock()

				eb.mu.Lock()
				eb.stats.EventsFailed++
				eb.mu.Unlock()

				eb.logger.Error("Handler failed",
					zap.String("handler", h.GetName()),
					zap.String("event_type", string(event.Type)),
					zap.Error(err))
			} else {
				eb.mu.Lock()
				eb.stats.EventsProcessed++
				eb.mu.Unlock()
			}
		}(handler)
	}

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("handler errors: %v", errors)
	}

	return nil
}

// asyncWorker processes events from the async queue.
func (eb *eventBus) asyncWorker(workerID int) {
	defer eb.wg.Done()

	eb.logger.Debug("Started async worker", zap.Int("worker_id", workerID))

	for {
		select {
		case event := <-eb.asyncQueue:
			handlers := eb.registry.GetHandlers(event.Type)
			if len(handlers) == 0 {
				eb.logger.Debug("No handlers for async event type", 
					zap.String("type", string(event.Type)),
					zap.Int("worker_id", workerID))
				continue
			}

			ctx := context.Background()
			err := eb.processEvent(ctx, event, handlers)
			if err != nil {
				eb.logger.Error("Async event processing failed",
					zap.String("type", string(event.Type)),
					zap.Int("worker_id", workerID),
					zap.Error(err))
			}
		case <-eb.stopChan:
			eb.logger.Debug("Stopping async worker", zap.Int("worker_id", workerID))
			return
		}
	}
}

// subscription implements the Subscription interface.
type subscription struct {
	id        string
	eventType EventType
	handler   Handler
	bus       *eventBus
}

// ID returns the subscription ID.
func (s *subscription) ID() string {
	return s.id
}

// EventType returns the subscribed event type.
func (s *subscription) EventType() EventType {
	return s.eventType
}

// Handler returns the event handler.
func (s *subscription) Handler() Handler {
	return s.handler
}

// Unsubscribe removes the subscription.
func (s *subscription) Unsubscribe() error {
	return s.bus.Unsubscribe(s)
}

// Config contains event bus configuration.
type Config struct {
	AsyncWorkers   int
	AsyncQueueSize int
	Logger         *zap.Logger
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

// Built-in middleware implementations

// LoggingMiddleware logs all events.
type LoggingMiddleware struct {
	logger *zap.Logger
}

// NewLoggingMiddleware creates a new logging middleware.
func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

// Process logs the event.
func (m *LoggingMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	m.logger.Info("Processing event",
		zap.String("id", event.ID),
		zap.String("type", string(event.Type)),
		zap.String("source", event.Source))

	start := time.Now()
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
	stats *Stats
}

// NewMetricsMiddleware creates a new metrics middleware.
func NewMetricsMiddleware(stats *Stats) *MetricsMiddleware {
	return &MetricsMiddleware{stats: stats}
}

// Process collects metrics for the event.
func (m *MetricsMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	start := time.Now()
	err := next(ctx, event)
	duration := time.Since(start)

	// Update metrics based on result
	if err != nil {
		// Metrics would be updated here
		_ = duration
	}

	return err
}

// ValidationMiddleware validates events before processing.
type ValidationMiddleware struct {
	validators map[EventType]func(*Event) error
}

// NewValidationMiddleware creates a new validation middleware.
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validators: make(map[EventType]func(*Event) error),
	}
}

// AddValidator adds a validator for an event type.
func (m *ValidationMiddleware) AddValidator(eventType EventType, validator func(*Event) error) {
	m.validators[eventType] = validator
}

// Process validates the event before processing.
func (m *ValidationMiddleware) Process(ctx context.Context, event *Event, next func(context.Context, *Event) error) error {
	if validator, exists := m.validators[event.Type]; exists {
		if err := validator(event); err != nil {
			return fmt.Errorf("event validation failed: %w", err)
		}
	}

	return next(ctx, event)
}