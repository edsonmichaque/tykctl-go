package alias

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// CommandBuilder helps build alias commands for Cobra
type CommandBuilder struct {
	manager *Manager
}

// NewCommandBuilder creates a new command builder
func NewCommandBuilder(manager *Manager) *CommandBuilder {
	return &CommandBuilder{
		manager: manager,
	}
}

// BuildAliasCommand creates the main alias command
func (cb *CommandBuilder) BuildAliasCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage command aliases",
		Long: `Manage command aliases for CLI commands.

Aliases allow you to create shortcuts for commonly used commands.
You can use parameters in aliases with $1, $2, etc.

Examples:
  tykctl alias set co "products checkout"
  tykctl alias set users "users list"
  tykctl alias set myproducts '!tykctl products list --json id,title --jq ".[] | \"#\(.id): \(.title)\""'
  tykctl alias list
  tykctl alias delete co`,
	}

	cmd.AddCommand(cb.BuildSetCommand())
	cmd.AddCommand(cb.BuildListCommand())
	cmd.AddCommand(cb.BuildDeleteCommand())
	cmd.AddCommand(cb.BuildEditCommand())
	cmd.AddCommand(cb.BuildShowCommand())

	return cmd
}

// BuildSetCommand creates the alias set command
func (cb *CommandBuilder) BuildSetCommand() *cobra.Command {
	var shell bool

	cmd := &cobra.Command{
		Use:   "set <name> <expansion>",
		Short: "Set an alias",
		Long: `Set an alias for a command.

The expansion can be a simple command or a shell command (prefixed with !).
Use $1, $2, etc. to reference arguments passed to the alias.

Examples:
  tykctl alias set co "products checkout"
  tykctl alias set users "users list"
  tykctl alias set myproducts '!tykctl products list --json id,title --jq ".[] | \"#\(.id): \(.title)\""'
  tykctl alias set --shell cleanup`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			if shell {
				return cb.setAliasShell(ctx, args[0])
			}

			if len(args) < 2 {
				return fmt.Errorf("alias expansion is required")
			}

			name := args[0]
			expansion := strings.Join(args[1:], " ")

			if err := cb.manager.SetAlias(ctx, name, expansion); err != nil {
				return fmt.Errorf("failed to set alias: %w", err)
			}

			fmt.Printf("Alias '%s' set to '%s'\n", name, expansion)
			return nil
		},
	}

	cmd.Flags().BoolVar(&shell, "shell", false, "Open editor for multi-line shell alias")

	return cmd
}

// BuildListCommand creates the alias list command
func (cb *CommandBuilder) BuildListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		Long:  `List all configured aliases.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			aliases := cb.manager.ListAliases(ctx)
			if len(aliases) == 0 {
				fmt.Println("No aliases configured.")
				return nil
			}

			fmt.Println("Configured aliases:")
			for name, expansion := range aliases {
				aliasType := cb.manager.GetAliasType(expansion)
				fmt.Printf("  %s (%s): %s\n", name, aliasType, expansion)
			}
			return nil
		},
	}
}

// BuildDeleteCommand creates the alias delete command
func (cb *CommandBuilder) BuildDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete an alias",
		Long:  `Delete an alias by name.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			name := args[0]

			if err := cb.manager.DeleteAlias(ctx, name); err != nil {
				return fmt.Errorf("failed to delete alias: %w", err)
			}

			fmt.Printf("Alias '%s' deleted\n", name)
			return nil
		},
	}
}

// BuildEditCommand creates the alias edit command
func (cb *CommandBuilder) BuildEditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit aliases in config file",
		Long:  `Open the configuration file in your default editor to edit aliases.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// This would need to be implemented by the specific extension
			// as it depends on the config file location
			return fmt.Errorf("alias edit command must be implemented by the extension")
		},
	}
}

// BuildShowCommand creates the alias show command
func (cb *CommandBuilder) BuildShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show alias expansion",
		Long:  `Show the expansion for a specific alias.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			name := args[0]
			expansion, exists := cb.manager.GetAlias(ctx, name)
			if !exists {
				return fmt.Errorf("alias '%s' does not exist", name)
			}

			aliasType := cb.manager.GetAliasType(expansion)
			fmt.Printf("Alias: %s (%s)\n", name, aliasType)
			fmt.Printf("Expansion: %s\n", expansion)

			// Show preview with sample arguments
			if len(cmd.ValidArgs) > 0 {
				sampleArgs := cmd.ValidArgs[:min(3, len(cmd.ValidArgs))]
				preview := cb.manager.ExpandAliasPreview(expansion, sampleArgs)
				fmt.Printf("Preview (with args %v): %s\n", sampleArgs, preview)
			}

			return nil
		},
	}
}

// setAliasShell opens an editor for creating a multi-line shell alias
func (cb *CommandBuilder) setAliasShell(ctx context.Context, name string) error {
	// Create a temporary file for the alias content
	tmpFile, err := os.CreateTemp("", "tykctl-alias-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write a template to the file
	template := `#!/bin/bash
# Multi-line shell alias for: ` + name + `
# Edit this file to define your alias expansion
# Use $1, $2, etc. for arguments
# Example:
# tykctl products list --status $1
`
	if _, err := tmpFile.WriteString(template); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}
	tmpFile.Close()

	// Get the default editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // fallback to vi
	}

	// Open the temporary file in the editor
	editCmd := exec.Command(editor, tmpFile.Name())
	editCmd.Stdin = os.Stdin
	editCmd.Stdout = os.Stdout
	editCmd.Stderr = os.Stderr

	if err := editCmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	// Read the edited content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read edited content: %w", err)
	}

	// Extract the actual command (skip comments and shebang)
	lines := strings.Split(string(content), "\n")
	var expansionLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "#!/") {
			continue
		}
		expansionLines = append(expansionLines, line)
	}

	if len(expansionLines) == 0 {
		return fmt.Errorf("no command found in alias definition")
	}

	expansion := "!" + strings.Join(expansionLines, "; ")
	if err := cb.manager.SetAlias(ctx, name, expansion); err != nil {
		return fmt.Errorf("failed to save alias: %w", err)
	}

	fmt.Printf("Shell alias '%s' set to '%s'\n", name, expansion)
	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}