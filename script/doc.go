// Package script provides dynamic script execution capabilities for tykctl extensions.
//
// Features:
//   - Script Management: Create, enable, disable, and manage scripts
//   - Custom Event Handling: Define custom events for script execution
//   - File-based Scripts: Execute actual script files with proper environment setup
//   - Context Passing: Rich context with custom data for extensions
//   - Error Handling: Proper error propagation and logging
//   - Flexible Execution: Extensions control when and how scripts are executed
//
// Example:
//   sm := script.NewScriptManager("/scripts")
//   s, err := sm.CreateScript("my-script", "Description", "#!/bin/bash\necho 'Hello'")
//   err = sm.EnableScript("my-script")
package script