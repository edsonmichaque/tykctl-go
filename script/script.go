// Package script provides a generic script-based hook system for extensions.
// Extensions can define their own events and implement custom script handlers.
//
// Example usage in an extension:
//
//	// Define custom events
//	const (
//		EventBeforeCreate ScriptEvent = "before-create"
//		EventAfterCreate  ScriptEvent = "after-create"
//		EventBeforeDelete ScriptEvent = "before-delete"
//	)
//
//	// Create script registry
//	registry := script.NewScriptRegistry()
//
//	// Register handlers
//	registry.RegisterHandler(EventBeforeCreate, func(ctx context.Context, scriptCtx *script.ScriptContext) error {
//		// Custom logic before create
//		return nil
//	})
//
//	// Execute scripts
//	scriptCtx := &script.ScriptContext{
//		Event: EventBeforeCreate,
//		Data:  map[string]interface{}{"resource": "user"},
//	}
//	registry.ExecuteHandlers(ctx, EventBeforeCreate, scriptCtx)
package script

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

// Script represents a script that can be executed
type Script struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Script      string            `json:"script"`
	Enabled     bool              `json:"enabled"`
	Timeout     time.Duration     `json:"timeout"`
	Environment map[string]string `json:"environment"`
	WorkingDir  string            `json:"working_dir"`
}

// ScriptManager manages scripts for the application
type ScriptManager struct {
	scriptDir string
	logger    *zap.Logger
}

// GetDefaultScriptDir returns the default script directory using XDG Base Directory
func GetDefaultScriptDir() string {
	return filepath.Join(xdg.ConfigHome, "tykctl", "scripts")
}

// ScriptValidator provides validation for script configurations
type ScriptValidator struct {
	logger *zap.Logger
}

// NewScriptValidator creates a new script validator
func NewScriptValidator(logger *zap.Logger) *ScriptValidator {
	return &ScriptValidator{logger: logger}
}

// ValidateScript validates a script configuration
func (sv *ScriptValidator) ValidateScript(script *Script) error {
	if script.Name == "" {
		return fmt.Errorf("script name cannot be empty")
	}

	if script.Script == "" {
		return fmt.Errorf("script script cannot be empty")
	}

	// Check if script file exists
	if _, err := os.Stat(script.Script); os.IsNotExist(err) {
		return fmt.Errorf("script script file does not exist: %s", script.Script)
	}

	// Check if script is executable
	if info, err := os.Stat(script.Script); err == nil {
		if info.Mode()&0111 == 0 {
			sv.logger.Warn("Script script is not executable", zap.String("script", script.Script))
		}
	}

	return nil
}

// ScriptExecutor provides execution capabilities for scripts
type ScriptExecutor struct {
	logger *zap.Logger
}

// NewScriptExecutor creates a new script executor
func NewScriptExecutor(logger *zap.Logger) *ScriptExecutor {
	return &ScriptExecutor{logger: logger}
}

// ExecuteScript executes a single script with the given context
func (se *ScriptExecutor) ExecuteScript(ctx context.Context, script *Script, scriptCtx *ScriptContext) error {
	if !script.Enabled {
		se.logger.Debug("Script is disabled, skipping", zap.String("script", script.Name))
		return nil
	}

	se.logger.Info("Executing script", zap.String("script", script.Name), zap.String("event", string(scriptCtx.Event)))

	// Set up environment variables
	env := os.Environ()
	for key, value := range script.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	// Add script context to environment
	env = append(env, fmt.Sprintf("TYKCTL_SCRIPT_EVENT=%s", scriptCtx.Event))
	env = append(env, fmt.Sprintf("TYKCTL_SCRIPT_COMMAND=%s", scriptCtx.Command))
	env = append(env, fmt.Sprintf("TYKCTL_SCRIPT_EXTENSION=%s", scriptCtx.Extension))
	env = append(env, fmt.Sprintf("TYKCTL_SCRIPT_WORKING_DIR=%s", scriptCtx.WorkingDir))

	// Execute the script script
	cmd := exec.CommandContext(ctx, script.Script)
	cmd.Env = env
	cmd.Dir = scriptCtx.WorkingDir

	output, err := cmd.Output()
	if err != nil {
		se.logger.Error("Script execution failed",
			zap.String("script", script.Name),
			zap.Error(err),
			zap.String("output", string(output)))
		return fmt.Errorf("script execution failed: %w", err)
	}

	se.logger.Info("Script executed successfully",
		zap.String("script", script.Name),
		zap.String("output", string(output)))

	return nil
}

