package hook

import (
	"context"
)

// Dispatcher defines the interface that all hook dispatchers must implement.
type Dispatcher interface {
	// Execute executes a hook of the given type with the provided data.
	Execute(ctx context.Context, hookType Type, data *Data) error
}
