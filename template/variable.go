package template

// Variable represents a template variable.
type Variable struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"`
	Required    bool        `yaml:"required"`
	Default     interface{} `yaml:"default"`
	Description string      `yaml:"description"`
	Validation  Validation  `yaml:"validation"`
}

// Validation represents variable validation rules.
type Validation struct {
	MinLength *int     `yaml:"min_length,omitempty"`
	MaxLength *int     `yaml:"max_length,omitempty"`
	Pattern   string   `yaml:"pattern,omitempty"`
	Enum      []string `yaml:"enum,omitempty"`
	MinValue  *float64 `yaml:"min_value,omitempty"`
	MaxValue  *float64 `yaml:"max_value,omitempty"`
}


