package template

import (
	"errors"
	"testing"
)

func TestTemplateError(t *testing.T) {
	// Test basic template error
	err := NewTemplateError(ErrorTypeLoadFailed, "file not found", nil)
	expected := "file not found"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Test template error with underlying error
	underlyingErr := errors.New("permission denied")
	err = NewTemplateError(ErrorTypeLoadFailed, "file not found", underlyingErr)
	expected = "file not found: permission denied"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Test unwrap
	if err.Unwrap() != underlyingErr {
		t.Error("Expected unwrap to return underlying error")
	}
}

func TestErrorTypeChecks(t *testing.T) {
	// Test IsTemplateError
	templateErr := NewTemplateError(ErrorTypeLoadFailed, "test", nil)
	if !IsTemplateError(templateErr, ErrorTypeLoadFailed) {
		t.Error("Expected IsTemplateError to return true")
	}
	if IsTemplateError(errors.New("other error"), ErrorTypeLoadFailed) {
		t.Error("Expected IsTemplateError to return false")
	}

	// Test different error types
	noSourceErr := NewTemplateError(ErrorTypeNoSource, "no source", nil)
	if !IsTemplateError(noSourceErr, ErrorTypeNoSource) {
		t.Error("Expected IsTemplateError to return true for ErrorTypeNoSource")
	}

	builtinNotFoundErr := NewTemplateError(ErrorTypeBuiltinNotFound, "not found", nil)
	if !IsTemplateError(builtinNotFoundErr, ErrorTypeBuiltinNotFound) {
		t.Error("Expected IsTemplateError to return true for ErrorTypeBuiltinNotFound")
	}

	invalidSourceErr := NewTemplateError(ErrorTypeInvalidSource, "invalid", nil)
	if !IsTemplateError(invalidSourceErr, ErrorTypeInvalidSource) {
		t.Error("Expected IsTemplateError to return true for ErrorTypeInvalidSource")
	}
}

func TestStandardErrors(t *testing.T) {
	// Test that standard errors are defined
	standardErrors := []error{
		ErrNoSourceSpecified,
		ErrBuiltinTemplateNotFound,
		ErrTemplateLoadFailed,
		ErrInvalidTemplateSource,
	}

	for _, err := range standardErrors {
		if err == nil {
			t.Error("Expected standard error to be non-nil")
		}
		if err.Error() == "" {
			t.Error("Expected standard error to have non-empty message")
		}
	}
}

func TestErrorTypes(t *testing.T) {
	// Test error type constants
	errorTypes := []ErrorType{
		ErrorTypeNoSource,
		ErrorTypeBuiltinNotFound,
		ErrorTypeLoadFailed,
		ErrorTypeInvalidSource,
	}

	for i, errorType := range errorTypes {
		if int(errorType) != i {
			t.Errorf("Expected error type %d to have value %d", i, int(errorType))
		}
	}
}
