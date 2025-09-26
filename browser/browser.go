package browser

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Browser represents a browser interface
type Browser interface {
	Open(url string) error
}

// DefaultBrowser implements the Browser interface
type DefaultBrowser struct{}

// New creates a new browser instance
func New() Browser {
	return &DefaultBrowser{}
}

// Open opens a URL in the default browser
func (b *DefaultBrowser) Open(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		// Try common Linux browsers
		browsers := []string{"xdg-open", "firefox", "google-chrome", "chromium-browser"}
		for _, browser := range browsers {
			if _, err := exec.LookPath(browser); err == nil {
				cmd = exec.Command(browser, url)
				break
			}
		}
		if cmd == nil {
			return fmt.Errorf("no suitable browser found")
		}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// OpenURL opens a URL in the default browser
func OpenURL(url string) error {
	browser := New()
	return browser.Open(url)
}

// OpenWithBrowser opens a URL with a specific browser
func OpenWithBrowser(browserName, url string) error {
	cmd := exec.Command(browserName, url)
	return cmd.Run()
}

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
		browsers := []string{"xdg-open", "firefox", "google-chrome", "chromium-browser"}
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

// OpenInBackground opens a URL in the background
func OpenInBackground(url string) error {
	cmd := exec.Command("nohup", GetDefaultBrowser(), url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
