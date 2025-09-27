package browser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Browser represents a browser interface
type Browser interface {
	Open(url string) error
	OpenWithContext(ctx context.Context, url string) error
	OpenInBackground(url string) error
	OpenInNewWindow(url string) error
	OpenInNewTab(url string) error
	IsAvailable() bool
	GetName() string
}

// Config represents browser configuration
type Config struct {
	BrowserName string        // Specific browser to use
	Timeout     time.Duration // Timeout for browser operations
	Background  bool          // Open in background
	NewWindow   bool          // Open in new window
	NewTab      bool          // Open in new tab
	Args        []string      // Additional browser arguments
}

// DefaultBrowser implements the Browser interface
type DefaultBrowser struct {
	config Config
}

// New creates a new browser instance with default configuration
func New() Browser {
	return &DefaultBrowser{
		config: Config{
			Timeout: 30 * time.Second,
		},
	}
}

// NewWithConfig creates a new browser instance with custom configuration
func NewWithConfig(config Config) Browser {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	return &DefaultBrowser{config: config}
}

// Open opens a URL in the default browser
func (b *DefaultBrowser) Open(url string) error {
	return b.OpenWithContext(context.Background(), url)
}

// OpenWithContext opens a URL with context support
func (b *DefaultBrowser) OpenWithContext(ctx context.Context, url string) error {
	if err := b.validateURL(url); err != nil {
		return err
	}

	cmd, err := b.createCommand(url)
	if err != nil {
		return err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, b.config.Timeout)
	defer cancel()

	// Execute command with context
	return cmd.Run()
}

// OpenInBackground opens a URL in the background
func (b *DefaultBrowser) OpenInBackground(url string) error {
	if err := b.validateURL(url); err != nil {
		return err
	}

	cmd, err := b.createCommand(url)
	if err != nil {
		return err
	}

	// Start command in background
	return cmd.Start()
}

// OpenInNewWindow opens a URL in a new window
func (b *DefaultBrowser) OpenInNewWindow(url string) error {
	config := b.config
	config.NewWindow = true
	config.NewTab = false
	
	browser := NewWithConfig(config)
	return browser.Open(url)
}

// OpenInNewTab opens a URL in a new tab
func (b *DefaultBrowser) OpenInNewTab(url string) error {
	config := b.config
	config.NewTab = true
	config.NewWindow = false
	
	browser := NewWithConfig(config)
	return browser.Open(url)
}

// IsAvailable checks if the browser is available
func (b *DefaultBrowser) IsAvailable() bool {
	if b.config.BrowserName != "" {
		return IsAvailable(b.config.BrowserName)
	}
	return GetDefaultBrowser() != ""
}

// GetName returns the browser name
func (b *DefaultBrowser) GetName() string {
	if b.config.BrowserName != "" {
		return b.config.BrowserName
	}
	return GetDefaultBrowser()
}

// validateURL validates the URL format
func (b *DefaultBrowser) validateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	
	// Basic URL validation
	if !strings.HasPrefix(url, "http://") && 
	   !strings.HasPrefix(url, "https://") && 
	   !strings.HasPrefix(url, "file://") &&
	   !strings.HasPrefix(url, "ftp://") {
		return fmt.Errorf("invalid URL format: %s", url)
	}
	
	return nil
}

// createCommand creates the appropriate command for the platform
func (b *DefaultBrowser) createCommand(url string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	var args []string

	// Use specific browser if configured
	if b.config.BrowserName != "" {
		if !IsAvailable(b.config.BrowserName) {
			return nil, fmt.Errorf("browser '%s' is not available", b.config.BrowserName)
		}
		
		args = append(args, url)
		if b.config.NewWindow {
			args = append([]string{"--new-window"}, args...)
		} else if b.config.NewTab {
			args = append([]string{"--new-tab"}, args...)
		}
		args = append(args, b.config.Args...)
		
		cmd = exec.Command(b.config.BrowserName, args...)
		return cmd, nil
	}

	// Use platform default
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		args = []string{url}
		if b.config.NewWindow {
			args = append([]string{"-n"}, args...)
		}
		cmd = exec.Command("open", args...)
	case "linux":
		// Try common Linux browsers
		browsers := []string{"xdg-open", "firefox", "google-chrome", "chromium-browser", "chromium"}
		for _, browser := range browsers {
			if IsAvailable(browser) {
				args = []string{url}
				if b.config.NewWindow && (browser == "firefox" || browser == "google-chrome" || browser == "chromium-browser" || browser == "chromium") {
					args = append([]string{"--new-window"}, args...)
				} else if b.config.NewTab && (browser == "firefox" || browser == "google-chrome" || browser == "chromium-browser" || browser == "chromium") {
					args = append([]string{"--new-tab"}, args...)
				}
				args = append(args, b.config.Args...)
				cmd = exec.Command(browser, args...)
				break
			}
		}
		if cmd == nil {
			return nil, fmt.Errorf("no suitable browser found")
		}
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd, nil
}

