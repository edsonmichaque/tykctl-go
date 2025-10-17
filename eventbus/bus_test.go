package eventbus

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

// Define test event types
const (
	TestEventTypeAPICreate    EventType = "api.create"
	TestEventTypeCommandStart EventType = "command.start"
)

func TestEventBus_Publish(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	bus := New(WithLogger(logger))
	defer bus.Close()

	// Subscribe to events
	receivedEvents := make([]*Event, 0)
	var mu sync.Mutex

	subscription, err := bus.Subscribe(TestEventTypeAPICreate, HandlerFunc(
		func(ctx context.Context, event *Event) error {
			mu.Lock()
			receivedEvents = append(receivedEvents, event)
			mu.Unlock()
			return nil
		},
	))
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	// Publish an event
	event := NewEvent(TestEventTypeAPICreate, map[string]interface{}{
		"api_id": "test-api",
		"name":   "Test API",
	}).WithSource("test")

	err = bus.Publish(event)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Check if event was received
	mu.Lock()
	if len(receivedEvents) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(receivedEvents))
	}

	receivedEvent := receivedEvents[0]
	if receivedEvent.Type != TestEventTypeAPICreate {
		t.Errorf("Expected event type %s, got %s", TestEventTypeAPICreate, receivedEvent.Type)
	}

	if receivedEvent.Source != "test" {
		t.Errorf("Expected source 'test', got %s", receivedEvent.Source)
	}
	mu.Unlock()
}

func TestEventBus_PublishAsync(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	bus := New(WithLogger(logger))
	defer bus.Close()

	// Subscribe to events
	receivedEvents := make([]*Event, 0)
	var mu sync.Mutex

	subscription, err := bus.Subscribe(TestEventTypeAPICreate, HandlerFunc(
		func(ctx context.Context, event *Event) error {
			mu.Lock()
			receivedEvents = append(receivedEvents, event)
			mu.Unlock()
			return nil
		},
	))
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	// Publish events asynchronously
	for i := 0; i < 5; i++ {
		event := NewEvent(TestEventTypeAPICreate, map[string]interface{}{
			"api_id": fmt.Sprintf("test-api-%d", i),
		}).WithSource("test")

		err = bus.PublishAsync(event)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Check if events were received
	mu.Lock()
	if len(receivedEvents) != 5 {
		t.Fatalf("Expected 5 events, got %d", len(receivedEvents))
	}
	mu.Unlock()
}

func TestEventBus_Subscribe(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	bus := New(WithLogger(logger))
	defer bus.Close()

	// Subscribe to events
	subscription, err := bus.Subscribe(TestEventTypeAPICreate, HandlerFunc(
		func(ctx context.Context, event *Event) error {
			return nil
		},
	))
	if err != nil {
		t.Fatal(err)
	}

	// Check subscription details
	if subscription.EventType() != TestEventTypeAPICreate {
		t.Errorf("Expected event type %s, got %s", TestEventTypeAPICreate, subscription.EventType())
	}

	// Unsubscribe
	err = subscription.Unsubscribe()
	if err != nil {
		t.Fatal(err)
	}
}

func TestEventBus_Stats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	bus := New(WithLogger(logger))
	defer bus.Close()

	// Subscribe to events
	subscription, err := bus.Subscribe(TestEventTypeAPICreate, HandlerFunc(
		func(ctx context.Context, event *Event) error {
			return nil
		},
	))
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	// Publish some events
	for i := 0; i < 3; i++ {
		event := NewEvent(TestEventTypeAPICreate, map[string]interface{}{"id": i})
		bus.Publish(event)
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Check stats
	stats := bus.GetStats()
	if stats.EventsPublished != 3 {
		t.Errorf("Expected 3 events published, got %d", stats.EventsPublished)
	}

	if stats.ActiveSubscriptions != 1 {
		t.Errorf("Expected 1 active subscription, got %d", stats.ActiveSubscriptions)
	}
}

func TestEventBus_Middleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	bus := New(WithLogger(logger))
	defer bus.Close()

	// Set up middleware
	loggingMiddleware := NewLoggingMiddleware(logger)
	metricsMiddleware := NewMetricsMiddleware()
	bus.SetMiddleware(loggingMiddleware, metricsMiddleware)

	// Subscribe to events
	subscription, err := bus.Subscribe(TestEventTypeAPICreate, HandlerFunc(
		func(ctx context.Context, event *Event) error {
			return nil
		},
	))
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	// Publish an event
	event := NewEvent(TestEventTypeAPICreate, map[string]interface{}{"api_id": "test"})
	err = bus.Publish(event)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Check that middleware was applied (no errors means it worked)
	stats := bus.GetStats()
	if stats.EventsPublished != 1 {
		t.Errorf("Expected 1 event published, got %d", stats.EventsPublished)
	}
}

