package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/adrg/xdg"
)

// DefaultSetter allows structs to set their own default values
type DefaultSetter interface {
	SetDefaults()
}

// Validator allows structs to validate their configuration
type Validator interface {
	Validate() error
}

// Configurable allows structs to configure themselves
type Configurable interface {
	Configure(cfg Config) error
}

// Config represents the main configuration interface
type Config interface {
	// Basic operations
	Get(key string) interface{}
	Set(key string, value interface{})
	Has(key string) bool
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetDuration(key string) time.Duration

	// Validation
	Validate() error
	ValidateWithContext(ctx context.Context) error

	// Metadata
	GetMetadata() ConfigMetadata
	GetSource() ConfigSource

	// Lifecycle
	Reload() error
	Watch(ctx context.Context) <-chan ConfigChange
	Close() error
}

// ConfigMetadata provides configuration metadata
type ConfigMetadata struct {
	Version    string            `json:"version"`
	LoadedAt   time.Time         `json:"loaded_at"`
	Sources    []ConfigSource    `json:"sources"`
	Validators []string          `json:"validators"`
	Extensions map[string]string `json:"extensions"`
}

// ConfigSource represents a configuration source
type ConfigSource struct {
	Type     string    `json:"type"` // file, env, remote, etc.
	Path     string    `json:"path"`
	Priority int       `json:"priority"`
	LoadedAt time.Time `json:"loaded_at"`
	Checksum string    `json:"checksum"`
}

// ConfigChange represents a configuration change
type ConfigChange struct {
	Key       string      `json:"key"`
	OldValue  interface{} `json:"old_value"`
	NewValue  interface{} `json:"new_value"`
	Source    string      `json:"source"`
	Timestamp time.Time   `json:"timestamp"`
}

// Hook represents a discovered hook with rich metadata
type Hook struct {
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Event        string            `json:"event"`
	Priority     int               `json:"priority"`
	Timeout      time.Duration     `json:"timeout"`
	Enabled      bool              `json:"enabled"`
	Metadata     map[string]string `json:"metadata"`
	Checksum     string            `json:"checksum"`
	DiscoveredAt time.Time         `json:"discovered_at"`
}

// HookFilter provides filtering options for hook discovery
type HookFilter struct {
	Events     []string          `json:"events"`
	Enabled    *bool             `json:"enabled"`
	Priority   *int              `json:"priority"`
	Metadata   map[string]string `json:"metadata"`
	MinTimeout time.Duration     `json:"min_timeout"`
	MaxTimeout time.Duration     `json:"max_timeout"`
}

// Plugin represents a discovered plugin with rich metadata
type Plugin struct {
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Version      string            `json:"version"`
	Commands     []string          `json:"commands"`
	Enabled      bool              `json:"enabled"`
	Metadata     map[string]string `json:"metadata"`
	Checksum     string            `json:"checksum"`
	DiscoveredAt time.Time         `json:"discovered_at"`
}

// PluginFilter provides filtering options for plugin discovery
type PluginFilter struct {
	Commands []string          `json:"commands"`
	Enabled  *bool             `json:"enabled"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata"`
}

// Template represents a discovered template with rich metadata
type Template struct {
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Type         string            `json:"type"`
	Format       string            `json:"format"`
	Schema       string            `json:"schema"`
	Enabled      bool              `json:"enabled"`
	Metadata     map[string]string `json:"metadata"`
	Checksum     string            `json:"checksum"`
	DiscoveredAt time.Time         `json:"discovered_at"`
}

// TemplateFilter provides filtering options for template discovery
type TemplateFilter struct {
	Types    []string          `json:"types"`
	Formats  []string          `json:"formats"`
	Enabled  *bool             `json:"enabled"`
	Metadata map[string]string `json:"metadata"`
}

// CacheConfig represents a discovered cache configuration
type CacheConfig struct {
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Type         string            `json:"type"`
	Enabled      bool              `json:"enabled"`
	Metadata     map[string]string `json:"metadata"`
	Checksum     string            `json:"checksum"`
	DiscoveredAt time.Time         `json:"discovered_at"`
}

