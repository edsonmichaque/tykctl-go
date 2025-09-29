package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
)

// SchemaExecutor handles JSON Schema validation for hooks.
type SchemaExecutor struct {
	schemaDir  string
	logger     *zap.Logger
	schemas    map[Type]map[string]interface{}
	discovered map[Type][]string // Cache for discovered schemas
}

// NewSchemaExecutor creates a new schema executor.
func NewSchemaExecutor(logger *zap.Logger, schemaDir string) *SchemaExecutor {
	return &SchemaExecutor{
		schemaDir:  schemaDir,
		logger:     logger,
		schemas:    make(map[Type]map[string]interface{}),
		discovered: make(map[Type][]string),
	}
}

// Execute executes JSON Schema validation for the given hook type and data.
func (se *SchemaExecutor) Execute(ctx context.Context, hookType Type, data *Data) error {
	// Load schemas if not already loaded
	if err := se.loadSchemas(ctx); err != nil {
		if se.logger != nil {
			se.logger.Error("Failed to load JSON schemas",
				zap.String("schema_dir", se.schemaDir),
				zap.Error(err),
			)
		}
		return fmt.Errorf("failed to load JSON schemas: %w", err)
	}

	// Get schema for the hook type
	schema, exists := se.schemas[hookType]
	if !exists {
		if se.logger != nil {
			se.logger.Debug("No JSON schema found for hook type",
				zap.String("hook_type", string(hookType)),
			)
		}
		return nil // No schema means no validation required
	}

	// Convert data to JSON for validation
	dataJSON, err := se.dataToJSON(data)
	if err != nil {
		return fmt.Errorf("failed to convert data to JSON: %w", err)
	}

	// Validate against JSON Schema
	if err := se.validateAgainstJSONSchema(dataJSON, schema); err != nil {
		return NewHookError(hookType, data.Extension, "JSON schema validation failed", err)
	}

	if se.logger != nil {
		se.logger.Debug("JSON schema validation passed",
			zap.String("hook_type", string(hookType)),
			zap.String("extension", data.Extension),
		)
	}

	return nil
}

// loadSchemas loads all JSON Schema files from the schema directory.
func (se *SchemaExecutor) loadSchemas(ctx context.Context) error {
	if se.schemaDir == "" {
		return nil // No schema directory specified
	}

	// Check if directory exists
	if _, err := os.Stat(se.schemaDir); os.IsNotExist(err) {
		return fmt.Errorf("JSON schema directory does not exist: %s", se.schemaDir)
	}

	// Read the schema directory to find hook types
	entries, err := os.ReadDir(se.schemaDir)
	if err != nil {
		return fmt.Errorf("failed to read schema directory %s: %w", se.schemaDir, err)
	}

	// Process each entry as a potential hook type
	for _, entry := range entries {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		entryName := entry.Name()

		// Skip hidden files and directories
		if strings.HasPrefix(entryName, ".") {
			continue
		}

		// Convert entry name to hook type
		hookType := Type(entryName)

		// Discover schemas for this hook type
		schemaFiles, err := se.findSchemaFiles(ctx, hookType)
		if err != nil {
			if se.logger != nil {
				se.logger.Warn("Failed to discover schemas for hook type",
					zap.String("hook_type", entryName),
					zap.Error(err),
				)
			}
			continue
		}

		// Load the first valid schema for this hook type
		if len(schemaFiles) > 0 {
			schema, err := se.loadJSONSchemaFile(schemaFiles[0])
			if err != nil {
				if se.logger != nil {
					se.logger.Warn("Failed to load JSON schema file",
						zap.String("hook_type", string(hookType)),
						zap.String("file", schemaFiles[0]),
						zap.Error(err),
					)
				}
				continue
			}

			// Validate that it's a proper JSON Schema using gojsonschema
			if err := se.validateJSONSchemaFile(schema); err != nil {
				if se.logger != nil {
					se.logger.Warn("Invalid JSON schema file",
						zap.String("hook_type", string(hookType)),
						zap.String("file", schemaFiles[0]),
						zap.Error(err),
					)
				}
				continue
			}

			// Store schema
			se.schemas[hookType] = schema

			if se.logger != nil {
				se.logger.Debug("Loaded JSON schema",
					zap.String("hook_type", string(hookType)),
					zap.String("file", schemaFiles[0]),
				)
			}
		}
	}

	return nil
}