// NewScriptManager creates a new script manager
func NewScriptManager(scriptDir string) *ScriptManager {
	return &ScriptManager{
		scriptDir: scriptDir,
		logger:    zap.NewNop(),
	}
}

// NewScriptManagerWithLogger creates a new script manager with custom logger
func NewScriptManagerWithLogger(scriptDir string, logger *zap.Logger) *ScriptManager {
	return &ScriptManager{
		scriptDir: scriptDir,
		logger:    logger,
	}
}

// ScriptEvent represents a custom event that can trigger scripts
// Extensions define their own event types
type ScriptEvent string

// ScriptContext provides context for script execution
type ScriptContext struct {
	Event       ScriptEvent
	Command     string
	Args        []string
	Extension   string
	WorkingDir  string
	Environment map[string]string
	Data        map[string]interface{} // Custom data for extensions
}

// ExecuteScript executes a script with the given context
func (sm *ScriptManager) ExecuteScript(ctx context.Context, script *Script, scriptCtx *ScriptContext) error {
	if !script.Enabled {
		sm.logger.Debug("Script disabled, skipping", zap.String("script", script.Name))
		return nil
	}

	sm.logger.Info("Executing script", zap.String("script", script.Name), zap.String("event", string(scriptCtx.Event)))

	// Create execution context with timeout
	execCtx, cancel := context.WithTimeout(ctx, script.Timeout)
	defer cancel()

	// Prepare environment
	env := make(map[string]string)
	for k, v := range script.Environment {
		env[k] = v
	}
	for k, v := range scriptCtx.Environment {
		env[k] = v
	}

	// Add standard script environment variables
	env["TYKCTL_SCRIPT_EVENT"] = string(scriptCtx.Event)
	env["TYKCTL_SCRIPT_COMMAND"] = scriptCtx.Command
	env["TYKCTL_SCRIPT_ARGS"] = strings.Join(scriptCtx.Args, " ")
	env["TYKCTL_SCRIPT_EXTENSION"] = scriptCtx.Extension
	env["TYKCTL_SCRIPT_WORKING_DIR"] = scriptCtx.WorkingDir

	// Execute the script script
	err := sm.executeScript(execCtx, script, scriptCtx, env)
	if err != nil {
		sm.logger.Error("Script execution failed", zap.String("script", script.Name), zap.Error(err))
		return fmt.Errorf("script %s failed: %w", script.Name, err)
	}

	sm.logger.Info("Script executed successfully", zap.String("script", script.Name))
	return nil
}

