// Package fs provides file system utilities including file operations, watching, and management.
//
// Features:
//   - File Operations: Create, read, write, and manage files
//   - File Watching: Monitor file system changes with fsnotify integration
//   - Directory Operations: Recursive directory operations and management
//   - Path Utilities: Cross-platform path handling and manipulation
//   - Context Support: Full context.Context integration for cancellation
//   - Event Handling: Rich event system for file system changes
//
// Example:
//   watcher := fs.NewWatcher()
//   watcher.Watch("/path/to/watch")
//   events := watcher.Events()
package fs