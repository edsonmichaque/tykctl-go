// Package telemetry provides anonymous usage analytics for tykctl-go.
package telemetry

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// Example_basicUsage demonstrates basic telemetry usage.
func Example_basicUsage() {
	// Create a telemetry client
	client, err := CreateDefaultClient()
	if err != nil {
		fmt.Printf("Failed to create telemetry client: %v\n", err)
		return
	}
	defer client.Close()
	
	// Track a command execution
	event := NewEventBuilder(EventTypeCommand).
		Command("tykctl deploy").
		Duration(2 * time.Second).
		Success(true).
		Properties(map[string]interface{}{
			"environment": "production",
			"region":      "us-east-1",
		}).
		Build()
	
	if err := client.Track(event); err != nil {
		fmt.Printf("Failed to track event: %v\n", err)
	}
	
	// Flush any pending events
	if err := client.Flush(); err != nil {
		fmt.Printf("Failed to flush events: %v\n", err)
	}
}

// Example_cobraIntegration demonstrates integrating telemetry with Cobra commands.
func Example_cobraIntegration() {
	// Create a telemetry client
	client, err := CreateDefaultClient()
	if err != nil {
		fmt.Printf("Failed to create telemetry client: %v\n", err)
		return
	}
	defer client.Close()
	
	// Create a Cobra command
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Deploying application...")
			time.Sleep(100 * time.Millisecond) // Simulate work
			return nil
		},
	}
	
	// Wrap the command with telemetry
	middleware := NewMiddleware(client)
	middleware.WrapCommand(cmd)
	
	// Execute the command (this will automatically track telemetry)
	cmd.Execute()
}

// Example_errorTracking demonstrates tracking errors.
func Example_errorTracking() {
	// Create a telemetry client
	client, err := CreateDefaultClient()
	if err != nil {
		fmt.Printf("Failed to create telemetry client: %v\n", err)
		return
	}
	defer client.Close()
	
	// Create an error tracker
	errorTracker := NewErrorTracker(client)
	
	// Track an error
	if err := errorTracker.Track("deployment_failed", "Failed to connect to Kubernetes cluster", map[string]interface{}{
		"cluster": "production",
		"retry_count": 3,
	}); err != nil {
		fmt.Printf("Failed to track error: %v\n", err)
	}
}

// Example_featureTracking demonstrates tracking feature usage.
func Example_featureTracking() {
	// Create a telemetry client
	client, err := CreateDefaultClient()
	if err != nil {
		fmt.Printf("Failed to create telemetry client: %v\n", err)
		return
	}
	defer client.Close()
	
	// Create a feature tracker
	featureTracker := NewFeatureTracker(client)
	
	// Track feature usage
	if err := featureTracker.Track("auto_scaling", map[string]interface{}{
		"min_replicas": 2,
		"max_replicas": 10,
		"cpu_threshold": 70,
	}); err != nil {
		fmt.Printf("Failed to track feature: %v\n", err)
	}
}

// Example_performanceTracking demonstrates tracking performance metrics.
func Example_performanceTracking() {
	// Create a telemetry client
	client, err := CreateDefaultClient()
	if err != nil {
		fmt.Printf("Failed to create telemetry client: %v\n", err)
		return
	}
	defer client.Close()
	
	// Create a performance tracker
	perfTracker := NewPerformanceTracker(client)
	
	// Track performance metrics
	start := time.Now()
	
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	
	duration := time.Since(start)
	
	if err := perfTracker.Track("api_call", duration, true, map[string]interface{}{
		"endpoint": "/api/v1/deployments",
		"method":   "POST",
		"status_code": 200,
	}); err != nil {
		fmt.Printf("Failed to track performance: %v\n", err)
	}
}