// findSchemaFiles finds JSON Schema files for a specific hook type.
func (se *SchemaExecutor) findSchemaFiles(ctx context.Context, hookType Type) ([]string, error) {
	// Check cache first
	if schemas, exists := se.discovered[hookType]; exists {
		return schemas, nil
	}

	var schemaFiles []string

	// Check if hookType is a direct file or directory
	schemaPath := filepath.Join(se.schemaDir, string(hookType))

	info, err := os.Stat(schemaPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Hook type doesn't exist, return empty list
			se.discovered[hookType] = []string{}
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to stat schema path %s: %w", schemaPath, err)
	}

	if info.IsDir() {
		// Hook type is a directory, process its entries in lexicographic order
		schemaFiles, err = se.processSchemaDirectory(ctx, schemaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to process schema directory %s: %w", schemaPath, err)
		}
	} else {
		// Hook type is a file, check if it's a valid JSON Schema
		if se.isValidJSONSchemaFile(info.Name(), schemaPath) {
			schemaFiles = []string{schemaPath}
		}
	}

	// Cache the discovered schemas
	se.discovered[hookType] = schemaFiles

	if se.logger != nil {
		se.logger.Debug("Discovered schemas",
			zap.String("hook_type", string(hookType)),
			zap.String("schema_path", schemaPath),
			zap.Bool("is_directory", info.IsDir()),
			zap.Int("count", len(schemaFiles)),
			zap.Strings("schemas", schemaFiles),
		)
	}

	return schemaFiles, nil
}

// processSchemaDirectory processes a schema directory and returns schemas in lexicographic order.
func (se *SchemaExecutor) processSchemaDirectory(ctx context.Context, dirPath string) ([]string, error) {
	var schemas []string

	// Read directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// Process entries in lexicographic order (ReadDir already sorts by name)
	for _, entry := range entries {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if entry.IsDir() {
			// Skip subdirectories for now
			continue
		}

		entryPath := filepath.Join(dirPath, entry.Name())

		// Check if it's a valid JSON Schema file
		if se.isValidJSONSchemaFile(entry.Name(), entryPath) {
			schemas = append(schemas, entryPath)
		}
	}

	return schemas, nil
}

// isValidJSONSchemaFile checks if a file is a valid JSON Schema file.
func (se *SchemaExecutor) isValidJSONSchemaFile(filename, fullPath string) bool {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".json" {
		return false
	}

	// Check if file exists and is readable
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}

	// Check if it's a file (not directory)
	if info.IsDir() {
		return false
	}

	// Basic check for JSON Schema content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return false
	}

	contentStr := string(content)

	// Check for JSON Schema keywords
	jsonSchemaKeywords := []string{"$schema", "type", "properties", "required", "items", "additionalProperties", "definitions", "$ref"}
	hasKeyword := false
	for _, keyword := range jsonSchemaKeywords {
		if strings.Contains(contentStr, keyword) {
			hasKeyword = true
			break
		}
	}

	return hasKeyword
}

// loadJSONSchemaFile loads a single JSON Schema file.
func (se *SchemaExecutor) loadJSONSchemaFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON schema file: %w", err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse JSON schema: %w", err)
	}

	return schema, nil
}

// validateJSONSchemaFile validates that the loaded file is a proper JSON Schema using gojsonschema.
func (se *SchemaExecutor) validateJSONSchemaFile(schema map[string]interface{}) error {
	// Convert schema to JSON bytes
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema to JSON: %w", err)
	}

	// Create schema loader
	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)

	// Validate the schema itself using gojsonschema
	_, err = gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("invalid JSON schema: %w", err)
	}

	return nil
}

