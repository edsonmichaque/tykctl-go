package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScriptValidator handles validation specific to script hooks.
type ScriptValidator struct {
	scriptDir string
}

// NewScriptValidator creates a new script hook validator.
func NewScriptValidator(scriptDir string) *ScriptValidator {
	return &ScriptValidator{
		scriptDir: scriptDir,
	}
}

// Validate validates script hook data.
func (v *ScriptValidator) Validate(data *Data) error {
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

	// Script-specific validations
	if err := v.validateScriptSpecific(data); err != nil {
		return err
	}

	return nil
}

// validateScriptSpecific performs script-specific validations.
func (v *ScriptValidator) validateScriptSpecific(data *Data) error {
	// Validate hook type is appropriate for scripts
	if !isValidScriptHookType(data.Type) {
		return NewValidationError("type", data.Type, "script", "invalid script hook type", nil)
	}

	return nil
}

// validateScriptFile validates that the script file exists and is executable.
func (v *ScriptValidator) validateScriptFile(scriptPath string) error {
	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script file does not exist: %s", scriptPath)
	}

	// Check if it's a valid script file extension
	ext := strings.ToLower(filepath.Ext(scriptPath))
	validExtensions := []string{".sh", ".bash", ".py", ".js", ".ts", ".go", ".rb", ".pl", ".ps1", ".bat", ".cmd"}

	isValidExt := false
	for _, validExt := range validExtensions {
		if ext == validExt {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return fmt.Errorf("unsupported script file extension: %s", ext)
	}

	// Check if script is within the allowed directory
	if v.scriptDir != "" {
		absScriptPath, err := filepath.Abs(scriptPath)
		if err != nil {
			return fmt.Errorf("invalid script path: %v", err)
		}

		absScriptDir, err := filepath.Abs(v.scriptDir)
		if err != nil {
			return fmt.Errorf("invalid script directory: %v", err)
		}

		if !strings.HasPrefix(absScriptPath, absScriptDir) {
			return fmt.Errorf("script path must be within script directory: %s", v.scriptDir)
		}
	}

	return nil
}

// validateMetadata validates hook metadata.
func (v *ScriptValidator) validateMetadata(metadata map[string]interface{}) error {
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

// isValidScriptHookType checks if the hook type is valid for scripts.
func isValidScriptHookType(hookType Type) bool {
	validTypes := map[Type]bool{
		"before-install":   true,
		"after-install":    true,
		"before-run":       true,
		"after-run":        true,
		"before-uninstall": true,
		"after-uninstall":  true,
		"before-update":    true,
		"after-update":     true,
		"pre-commit":       true,
		"post-commit":      true,
		"pre-push":         true,
		"post-push":        true,
	}
	return validTypes[hookType]
}
