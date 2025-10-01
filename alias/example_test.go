package alias

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliasManager(t *testing.T) {
	// Create in-memory provider for testing
	provider := NewInMemoryConfigProvider()
	manager := NewManager(provider, []string{"help", "version"})

	ctx := context.Background()

	// Test setting alias
	err := manager.SetAlias(ctx, "test", "echo hello")
	assert.NoError(t, err)

	// Test getting alias
	expansion, exists := manager.GetAlias(ctx, "test")
	assert.True(t, exists)
	assert.Equal(t, "echo hello", expansion)

	// Test listing aliases
	aliases := manager.ListAliases(ctx)
	assert.Len(t, aliases, 1)
	assert.Equal(t, "echo hello", aliases["test"])

	// Test deleting alias
	err = manager.DeleteAlias(ctx, "test")
	assert.NoError(t, err)

	// Verify deletion
	_, exists = manager.GetAlias(ctx, "test")
	assert.False(t, exists)
}

func TestParameterExpansion(t *testing.T) {
	provider := NewInMemoryConfigProvider()
	manager := NewManager(provider, []string{})

	// Test parameter expansion
	expansion := "echo $1 $2"
	args := []string{"hello", "world"}
	result := manager.ExpandAliasPreview(expansion, args)
	assert.Equal(t, "echo hello world", result)

	// Test $* expansion
	expansion = "echo $*"
	result = manager.ExpandAliasPreview(expansion, args)
	assert.Equal(t, "echo hello world", result)

	// Test $@ expansion
	expansion = "echo $@"
	result = manager.ExpandAliasPreview(expansion, args)
	assert.Equal(t, "echo hello world", result)
}

func TestAliasValidation(t *testing.T) {
	provider := NewInMemoryConfigProvider()
	manager := NewManager(provider, []string{"help", "version"})

	ctx := context.Background()

	// Test valid alias
	err := manager.SetAlias(ctx, "valid", "echo hello")
	assert.NoError(t, err)

	// Test reserved name
	err = manager.SetAlias(ctx, "help", "echo hello")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reserved command name")

	// Test empty name
	err = manager.SetAlias(ctx, "", "echo hello")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	// Test whitespace in name
	err = manager.SetAlias(ctx, "invalid name", "echo hello")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot contain whitespace")

	// Test shell metacharacters
	err = manager.SetAlias(ctx, "invalid&name", "echo hello")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot contain shell metacharacters")
}

func TestAliasTypes(t *testing.T) {
	provider := NewInMemoryConfigProvider()
	manager := NewManager(provider, []string{})

	// Test command alias
	expansion := "echo hello"
	assert.False(t, manager.IsShellAlias(expansion))
	assert.Equal(t, "command", manager.GetAliasType(expansion))

	// Test shell alias
	expansion = "!echo hello"
	assert.True(t, manager.IsShellAlias(expansion))
	assert.Equal(t, "shell", manager.GetAliasType(expansion))
}

func TestConfigProviderBuilder(t *testing.T) {
	// Test in-memory provider
	provider := NewConfigProviderBuilder().
		WithInMemory().
		Build()

	assert.NotNil(t, provider)

	ctx := context.Background()
	err := provider.SetAlias(ctx, "test", "echo hello")
	assert.NoError(t, err)

	expansion, exists := provider.GetAlias(ctx, "test")
	assert.True(t, exists)
	assert.Equal(t, "echo hello", expansion)
}

func TestExtensionConfigProvider(t *testing.T) {
	// Mock functions
	setAliasFunc := func(ctx context.Context, name, expansion string) error {
		return nil
	}

	getAliasFunc := func(ctx context.Context, name string) (string, bool) {
		if name == "test" {
			return "echo hello", true
		}
		return "", false
	}

	deleteAliasFunc := func(ctx context.Context, name string) error {
		return nil
	}

	listAliasesFunc := func(ctx context.Context) map[string]string {
		return map[string]string{"test": "echo hello"}
	}

	// Create extension config provider
	provider := NewExtensionConfigProvider(
		setAliasFunc,
		getAliasFunc,
		deleteAliasFunc,
		listAliasesFunc,
	)

	ctx := context.Background()

	// Test operations
	err := provider.SetAlias(ctx, "test", "echo hello")
	assert.NoError(t, err)

	expansion, exists := provider.GetAlias(ctx, "test")
	assert.True(t, exists)
	assert.Equal(t, "echo hello", expansion)

	aliases := provider.ListAliases(ctx)
	assert.Len(t, aliases, 1)
	assert.Equal(t, "echo hello", aliases["test"])

	err = provider.DeleteAlias(ctx, "test")
	assert.NoError(t, err)
}