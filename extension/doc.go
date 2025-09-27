// Package extension provides core extension management functionality for installing, managing, and running Tyk CLI extensions.
//
// Features:
//   - Extension Discovery: Search and discover Tyk CLI extensions from GitHub
//   - Extension Installation: Install extensions from GitHub repositories
//   - Extension Management: List, remove, and manage installed extensions
//   - GitHub Integration: Full GitHub API integration for extension discovery
//   - Hook Integration: Built-in hook system for extension lifecycle events
//   - Configuration Management: Persistent configuration and metadata storage
//   - Context Support: Full context.Context integration for all operations
//
// Example:
//   installer := extension.NewInstaller("/config/dir")
//   extensions, err := installer.SearchExtensions(ctx, "tyk", 10)
//   err = installer.InstallExtension(ctx, "owner", "repo")
package extension