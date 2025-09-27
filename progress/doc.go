// Package progress provides progress indicators including spinners and progress bars for CLI applications.
//
// Features:
//   - Spinner Support: Animated spinners for long-running operations
//   - Progress Bars: Visual progress bars with percentage and status
//   - Context Support: Full context.Context integration for cancellation
//   - Bubble Tea Integration: Built on the Bubble Tea TUI framework
//   - Customizable: Configurable messages, characters, and styling
//   - Thread-safe: Safe for concurrent use
//
// Example:
//   spinner := progress.New()
//   err := spinner.WithContext(ctx, "Processing...", func() error {
//       // Long-running operation
//       return nil
//   })
package progress