// Package jsonschema provides JSON Schema validation using the gojsonschema library with context support.
//
// Features:
//   - JSON Schema Validation: Validate JSON data against JSON Schema specifications
//   - Context Support: Full context.Context integration for cancellation and timeouts
//   - Multiple Input Sources: Support for strings, files, URLs, and Go objects
//   - Detailed Error Reporting: Comprehensive validation error information
//   - Schema Management: Extract metadata from schemas (version, title, description)
//   - Directory Validation: Validate all JSON files in a directory
//   - External Library Integration: Uses gojsonschema for robust validation
//
// Example:
//   validator, err := jsonschema.New(schemaString)
//   result, err := validator.ValidateString(ctx, jsonData)
//   isValid, err := jsonschema.IsValidString(ctx, jsonData, schemaString)
package jsonschema