// Convenience functions for easy usage

// OpenURL opens a URL in the default browser (convenience function)
func OpenURL(url string) error {
	browser := New()
	return browser.Open(url)
}

// OpenURLWithTimeout opens a URL with a timeout
func OpenURLWithTimeout(url string, timeout time.Duration) error {
	config := Config{Timeout: timeout}
	browser := NewWithConfig(config)
	return browser.Open(url)
}

// OpenURLInBackground opens a URL in the background
func OpenURLInBackground(url string) error {
	browser := New()
	return browser.OpenInBackground(url)
}

// OpenURLInNewWindow opens a URL in a new window
func OpenURLInNewWindow(url string) error {
	browser := New()
	return browser.OpenInNewWindow(url)
}

// OpenURLInNewTab opens a URL in a new tab
func OpenURLInNewTab(url string) error {
	browser := New()
	return browser.OpenInNewTab(url)
}

// OpenWithBrowser opens a URL with a specific browser
func OpenWithBrowser(browserName, url string) error {
	config := Config{BrowserName: browserName}
	browser := NewWithConfig(config)
	return browser.Open(url)
}

// OpenWithBrowserAndTimeout opens a URL with a specific browser and timeout
func OpenWithBrowserAndTimeout(browserName, url string, timeout time.Duration) error {
	config := Config{
		BrowserName: browserName,
		Timeout:     timeout,
	}
	browser := NewWithConfig(config)
	return browser.Open(url)
}

// Utility functions

// IsAvailable checks if a browser is available
func IsAvailable(browserName string) bool {
	_, err := exec.LookPath(browserName)
	return err == nil
}

// GetDefaultBrowser returns the default browser for the platform
func GetDefaultBrowser() string {
	switch runtime.GOOS {
	case "windows":
		return "rundll32"
	case "darwin":
		return "open"
	case "linux":
		// Try common Linux browsers in order of preference
		browsers := []string{"xdg-open", "firefox", "google-chrome", "chromium-browser", "chromium"}
		for _, browser := range browsers {
			if IsAvailable(browser) {
				return browser
			}
		}
		return "xdg-open"
	default:
		return ""
	}
}

// GetAvailableBrowsers returns a list of available browsers on the system
func GetAvailableBrowsers() []string {
	var browsers []string
	
	switch runtime.GOOS {
	case "windows":
		candidates := []string{"chrome", "firefox", "edge", "iexplore"}
		for _, browser := range candidates {
			if IsAvailable(browser) {
				browsers = append(browsers, browser)
			}
		}
	case "darwin":
		candidates := []string{"open", "chrome", "firefox", "safari"}
		for _, browser := range candidates {
			if IsAvailable(browser) {
				browsers = append(browsers, browser)
			}
		}
	case "linux":
		candidates := []string{"xdg-open", "firefox", "google-chrome", "chromium-browser", "chromium", "konqueror", "opera"}
		for _, browser := range candidates {
			if IsAvailable(browser) {
				browsers = append(browsers, browser)
			}
		}
	}
	
	return browsers
}

// ValidateURL validates a URL format
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	
	// Basic URL validation
	if !strings.HasPrefix(url, "http://") && 
	   !strings.HasPrefix(url, "https://") && 
	   !strings.HasPrefix(url, "file://") &&
	   !strings.HasPrefix(url, "ftp://") {
		return fmt.Errorf("invalid URL format: %s", url)
	}
	
	return nil
}

// BrowserInfo represents information about a browser
type BrowserInfo struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
	Default   bool   `json:"default"`
}

// GetBrowserInfo returns information about browsers on the system
func GetBrowserInfo() []BrowserInfo {
	var info []BrowserInfo
	
	// Get default browser
	defaultBrowser := GetDefaultBrowser()
	
	// Get available browsers
	availableBrowsers := GetAvailableBrowsers()
	
	// Create info for each browser
	for _, browser := range availableBrowsers {
		info = append(info, BrowserInfo{
			Name:      browser,
			Available: true,
			Default:   browser == defaultBrowser,
		})
	}
	
	return info
}

// Error types for better error handling

// BrowserError represents a browser-specific error
type BrowserError struct {
	Browser string
	URL     string
	Err     error
}

func (e *BrowserError) Error() string {
	return fmt.Sprintf("browser error (%s): %v", e.Browser, e.Err)
}

func (e *BrowserError) Unwrap() error {
	return e.Err
}

// URLValidationError represents a URL validation error
type URLValidationError struct {
	URL string
	Err error
}

func (e *URLValidationError) Error() string {
	return fmt.Sprintf("URL validation error (%s): %v", e.URL, e.Err)
}

func (e *URLValidationError) Unwrap() error {
	return e.Err
}
