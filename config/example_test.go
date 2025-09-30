package config

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
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

// TestBasicUsage demonstrates basic usage of the config package
func TestBasicUsage(t *testing.T) {
	ctx := context.Background()

	// Create a simple loader
	loader, err := NewLoader(ctx, LoaderOptions{
		Extension:      "example",
		CacheEnabled:   true,
		CacheTTL:       5 * time.Minute,
		LogLevel:       LogLevelInfo,
		MetricsEnabled: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer loader.Close()

	// Load configuration
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

	// Discover hooks
	hooks, err := loader.DiscoverHooks(ctx, HookFilter{
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
	plugins, err := loader.DiscoverPlugins(ctx, PluginFilter{
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

	// Discover templates
	templates, err := loader.DiscoverTemplates(ctx, TemplateFilter{
		Types:   []string{"output", "input"},
		Formats: []string{"json", "yaml"},
	})
	if err != nil {
		log.Printf("Failed to discover templates: %v", err)
	} else {
		fmt.Printf("Found %d templates:\n", len(templates))
		for _, template := range templates {
			fmt.Printf("  - %s: %s (type: %s, format: %s)\n",
				template.Name, template.Path, template.Type, template.Format)
		}
	}
}

// TestContextStore demonstrates context management
func TestContextStore(t *testing.T) {
	// Create context store
	ctxStore, err := NewContextStore(ContextOptions{
		ConfigPath: "/tmp/tykctl-contexts",
		Logger:     &simpleLogger{level: LogLevelInfo, logger: log.New(log.Writer(), "", log.LstdFlags)},
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

		// Test GetContextConfig with error handling
		if url, err := ctxStore.GetContextConfig("dev", "url"); err != nil {
			log.Printf("Failed to get URL config: %v", err)
		} else {
			fmt.Printf("URL from GetContextConfig: %v\n", url)
		}

		// Test SetContextConfig with error handling
		if err := ctxStore.SetContextConfig("dev", "new_key", "new_value"); err != nil {
			log.Printf("Failed to set config: %v", err)
		} else {
			fmt.Printf("Successfully set new config key\n")
		}
	}

	// List all contexts
	contexts := ctxStore.ListContexts()
	fmt.Printf("Available contexts:\n")
	for _, context := range contexts {
		fmt.Printf("  - %s: %s\n", context.Name, context.Description)
	}
}

// TestCustomLoader demonstrates configurable loader properties
func TestCustomLoader(t *testing.T) {
	ctx := context.Background()

	// Create loader with custom properties
	loader, err := NewLoader(ctx, LoaderOptions{
		Extension:      "my-app",
		Context:        "dev",
		CacheEnabled:   true,
		CacheTTL:       10 * time.Minute,
		LogLevel:       LogLevelDebug,
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
	defer loader.Close()

	fmt.Printf("Loader created with custom properties:\n")
	fmt.Printf("  Extension: %s\n", loader.extension)
	fmt.Printf("  Context: %s\n", loader.context)
	fmt.Printf("  Env Prefix: %s\n", loader.envPrefix)
	fmt.Printf("  Config Formats: %v\n", loader.configFormats)
	fmt.Printf("  Config Paths: %v\n", loader.configPaths)
	fmt.Printf("  Context Paths: %v\n", loader.contextPaths)
}
