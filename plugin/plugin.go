package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/edsonmichaque/tykctl-go/config"
)

// Plugin represents a plugin executable
type Plugin struct {
	Name     string
	Path     string
	Extension string
}

// Manager handles plugin operations for an extension
type Manager struct {
	extension string
	config    ConfigProvider
	DefaultTimeout time.Duration  // Default timeout for plugin execution
}

// ConfigProvider provides configuration access for plugins
type ConfigProvider interface {
	GetConfigDir() string
	GetPluginDir(ctx context.Context) string
	GetPluginDiscoveryPaths(ctx context.Context) []string
}

// NewManager creates a new plugin manager for an extension
func NewManager(extension string, configProvider ConfigProvider) *Manager {
	return &Manager{
		extension: extension,
		config:    configProvider,
		DefaultTimeout: 0, // Will be set by GetConfiguredTimeout()
	}
}

// InstallFromDirectory installs plugins from a directory
func (m *Manager) InstallFromDirectory(ctx context.Context, sourceDir, pluginDir string, customName string) error {
	// Read the source directory to find executable files
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	var pluginExecutables []string
	var otherExecutables []string

	for _, entry := range entries {
		if !entry.IsDir() {
			// Check if file is executable by checking its mode
			info, err := entry.Info()
			if err == nil && isExecutable(info.Mode()) {
				pluginPrefix := fmt.Sprintf("tykctl-%s-", m.extension)
				if strings.HasPrefix(entry.Name(), pluginPrefix) {
					pluginExecutables = append(pluginExecutables, entry.Name())
				} else {
					otherExecutables = append(otherExecutables, entry.Name())
				}
			}
		}
	}

	if len(pluginExecutables) == 0 && len(otherExecutables) == 0 {
		return fmt.Errorf("no executable files found in directory %s", sourceDir)
	}

	// If we found plugin executables, install them directly
	if len(pluginExecutables) > 0 {
		return m.installPluginExecutables(ctx, sourceDir, pluginDir, pluginExecutables)
	}

	// Otherwise, treat as a plugin directory with multiple executables
	pluginName := filepath.Base(sourceDir)
	if customName != "" {
		pluginName = customName
	}

	pluginFileName := fmt.Sprintf("tykctl-%s-%s", m.extension, pluginName)
	pluginPath := filepath.Join(pluginDir, pluginFileName)

	// Check if plugin already exists
	if _, err := os.Stat(pluginPath); err == nil {
		return fmt.Errorf("plugin %s already exists at %s", pluginName, pluginPath)
	}

	if len(otherExecutables) == 1 {
		// Single executable file - copy it directly
		sourceFile := filepath.Join(sourceDir, otherExecutables[0])
		return m.copyExecutableFile(sourceFile, pluginPath)
	}

	// Multiple executable files - create a wrapper script
	return m.createPluginWrapper(sourceDir, pluginPath, otherExecutables)
}

// InstallFromFile installs a plugin from a single file
func (m *Manager) InstallFromFile(ctx context.Context, sourceFile, pluginDir, customName string) error {
	// Get the file name without extension as the plugin name
	pluginName := strings.TrimSuffix(filepath.Base(sourceFile), filepath.Ext(sourceFile))

	// Use custom name if provided
	if customName != "" {
		pluginName = customName
	}

	// Plugin files should follow the naming convention: tykctl-<extension>-<name>
	pluginFileName := fmt.Sprintf("tykctl-%s-%s", m.extension, pluginName)
	pluginPath := filepath.Join(pluginDir, pluginFileName)

	// Check if plugin already exists
	if _, err := os.Stat(pluginPath); err == nil {
		return fmt.Errorf("plugin %s already exists at %s", pluginName, pluginPath)
	}

	// Copy the file (will make it executable regardless of source permissions)
	return m.copyExecutableFile(sourceFile, pluginPath)
}

// CreateTemplate creates a new plugin template
func (m *Manager) CreateTemplate(pluginName, pluginDir string) error {
	pluginFileName := fmt.Sprintf("tykctl-%s-%s", m.extension, pluginName)
	pluginPath := filepath.Join(pluginDir, pluginFileName)

	// Check if plugin already exists
	if _, err := os.Stat(pluginPath); err == nil {
		return fmt.Errorf("plugin %s already exists at %s", pluginName, pluginPath)
	}

	// Create plugin template content
	templateContent := m.generatePluginTemplate(pluginName)

	// Write the template file
	if err := os.WriteFile(pluginPath, []byte(templateContent), 0755); err != nil {
		return fmt.Errorf("failed to create plugin template: %w", err)
	}

	fmt.Printf("Plugin %s created successfully at %s\n", pluginName, pluginPath)
	fmt.Printf("Plugin directory: %s\n", pluginDir)

	return nil
}

