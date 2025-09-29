package template

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func ExampleFileLoader() {
	// Create a file loader
	loader := NewFileLoader("template.yaml")

	// Load the template
	ctx := context.Background()
	template, err := loader.Load(ctx)
	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
		return
	}

	fmt.Printf("Loaded template: %s\n", template.Name)
}

func ExampleProcessor_process() {
	// Create a template
	template := &Template{
		Name:         "example",
		ResourceType: "product",
		Content: map[string]interface{}{
			"name":        "{{.Name}}",
			"description": "{{.Description}}",
		},
	}

	// Create processor with mock loader and validator
	loader := NewFileLoader("dummy.yaml")
	validator := NewValidator()
	processor := NewProcessor(loader, validator)

	// Process template with variables
	variables := map[string]interface{}{
		"Name":        "My Product",
		"Description": "A great product",
	}

	ctx := context.Background()
	result, err := processor.process(ctx, template, variables)
	if err != nil {
		fmt.Printf("Error processing template: %v\n", err)
		return
	}

	fmt.Printf("Processed result: %s\n", string(result))
}

func ExampleTemplateValidator_Validate() {
	// Create a template
	template := &Template{
		Name:         "example",
		ResourceType: "product",
		Content: map[string]interface{}{
			"name": "test",
		},
		Variables: []Variable{
			{
				Name:     "Name",
				Type:     "string",
				Required: true,
				Validation: Validation{
					MinLength: intPtr(1),
					MaxLength: intPtr(100),
				},
			},
		},
	}

	// Create validator
	validator := NewValidator()

	// Validate template
	err := validator.Validate(template)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	fmt.Println("Template is valid")
}

func TestTemplateValidation(t *testing.T) {
	tests := []struct {
		name     string
		template *Template
		wantErr  bool
	}{
		{
			name: "valid template",
			template: &Template{
				Name:         "test",
				ResourceType: "product",
				Content:      map[string]interface{}{"name": "test"},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			template: &Template{
				ResourceType: "product",
				Content:      map[string]interface{}{"name": "test"},
			},
			wantErr: true,
		},
		{
			name: "missing resource type",
			template: &Template{
				Name:    "test",
				Content: map[string]interface{}{"name": "test"},
			},
			wantErr: true,
		},
		{
			name: "missing content",
			template: &Template{
				Name:         "test",
				ResourceType: "product",
			},
			wantErr: true,
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileLoader(t *testing.T) {
	// Create a temporary template file
	tempFile := "test_template.yaml"
	content := `name: "test template"
resource_type: "product"
content:
  name: "{{.Name}}"
  description: "{{.Description}}"
`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile)

	// Test file loader
	loader := NewFileLoader(tempFile)
	ctx := context.Background()

	template, err := loader.Load(ctx)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	if template.Name != "test template" {
		t.Errorf("Expected name 'test template', got '%s'", template.Name)
	}

	if template.ResourceType != "product" {
		t.Errorf("Expected resource_type 'product', got '%s'", template.ResourceType)
	}
}

// Helper function for tests
func intPtr(i int) *int {
	return &i
}
