package main

import (
	"context"
	"fmt"
	"os"

	"github.com/edsonmichaque/tykctl-go/extension"
	"go.uber.org/zap"
)

func main() {
	// Create extension installer with functional options
	configDir := "/tmp/tykctl-config" // In real usage, this would come from config

	// Example 1: Basic installer
	installer := extension.NewInstaller(configDir)

	// Example 2: Installer with GitHub token (commented out)
	// installer := extension.NewInstaller(configDir, extension.WithGitHubToken("your-token-here"))

	// Example 3: Installer with custom logger (commented out)
	// customLogger := zap.NewExample()
	// installer := extension.NewInstaller(configDir, extension.WithLogger(customLogger))

	// Note: The zap import is used in the commented examples above
	_ = zap.NewExample

	ctx := context.Background()

	// Example: Search for extensions
	fmt.Println("Searching for extensions...")
	extensions, err := installer.SearchExtensions(ctx, "tyk", 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching extensions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d extensions:\n", len(extensions))
	for _, ext := range extensions {
		fmt.Printf("- %s: %s (%d stars)\n", ext.Name, ext.Description, ext.Stars)
	}

	// Example: List installed extensions
	fmt.Println("\nInstalled extensions:")
	installed, err := installer.ListInstalledExtensions(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing installed extensions: %v\n", err)
		os.Exit(1)
	}

	if len(installed) == 0 {
		fmt.Println("No extensions installed")
	} else {
		for _, ext := range installed {
			fmt.Printf("- %s (v%s) - %s\n", ext.Name, ext.Version, ext.Repository)
		}
	}

	// Example: Install an extension (commented out to avoid side effects)
	// fmt.Println("\nInstalling example extension...")
	// if err := installer.InstallExtension(ctx, "example", "tykctl-demo"); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error installing extension: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Extension installed successfully!")

	// Example: Create extension runner
	fmt.Println("\nCreating extension runner...")
	runner := extension.NewRunner(configDir)

	// Example: List available extensions
	fmt.Println("Available extensions:")
	available, err := runner.ListAvailableExtensions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing available extensions: %v\n", err)
	} else {
		if len(available) == 0 {
			fmt.Println("No extensions available")
		} else {
			for _, ext := range available {
				fmt.Printf("- %s\n", ext)
			}
		}
	}

	// Example: Check if a specific extension is available
	extensionName := "example"
	if runner.IsExtensionAvailable(extensionName) {
		fmt.Printf("\nExtension '%s' is available for execution\n", extensionName)

		// Example: Run extension (commented out to avoid side effects)
		// fmt.Printf("Running extension '%s'...\n", extensionName)
		// if err := runner.RunExtension(ctx, extensionName, []string{"--help"}); err != nil {
		// 	fmt.Fprintf(os.Stderr, "Error running extension: %v\n", err)
		// }
	} else {
		fmt.Printf("\nExtension '%s' is not available\n", extensionName)
	}

	fmt.Println("\nExtension management and execution example completed!")
}
