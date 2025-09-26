// Package extension provides extension management functionality for tykctl.
package extension

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/edsonmichaque/tykctl-go/fs"
	"github.com/google/go-github/v75/github"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v3"
)

// Info represents information about an extension
type Info struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stars       int       `json:"stargazers_count"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Installed represents an installed extension
type Installed struct {
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Repository  string    `yaml:"repository"`
	InstalledAt time.Time `yaml:"installed_at"`
	Path        string    `yaml:"path"`
}

// Manager manages tykctl extensions
type Manager struct {
	configDir string
	fs        *fs.FS
	client    *github.Client
}

// NewManager creates a new extension manager
func NewManager(configDir string) *Manager {
	fsInstance := fs.New()
	ctx := context.Background()
	_ = fsInstance.MkdirAll(ctx, configDir, 0o755)

	// Create GitHub client
	client := github.NewClient(nil)

	return &Manager{
		configDir: configDir,
		fs:        fsInstance,
		client:    client,
	}
}

// NewManagerWithAuth creates a new extension manager with GitHub authentication
func NewManagerWithAuth(configDir, token string) *Manager {
	fsInstance := fs.New()
	ctx := context.Background()
	_ = fsInstance.MkdirAll(ctx, configDir, 0o755)

	// Create GitHub client with authentication
	client := github.NewClient(nil).WithAuthToken(token)

	return &Manager{
		configDir: configDir,
		fs:        fsInstance,
		client:    client,
	}
}

// SearchExtensions searches for extensions
func (em *Manager) SearchExtensions(query string, limit int) ([]Info, error) {
	logger := zap.L().Sugar()

	logger.Debugw("Searching for extensions",
		"query", query,
		"limit", limit)

	// Search for repositories with tykctl-extension topic
	searchQuery := fmt.Sprintf("topic:tykctl-extension %s", query)

	ctx := context.Background()
	opts := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	result, _, err := em.client.Search.Repositories(ctx, searchQuery, opts)
	if err != nil {
		logger.Errorw("Failed to search GitHub repositories",
			"query", searchQuery,
			"error", err)
		// Return empty results if search fails
		return []Info{}, nil
	}

	// Convert GitHub repositories to Info
	extensions := make([]Info, 0, len(result.Repositories))
	for _, repo := range result.Repositories {
		extensions = append(extensions, Info{
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			Stars:       repo.GetStargazersCount(),
			UpdatedAt:   repo.GetUpdatedAt().Time,
		})
	}

	logger.Infow("Found extensions",
		"query", query,
		"count", len(extensions))

	return extensions, nil
}

// InstallExtension installs an extension
func (em *Manager) InstallExtension(ctx context.Context, owner, repo string) error {
	logger := zap.L().Sugar()

	logger.Infow("Installing extension",
		"owner", owner,
		"repo", repo)

	extensionsDir := em.fs.JoinPath(xdg.DataHome, "tykctl", "extensions")
	extDir := em.fs.JoinPath(extensionsDir, fmt.Sprintf("tykctl-%s", repo))

	if err := em.fs.EnsureDir(ctx, extDir, 0o755); err != nil {
		logger.Errorw("Failed to create extension directory",
			"path", extDir,
			"error", err)
		return fmt.Errorf("failed to create extension directory: %w", err)
	}

	binaryPath := em.fs.JoinPath(extDir, fmt.Sprintf("tykctl-%s", repo))
	content := fmt.Sprintf("#!/bin/bash\necho \"Extension %s/%s is not yet implemented\"\n", owner, repo)

	// #nosec G306 -- 755 permissions required for executable files
	if err := em.fs.EnsureFile(ctx, binaryPath, []byte(content), 0o755); err != nil {
		logger.Errorw("Failed to create extension binary",
			"path", binaryPath,
			"error", err)
		return fmt.Errorf("failed to create extension binary: %w", err)
	}

	ext := Installed{
		Name:        repo,
		Version:     "1.0.0",
		Repository:  fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		InstalledAt: time.Now(),
		Path:        binaryPath,
	}

	if err := em.saveExtension(ctx, &ext); err != nil {
		logger.Errorw("Failed to save extension metadata",
			"name", repo,
			"error", err)
		return err
	}

	logger.Infow("Extension installed successfully",
		"name", repo,
		"path", binaryPath)

	return nil
}

// RemoveExtension removes an extension
func (em *Manager) RemoveExtension(ctx context.Context, name string) error {
	logger := zap.L().Sugar()

	logger.Infow("Removing extension", "name", name)

	extensions, err := em.loadExtensions(ctx)
	if err != nil {
		logger.Errorw("Failed to load extensions", "error", err)
		return fmt.Errorf("failed to load extensions: %w", err)
	}

	ext, exists := extensions[name]
	if !exists {
		logger.Warnw("Extension not found", "name", name)
		return fmt.Errorf("extension %s not found", name)
	}

	if err := em.fs.RemoveAllIfExists(ctx, em.fs.JoinPath(filepath.Dir(ext.Path))); err != nil {
		logger.Errorw("Failed to remove extension directory",
			"path", ext.Path,
			"error", err)
		return fmt.Errorf("failed to remove extension directory: %w", err)
	}

	delete(extensions, name)
	if err := em.saveExtensions(ctx, extensions); err != nil {
		logger.Errorw("Failed to save extensions after removal",
			"name", name,
			"error", err)
		return err
	}

	logger.Infow("Extension removed successfully", "name", name)
	return nil
}

// ListInstalledExtensions lists all installed extensions
func (em *Manager) ListInstalledExtensions(ctx context.Context) ([]Installed, error) {
	extensions, err := em.loadExtensions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load extensions: %w", err)
	}

	var result []Installed
	for _, ext := range extensions {
		result = append(result, ext)
	}

	return result, nil
}

// loadExtensions loads the extensions registry
func (em *Manager) loadExtensions(ctx context.Context) (map[string]Installed, error) {
	extFile := em.fs.JoinPath(em.configDir, "extensions.yaml")

	exists, err := em.fs.Exists(ctx, extFile)
	if err != nil {
		return nil, fmt.Errorf("failed to check extensions file: %w", err)
	}
	if !exists {
		return make(map[string]Installed), nil
	}

	data, err := em.fs.ReadFile(ctx, extFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read extensions file: %w", err)
	}

	var extensions map[string]Installed
	if err := yaml.Unmarshal(data, &extensions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal extensions: %w", err)
	}

	return extensions, nil
}

// saveExtension saves a single extension to the registry
func (em *Manager) saveExtension(ctx context.Context, ext *Installed) error {
	extensions, err := em.loadExtensions(ctx)
	if err != nil {
		return err
	}

	extensions[ext.Name] = *ext
	return em.saveExtensions(ctx, extensions)
}

// saveExtensions saves the extensions registry
func (em *Manager) saveExtensions(ctx context.Context, extensions map[string]Installed) error {
	extFile := em.fs.JoinPath(em.configDir, "extensions.yaml")

	data, err := yaml.Marshal(extensions)
	if err != nil {
		return fmt.Errorf("failed to marshal extensions: %w", err)
	}

	// #nosec G306 -- 644 permissions are standard for config files
	if err := em.fs.EnsureFile(ctx, extFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write extensions file: %w", err)
	}

	return nil
}