// Execute executes a plugin with the given arguments and environment setup
func (m *Manager) Execute(ctx context.Context, pluginPath string, args []string) error {
	timeout := m.GetConfiguredTimeout()
	return m.ExecuteWithTimeout(ctx, pluginPath, args, timeout)
}

// GetConfiguredTimeout returns the configured timeout for plugin execution
func (m *Manager) GetConfiguredTimeout() time.Duration {
	// Check extension-specific timeout first
	extensionUpper := strings.ToUpper(m.extension)
	if timeoutStr := os.Getenv(fmt.Sprintf("TYKCTL_%s_PLUGIN_TIMEOUT", extensionUpper)); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			return timeout
		}
	}
	
	// Check global plugin timeout
	if timeoutStr := os.Getenv("TYKCTL_PLUGIN_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			return timeout
		}
	}
	
	// Return default timeout (0 means no timeout)
	return m.DefaultTimeout
}

// ExecuteWithTimeout executes a plugin with a specific timeout
func (m *Manager) ExecuteWithTimeout(ctx context.Context, pluginPath string, args []string, timeout time.Duration) error {
	var execCtx context.Context
	var cancel context.CancelFunc

	// If timeout is 0, use the original context (no timeout)
	if timeout == 0 {
		execCtx = ctx
		cancel = func() {} // No-op cancel function
	} else {
		// Create a timeout context
		execCtx, cancel = context.WithTimeout(ctx, timeout)
	}
	defer cancel()

	// Create command to execute the plugin
	cmd := exec.CommandContext(execCtx, pluginPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set up environment variables for the plugin
	pluginEnv := m.setupPluginEnvironment(ctx, pluginPath)
	cmd.Env = append(os.Environ(), pluginEnv...)

	// Execute the plugin
	err := cmd.Run()
	if err != nil {
		// Check if the error is due to timeout (only if timeout was set)
		if timeout > 0 && execCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("plugin execution timed out after %v: %w", timeout, err)
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		return fmt.Errorf("failed to execute plugin: %w", err)
	}

	return nil
}

// DiscoverPlugins discovers all plugins for the extension
func (m *Manager) DiscoverPlugins(ctx context.Context) ([]Plugin, error) {
	var plugins []Plugin
	paths := m.config.GetPluginDiscoveryPaths(ctx)

	for _, path := range paths {
		if entries, err := os.ReadDir(path); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					// Check if file is executable by checking its mode
					info, err := entry.Info()
					if err == nil && isExecutable(info.Mode()) {
						pluginPrefix := fmt.Sprintf("tykctl-%s-", m.extension)
						if strings.HasPrefix(entry.Name(), pluginPrefix) {
							pluginName := strings.TrimPrefix(entry.Name(), pluginPrefix)
							plugins = append(plugins, Plugin{
								Name:      pluginName,
								Path:      filepath.Join(path, entry.Name()),
								Extension: m.extension,
							})
						}
					}
				}
			}
		}
	}

	return plugins, nil
}

// RemovePlugin removes a plugin from the configured directory
func (m *Manager) RemovePlugin(ctx context.Context, pluginName, pluginDir string) error {
	pluginFileName := fmt.Sprintf("tykctl-%s-%s", m.extension, pluginName)
	pluginPath := filepath.Join(pluginDir, pluginFileName)

	// Check if plugin exists
	if _, err := os.Stat(pluginPath); err != nil {
		return fmt.Errorf("plugin %s not found at %s", pluginName, pluginPath)
	}

	// Remove the plugin file
	if err := os.Remove(pluginPath); err != nil {
		return fmt.Errorf("failed to remove plugin: %w", err)
	}

	fmt.Printf("Plugin %s removed successfully from %s\n", pluginName, pluginPath)
	return nil
}