// Example_timingHelper demonstrates using the timing helper.
func Example_timingHelper() {
	// Create a telemetry client
	client, err := CreateDefaultClient()
	if err != nil {
		fmt.Printf("Failed to create telemetry client: %v\n", err)
		return
	}
	defer client.Close()
	
	// Create a performance tracker
	perfTracker := NewPerformanceTracker(client)
	
	// Start timing an operation
	timer := StartTiming(perfTracker, "database_query", map[string]interface{}{
		"table": "users",
		"query_type": "SELECT",
	})
	
	// Simulate some work
	time.Sleep(50 * time.Millisecond)
	
	// Finish timing (success case)
	if err := timer.Finish(true); err != nil {
		fmt.Printf("Failed to track timing: %v\n", err)
	}
}

// Example_timingHelperWithError demonstrates using the timing helper with error handling.
func Example_timingHelperWithError() {
	// Create a telemetry client
	client, err := CreateDefaultClient()
	if err != nil {
		fmt.Printf("Failed to create telemetry client: %v\n", err)
		return
	}
	defer client.Close()
	
	// Create a performance tracker
	perfTracker := NewPerformanceTracker(client)
	
	// Start timing an operation
	timer := StartTiming(perfTracker, "api_call", map[string]interface{}{
		"endpoint": "/api/v1/users",
		"method":   "GET",
	})
	
	// Simulate some work
	time.Sleep(50 * time.Millisecond)
	
	// Simulate an error
	err = fmt.Errorf("connection timeout")
	
	// Finish timing with error
	if err := timer.FinishWithError(err); err != nil {
		fmt.Printf("Failed to track timing: %v\n", err)
	}
}

// Example_configurationManagement demonstrates configuration management.
func Example_configurationManagement() {
	// Create a configuration manager
	configManager := NewConfigManager()
	
	// Load existing configuration
	if err := configManager.Load(); err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}
	
	// Check if telemetry is enabled
	if configManager.IsEnabled() {
		fmt.Println("Telemetry is enabled")
	} else {
		fmt.Println("Telemetry is disabled")
	}
	
	// Disable telemetry
	if err := configManager.SetEnabled(false); err != nil {
		fmt.Printf("Failed to disable telemetry: %v\n", err)
		return
	}
	
	fmt.Println("Telemetry has been disabled")
}

// Example_customTransport demonstrates using a custom transport.
func Example_customTransport() {
	// Create a custom configuration
	config := &Config{
		Enabled:        true,
		Endpoint:       "https://custom-telemetry.example.com/v1/events",
		BatchSize:      50,
		FlushInterval:  2 * time.Minute,
		RetryAttempts: 5,
		RetryDelay:     2 * time.Second,
		Timeout:        60 * time.Second,
		UserAgent:      "tykctl-go-custom/1.0",
	}
	
	// Create custom storage
	storage := NewMemoryStorage()
	
	// Create client with custom configuration
	client := CreateClientFromConfig(config, storage)
	defer client.Close()
	
	// Track an event
	event := NewEventBuilder(EventTypeCommand).
		Command("tykctl custom-command").
		Success(true).
		Build()
	
	if err := client.Track(event); err != nil {
		fmt.Printf("Failed to track event: %v\n", err)
	}
}

// Example_mockTransport demonstrates using a mock transport for testing.
func Example_mockTransport() {
	// Create a mock transport
	mockTransport := NewMockTransport()
	
	// Create a configuration
	config := DefaultConfig()
	config.Endpoint = "https://test.example.com"
	
	// Create storage
	storage := NewMemoryStorage()
	
	// Create client with mock transport
	client := NewClient(config, mockTransport, storage)
	defer client.Close()
	
	// Track some events
	events := []*Event{
		NewEventBuilder(EventTypeCommand).Command("test1").Success(true).Build(),
		NewEventBuilder(EventTypeCommand).Command("test2").Success(false).Build(),
	}
	
	for _, event := range events {
		if err := client.Track(event); err != nil {
			fmt.Printf("Failed to track event: %v\n", err)
		}
	}
	
	// Flush events
	if err := client.Flush(); err != nil {
		fmt.Printf("Failed to flush events: %v\n", err)
	}
	
	// Check what was sent
	sentEvents := mockTransport.GetEvents()
	fmt.Printf("Sent %d batches of events\n", len(sentEvents))
	
	for i, batch := range sentEvents {
		fmt.Printf("Batch %d: %d events\n", i+1, len(batch))
	}
}