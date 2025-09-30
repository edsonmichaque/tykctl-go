package hook

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// SchemaDispatcher handles only JSON Schema validation hooks with schema-specific functionality.
type SchemaDispatcher struct {
	schemaExecutor *SchemaExecutor
	validator      Validator
	logger         *zap.Logger
}

// NewSchemaDispatcher creates a new JSON Schema-only dispatcher.
func NewSchemaDispatcher(logger *zap.Logger, schemaDir string) *SchemaDispatcher {
	return &SchemaDispatcher{
		schemaExecutor: NewSchemaExecutor(logger, schemaDir),
		validator:      NewSchemaValidator(schemaDir),
		logger:         logger,
	}
}

// Execute executes only JSON Schema validation hooks.
func (sp *SchemaDispatcher) Execute(ctx context.Context, hookType Type, data *Data) error {
	// Validate hook data
	if err := sp.validator.Validate(data); err != nil {
		return fmt.Errorf("hook validation failed: %w", err)
	}

	// Execute JSON Schema validation hooks
	if err := sp.schemaExecutor.Execute(ctx, hookType, data); err != nil {
		if sp.logger != nil {
			sp.logger.Error("JSON Schema validation failed",
				zap.String("hook_type", string(hookType)),
				zap.String("extension", data.Extension),
				zap.Error(err),
			)
		}
		return fmt.Errorf("JSON schema validation failed: %w", err)
	}

	return nil
}

// GetSchemaExecutor returns the underlying JSON Schema executor for advanced usage.
func (sp *SchemaDispatcher) GetSchemaExecutor() *SchemaExecutor {
	return sp.schemaExecutor
}

// ValidateJSONSchema validates a JSON Schema file without executing hooks.
func (sp *SchemaDispatcher) ValidateJSONSchema(ctx context.Context, schemaFile string) error {
	if sp.schemaExecutor == nil {
		return fmt.Errorf("JSON schema executor not available")
	}

	// This would need to be implemented in SchemaExecutor
	// For now, return a placeholder
	return fmt.Errorf("JSON schema validation not yet implemented")
}

// ListJSONSchemas returns a list of available JSON Schema files for a given hook type.
// This method adds dispatcher-level validation and error handling.
func (sp *SchemaDispatcher) ListJSONSchemas(ctx context.Context, hookType Type) ([]string, error) {
	if sp.schemaExecutor == nil {
		return nil, fmt.Errorf("schema executor not available")
	}

	// Validate hook type before listing schemas
	if hookType == "" {
		return nil, fmt.Errorf("hook type cannot be empty")
	}

	allSchemas, err := sp.schemaExecutor.discoverSchemas(ctx, hookType)
	if err != nil {
		return nil, fmt.Errorf("failed to discover schemas: %w", err)
	}
	schemas := allSchemas[hookType]

	if sp.logger != nil {
		sp.logger.Debug("Listed schemas for hook type",
			zap.String("hook_type", string(hookType)),
			zap.Int("count", len(schemas)),
		)
	}

	return schemas, nil
}

// DiscoverAllJSONSchemas discovers all available JSON Schema files.
// This method adds dispatcher-level validation and error handling.
func (sp *SchemaDispatcher) DiscoverAllJSONSchemas() (map[Type][]string, error) {
	if sp.schemaExecutor == nil {
		return nil, fmt.Errorf("schema executor not available")
	}

	allSchemas, err := sp.schemaExecutor.discoverSchemas(context.Background())
	if err != nil {
		if sp.logger != nil {
			sp.logger.Error("Failed to discover all schemas",
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("failed to discover all schemas: %w", err)
	}

	if sp.logger != nil {
		totalSchemas := 0
		for _, schemas := range allSchemas {
			totalSchemas += len(schemas)
		}
		sp.logger.Debug("Discovered all schemas",
			zap.Int("total_schemas", totalSchemas),
			zap.Int("hook_types", len(allSchemas)),
		)
	}

	return allSchemas, nil
}

// CountJSONSchemas returns the count of available JSON Schema files for a given hook type.
// This method adds dispatcher-level validation and error handling.
func (sp *SchemaDispatcher) CountJSONSchemas(ctx context.Context, hookType Type) (int, error) {
	if sp.schemaExecutor == nil {
		return 0, fmt.Errorf("schema executor not available")
	}

	// Validate hook type before counting schemas
	if hookType == "" {
		return 0, fmt.Errorf("hook type cannot be empty")
	}

	allSchemas, err := sp.schemaExecutor.discoverSchemas(ctx, hookType)
	if err != nil {
		return 0, fmt.Errorf("failed to discover schemas: %w", err)
	}
	schemas := allSchemas[hookType]

	if sp.logger != nil {
		sp.logger.Debug("Counted schemas for hook type",
			zap.String("hook_type", string(hookType)),
			zap.Int("count", len(schemas)),
		)
	}

	return len(schemas), nil
}

// ReloadJSONSchemas reloads all JSON schemas from the schema directory.
func (sp *SchemaDispatcher) ReloadJSONSchemas(ctx context.Context) error {
	if sp.schemaExecutor == nil {
		return fmt.Errorf("JSON schema executor not available")
	}

	return sp.schemaExecutor.ReloadSchemas(ctx)
}

// GetJSONSchemaDirectory returns the JSON schema directory path.
// This method adds dispatcher-level validation and error handling.
func (sp *SchemaDispatcher) GetJSONSchemaDirectory() (string, error) {
	if sp.schemaExecutor == nil {
		return "", fmt.Errorf("schema executor not available")
	}

	schemaDir := sp.schemaExecutor.getSchemaDir()

	if sp.logger != nil {
		sp.logger.Debug("Retrieved schema directory",
			zap.String("schema_dir", schemaDir),
		)
	}

	return schemaDir, nil
}

// GetJSONSchemaForHookType returns the JSON schema for a specific hook type.
func (sp *SchemaDispatcher) GetJSONSchemaForHookType(hookType Type) (map[string]interface{}, error) {
	if sp.schemaExecutor == nil {
		return nil, fmt.Errorf("JSON schema executor not available")
	}

	schemas := sp.schemaExecutor.GetLoadedSchemas()
	schema, exists := schemas[hookType]
	if !exists {
		return nil, fmt.Errorf("no JSON schema found for hook type: %s", hookType)
	}

	return schema, nil
}

// ValidateDataAgainstJSONSchema validates data against a specific JSON schema.
func (sp *SchemaDispatcher) ValidateDataAgainstJSONSchema(ctx context.Context, data *Data, schema map[string]interface{}) error {
	if sp.schemaExecutor == nil {
		return fmt.Errorf("JSON schema executor not available")
	}

	// Convert data to JSON for validation
	dataJSON, err := sp.schemaExecutor.dataToJSON(data)
	if err != nil {
		return fmt.Errorf("failed to convert data to JSON: %w", err)
	}

	// Validate against JSON schema
	if err := sp.schemaExecutor.validateAgainstJSONSchema(dataJSON, schema); err != nil {
		return fmt.Errorf("JSON schema validation failed: %w", err)
	}

	return nil
}

// getJSONSchemaKeys extracts keys from a JSON schema map for logging.
func getJSONSchemaKeys(schema map[string]interface{}) []string {
	var keys []string
	for key := range schema {
		keys = append(keys, key)
	}
	return keys
}