// installPluginExecutables installs multiple plugin executables from a directory
func (m *Manager) installPluginExecutables(ctx context.Context, sourceDir, pluginDir string, pluginExecutables []string) error {
	var installedPlugins []string

	for _, executable := range pluginExecutables {
		// Extract plugin name from the executable name
		pluginPrefix := fmt.Sprintf("tykctl-%s-", m.extension)
		pluginName := strings.TrimPrefix(executable, pluginPrefix)

		sourceFile := filepath.Join(sourceDir, executable)
		destFile := filepath.Join(pluginDir, executable)

		// Check if plugin already exists
		if _, err := os.Stat(destFile); err == nil {
			fmt.Printf("Warning: Plugin %s already exists, skipping\n", pluginName)
			continue
		}

		// Copy the executable
		if err := m.copyExecutableFile(sourceFile, destFile); err != nil {
			return fmt.Errorf("failed to install plugin %s: %w", pluginName, err)
		}

		installedPlugins = append(installedPlugins, pluginName)
	}

	if len(installedPlugins) > 0 {
		fmt.Printf("Successfully installed plugins: %s\n", strings.Join(installedPlugins, ", "))
	}

	return nil
}

// copyExecutableFile copies an executable file and makes it executable
func (m *Manager) copyExecutableFile(source, dest string) error {
	// Read the source file
	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Write to destination with executable permissions
	if err := os.WriteFile(dest, data, 0755); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	return nil
}

// createPluginWrapper creates a wrapper script for multiple executables
func (m *Manager) createPluginWrapper(sourceDir, pluginPath string, executables []string) error {
	// Generate wrapper script content
	wrapperContent := m.generateWrapperScript(sourceDir, executables)

	// Write the wrapper script
	if err := os.WriteFile(pluginPath, []byte(wrapperContent), 0755); err != nil {
		return fmt.Errorf("failed to create wrapper script: %w", err)
	}

	fmt.Printf("Plugin wrapper created successfully at %s\n", pluginPath)
	return nil
}

// generateWrapperScript generates a cross-platform wrapper script for multiple executables
func (m *Manager) generateWrapperScript(sourceDir string, executables []string) string {
	switch runtime.GOOS {
	case "windows":
		return m.generateWindowsWrapperScript(sourceDir, executables)
	default:
		return m.generateUnixWrapperScript(sourceDir, executables)
	}
}

// generateUnixWrapperScript generates a bash wrapper script for Unix-like systems
func (m *Manager) generateUnixWrapperScript(sourceDir string, executables []string) string {
	var script strings.Builder

	script.WriteString("#!/bin/bash\n")
	script.WriteString("# Plugin wrapper\n")
	script.WriteString("set -euo pipefail\n\n")

	script.WriteString(fmt.Sprintf("PLUGIN_DIR=\"%s\"\n", sourceDir))
	script.WriteString("COMMAND=\"${1:-help}\"\n\n")

	script.WriteString("case \"$COMMAND\" in\n")
	script.WriteString("    \"help\")\n")
	script.WriteString("        echo \"Available commands:\"\n")

	for _, exec := range executables {
		script.WriteString(fmt.Sprintf("        echo \"  - %s\"\n", exec))
	}

	script.WriteString("        ;;\n")

	for _, exec := range executables {
		script.WriteString(fmt.Sprintf("    \"%s\")\n", exec))
		script.WriteString(fmt.Sprintf("        exec \"$PLUGIN_DIR/%s\" \"$@\"\n", exec))
		script.WriteString("        ;;\n")
	}

	script.WriteString("    *)\n")
	script.WriteString("        echo \"Unknown command: $COMMAND\"\n")
	script.WriteString("        exit 1\n")
	script.WriteString("        ;;\n")
	script.WriteString("esac\n")

	return script.String()
}

// generateWindowsWrapperScript generates a batch wrapper script for Windows
func (m *Manager) generateWindowsWrapperScript(sourceDir string, executables []string) string {
	var script strings.Builder

	script.WriteString("@echo off\n")
	script.WriteString("REM Plugin wrapper\n")
	script.WriteString("setlocal enabledelayedexpansion\n\n")

	script.WriteString(fmt.Sprintf("set PLUGIN_DIR=%s\n", sourceDir))
	script.WriteString("set COMMAND=%1\n")
	script.WriteString("if \"%COMMAND%\"==\"\" set COMMAND=help\n\n")

	script.WriteString("if \"%COMMAND%\"==\"help\" (\n")
	script.WriteString("    echo Available commands:\n")

	for _, exec := range executables {
		script.WriteString(fmt.Sprintf("    echo   - %s\n", exec))
	}

	script.WriteString("    goto :eof\n")
	script.WriteString(")\n\n")

	for _, exec := range executables {
		script.WriteString(fmt.Sprintf("if \"%%COMMAND%%\"==\"%s\" (\n", exec))
		script.WriteString(fmt.Sprintf("    \"%%PLUGIN_DIR%%\\%s\" %%*\n", exec))
		script.WriteString("    goto :eof\n")
		script.WriteString(")\n\n")
	}

	script.WriteString("echo Unknown command: %COMMAND%\n")
	script.WriteString("exit /b 1\n")

	return script.String()
}

