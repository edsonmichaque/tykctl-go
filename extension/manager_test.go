package extension

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	configDir := filepath.Join(os.TempDir(), "tykctl-test")
	defer os.RemoveAll(configDir)

	manager := NewManager(configDir)
	if manager == nil {
		t.Fatal("Manager should not be nil")
	}

	if manager.configDir != configDir {
		t.Fatalf("Expected configDir %s, got %s", configDir, manager.configDir)
	}

	if manager.fs == nil {
		t.Fatal("FileSystem should not be nil")
	}
}

func TestManager_InstallExtension(t *testing.T) {
	configDir := filepath.Join(os.TempDir(), "tykctl-test")
	defer os.RemoveAll(configDir)

	manager := NewManager(configDir)
	ctx := context.Background()

	err := manager.InstallExtension(ctx, "test", "extension")
	if err != nil {
		t.Fatalf("InstallExtension failed: %v", err)
	}

	// Check if extension was installed
	extensions, err := manager.ListInstalledExtensions(ctx)
	if err != nil {
		t.Fatalf("ListInstalledExtensions failed: %v", err)
	}

	if len(extensions) != 1 {
		t.Fatalf("Expected 1 extension, got %d", len(extensions))
	}

	ext := extensions[0]
	if ext.Name != "extension" {
		t.Fatalf("Expected name 'extension', got '%s'", ext.Name)
	}

	if ext.Version != "1.0.0" {
		t.Fatalf("Expected version '1.0.0', got '%s'", ext.Version)
	}
}

func TestManager_RemoveExtension(t *testing.T) {
	configDir := filepath.Join(os.TempDir(), "tykctl-test")
	defer os.RemoveAll(configDir)

	manager := NewManager(configDir)
	ctx := context.Background()

	// Install an extension first
	err := manager.InstallExtension(ctx, "test", "extension")
	if err != nil {
		t.Fatalf("InstallExtension failed: %v", err)
	}

	// Remove the extension
	err = manager.RemoveExtension(ctx, "extension")
	if err != nil {
		t.Fatalf("RemoveExtension failed: %v", err)
	}

	// Check if extension was removed
	extensions, err := manager.ListInstalledExtensions(ctx)
	if err != nil {
		t.Fatalf("ListInstalledExtensions failed: %v", err)
	}

	if len(extensions) != 0 {
		t.Fatalf("Expected 0 extensions, got %d", len(extensions))
	}
}

func TestManager_ListInstalledExtensions(t *testing.T) {
	configDir := filepath.Join(os.TempDir(), "tykctl-test")
	defer os.RemoveAll(configDir)

	manager := NewManager(configDir)
	ctx := context.Background()

	// Initially should be empty
	extensions, err := manager.ListInstalledExtensions(ctx)
	if err != nil {
		t.Fatalf("ListInstalledExtensions failed: %v", err)
	}

	if len(extensions) != 0 {
		t.Fatalf("Expected 0 extensions initially, got %d", len(extensions))
	}

	// Install an extension
	err = manager.InstallExtension(ctx, "test", "extension")
	if err != nil {
		t.Fatalf("InstallExtension failed: %v", err)
	}

	// Should now have 1 extension
	extensions, err = manager.ListInstalledExtensions(ctx)
	if err != nil {
		t.Fatalf("ListInstalledExtensions failed: %v", err)
	}

	if len(extensions) != 1 {
		t.Fatalf("Expected 1 extension, got %d", len(extensions))
	}
}

func TestManager_SearchExtensions(t *testing.T) {
	configDir := filepath.Join(os.TempDir(), "tykctl-test")
	defer os.RemoveAll(configDir)

	manager := NewManager(configDir)

	// This will likely return empty results since we're not mocking GitHub API
	extensions, err := manager.SearchExtensions("test", 10)
	if err != nil {
		t.Fatalf("SearchExtensions failed: %v", err)
	}

	// Should not panic and return empty results gracefully
	if extensions == nil {
		t.Fatal("Extensions should not be nil")
	}
}

func TestInstalled(t *testing.T) {
	ext := Installed{
		Name:        "test-extension",
		Version:     "1.0.0",
		Repository:  "https://github.com/test/test-extension",
		InstalledAt: time.Now(),
		Path:        "/path/to/extension",
	}

	if ext.Name != "test-extension" {
		t.Fatal("Name should be set correctly")
	}

	if ext.Version != "1.0.0" {
		t.Fatal("Version should be set correctly")
	}

	if ext.Repository != "https://github.com/test/test-extension" {
		t.Fatal("Repository should be set correctly")
	}

	if ext.Path != "/path/to/extension" {
		t.Fatal("Path should be set correctly")
	}
}

func TestInfo(t *testing.T) {
	info := Info{
		Name:        "test-extension",
		Description: "A test extension",
		Stars:       42,
		UpdatedAt:   time.Now(),
	}

	if info.Name != "test-extension" {
		t.Fatal("Name should be set correctly")
	}

	if info.Description != "A test extension" {
		t.Fatal("Description should be set correctly")
	}

	if info.Stars != 42 {
		t.Fatal("Stars should be set correctly")
	}

	if info.UpdatedAt.IsZero() {
		t.Fatal("UpdatedAt should be set")
	}
}
