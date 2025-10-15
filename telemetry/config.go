// Package telemetry provides anonymous usage analytics for tykctl-go.
package telemetry

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v3"
)

// ConfigManager manages telemetry configuration.
type ConfigManager struct {
	configPath string
	config     *Config
}

// NewConfigManager creates a new configuration manager.
func NewConfigManager() *ConfigManager {
	configPath := filepath.Join(xdg.ConfigHome, "tykctl", "telemetry.yaml")
	
	return &ConfigManager{
		configPath: configPath,
		config:     DefaultConfig(),
	}
}

// Load loads the configuration from the file.
func (cm *ConfigManager) Load() error {
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		return cm.Save()
	}
	
	// Read config file
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Unmarshal YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	cm.config = &config
	return nil
}

// Save saves the configuration to the file.
func (cm *ConfigManager) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Marshal to YAML
	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// GetConfig returns the current configuration.
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// SetConfig sets the configuration.
func (cm *ConfigManager) SetConfig(config *Config) {
	cm.config = config
}

// SetEnabled enables or disables telemetry.
func (cm *ConfigManager) SetEnabled(enabled bool) error {
	cm.config.Enabled = enabled
	return cm.Save()
}

// IsEnabled returns whether telemetry is enabled.
func (cm *ConfigManager) IsEnabled() bool {
	return cm.config.Enabled
}

// GetConfigPath returns the configuration file path.
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}

// EnvironmentConfigManager manages configuration from environment variables.
type EnvironmentConfigManager struct {
	config *Config
}

// NewEnvironmentConfigManager creates a new environment-based configuration manager.
func NewEnvironmentConfigManager() *EnvironmentConfigManager {
	config := DefaultConfig()
	
	// Override with environment variables if set
	if enabled := os.Getenv("TYKCTL_TELEMETRY_ENABLED"); enabled != "" {
		config.Enabled = enabled == "true" || enabled == "1"
	}
	
	if endpoint := os.Getenv("TYKCTL_TELEMETRY_ENDPOINT"); endpoint != "" {
		config.Endpoint = endpoint
	}
	
	if userAgent := os.Getenv("TYKCTL_TELEMETRY_USER_AGENT"); userAgent != "" {
		config.UserAgent = userAgent
	}
	
	return &EnvironmentConfigManager{
		config: config,
	}
}

// Load does nothing for environment config manager.
func (ecm *EnvironmentConfigManager) Load() error {
	return nil
}

// Save does nothing for environment config manager.
func (ecm *EnvironmentConfigManager) Save() error {
	return nil
}

// GetConfig returns the current configuration.
func (ecm *EnvironmentConfigManager) GetConfig() *Config {
	return ecm.config
}

// SetConfig sets the configuration.
func (ecm *EnvironmentConfigManager) SetConfig(config *Config) {
	ecm.config = config
}

// SetEnabled enables or disables telemetry.
func (ecm *EnvironmentConfigManager) SetEnabled(enabled bool) error {
	ecm.config.Enabled = enabled
	return nil
}

// IsEnabled returns whether telemetry is enabled.
func (ecm *EnvironmentConfigManager) IsEnabled() bool {
	return ecm.config.Enabled
}

// GetConfigPath returns an empty string for environment config manager.
func (ecm *EnvironmentConfigManager) GetConfigPath() string {
	return ""
}

// ConfigManagerInterface defines the interface for configuration managers.
type ConfigManagerInterface interface {
	Load() error
	Save() error
	GetConfig() *Config
	SetConfig(config *Config)
	SetEnabled(enabled bool) error
	IsEnabled() bool
	GetConfigPath() string
}

// GetDefaultStoragePath returns the default storage path for telemetry events.
func GetDefaultStoragePath() string {
	return filepath.Join(xdg.DataHome, "tykctl", "telemetry.json")
}

// GetDefaultConfigPath returns the default configuration path.
func GetDefaultConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "tykctl", "telemetry.yaml")
}

// GetSystemInfo returns system information for telemetry events.
func GetSystemInfo() map[string]string {
	return map[string]string{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}
}

// GetCLIVersion returns the CLI version information.
func GetCLIVersion() string {
	// This would typically be set at build time
	// For now, return a placeholder
	return "1.0.0"
}

// CreateDefaultClient creates a default telemetry client with file-based storage.
func CreateDefaultClient() (Client, error) {
	// Create config manager
	configManager := NewConfigManager()
	
	// Load configuration
	if err := configManager.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	config := configManager.GetConfig()
	
	// Create storage
	storagePath := GetDefaultStoragePath()
	storage := NewFileStorage(storagePath)
	
	// Create transport
	transport := NewHTTPTransport(
		config.Endpoint,
		config.UserAgent,
		config.Timeout,
	)
	
	// Create client
	client := NewClient(config, transport, storage)
	
	return client, nil
}

// CreateNoOpClient creates a no-operation telemetry client.
func CreateNoOpClient() Client {
	return NewNoOpClient()
}

// CreateClientFromConfig creates a telemetry client from a configuration.
func CreateClientFromConfig(config *Config, storage Storage) Client {
	transport := NewHTTPTransport(
		config.Endpoint,
		config.UserAgent,
		config.Timeout,
	)
	
	return NewClient(config, transport, storage)
}

// IsTelemetryDisabled returns true if telemetry is disabled via environment variable.
func IsTelemetryDisabled() bool {
	return os.Getenv("TYKCTL_NO_TELEMETRY") != "" || 
		   os.Getenv("TYKCTL_TELEMETRY_ENABLED") == "false" ||
		   os.Getenv("TYKCTL_TELEMETRY_ENABLED") == "0"
}

// GetTelemetryClient returns an appropriate telemetry client based on environment and configuration.
func GetTelemetryClient() (Client, error) {
	// Check if telemetry is disabled via environment
	if IsTelemetryDisabled() {
		return CreateNoOpClient(), nil
	}
	
	// Try to create default client
	client, err := CreateDefaultClient()
	if err != nil {
		// If there's an error, return no-op client
		return CreateNoOpClient(), nil
	}
	
	return client, nil
}