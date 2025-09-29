package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PolicyValidator handles validation specific to Rego policy hooks.
type PolicyValidator struct {
	policyDir string
}

// NewPolicyValidator creates a new policy hook validator.
func NewPolicyValidator(policyDir string) *PolicyValidator {
	return &PolicyValidator{
		policyDir: policyDir,
	}
}

// Validate validates policy hook data.
func (v *PolicyValidator) Validate(data *Data) error {
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

	// Policy-specific validations
	if err := v.validatePolicySpecific(data); err != nil {
		return err
	}

	return nil
}

// validatePolicySpecific performs policy-specific validations.
func (v *PolicyValidator) validatePolicySpecific(data *Data) error {
	// Validate hook type is appropriate for policies
	if !isValidPolicyHookType(data.Type) {
		return NewValidationError("type", data.Type, "policy", "invalid policy hook type", nil)
	}

	// Validate policy-specific metadata
	if err := v.validatePolicyMetadata(data.Metadata); err != nil {
		return NewValidationError("metadata", "", "policy", err.Error(), nil)
	}

	return nil
}

// validatePolicyFile validates that the policy file exists and is valid.
func (v *PolicyValidator) validatePolicyFile(policyPath string) error {
	// Check if policy file exists
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		return fmt.Errorf("policy file does not exist: %s", policyPath)
	}

	// Check if it's a valid policy file extension
	ext := strings.ToLower(filepath.Ext(policyPath))
	if ext != ".rego" {
		return fmt.Errorf("policy file must have .rego extension, got: %s", ext)
	}

	// Check if policy is within the allowed directory
	if v.policyDir != "" {
		absPolicyPath, err := filepath.Abs(policyPath)
		if err != nil {
			return fmt.Errorf("invalid policy path: %v", err)
		}

		absPolicyDir, err := filepath.Abs(v.policyDir)
		if err != nil {
			return fmt.Errorf("invalid policy directory: %v", err)
		}

		if !strings.HasPrefix(absPolicyPath, absPolicyDir) {
			return fmt.Errorf("policy path must be within policy directory: %s", v.policyDir)
		}
	}

	// Basic Rego syntax validation
	if err := v.validateRegoSyntax(policyPath); err != nil {
		return fmt.Errorf("invalid Rego syntax: %v", err)
	}

	return nil
}

// validateRegoSyntax performs basic Rego syntax validation.
func (v *PolicyValidator) validateRegoSyntax(policyPath string) error {
	content, err := os.ReadFile(policyPath)
	if err != nil {
		return fmt.Errorf("failed to read policy file: %v", err)
	}

	// Basic checks for Rego syntax
	contentStr := string(content)

	// Check for package declaration
	if !strings.Contains(contentStr, "package ") {
		return fmt.Errorf("policy file must contain a package declaration")
	}

	// Check for basic Rego keywords
	regoKeywords := []string{"allow", "deny", "input", "data"}
	hasKeyword := false
	for _, keyword := range regoKeywords {
		if strings.Contains(contentStr, keyword) {
			hasKeyword = true
			break
		}
	}

	if !hasKeyword {
		return fmt.Errorf("policy file must contain at least one Rego keyword (allow, deny, input, data)")
	}

	return nil
}

// validatePolicyMetadata validates policy-specific metadata.
func (v *PolicyValidator) validatePolicyMetadata(metadata map[string]interface{}) error {
	if metadata == nil {
		return nil
	}

	// Check for required policy metadata
	requiredKeys := []string{"policy_name", "version"}
	for _, key := range requiredKeys {
		if _, exists := metadata[key]; !exists {
			return fmt.Errorf("policy metadata must contain key: %s", key)
		}
	}

	// Validate policy name format
	if policyName, exists := metadata["policy_name"]; exists {
		if nameStr, ok := policyName.(string); ok {
			if !isValidPolicyName(nameStr) {
				return fmt.Errorf("invalid policy name format: %s", nameStr)
			}
		} else {
			return fmt.Errorf("policy_name must be a string")
		}
	}

	// Validate version format
	if version, exists := metadata["version"]; exists {
		if versionStr, ok := version.(string); ok {
			if !isValidVersion(versionStr) {
				return fmt.Errorf("invalid version format: %s", versionStr)
			}
		} else {
			return fmt.Errorf("version must be a string")
		}
	}

	return nil
}

// validateMetadata validates hook metadata.
func (v *PolicyValidator) validateMetadata(metadata map[string]interface{}) error {
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

// isValidPolicyHookType checks if the hook type is valid for policies.
func isValidPolicyHookType(hookType Type) bool {
	validTypes := map[Type]bool{
		"before-install":   true,
		"after-install":    true,
		"before-run":       true,
		"after-run":        true,
		"before-uninstall": true,
		"after-uninstall":  true,
		"before-update":    true,
		"after-update":     true,
		"authorization":    true,
		"validation":       true,
		"audit":            true,
	}
	return validTypes[hookType]
}

// isValidPolicyName checks if the policy name is valid.
func isValidPolicyName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched && len(name) > 0 && len(name) <= 100
}

// isValidVersion checks if the version string is valid.
func isValidVersion(version string) bool {
	// Basic semantic version validation
	matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+(-[a-zA-Z0-9-]+)?(\+[a-zA-Z0-9-]+)?$`, version)
	return matched && len(version) > 0 && len(version) <= 20
}
