// Package telemetry provides anonymous usage analytics for tykctl-go.
//
// The telemetry package collects anonymous usage statistics to help improve
// the CLI tool while respecting user privacy and providing opt-out capabilities.
//
// Key Features:
//   - Anonymous data collection (no personal information)
//   - Privacy-first design (no sensitive data transmitted)
//   - Opt-out capability (easily disabled by users)
//   - Configurable (flexible configuration options)
//   - Multiple transports (HTTP, file, mock)
//   - Cobra integration (seamless command tracking)
//   - Performance tracking (built-in metrics collection)
//   - Error tracking (automatic error reporting)
//
// Basic Usage:
//
//	// Create a telemetry client
//	client, err := telemetry.CreateDefaultClient()
//	if err != nil {
//		// Handle error or use no-op client
//		client = telemetry.CreateNoOpClient()
//	}
//	defer client.Close()
//
//	// Track a command execution
//	event := telemetry.NewEventBuilder(telemetry.EventTypeCommand).
//		Command("tykctl deploy").
//		Duration(2 * time.Second).
//		Success(true).
//		Properties(map[string]interface{}{
//			"environment": "production",
//			"region":      "us-east-1",
//		}).
//		Build()
//
//	if err := client.Track(event); err != nil {
//		// Log error but don't fail the operation
//		fmt.Printf("Warning: failed to track telemetry: %v\n", err)
//	}
//
// Cobra Integration:
//
//	// Create a Cobra command
//	cmd := &cobra.Command{
//		Use:   "deploy",
//		Short: "Deploy an application",
//		RunE: func(cmd *cobra.Command, args []string) error {
//			// Your command logic here
//			return nil
//		},
//	}
//
//	// Wrap the command with telemetry
//	middleware := telemetry.NewMiddleware(client)
//	middleware.WrapCommand(cmd)
//
// Configuration:
//
// The telemetry package uses YAML configuration files located at:
//   - Linux/macOS: ~/.config/tykctl/telemetry.yaml
//   - Windows: %APPDATA%\tykctl\telemetry.yaml
//
// Example configuration:
//
//	enabled: true
//	endpoint: "https://telemetry.tyk.io/v1/events"
//	batch_size: 100
//	flush_interval: "5m"
//	retry_attempts: 3
//	retry_delay: "1s"
//	timeout: "30s"
//	user_agent: "tykctl-go-telemetry/1.0"
//
// Environment Variables:
//
// You can override configuration using environment variables:
//
//	TYKCTL_NO_TELEMETRY=1                    # Disable telemetry
//	TYKCTL_TELEMETRY_ENABLED=false          # Disable telemetry
//	TYKCTL_TELEMETRY_ENDPOINT="https://..." # Custom endpoint
//	TYKCTL_TELEMETRY_USER_AGENT="my-app/1.0" # Custom user agent
//
// Privacy and Security:
//
// The telemetry package collects anonymous data including:
//   - Command usage patterns
//   - Performance metrics (execution times)
//   - Error types and frequencies
//   - Feature usage statistics
//   - System information (OS, architecture)
//   - CLI version information
//
// The telemetry package explicitly does NOT collect:
//   - Personal information (usernames, email addresses)
//   - Sensitive data (API keys, tokens, passwords)
//   - Repository names or content
//   - File contents or paths
//   - Network traffic or request/response data
//
// All events are automatically sanitized to remove sensitive information.
//
// Testing:
//
// Use mock transports and storage for testing:
//
//	mockTransport := telemetry.NewMockTransport()
//	mockStorage := telemetry.NewMockStorage()
//	config := telemetry.DefaultConfig()
//	client := telemetry.NewClient(config, mockTransport, mockStorage)
//
// For more examples and detailed documentation, see the README.md file.
package telemetry