package hook

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// BuiltinHook represents a Go function that can be executed as a hook
type BuiltinHook func(ctx context.Context, data *HookData) error

// BuiltinManager manages builtin Go hooks
type BuiltinManager struct {
	hooks map[HookType][]BuiltinHook
	mutex sync.RWMutex
}

// NewBuiltinManager creates a new builtin hook manager
func NewBuiltinManager() *BuiltinManager {
	return &BuiltinManager{
		hooks: make(map[HookType][]BuiltinHook),
	}
}

// Register registers a builtin hook for a specific type
func (bm *BuiltinManager) Register(ctx context.Context, hookType HookType, hook BuiltinHook) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	bm.hooks[hookType] = append(bm.hooks[hookType], hook)
}

// Unregister removes a builtin hook for a specific type
func (bm *BuiltinManager) Unregister(ctx context.Context, hookType HookType, hook BuiltinHook) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	hooks := bm.hooks[hookType]
	for i, h := range hooks {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", hook) {
			bm.hooks[hookType] = append(hooks[:i], hooks[i+1:]...)
			break
		}
	}
}

// Execute executes all builtin hooks for a specific type
func (bm *BuiltinManager) Execute(ctx context.Context, hookType HookType, data *HookData) error {
	bm.mutex.RLock()
	hooks := make([]BuiltinHook, len(bm.hooks[hookType]))
	copy(hooks, bm.hooks[hookType])
	bm.mutex.RUnlock()

	for _, hook := range hooks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := hook(ctx, data); err != nil {
				return fmt.Errorf("builtin hook %s failed: %w", hookType, err)
			}
		}
	}

	return nil
}

// ExecuteAsync executes all builtin hooks for a specific type asynchronously
func (bm *BuiltinManager) ExecuteAsync(ctx context.Context, hookType HookType, data *HookData) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		bm.mutex.RLock()
		hooks := make([]BuiltinHook, len(bm.hooks[hookType]))
		copy(hooks, bm.hooks[hookType])
		bm.mutex.RUnlock()

		var wg sync.WaitGroup
		var mu sync.Mutex
		var errors []error

		for _, hook := range hooks {
			wg.Add(1)
			go func(h BuiltinHook) {
				defer wg.Done()

				select {
				case <-ctx.Done():
					mu.Lock()
					errors = append(errors, ctx.Err())
					mu.Unlock()
					return
				default:
					if err := h(ctx, data); err != nil {
						mu.Lock()
						errors = append(errors, fmt.Errorf("builtin hook %s failed: %w", hookType, err))
						mu.Unlock()
					}
				}
			}(hook)
		}

		wg.Wait()

		if len(errors) > 0 {
			errChan <- fmt.Errorf("multiple builtin hook errors: %v", errors)
		}
	}()

	return errChan
}

// List returns all registered builtin hooks for a specific type
func (bm *BuiltinManager) List(ctx context.Context, hookType HookType) []BuiltinHook {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	hooks := make([]BuiltinHook, len(bm.hooks[hookType]))
	copy(hooks, bm.hooks[hookType])
	return hooks
}

// Clear removes all builtin hooks for a specific type
func (bm *BuiltinManager) Clear(ctx context.Context, hookType HookType) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	delete(bm.hooks, hookType)
}

// ClearAll removes all builtin hooks
func (bm *BuiltinManager) ClearAll(ctx context.Context) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	bm.hooks = make(map[HookType][]BuiltinHook)
}

// Count returns the number of builtin hooks for a specific type
func (bm *BuiltinManager) Count(ctx context.Context, hookType HookType) int {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	return len(bm.hooks[hookType])
}

// HasHooks returns true if there are builtin hooks for a specific type
func (bm *BuiltinManager) HasHooks(ctx context.Context, hookType HookType) bool {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	return len(bm.hooks[hookType]) > 0
}

// HookTypes returns all registered hook types
func (bm *BuiltinManager) HookTypes(ctx context.Context) []HookType {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	types := make([]HookType, 0, len(bm.hooks))
	for hookType := range bm.hooks {
		types = append(types, hookType)
	}

	return types
}

// Predefined builtin hooks
func LoggingHook(hookType HookType) BuiltinHook {
	return func(ctx context.Context, data *HookData) error {
		log.Printf("[BUILTIN HOOK] %s executed", hookType)
		return nil
	}
}

func TimingHook(hookType HookType) BuiltinHook {
	return func(ctx context.Context, data *HookData) error {
		start := time.Now()
		defer func() {
			log.Printf("[BUILTIN HOOK] %s took %v", hookType, time.Since(start))
		}()
		return nil
	}
}

func ValidationHook(validator func(*HookData) error) BuiltinHook {
	return func(ctx context.Context, data *HookData) error {
		return validator(data)
	}
}

func MetricsHook(collector func(string, map[string]interface{}) error) BuiltinHook {
	return func(ctx context.Context, data *HookData) error {
		metrics := map[string]interface{}{
			"extension_name": data.ExtensionName,
			"extension_path": data.ExtensionPath,
			"timestamp":      time.Now().Unix(),
		}

		if data.Metadata != nil {
			for k, v := range data.Metadata {
				metrics[k] = v
			}
		}

		return collector("extension_operation", metrics)
	}
}

func NotificationHook(notifier func(string, string) error) BuiltinHook {
	return func(ctx context.Context, data *HookData) error {
		message := fmt.Sprintf("Extension %s operation completed", data.ExtensionName)
		return notifier("TykCLI", message)
	}
}

func RetryHook(maxRetries int, delay time.Duration) BuiltinHook {
	return func(ctx context.Context, data *HookData) error {
		if data.Error == nil {
			return nil
		}

		for i := 0; i < maxRetries; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				log.Printf("[BUILTIN HOOK] Retrying operation (attempt %d/%d)", i+1, maxRetries)
			}
		}

		return fmt.Errorf("max retries exceeded: %w", data.Error)
	}
}
