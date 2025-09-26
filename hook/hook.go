package hook

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Hook represents a function that can be executed at specific points
type Hook func(ctx context.Context, data interface{}) error

// HookType represents the type of hook
type HookType string

const (
	// HookTypeBeforeInstall is executed before installing an extension
	HookTypeBeforeInstall HookType = "before_install"
	// HookTypeAfterInstall is executed after installing an extension
	HookTypeAfterInstall HookType = "after_install"
	// HookTypeBeforeUninstall is executed before uninstalling an extension
	HookTypeBeforeUninstall HookType = "before_uninstall"
	// HookTypeAfterUninstall is executed after uninstalling an extension
	HookTypeAfterUninstall HookType = "after_uninstall"
	// HookTypeBeforeRun is executed before running an extension
	HookTypeBeforeRun HookType = "before_run"
	// HookTypeAfterRun is executed after running an extension
	HookTypeAfterRun HookType = "after_run"
	// HookTypeOnError is executed when an error occurs
	HookTypeOnError HookType = "on_error"
)

// HookData contains data passed to hooks
type HookData struct {
	ExtensionName string
	ExtensionPath string
	Error         error
	Metadata      map[string]interface{}
}

// Manager manages builtin, external, and Rego hooks for the tykctl-go SDK
type Manager struct {
	builtin  *BuiltinManager
	external *ExternalManager
	rego     *RegoHookManager
}

// New creates a new hook manager
func New() *Manager {
	return &Manager{
		builtin:  NewBuiltinManager(),
		external: NewExternalManager(GetDefaultHookDir()),
		rego:     NewRegoHookManager(nil),
	}
}

// NewWithExternalDir creates a new hook manager with custom external hook directory
func NewWithExternalDir(externalHookDir string) *Manager {
	return &Manager{
		builtin:  NewBuiltinManager(),
		external: NewExternalManager(externalHookDir),
		rego:     NewRegoHookManager(nil),
	}
}

// NewWithLogger creates a new hook manager with custom logger
func NewWithLogger(externalHookDir string, logger *zap.Logger) *Manager {
	return &Manager{
		builtin:  NewBuiltinManager(),
		external: NewExternalManagerWithLogger(externalHookDir, logger),
		rego:     NewRegoHookManager(logger),
	}
}

// RegisterBuiltin registers a builtin hook for a specific type
func (m *Manager) RegisterBuiltin(ctx context.Context, hookType HookType, hook BuiltinHook) {
	m.builtin.Register(ctx, hookType, hook)
}

// UnregisterBuiltin removes a builtin hook for a specific type
func (m *Manager) UnregisterBuiltin(ctx context.Context, hookType HookType, hook BuiltinHook) {
	m.builtin.Unregister(ctx, hookType, hook)
}

// Register registers a builtin hook (for backward compatibility)
func (m *Manager) Register(ctx context.Context, hookType HookType, hook Hook) {
	// Convert old Hook type to BuiltinHook
	builtinHook := func(ctx context.Context, data *HookData) error {
		return hook(ctx, data)
	}
	m.builtin.Register(ctx, hookType, builtinHook)
}

// Unregister removes a builtin hook (for backward compatibility)
func (m *Manager) Unregister(ctx context.Context, hookType HookType, hook Hook) {
	// Convert old Hook type to BuiltinHook
	builtinHook := func(ctx context.Context, data *HookData) error {
		return hook(ctx, data)
	}
	m.builtin.Unregister(ctx, hookType, builtinHook)
}

// Execute executes all hooks (both builtin and external) for a specific type
func (m *Manager) Execute(ctx context.Context, hookType HookType, data interface{}) error {
	hookData, ok := data.(*HookData)
	if !ok {
		return fmt.Errorf("invalid hook data type, expected *HookData")
	}

	// Execute builtin hooks first
	if err := m.builtin.Execute(ctx, hookType, hookData); err != nil {
		return fmt.Errorf("builtin hooks failed: %w", err)
	}

	// Execute external hooks
	if err := m.external.ExecuteExternalHooksForType(ctx, hookType, hookData); err != nil {
		return fmt.Errorf("external hooks failed: %w", err)
	}

	return nil
}

// ExecuteAsync executes all hooks for a specific type asynchronously
func (m *Manager) ExecuteAsync(ctx context.Context, hookType HookType, data interface{}) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		hookData, ok := data.(*HookData)
		if !ok {
			errChan <- fmt.Errorf("invalid hook data type, expected *HookData")
			return
		}

		// Execute builtin hooks asynchronously
		builtinErrChan := m.builtin.ExecuteAsync(ctx, hookType, hookData)

		// Execute external hooks synchronously (they're already file-based)
		externalErr := m.external.ExecuteExternalHooksForType(ctx, hookType, hookData)

		// Wait for builtin hooks to complete
		var builtinErr error
		select {
		case builtinErr = <-builtinErrChan:
		case <-ctx.Done():
			builtinErr = ctx.Err()
		}

		// Combine errors
		var errors []error
		if builtinErr != nil {
			errors = append(errors, builtinErr)
		}
		if externalErr != nil {
			errors = append(errors, externalErr)
		}

		if len(errors) > 0 {
			errChan <- fmt.Errorf("multiple hook errors: %v", errors)
		}
	}()

	return errChan
}

// ListBuiltin returns all registered builtin hooks for a specific type
func (m *Manager) ListBuiltin(ctx context.Context, hookType HookType) []BuiltinHook {
	return m.builtin.List(ctx, hookType)
}

// ListExternal returns all external hooks
func (m *Manager) ListExternal(ctx context.Context) ([]*ExternalHook, error) {
	return m.external.ListExternalHooks(ctx)
}

