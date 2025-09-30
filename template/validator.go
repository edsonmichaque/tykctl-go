package template

import (
	"fmt"
	"regexp"
)

// Validator defines the interface for validating templates.
type Validator interface {
	Validate(template *Template) error
}

// TemplateValidator is the concrete implementation of the Validator interface.
type TemplateValidator struct {
	// validation configuration
}

// NewValidator creates a new template validator.
func NewValidator() *TemplateValidator {
	return &TemplateValidator{}
}

// Validate validates a template.
func (v *TemplateValidator) Validate(template *Template) error {
	// Validate template name
	if template.Name == "" {
		return fmt.Errorf("template name is required")
	}

	// Validate resource type
	if template.ResourceType == "" {
		return fmt.Errorf("template resource_type is required")
	}

	// Validate content
	if len(template.Content) == 0 {
		return fmt.Errorf("template content is required")
	}

	// Validate variables
	for _, variable := range template.Variables {
		if err := v.validateVariable(variable); err != nil {
			return fmt.Errorf("variable %s validation failed: %w", variable.Name, err)
		}
	}

	return nil
}

// validateVariable validates a single variable.
func (v *TemplateValidator) validateVariable(variable Variable) error {
	// Validate variable name
	if variable.Name == "" {
		return fmt.Errorf("variable name is required")
	}

	// Validate variable type
	validTypes := []string{"string", "integer", "number", "boolean", "array", "object"}
	if variable.Type != "" {
		valid := false
		for _, validType := range validTypes {
			if variable.Type == validType {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid variable type: %s", variable.Type)
		}
	}

	// Validate validation rules
	if err := v.validateValidationRules(variable.Validation); err != nil {
		return fmt.Errorf("validation rules failed for variable %s: %w", variable.Name, err)
	}

	return nil
}

// validateValidationRules validates variable validation rules.
func (v *TemplateValidator) validateValidationRules(validation Validation) error {
	// Validate min/max length
	if validation.MinLength != nil && validation.MaxLength != nil {
		if *validation.MinLength > *validation.MaxLength {
			return fmt.Errorf("min_length cannot be greater than max_length")
		}
	}

	// Validate min/max value
	if validation.MinValue != nil && validation.MaxValue != nil {
		if *validation.MinValue > *validation.MaxValue {
			return fmt.Errorf("min_value cannot be greater than max_value")
		}
	}

	// Validate pattern
	if validation.Pattern != "" {
		if _, err := regexp.Compile(validation.Pattern); err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	}

	// Validate enum
	if len(validation.Enum) == 0 && validation.Pattern == "" {
		// No validation rules specified, this is okay
		return nil
	}

	return nil
}


