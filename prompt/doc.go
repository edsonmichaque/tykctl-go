// Package prompt provides interactive user input components for CLI applications.
//
// Features:
//   - User Input: Interactive prompts for user input
//   - Validation: Input validation and error handling
//   - Multiple Input Types: Support for various input formats
//   - Context Support: Full context.Context integration for cancellation
//   - Customizable: Configurable prompts and validation rules
//
// Example:
//   answer, err := prompt.Ask("What's your name?")
//   if err != nil {
//       log.Fatal(err)
//   }
package prompt