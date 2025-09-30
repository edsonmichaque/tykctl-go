package hook

import (
	"fmt"
)

// BuiltinValidator handles validation specific to builtin hooks.
type BuiltinValidator struct {
	// Builtin-specific validation configuration
}

// NewBuiltinValidator creates a new builtin hook validator.
func NewBuiltinValidator() *BuiltinValidator {
	return &BuiltinValidator{}
}

// Validate validates builtin hook data.
func (v *BuiltinValidator) Validate(data *Data) error {
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

	// Builtin-specific validations
	if err := v.validateBuiltinSpecific(data); err != nil {
		return err
	}

	return nil
}

// validateBuiltinSpecific performs builtin-specific validations.
func (v *BuiltinValidator) validateBuiltinSpecific(data *Data) error {
	// Validate hook type is a standard builtin type
	if !isValidBuiltinHookType(data.Type) {
		return NewValidationError("type", data.Type, "builtin", "invalid builtin hook type", nil)
	}

	return nil
}

// validateMetadata validates hook metadata.
func (v *BuiltinValidator) validateMetadata(metadata map[string]interface{}) error {
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

// isValidBuiltinHookType checks if the hook type is a valid builtin type.
func isValidBuiltinHookType(hookType Type) bool {
	validTypes := map[Type]bool{
		"before-install":   true,
		"after-install":    true,
		"before-run":       true,
		"after-run":        true,
		"before-uninstall": true,
		"after-uninstall":  true,
		"before-update":    true,
		"after-update":     true,
		// Auth hooks
		"before-login":     true,
		"after-login":      true,
		"before-logout":    true,
		"after-logout":     true,
	}
	return validTypes[hookType]
}
