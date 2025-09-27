// Package command provides structured command handling with Cobra integration.
//
// Features:
//   - Cobra Integration: Built on top of the popular Cobra CLI library
//   - Logger Support: Built-in logging capabilities
//   - Context Support: Full context.Context integration
//   - Command Creation: Helper functions for creating commands
//   - Long Description Support: Support for detailed command descriptions
//
// Example:
//   cmd := command.New("myapp", "Short description", handler)
//   cmd.SetLogger(logger)
//   cmd.SetContext(ctx)
package command