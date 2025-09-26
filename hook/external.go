package hook

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"go.uber.org/zap"
)

// ExternalHook represents a file-based hook similar to Git hooks
type ExternalHook struct {
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Enabled     bool              `json:"enabled"`
	Timeout     time.Duration     `json:"timeout"`
	Environment map[string]string `json:"environment"`
}

// ExternalManager manages external file-based hooks
type ExternalManager struct {
	hookDir string
	logger  *zap.Logger
}

// GetDefaultHookDir returns the default hook directory using XDG Base Directory
func GetDefaultHookDir() string {
	return filepath.Join(xdg.ConfigHome, "tykctl", "hooks")
}

// NewExternalManager creates a new external hook manager
func NewExternalManager(hookDir string) *ExternalManager {
	return &ExternalManager{
		hookDir: hookDir,
		logger:  zap.NewNop(),
	}
}

// NewExternalManagerWithLogger creates a new external hook manager with logger
func NewExternalManagerWithLogger(hookDir string, logger *zap.Logger) *ExternalManager {
	return &ExternalManager{
		hookDir: hookDir,
		logger:  logger,
	}
}

// ListExternalHooks returns all available external hooks
func (em *ExternalManager) ListExternalHooks(ctx context.Context) ([]*ExternalHook, error) {
	hooks := []*ExternalHook{}

	// Check if hook directory exists
	if _, err := os.Stat(em.hookDir); os.IsNotExist(err) {
		em.logger.Debug("Hook directory does not exist", zap.String("dir", em.hookDir))
		return hooks, nil
	}

	// Read hook directory
	entries, err := os.ReadDir(em.hookDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read hook directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		hookName := entry.Name()
		if strings.HasPrefix(hookName, ".") {
			continue
		}

		hook := &ExternalHook{
			Name:        hookName,
			Path:        filepath.Join(em.hookDir, hookName),
			Enabled:     true,
			Timeout:     30 * time.Second,
			Environment: make(map[string]string),
		}

		hooks = append(hooks, hook)
	}

	return hooks, nil
}

// GetExternalHook returns a specific external hook by name
func (em *ExternalManager) GetExternalHook(ctx context.Context, name string) (*ExternalHook, error) {
	hooks, err := em.ListExternalHooks(ctx)
	if err != nil {
		return nil, err
	}

	for _, hook := range hooks {
		if hook.Name == name {
			return hook, nil
		}
	}

	return nil, fmt.Errorf("external hook not found: %s", name)
}

// ExecuteExternalHook executes a single external hook
func (em *ExternalManager) ExecuteExternalHook(ctx context.Context, hook *ExternalHook, data *HookData) error {
	if !hook.Enabled {
		em.logger.Debug("External hook disabled, skipping", zap.String("hook", hook.Name))
		return nil
	}

	em.logger.Info("Executing external hook", zap.String("hook", hook.Name))

	// Create execution context with timeout
	execCtx, cancel := context.WithTimeout(ctx, hook.Timeout)
	defer cancel()

	// Prepare environment
	env := os.Environ()
	for k, v := range hook.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Add hook context to environment
	env = append(env, fmt.Sprintf("TYKCTL_HOOK_EVENT=%s", data.ExtensionName))
	env = append(env, fmt.Sprintf("TYKCTL_HOOK_EXTENSION=%s", data.ExtensionName))
	env = append(env, fmt.Sprintf("TYKCTL_HOOK_PATH=%s", data.ExtensionPath))
	env = append(env, fmt.Sprintf("TYKCTL_HOOK_WORKING_DIR=%s", os.Getenv("PWD")))

	// Add metadata as environment variables
	if data.Metadata != nil {
		for k, v := range data.Metadata {
			env = append(env, fmt.Sprintf("TYKCTL_HOOK_%s=%v", strings.ToUpper(k), v))
		}
	}

	// Execute the hook script
	cmd := exec.CommandContext(execCtx, hook.Path)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		em.logger.Error("External hook execution failed",
			zap.String("hook", hook.Name),
			zap.Error(err))
		return fmt.Errorf("external hook %s failed: %w", hook.Name, err)
	}

	em.logger.Info("External hook executed successfully", zap.String("hook", hook.Name))
	return nil
}

