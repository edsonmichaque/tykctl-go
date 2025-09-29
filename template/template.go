package template

// Template represents a template structure.
type Template struct {
	Name         string                 `yaml:"name"`
	Description  string                 `yaml:"description"`
	ResourceType string                 `yaml:"resource_type"`
	Version      string                 `yaml:"version"`
	Author       string                 `yaml:"author"`
	Tags         []string               `yaml:"tags"`
	Variables    []Variable             `yaml:"variables"`
	Content      map[string]interface{} `yaml:"content"`
	Metadata     map[string]interface{} `yaml:"metadata"`
}

