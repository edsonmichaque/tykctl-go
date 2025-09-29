package template

import "context"

// Loader defines the interface for loading templates from various sources.
type Loader interface {
	Load(ctx context.Context) (*Template, error)
}

