package hook

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// PolicyDispatcher handles only Rego policy hooks with policy-specific functionality.
type PolicyDispatcher struct {
	regoExecutor *RegoExecutor
	validator    Validator
	logger       *zap.Logger
}

// NewPolicyDispatcher creates a new policy-only dispatcher.
func NewPolicyDispatcher(logger *zap.Logger, policyDir string) *PolicyDispatcher {
	return &PolicyDispatcher{
		regoExecutor: NewRegoExecutor(logger, policyDir),
		validator:    NewPolicyValidator(policyDir),
		logger:       logger,
	}
}

// Execute executes only Rego policy hooks.
func (pp *PolicyDispatcher) Execute(ctx context.Context, hookType Type, data *Data) error {
	// Validate hook data
	if err := pp.validator.Validate(data); err != nil {
		return fmt.Errorf("hook validation failed: %w", err)
	}

	// Execute Rego policy hooks
	if err := pp.regoExecutor.Execute(ctx, hookType, data); err != nil {
		if pp.logger != nil {
			pp.logger.Error("Rego policy execution failed",
				zap.String("hook_type", string(hookType)),
				zap.String("extension", data.Extension),
				zap.Error(err),
			)
		}
		return fmt.Errorf("rego policy execution failed: %w", err)
	}

	return nil
}

// GetRegoExecutor returns the underlying Rego executor for advanced usage.
func (pp *PolicyDispatcher) GetRegoExecutor() *RegoExecutor {
	return pp.regoExecutor
}

// ValidatePolicy validates a policy file without executing hooks.
func (pp *PolicyDispatcher) ValidatePolicy(ctx context.Context, policyFile string) error {
	if pp.regoExecutor == nil {
		return fmt.Errorf("rego executor not available")
	}

	// This would need to be implemented in RegoExecutor
	// For now, return a placeholder
	return fmt.Errorf("policy validation not yet implemented")
}

// ListPolicies returns a list of available policy files for a given hook type.
// This method adds dispatcher-level validation and error handling.
func (pp *PolicyDispatcher) ListPolicies(ctx context.Context, hookType Type) ([]string, error) {
	if pp.regoExecutor == nil {
		return nil, fmt.Errorf("rego executor not available")
	}

	// Validate hook type before listing policies
	if hookType == "" {
		return nil, fmt.Errorf("hook type cannot be empty")
	}

	allPolicies, err := pp.regoExecutor.discoverPolicies(ctx, hookType)
	if err != nil {
		return nil, fmt.Errorf("failed to discover policies: %w", err)
	}
	policies := allPolicies[hookType]

	if pp.logger != nil {
		pp.logger.Debug("Listed policies for hook type",
			zap.String("hook_type", string(hookType)),
			zap.Int("count", len(policies)),
		)
	}

	return policies, nil
}

// DiscoverAllPolicies discovers all available policy files.
// This method adds dispatcher-level validation and error handling.
func (pp *PolicyDispatcher) DiscoverAllPolicies() (map[Type][]string, error) {
	if pp.regoExecutor == nil {
		return nil, fmt.Errorf("rego executor not available")
	}

	allPolicies, err := pp.regoExecutor.discoverPolicies(context.Background())
	if err != nil {
		if pp.logger != nil {
			pp.logger.Error("Failed to discover all policies",
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("failed to discover all policies: %w", err)
	}

	if pp.logger != nil {
		totalPolicies := 0
		for _, policies := range allPolicies {
			totalPolicies += len(policies)
		}
		pp.logger.Debug("Discovered all policies",
			zap.Int("total_policies", totalPolicies),
			zap.Int("hook_types", len(allPolicies)),
		)
	}

	return allPolicies, nil
}

// CountPolicies returns the count of available policy files for a given hook type.
// This method adds dispatcher-level validation and error handling.
func (pp *PolicyDispatcher) CountPolicies(ctx context.Context, hookType Type) (int, error) {
	if pp.regoExecutor == nil {
		return 0, fmt.Errorf("rego executor not available")
	}

	// Validate hook type before counting policies
	if hookType == "" {
		return 0, fmt.Errorf("hook type cannot be empty")
	}

	allPolicies, err := pp.regoExecutor.discoverPolicies(ctx, hookType)
	if err != nil {
		return 0, fmt.Errorf("failed to discover policies: %w", err)
	}
	policies := allPolicies[hookType]

	if pp.logger != nil {
		pp.logger.Debug("Counted policies for hook type",
			zap.String("hook_type", string(hookType)),
			zap.Int("count", len(policies)),
		)
	}

	return len(policies), nil
}

// GetPolicyDirectory returns the policy directory path.
// This method adds dispatcher-level validation and error handling.
func (pp *PolicyDispatcher) GetPolicyDirectory() (string, error) {
	if pp.regoExecutor == nil {
		return "", fmt.Errorf("rego executor not available")
	}

	policyDir := pp.regoExecutor.getPolicyDirectory()

	if pp.logger != nil {
		pp.logger.Debug("Retrieved policy directory",
			zap.String("policy_dir", policyDir),
		)
	}

	return policyDir, nil
}