// generatePluginTemplate generates a cross-platform plugin template
func (m *Manager) generatePluginTemplate(pluginName string) string {
	switch runtime.GOOS {
	case "windows":
		return m.generateWindowsPluginTemplate(pluginName)
	default:
		return m.generateUnixPluginTemplate(pluginName)
	}
}

// generateUnixPluginTemplate generates a bash plugin template for Unix-like systems
func (m *Manager) generateUnixPluginTemplate(pluginName string) string {
	var template strings.Builder

	template.WriteString("#!/bin/bash\n")
	template.WriteString("# Plugin Template\n")
	template.WriteString("set -euo pipefail\n\n")

	template.WriteString(fmt.Sprintf("PLUGIN_NAME=\"%s\"\n", pluginName))
	template.WriteString("PLUGIN_VERSION=\"1.0.0\"\n\n")

	template.WriteString("log() {\n")
	template.WriteString("    echo \"[$(date '+%Y-%m-%d %H:%M:%S')] [$PLUGIN_NAME] $*\" >&2\n")
	template.WriteString("}\n\n")

	template.WriteString("case \"${1:-}\" in\n")
	template.WriteString("    \"version\")\n")
	template.WriteString("        echo \"$PLUGIN_VERSION\"\n")
	template.WriteString("        ;;\n")
	template.WriteString("    \"info\")\n")
	template.WriteString("        echo \"Name: $PLUGIN_NAME\"\n")
	template.WriteString("        echo \"Version: $PLUGIN_VERSION\"\n")
	template.WriteString(fmt.Sprintf("        echo \"Description: %s plugin for tykctl-%s\"\n", strings.Title(pluginName), m.extension))
	template.WriteString("        ;;\n")
	template.WriteString("    *)\n")
	template.WriteString("        echo \"Usage: $0 {version|info}\"\n")
	template.WriteString("        exit 1\n")
	template.WriteString("        ;;\n")
	template.WriteString("esac\n")

	return template.String()
}

// generateWindowsPluginTemplate generates a batch plugin template for Windows
func (m *Manager) generateWindowsPluginTemplate(pluginName string) string {
	var template strings.Builder

	template.WriteString("@echo off\n")
	template.WriteString("REM Plugin Template\n")
	template.WriteString("setlocal enabledelayedexpansion\n\n")

	template.WriteString(fmt.Sprintf("set PLUGIN_NAME=%s\n", pluginName))
	template.WriteString("set PLUGIN_VERSION=1.0.0\n\n")

	template.WriteString("set COMMAND=%1\n")
	template.WriteString("if \"%COMMAND%\"==\"\" set COMMAND=help\n\n")

	template.WriteString("if \"%COMMAND%\"==\"version\" (\n")
	template.WriteString("    echo %PLUGIN_VERSION%\n")
	template.WriteString("    goto :eof\n")
	template.WriteString(")\n\n")

	template.WriteString("if \"%COMMAND%\"==\"info\" (\n")
	template.WriteString("    echo Name: %PLUGIN_NAME%\n")
	template.WriteString("    echo Version: %PLUGIN_VERSION%\n")
	template.WriteString(fmt.Sprintf("    echo Description: %s plugin for tykctl-%s\n", strings.Title(pluginName), m.extension))
	template.WriteString("    goto :eof\n")
	template.WriteString(")\n\n")

	template.WriteString("if \"%COMMAND%\"==\"help\" (\n")
	template.WriteString("    echo Usage: %0 {version^|info}\n")
	template.WriteString("    goto :eof\n")
	template.WriteString(")\n\n")

	template.WriteString("echo Unknown command: %COMMAND%\n")
	template.WriteString("echo Usage: %0 {version^|info}\n")
	template.WriteString("exit /b 1\n")

	return template.String()
}

