package browser

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example demonstrates various ways to use the improved browser package
func Example() {
	fmt.Println("=== Browser Package Examples ===\n")

	// Example 1: Basic usage - simplest way to open a URL
	fmt.Println("1. Basic Usage:")
	err := OpenURL("https://example.com")
	if err != nil {
		log.Printf("Error opening URL: %v", err)
	} else {
		fmt.Println("✓ URL opened successfully")
	}
	fmt.Println()

	// Example 2: Using the Browser interface
	fmt.Println("2. Using Browser Interface:")
	b := New()
	if b.IsAvailable() {
		fmt.Printf("✓ Browser available: %s\n", b.GetName())
		err := b.Open("https://github.com")
		if err != nil {
			log.Printf("Error: %v", err)
		}
	} else {
		fmt.Println("✗ No browser available")
	}
	fmt.Println()

	// Example 3: Opening with timeout
	fmt.Println("3. Opening with Timeout:")
	err = OpenURLWithTimeout("https://google.com", 10*time.Second)
	if err != nil {
		log.Printf("Error with timeout: %v", err)
	} else {
		fmt.Println("✓ URL opened with timeout")
	}
	fmt.Println()

	// Example 4: Opening in background
	fmt.Println("4. Opening in Background:")
	err = OpenURLInBackground("https://stackoverflow.com")
	if err != nil {
		log.Printf("Error opening in background: %v", err)
	} else {
		fmt.Println("✓ URL opened in background")
	}
	fmt.Println()

	// Example 5: Opening in new window/tab
	fmt.Println("5. Opening in New Window/Tab:")
	err = OpenURLInNewWindow("https://golang.org")
	if err != nil {
		log.Printf("Error opening in new window: %v", err)
	} else {
		fmt.Println("✓ URL opened in new window")
	}

	err = OpenURLInNewTab("https://pkg.go.dev")
	if err != nil {
		log.Printf("Error opening in new tab: %v", err)
	} else {
		fmt.Println("✓ URL opened in new tab")
	}
	fmt.Println()

	// Example 6: Using specific browser
	fmt.Println("6. Using Specific Browser:")
	if IsAvailable("firefox") {
		err = OpenWithBrowser("firefox", "https://mozilla.org")
		if err != nil {
			log.Printf("Error opening with Firefox: %v", err)
		} else {
			fmt.Println("✓ URL opened with Firefox")
		}
	} else {
		fmt.Println("✗ Firefox not available")
	}
	fmt.Println()

	// Example 7: Custom configuration
	fmt.Println("7. Custom Configuration:")
	config := Config{
		BrowserName: "chrome",
		Timeout:     15 * time.Second,
		NewTab:      true,
		Args:        []string{"--incognito"},
	}
	
	browser := NewWithConfig(config)
	if browser.IsAvailable() {
		err := browser.Open("https://example.com")
		if err != nil {
			log.Printf("Error with custom config: %v", err)
		} else {
			fmt.Println("✓ URL opened with custom configuration")
		}
	} else {
		fmt.Println("✗ Chrome not available")
	}
	fmt.Println()

	// Example 8: Context-aware opening
	fmt.Println("8. Context-Aware Opening:")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	browser2 := New()
	err = browser2.OpenWithContext(ctx, "https://example.com")
	if err != nil {
		log.Printf("Error with context: %v", err)
	} else {
		fmt.Println("✓ URL opened with context")
	}
	fmt.Println()

	// Example 9: Browser information
	fmt.Println("9. Browser Information:")
	fmt.Printf("Default browser: %s\n", GetDefaultBrowser())
	
	availableBrowsers := GetAvailableBrowsers()
	fmt.Printf("Available browsers: %v\n", availableBrowsers)
	
	browserInfo := GetBrowserInfo()
	fmt.Println("Browser details:")
	for _, info := range browserInfo {
		defaultMark := ""
		if info.Default {
			defaultMark = " (default)"
		}
		fmt.Printf("  - %s: available=%t%s\n", info.Name, info.Available, defaultMark)
	}
	fmt.Println()

	// Example 10: URL validation
	fmt.Println("10. URL Validation:")
	testURLs := []string{
		"https://example.com",
		"http://localhost:8080",
		"file:///path/to/file.html",
		"invalid-url",
		"",
	}
	
	for _, url := range testURLs {
		err := ValidateURL(url)
		if err != nil {
			fmt.Printf("✗ Invalid URL '%s': %v\n", url, err)
		} else {
			fmt.Printf("✓ Valid URL: %s\n", url)
		}
	}
	fmt.Println()

	// Example 11: Error handling
	fmt.Println("11. Error Handling:")
	err = OpenURL("invalid-url")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
	
	err = OpenWithBrowser("nonexistent-browser", "https://example.com")
	if err != nil {
		fmt.Printf("Expected browser error: %v\n", err)
	}
	fmt.Println()

	fmt.Println("=== Examples Complete ===")
}

