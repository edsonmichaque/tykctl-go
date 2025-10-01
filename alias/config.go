package alias

import (
	"context"
	"fmt"
	"sync"
)

// InMemoryConfigProvider provides in-memory alias storage
type InMemoryConfigProvider struct {
	aliases map[string]string
	mutex   sync.RWMutex
}

// NewInMemoryConfigProvider creates a new in-memory config provider
func NewInMemoryConfigProvider() *InMemoryConfigProvider {
	return &InMemoryConfigProvider{
		aliases: make(map[string]string),
	}
}

// SetAlias sets an alias in memory
func (p *InMemoryConfigProvider) SetAlias(ctx context.Context, name, expansion string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.aliases[name] = expansion
	return nil
}

// GetAlias gets an alias from memory
func (p *InMemoryConfigProvider) GetAlias(ctx context.Context, name string) (string, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	expansion, exists := p.aliases[name]
	return expansion, exists
}

// DeleteAlias deletes an alias from memory
func (p *InMemoryConfigProvider) DeleteAlias(ctx context.Context, name string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	delete(p.aliases, name)
	return nil
}

// ListAliases lists all aliases from memory
func (p *InMemoryConfigProvider) ListAliases(ctx context.Context) map[string]string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]string)
	for name, expansion := range p.aliases {
		result[name] = expansion
	}
	return result
}

// ConfigProviderBuilder helps build config providers
type ConfigProviderBuilder struct {
	provider ConfigProvider
}

// NewConfigProviderBuilder creates a new config provider builder
func NewConfigProviderBuilder() *ConfigProviderBuilder {
	return &ConfigProviderBuilder{}
}

// WithInMemory sets an in-memory config provider
func (b *ConfigProviderBuilder) WithInMemory() *ConfigProviderBuilder {
	b.provider = NewInMemoryConfigProvider()
	return b
}

// WithCustom sets a custom config provider
func (b *ConfigProviderBuilder) WithCustom(provider ConfigProvider) *ConfigProviderBuilder {
	b.provider = provider
	return b
}

// Build builds the config provider
func (b *ConfigProviderBuilder) Build() ConfigProvider {
	if b.provider == nil {
		// Default to in-memory provider
		return NewInMemoryConfigProvider()
	}
	return b.provider
}

// ExtensionConfigProvider is a config provider that integrates with extension configuration
type ExtensionConfigProvider struct {
	configManager interface{} // This would be the actual config manager type from the extension
	setAliasFunc  func(ctx context.Context, name, expansion string) error
	getAliasFunc  func(ctx context.Context, name string) (string, bool)
	deleteAliasFunc func(ctx context.Context, name string) error
	listAliasesFunc func(ctx context.Context) map[string]string
}

// NewExtensionConfigProvider creates a new extension config provider
func NewExtensionConfigProvider(
	setAliasFunc func(ctx context.Context, name, expansion string) error,
	getAliasFunc func(ctx context.Context, name string) (string, bool),
	deleteAliasFunc func(ctx context.Context, name string) error,
	listAliasesFunc func(ctx context.Context) map[string]string,
) *ExtensionConfigProvider {
	return &ExtensionConfigProvider{
		setAliasFunc:    setAliasFunc,
		getAliasFunc:    getAliasFunc,
		deleteAliasFunc: deleteAliasFunc,
		listAliasesFunc: listAliasesFunc,
	}
}

// SetAlias sets an alias using the extension's config
func (p *ExtensionConfigProvider) SetAlias(ctx context.Context, name, expansion string) error {
	if p.setAliasFunc == nil {
		return fmt.Errorf("set alias function not configured")
	}
	return p.setAliasFunc(ctx, name, expansion)
}

// GetAlias gets an alias using the extension's config
func (p *ExtensionConfigProvider) GetAlias(ctx context.Context, name string) (string, bool) {
	if p.getAliasFunc == nil {
		return "", false
	}
	return p.getAliasFunc(ctx, name)
}

// DeleteAlias deletes an alias using the extension's config
func (p *ExtensionConfigProvider) DeleteAlias(ctx context.Context, name string) error {
	if p.deleteAliasFunc == nil {
		return fmt.Errorf("delete alias function not configured")
	}
	return p.deleteAliasFunc(ctx, name)
}

// ListAliases lists all aliases using the extension's config
func (p *ExtensionConfigProvider) ListAliases(ctx context.Context) map[string]string {
	if p.listAliasesFunc == nil {
		return make(map[string]string)
	}
	return p.listAliasesFunc(ctx)
}