// discoverSchemas discovers JSON Schema files for specific hook types or all hook types.
// If no hook types are provided, it discovers all schemas. Otherwise, it discovers schemas for the specified types.
func (se *SchemaExecutor) discoverSchemas(ctx context.Context, hookTypes ...Type) (map[Type][]string, error) {
	allSchemas := make(map[Type][]string)

	if se.schemaDir == "" {
		return allSchemas, nil
	}

	// If specific hook types requested, check cache first
	if len(hookTypes) > 0 {
		// Check if all requested types are in cache
		allCached := true
		for _, hookType := range hookTypes {
			if schemas, exists := se.discovered[hookType]; exists {
				allSchemas[hookType] = schemas
			} else {
				allCached = false
				break
			}
		}
		if allCached {
			return allSchemas, nil
		}
	}

	// Read the schema directory to find hook types
	entries, err := os.ReadDir(se.schemaDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema directory %s: %w", se.schemaDir, err)
	}

	// Process each entry as a potential hook type
	for _, entry := range entries {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		entryName := entry.Name()

		// Skip hidden files and directories
		if strings.HasPrefix(entryName, ".") {
			continue
		}

		// Convert entry name to hook type
		currentHookType := Type(entryName)

		// If specific hook types requested, only process those types
		if len(hookTypes) > 0 {
			found := false
			for _, hookType := range hookTypes {
				if currentHookType == hookType {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Discover schemas for this hook type
		schemas, err := se.findSchemaFiles(ctx, currentHookType)
		if err != nil {
			if se.logger != nil {
				se.logger.Warn("Failed to discover schemas for hook type",
					zap.String("hook_type", entryName),
					zap.Error(err),
				)
			}
			continue
		}

		if len(schemas) > 0 {
			allSchemas[currentHookType] = schemas
		}
	}

	// Update cache
	se.discovered = allSchemas

	if se.logger != nil {
		totalSchemas := 0
		for _, schemas := range allSchemas {
			totalSchemas += len(schemas)
		}
		se.logger.Debug("Discovered schemas",
			zap.String("schema_dir", se.schemaDir),
			zap.Strings("requested_types", hookTypesToStrings(hookTypes)),
			zap.Int("total_schemas", totalSchemas),
			zap.Int("hook_types", len(allSchemas)),
		)
	}

	return allSchemas, nil
}

// listDiscoveredSchemas returns all discovered schemas for a hook type.
func (se *SchemaExecutor) listDiscoveredSchemas(hookType Type) []string {
	if schemas, exists := se.discovered[hookType]; exists {
		return schemas
	}
	return []string{}
}

// refreshDiscovery refreshes the schema discovery cache.
func (se *SchemaExecutor) refreshDiscovery() error {
	_, err := se.discoverSchemas(context.Background())
	return err
}

// countDiscoveredSchemas returns the count of discovered schemas for a hook type.
func (se *SchemaExecutor) countDiscoveredSchemas(hookType Type) int {
	if schemas, exists := se.discovered[hookType]; exists {
		return len(schemas)
	}
	return 0
}

// dataToJSON converts hook data to JSON for validation.
func (se *SchemaExecutor) dataToJSON(data *Data) (map[string]interface{}, error) {
	jsonData := map[string]interface{}{
		"type":      data.Type,
		"extension": data.Extension,
		"error":     data.Error,
		"metadata":  data.Metadata,
	}

	return jsonData, nil
}

// validateAgainstJSONSchema validates data against a JSON Schema using gojsonschema.
func (se *SchemaExecutor) validateAgainstJSONSchema(data map[string]interface{}, schema map[string]interface{}) error {
	// Convert schema to JSON bytes
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema to JSON: %w", err)
	}

	// Convert data to JSON bytes
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Create schema loader
	schemaLoader := gojsonschema.NewBytesLoader(schemaBytes)
	documentLoader := gojsonschema.NewBytesLoader(dataBytes)

	// Validate using gojsonschema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("JSON schema validation error: %w", err)
	}

	// Check validation result
	if !result.Valid() {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
		return fmt.Errorf("JSON schema validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// getSchemaDir returns the schema directory path.
func (se *SchemaExecutor) getSchemaDir() string {
	return se.schemaDir
}

// GetLoadedSchemas returns the currently loaded schemas.
func (se *SchemaExecutor) GetLoadedSchemas() map[Type]map[string]interface{} {
	return se.schemas
}

// ReloadSchemas reloads all schemas from the schema directory.
func (se *SchemaExecutor) ReloadSchemas(ctx context.Context) error {
	se.schemas = make(map[Type]map[string]interface{})
	se.discovered = make(map[Type][]string) // Clear discovery cache
	return se.loadSchemas(ctx)
}
