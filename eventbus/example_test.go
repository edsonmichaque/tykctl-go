package eventbus_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/edsonmichaque/tykctl-go/eventbus"
	"go.uber.org/zap"
)

// Define your own event types
const (
	EventTypeAPICreate eventbus.EventType = "api.create"
	EventTypeAPIUpdate eventbus.EventType = "api.update"
	EventTypeAPIDelete eventbus.EventType = "api.delete"
	EventTypeCommandStart eventbus.EventType = "command.start"
	EventTypeCommandComplete eventbus.EventType = "command.complete"
)

func ExampleEventBus_basicUsage() {
	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create event bus
	bus := eventbus.New(
		eventbus.WithLogger(logger),
		eventbus.WithAsyncWorkers(2),
		eventbus.WithAsyncQueueSize(100),
	)
	defer bus.Close()

	// Subscribe to API create events
	subscription, err := bus.Subscribe(EventTypeAPICreate, eventbus.HandlerFunc(
		func(ctx context.Context, event *eventbus.Event) error {
			fmt.Printf("API created: %+v\n", event.Data)
			return nil
		},
	))
	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Unsubscribe()

	// Publish an API create event
	event := eventbus.NewEvent(EventTypeAPICreate, map[string]interface{}{
		"api_id": "api-123",
		"name":   "My API",
	}).WithSource("tykctl-gateway")

	// Synchronous publishing
	err = bus.Publish(event)
	if err != nil {
		log.Fatal(err)
	}

	// Asynchronous publishing
	err = bus.PublishAsync(event)
	if err != nil {
		log.Fatal(err)
	}

	// Wait a bit for async processing
	time.Sleep(100 * time.Millisecond)

	// Output:
	// API created: map[api_id:api-123 name:My API]
	// API created: map[api_id:api-123 name:My API]
}

func ExampleEventBus_withMiddleware() {
	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create event bus with middleware
	bus := eventbus.New(
		eventbus.WithLogger(logger),
		eventbus.WithAsyncWorkers(1),
	)

	// Set up middleware
	loggingMiddleware := eventbus.NewLoggingMiddleware(logger)
	metricsMiddleware := eventbus.NewMetricsMiddleware()
	validationMiddleware := eventbus.NewValidationMiddleware()

	// Add validation for API create events
	validationMiddleware.AddValidator(EventTypeAPICreate, func(event *eventbus.Event) error {
		if event.Data == nil {
			return fmt.Errorf("event data is required")
		}
		return nil
	})

	bus.SetMiddleware(loggingMiddleware, metricsMiddleware, validationMiddleware)

	// Subscribe to events
	subscription, _ := bus.Subscribe(EventTypeAPICreate, eventbus.HandlerFunc(
		func(ctx context.Context, event *eventbus.Event) error {
			fmt.Printf("Processing API: %+v\n", event.Data)
			return nil
		},
	))
	defer subscription.Unsubscribe()

	// Publish valid event
	validEvent := eventbus.NewEvent(EventTypeAPICreate, map[string]interface{}{
		"api_id": "api-456",
		"name":   "Valid API",
	})
	bus.Publish(validEvent)

	// Publish invalid event (will fail validation)
	invalidEvent := eventbus.NewEvent(EventTypeAPICreate, nil)
	bus.Publish(invalidEvent)

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Output:
	// Processing API: map[api_id:api-456 name:Valid API]
}

func ExampleEventBus_withRetry() {
	logger, _ := zap.NewDevelopment()
	bus := eventbus.New(eventbus.WithLogger(logger))

	// Create a handler that fails the first two times
	attemptCount := 0
	handler := eventbus.HandlerFunc(func(ctx context.Context, event *eventbus.Event) error {
		attemptCount++
		fmt.Printf("Attempt %d for event %s\n", attemptCount, event.Type)
		
		if attemptCount < 3 {
			return fmt.Errorf("simulated failure")
		}
		
		fmt.Printf("Successfully processed event %s\n", event.Type)
		return nil
	})

	// Wrap with retry middleware
	retryHandler := eventbus.NewRetryHandler(handler, 3, 100*time.Millisecond)

	// Subscribe
	subscription, _ := bus.Subscribe(EventTypeAPICreate, retryHandler)
	defer subscription.Unsubscribe()

	// Publish event
	event := eventbus.NewEvent(EventTypeAPICreate, map[string]interface{}{
		"api_id": "api-retry-test",
	})
	bus.Publish(event)

	// Wait for processing
	time.Sleep(1 * time.Second)

	// Output:
	// Attempt 1 for event api.create
	// Attempt 2 for event api.create
	// Attempt 3 for event api.create
	// Successfully processed event api.create
}

func ExampleEventBus_withFiltering() {
	logger, _ := zap.NewDevelopment()
	bus := eventbus.New(eventbus.WithLogger(logger))

	// Create a handler that only processes events from specific sources
	handler := eventbus.HandlerFunc(func(ctx context.Context, event *eventbus.Event) error {
		fmt.Printf("Processing event from %s: %s\n", event.Source, event.Type)
		return nil
	})

	// Create filter for specific sources
	filter := eventbus.EventFilter{
		Sources: []string{"tykctl-gateway"},
		Types:   []eventbus.EventType{EventTypeAPICreate},
	}

	// Wrap with filtered handler
	filteredHandler := eventbus.NewFilteredHandler(handler, filter)

	// Subscribe
	subscription, _ := bus.Subscribe(EventTypeAPICreate, filteredHandler)
	defer subscription.Unsubscribe()

	// Publish events from different sources
	event1 := eventbus.NewEvent(EventTypeAPICreate, map[string]interface{}{"api_id": "api-1"})
	event1.WithSource("tykctl-gateway")

	event2 := eventbus.NewEvent(EventTypeAPICreate, map[string]interface{}{"api_id": "api-2"})
	event2.WithSource("tykctl-portal")

	bus.Publish(event1) // Will be processed
	bus.Publish(event2) // Will be filtered out

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Output:
	// Processing event from tykctl-gateway: api.create
}