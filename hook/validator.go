package hook

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator defines the interface for validating hook data.
type Validator interface {
	Validate(data *Data) error
}

// HookValidator is the concrete implementation of the Validator interface.
type HookValidator struct {
	// validation configuration
}

// NewValidator creates a new hook validator.
func NewValidator() *HookValidator {
	return &HookValidator{}
}

// Validate validates hook data.
func (v *HookValidator) Validate(data *Data) error {
	// Validate hook type
	if data.Type == "" {
		return NewValidationError("type", "", "required", "hook type is required", nil)
	}

	// Validate extension name
	if data.Extension == "" {
		return NewValidationError("extension", "", "required", "extension name is required", nil)
	}

	// Validate extension name format (alphanumeric, hyphens, underscores)
	if !isValidExtensionName(data.Extension) {
		return NewValidationError("extension", data.Extension, "format", "extension name must be alphanumeric with hyphens and underscores only", nil)
	}

	// Validate metadata if provided
	if data.Metadata != nil {
		if err := v.validateMetadata(data.Metadata); err != nil {
			return NewValidationError("metadata", "", "validation", "metadata validation failed", err)
		}
	}

	return nil
}

// validateMetadata validates hook metadata.
func (v *HookValidator) validateMetadata(metadata map[string]interface{}) error {
	for key, value := range metadata {
		// Validate key format
		if !isValidMetadataKey(key) {
			return fmt.Errorf("invalid metadata key: %s", key)
		}

		// Validate value type
		if !isValidMetadataValue(value) {
			return fmt.Errorf("invalid metadata value type for key %s: %T", key, value)
		}
	}

	return nil
}

// isValidExtensionName checks if the extension name is valid.
func isValidExtensionName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched && len(name) > 0 && len(name) <= 100
}

// isValidPath checks if the path is valid.
func isValidPath(path string) bool {
	// Basic path validation - should not be empty and should not contain invalid characters
	return len(path) > 0 && !strings.Contains(path, "..") && !strings.Contains(path, "\x00")
}

// isValidMetadataKey checks if the metadata key is valid.
func isValidMetadataKey(key string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, key)
	return matched && len(key) > 0 && len(key) <= 50
}

// isValidMetadataValue checks if the metadata value is valid.
func isValidMetadataValue(value interface{}) bool {
	switch value.(type) {
	case string, int, int64, float64, bool, []string, map[string]interface{}:
		return true
	default:
		return false
	}
}
