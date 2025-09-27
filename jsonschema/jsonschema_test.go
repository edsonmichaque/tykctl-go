package jsonschema

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const validSchema = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"name": {
			"type": "string",
			"minLength": 1
		},
		"age": {
			"type": "integer",
			"minimum": 0,
			"maximum": 150
		},
		"email": {
			"type": "string",
			"format": "email"
		}
	},
	"required": ["name", "age"]
}`

const invalidSchema = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
		"name": {
			"type": "invalid_type"
		}
	}
}`

func TestNew(t *testing.T) {
	validator, err := New(validSchema)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if validator == nil {
		t.Fatal("Validator is nil")
	}
}

func TestNewInvalidSchema(t *testing.T) {
	_, err := New(invalidSchema)
	if err == nil {
		t.Error("Expected error for invalid schema")
	}
}

func TestValidateValidData(t *testing.T) {
	validator, err := New(validSchema)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	validData := `{"name": "John", "age": 30, "email": "john@example.com"}`
	result, err := validator.ValidateString(context.Background(), validData)
	if err != nil {
		t.Fatalf("ValidateString failed: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid, got invalid. Errors: %v", result.Errors)
	}
}

func TestValidateInvalidData(t *testing.T) {
	validator, err := New(validSchema)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	invalidData := `{"name": "", "age": -5}`
	result, err := validator.ValidateString(context.Background(), invalidData)
	if err != nil {
		t.Fatalf("ValidateString failed: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid, got valid")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestValidateObject(t *testing.T) {
	validator, err := New(validSchema)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	data := map[string]interface{}{
		"name":  "Jane",
		"age":   25,
		"email": "jane@example.com",
	}

	result, err := validator.ValidateObject(context.Background(), data)
	if err != nil {
		t.Fatalf("ValidateObject failed: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid, got invalid. Errors: %v", result.Errors)
	}
}

func TestValidateFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write valid JSON to file
	validData := `{"name": "Test", "age": 30}`
	if _, err := tmpFile.WriteString(validData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	validator, err := New(validSchema)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	result, err := validator.ValidateFile(context.Background(), tmpFile.Name())
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected valid, got invalid. Errors: %v", result.Errors)
	}
}

func TestContextCancellation(t *testing.T) {
	validator, err := New(validSchema)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	validData := `{"name": "John", "age": 30}`
	_, err = validator.ValidateString(ctx, validData)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestContextTimeout(t *testing.T) {
	validator, err := New(validSchema)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(1 * time.Millisecond)

	validData := `{"name": "John", "age": 30}`
	_, err = validator.ValidateString(ctx, validData)
	if err == nil {
		t.Error("Expected error for timed out context")
	}

	if err != context.DeadlineExceeded {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	validData := `{"name": "Test", "age": 30}`

	// Test Validate function
	result, err := Validate(context.Background(), []byte(validData), validSchema)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if !result.Valid {
		t.Error("Expected valid result")
	}

	// Test ValidateString function
	result, err = ValidateString(context.Background(), validData, validSchema)
	if err != nil {
		t.Fatalf("ValidateString failed: %v", err)
	}
	if !result.Valid {
		t.Error("Expected valid result")
	}

	// Test ValidateObject function
	data := map[string]interface{}{
		"name": "Test",
		"age":  30,
	}
	result, err = ValidateObject(context.Background(), data, validSchema)
	if err != nil {
		t.Fatalf("ValidateObject failed: %v", err)
	}
	if !result.Valid {
		t.Error("Expected valid result")
	}
}

func TestIsValidFunctions(t *testing.T) {
	validData := `{"name": "Test", "age": 30}`
	invalidData := `{"name": "", "age": -5}`

	// Test IsValid
	valid, err := IsValid(context.Background(), []byte(validData), validSchema)
	if err != nil {
		t.Fatalf("IsValid failed: %v", err)
	}
	if !valid {
		t.Error("Expected valid")
	}

	valid, err = IsValid(context.Background(), []byte(invalidData), validSchema)
	if err != nil {
		t.Fatalf("IsValid failed: %v", err)
	}
	if valid {
		t.Error("Expected invalid")
	}

	// Test IsValidString
	valid, err = IsValidString(context.Background(), validData, validSchema)
	if err != nil {
		t.Fatalf("IsValidString failed: %v", err)
	}
	if !valid {
		t.Error("Expected valid")
	}

	// Test IsValidObject
	data := map[string]interface{}{
		"name": "Test",
		"age":  30,
	}
	valid, err = IsValidObject(context.Background(), data, validSchema)
	if err != nil {
		t.Fatalf("IsValidObject failed: %v", err)
	}
	if !valid {
		t.Error("Expected valid")
	}
}

func TestGetSchemaVersion(t *testing.T) {
	version, err := GetSchemaVersion(validSchema)
	if err != nil {
		t.Fatalf("GetSchemaVersion failed: %v", err)
	}

	expected := "http://json-schema.org/draft-07/schema#"
	if version != expected {
		t.Errorf("Expected %s, got %s", expected, version)
	}
}

func TestGetSchemaTitle(t *testing.T) {
	schemaWithTitle := `{
		"title": "Test Schema",
		"type": "object"
	}`

	title, err := GetSchemaTitle(schemaWithTitle)
	if err != nil {
		t.Fatalf("GetSchemaTitle failed: %v", err)
	}

	expected := "Test Schema"
	if title != expected {
		t.Errorf("Expected %s, got %s", expected, title)
	}
}

func TestGetSchemaDescription(t *testing.T) {
	schemaWithDescription := `{
		"description": "A test schema",
		"type": "object"
	}`

	description, err := GetSchemaDescription(schemaWithDescription)
	if err != nil {
		t.Fatalf("GetSchemaDescription failed: %v", err)
	}

	expected := "A test schema"
	if description != expected {
		t.Errorf("Expected %s, got %s", expected, description)
	}
}

func TestValidateDirectory(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test JSON files
	files := []string{
		`{"name": "Test1", "age": 30}`,
		`{"name": "Test2", "age": 25}`,
		`{"name": "", "age": -5}`, // Invalid
	}

	for i, content := range files {
		filename := filepath.Join(tmpDir, fmt.Sprintf("test%d.json", i+1))
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	results, err := ValidateDirectory(context.Background(), tmpDir, validSchema)
	if err != nil {
		t.Fatalf("ValidateDirectory failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Check that we have both valid and invalid results
	validCount := 0
	invalidCount := 0
	for _, result := range results {
		if result.Valid {
			validCount++
		} else {
			invalidCount++
		}
	}

	if validCount != 2 {
		t.Errorf("Expected 2 valid files, got %d", validCount)
	}
	if invalidCount != 1 {
		t.Errorf("Expected 1 invalid file, got %d", invalidCount)
	}
}

func TestValidateDirectoryWithCancellation(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test JSON files
	files := []string{
		`{"name": "Test1", "age": 30}`,
		`{"name": "Test2", "age": 25}`,
		`{"name": "Test3", "age": 35}`,
	}

	for i, content := range files {
		filename := filepath.Join(tmpDir, fmt.Sprintf("test%d.json", i+1))
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = ValidateDirectory(ctx, tmpDir, validSchema)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestValidationError(t *testing.T) {
	validator, err := New(validSchema)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	invalidData := `{"name": "", "age": -5}`
	result, err := validator.ValidateString(context.Background(), invalidData)
	if err != nil {
		t.Fatalf("ValidateString failed: %v", err)
	}

	if result.Valid {
		t.Error("Expected invalid result")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}

	// Check error structure
	for _, err := range result.Errors {
		if err.Field == "" {
			t.Error("Expected non-empty field")
		}
		if err.Description == "" {
			t.Error("Expected non-empty description")
		}
	}
}