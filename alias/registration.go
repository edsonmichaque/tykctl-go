package alias

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Registrar handles alias registration for Cobra commands
type Registrar struct {
	manager *Manager
}

// NewRegistrar creates a new alias registrar
func NewRegistrar(manager *Manager) *Registrar {
	return &Registrar{
		manager: manager,
	}
}

// RegisterAliases loads aliases from config and registers them as subcommands
func (r *Registrar) RegisterAliases(ctx context.Context, rootCmd *cobra.Command) error {
	// Load aliases from config
	aliases := r.manager.ListAliases(ctx)

	// Get existing command names to avoid conflicts
	existingCommands := r.getExistingCommandNames(rootCmd)

	for aliasName, expansion := range aliases {
		// Skip if there's a command conflict
		if existingCommands[aliasName] {
			continue
		}

		// Create a subcommand for each alias
		aliasCmd := &cobra.Command{
			Use:   aliasName,
			Short: fmt.Sprintf("Alias for: %s", expansion),
			Long:  fmt.Sprintf("This is an alias for: %s", expansion),
			RunE: func(cmd *cobra.Command, args []string) error {
				return r.manager.ExecuteAlias(ctx, aliasName, args)
			},
		}

		// Add the alias as a subcommand
		rootCmd.AddCommand(aliasCmd)
	}

	return nil
}

// RegisterAliasesWithValidation registers aliases with validation
func (r *Registrar) RegisterAliasesWithValidation(ctx context.Context, rootCmd *cobra.Command, reservedNames []string) error {
	// Load aliases from config
	aliases := r.manager.ListAliases(ctx)

	// Get existing command names to avoid conflicts
	existingCommands := r.getExistingCommandNames(rootCmd)

	// Add reserved names to existing commands
	for _, reserved := range reservedNames {
		existingCommands[reserved] = true
	}

	var conflicts []string
	var registered int

	for aliasName, expansion := range aliases {
		// Check for conflicts
		if existingCommands[aliasName] {
			conflicts = append(conflicts, aliasName)
			continue
		}

		// Validate alias before registration
		if err := r.manager.validateAliasName(aliasName); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Invalid alias '%s': %v\n", aliasName, err)
			continue
		}

		// Create a subcommand for each alias
		aliasCmd := &cobra.Command{
			Use:   aliasName,
			Short: fmt.Sprintf("Alias for: %s", expansion),
			Long:  fmt.Sprintf("This is an alias for: %s", expansion),
			RunE: func(cmd *cobra.Command, args []string) error {
				return r.manager.ExecuteAlias(ctx, aliasName, args)
			},
		}

		// Add the alias as a subcommand
		rootCmd.AddCommand(aliasCmd)
		registered++
	}

	// Report conflicts if any
	if len(conflicts) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: Skipped %d aliases due to command conflicts: %v\n", len(conflicts), conflicts)
	}

	// Report registration summary
	if registered > 0 {
		fmt.Fprintf(os.Stderr, "Registered %d aliases\n", registered)
	}

	return nil
}

// getExistingCommandNames returns a map of existing command names
func (r *Registrar) getExistingCommandNames(rootCmd *cobra.Command) map[string]bool {
	existingCommands := make(map[string]bool)

	// Add root command name
	existingCommands[rootCmd.Name()] = true

	// Add all subcommand names and aliases
	for _, cmd := range rootCmd.Commands() {
		existingCommands[cmd.Name()] = true
		// Also check aliases
		for _, alias := range cmd.Aliases {
			existingCommands[alias] = true
		}
	}

	return existingCommands
}

// ValidateAliases validates all configured aliases
func (r *Registrar) ValidateAliases(ctx context.Context, reservedNames []string) []ValidationError {
	var errors []ValidationError

	aliases := r.manager.ListAliases(ctx)

	for aliasName, expansion := range aliases {
		// Check for reserved names
		for _, reserved := range reservedNames {
			if aliasName == reserved {
				errors = append(errors, ValidationError{
					AliasName: aliasName,
					Expansion: expansion,
					Error:     fmt.Sprintf("'%s' is a reserved command name", aliasName),
				})
				continue
			}
		}

		// Validate alias name
		if err := r.manager.validateAliasName(aliasName); err != nil {
			errors = append(errors, ValidationError{
				AliasName: aliasName,
				Expansion: expansion,
				Error:     err.Error(),
			})
		}

		// Validate expansion
		if err := r.manager.validateExpansion(expansion); err != nil {
			errors = append(errors, ValidationError{
				AliasName: aliasName,
				Expansion: expansion,
				Error:     err.Error(),
			})
		}
	}

	return errors
}

// ValidationError represents an alias validation error
type ValidationError struct {
	AliasName string
	Expansion string
	Error     string
}

// String returns a string representation of the validation error
func (ve ValidationError) String() string {
	return fmt.Sprintf("Alias '%s' (%s): %s", ve.AliasName, ve.Expansion, ve.Error)
}