// ExecuteExternalHooksForType executes all external hooks for a specific hook type
func (em *ExternalManager) ExecuteExternalHooksForType(ctx context.Context, hookType HookType, data *HookData) error {
	hooks, err := em.ListExternalHooks(ctx)
	if err != nil {
		return err
	}

	// Filter hooks by type (hooks can be named with type prefix)
	typePrefix := string(hookType) + "-"

	for _, hook := range hooks {
		if !hook.Enabled {
			continue
		}

		// Check if hook matches the type (either exact match or prefixed)
		if hook.Name == string(hookType) || strings.HasPrefix(hook.Name, typePrefix) {
			if err := em.ExecuteExternalHook(ctx, hook, data); err != nil {
				em.logger.Error("External hook execution failed",
					zap.String("hook", hook.Name),
					zap.Error(err))
				// Continue with other hooks even if one fails
			}
		}
	}

	return nil
}

// CreateExternalHook creates a new external hook file
func (em *ExternalManager) CreateExternalHook(ctx context.Context, name, content string) (*ExternalHook, error) {
	// Ensure hook directory exists
	if err := os.MkdirAll(em.hookDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create hook directory: %w", err)
	}

	hookPath := filepath.Join(em.hookDir, name)
	if err := os.WriteFile(hookPath, []byte(content), 0755); err != nil {
		return nil, fmt.Errorf("failed to create hook file: %w", err)
	}

	hook := &ExternalHook{
		Name:        name,
		Path:        hookPath,
		Enabled:     true,
		Timeout:     30 * time.Second,
		Environment: make(map[string]string),
	}

	em.logger.Info("External hook created", zap.String("hook", name))
	return hook, nil
}

// DeleteExternalHook deletes an external hook
func (em *ExternalManager) DeleteExternalHook(ctx context.Context, name string) error {
	hookPath := filepath.Join(em.hookDir, name)

	if err := os.Remove(hookPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("external hook not found: %s", name)
		}
		return fmt.Errorf("failed to delete external hook: %w", err)
	}

	em.logger.Info("External hook deleted", zap.String("hook", name))
	return nil
}

// EnableExternalHook enables an external hook
func (em *ExternalManager) EnableExternalHook(ctx context.Context, name string) error {
	hook, err := em.GetExternalHook(ctx, name)
	if err != nil {
		return err
	}

	hook.Enabled = true
	em.logger.Info("External hook enabled", zap.String("hook", name))
	return nil
}

// DisableExternalHook disables an external hook
func (em *ExternalManager) DisableExternalHook(ctx context.Context, name string) error {
	hook, err := em.GetExternalHook(ctx, name)
	if err != nil {
		return err
	}

	hook.Enabled = false
	em.logger.Info("External hook disabled", zap.String("hook", name))
	return nil
}

// ValidateExternalHook validates an external hook
func (em *ExternalManager) ValidateExternalHook(hook *ExternalHook) error {
	if hook.Name == "" {
		return fmt.Errorf("hook name cannot be empty")
	}

	if hook.Path == "" {
		return fmt.Errorf("hook path cannot be empty")
	}

	// Check if hook file exists
	if _, err := os.Stat(hook.Path); os.IsNotExist(err) {
		return fmt.Errorf("hook file does not exist: %s", hook.Path)
	}

	// Check if hook is executable
	if info, err := os.Stat(hook.Path); err == nil {
		if info.Mode()&0111 == 0 {
			em.logger.Warn("Hook file is not executable", zap.String("path", hook.Path))
		}
	}

	return nil
}
