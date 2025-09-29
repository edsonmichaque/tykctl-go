package template

// resolveBuiltinTemplate resolves a built-in template by name
func resolveBuiltinTemplate(name string) (*Template, error) {
	// This is a placeholder implementation
	// In a real implementation, this would return built-in templates
	// For now, we'll return an error indicating that built-in templates
	// should be handled by the specific command implementations
	return nil, NewTemplateError(ErrorTypeBuiltinNotFound, "", ErrBuiltinTemplateNotFound)
}
