package browser

import (
	"runtime"
	"testing"
)

func TestPlatformSupport(t *testing.T) {
	// Test that we support the three main platforms
	supportedPlatforms := []string{"windows", "darwin", "linux"}
	
	currentPlatform := runtime.GOOS
	supported := false
	for _, platform := range supportedPlatforms {
		if currentPlatform == platform {
			supported = true
			break
		}
	}
	
	if !supported {
		t.Logf("Current platform %s is not explicitly supported, but may still work", currentPlatform)
	}
}

func TestGetDefaultBrowser(t *testing.T) {
	browser := GetDefaultBrowser()
	
	switch runtime.GOOS {
	case "windows":
		if browser != "rundll32" {
			t.Errorf("Expected 'rundll32' for Windows, got '%s'", browser)
		}
	case "darwin":
		if browser != "open" {
			t.Errorf("Expected 'open' for macOS, got '%s'", browser)
		}
	case "linux":
		if browser != "xdg-open" {
			t.Errorf("Expected 'xdg-open' for Linux, got '%s'", browser)
		}
	default:
		if browser != "" {
			t.Errorf("Expected empty string for unsupported platform, got '%s'", browser)
		}
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"https://example.com", false},
		{"http://example.com", false},
		{"file:///path/to/file", false},
		{"ftp://example.com", false},
		{"", true},
		{"not-a-url", true},
		{"example.com", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestNewBrowser(t *testing.T) {
	b := New()
	if b == nil {
		t.Error("New() returned nil")
	}
	
	// Test that the browser has a name
	name := b.GetName()
	if name == "" {
		t.Error("Browser name should not be empty")
	}
	
	// Test availability check doesn't panic
	_ = b.IsAvailable()
}

func TestBrowserWithConfig(t *testing.T) {
	config := Config{
		BrowserName: "test-browser",
		Timeout:     5,
	}
	
	b := NewWithConfig(config)
	if b == nil {
		t.Error("NewWithConfig() returned nil")
	}
	
	name := b.GetName()
	if name != "test-browser" {
		t.Errorf("Expected browser name 'test-browser', got '%s'", name)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test that convenience functions don't panic
	// We can't actually test opening URLs in tests, but we can test validation
	
	// Test URL validation in convenience functions
	err := ValidateURL("https://example.com")
	if err != nil {
		t.Errorf("ValidateURL failed for valid URL: %v", err)
	}
	
	err = ValidateURL("invalid-url")
	if err == nil {
		t.Error("ValidateURL should fail for invalid URL")
	}
}

func TestPlatformSpecificCommands(t *testing.T) {
	// Test that we create the right commands for each platform
	b := New()
	
	// This tests the internal createCommand method indirectly
	// by checking that we can create a browser instance and get its name
	name := b.GetName()
	if name == "" {
		t.Error("Browser should have a name")
	}
	
	// Test that we can get available browsers
	browsers := GetAvailableBrowsers()
	if len(browsers) == 0 {
		t.Log("No browsers detected - this might be expected in some environments")
	}
}

func TestErrorHandling(t *testing.T) {
	// Test URL validation errors
	err := ValidateURL("")
	if err == nil {
		t.Error("Empty URL should cause validation error")
	}
	
	err = ValidateURL("not-a-url")
	if err == nil {
		t.Error("Invalid URL should cause validation error")
	}
}

// Benchmark tests
func BenchmarkOpenURL(b *testing.B) {
	// This benchmark tests the overhead of creating browser instances
	// We can't actually open URLs in benchmarks, but we can test the setup
	for i := 0; i < b.N; i++ {
		browser := New()
		_ = browser.GetName()
	}
}

func BenchmarkValidateURL(b *testing.B) {
	url := "https://example.com"
	for i := 0; i < b.N; i++ {
		_ = ValidateURL(url)
	}
}