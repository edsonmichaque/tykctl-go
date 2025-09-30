package hook

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// BuiltinDispatcher handles only builtin hooks with builtin-specific functionality.
type BuiltinDispatcher struct {
	builtinExecutor *BuiltinExecutor
	validator       Validator
	logger          *zap.Logger
}

// NewBuiltinDispatcher creates a new builtin-only dispatcher.
func NewBuiltinDispatcher(logger *zap.Logger) *BuiltinDispatcher {
	return &BuiltinDispatcher{
		builtinExecutor: NewBuiltinExecutor(),
		validator:       NewBuiltinValidator(),
		logger:          logger,
	}
}

// Execute executes only builtin hooks.
func (bp *BuiltinDispatcher) Execute(ctx context.Context, hookType Type, data *Data) error {
	// Validate hook data
	if err := bp.validator.Validate(data); err != nil {
		return fmt.Errorf("hook validation failed: %w", err)
	}

	// Execute builtin hooks
	if err := bp.builtinExecutor.Execute(ctx, hookType, data); err != nil {
		if bp.logger != nil {
			bp.logger.Error("Builtin hook execution failed",
				zap.String("hook_type", string(hookType)),
				zap.String("extension", data.Extension),
				zap.Error(err),
			)
		}
		return fmt.Errorf("builtin hook execution failed: %w", err)
	}

	return nil
}

// Register registers a builtin hook.
func (bp *BuiltinDispatcher) Register(hookType Type, hook BuiltinHook) error {
	bp.builtinExecutor.Register(context.Background(), hookType, hook)
	return nil
}

// Unregister unregisters a builtin hook.
func (bp *BuiltinDispatcher) Unregister(hookType Type, hook BuiltinHook) error {
	bp.builtinExecutor.Unregister(context.Background(), hookType, hook)
	return nil
}

// ListBuiltinHooks returns a list of registered builtin hook types.
func (bp *BuiltinDispatcher) ListBuiltinHooks() []Type {
	return bp.builtinExecutor.ListBuiltin()
}

// CountBuiltinHooks returns the number of registered builtin hooks for a given type.
func (bp *BuiltinDispatcher) CountBuiltinHooks(hookType Type) int {
	return bp.builtinExecutor.CountBuiltin(hookType)
}

// GetBuiltinExecutor returns the underlying builtin executor for advanced usage.
func (bp *BuiltinDispatcher) GetBuiltinExecutor() *BuiltinExecutor {
	return bp.builtinExecutor
}
