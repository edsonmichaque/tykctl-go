package hook

import (
	"context"
)

// Processor defines the interface that all hook processors must implement.
type Processor interface {
	// Execute executes a hook of the given type with the provided data.
	Execute(ctx context.Context, hookType Type, data *Data) error
}
