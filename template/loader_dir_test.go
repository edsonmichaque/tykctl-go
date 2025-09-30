package template

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDirLoader(t *testing.T) {
	// Create test directory structure
	testDir := "test_templates"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create test template files
	templates := map[string]string{
		"product.yaml": `name: "Product Template"
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
    required: false`,

		"api.yml": `name: "API Template"
resource_type: "api"
content:
  endpoint: "{{.Endpoint}}"
  method: "{{.Method}}"
variables:
  - name: "Endpoint"
    type: "string"
    required: true
  - name: "Method"
    type: "string"
    required: true`,

		"user/template.yaml": `name: "User Template"
resource_type: "user"
content:
  username: "{{.Username}}"
  email: "{{.Email}}"
variables:
  - name: "Username"
    type: "string"
    required: true
  - name: "Email"
    type: "string"
    required: true`,
	}

	// Write template files
	for path, content := range templates {
		fullPath := filepath.Join(testDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write template %s: %v", path, err)
		}
	}

	t.Run("LoadTemplate", func(t *testing.T) {
		loader := NewDirLoader(testDir, "product", nil)
		ctx := context.Background()

		template, err := loader.Load(ctx)
		if err != nil {
			t.Fatalf("Failed to load template: %v", err)
		}

		if template.Name != "Product Template" {
			t.Errorf("Expected name 'Product Template', got '%s'", template.Name)
		}

		if template.ResourceType != "product" {
			t.Errorf("Expected resource_type 'product', got '%s'", template.ResourceType)
		}
	})

	t.Run("LoadTemplateWithYmlExtension", func(t *testing.T) {
		loader := NewDirLoader(testDir, "api", nil)
		ctx := context.Background()

		template, err := loader.Load(ctx)
		if err != nil {
			t.Fatalf("Failed to load template: %v", err)
		}

		if template.Name != "API Template" {
			t.Errorf("Expected name 'API Template', got '%s'", template.Name)
		}
	})

	t.Run("LoadTemplateFromSubdirectory", func(t *testing.T) {
		loader := NewDirLoader(testDir, "user", nil)
		ctx := context.Background()

		template, err := loader.Load(ctx)
		if err != nil {
			t.Fatalf("Failed to load template: %v", err)
		}

		if template.Name != "User Template" {
			t.Errorf("Expected name 'User Template', got '%s'", template.Name)
		}
	})

	t.Run("ListTemplates", func(t *testing.T) {
		loader := NewDirLoader(testDir, "", nil)

		templates, err := loader.ListTemplates()
		if err != nil {
			t.Fatalf("Failed to list templates: %v", err)
		}

		expected := []string{"api", "product", "user"}
		if len(templates) != len(expected) {
			t.Errorf("Expected %d templates, got %d", len(expected), len(templates))
		}

		// Check that all expected templates are present
		for _, exp := range expected {
			found := false
			for _, tmpl := range templates {
				if tmpl == exp {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected template '%s' not found in list", exp)
			}
		}
	})

	t.Run("GetTemplatePath", func(t *testing.T) {
		loader := NewDirLoader(testDir, "product", nil)

		path, err := loader.GetTemplatePath("product")
		if err != nil {
			t.Fatalf("Failed to get template path: %v", err)
		}

		expectedPath := filepath.Join(testDir, "product.yaml")
		if path != expectedPath {
			t.Errorf("Expected path '%s', got '%s'", expectedPath, path)
		}
	})

	t.Run("TemplateNotFound", func(t *testing.T) {
		loader := NewDirLoader(testDir, "nonexistent", nil)
		ctx := context.Background()

		_, err := loader.Load(ctx)
		if err == nil {
			t.Error("Expected error for nonexistent template")
		}
	})
}

func TestCustomResolver(t *testing.T) {
	testDir := "test_custom_resolver"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create a template with custom naming
	templateContent := `name: "Custom Template"
resource_type: "custom"
content:
  value: "{{.Value}}"
variables:
  - name: "Value"
    type: "string"
    required: true`

	err = os.WriteFile(filepath.Join(testDir, "custom.tmpl"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	// Create custom resolver
	resolver := NewCustomResolver(func(dir, name string) (string, error) {
		path := filepath.Join(dir, name+".tmpl")
		if _, err := os.Stat(path); err != nil {
			return "", err
		}
		return path, nil
	})

	loader := NewDirLoader(testDir, "custom", resolver)
	ctx := context.Background()

	template, err := loader.Load(ctx)
	if err != nil {
		t.Fatalf("Failed to load template with custom resolver: %v", err)
	}

	if template.Name != "Custom Template" {
		t.Errorf("Expected name 'Custom Template', got '%s'", template.Name)
	}
}

func TestDefaultResolver(t *testing.T) {
	testDir := "test_default_resolver"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create template with .yaml extension
	templateContent := `name: "YAML Template"
resource_type: "yaml"
content:
  data: "{{.Data}}"
variables:
  - name: "Data"
    type: "string"
    required: true`

	err = os.WriteFile(filepath.Join(testDir, "test.yaml"), []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write template: %v", err)
	}

	resolver := NewDefaultResolver()
	path, err := resolver.Resolve(testDir, "test")
	if err != nil {
		t.Fatalf("Failed to resolve template: %v", err)
	}

	expectedPath := filepath.Join(testDir, "test.yaml")
	if path != expectedPath {
		t.Errorf("Expected path '%s', got '%s'", expectedPath, path)
	}
}
