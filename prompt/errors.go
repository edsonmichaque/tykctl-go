package prompt

import (
	"errors"
	"fmt"
)

// Standard prompt errors
var (
	// Input errors
	ErrInputFailed        = errors.New("input failed")
	ErrInputTimeout       = errors.New("input timeout")
	ErrInputCancelled     = errors.New("input cancelled")
	ErrInputEOF           = errors.New("end of input")
	ErrInputInvalid       = errors.New("invalid input")
	ErrInputEmpty         = errors.New("input is empty")
	ErrInputTooLong       = errors.New("input too long")

	// Validation errors
	ErrValidationFailed   = errors.New("validation failed")
	ErrInvalidFormat      = errors.New("invalid format")
	ErrInvalidChoice      = errors.New("invalid choice")
	ErrChoiceOutOfRange   = errors.New("choice out of range")
	ErrInvalidNumber      = errors.New("invalid number")
	ErrInvalidBoolean     = errors.New("invalid boolean")
	ErrInvalidSelection   = errors.New("invalid selection")

	// Option errors
	ErrNoOptions          = errors.New("no options provided")
	ErrEmptyOptions       = errors.New("options list is empty")
	ErrInvalidOption      = errors.New("invalid option")
	ErrDuplicateOptions   = errors.New("duplicate options")

	// Terminal errors
	ErrNotInteractive     = errors.New("not interactive")
	ErrTerminalError      = errors.New("terminal error")
	ErrTerminalNotTTY     = errors.New("terminal is not a TTY")
	ErrTerminalUnavailable = errors.New("terminal unavailable")

	// Reader errors
	ErrReaderFailed       = errors.New("reader failed")
	ErrReaderClosed       = errors.New("reader closed")
	ErrReaderTimeout      = errors.New("reader timeout")

	// Password errors
	ErrPasswordFailed     = errors.New("password input failed")
	ErrPasswordEmpty      = errors.New("password cannot be empty")
	ErrPasswordTooShort   = errors.New("password too short")
	ErrPasswordTooLong    = errors.New("password too long")

	// Confirmation errors
	ErrConfirmationFailed = errors.New("confirmation failed")
	ErrInvalidConfirmation = errors.New("invalid confirmation")
	ErrConfirmationRequired = errors.New("confirmation required")

	// Multi-select errors
	ErrMultiSelectFailed  = errors.New("multi-select failed")
	ErrNoSelections       = errors.New("no selections made")
	ErrTooManySelections  = errors.New("too many selections")
	ErrInvalidSelections  = errors.New("invalid selections")
)

// PromptError represents a prompt-specific error with context
type PromptError struct {
	Type      string
	Message   string
	Question  string
	Input     string
	Err       error
}

// Error implements the error interface
func (e *PromptError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (question: %q, input: %q): %v", 
			e.Type, e.Message, e.Question, e.Input, e.Err)
	}
	return fmt.Sprintf("%s: %s (question: %q, input: %q)", 
		e.Type, e.Message, e.Question, e.Input)
}

// Unwrap returns the underlying error
func (e *PromptError) Unwrap() error {
	return e.Err
}

// NewPromptError creates a new prompt error
func NewPromptError(errType, message, question, input string, err error) *PromptError {
	return &PromptError{
		Type:     errType,
		Message:  message,
		Question: question,
		Input:    input,
		Err:      err,
	}
}

// ValidationError represents a validation-specific error
type ValidationError struct {
	Field     string
	Value     interface{}
	Rule      string
	Message   string
	Question  string
	Err       error
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation failed for field '%s': %s (rule: %s, value: %v, question: %q): %v",
			e.Field, e.Message, e.Rule, e.Value, e.Question, e.Err)
	}
	return fmt.Sprintf("validation failed for field '%s': %s (rule: %s, value: %v, question: %q)",
		e.Field, e.Message, e.Rule, e.Value, e.Question)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, rule, message, question string, err error) *ValidationError {
	return &ValidationError{
		Field:    field,
		Value:    value,
		Rule:     rule,
		Message:  message,
		Question: question,
		Err:      err,
	}
}

// ChoiceError represents a choice-specific error
type ChoiceError struct {
	Choice    int
	MaxChoice int
	Options   []string
	Message   string
	Question  string
	Err       error
}

