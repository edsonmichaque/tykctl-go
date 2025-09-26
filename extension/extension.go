package extension

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/edsonmichaque/tykctl-go/hook"
	"github.com/google/go-github/v75/github"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v3"
)

// Extension-specific hook types
const (
	HookTypeBeforeInstall   hook.HookType = "extension-before-install"
	HookTypeAfterInstall    hook.HookType = "extension-after-install"
	HookTypeBeforeUninstall hook.HookType = "extension-before-uninstall"
	HookTypeAfterUninstall  hook.HookType = "extension-after-uninstall"
	HookTypeBeforeRun       hook.HookType = "extension-before-run"
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

// Installer manages tykctl extensions
type Installer struct {
	configDir string
	client    *github.Client
	logger    *zap.Logger
	hooks     *hook.Manager
}

// InstallerOption defines a functional option for configuring an Installer
type InstallerOption func(*Installer)

// WithGitHubToken sets the GitHub authentication token
func WithGitHubToken(token string) InstallerOption {
	return func(i *Installer) {
		i.client = github.NewClient(nil).WithAuthToken(token)
	}
}

// WithLogger sets a custom logger
func WithLogger(logger *zap.Logger) InstallerOption {
	return func(i *Installer) {
		i.logger = logger
	}
}

// WithHooks sets a custom hook manager
func WithHooks(hooks *hook.Manager) InstallerOption {
	return func(i *Installer) {
		i.hooks = hooks
	}
}

// NewInstaller creates a new extension installer with the given config directory and options
func NewInstaller(configDir string, opts ...InstallerOption) *Installer {
	// Create default GitHub client
	client := github.NewClient(nil)
	logger := zap.L()

	// Create default hook manager
	hooks := hook.New()

	installer := &Installer{
		configDir: configDir,
		client:    client,
		logger:    logger,
		hooks:     hooks,
	}

	// Apply options
	for _, opt := range opts {
		opt(installer)
	}

	return installer
}

// SearchExtensions searches for extensions
func (i *Installer) SearchExtensions(ctx context.Context, query string, limit int) ([]Info, error) {

	// Search for repositories with tykctl-extension topic
	searchQuery := "topic:tykctl-extension"
	if query != "" {
		searchQuery += " " + query
	}

	opts := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	repos, _, err := i.client.Search.Repositories(ctx, searchQuery, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	var extensions []Info
	for _, repo := range repos.Repositories {
		ext := Info{
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			Stars:       repo.GetStargazersCount(),
			UpdatedAt:   repo.GetUpdatedAt().Time,
		}
		extensions = append(extensions, ext)
	}

	return extensions, nil
}

// InstallExtension installs an extension
func (i *Installer) InstallExtension(ctx context.Context, owner, repo string) error {
	i.logger.Info("Installing extension",
		zap.String("owner", owner),
		zap.String("repo", repo))

	// Execute before install hooks
	hookData := &hook.HookData{
		ExtensionName: repo,
		ExtensionPath: "",
		Metadata: map[string]interface{}{
			"owner": owner,
			"repo":  repo,
		},
	}

	if err := i.hooks.Execute(ctx, HookTypeBeforeInstall, hookData); err != nil {
		i.logger.Error("Before install hook failed", zap.Error(err))
		return fmt.Errorf("before install hook failed: %w", err)
	}

	extensionsDir := filepath.Join(xdg.DataHome, "tykctl", "extensions")
	extDir := filepath.Join(extensionsDir, fmt.Sprintf("tykctl-%s", repo))

	// Create extension directory
	if err := os.MkdirAll(extDir, 0755); err != nil {
		i.logger.Error("Failed to create extension directory",
			zap.String("path", extDir),
			zap.Error(err))
		return fmt.Errorf("failed to create extension directory: %w", err)
	}

	binaryPath := filepath.Join(extDir, fmt.Sprintf("tykctl-%s", repo))
	content := fmt.Sprintf("#!/bin/bash\necho \"Extension %s/%s is not yet implemented\"\n", owner, repo)

	// Create extension binary
	if err := os.WriteFile(binaryPath, []byte(content), 0755); err != nil {
		i.logger.Error("Failed to create extension binary",
			zap.String("path", binaryPath),
			zap.Error(err))
		return fmt.Errorf("failed to create extension binary: %w", err)
	}

	ext := Installed{
		Name:        repo,
		Version:     "1.0.0",
		Repository:  fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		InstalledAt: time.Now(),
		Path:        binaryPath,
	}

	if err := i.saveExtension(ctx, &ext); err != nil {
		i.logger.Error("Failed to save extension metadata",
			zap.String("name", repo),
			zap.Error(err))
		return err
	}

	i.logger.Info("Extension installed successfully",
		zap.String("name", repo),
		zap.String("path", binaryPath))

	// Execute after install hooks
	hookData.ExtensionPath = binaryPath
	if err := i.hooks.Execute(ctx, HookTypeAfterInstall, hookData); err != nil {
		i.logger.Error("After install hook failed", zap.Error(err))
		// Don't fail the installation if after hooks fail
	}

	return nil
}

// RemoveExtension removes an extension
func (i *Installer) RemoveExtension(ctx context.Context, name string) error {
	i.logger.Info("Removing extension", zap.String("name", name))

	extensions, err := i.loadExtensions(ctx)
	if err != nil {
		return err
	}

	ext, exists := extensions[name]
	if !exists {
		return fmt.Errorf("extension %s not found", name)
	}

	// Execute before uninstall hooks
	hookData := &hook.HookData{
		ExtensionName: name,
		ExtensionPath: ext.Path,
		Metadata: map[string]interface{}{
			"version":    ext.Version,
			"repository": ext.Repository,
		},
	}

	if err := i.hooks.Execute(ctx, HookTypeBeforeUninstall, hookData); err != nil {
		i.logger.Error("Before uninstall hook failed", zap.Error(err))
		return fmt.Errorf("before uninstall hook failed: %w", err)
	}

	// Remove extension binary
	if err := os.Remove(ext.Path); err != nil {
		i.logger.Warn("Failed to remove extension binary",
			zap.String("path", ext.Path),
			zap.Error(err))
	}

	// Remove from registry
	delete(extensions, name)
	if err := i.saveExtensions(ctx, extensions); err != nil {
		return err
	}

	i.logger.Info("Extension removed successfully", zap.String("name", name))

	// Execute after uninstall hooks
	if err := i.hooks.Execute(ctx, HookTypeAfterUninstall, hookData); err != nil {
		i.logger.Error("After uninstall hook failed", zap.Error(err))
		// Don't fail the removal if after hooks fail
	}

	return nil
}

// ListInstalledExtensions lists all installed extensions
func (i *Installer) ListInstalledExtensions(ctx context.Context) ([]Installed, error) {
	extensions, err := i.loadExtensions(ctx)
	if err != nil {
		return nil, err
	}

	var installed []Installed
	for _, ext := range extensions {
		installed = append(installed, ext)
	}

	return installed, nil
}

// loadExtensions loads the extensions registry
func (i *Installer) loadExtensions(ctx context.Context) (map[string]Installed, error) {
	extFile := filepath.Join(i.configDir, "extensions.yaml")

	if _, err := os.Stat(extFile); os.IsNotExist(err) {
		return make(map[string]Installed), nil
	}

	data, err := os.ReadFile(extFile)
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
func (i *Installer) saveExtension(ctx context.Context, ext *Installed) error {
	extensions, err := i.loadExtensions(ctx)
	if err != nil {
		return err
	}

	extensions[ext.Name] = *ext
	return i.saveExtensions(ctx, extensions)
}

// saveExtensions saves the extensions registry
func (i *Installer) saveExtensions(ctx context.Context, extensions map[string]Installed) error {
	extFile := filepath.Join(i.configDir, "extensions.yaml")

	data, err := yaml.Marshal(extensions)
	if err != nil {
		return fmt.Errorf("failed to marshal extensions: %w", err)
	}

	if err := os.WriteFile(extFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write extensions file: %w", err)
	}

	return nil
}
