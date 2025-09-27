// Package terminal provides terminal utilities and styling for CLI applications.
//
// Features:
//   - Terminal Detection: Detect terminal capabilities and features
//   - Color Support: Cross-platform color support for terminal output
//   - Styling: Text styling including bold, italic, underline
//   - ANSI Escape Codes: Support for ANSI escape sequences
//   - Cross-platform: Works on Windows, macOS, and Linux
//
// Example:
//   term := terminal.New()
//   styled := term.Bold("Important text")
//   colored := term.Red("Error message")
package terminal