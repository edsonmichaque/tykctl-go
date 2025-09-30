package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ContextStore manages multiple contexts
type ContextStore struct {
	contexts   map[string]*Context
	current    string
	configPath string
	logger     Logger
	mu         sync.RWMutex
}

// ContextOptions provides configuration for the context store
type ContextOptions struct {
	ConfigPath string
	Logger     Logger
}

// NewContextStore creates a new context store
func NewContextStore(opts ContextOptions) (*ContextStore, error) {
	cs := &ContextStore{
		contexts:   make(map[string]*Context),
		configPath: opts.ConfigPath,
		logger:     opts.Logger,
	}

	// Load existing contexts
	if err := cs.loadContexts(); err != nil {
		return nil, fmt.Errorf("failed to load contexts: %w", err)
	}

	return cs, nil
}

// CreateContext creates a new context
func (cs *ContextStore) CreateContext(name string, config map[string]interface{}) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, exists := cs.contexts[name]; exists {
		return fmt.Errorf("context %s already exists", name)
	}

	// Validate name
	if name == "" {
		return fmt.Errorf("context name cannot be empty")
	}

	// Validate config
	if config == nil {
		config = make(map[string]interface{})
	}

	context := &Context{
		Name:        name,
		Description: fmt.Sprintf("Context for %s", name),
		Config:      config,
		Resources: ContextResources{
			Hooks:     make(map[string][]Hook),
			Plugins:   make([]Plugin, 0),
			Templates: make([]Template, 0),
			Cache:     make([]CacheConfig, 0),
		},
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate context
	if err := cs.validateContext(context); err != nil {
		return fmt.Errorf("context validation failed: %w", err)
	}

	cs.contexts[name] = context

	// Save contexts
	if err := cs.saveContexts(); err != nil {
		// Rollback the context creation if save fails
		delete(cs.contexts, name)
		return fmt.Errorf("failed to save contexts: %w", err)
	}

	cs.logger.Info("Context created", "name", name)
	return nil
}

// SwitchContext switches to a different context
func (cs *ContextStore) SwitchContext(name string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if name == "" {
		return fmt.Errorf("context name cannot be empty")
	}

	if _, exists := cs.contexts[name]; !exists {
		return fmt.Errorf("context %s does not exist", name)
	}

	cs.current = name

	// Save the current context
	if err := cs.saveContexts(); err != nil {
		return fmt.Errorf("failed to save current context: %w", err)
	}

	cs.logger.Info("Context switched", "name", name)
	return nil
}

// DeleteContext deletes a context
func (cs *ContextStore) DeleteContext(name string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if name == "" {
		return fmt.Errorf("context name cannot be empty")
	}

	if _, exists := cs.contexts[name]; !exists {
		return fmt.Errorf("context %s does not exist", name)
	}

	// Prevent deleting the last context
	if len(cs.contexts) == 1 {
		return fmt.Errorf("cannot delete the last remaining context")
	}

	delete(cs.contexts, name)

	// If we deleted the current context, switch to the first available one
	if cs.current == name {
		for contextName := range cs.contexts {
			cs.current = contextName
			break
		}
	}

	// Save contexts
	if err := cs.saveContexts(); err != nil {
		// Rollback the deletion if save fails
		cs.contexts[name] = &Context{Name: name} // Restore with minimal context
		return fmt.Errorf("failed to save contexts: %w", err)
	}

	cs.logger.Info("Context deleted", "name", name)
	return nil
}

// ListContexts returns a list of all contexts
func (cs *ContextStore) ListContexts() []Context {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	contexts := make([]Context, 0, len(cs.contexts))
	for _, context := range cs.contexts {
		contexts = append(contexts, *context)
	}

	return contexts
}

// GetCurrentContext returns the current context
func (cs *ContextStore) GetCurrentContext() *Context {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if cs.current == "" {
		return nil
	}

	if context, exists := cs.contexts[cs.current]; exists {
		return context
	}

	return nil
}

// SetContextConfig sets a configuration value for a context
func (cs *ContextStore) SetContextConfig(name string, key string, value interface{}) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if name == "" {
		return fmt.Errorf("context name cannot be empty")
	}

	if key == "" {
		return fmt.Errorf("config key cannot be empty")
	}

	context, exists := cs.contexts[name]
	if !exists {
		return fmt.Errorf("context %s does not exist", name)
	}

	// Store old value for potential rollback
	oldValue := context.Config[key]
	context.Config[key] = value
	context.UpdatedAt = time.Now()

	// Save contexts
	if err := cs.saveContexts(); err != nil {
		// Rollback the config change if save fails
		context.Config[key] = oldValue
		return fmt.Errorf("failed to save contexts: %w", err)
	}

	cs.logger.Info("Context config updated", "name", name, "key", key)
	return nil
}

