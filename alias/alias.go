package alias

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Manager handles alias operations for CLI extensions
type Manager struct {
	configProvider ConfigProvider
	reservedNames  []string
}

// ConfigProvider provides configuration access for aliases
type ConfigProvider interface {
	SetAlias(ctx context.Context, name, expansion string) error
	GetAlias(ctx context.Context, name string) (string, bool)
	DeleteAlias(ctx context.Context, name string) error
	ListAliases(ctx context.Context) map[string]string
}

// NewManager creates a new alias manager
func NewManager(configProvider ConfigProvider, reservedNames []string) *Manager {
	return &Manager{
		configProvider: configProvider,
		reservedNames:  reservedNames,
	}
}

// SetAlias sets an alias with the given name and expansion
func (m *Manager) SetAlias(ctx context.Context, name, expansion string) error {
	// Validate alias name
	if err := m.validateAliasName(name); err != nil {
		return err
	}

	// Validate expansion
	if err := m.validateExpansion(expansion); err != nil {
		return err
	}

	return m.configProvider.SetAlias(ctx, name, expansion)
}

// GetAlias retrieves an alias by name
func (m *Manager) GetAlias(ctx context.Context, name string) (string, bool) {
	return m.configProvider.GetAlias(ctx, name)
}

// DeleteAlias deletes an alias by name
func (m *Manager) DeleteAlias(ctx context.Context, name string) error {
	// Check if alias exists
	if _, exists := m.GetAlias(ctx, name); !exists {
		return fmt.Errorf("alias '%s' does not exist", name)
	}

	return m.configProvider.DeleteAlias(ctx, name)
}

// ListAliases returns all configured aliases
func (m *Manager) ListAliases(ctx context.Context) map[string]string {
	return m.configProvider.ListAliases(ctx)
}

// ExecuteAlias executes an alias with the given arguments
func (m *Manager) ExecuteAlias(ctx context.Context, aliasName string, args []string) error {
	expansion, exists := m.GetAlias(ctx, aliasName)
	if !exists {
		return fmt.Errorf("alias '%s' not found", aliasName)
	}

	return m.executeExpansion(ctx, aliasName, expansion, args)
}

// executeExpansion executes the alias expansion
func (m *Manager) executeExpansion(ctx context.Context, aliasName, expansion string, args []string) error {
	// Check if it's a shell alias (prefixed with !)
	if strings.HasPrefix(expansion, "!") {
		return m.executeShellAlias(ctx, aliasName, expansion[1:], args)
	}

	// Regular alias - expand parameters and parse as new command
	expandedCmd := m.expandParameters(expansion, args)

	// Split the expanded command into parts
	cmdParts := strings.Fields(expandedCmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("empty alias expansion")
	}

	// Create a new command with the expanded arguments
	// We need to find the current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	execCmd := exec.CommandContext(ctx, execPath, cmdParts...)
	execCmd.Env = os.Environ()
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	err = execCmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		return err
	}

	return nil
}

// executeShellAlias executes a shell alias
func (m *Manager) executeShellAlias(ctx context.Context, aliasName, shellCmd string, args []string) error {
	// Expand parameters in the shell command
	expandedCmd := m.expandParameters(shellCmd, args)

	// Determine shell to use
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	// Run the shell command
	cmd := exec.CommandContext(ctx, shell, "-c", expandedCmd)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		return err
	}

	return nil
}

// expandParameters expands $1, $2, etc. in alias expansion
func (m *Manager) expandParameters(expansion string, args []string) string {
	// Replace $1, $2, etc. with actual arguments
	result := expansion

	// Replace numbered parameters ($1, $2, etc.)
	for i, arg := range args {
		param := fmt.Sprintf("$%d", i+1)
		result = strings.ReplaceAll(result, param, arg)
	}

	// Replace $* with all arguments
	if strings.Contains(result, "$*") {
		allArgs := strings.Join(args, " ")
		result = strings.ReplaceAll(result, "$*", allArgs)
	}

	// Replace $@ with all arguments as separate words
	if strings.Contains(result, "$@") {
		allArgs := strings.Join(args, " ")
		result = strings.ReplaceAll(result, "$@", allArgs)
	}

	return result
}

// validateAliasName validates an alias name
func (m *Manager) validateAliasName(name string) error {
	if name == "" {
		return fmt.Errorf("alias name cannot be empty")
	}

	// Check for reserved names
	for _, reserved := range m.reservedNames {
		if name == reserved {
			return fmt.Errorf("'%s' is a reserved command name", name)
		}
	}

	// Check for invalid characters
	if strings.ContainsAny(name, " \t\n\r") {
		return fmt.Errorf("alias name cannot contain whitespace")
	}

	// Check for shell metacharacters
	if strings.ContainsAny(name, "&|;()<>") {
		return fmt.Errorf("alias name cannot contain shell metacharacters")
	}

	return nil
}

// validateExpansion validates an alias expansion
func (m *Manager) validateExpansion(expansion string) error {
	if expansion == "" {
		return fmt.Errorf("alias expansion cannot be empty")
	}

	// Check for basic syntax issues
	if strings.HasPrefix(expansion, " ") {
		return fmt.Errorf("alias expansion cannot start with whitespace")
	}

	return nil
}

// IsShellAlias checks if an expansion is a shell alias
func (m *Manager) IsShellAlias(expansion string) bool {
	return strings.HasPrefix(expansion, "!")
}

// GetAliasType returns the type of alias (shell or command)
func (m *Manager) GetAliasType(expansion string) string {
	if m.IsShellAlias(expansion) {
		return "shell"
	}
	return "command"
}

// ExpandAliasPreview shows how an alias would be expanded with given arguments
func (m *Manager) ExpandAliasPreview(expansion string, args []string) string {
	return m.expandParameters(expansion, args)
}