// executeScript executes the script script
func (sm *ScriptManager) executeScript(ctx context.Context, script *Script, scriptCtx *ScriptContext, env map[string]string) error {
	sm.logger.Debug("Executing script script",
		zap.String("script", script.Script),
		zap.String("working_dir", script.WorkingDir),
		zap.Duration("timeout", script.Timeout))

	// Convert environment map to slice
	envSlice := os.Environ()
	for k, v := range env {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	// Execute the actual script
	cmd := exec.CommandContext(ctx, script.Script)
	cmd.Env = envSlice
	cmd.Dir = scriptCtx.WorkingDir

	output, err := cmd.Output()
	if err != nil {
		sm.logger.Error("Script execution failed",
			zap.String("script", script.Name),
			zap.Error(err),
			zap.String("output", string(output)))
		return fmt.Errorf("script execution failed: %w", err)
	}

	sm.logger.Debug("Script script completed",
		zap.String("script", script.Name),
		zap.String("output", string(output)))
	return nil
}

// ListScripts returns all available scripts
func (sm *ScriptManager) ListScripts() ([]*Script, error) {
	scripts := []*Script{}

	// Check if script directory exists
	if _, err := os.Stat(sm.scriptDir); os.IsNotExist(err) {
		sm.logger.Debug("Script directory does not exist", zap.String("dir", sm.scriptDir))
		return scripts, nil
	}

	// Read script directory
	entries, err := os.ReadDir(sm.scriptDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read script directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		scriptName := entry.Name()
		if strings.HasPrefix(scriptName, ".") {
			continue
		}

		script := &Script{
			Name:        scriptName,
			Description: fmt.Sprintf("Script: %s", scriptName),
			Script:      filepath.Join(sm.scriptDir, scriptName),
			Enabled:     true,
			Timeout:     30 * time.Second,
			Environment: make(map[string]string),
			WorkingDir:  sm.scriptDir,
		}

		scripts = append(scripts, script)
	}

	return scripts, nil
}

// GetScript returns a specific script by name
func (sm *ScriptManager) GetScript(name string) (*Script, error) {
	scripts, err := sm.ListScripts()
	if err != nil {
		return nil, err
	}

	for _, script := range scripts {
		if script.Name == name {
			return script, nil
		}
	}

	return nil, fmt.Errorf("script not found: %s", name)
}

// EnableScript enables a script
func (sm *ScriptManager) EnableScript(name string) error {
	script, err := sm.GetScript(name)
	if err != nil {
		return err
	}

	script.Enabled = true
	sm.logger.Info("Script enabled", zap.String("script", name))
	return nil
}

// DisableScript disables a script
func (sm *ScriptManager) DisableScript(name string) error {
	script, err := sm.GetScript(name)
	if err != nil {
		return err
	}

	script.Enabled = false
	sm.logger.Info("Script disabled", zap.String("script", name))
	return nil
}

// CreateScript creates a new script
func (sm *ScriptManager) CreateScript(name, description, scriptContent string) (*Script, error) {
	// Ensure script directory exists
	if err := os.MkdirAll(sm.scriptDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create script directory: %w", err)
	}

	script := &Script{
		Name:        name,
		Description: description,
		Script:      filepath.Join(sm.scriptDir, name),
		Enabled:     true,
		Timeout:     30 * time.Second,
		Environment: make(map[string]string),
		WorkingDir:  sm.scriptDir,
	}

	// Create script script file
	scriptPath := filepath.Join(sm.scriptDir, name)
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return nil, fmt.Errorf("failed to create script script: %w", err)
	}

	sm.logger.Info("Script created", zap.String("script", name))
	return script, nil
}

// DeleteScript deletes a script
func (sm *ScriptManager) DeleteScript(name string) error {
	scriptPath := filepath.Join(sm.scriptDir, name)

	if err := os.Remove(scriptPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("script not found: %s", name)
		}
		return fmt.Errorf("failed to delete script: %w", err)
	}

	sm.logger.Info("Script deleted", zap.String("script", name))
	return nil
}

// ExecuteScriptsForEvent executes all scripts for a specific event
// Extensions can implement their own event filtering logic
func (sm *ScriptManager) ExecuteScriptsForEvent(ctx context.Context, event ScriptEvent, scriptCtx *ScriptContext) error {
	scripts, err := sm.ListScripts()
	if err != nil {
		return err
	}

	// Execute all enabled scripts - extensions can implement their own filtering
	for _, script := range scripts {
		if script.Enabled {
			if err := sm.ExecuteScript(ctx, script, scriptCtx); err != nil {
				sm.logger.Error("Script execution failed", zap.String("script", script.Name), zap.Error(err))
				// Continue with other scripts even if one fails
			}
		}
	}

	return nil
}

// ScriptHandler defines a function that can handle script events
type ScriptHandler func(ctx context.Context, scriptCtx *ScriptContext) error

// ScriptRegistry allows extensions to register custom event handlers
type ScriptRegistry struct {
	handlers map[ScriptEvent][]ScriptHandler
}

// NewScriptRegistry creates a new script registry
func NewScriptRegistry() *ScriptRegistry {
	return &ScriptRegistry{
		handlers: make(map[ScriptEvent][]ScriptHandler),
	}
}

// RegisterHandler registers a handler for a specific event
func (sr *ScriptRegistry) RegisterHandler(event ScriptEvent, handler ScriptHandler) {
	sr.handlers[event] = append(sr.handlers[event], handler)
}

// ExecuteHandlers executes all registered handlers for an event
func (sr *ScriptRegistry) ExecuteHandlers(ctx context.Context, event ScriptEvent, scriptCtx *ScriptContext) error {
	handlers, exists := sr.handlers[event]
	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, scriptCtx); err != nil {
			return fmt.Errorf("script handler failed: %w", err)
		}
	}

	return nil
}

// GetRegisteredEvents returns all events that have registered handlers
func (sr *ScriptRegistry) GetRegisteredEvents() []ScriptEvent {
	events := make([]ScriptEvent, 0, len(sr.handlers))
	for event := range sr.handlers {
		events = append(events, event)
	}
	return events
}
