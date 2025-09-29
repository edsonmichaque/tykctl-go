package hook

import (
	"context"
	"fmt"
	"sync"
)

// Executor defines the interface for executing hooks.
type Executor interface {
	Execute(ctx context.Context, hookType Type, data *Data) error
}

// BuiltinExecutor manages and executes builtin Go hooks.
type BuiltinExecutor struct {
	hooks map[Type][]BuiltinHook
	mutex sync.RWMutex
}

// NewBuiltinExecutor creates a new builtin hook executor.
func NewBuiltinExecutor() *BuiltinExecutor {
	return &BuiltinExecutor{
		hooks: make(map[Type][]BuiltinHook),
	}
}

// Register registers a builtin hook for a specific type.
func (be *BuiltinExecutor) Register(ctx context.Context, hookType Type, hook BuiltinHook) {
	be.mutex.Lock()
	defer be.mutex.Unlock()
	be.hooks[hookType] = append(be.hooks[hookType], hook)
}

// Unregister removes a builtin hook for a specific type.
func (be *BuiltinExecutor) Unregister(ctx context.Context, hookType Type, hook BuiltinHook) {
	be.mutex.Lock()
	defer be.mutex.Unlock()
	hooks := be.hooks[hookType]
	for i, h := range hooks {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", hook) {
			be.hooks[hookType] = append(hooks[:i], hooks[i+1:]...)
			break
		}
	}
}

// Execute executes all builtin hooks for a specific type.
func (be *BuiltinExecutor) Execute(ctx context.Context, hookType Type, data *Data) error {
	be.mutex.RLock()
	defer be.mutex.RUnlock()

	hooks := be.hooks[hookType]
	for _, hook := range hooks {
		select {
		case <-ctx.Done():
			return ErrWithContext(ErrContextCancelled, "builtin hook execution")
		default:
			if err := hook(ctx, data); err != nil {
				return NewHookError(hookType, data.Extension, "builtin execution failed", err)
			}
		}
	}
	return nil
}

// ListBuiltin returns a list of registered builtin hook types.
func (be *BuiltinExecutor) ListBuiltin() []Type {
	be.mutex.RLock()
	defer be.mutex.RUnlock()
	var types []Type
	for t := range be.hooks {
		types = append(types, t)
	}
	return types
}

// CountBuiltin returns the number of registered builtin hooks for a given type.
func (be *BuiltinExecutor) CountBuiltin(hookType Type) int {
	be.mutex.RLock()
	defer be.mutex.RUnlock()
	return len(be.hooks[hookType])
}

// BuiltinHook represents a Go function that can be executed as a hook.
type BuiltinHook func(ctx context.Context, data *Data) error
