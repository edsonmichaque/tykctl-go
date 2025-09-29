package hook

import (
	"errors"
	"fmt"
)

// Standard hook errors
var (
	// Hook errors
	ErrHookNotFound         = errors.New("hook not found")
	ErrHookInvalid          = errors.New("invalid hook format")
	ErrHookExecutionFailed  = errors.New("hook execution failed")
	ErrHookTimeout          = errors.New("hook execution timeout")
	ErrHookPermissionDenied = errors.New("hook permission denied")
	ErrContextCancelled     = errors.New("context cancelled")

	// Validation errors
	ErrValidationFailed     = errors.New("validation failed")
	ErrInvalidHookType      = errors.New("invalid hook type")
	ErrInvalidExtensionName = errors.New("invalid extension name")
	ErrInvalidMetadata      = errors.New("invalid metadata")

	// Executor errors
	ErrExecutorFailed     = errors.New("executor failed")
	ErrExternalHookFailed = errors.New("external hook failed")
	ErrRegoPolicyFailed   = errors.New("rego policy failed")

	// File system errors
	ErrHookDirNotFound   = errors.New("hook directory not found")
	ErrPolicyDirNotFound = errors.New("policy directory not found")
	ErrFileReadFailed    = errors.New("failed to read file")
	ErrFileWriteFailed   = errors.New("failed to write file")
)

// HookError represents a hook-specific error with context
type HookError struct {
	Type          Type
	ExtensionName string
	Message       string
	Err           error
}

// Error implements the error interface
func (e *HookError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("hook '%s' for extension '%s' failed: %s: %v", e.Type, e.ExtensionName, e.Message, e.Err)
	}
	return fmt.Sprintf("hook '%s' for extension '%s' failed: %s", e.Type, e.ExtensionName, e.Message)
}

// Unwrap returns the underlying error
func (e *HookError) Unwrap() error {
	return e.Err
}

// NewHookError creates a new hook error
func NewHookError(hookType Type, extensionName, message string, err error) *HookError {
	return &HookError{
		Type:          hookType,
		ExtensionName: extensionName,
		Message:       message,
		Err:           err,
	}
}

// ValidationError represents a validation-specific error
type ValidationError struct {
	Field   string
	Value   interface{}
	Rule    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation failed for field '%s': %s (rule: %s, value: %v): %v",
			e.Field, e.Message, e.Rule, e.Value, e.Err)
	}
	return fmt.Sprintf("validation failed for field '%s': %s (rule: %s, value: %v)",
		e.Field, e.Message, e.Rule, e.Value)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, rule, message string, err error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Rule:    rule,
		Message: message,
		Err:     err,
	}
}

// ExecutorError represents an executor-specific error
type ExecutorError struct {
	Executor string
	Type     Type
	Message  string
	Err      error
}

// Error implements the error interface
func (e *ExecutorError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("executor error (%s) for hook type '%s': %s: %v",
			e.Executor, e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("executor error (%s) for hook type '%s': %s",
		e.Executor, e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *ExecutorError) Unwrap() error {
	return e.Err
}

// NewExecutorError creates a new executor error
func NewExecutorError(executor string, hookType Type, message string, err error) *ExecutorError {
	return &ExecutorError{
		Executor: executor,
		Type:     hookType,
		Message:  message,
		Err:      err,
	}
}

// Error helpers

// IsHookError checks if an error is a hook error
func IsHookError(err error) bool {
	var hookErr *HookError
	return errors.As(err, &hookErr)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsExecutorError checks if an error is an executor error
func IsExecutorError(err error) bool {
	var executorErr *ExecutorError
	return errors.As(err, &executorErr)
}

// WrapError wraps an error with additional context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorf wraps an error with formatted additional context
func WrapErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

// ErrWithContext creates an error with context information
func ErrWithContext(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