// GetContextConfig gets a configuration value from a context
func (cs *ContextStore) GetContextConfig(name string, key string) (interface{}, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if name == "" {
		return nil, fmt.Errorf("context name cannot be empty")
	}

	if key == "" {
		return nil, fmt.Errorf("config key cannot be empty")
	}

	context, exists := cs.contexts[name]
	if !exists {
		return nil, fmt.Errorf("context %s does not exist", name)
	}

	value, exists := context.Config[key]
	if !exists {
		return nil, fmt.Errorf("config key %s does not exist in context %s", key, name)
	}

	return value, nil
}

// MigrateContext migrates a context from old format
func (cs *ContextStore) MigrateContext(oldContext map[string]interface{}) (*Context, error) {
	if oldContext == nil {
		return nil, fmt.Errorf("old context cannot be nil")
	}

	// Validate required fields
	name, ok := oldContext["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("invalid or missing context name")
	}

	description, ok := oldContext["description"].(string)
	if !ok {
		description = fmt.Sprintf("Migrated context for %s", name)
	}

	config, ok := oldContext["config"].(map[string]interface{})
	if !ok {
		config = make(map[string]interface{})
	}

	newContext := &Context{
		Name:        name,
		Description: description,
		Config:      config,
		Resources: ContextResources{
			Hooks:     make(map[string][]Hook),
			Plugins:   make([]Plugin, 0),
			Templates: make([]Template, 0),
			Cache:     make([]CacheConfig, 0),
		},
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Migrate metadata
	if oldMetadata, ok := oldContext["metadata"].(map[string]interface{}); ok {
		for k, v := range oldMetadata {
			if str, ok := v.(string); ok {
				newContext.Metadata[k] = str
			}
		}
	}

	// Validate the migrated context
	if err := cs.validateContext(newContext); err != nil {
		return nil, fmt.Errorf("migrated context validation failed: %w", err)
	}

	return newContext, nil
}

// Helper functions
func (cs *ContextStore) loadContexts() error {
	contextsFile := filepath.Join(cs.configPath, "contexts.json")

	if _, err := os.Stat(contextsFile); os.IsNotExist(err) {
		// Create contexts directory if it doesn't exist
		if err := os.MkdirAll(cs.configPath, 0755); err != nil {
			return fmt.Errorf("failed to create contexts directory: %w", err)
		}
		return nil
	}

	data, err := os.ReadFile(contextsFile)
	if err != nil {
		return fmt.Errorf("failed to read contexts file: %w", err)
	}

	var contextsData struct {
		Current  string              `json:"current"`
		Contexts map[string]*Context `json:"contexts"`
	}

	if err := json.Unmarshal(data, &contextsData); err != nil {
		return fmt.Errorf("failed to unmarshal contexts: %w", err)
	}

	cs.current = contextsData.Current
	cs.contexts = contextsData.Contexts

	return nil
}

func (cs *ContextStore) saveContexts() error {
	contextsFile := filepath.Join(cs.configPath, "contexts.json")

	contextsData := struct {
		Current  string              `json:"current"`
		Contexts map[string]*Context `json:"contexts"`
	}{
		Current:  cs.current,
		Contexts: cs.contexts,
	}

	data, err := json.MarshalIndent(contextsData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal contexts: %w", err)
	}

	if err := os.WriteFile(contextsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write contexts file: %w", err)
	}

	return nil
}

func (cs *ContextStore) validateContext(context *Context) error {
	// Validate required fields
	if context.Name == "" {
		return fmt.Errorf("context name is required")
	}

	// Validate configuration
	if context.Config == nil {
		context.Config = make(map[string]interface{})
	}

	// Validate resources
	if context.Resources.Hooks == nil {
		context.Resources.Hooks = make(map[string][]Hook)
	}
	if context.Resources.Plugins == nil {
		context.Resources.Plugins = make([]Plugin, 0)
	}
	if context.Resources.Templates == nil {
		context.Resources.Templates = make([]Template, 0)
	}
	if context.Resources.Cache == nil {
		context.Resources.Cache = make([]CacheConfig, 0)
	}

	// Validate metadata
	if context.Metadata == nil {
		context.Metadata = make(map[string]string)
	}

	return nil
}
