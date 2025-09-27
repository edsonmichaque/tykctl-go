// Package jq provides pure Go JQ integration for JSON manipulation and processing using gojq.
//
// Features:
//   - Pure Go Implementation: No external jq binary dependency
//   - JSON Processing: Process JSON data with jq programs
//   - Multiple Input Types: Support for strings, bytes, and Go objects
//   - Complex Queries: Full jq language support for advanced JSON manipulation
//   - Error Handling: Comprehensive error handling with Go error wrapping
//   - Cross-platform: Works consistently across all platforms
//
// Example:
//   result, err := jq.ProcessString(jsonData, ".users[0].name")
//   data, err := jq.Process(jsonBytes, ".field")
//   obj, err := jq.ProcessObject(myObject, ".property")
package jq