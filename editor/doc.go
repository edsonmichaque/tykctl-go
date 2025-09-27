// Package editor provides file editing capabilities with external editors.
//
// Features:
//   - External Editor Integration: Use system editors (vim, nano, VS Code, etc.)
//   - Environment Variable Support: Respects EDITOR, VISUAL, and TYKCTL_EDITOR
//   - String Editing: Edit strings in temporary files
//   - File Editing: Edit existing files
//   - Context Support: Full context.Context integration for cancellation
//   - Timeout Support: Configurable timeouts for editor sessions
//
// Example:
//   editor := editor.New()
//   content, err := editor.EditString(ctx, "Initial content")
package editor