// Package version provides semantic versioning utilities for CLI applications.
//
// Features:
//   - Semantic Versioning: Support for semantic version parsing and comparison
//   - Version Parsing: Parse version strings into structured data
//   - Version Comparison: Compare versions using semantic versioning rules
//   - Version Validation: Validate version string formats
//   - String Formatting: Format versions for display
//
// Example:
//   v := version.New("1.2.3")
//   fmt.Println(v.String()) // "1.2.3"
//   isNewer := v.IsNewerThan("1.0.0")
package version