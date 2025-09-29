package hook

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ScriptExecutor executes external script hooks.
type ScriptExecutor struct {
	hookDir    string
	logger     *zap.Logger
	discovered map[Type][]string // Cache for discovered scripts
}

// NewScriptExecutor creates a new script hook executor.
func NewScriptExecutor(logger *zap.Logger, hookDir string) *ScriptExecutor {
	return &ScriptExecutor{
		hookDir:    hookDir,
		logger:     logger,
		discovered: make(map[Type][]string),
	}
}

// Execute executes script hooks for a specific type.
func (e *ScriptExecutor) Execute(ctx context.Context, hookType Type, data *Data) error {
	if e.hookDir == "" {
		return nil
	}

	// Look for hook scripts in the directory
	hookFiles, err := e.findHookFiles(ctx, hookType)
	if err != nil {
		return fmt.Errorf("failed to find hook files: %w", err)
	}

	if len(hookFiles) == 0 {
		return nil
	}

	// Execute each hook file
	for _, hookFile := range hookFiles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := e.executeHookFile(ctx, hookFile, data); err != nil {
			if e.logger != nil {
				e.logger.Error("Hook execution failed",
					zap.String("hook_file", hookFile),
					zap.String("hook_type", string(hookType)),
					zap.Error(err),
				)
			}
			return err
		}
	}

	return nil
}

// findHookFiles finds hook files for a specific hook type.
func (e *ScriptExecutor) findHookFiles(ctx context.Context, hookType Type) ([]string, error) {
	// Check cache first
	if scripts, exists := e.discovered[hookType]; exists {
		return scripts, nil
	}

	var hookFiles []string

	// Check if hookType is a direct file or directory
	hookPath := filepath.Join(e.hookDir, string(hookType))

	info, err := os.Stat(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Hook type doesn't exist, return empty list
			e.discovered[hookType] = []string{}
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to stat hook path %s: %w", hookPath, err)
	}

	if info.IsDir() {
		// Hook type is a directory, process its entries in lexicographic order
		hookFiles, err = e.processHookDirectory(ctx, hookPath)
		if err != nil {
			return nil, fmt.Errorf("failed to process hook directory %s: %w", hookPath, err)
		}
	} else {
		// Hook type is a file, check if it's a valid script
		if e.isValidScript(info.Name(), hookPath) {
			hookFiles = []string{hookPath}
		}
	}

	// Cache the discovered scripts
	e.discovered[hookType] = hookFiles

	if e.logger != nil {
		e.logger.Debug("Discovered scripts",
			zap.String("hook_type", string(hookType)),
			zap.String("hook_path", hookPath),
			zap.Bool("is_directory", info.IsDir()),
			zap.Int("count", len(hookFiles)),
			zap.Strings("scripts", hookFiles),
		)
	}

	return hookFiles, nil
}

// processHookDirectory processes a hook directory and returns scripts in lexicographic order.
func (e *ScriptExecutor) processHookDirectory(ctx context.Context, dirPath string) ([]string, error) {
	var scripts []string

	// Read directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// Sort entries by name for lexicographic order
	// Go's ReadDir already returns entries sorted by name, but let's be explicit
	for _, entry := range entries {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if entry.IsDir() {
			// Skip subdirectories for now
			continue
		}

		entryPath := filepath.Join(dirPath, entry.Name())

		// Check if it's a valid script
		if e.isValidScript(entry.Name(), entryPath) {
			scripts = append(scripts, entryPath)
		}
	}

	return scripts, nil
}

// isValidScript checks if a file is a valid executable script.
func (e *ScriptExecutor) isValidScript(filename, fullPath string) bool {
	// Script names should be hook type names (with dashes), no extension, and executable

	// Check if file has no extension
	if filepath.Ext(filename) != "" {
		return false
	}

	// Check if file is executable
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}

	// Check if it's a file (not directory)
	if info.IsDir() {
		return false
	}

	// Check if file has executable permissions
	return (info.Mode().Perm() & 0111) != 0
}

// executeHookFile executes a single hook file.
func (e *ScriptExecutor) executeHookFile(ctx context.Context, hookFile string, data *Data) error {
	// Create command with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, hookFile)

	// Set environment variables for the hook
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("HOOK_TYPE=%s", data.Type),
		fmt.Sprintf("EXTENSION=%s", data.Extension),
	)

	// Add metadata as environment variables
	for key, value := range data.Metadata {
		cmd.Env = append(cmd.Env, fmt.Sprintf("HOOK_%s=%v", strings.ToUpper(key), value))
	}

	// Execute the hook
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("hook execution failed: %w, output: %s", err, string(output))
	}

	if e.logger != nil {
		e.logger.Debug("Hook executed successfully",
			zap.String("hook_file", hookFile),
			zap.String("output", string(output)),
		)
	}

	return nil
}

// discoverScripts discovers scripts for specific hook types or all hook types.
// If no hook types are provided, it discovers all scripts. Otherwise, it discovers scripts for the specified types.
func (e *ScriptExecutor) discoverScripts(ctx context.Context, hookTypes ...Type) (map[Type][]string, error) {
	allScripts := make(map[Type][]string)

	if e.hookDir == "" {
		return allScripts, nil
	}

	// If specific hook types requested, check cache first
	if len(hookTypes) > 0 {
		// Check if all requested types are in cache
		allCached := true
		for _, hookType := range hookTypes {
			if scripts, exists := e.discovered[hookType]; exists {
				allScripts[hookType] = scripts
			} else {
				allCached = false
				break
			}
		}
		if allCached {
			return allScripts, nil
		}
	}

	// Read the hook directory to find hook types
	entries, err := os.ReadDir(e.hookDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read hook directory %s: %w", e.hookDir, err)
	}

	// Process each entry as a potential hook type
	for _, entry := range entries {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		entryName := entry.Name()

		// Skip hidden files and directories
		if strings.HasPrefix(entryName, ".") {
			continue
		}

		// Convert entry name to hook type
		currentHookType := Type(entryName)

		// If specific hook types requested, only process those types
		if len(hookTypes) > 0 {
			found := false
			for _, hookType := range hookTypes {
				if currentHookType == hookType {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Discover scripts for this hook type
		scripts, err := e.findHookFiles(ctx, currentHookType)
		if err != nil {
			if e.logger != nil {
				e.logger.Warn("Failed to discover scripts for hook type",
					zap.String("hook_type", entryName),
					zap.Error(err),
				)
			}
			continue
		}

		if len(scripts) > 0 {
			allScripts[currentHookType] = scripts
			// Update cache for this hook type
			e.discovered[currentHookType] = scripts
		}
	}

	if e.logger != nil {
		totalScripts := 0
		for _, scripts := range allScripts {
			totalScripts += len(scripts)
		}
		e.logger.Debug("Discovered scripts",
			zap.String("hook_dir", e.hookDir),
			zap.Strings("requested_types", hookTypesToStrings(hookTypes)),
			zap.Int("total_scripts", totalScripts),
			zap.Int("hook_types", len(allScripts)),
		)
	}

	return allScripts, nil
}

// getScriptDirectory returns the script directory path.
func (e *ScriptExecutor) getScriptDirectory() string {
	return e.hookDir
}
