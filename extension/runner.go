package extension

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/edsonmichaque/tykctl-go/hook"
	"go.uber.org/zap"
)

// Runner executes tykctl extensions
type Runner struct {
	configDir string
	logger    *zap.Logger
	hooks     *hook.Manager
}

// NewRunner creates a new extension runner
func NewRunner(configDir string) *Runner {
	logger := zap.L()
	return &Runner{
		configDir: configDir,
		logger:    logger,
		hooks:     hook.New(),
	}
}

// NewRunnerWithHooks creates a new extension runner with custom hooks
func NewRunnerWithHooks(configDir string, hooks *hook.Manager) *Runner {
	logger := zap.L()
	return &Runner{
		configDir: configDir,
		logger:    logger,
		hooks:     hooks,
	}
}

// RunExtension executes an extension with optional custom environment variables
func (r *Runner) RunExtension(ctx context.Context, extensionName string, args []string, envVars ...map[string]string) error {
	extensionPath := r.findExtension(extensionName)
	if extensionPath == "" {
		return fmt.Errorf("extension '%s' not found. Run 'tykctl extension list' to see available extensions", extensionName)
	}

	// Execute before run hooks
	hookData := &hook.HookData{
		ExtensionName: extensionName,
		ExtensionPath: extensionPath,
		Metadata: map[string]interface{}{
			"args": args,
		},
	}

	if err := r.hooks.Execute(ctx, hook.HookTypeBeforeRun, hookData); err != nil {
		r.logger.Error("Before run hook failed", zap.Error(err))
		return fmt.Errorf("before run hook failed: %w", err)
	}

	r.logger.Info("Running extension",
		zap.String("name", extensionName),
		zap.String("path", extensionPath),
		zap.Strings("args", args))

	// Set up environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("TYK_CLI_CONFIG=%s", r.configDir))

	// Add custom environment variables if provided
	if len(envVars) > 0 && envVars[0] != nil {
		for key, value := range envVars[0] {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Pass debug flag to extension via environment variable
	if debug := os.Getenv("TYKCTL_DEBUG"); debug == "true" {
		env = append(env, "TYKCTL_DEBUG=true")
	} else {
		env = append(env, "TYKCTL_DEBUG=false")
	}

	// Execute the extension
	cmd := exec.CommandContext(ctx, extensionPath, args...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			r.logger.Error("Extension exited with error",
				zap.String("name", extensionName),
				zap.Int("exit_code", exitError.ExitCode()))
			return fmt.Errorf("extension %s failed with exit code %d", extensionName, exitError.ExitCode())
		}
		return fmt.Errorf("failed to run extension %s: %w", extensionName, err)
	}

	r.logger.Info("Extension completed successfully", zap.String("name", extensionName))
	return nil
}

// RunExtensionWithOutput executes an extension and captures its output with optional custom environment variables
func (r *Runner) RunExtensionWithOutput(ctx context.Context, extensionName string, args []string, envVars ...map[string]string) ([]byte, error) {
	extensionPath := r.findExtension(extensionName)
	if extensionPath == "" {
		return nil, fmt.Errorf("extension '%s' not found. Run 'tykctl extension list' to see available extensions", extensionName)
	}

	r.logger.Info("Running extension with output capture",
		zap.String("name", extensionName),
		zap.String("path", extensionPath),
		zap.Strings("args", args))

	// Set up environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("TYK_CLI_CONFIG=%s", r.configDir))

	// Add custom environment variables if provided
	if len(envVars) > 0 && envVars[0] != nil {
		for key, value := range envVars[0] {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Pass debug flag to extension via environment variable
	if debug := os.Getenv("TYKCTL_DEBUG"); debug == "true" {
		env = append(env, "TYKCTL_DEBUG=true")
	} else {
		env = append(env, "TYKCTL_DEBUG=false")
	}

	// Execute the extension and capture output
	cmd := exec.CommandContext(ctx, extensionPath, args...)
	cmd.Env = env
	cmd.Stdin = os.Stdin

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			r.logger.Error("Extension exited with error",
				zap.String("name", extensionName),
				zap.Int("exit_code", exitError.ExitCode()),
				zap.String("stderr", string(exitError.Stderr)))
			return output, fmt.Errorf("extension %s failed with exit code %d: %s", extensionName, exitError.ExitCode(), string(exitError.Stderr))
		}
		return output, fmt.Errorf("failed to run extension %s: %w", extensionName, err)
	}

	r.logger.Info("Extension completed successfully", zap.String("name", extensionName))
	return output, nil
}

// IsExtensionAvailable checks if an extension is available for execution
func (r *Runner) IsExtensionAvailable(extensionName string) bool {
	return r.findExtension(extensionName) != ""
}

// ListAvailableExtensions returns a list of all available extensions
func (r *Runner) ListAvailableExtensions() ([]string, error) {
	var extensions []string

	// Check XDG data directory
	dataDir := filepath.Join(xdg.DataHome, "tykctl", "extensions")
	if entries, err := os.ReadDir(dataDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				extName := entry.Name()
				if len(extName) > 8 && extName[:8] == "tykctl-" {
					extensions = append(extensions, extName[8:])
				}
			}
		}
	}

	// Check legacy config directory
	configDir := filepath.Join(xdg.ConfigHome, "tykctl", "extensions")
	if entries, err := os.ReadDir(configDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				extName := entry.Name()
				if len(extName) > 8 && extName[:8] == "tykctl-" {
					// Avoid duplicates
					found := false
					for _, existing := range extensions {
						if existing == extName[8:] {
							found = true
							break
						}
					}
					if !found {
						extensions = append(extensions, extName[8:])
					}
				}
			}
		}
	}

	return extensions, nil
}

// findExtension finds an extension binary
func (r *Runner) findExtension(name string) string {
	// Check XDG data directory
	dataDir := filepath.Join(xdg.DataHome, "tykctl", "extensions")
	xdgPath := filepath.Join(dataDir, fmt.Sprintf("tykctl-%s", name))
	if _, err := os.Stat(xdgPath); err == nil {
		return xdgPath
	}

	// Check legacy config directory
	configDir := filepath.Join(xdg.ConfigHome, "tykctl", "extensions")
	legacyPath := filepath.Join(configDir, fmt.Sprintf("tykctl-%s", name))
	if _, err := os.Stat(legacyPath); err == nil {
		return legacyPath
	}

	// Check PATH
	path, err := exec.LookPath(fmt.Sprintf("tykctl-%s", name))
	if err == nil {
		return path
	}

	return ""
}