// List returns all registered builtin hooks for a specific type (for backward compatibility)
func (m *Manager) List(ctx context.Context, hookType HookType) []Hook {
	builtinHooks := m.builtin.List(ctx, hookType)
	hooks := make([]Hook, len(builtinHooks))
	for i, bh := range builtinHooks {
		hooks[i] = func(ctx context.Context, data interface{}) error {
			hookData, ok := data.(*HookData)
			if !ok {
				return fmt.Errorf("invalid hook data type")
			}
			return bh(ctx, hookData)
		}
	}
	return hooks
}

// ClearBuiltin removes all builtin hooks for a specific type
func (m *Manager) ClearBuiltin(ctx context.Context, hookType HookType) {
	m.builtin.Clear(ctx, hookType)
}

// ClearAllBuiltin removes all builtin hooks
func (m *Manager) ClearAllBuiltin(ctx context.Context) {
	m.builtin.ClearAll(ctx)
}

// Clear removes all builtin hooks for a specific type (for backward compatibility)
func (m *Manager) Clear(ctx context.Context, hookType HookType) {
	m.builtin.Clear(ctx, hookType)
}

// ClearAll removes all builtin hooks (for backward compatibility)
func (m *Manager) ClearAll(ctx context.Context) {
	m.builtin.ClearAll(ctx)
}

// CountBuiltin returns the number of builtin hooks for a specific type
func (m *Manager) CountBuiltin(ctx context.Context, hookType HookType) int {
	return m.builtin.Count(ctx, hookType)
}

// CountExternal returns the number of external hooks
func (m *Manager) CountExternal(ctx context.Context) (int, error) {
	hooks, err := m.external.ListExternalHooks(ctx)
	if err != nil {
		return 0, err
	}
	return len(hooks), nil
}

// Count returns the number of builtin hooks for a specific type (for backward compatibility)
func (m *Manager) Count(ctx context.Context, hookType HookType) int {
	return m.builtin.Count(ctx, hookType)
}

// HasBuiltinHooks returns true if there are builtin hooks for a specific type
func (m *Manager) HasBuiltinHooks(ctx context.Context, hookType HookType) bool {
	return m.builtin.HasHooks(ctx, hookType)
}

// HasExternalHooks returns true if there are external hooks
func (m *Manager) HasExternalHooks(ctx context.Context) (bool, error) {
	hooks, err := m.external.ListExternalHooks(ctx)
	if err != nil {
		return false, err
	}
	return len(hooks) > 0, nil
}

// HasHooks returns true if there are builtin hooks for a specific type (for backward compatibility)
func (m *Manager) HasHooks(ctx context.Context, hookType HookType) bool {
	return m.builtin.HasHooks(ctx, hookType)
}

// HookTypes returns all registered builtin hook types
func (m *Manager) HookTypes(ctx context.Context) []HookType {
	return m.builtin.HookTypes(ctx)
}

// External Hook Management Methods

// CreateExternalHook creates a new external hook
func (m *Manager) CreateExternalHook(ctx context.Context, name, content string) (*ExternalHook, error) {
	return m.external.CreateExternalHook(ctx, name, content)
}

// DeleteExternalHook deletes an external hook
func (m *Manager) DeleteExternalHook(ctx context.Context, name string) error {
	return m.external.DeleteExternalHook(ctx, name)
}

// EnableExternalHook enables an external hook
func (m *Manager) EnableExternalHook(ctx context.Context, name string) error {
	return m.external.EnableExternalHook(ctx, name)
}

// DisableExternalHook disables an external hook
func (m *Manager) DisableExternalHook(ctx context.Context, name string) error {
	return m.external.DisableExternalHook(ctx, name)
}

// GetExternalHook returns a specific external hook
func (m *Manager) GetExternalHook(ctx context.Context, name string) (*ExternalHook, error) {
	return m.external.GetExternalHook(ctx, name)
}

// ValidateExternalHook validates an external hook
func (m *Manager) ValidateExternalHook(ctx context.Context, hook *ExternalHook) error {
	return m.external.ValidateExternalHook(hook)
}

// Rego Hook Management Methods

// RegisterRegoHook registers a Rego hook
func (m *Manager) RegisterRegoHook(ctx context.Context, hook *RegoHook) error {
	return m.rego.RegisterRegoHook(ctx, hook)
}

// UnregisterRegoHook removes a Rego hook
func (m *Manager) UnregisterRegoHook(ctx context.Context, name string) error {
	return m.rego.UnregisterRegoHook(ctx, name)
}

// ExecuteRegoHook executes a Rego hook
func (m *Manager) ExecuteRegoHook(ctx context.Context, name string, input map[string]interface{}) (*RegoResult, error) {
	return m.rego.ExecuteRegoHook(ctx, name, input)
}

// ListRegoHooks returns all registered Rego hooks
func (m *Manager) ListRegoHooks(ctx context.Context) []*RegoHook {
	return m.rego.ListRegoHooks(ctx)
}

// GetRegoHook returns a specific Rego hook
func (m *Manager) GetRegoHook(ctx context.Context, name string) (*RegoHook, error) {
	return m.rego.GetRegoHook(ctx, name)
}

// EnableRegoHook enables a Rego hook
func (m *Manager) EnableRegoHook(ctx context.Context, name string) error {
	return m.rego.EnableRegoHook(ctx, name)
}

// DisableRegoHook disables a Rego hook
func (m *Manager) DisableRegoHook(ctx context.Context, name string) error {
	return m.rego.DisableRegoHook(ctx, name)
}

// LoadRegoHooksFromDirectory loads Rego hooks from a directory
func (m *Manager) LoadRegoHooksFromDirectory(ctx context.Context, dir string) error {
	return m.rego.LoadRegoHooksFromDirectory(ctx, dir)
}
