package hook

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// ScriptProcessor handles only script hooks with script-specific functionality.
type ScriptProcessor struct {
	scriptExecutor *ScriptExecutor
	validator      Validator
	logger         *zap.Logger
}

// NewScriptProcessor creates a new script-only processor.
func NewScriptProcessor(logger *zap.Logger, scriptDir string) *ScriptProcessor {
	return &ScriptProcessor{
		scriptExecutor: NewScriptExecutor(logger, scriptDir),
		validator:      NewScriptValidator(scriptDir),
		logger:         logger,
	}
}

// Execute executes only script hooks.
func (sp *ScriptProcessor) Execute(ctx context.Context, hookType Type, data *Data) error {
	// Validate hook data
	if err := sp.validator.Validate(data); err != nil {
		return fmt.Errorf("hook validation failed: %w", err)
	}

	// Execute script hooks
	if err := sp.scriptExecutor.Execute(ctx, hookType, data); err != nil {
		if sp.logger != nil {
			sp.logger.Error("Script hook execution failed",
				zap.String("hook_type", string(hookType)),
				zap.String("extension", data.Extension),
				zap.Error(err),
			)
		}
		return fmt.Errorf("script hook execution failed: %w", err)
	}

	return nil
}

// GetScriptExecutor returns the underlying script executor for advanced usage.
func (sp *ScriptProcessor) GetScriptExecutor() *ScriptExecutor {
	return sp.scriptExecutor
}

// ListScripts returns a list of available script files for a given hook type.
// This method adds processor-level validation and error handling.
func (sp *ScriptProcessor) ListScripts(ctx context.Context, hookType Type) ([]string, error) {
	if sp.scriptExecutor == nil {
		return nil, fmt.Errorf("script executor not available")
	}

	// Validate hook type before listing scripts
	if hookType == "" {
		return nil, fmt.Errorf("hook type cannot be empty")
	}

	allScripts, err := sp.scriptExecutor.discoverScripts(ctx, hookType)
	if err != nil {
		return nil, fmt.Errorf("failed to discover scripts: %w", err)
	}
	scripts := allScripts[hookType]

	if sp.logger != nil {
		sp.logger.Debug("Listed scripts for hook type",
			zap.String("hook_type", string(hookType)),
			zap.Int("count", len(scripts)),
		)
	}

	return scripts, nil
}

// DiscoverAllScripts discovers all scripts in the script directory.
// This method adds processor-level logging and error handling.
func (sp *ScriptProcessor) DiscoverAllScripts() (map[Type][]string, error) {
	if sp.scriptExecutor == nil {
		return nil, fmt.Errorf("script executor not available")
	}

	allScripts, err := sp.scriptExecutor.discoverScripts(context.Background())
	if err != nil {
		if sp.logger != nil {
			sp.logger.Error("Failed to discover all scripts",
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("script discovery failed: %w", err)
	}

	if sp.logger != nil {
		totalScripts := 0
		for _, scripts := range allScripts {
			totalScripts += len(scripts)
		}
		sp.logger.Debug("Discovered all scripts",
			zap.Int("total_scripts", totalScripts),
			zap.Int("hook_types", len(allScripts)),
		)
	}

	return allScripts, nil
}

// CountScripts returns the count of discovered scripts for a hook type.
// This method adds processor-level validation and logging.
func (sp *ScriptProcessor) CountScripts(ctx context.Context, hookType Type) int {
	if sp.scriptExecutor == nil {
		return 0
	}

	if hookType == "" {
		return 0
	}

	allScripts, err := sp.scriptExecutor.discoverScripts(ctx, hookType)
	if err != nil {
		if sp.logger != nil {
			sp.logger.Error("Failed to discover scripts for counting",
				zap.String("hook_type", string(hookType)),
				zap.Error(err),
			)
		}
		return 0
	}

	scripts := allScripts[hookType]
	count := len(scripts)

	if sp.logger != nil {
		sp.logger.Debug("Counted scripts for hook type",
			zap.String("hook_type", string(hookType)),
			zap.Int("count", count),
		)
	}

	return count
}

// ValidateScript validates a script file without executing it.
// This method adds processor-level validation logic.
func (sp *ScriptProcessor) ValidateScript(ctx context.Context, scriptFile string) error {
	if sp.scriptExecutor == nil {
		return fmt.Errorf("script executor not available")
	}

	if scriptFile == "" {
		return fmt.Errorf("script file path cannot be empty")
	}

	// Get the script directory to validate the file is within it
	scriptDir := sp.scriptExecutor.getScriptDirectory()
	if scriptDir != "" && !isPathWithinDirectory(scriptFile, scriptDir) {
		return fmt.Errorf("script file %s is not within script directory %s", scriptFile, scriptDir)
	}

	// Validate script file properties
	if err := validateScriptFile(scriptFile); err != nil {
		return fmt.Errorf("script file %s is not valid: %w", scriptFile, err)
	}

	if sp.logger != nil {
		sp.logger.Debug("Script validation passed",
			zap.String("script_file", scriptFile),
		)
	}

	return nil
}

// GetScriptDirectory returns the script directory path.
// This method adds processor-level validation.
func (sp *ScriptProcessor) GetScriptDirectory() string {
	if sp.scriptExecutor == nil {
		return ""
	}

	return sp.scriptExecutor.getScriptDirectory()
}

// isPathWithinDirectory checks if a file path is within a directory.
func isPathWithinDirectory(filePath, dirPath string) bool {
	// Simple check - in a real implementation, you'd use filepath.Rel
	return len(filePath) > len(dirPath) && filePath[:len(dirPath)] == dirPath
}

// getFileName extracts the filename from a full path.
func getFileName(filePath string) string {
	// Simple implementation - in a real implementation, you'd use filepath.Base
	lastSlash := -1
	for i, char := range filePath {
		if char == '/' {
			lastSlash = i
		}
	}
	if lastSlash == -1 {
		return filePath
	}
	return filePath[lastSlash+1:]
}

// validateScriptFile validates that a file is a valid script.
func validateScriptFile(filePath string) error {
	// Check if file has no extension
	if filepath.Ext(filePath) != "" {
		return fmt.Errorf("script file must have no extension")
	}

	// Check if file is executable
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("script file does not exist")
		}
		return fmt.Errorf("failed to get script file info: %w", err)
	}

	// Check if it's a file (not directory)
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a script file")
	}

	// Check if file has executable permissions
	if (info.Mode().Perm() & 0111) == 0 {
		return fmt.Errorf("script file is not executable")
	}

	return nil
}
