package template

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCompleteWorkflow(t *testing.T) {
	// Create a test template file
	tempFile := "test_workflow.yaml"
	content := `name: "Test Template"
resource_type: "product"
content:
  name: "{{.Name}}"
  description: "{{.Description}}"
variables:
  - name: "Name"
    type: "string"
    required: true
    validation:
      min_length: 1
      max_length: 100
  - name: "Description"
    type: "string"
    required: false
`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	// Test complete workflow
	ctx := context.Background()

	// 1. Create loader
	loader := NewFileLoader(tempFile)

	// 2. Load template
	tmpl, err := loader.Load(ctx)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// 3. Validate template
	validator := NewValidator()
	err = validator.Validate(tmpl)
	if err != nil {
		t.Fatalf("Template validation failed: %v", err)
	}

	// 4. Process template
	processor := NewProcessor(loader, validator)
	variables := map[string]interface{}{
		"Name":        "Test Product",
		"Description": "A test product",
	}

	result, err := processor.Process(ctx, variables)
	if err != nil {
		t.Fatalf("Template processing failed: %v", err)
	}

	// 5. Verify result
	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}

	t.Logf("Processed template result: %s", string(result))
}

func TestFileLoaderWorkflow(t *testing.T) {
	// Test file loader creation
	loader := NewFileLoader("nonexistent.yaml")
	if loader == nil {
		t.Error("Expected non-nil file loader")
	}

	// Test loading non-existent file
	ctx := context.Background()
	_, err := loader.Load(ctx)
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestDirLoaderWorkflow(t *testing.T) {
	// Create test directory
	testDir := "test_dir_workflow"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create test template
	templateContent := `name: "Dir Template"
resource_type: "product"
content:
  name: "{{.Name}}"
  description: "{{.Description}}"
variables:
  - name: "Name"
    type: "string"
    required: true
  - name: "Description"
    type: "string"
    required: false`

	err = os.WriteFile(filepath.Join(testDir, "product.yaml"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Test dir loader creation
	loader := NewDirLoader(testDir, "product", nil)
	if loader == nil {
		t.Error("Expected non-nil dir loader")
	}

	// Test loading template
	ctx := context.Background()
	template, err := loader.Load(ctx)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	if template.Name != "Dir Template" {
		t.Errorf("Expected name 'Dir Template', got '%s'", template.Name)
	}

	// Test listing templates
	templates, err := loader.ListTemplates()
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	if len(templates) != 1 || templates[0] != "product" {
		t.Errorf("Expected ['product'], got %v", templates)
	}
}

func TestProcess(t *testing.T) {
	// Create a test template file
	tempFile := "test_process_from_source.yaml"
	content := `name: "Source Template"
resource_type: "product"
content:
  name: "{{.Name}}"
  description: "{{.Description}}"
variables:
  - name: "Name"
    type: "string"
    required: true
  - name: "Description"
    type: "string"
    required: false
`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	// Test Process
	ctx := context.Background()
	loader := NewFileLoader(tempFile)

	validator := NewValidator()
	processor := NewProcessor(loader, validator)

	variables := map[string]interface{}{
		"Name":        "Source Product",
		"Description": "A product from source",
	}

	result, err := processor.Process(ctx, variables)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if len(result) == 0 {
		t.Error("Expected non-empty result")
	}

	t.Logf("Process result: %s", string(result))
}

func TestTemplateValidationWorkflow(t *testing.T) {
	// Test valid template
	validTemplate := &Template{
		Name:         "Test Template",
		ResourceType: "product",
		Content: map[string]interface{}{
			"name": "{{.Name}}",
		},
		Variables: []Variable{
			{
				Name:     "Name",
				Type:     "string",
				Required: true,
			},
		},
	}

	validator := NewValidator()
	err := validator.Validate(validTemplate)
	if err != nil {
		t.Errorf("Expected valid template to pass validation, got error: %v", err)
	}

	// Test invalid template (missing name)
	invalidTemplate := &Template{
		ResourceType: "product",
		Content: map[string]interface{}{
			"name": "{{.Name}}",
		},
	}

	err = validator.Validate(invalidTemplate)
	if err == nil {
		t.Error("Expected invalid template to fail validation")
	}

	// Test invalid template (missing resource type)
	invalidTemplate2 := &Template{
		Name: "Test Template",
		Content: map[string]interface{}{
			"name": "{{.Name}}",
		},
	}

	err = validator.Validate(invalidTemplate2)
	if err == nil {
		t.Error("Expected invalid template to fail validation")
	}

	// Test invalid template (missing content)
	invalidTemplate3 := &Template{
		Name:         "Test Template",
		ResourceType: "product",
	}

	err = validator.Validate(invalidTemplate3)
	if err == nil {
		t.Error("Expected invalid template to fail validation")
	}
}