// setupPluginEnvironment sets up environment variables for plugin execution
func (m *Manager) setupPluginEnvironment(ctx context.Context, pluginPath string) []string {
	var envVars []string

	// Extract plugin name from path
	pluginName := strings.TrimSuffix(filepath.Base(pluginPath), filepath.Ext(pluginPath))
	pluginPrefix := fmt.Sprintf("tykctl-%s-", m.extension)
	if strings.HasPrefix(pluginName, pluginPrefix) {
		pluginName = strings.TrimPrefix(pluginName, pluginPrefix)
	}

	// Plugin identification
	envVars = append(envVars, fmt.Sprintf("TYKCTL_PLUGIN_NAME=%s", pluginName))
	envVars = append(envVars, fmt.Sprintf("TYKCTL_PLUGIN_PATH=%s", pluginPath))
	envVars = append(envVars, fmt.Sprintf("TYKCTL_PLUGIN_EXTENSION=%s", m.extension))

	// Plugin directories
	currentPluginDir := filepath.Dir(pluginPath)
	envVars = append(envVars, fmt.Sprintf("TYKCTL_PLUGIN_DIR=%s", currentPluginDir))

	// Extension-specific directories
	configDir := m.config.GetConfigDir()
	extensionPluginDir := m.config.GetPluginDir(ctx)

	extensionUpper := strings.ToUpper(m.extension)
	envVars = append(envVars, fmt.Sprintf("TYKCTL_%s_CONFIG_DIR=%s", extensionUpper, configDir))
	envVars = append(envVars, fmt.Sprintf("TYKCTL_%s_PLUGIN_DIR=%s", extensionUpper, extensionPluginDir))

	// Global TYKCTL directories
	envVars = append(envVars, fmt.Sprintf("TYKCTL_%s_GLOBAL_CONFIG_DIR=%s", extensionUpper, config.GetConfigHome()))

	// Extension-specific API configuration (if available)
	apiURLKey := fmt.Sprintf("TYK_%s_URL", extensionUpper)
	apiTokenKey := fmt.Sprintf("TYK_%s_TOKEN", extensionUpper)
	
	if apiURL := os.Getenv(apiURLKey); apiURL != "" {
		envVars = append(envVars, fmt.Sprintf("%s=%s", apiURLKey, apiURL))
	}
	if apiToken := os.Getenv(apiTokenKey); apiToken != "" {
		envVars = append(envVars, fmt.Sprintf("%s=%s", apiTokenKey, apiToken))
	}

	// Context information
	if contextName := os.Getenv("TYKCTL_CONTEXT"); contextName != "" {
		envVars = append(envVars, fmt.Sprintf("TYKCTL_%s_CONTEXT=%s", extensionUpper, contextName))
	}

	// Debug and verbose flags
	if debug := os.Getenv("TYKCTL_DEBUG"); debug != "" {
		envVars = append(envVars, fmt.Sprintf("TYKCTL_%s_DEBUG=%s", extensionUpper, debug))
	}
	if verbose := os.Getenv("TYKCTL_VERBOSE"); verbose != "" {
		envVars = append(envVars, fmt.Sprintf("TYKCTL_%s_VERBOSE=%s", extensionUpper, verbose))
	}

	// Plugin discovery paths
	discoveryPaths := m.config.GetPluginDiscoveryPaths(ctx)
	envVars = append(envVars, fmt.Sprintf("TYKCTL_%s_PLUGIN_DISCOVERY_PATHS=%s", extensionUpper, strings.Join(discoveryPaths, ":")))

	return envVars
}

// isExecutable checks if a file is executable across different platforms
func isExecutable(mode os.FileMode) bool {
	switch runtime.GOOS {
	case "windows":
		// On Windows, check for .exe, .bat, .cmd, .ps1 extensions
		// or if the file has executable permissions
		return mode&0111 != 0
	case "darwin", "linux", "freebsd", "openbsd", "netbsd":
		// On Unix-like systems, check executable permissions
		return mode&0111 != 0
	default:
		// For other platforms, use Unix-like behavior
		return mode&0111 != 0
	}
}

// getExecutableExtension returns the appropriate executable extension for the platform
func getExecutableExtension() string {
	switch runtime.GOOS {
	case "windows":
		return ".exe"
	default:
		return ""
	}
}

// getScriptExtension returns the appropriate script extension for the platform
func getScriptExtension() string {
	switch runtime.GOOS {
	case "windows":
		return ".bat"
	default:
		return ""
	}
}