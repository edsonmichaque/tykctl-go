package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Rule    string      `json:"rule"`
	Message string      `json:"message"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors"`
	Warnings []ValidationError `json:"warnings"`
}

// StructValidator provides validation capabilities for structs
type StructValidator interface {
	Validate(ctx context.Context, value interface{}) error
	GetName() string
	GetDescription() string
}

// Built-in validators
type RequiredValidator struct{}

func (v RequiredValidator) Validate(ctx context.Context, value interface{}) error {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.String && value == "") {
		return ValidationError{Message: "field is required"}
	}
	return nil
}

func (v RequiredValidator) GetName() string        { return "required" }
func (v RequiredValidator) GetDescription() string { return "Field is required" }

type MinLengthValidator struct {
	Min int
}

func (v MinLengthValidator) Validate(ctx context.Context, value interface{}) error {
	if str, ok := value.(string); ok && len(str) < v.Min {
		return ValidationError{Message: fmt.Sprintf("minimum length is %d", v.Min)}
	}
	return nil
}

func (v MinLengthValidator) GetName() string        { return "min_length" }
func (v MinLengthValidator) GetDescription() string { return fmt.Sprintf("Minimum length is %d", v.Min) }

type MaxLengthValidator struct {
	Max int
}

func (v MaxLengthValidator) Validate(ctx context.Context, value interface{}) error {
	if str, ok := value.(string); ok && len(str) > v.Max {
		return ValidationError{Message: fmt.Sprintf("maximum length is %d", v.Max)}
	}
	return nil
}

func (v MaxLengthValidator) GetName() string        { return "max_length" }
func (v MaxLengthValidator) GetDescription() string { return fmt.Sprintf("Maximum length is %d", v.Max) }

type RangeValidator struct {
	Min, Max int
}

func (v RangeValidator) Validate(ctx context.Context, value interface{}) error {
	if num, ok := value.(int); ok && (num < v.Min || num > v.Max) {
		return ValidationError{Message: fmt.Sprintf("value must be between %d and %d", v.Min, v.Max)}
	}
	return nil
}

func (v RangeValidator) GetName() string        { return "range" }
func (v RangeValidator) GetDescription() string { return fmt.Sprintf("Value must be between %d and %d", v.Min, v.Max) }

type URLValidator struct{}

func (v URLValidator) Validate(ctx context.Context, value interface{}) error {
	if str, ok := value.(string); ok && str != "" {
		if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
			return ValidationError{Message: "must be a valid URL"}
		}
	}
	return nil
}

func (v URLValidator) GetName() string        { return "url" }
func (v URLValidator) GetDescription() string { return "Must be a valid URL" }

// ValidatorRegistry manages validators
type ValidatorRegistry struct {
	validators map[string]StructValidator
}

func NewValidatorRegistry() *ValidatorRegistry {
	registry := &ValidatorRegistry{
		validators: make(map[string]StructValidator),
	}

	// Register built-in validators
	registry.Register(RequiredValidator{})
	registry.Register(MinLengthValidator{Min: 1})
	registry.Register(MaxLengthValidator{Max: 255})
	registry.Register(RangeValidator{Min: 0, Max: 100})
	registry.Register(URLValidator{})

	return registry
}

func (r *ValidatorRegistry) Register(validator StructValidator) {
	r.validators[validator.GetName()] = validator
}

func (r *ValidatorRegistry) Get(name string) (StructValidator, bool) {
	validator, exists := r.validators[name]
	return validator, exists
}

// ValidateStruct validates a struct using tags
func ValidateStruct(ctx context.Context, target interface{}) ValidationResult {
	result := ValidationResult{Valid: true}

	v := reflect.ValueOf(target)
	t := reflect.TypeOf(target)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get validation tags
		tags := fieldType.Tag.Get("validate")
		if tags == "" {
			continue
		}

		// Parse validation rules
		rules := strings.Split(tags, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)

			// Apply validation
			if err := applyValidationRule(ctx, rule, field.Interface()); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   fieldType.Name,
					Value:   field.Interface(),
					Rule:    rule,
					Message: err.Error(),
				})
			}
		}
	}

	return result
}

func applyValidationRule(ctx context.Context, rule string, value interface{}) error {
	// Simple implementation - in a real implementation, this would be more sophisticated
	registry := NewValidatorRegistry()
	
	if validator, exists := registry.Get(rule); exists {
		return validator.Validate(ctx, value)
	}
	
	return fmt.Errorf("unknown validation rule: %s", rule)
}