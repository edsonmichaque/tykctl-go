package hook

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-policy-agent/opa/rego"
	"go.uber.org/zap"
)

// RegoHook represents an OPA/Rego-based hook
type RegoHook struct {
	Name        string
	Description string
	Policy      string
	Input       map[string]interface{}
	Query       string
	Enabled     bool
	Timeout     int
	Logger      *zap.Logger
}

// RegoHookManager manages OPA/Rego hooks
type RegoHookManager struct {
	hooks  map[string]*RegoHook
	logger *zap.Logger
}

// NewRegoHookManager creates a new Rego hook manager
func NewRegoHookManager(logger *zap.Logger) *RegoHookManager {
	return &RegoHookManager{
		hooks:  make(map[string]*RegoHook),
		logger: logger,
	}
}

// RegisterRegoHook registers a new Rego hook
func (m *RegoHookManager) RegisterRegoHook(ctx context.Context, hook *RegoHook) error {
	if hook.Name == "" {
		return fmt.Errorf("hook name cannot be empty")
	}
	
	if hook.Policy == "" {
		return fmt.Errorf("hook policy cannot be empty")
	}
	
	if hook.Query == "" {
		hook.Query = "data.policy.allow" // Default query
	}
	
	if hook.Timeout == 0 {
		hook.Timeout = 30 // Default 30 seconds
	}
	
	if hook.Logger == nil {
		hook.Logger = m.logger
	}
	
	m.hooks[hook.Name] = hook
	m.logger.Info("Registered Rego hook", zap.String("name", hook.Name))
	
	return nil
}

// UnregisterRegoHook removes a Rego hook
func (m *RegoHookManager) UnregisterRegoHook(ctx context.Context, name string) error {
	if _, exists := m.hooks[name]; !exists {
		return fmt.Errorf("hook %s not found", name)
	}
	
	delete(m.hooks, name)
	m.logger.Info("Unregistered Rego hook", zap.String("name", name))
	
	return nil
}

// ExecuteRegoHook executes a Rego hook
func (m *RegoHookManager) ExecuteRegoHook(ctx context.Context, name string, input map[string]interface{}) (*RegoResult, error) {
	hook, exists := m.hooks[name]
	if !exists {
		return nil, fmt.Errorf("hook %s not found", name)
	}
	
	if !hook.Enabled {
		return &RegoResult{
			Allowed: false,
			Reason:  "hook disabled",
		}, nil
	}
	
	// Merge hook input with provided input
	mergedInput := make(map[string]interface{})
	for k, v := range hook.Input {
		mergedInput[k] = v
	}
	for k, v := range input {
		mergedInput[k] = v
	}
	
	// Create Rego query
	query := rego.New(
		rego.Query(hook.Query),
		rego.Module(hook.Name, hook.Policy),
	)
	
	// Prepare query
	preparedQuery, err := query.PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare Rego query: %w", err)
	}
	
	// Execute query
	results, err := preparedQuery.Eval(ctx, rego.EvalInput(mergedInput))
	if err != nil {
		return nil, fmt.Errorf("failed to execute Rego query: %w", err)
	}
	
	// Process results
	result := &RegoResult{
		Allowed: false,
		Reason:  "no decision made",
		Data:    make(map[string]interface{}),
	}
	
	if len(results) > 0 {
		for _, r := range results {
			if r.Expressions != nil {
				for _, expr := range r.Expressions {
					if expr.Value != nil {
						// Handle boolean result
						if allowed, ok := expr.Value.(bool); ok {
							result.Allowed = allowed
							if allowed {
								result.Reason = "policy allowed"
							} else {
								result.Reason = "policy denied"
							}
						}
						
						// Handle object result
						if obj, ok := expr.Value.(map[string]interface{}); ok {
							if allowed, exists := obj["allow"]; exists {
								if allow, ok := allowed.(bool); ok {
									result.Allowed = allow
									if allow {
										result.Reason = "policy allowed"
									} else {
										result.Reason = "policy denied"
									}
								}
							}
							
							if reason, exists := obj["reason"]; exists {
								if reasonStr, ok := reason.(string); ok {
									result.Reason = reasonStr
								}
							}
							
							// Store additional data
							for k, v := range obj {
								if k != "allow" && k != "reason" {
									result.Data[k] = v
								}
							}
						}
					}
				}
			}
		}
	}
	
	hook.Logger.Info("Executed Rego hook",
		zap.String("name", name),
		zap.Bool("allowed", result.Allowed),
		zap.String("reason", result.Reason),
	)
	
	return result, nil
}

