package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/edsonmichaque/tykctl-go/config"
)

// ExampleConfig represents a sample configuration struct
type ExampleConfig struct {
	URL          string        `mapstructure:"url" validate:"required,url"`
	Token        string        `mapstructure:"token" validate:"required,min_length=10"`
	MaxTokens    int           `mapstructure:"max_tokens" validate:"range=1,10000"`
	Temperature  float64       `mapstructure:"temperature" validate:"range=0,2"`
	CacheEnabled bool          `mapstructure:"cache_enabled"`
	CacheTTL     time.Duration `mapstructure:"cache_ttl"`
}

// SetDefaults implements DefaultSetter interface
func (c *ExampleConfig) SetDefaults() {
	if c.MaxTokens == 0 {
		c.MaxTokens = 512
	}
	if c.Temperature == 0 {
		c.Temperature = 0.7
	}
	if c.CacheTTL == 0 {
		c.CacheTTL = time.Hour
	}
}

// Validate implements Validator interface
func (c *ExampleConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("url is required")
	}
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}

func main() {
	ctx := context.Background()

	// Example 1: Basic configuration loading
	fmt.Println("=== Example 1: Basic Configuration Loading ===")
	loader, err := config.NewLoader(ctx, config.LoaderOptions{
		Extension:      "example",
		CacheEnabled:   true,
		CacheTTL:       5 * time.Minute,
		LogLevel:       config.LogLevelInfo,
		MetricsEnabled: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer loader.Close()

	var cfg ExampleConfig
	if err := loader.Load(ctx, &cfg); err != nil {
		log.Printf("Failed to load config: %v", err)
		// Use defaults
		cfg.SetDefaults()
	}

	fmt.Printf("URL: %s\n", cfg.URL)
	fmt.Printf("Max Tokens: %d\n", cfg.MaxTokens)
	fmt.Printf("Temperature: %.2f\n", cfg.Temperature)
	fmt.Printf("Cache TTL: %v\n", cfg.CacheTTL)

	// Example 2: Context management
	fmt.Println("\n=== Example 2: Context Management ===")
	// Create a simple logger for the context store
	logger, err := config.NewLogger(config.LoggerOptions{
		Level: config.LogLevelInfo,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctxStore, err := config.NewContextStore(config.ContextOptions{
		ConfigPath: "/tmp/tykctl-contexts",
		Logger:     logger,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create development context
	devConfig := map[string]interface{}{
		"url":         "https://dev-api.example.com",
		"max_tokens":  256,
		"temperature": 0.8,
	}

	if err := ctxStore.CreateContext("dev", devConfig); err != nil {
		log.Printf("Failed to create dev context: %v", err)
	}

	// Create production context
	prodConfig := map[string]interface{}{
		"url":         "https://api.example.com",
		"max_tokens":  1024,
		"temperature": 0.3,
	}

	if err := ctxStore.CreateContext("prod", prodConfig); err != nil {
		log.Printf("Failed to create prod context: %v", err)
	}

	// Switch to development context
	if err := ctxStore.SwitchContext("dev"); err != nil {
		log.Printf("Failed to switch to dev context: %v", err)
	}

	// Get current context
	currentCtx := ctxStore.GetCurrentContext()
	if currentCtx != nil {
		fmt.Printf("Current context: %s\n", currentCtx.Name)
		fmt.Printf("URL: %v\n", currentCtx.Config["url"])
		fmt.Printf("Max Tokens: %v\n", currentCtx.Config["max_tokens"])
	}

	// Example 3: Resource discovery
	fmt.Println("\n=== Example 3: Resource Discovery ===")

	// Discover hooks
	hooks, err := loader.DiscoverHooks(ctx, config.HookFilter{
		Events: []string{"pre-command", "post-command"},
	})
	if err != nil {
		log.Printf("Failed to discover hooks: %v", err)
	} else {
		fmt.Printf("Found hooks:\n")
		for event, eventHooks := range hooks {
			fmt.Printf("  %s (%d hooks):\n", event, len(eventHooks))
			for _, hook := range eventHooks {
				fmt.Printf("    - %s: %s\n", hook.Name, hook.Path)
			}
		}
	}

	// Discover plugins
	plugins, err := loader.DiscoverPlugins(ctx, config.PluginFilter{
		Commands: []string{"train", "evaluate"},
	})
	if err != nil {
		log.Printf("Failed to discover plugins: %v", err)
	} else {
		fmt.Printf("Found %d plugins:\n", len(plugins))
		for _, plugin := range plugins {
			fmt.Printf("  - %s: %s\n", plugin.Name, plugin.Path)
		}
	}

	// Example 4: Custom loader configuration
	fmt.Println("\n=== Example 4: Custom Loader Configuration ===")
	customLoader, err := config.NewLoader(ctx, config.LoaderOptions{
		Extension:      "my-app",
		Context:        "dev",
		CacheEnabled:   true,
		CacheTTL:       10 * time.Minute,
		LogLevel:       config.LogLevelDebug,
		MetricsEnabled: true,

		// Custom properties
		EnvPrefix:     "MYAPP",                     // Custom env prefix: MYAPP_*
		ConfigFormats: []string{"yaml", "json"},    // Only YAML and JSON
		ConfigPaths:   []string{"/custom/path"},    // Additional config paths
		ContextPaths:  []string{"/custom/context"}, // Additional context paths
	})
	if err != nil {
		log.Fatal(err)
	}
	defer customLoader.Close()

	fmt.Printf("Custom loader created successfully!\n")
	fmt.Printf("  Extension: my-app\n")
	fmt.Printf("  Context: dev\n")
	fmt.Printf("  Env Prefix: MYAPP\n")
	fmt.Printf("  Config Formats: [yaml json]\n")
	fmt.Printf("  Config Paths: [/custom/path]\n")
	fmt.Printf("  Context Paths: [/custom/context]\n")

	fmt.Println("\n=== All Examples Completed Successfully! ===")
}