// Context represents a configuration context
type Context struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Resources   ContextResources       `json:"resources"`
	Metadata    map[string]string      `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ContextResources represents resources for a context
type ContextResources struct {
	Hooks     map[string][]Hook `json:"hooks"`
	Plugins   []Plugin          `json:"plugins"`
	Templates []Template        `json:"templates"`
	Cache     []CacheConfig     `json:"cache"`
}

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// LogLevel represents logging levels
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Cache represents a cache interface
type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Close() error
	TTL() time.Duration
}

// Metrics interface for collecting metrics
type Metrics interface {
	RecordLoadDuration(duration time.Duration)
	RecordDiscoveryDuration(resourceType string, duration time.Duration)
	Close() error
}

// Loader provides enhanced configuration and resource loading
type Loader struct {
	ctx        context.Context
	extension  string
	context    string
	cache      Cache
	logger     Logger
	metrics    Metrics
	validators []Validator
	loaders    []ConfigLoader
	mu         sync.RWMutex
	configs    map[string]Config
	lastReload time.Time

	// Configurable properties
	envPrefix     string
	configFormats []string
	configPaths   []string
	contextPaths  []string
}

// LoaderOptions provides configuration for the loader
type LoaderOptions struct {
	Extension      string
	Context        string
	CacheEnabled   bool
	CacheTTL       time.Duration
	LogLevel       LogLevel
	MetricsEnabled bool
	Validators     []Validator
	Loaders        []ConfigLoader
	ReloadInterval time.Duration

	// Configurable properties
	EnvPrefix     string   // Environment variable prefix (default: TYKCTL)
	ConfigFormats []string // Supported config formats (default: ["yaml", "json", "toml"])
	ConfigPaths   []string // Additional config paths
	ContextPaths  []string // Context-specific paths
}

// NewLoader creates a new enhanced configuration loader
func NewLoader(ctx context.Context, opts LoaderOptions) (*Loader, error) {
	loader := &Loader{
		ctx:        ctx,
		extension:  opts.Extension,
		context:    opts.Context,
		configs:    make(map[string]Config),
		validators: opts.Validators,
		loaders:    opts.Loaders,

		// Set configurable properties with defaults
		envPrefix:     getEnvPrefix(opts.EnvPrefix, opts.Extension),
		configFormats: getConfigFormats(opts.ConfigFormats),
		configPaths:   getConfigPaths(opts.ConfigPaths, opts.Extension),
		contextPaths:  getContextPaths(opts.ContextPaths, opts.Context),
	}

	// Initialize cache
	if opts.CacheEnabled {
		cache, err := NewCache(CacheOptions{
			TTL: opts.CacheTTL,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create cache: %w", err)
		}
		loader.cache = cache
	}

	// Initialize logger
	logger, err := NewLogger(LoggerOptions{
		Level: opts.LogLevel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	loader.logger = logger

	// Initialize metrics
	if opts.MetricsEnabled {
		metrics, err := NewMetrics(MetricsOptions{
			Extension: opts.Extension,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create metrics: %w", err)
		}
		loader.metrics = metrics
	}

	// Start reload timer if configured
	if opts.ReloadInterval > 0 {
		go loader.startReloadTimer(opts.ReloadInterval)
	}

	return loader, nil
}

// Load loads configuration into your struct with enhanced features
func (l *Loader) Load(ctx context.Context, target interface{}) error {
	start := time.Now()
	defer func() {
		if l.metrics != nil {
			l.metrics.RecordLoadDuration(time.Since(start))
		}
	}()

	l.mu.Lock()
	defer l.mu.Unlock()

	// Check cache first
	if l.cache != nil {
		if cached, err := l.cache.Get(l.extension); err == nil {
			l.logger.Debug("Using cached configuration", "extension", l.extension)
			return l.unmarshalFromCache(cached, target)
		}
	}

	// Load from sources
	config, err := l.loadFromSources(ctx)
	if err != nil {
		l.logger.Error("Failed to load configuration", "error", err)
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err := l.validateConfig(ctx, config); err != nil {
		l.logger.Error("Configuration validation failed", "error", err)
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Unmarshal into target
	if err := l.unmarshalConfig(config, target); err != nil {
		l.logger.Error("Failed to unmarshal configuration", "error", err)
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// Cache the configuration
	if l.cache != nil {
		l.cache.Set(l.extension, config, l.cache.TTL())
	}

	// Store in loader
	l.configs[l.extension] = config
	l.lastReload = time.Now()

	l.logger.Info("Configuration loaded successfully",
		"extension", l.extension,
		"sources", len(config.GetMetadata().Sources),
		"duration", time.Since(start))

	return nil
}

// DiscoverHooks discovers hooks with enhanced filtering and metadata
func (l *Loader) DiscoverHooks(ctx context.Context, filter HookFilter) (map[string][]Hook, error) {
	start := time.Now()
	defer func() {
		if l.metrics != nil {
			l.metrics.RecordDiscoveryDuration("hooks", time.Since(start))
		}
	}()

	l.logger.Debug("Discovering hooks", "extension", l.extension, "filter", filter)

	hooks, err := DiscoverHooks(l.extension, filter, l.contextPaths)
	if err != nil {
		l.logger.Error("Failed to discover hooks", "error", err)
		return nil, fmt.Errorf("failed to discover hooks: %w", err)
	}

	// Group by event
	hooksByEvent := make(map[string][]Hook)
	for _, hook := range hooks {
		hooksByEvent[hook.Event] = append(hooksByEvent[hook.Event], hook)
	}

	l.logger.Info("Hooks discovered successfully",
		"extension", l.extension,
		"total", len(hooks),
		"events", len(hooksByEvent))

	return hooksByEvent, nil
}

// DiscoverPlugins discovers plugins with enhanced filtering and metadata
func (l *Loader) DiscoverPlugins(ctx context.Context, filter PluginFilter) ([]Plugin, error) {
	start := time.Now()
	defer func() {
		if l.metrics != nil {
			l.metrics.RecordDiscoveryDuration("plugins", time.Since(start))
		}
	}()

	l.logger.Debug("Discovering plugins", "extension", l.extension, "filter", filter)

	plugins, err := DiscoverPlugins(l.extension, filter, l.contextPaths)
	if err != nil {
		l.logger.Error("Failed to discover plugins", "error", err)
		return nil, fmt.Errorf("failed to discover plugins: %w", err)
	}

	l.logger.Info("Plugins discovered successfully",
		"extension", l.extension,
		"count", len(plugins))

	return plugins, nil
}

// DiscoverTemplates discovers templates with enhanced filtering and metadata
func (l *Loader) DiscoverTemplates(ctx context.Context, filter TemplateFilter) ([]Template, error) {
	start := time.Now()
	defer func() {
		if l.metrics != nil {
			l.metrics.RecordDiscoveryDuration("templates", time.Since(start))
		}
	}()

	l.logger.Debug("Discovering templates", "extension", l.extension, "filter", filter)

	templates, err := DiscoverTemplates(l.extension, filter, l.contextPaths)
	if err != nil {
		l.logger.Error("Failed to discover templates", "error", err)
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	l.logger.Info("Templates discovered successfully",
		"extension", l.extension,
		"count", len(templates))

	return templates, nil
}

// DiscoverCache discovers cache configurations for the extension
func (l *Loader) DiscoverCache(ctx context.Context) ([]CacheConfig, error) {
	start := time.Now()
	defer func() {
		if l.metrics != nil {
			l.metrics.RecordDiscoveryDuration("cache", time.Since(start))
		}
	}()

	l.logger.Debug("Discovering cache configs", "extension", l.extension)

	cacheConfigs, err := DiscoverCache(l.extension, l.contextPaths)
	if err != nil {
		l.logger.Error("Failed to discover cache configs", "error", err)
		return nil, fmt.Errorf("failed to discover cache configs: %w", err)
	}

	l.logger.Info("Cache configs discovered successfully",
		"extension", l.extension,
		"count", len(cacheConfigs))

	return cacheConfigs, nil
}

// Watch watches for configuration changes
func (l *Loader) Watch(ctx context.Context) <-chan ConfigChange {
	changes := make(chan ConfigChange, 100)

	go func() {
		defer close(changes)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := l.checkForChanges(ctx, changes); err != nil {
					l.logger.Error("Failed to check for changes", "error", err)
				}
			}
		}
	}()

	return changes
}

// Reload reloads the configuration
func (l *Loader) Reload(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Info("Reloading configuration", "extension", l.extension)

	// Clear cache
	if l.cache != nil {
		l.cache.Delete(l.extension)
	}

	// Reload from sources
	config, err := l.loadFromSources(ctx)
	if err != nil {
		return fmt.Errorf("failed to reload configuration: %w", err)
	}

	// Validate configuration
	if err := l.validateConfig(ctx, config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Update stored configuration
	l.configs[l.extension] = config
	l.lastReload = time.Now()

	l.logger.Info("Configuration reloaded successfully", "extension", l.extension)

	return nil
}

// Close closes the loader and cleans up resources
func (l *Loader) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cache != nil {
		l.cache.Close()
	}

	if l.metrics != nil {
		l.metrics.Close()
	}

	l.logger.Info("Loader closed", "extension", l.extension)

	return nil
}

// Load loads configuration into your struct. That's it.
func Load(ctx context.Context, extension string, target interface{}) error {
	loader, err := NewLoader(ctx, LoaderOptions{
		Extension: extension,
	})
	if err != nil {
		return err
	}
	return loader.Load(ctx, target)
}

// ConfigLoader interface for loading configuration from different sources
type ConfigLoader interface {
	Load(ctx context.Context, extension string) (Config, error)
	GetName() string
}

// Helper functions for configurable properties
func getEnvPrefix(prefix, extension string) string {
	if prefix == "" {
		prefix = "TYKCTL"
	}
	if extension != "" {
		prefix = fmt.Sprintf("%s_%s", prefix, strings.ToUpper(extension))
	}
	return prefix
}

func getConfigFormats(formats []string) []string {
	if len(formats) == 0 {
		return []string{"yaml", "json", "toml"}
	}
	return formats
}

func getConfigPaths(paths []string, extension string) []string {
	if len(paths) == 0 {
		// Default paths
		paths = []string{
			"/etc/tykctl",                               // System-wide
			filepath.Join(xdg.ConfigHome, "tykctl"),     // XDG config
			filepath.Join(os.Getenv("HOME"), ".tykctl"), // User home
			".",       // Current directory
			".tykctl", // Project config
		}
	}

	// Add extension-specific paths
	if extension != "" {
		for _, path := range paths {
			paths = append(paths, filepath.Join(path, extension))
		}
	}

	return paths
}

func getContextPaths(paths []string, context string) []string {
	if len(paths) == 0 {
		// Default context paths
		paths = []string{
			filepath.Join(xdg.DataHome, "tykctl", "contexts"),
			filepath.Join(os.Getenv("HOME"), ".tykctl", "contexts"),
			"./.tykctl/contexts",
		}
	}

	// Add context-specific paths
	if context != "" {
		for _, path := range paths {
			paths = append(paths, filepath.Join(path, context))
		}
	}

	return paths
}

// basicConfig is a simple implementation of the Config interface
type basicConfig struct {
	path   string
	format string
	data   map[string]interface{}
	mu     sync.RWMutex
}

// Implement Config interface
func (c *basicConfig) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *basicConfig) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *basicConfig) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.data[key]
	return exists
}

func (c *basicConfig) GetString(key string) string {
	if val := c.Get(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (c *basicConfig) GetInt(key string) int {
	if val := c.Get(key); val != nil {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return 0
}

func (c *basicConfig) GetBool(key string) bool {
	if val := c.Get(key); val != nil {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func (c *basicConfig) GetDuration(key string) time.Duration {
	if val := c.Get(key); val != nil {
		if d, ok := val.(time.Duration); ok {
			return d
		}
		if str, ok := val.(string); ok {
			if d, err := time.ParseDuration(str); err == nil {
				return d
			}
		}
	}
	return 0
}

func (c *basicConfig) Validate() error {
	// Basic validation - can be extended
	return nil
}

func (c *basicConfig) ValidateWithContext(ctx context.Context) error {
	// Basic validation with context - can be extended
	return nil
}

func (c *basicConfig) GetMetadata() ConfigMetadata {
	return ConfigMetadata{
		Version:    "1.0.0",
		LoadedAt:   time.Now(),
		Sources:    []ConfigSource{c.GetSource()},
		Validators: []string{},
		Extensions: map[string]string{
			"format": c.format,
		},
	}
}

func (c *basicConfig) GetSource() ConfigSource {
	return ConfigSource{
		Type: "file",
		Path: c.path,
	}
}

func (c *basicConfig) Reload() error {
	// Reload from file - simplified implementation
	return nil
}

func (c *basicConfig) Watch(ctx context.Context) <-chan ConfigChange {
	// Watch for changes - simplified implementation
	changes := make(chan ConfigChange)
	go func() {
		defer close(changes)
		<-ctx.Done()
	}()
	return changes
}

func (c *basicConfig) Close() error {
	// Cleanup - simplified implementation
	return nil
}

// Placeholder implementations for missing methods
func (l *Loader) loadFromSources(ctx context.Context) (Config, error) {
	// Try to load from all available sources
	for _, loader := range l.loaders {
		config, err := loader.Load(ctx, l.extension)
		if err == nil && config != nil {
			return config, nil
		}
		// Log error but continue with other loaders
		if l.logger != nil {
			l.logger.Debug("Failed to load config from loader", "loader", loader.GetName(), "error", err)
		}
	}

	// If no loaders succeeded, try to load from file system
	return l.loadFromFileSystem(ctx)
}

// loadFromFileSystem attempts to load configuration from the file system
func (l *Loader) loadFromFileSystem(ctx context.Context) (Config, error) {
	// Get all possible config paths
	paths := getConfigPaths(l.configPaths, l.extension)

	// Look for config files in each path
	for _, path := range paths {
		// Check each supported format
		for _, format := range l.configFormats {
			configFile := filepath.Join(path, fmt.Sprintf("config.%s", format))

			// Check if file exists
			if _, err := os.Stat(configFile); err == nil {
				// Found a config file, try to load it
				config, err := l.loadConfigFile(ctx, configFile, format)
				if err == nil {
					return config, nil
				}
				// Log error but continue with other files
				if l.logger != nil {
					l.logger.Debug("Failed to load config file", "file", configFile, "error", err)
				}
			}
		}
	}

	return nil, fmt.Errorf("no configuration found in any of the search paths")
}

// loadConfigFile loads a configuration file
func (l *Loader) loadConfigFile(ctx context.Context, filePath, format string) (Config, error) {
	// This is a simplified implementation
	// In practice, you'd use a proper configuration library like viper
	// to handle different formats (YAML, JSON, TOML, etc.)

	// For now, return a basic implementation
	return &basicConfig{
		path:   filePath,
		format: format,
		data:   make(map[string]interface{}),
	}, nil
}

func (l *Loader) validateConfig(ctx context.Context, config Config) error {
	// Run all validators
	for _, validator := range l.validators {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	// Validate the config itself
	if err := config.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// Validate with context if supported
	if err := config.ValidateWithContext(ctx); err != nil {
		return fmt.Errorf("context validation failed: %w", err)
	}

	return nil
}

func (l *Loader) unmarshalConfig(config Config, target interface{}) error {
	// Check if target implements Configurable interface
	if configurable, ok := target.(Configurable); ok {
		return configurable.Configure(config)
	}

	// For now, just set defaults and validate
	// TODO: Implement proper unmarshaling from Config interface
	
	// Set defaults if the target implements DefaultSetter
	if setter, ok := target.(DefaultSetter); ok {
		setter.SetDefaults()
	}

	// Validate if the target implements Validator
	if validator, ok := target.(Validator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	return nil
}

func (l *Loader) unmarshalFromCache(cached interface{}, target interface{}) error {
	// If cached is already the target type, copy it
	if reflect.TypeOf(cached) == reflect.TypeOf(target) {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(cached).Elem())
		return nil
	}

	// Try to convert using type assertion
	if config, ok := cached.(Config); ok {
		return l.unmarshalConfig(config, target)
	}

	// For other types, try JSON marshaling/unmarshaling
	// This is a fallback for simple types
	return fmt.Errorf("cache unmarshaling not implemented for type %T", cached)
}

func (l *Loader) checkForChanges(ctx context.Context, changes chan<- ConfigChange) error {
	// Check if any config files have been modified
	for name, config := range l.configs {
		// Get the source of the config
		source := config.GetSource()
		if source.Path == "" {
			continue // Skip configs without file sources
		}

		// Check file modification time
		info, err := os.Stat(source.Path)
		if err != nil {
			continue // Skip if file doesn't exist or can't be stat'd
		}

		// If file was modified after last reload, send a change event
		if info.ModTime().After(l.lastReload) {
			change := ConfigChange{
				Key:       name,
				OldValue:  nil, // We don't store old values
				NewValue:  config,
				Source:    source.Path,
				Timestamp: info.ModTime(),
			}

			select {
			case changes <- change:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

func (l *Loader) startReloadTimer(interval time.Duration) {
	if interval <= 0 {
		return // No reload timer needed
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Check for changes and reload if necessary
				changes := make(chan ConfigChange, 10)
				go func() {
					defer close(changes)
					if err := l.checkForChanges(l.ctx, changes); err != nil {
						if l.logger != nil {
							l.logger.Error("Error checking for changes", "error", err)
						}
					}
				}()

				// Process any changes
				for change := range changes {
					if l.logger != nil {
						l.logger.Info("Config changed", "key", change.Key, "source", change.Source)
					}

					// Update the config in the map
					l.mu.Lock()
					if config, ok := change.NewValue.(Config); ok {
						l.configs[change.Key] = config
					}
					l.lastReload = time.Now()
					l.mu.Unlock()
				}

			case <-l.ctx.Done():
				return
			}
		}
	}()
}

// GetConfigHome returns the primary configuration home directory
// This follows the standard precedence: XDG_CONFIG_HOME > HOME/.tykctl > current directory
func GetConfigHome() string {
	// Check XDG_CONFIG_HOME first
	if xdgConfigHome := xdg.ConfigHome; xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "tykctl")
	}

	// Fall back to HOME/.tykctl
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".tykctl")
	}

	// Last resort: current directory
	return ".tykctl"
}