// ExampleAdvanced demonstrates advanced usage patterns
func ExampleAdvanced() {
	fmt.Println("=== Advanced Browser Usage ===\n")

	// Example 1: Browser selection with fallback
	fmt.Println("1. Browser Selection with Fallback:")
	preferredBrowsers := []string{"chrome", "firefox", "safari"}
	var selectedBrowser string
	
	for _, browser := range preferredBrowsers {
		if IsAvailable(browser) {
			selectedBrowser = browser
			break
		}
	}
	
	if selectedBrowser != "" {
		fmt.Printf("Using preferred browser: %s\n", selectedBrowser)
		err := OpenWithBrowser(selectedBrowser, "https://example.com")
		if err != nil {
			log.Printf("Error: %v", err)
		}
	} else {
		fmt.Println("No preferred browser available, using default")
		err := OpenURL("https://example.com")
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
	fmt.Println()

	// Example 2: Batch URL opening
	fmt.Println("2. Batch URL Opening:")
	urls := []string{
		"https://github.com",
		"https://stackoverflow.com",
		"https://golang.org",
	}
	
	browser := New()
	for i, url := range urls {
		if i == 0 {
			// Open first URL normally
			err := browser.Open(url)
			if err != nil {
				log.Printf("Error opening %s: %v", url, err)
			}
		} else {
			// Open subsequent URLs in new tabs
			err := browser.OpenInNewTab(url)
			if err != nil {
				log.Printf("Error opening %s in new tab: %v", url, err)
			}
		}
		time.Sleep(500 * time.Millisecond) // Small delay between opens
	}
	fmt.Println("✓ Batch URL opening completed")
	fmt.Println()

	// Example 3: Conditional opening based on environment
	fmt.Println("3. Conditional Opening:")
	config := Config{
		Timeout: 10 * time.Second,
	}
	
	// Check if we're in a headless environment
	if GetDefaultBrowser() == "" {
		fmt.Println("No browser available in headless environment")
		return
	}
	
	// Open with appropriate settings
	if IsAvailable("chrome") {
		config.BrowserName = "chrome"
		config.Args = []string{"--new-window"}
	} else {
		config.NewWindow = true
	}
	
	browser2 := NewWithConfig(config)
	err := browser2.Open("https://example.com")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("✓ URL opened with conditional configuration")
	}
	fmt.Println()

	fmt.Println("=== Advanced Examples Complete ===")
}

// ExampleErrorHandling demonstrates comprehensive error handling
func ExampleErrorHandling() {
	fmt.Println("=== Error Handling Examples ===\n")

	// Example 1: URL validation errors
	fmt.Println("1. URL Validation Errors:")
	invalidURLs := []string{
		"",
		"not-a-url",
		"ftp://example.com", // Valid but might not be supported by browser
	}
	
	for _, url := range invalidURLs {
		err := ValidateURL(url)
		if err != nil {
			fmt.Printf("Validation error for '%s': %v\n", url, err)
		}
	}
	fmt.Println()

	// Example 2: Browser availability errors
	fmt.Println("2. Browser Availability Errors:")
	unavailableBrowsers := []string{"nonexistent-browser", "fake-browser"}
	
	for _, browser := range unavailableBrowsers {
		if !IsAvailable(browser) {
			fmt.Printf("Browser '%s' is not available\n", browser)
		}
	}
	fmt.Println()

	// Example 3: Timeout errors
	fmt.Println("3. Timeout Errors:")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	browser := New()
	err := browser.OpenWithContext(ctx, "https://example.com")
	if err != nil {
		fmt.Printf("Timeout error: %v\n", err)
	}
	fmt.Println()

	// Example 4: Graceful error handling
	fmt.Println("4. Graceful Error Handling:")
	err = OpenURL("https://example.com")
	if err != nil {
		// Handle different types of errors
		switch {
		case err.Error() == "URL cannot be empty":
			fmt.Println("Empty URL provided")
		case err.Error() == "no suitable browser found":
			fmt.Println("No browser available on this system")
		default:
			fmt.Printf("Unexpected error: %v\n", err)
		}
	} else {
		fmt.Println("URL opened successfully")
	}
	fmt.Println()

	fmt.Println("=== Error Handling Examples Complete ===")
}