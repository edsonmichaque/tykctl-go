package template

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Processor handles template processing.
type Processor struct {
	loader    Loader
	validator Validator
}

// NewProcessor creates a new template processor.
func NewProcessor(loader Loader, validator Validator) *Processor {
	return &Processor{
		loader:    loader,
		validator: validator,
	}
}

// Process loads, validates, and processes a template from a source.
func (p *Processor) Process(ctx context.Context, variables map[string]interface{}) ([]byte, error) {
	// Load template
	tmpl, err := p.loader.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	// Validate template
	if err := p.validator.Validate(tmpl); err != nil {
		return nil, fmt.Errorf("template validation failed: %w", err)
	}

	// Process template
	return p.process(ctx, tmpl, variables)
}

// process processes a template with the given variables.
func (p *Processor) process(ctx context.Context, tmpl *Template, variables map[string]interface{}) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Convert template content to YAML string for processing
	contentBytes, err := templateToYAML(tmpl.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to convert template content to YAML: %w", err)
	}

	// Create Go template from content
	goTemplate, err := template.New(tmpl.Name).Parse(string(contentBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with variables
	var buf bytes.Buffer
	if err := goTemplate.Execute(&buf, variables); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// templateToYAML converts a map to YAML string
func templateToYAML(content map[string]interface{}) ([]byte, error) {
	return yaml.Marshal(content)
}



