package hook

// Type represents the type of hook
type Type string

// Data contains data passed to hooks
type Data struct {
	Type      Type                   `json:"type"`
	Extension string                 `json:"extension"`
	Error     error                  `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewData creates a new Data instance
func NewData(hookType Type, extensionName string) *Data {
	return &Data{
		Type:      hookType,
		Extension: extensionName,
		Metadata:  make(map[string]interface{}),
	}
}

// WithError sets the error
func (h *Data) WithError(err error) *Data {
	h.Error = err
	return h
}

// WithMetadata sets metadata
func (h *Data) WithMetadata(key string, value interface{}) *Data {
	if h.Metadata == nil {
		h.Metadata = make(map[string]interface{})
	}
	h.Metadata[key] = value
	return h
}

// WithMetadataMap sets multiple metadata entries
func (h *Data) WithMetadataMap(metadata map[string]interface{}) *Data {
	if h.Metadata == nil {
		h.Metadata = make(map[string]interface{})
	}
	for key, value := range metadata {
		h.Metadata[key] = value
	}
	return h
}
