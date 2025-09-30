# TykCtl Config Package Examples

This directory contains examples demonstrating how to use the TykCtl configuration package.

## Running the Examples

```bash
# Run the comprehensive example
go run example_usage.go
```

## What the Examples Demonstrate

### 1. Basic Configuration Loading
- Creating a config loader with caching, logging, and metrics
- Loading configuration into a struct with default values
- Error handling for configuration loading

### 2. Context Management
- Creating a context store for managing multiple environments
- Creating development and production contexts
- Switching between contexts
- Getting current context information

### 3. Resource Discovery
- Discovering hooks with filtering
- Discovering plugins with filtering
- Discovering templates with filtering

### 4. Custom Loader Configuration
- Configuring custom environment prefixes
- Setting custom configuration formats
- Adding custom configuration paths
- Adding custom context paths

## Key Features Demonstrated

- **Multi-Context Support**: Like kubectl, gh, aws CLI contexts
- **Configurable Properties**: Custom environment prefixes, formats, and paths
- **Error Handling**: Comprehensive error handling with rollback mechanisms
- **Resource Discovery**: Automatic discovery of hooks, plugins, templates, and cache configs
- **Validation**: Built-in validation with custom validators
- **Caching**: Performance optimization with caching
- **Logging**: Structured logging for debugging and monitoring
- **Metrics**: Built-in metrics collection for observability

## Expected Output

The example will show:
- Configuration loading with defaults
- Context creation and switching
- Resource discovery results
- Custom loader configuration details

Note: Some operations may show "not implemented" errors as this is a demonstration package with placeholder implementations for the core configuration loading logic.