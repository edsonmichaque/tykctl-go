package template

import "errors"

// Template resolution error types
var (
	// ErrNoSourceSpecified indicates no template source was provided
	ErrNoSourceSpecified = errors.New("no template source specified: use WithName, WithURL, or WithFile")

	// ErrBuiltinTemplateNotFound indicates a built-in template was not found
	ErrBuiltinTemplateNotFound = errors.New("built-in template not found")

	// ErrTemplateLoadFailed indicates template loading failed
	ErrTemplateLoadFailed = errors.New("template load failed")

	// ErrInvalidTemplateSource indicates an invalid template source was provided
	ErrInvalidTemplateSource = errors.New("invalid template source")
)

// TemplateError represents a template resolution error with context
type TemplateError struct {
	Type    ErrorType
	Message string
	Err     error
}

// ErrorType represents the type of template error
type ErrorType int

const (
	// ErrorTypeNoSource indicates no source was specified
	ErrorTypeNoSource ErrorType = iota

	// ErrorTypeBuiltinNotFound indicates built-in template not found
	ErrorTypeBuiltinNotFound

	// ErrorTypeLoadFailed indicates template loading failed
	ErrorTypeLoadFailed

	// ErrorTypeInvalidSource indicates invalid template source
	ErrorTypeInvalidSource
)

// Error implements the error interface
func (e *TemplateError) Error() string {
	if e.Message != "" && e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "template error"
}

// Unwrap returns the underlying error
func (e *TemplateError) Unwrap() error {
	return e.Err
}

// NewTemplateError creates a new template error
func NewTemplateError(errorType ErrorType, message string, err error) *TemplateError {
	return &TemplateError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

// IsTemplateError checks if an error is a template error of a specific type
func IsTemplateError(err error, errorType ErrorType) bool {
	var templateErr *TemplateError
	if errors.As(err, &templateErr) {
		return templateErr.Type == errorType
	}
	return false
}