func TestEventBus_HandlerPriority(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	bus := New(WithLogger(logger))
	defer bus.Close()

	// Track handler execution order
	executionOrder := make([]string, 0)
	var mu sync.Mutex

	// Create handlers with different priorities
	handler1 := &BaseHandler{
		name:      "handler1",
		priority:  100,
		timeout:   1 * time.Second,
		canHandle: func(EventType) bool { return true },
		handle: func(ctx context.Context, event *Event) error {
			mu.Lock()
			executionOrder = append(executionOrder, "handler1")
			mu.Unlock()
			return nil
		},
	}

	handler2 := &BaseHandler{
		name:      "handler2",
		priority:  50,
		timeout:   1 * time.Second,
		canHandle: func(EventType) bool { return true },
		handle: func(ctx context.Context, event *Event) error {
			mu.Lock()
			executionOrder = append(executionOrder, "handler2")
			mu.Unlock()
			return nil
		},
	}

	// Subscribe handlers
	subscription1, _ := bus.Subscribe(TestEventTypeAPICreate, handler1)
	defer subscription1.Unsubscribe()

	subscription2, _ := bus.Subscribe(TestEventTypeAPICreate, handler2)
	defer subscription2.Unsubscribe()

	// Publish an event
	event := NewEvent(TestEventTypeAPICreate, map[string]interface{}{"api_id": "test"})
	bus.Publish(event)

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Check execution order (higher priority should execute first)
	mu.Lock()
	if len(executionOrder) != 2 {
		t.Fatalf("Expected 2 handlers to execute, got %d", len(executionOrder))
	}

	if executionOrder[0] != "handler1" {
		t.Errorf("Expected handler1 to execute first, got %s", executionOrder[0])
	}

	if executionOrder[1] != "handler2" {
		t.Errorf("Expected handler2 to execute second, got %s", executionOrder[1])
	}
	mu.Unlock()
}

func TestEventBus_ErrorHandling(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	bus := New(WithLogger(logger))
	defer bus.Close()

	// Create a handler that returns an error
	subscription, err := bus.Subscribe(TestEventTypeAPICreate, HandlerFunc(
		func(ctx context.Context, event *Event) error {
			return fmt.Errorf("handler error")
		},
	))
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	// Publish an event
	event := NewEvent(TestEventTypeAPICreate, map[string]interface{}{"api_id": "test"})
	err = bus.Publish(event)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Check stats - should show failed events
	stats := bus.GetStats()
	if stats.EventsFailed == 0 {
		t.Error("Expected events to fail, but none failed")
	}
}

func TestEventBus_Close(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	bus := New(WithLogger(logger))

	// Subscribe to events
	subscription, err := bus.Subscribe(TestEventTypeAPICreate, HandlerFunc(
		func(ctx context.Context, event *Event) error {
			return nil
		},
	))
	if err != nil {
		t.Fatal(err)
	}

	// Close the bus
	err = bus.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Try to unsubscribe after close (should not panic)
	err = subscription.Unsubscribe()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlerRegistry_Walk(t *testing.T) {
	registry := NewHandlerRegistry()

	// Register some handlers
	handler1 := HandlerFunc(func(ctx context.Context, event *Event) error { return nil })
	handler2 := HandlerFunc(func(ctx context.Context, event *Event) error { return nil })
	handler3 := HandlerFunc(func(ctx context.Context, event *Event) error { return nil })

	registry.Register(TestEventTypeAPICreate, handler1)
	registry.Register(TestEventTypeAPICreate, handler2)
	registry.Register(TestEventTypeCommandStart, handler3)

	// Test Walk
	visited := make(map[EventType][]Handler)
	registry.Walk(func(eventType EventType, handlers []Handler) bool {
		visited[eventType] = handlers
		return true
	})

	// Verify all event types were visited
	if len(visited) != 2 {
		t.Fatalf("Expected 2 event types, got %d", len(visited))
	}

	// Verify handlers for each event type
	if len(visited[TestEventTypeAPICreate]) != 2 {
		t.Fatalf("Expected 2 handlers for APICreate, got %d", len(visited[TestEventTypeAPICreate]))
	}

	if len(visited[TestEventTypeCommandStart]) != 1 {
		t.Fatalf("Expected 1 handler for CommandStart, got %d", len(visited[TestEventTypeCommandStart]))
	}

	// Test early termination
	visited = make(map[EventType][]Handler)
	registry.Walk(func(eventType EventType, handlers []Handler) bool {
		visited[eventType] = handlers
		return false // Stop after first iteration
	})

	if len(visited) != 1 {
		t.Fatalf("Expected 1 event type (early termination), got %d", len(visited))
	}
}
