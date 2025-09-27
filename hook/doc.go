// Package hook provides a flexible hook system for tykctl extensions with support for builtin, external, and Rego hooks.
//
// Features:
//   - Multiple Hook Types: Support for builtin, external, and Rego policy hooks
//   - Functional Options: Clean configuration using functional options pattern
//   - Generic Event System: Extensions define their own event types
//   - Context Passing: Rich context with custom data for extensions
//   - Error Handling: Proper error propagation and logging
//   - Flexible Execution: Extensions control when and how hooks are executed
//   - Rego Policy Support: Integration with Open Policy Agent (OPA) for policy-based hooks
//
// Example:
//   hm := hook.New()
//   hm.Register(ctx, "before-save", func(ctx context.Context, data interface{}) error {
//       // Custom logic
//       return nil
//   })
package hook