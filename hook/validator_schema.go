package hook

import (
	"fmt"
	"regexp"
)

// SchemaValidator handles validation specific to JSON Schema validation hooks.
type SchemaValidator struct {
	schemaDir string
}

// NewSchemaValidator creates a new JSON Schema hook validator.
func NewSchemaValidator(schemaDir string) *SchemaValidator {
	return &SchemaValidator{
		schemaDir: schemaDir,
	}
}

// Validate validates JSON Schema hook data.
func (v *SchemaValidator) Validate(data *Data) error {
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

	// JSON Schema-specific validations
	if err := v.validateJSONSchemaSpecific(data); err != nil {
		return err
	}

	return nil
}

// validateJSONSchemaSpecific performs JSON Schema-specific validations.
func (v *SchemaValidator) validateJSONSchemaSpecific(data *Data) error {
	// Validate hook type is appropriate for JSON schemas
	if !isValidJSONSchemaHookType(data.Type) {
		return NewValidationError("type", data.Type, "schema", "invalid JSON schema hook type", nil)
	}

	// Validate JSON Schema-specific metadata
	if err := v.validateJSONSchemaMetadata(data.Metadata); err != nil {
		return NewValidationError("metadata", "", "schema", err.Error(), nil)
	}

	return nil
}

// validateJSONSchemaMetadata validates JSON Schema-specific metadata.
func (v *SchemaValidator) validateJSONSchemaMetadata(metadata map[string]interface{}) error {
	if metadata == nil {
		return nil
	}

	// Check for required JSON Schema metadata
	requiredKeys := []string{"schema_name", "version"}
	for _, key := range requiredKeys {
		if _, exists := metadata[key]; !exists {
			return fmt.Errorf("JSON schema metadata must contain key: %s", key)
		}
	}

	// Validate schema name format
	if schemaName, exists := metadata["schema_name"]; exists {
		if nameStr, ok := schemaName.(string); ok {
			if !isValidJSONSchemaName(nameStr) {
				return fmt.Errorf("invalid JSON schema name format: %s", nameStr)
			}
		} else {
			return fmt.Errorf("schema_name must be a string")
		}
	}

	// Validate version format
	if version, exists := metadata["version"]; exists {
		if versionStr, ok := version.(string); ok {
			if !isValidSemanticVersion(versionStr) {
				return fmt.Errorf("invalid version format: %s", versionStr)
			}
		} else {
			return fmt.Errorf("version must be a string")
		}
	}

	// Validate JSON Schema type if provided
	if schemaType, exists := metadata["schema_type"]; exists {
		if typeStr, ok := schemaType.(string); ok {
			if !isValidJSONSchemaType(typeStr) {
				return fmt.Errorf("invalid JSON schema type: %s", typeStr)
			}
		} else {
			return fmt.Errorf("schema_type must be a string")
		}
	}

	return nil
}

// validateMetadata validates hook metadata.
func (v *SchemaValidator) validateMetadata(metadata map[string]interface{}) error {
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

// isValidJSONSchemaHookType checks if the hook type is valid for JSON schemas.
func isValidJSONSchemaHookType(hookType Type) bool {
	validTypes := map[Type]bool{
		"before-install":   true,
		"after-install":    true,
		"before-run":       true,
		"after-run":        true,
		"before-uninstall": true,
		"after-uninstall":  true,
		"before-update":    true,
		"after-update":     true,
		"validation":       true,
		"pre-commit":       true,
		"post-commit":      true,
		"pre-push":         true,
		"post-push":        true,
	}
	return validTypes[hookType]
}

// isValidJSONSchemaName checks if the JSON schema name is valid.
func isValidJSONSchemaName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched && len(name) > 0 && len(name) <= 100
}

// isValidJSONSchemaType checks if the JSON schema type is valid.
func isValidJSONSchemaType(schemaType string) bool {
	validTypes := []string{
		"object",
		"array",
		"string",
		"number",
		"integer",
		"boolean",
		"null",
	}

	for _, validType := range validTypes {
		if schemaType == validType {
			return true
		}
	}
	return false
}

// isValidSemanticVersion checks if the version string is a valid semantic version.
func isValidSemanticVersion(version string) bool {
	// Basic semantic version validation
	matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+(-[a-zA-Z0-9-]+)?(\+[a-zA-Z0-9-]+)?$`, version)
	return matched && len(version) > 0 && len(version) <= 20
}
