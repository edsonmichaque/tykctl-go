package template

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FileLoader loads templates from local files.
type FileLoader struct {
	filePath string
}

// NewFileLoader creates a new file loader.
func NewFileLoader(filePath string) *FileLoader {
	return &FileLoader{
		filePath: filePath,
	}
}

// Load loads a template from a file.
func (f *FileLoader) Load(ctx context.Context) (*Template, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var template Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &template, nil
}

