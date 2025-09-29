package hook

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// BuiltinProcessor handles only builtin hooks with builtin-specific functionality.
type BuiltinProcessor struct {
	builtinExecutor *BuiltinExecutor
	validator       Validator
	logger          *zap.Logger
}

// NewBuiltinProcessor creates a new builtin-only processor.
func NewBuiltinProcessor(logger *zap.Logger) *BuiltinProcessor {
	return &BuiltinProcessor{
		builtinExecutor: NewBuiltinExecutor(),
		validator:       NewBuiltinValidator(),
		logger:          logger,
	}
}

// Execute executes only builtin hooks.
func (bp *BuiltinProcessor) Execute(ctx context.Context, hookType Type, data *Data) error {
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
func (bp *BuiltinProcessor) Register(hookType Type, hook BuiltinHook) error {
	bp.builtinExecutor.Register(context.Background(), hookType, hook)
	return nil
}

// Unregister unregisters a builtin hook.
func (bp *BuiltinProcessor) Unregister(hookType Type, hook BuiltinHook) error {
	bp.builtinExecutor.Unregister(context.Background(), hookType, hook)
	return nil
}

// ListBuiltinHooks returns a list of registered builtin hook types.
func (bp *BuiltinProcessor) ListBuiltinHooks() []Type {
	return bp.builtinExecutor.ListBuiltin()
}

// CountBuiltinHooks returns the number of registered builtin hooks for a given type.
func (bp *BuiltinProcessor) CountBuiltinHooks(hookType Type) int {
	return bp.builtinExecutor.CountBuiltin(hookType)
}

// GetBuiltinExecutor returns the underlying builtin executor for advanced usage.
func (bp *BuiltinProcessor) GetBuiltinExecutor() *BuiltinExecutor {
	return bp.builtinExecutor
}
