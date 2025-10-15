// Package telemetry provides anonymous usage analytics for tykctl-go.
package telemetry

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Middleware provides telemetry middleware for Cobra commands.
type Middleware struct {
	client Client
}

// NewMiddleware creates a new telemetry middleware.
func NewMiddleware(client Client) *Middleware {
	return &Middleware{
		client: client,
	}
}

// WrapCommand wraps a Cobra command with telemetry tracking.
func (m *Middleware) WrapCommand(cmd *cobra.Command) *cobra.Command {
	// Store the original RunE function
	originalRunE := cmd.RunE
	originalRun := cmd.Run
	
	// Wrap RunE if it exists
	if originalRunE != nil {
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return m.trackCommand(cmd, args, func() error {
				return originalRunE(cmd, args)
			})
		}
	}
	
	// Wrap Run if it exists
	if originalRun != nil {
		cmd.Run = func(cmd *cobra.Command, args []string) {
			m.trackCommand(cmd, args, func() error {
				originalRun(cmd, args)
				return nil
			})
		}
	}
	
	// Wrap all subcommands
	for _, subcmd := range cmd.Commands() {
		m.WrapCommand(subcmd)
	}
	
	return cmd
}

// trackCommand tracks command execution with telemetry.
func (m *Middleware) trackCommand(cmd *cobra.Command, args []string, fn func() error) error {
	start := time.Now()
	
	// Create command event
	event := NewEventBuilder(EventTypeCommand).
		Command(cmd.CommandPath()).
		Properties(map[string]interface{}{
			"args_count": len(args),
			"flags":      m.extractFlags(cmd),
		}).
		Build()
	
	// Set system information
	event.CLIVersion = GetCLIVersion()
	event.OS = runtime.GOOS
	event.Arch = runtime.GOARCH
	
	// Execute the command
	err := fn()
	
	// Set execution results
	event.Duration = time.Since(start).Milliseconds()
	event.Success = err == nil
	
	if err != nil {
		event.ErrorType = "command_error"
		event.ErrorMessage = sanitizeErrorMessage(err.Error())
	}
	
	// Track the event
	if trackErr := m.client.Track(event); trackErr != nil {
		// Log telemetry error but don't fail the command
		fmt.Printf("Warning: failed to track telemetry: %v\n", trackErr)
	}
	
	return err
}

// extractFlags extracts flag information from the command.
func (m *Middleware) extractFlags(cmd *cobra.Command) map[string]interface{} {
	flags := make(map[string]interface{})
	
	// Get local flags
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			flags[f.Name] = true
		}
	})
	
	// Get persistent flags
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			flags[f.Name] = true
		}
	})
	
	return flags
}

// TrackError tracks an error event.
func (m *Middleware) TrackError(errorType, message string, properties map[string]interface{}) error {
	event := NewEventBuilder(EventTypeError).
		Error(errorType, message).
		Properties(properties).
		Build()
	
	// Set system information
	event.CLIVersion = GetCLIVersion()
	event.OS = runtime.GOOS
	event.Arch = runtime.GOARCH
	
	return m.client.Track(event)
}

// TrackFeature tracks a feature usage event.
func (m *Middleware) TrackFeature(feature string, properties map[string]interface{}) error {
	event := NewEventBuilder(EventTypeFeature).
		Feature(feature).
		Success(true).
		Properties(properties).
		Build()
	
	// Set system information
	event.CLIVersion = GetCLIVersion()
	event.OS = runtime.GOOS
	event.Arch = runtime.GOARCH
	
	return m.client.Track(event)
}

// TrackPerformance tracks a performance metric event.
func (m *Middleware) TrackPerformance(operation string, duration time.Duration, success bool, properties map[string]interface{}) error {
	event := NewEventBuilder(EventTypePerformance).
		Command(operation).
		Duration(duration).
		Success(success).
		Properties(properties).
		Build()
	
	// Set system information
	event.CLIVersion = GetCLIVersion()
	event.OS = runtime.GOOS
	event.Arch = runtime.GOARCH
	
	return m.client.Track(event)
}

// CommandWrapper is a helper function to wrap a command with telemetry.
func CommandWrapper(client Client) func(*cobra.Command) *cobra.Command {
	middleware := NewMiddleware(client)
	return middleware.WrapCommand
}

// ErrorTracker is a helper for tracking errors.
type ErrorTracker struct {
	middleware *Middleware
}

// NewErrorTracker creates a new error tracker.
func NewErrorTracker(client Client) *ErrorTracker {
	return &ErrorTracker{
		middleware: NewMiddleware(client),
	}
}

// Track tracks an error.
func (et *ErrorTracker) Track(errorType, message string, properties map[string]interface{}) error {
	return et.middleware.TrackError(errorType, message, properties)
}

// FeatureTracker is a helper for tracking feature usage.
type FeatureTracker struct {
	middleware *Middleware
}

// NewFeatureTracker creates a new feature tracker.
func NewFeatureTracker(client Client) *FeatureTracker {
	return &FeatureTracker{
		middleware: NewMiddleware(client),
	}
}

// Track tracks a feature usage.
func (ft *FeatureTracker) Track(feature string, properties map[string]interface{}) error {
	return ft.middleware.TrackFeature(feature, properties)
}

// PerformanceTracker is a helper for tracking performance metrics.
type PerformanceTracker struct {
	middleware *Middleware
}

// NewPerformanceTracker creates a new performance tracker.
func NewPerformanceTracker(client Client) *PerformanceTracker {
	return &PerformanceTracker{
		middleware: NewMiddleware(client),
	}
}

// Track tracks a performance metric.
func (pt *PerformanceTracker) Track(operation string, duration time.Duration, success bool, properties map[string]interface{}) error {
	return pt.middleware.TrackPerformance(operation, duration, success, properties)
}

// TimingHelper provides a convenient way to track timing of operations.
type TimingHelper struct {
	start      time.Time
	operation  string
	tracker    *PerformanceTracker
	properties map[string]interface{}
}

// StartTiming starts timing an operation.
func StartTiming(tracker *PerformanceTracker, operation string, properties map[string]interface{}) *TimingHelper {
	return &TimingHelper{
		start:      time.Now(),
		operation:  operation,
		tracker:    tracker,
		properties: properties,
	}
}

// Finish finishes timing and tracks the performance metric.
func (th *TimingHelper) Finish(success bool) error {
	duration := time.Since(th.start)
	return th.tracker.Track(th.operation, duration, success, th.properties)
}

// FinishWithError finishes timing and tracks the performance metric with error information.
func (th *TimingHelper) FinishWithError(err error) error {
	success := err == nil
	properties := th.properties
	
	if err != nil {
		if properties == nil {
			properties = make(map[string]interface{})
		}
		properties["error_type"] = "operation_error"
		properties["error_message"] = sanitizeErrorMessage(err.Error())
	}
	
	return th.Finish(success)
}