// ListRegoHooks returns all registered Rego hooks
func (m *RegoHookManager) ListRegoHooks(ctx context.Context) []*RegoHook {
	hooks := make([]*RegoHook, 0, len(m.hooks))
	for _, hook := range m.hooks {
		hooks = append(hooks, hook)
	}
	return hooks
}

// GetRegoHook returns a specific Rego hook
func (m *RegoHookManager) GetRegoHook(ctx context.Context, name string) (*RegoHook, error) {
	hook, exists := m.hooks[name]
	if !exists {
		return nil, fmt.Errorf("hook %s not found", name)
	}
	return hook, nil
}

// EnableRegoHook enables a Rego hook
func (m *RegoHookManager) EnableRegoHook(ctx context.Context, name string) error {
	hook, exists := m.hooks[name]
	if !exists {
		return fmt.Errorf("hook %s not found", name)
	}
	
	hook.Enabled = true
	m.logger.Info("Enabled Rego hook", zap.String("name", name))
	
	return nil
}

// DisableRegoHook disables a Rego hook
func (m *RegoHookManager) DisableRegoHook(ctx context.Context, name string) error {
	hook, exists := m.hooks[name]
	if !exists {
		return fmt.Errorf("hook %s not found", name)
	}
	
	hook.Enabled = false
	m.logger.Info("Disabled Rego hook", zap.String("name", name))
	
	return nil
}

// LoadRegoHooksFromDirectory loads Rego hooks from a directory
func (m *RegoHookManager) LoadRegoHooksFromDirectory(ctx context.Context, dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", dir)
	}
	
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".rego") {
			return nil
		}
		
		// Read policy file
		policy, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read policy file %s: %w", path, err)
		}
		
		// Extract hook name from filename
		name := strings.TrimSuffix(filepath.Base(path), ".rego")
		
		// Create hook
		hook := &RegoHook{
			Name:        name,
			Description: fmt.Sprintf("Rego policy loaded from %s", path),
			Policy:      string(policy),
			Query:       "data.policy.allow",
			Enabled:     true,
			Timeout:     30,
			Logger:      m.logger,
		}
		
		// Register hook
		if err := m.RegisterRegoHook(ctx, hook); err != nil {
			return fmt.Errorf("failed to register hook %s: %w", name, err)
		}
		
		return nil
	})
}

// RegoResult represents the result of a Rego hook execution
type RegoResult struct {
	Allowed bool                   `json:"allowed"`
	Reason  string                 `json:"reason"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ToJSON returns the result as JSON
func (r *RegoResult) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// String returns a string representation of the result
func (r *RegoResult) String() string {
	if r.Allowed {
		return fmt.Sprintf("ALLOWED: %s", r.Reason)
	}
	return fmt.Sprintf("DENIED: %s", r.Reason)
}

// IsAllowed returns true if the result allows the action
func (r *RegoResult) IsAllowed() bool {
	return r.Allowed
}

// GetData returns additional data from the result
func (r *RegoResult) GetData(key string) interface{} {
	return r.Data[key]
}

// SetData sets additional data in the result
func (r *RegoResult) SetData(key string, value interface{}) {
	if r.Data == nil {
		r.Data = make(map[string]interface{})
	}
	r.Data[key] = value
}
