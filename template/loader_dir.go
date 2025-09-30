package template

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Resolver resolves template names to file paths.
type Resolver struct {
	Extensions  []string
	ResolveFunc func(dir, name string) (string, error)
}

// NewDefaultResolver creates a default resolver with common extensions.
func NewDefaultResolver() *Resolver {
	return &Resolver{
		Extensions:  []string{".yaml", ".yml", ".json"},
		ResolveFunc: defaultResolveFunc,
	}
}

// defaultResolveFunc is the default resolution function.
func defaultResolveFunc(dir, name string) (string, error) {
	extensions := []string{".yaml", ".yml", ".json"}

	// Try different extensions
	for _, ext := range extensions {
		path := filepath.Join(dir, name+ext)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Try with directory name as prefix
	for _, ext := range extensions {
		path := filepath.Join(dir, name, name+ext)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Try with template name as directory
	for _, ext := range extensions {
		path := filepath.Join(dir, name, "template"+ext)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("template not found: %s", filepath.Join(dir, name))
}

// Resolve resolves a template name to a file path.
func (r *Resolver) Resolve(dir, name string) (string, error) {
	if r.ResolveFunc != nil {
		return r.ResolveFunc(dir, name)
	}
	return defaultResolveFunc(dir, name)
}

// NewCustomResolver creates a resolver with custom resolution logic.
func NewCustomResolver(resolveFunc func(dir, name string) (string, error)) *Resolver {
	return &Resolver{
		ResolveFunc: resolveFunc,
	}
}

// DirLoader loads templates from a directory by name.
type DirLoader struct {
	dir      string
	name     string
	resolver *Resolver
}

// NewDirLoader creates a new directory loader.
func NewDirLoader(dir, name string, resolver *Resolver) *DirLoader {
	if resolver == nil {
		resolver = NewDefaultResolver()
	}
	return &DirLoader{
		dir:      dir,
		name:     name,
		resolver: resolver,
	}
}

// Load loads a template from the directory.
func (d *DirLoader) Load(ctx context.Context) (*Template, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled: directory loader load")
	default:
	}

	// Resolve template name to file path
	filePath, err := d.resolver.Resolve(d.dir, d.name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template %s: %w", filepath.Join(d.dir, d.name), err)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Parse template
	var template Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", filePath, err)
	}

	return &template, nil
}

// ListTemplates lists available templates in the directory.
func (d *DirLoader) ListTemplates() ([]string, error) {
	var templates []string

	err := filepath.Walk(d.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if file has a supported extension
		ext := strings.ToLower(filepath.Ext(path))
		for _, supportedExt := range []string{".yaml", ".yml", ".json"} {
			if ext == supportedExt {
				// Extract template name
				relPath, err := filepath.Rel(d.dir, path)
				if err != nil {
					return err
				}

				// Remove extension
				name := strings.TrimSuffix(relPath, ext)

				// If the file is in a subdirectory, use the directory name as template name
				if filepath.Dir(name) != "." {
					name = filepath.Dir(name)
				}

				templates = append(templates, name)
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	return templates, nil
}

// GetTemplatePath returns the resolved path for a template name.
func (d *DirLoader) GetTemplatePath(name string) (string, error) {
	return d.resolver.Resolve(d.dir, name)
}
