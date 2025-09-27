package jsonschema

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

// Validator represents a JSON Schema validator
type Validator struct {
	schema *gojsonschema.Schema
}

// ValidationError represents a validation error
type ValidationError struct {
	Field       string `json:"field"`
	Description string `json:"description"`
	Context     string `json:"context"`
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// New creates a new validator from a schema string
func New(schema string) (*Validator, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	schemaObj, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &Validator{schema: schemaObj}, nil
}

// NewFromFile creates a new validator from a schema file
func NewFromFile(schemaPath string) (*Validator, error) {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	return New(string(schemaBytes))
}

// NewFromURL creates a new validator from a schema URL
func NewFromURL(schemaURL string) (*Validator, error) {
	resp, err := http.Get(schemaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schema from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch schema: HTTP %d", resp.StatusCode)
	}

	schemaBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema response: %w", err)
	}

	return New(string(schemaBytes))
}

// Validate validates JSON data against the schema
func (v *Validator) Validate(ctx context.Context, data []byte) (*ValidationResult, error) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	documentLoader := gojsonschema.NewBytesLoader(data)
	result, err := v.schema.Validate(documentLoader)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	validationResult := &ValidationResult{
		Valid:  result.Valid(),
		Errors: make([]ValidationError, 0),
	}

	if !result.Valid() {
		for _, err := range result.Errors() {
			validationResult.Errors = append(validationResult.Errors, ValidationError{
				Field:       err.Field(),
				Description: err.Description(),
				Context:     err.Context().String(),
			})
		}
	}

	return validationResult, nil
}

// ValidateString validates a JSON string against the schema
func (v *Validator) ValidateString(ctx context.Context, data string) (*ValidationResult, error) {
	return v.Validate(ctx, []byte(data))
}

// ValidateObject validates a Go object against the schema
func (v *Validator) ValidateObject(ctx context.Context, data interface{}) (*ValidationResult, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}

	return v.Validate(ctx, jsonData)
}

// ValidateFile validates a JSON file against the schema
func (v *Validator) ValidateFile(ctx context.Context, filePath string) (*ValidationResult, error) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return v.Validate(ctx, data)
}

// ValidateURL validates JSON data from a URL against the schema
func (v *Validator) ValidateURL(ctx context.Context, url string) (*ValidationResult, error) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return v.Validate(ctx, data)
}

// Convenience functions for common use cases

// Validate validates JSON data against a schema string
func Validate(ctx context.Context, data []byte, schema string) (*ValidationResult, error) {
	validator, err := New(schema)
	if err != nil {
		return nil, err
	}

	return validator.Validate(ctx, data)
}

// ValidateString validates a JSON string against a schema string
func ValidateString(ctx context.Context, data string, schema string) (*ValidationResult, error) {
	validator, err := New(schema)
	if err != nil {
		return nil, err
	}

	return validator.ValidateString(ctx, data)
}

// ValidateObject validates a Go object against a schema string
func ValidateObject(ctx context.Context, data interface{}, schema string) (*ValidationResult, error) {
	validator, err := New(schema)
	if err != nil {
		return nil, err
	}

	return validator.ValidateObject(ctx, data)
}

// ValidateFile validates a JSON file against a schema file
func ValidateFile(ctx context.Context, dataPath string, schemaPath string) (*ValidationResult, error) {
	validator, err := NewFromFile(schemaPath)
	if err != nil {
		return nil, err
	}

	return validator.ValidateFile(ctx, dataPath)
}

// ValidateURL validates JSON data from a URL against a schema URL
func ValidateURL(ctx context.Context, dataURL string, schemaURL string) (*ValidationResult, error) {
	validator, err := NewFromURL(schemaURL)
	if err != nil {
		return nil, err
	}

	return validator.ValidateURL(ctx, dataURL)
}

// IsValid checks if JSON data is valid against a schema (returns only boolean)
func IsValid(ctx context.Context, data []byte, schema string) (bool, error) {
	result, err := Validate(ctx, data, schema)
	if err != nil {
		return false, err
	}
	return result.Valid, nil
}

// IsValidString checks if a JSON string is valid against a schema (returns only boolean)
func IsValidString(ctx context.Context, data string, schema string) (bool, error) {
	result, err := ValidateString(ctx, data, schema)
	if err != nil {
		return false, err
	}
	return result.Valid, nil
}

// IsValidObject checks if a Go object is valid against a schema (returns only boolean)
func IsValidObject(ctx context.Context, data interface{}, schema string) (bool, error) {
	result, err := ValidateObject(ctx, data, schema)
	if err != nil {
		return false, err
	}
	return result.Valid, nil
}

// GetSchemaVersion extracts the JSON Schema version from a schema
func GetSchemaVersion(schema string) (string, error) {
	var schemaMap map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &schemaMap); err != nil {
		return "", fmt.Errorf("failed to parse schema: %w", err)
	}

	if version, ok := schemaMap["$schema"].(string); ok {
		return version, nil
	}

	return "", fmt.Errorf("schema version not found")
}

// GetSchemaTitle extracts the title from a schema
func GetSchemaTitle(schema string) (string, error) {
	var schemaMap map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &schemaMap); err != nil {
		return "", fmt.Errorf("failed to parse schema: %w", err)
	}

	if title, ok := schemaMap["title"].(string); ok {
		return title, nil
	}

	return "", fmt.Errorf("schema title not found")
}

// GetSchemaDescription extracts the description from a schema
func GetSchemaDescription(schema string) (string, error) {
	var schemaMap map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &schemaMap); err != nil {
		return "", fmt.Errorf("failed to parse schema: %w", err)
	}

	if description, ok := schemaMap["description"].(string); ok {
		return description, nil
	}

	return "", fmt.Errorf("schema description not found")
}

// ValidateDirectory validates all JSON files in a directory against a schema
func ValidateDirectory(ctx context.Context, dirPath string, schema string) (map[string]*ValidationResult, error) {
	validator, err := New(schema)
	if err != nil {
		return nil, err
	}

	results := make(map[string]*ValidationResult)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			result, err := validator.ValidateFile(ctx, path)
			if err != nil {
				return fmt.Errorf("failed to validate %s: %w", path, err)
			}
			results[path] = result
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}