// Error implements the error interface
func (e *ChoiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("choice error: %s (choice: %d, max: %d, question: %q): %v",
			e.Message, e.Choice, e.MaxChoice, e.Question, e.Err)
	}
	return fmt.Sprintf("choice error: %s (choice: %d, max: %d, question: %q)",
		e.Message, e.Choice, e.MaxChoice, e.Question)
}

// Unwrap returns the underlying error
func (e *ChoiceError) Unwrap() error {
	return e.Err
}

// NewChoiceError creates a new choice error
func NewChoiceError(choice, maxChoice int, options []string, message, question string, err error) *ChoiceError {
	return &ChoiceError{
		Choice:    choice,
		MaxChoice: maxChoice,
		Options:   options,
		Message:   message,
		Question:  question,
		Err:       err,
	}
}

// InputError represents an input-specific error
type InputError struct {
	Input     string
	Expected  string
	Message   string
	Question  string
	Err       error
}

// Error implements the error interface
func (e *InputError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("input error: %s (input: %q, expected: %q, question: %q): %v",
			e.Message, e.Input, e.Expected, e.Question, e.Err)
	}
	return fmt.Sprintf("input error: %s (input: %q, expected: %q, question: %q)",
		e.Message, e.Input, e.Expected, e.Question)
}

// Unwrap returns the underlying error
func (e *InputError) Unwrap() error {
	return e.Err
}

// NewInputError creates a new input error
func NewInputError(input, expected, message, question string, err error) *InputError {
	return &InputError{
		Input:    input,
		Expected: expected,
		Message:  message,
		Question: question,
		Err:      err,
	}
}

// Error helpers

// IsPromptError checks if an error is a prompt error
func IsPromptError(err error) bool {
	var promptErr *PromptError
	return errors.As(err, &promptErr)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsChoiceError checks if an error is a choice error
func IsChoiceError(err error) bool {
	var choiceErr *ChoiceError
	return errors.As(err, &choiceErr)
}

// IsInputError checks if an error is an input error
func IsInputError(err error) bool {
	var inputErr *InputError
	return errors.As(err, &inputErr)
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

// Error constructors for common scenarios

// NewInputFailedError creates an error for input failure
func NewInputFailedError(question, input string, err error) *PromptError {
	return NewPromptError("input_failed", "failed to read input", question, input, err)
}

// NewInvalidChoiceError creates an error for invalid choice
func NewInvalidChoiceError(choice int, maxChoice int, options []string, question string) *ChoiceError {
	return NewChoiceError(choice, maxChoice, options, "invalid choice", question, nil)
}

// NewChoiceOutOfRangeError creates an error for choice out of range
func NewChoiceOutOfRangeError(choice int, maxChoice int, options []string, question string) *ChoiceError {
	return NewChoiceError(choice, maxChoice, options, "choice out of range", question, nil)
}

// NewInvalidNumberError creates an error for invalid number
func NewInvalidNumberError(input, question string, err error) *InputError {
	return NewInputError(input, "valid number", "invalid number format", question, err)
}

// NewInvalidBooleanError creates an error for invalid boolean
func NewInvalidBooleanError(input, question string) *InputError {
	return NewInputError(input, "y/n or yes/no", "invalid boolean format", question, nil)
}

// NewNoOptionsError creates an error for no options provided
func NewNoOptionsError(question string) *PromptError {
	return NewPromptError("no_options", "no options provided", question, "", nil)
}

// NewEmptyInputError creates an error for empty input
func NewEmptyInputError(question string) *InputError {
	return NewInputError("", "non-empty input", "input cannot be empty", question, nil)
}

// NewNotInteractiveError creates an error for non-interactive terminal
func NewNotInteractiveError(question string) *PromptError {
	return NewPromptError("not_interactive", "terminal is not interactive", question, "", nil)
}

// NewPasswordError creates an error for password input failure
func NewPasswordError(question string, err error) *PromptError {
	return NewPromptError("password_failed", "password input failed", question, "", err)
}

// NewConfirmationError creates an error for confirmation failure
func NewConfirmationError(input, question string) *InputError {
	return NewInputError(input, "y/n or yes/no", "invalid confirmation", question, nil)
}
