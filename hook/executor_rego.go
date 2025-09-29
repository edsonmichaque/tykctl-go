package hook

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

// RegoExecutor executes Rego policy hooks.
type RegoExecutor struct {
	policyDir  string
	logger     *zap.Logger
	discovered map[Type][]string // Cache for discovered policies
}

// NewRegoExecutor creates a new Rego hook executor.
func NewRegoExecutor(logger *zap.Logger, policyDir string) *RegoExecutor {
	return &RegoExecutor{
		policyDir:  policyDir,
		logger:     logger,
		discovered: make(map[Type][]string),
	}
}

// Execute executes Rego policy hooks for a specific type.
func (e *RegoExecutor) Execute(ctx context.Context, hookType Type, data *Data) error {
	if e.policyDir == "" {
		return nil
	}

	// Look for Rego policy files
	policyFiles, err := e.findPolicyFiles(ctx, hookType)
	if err != nil {
		return fmt.Errorf("failed to find policy files: %w", err)
	}

	if len(policyFiles) == 0 {
		return nil
	}

	// Execute each policy file
	for _, policyFile := range policyFiles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := e.executePolicyFile(ctx, policyFile, data); err != nil {
			if e.logger != nil {
				e.logger.Error("Policy execution failed",
					zap.String("policy_file", policyFile),
					zap.String("hook_type", string(hookType)),
					zap.Error(err),
				)
			}
			return err
		}
	}

	return nil
}

// findPolicyFiles finds Rego policy files for a specific hook type.
func (e *RegoExecutor) findPolicyFiles(ctx context.Context, hookType Type) ([]string, error) {
	// Check cache first
	if policies, exists := e.discovered[hookType]; exists {
		return policies, nil
	}

	var policyFiles []string

	// Check if hookType is a direct file or directory
	policyPath := filepath.Join(e.policyDir, string(hookType))

	info, err := os.Stat(policyPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Hook type doesn't exist, return empty list
			e.discovered[hookType] = []string{}
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to stat policy path %s: %w", policyPath, err)
	}

	if info.IsDir() {
		// Hook type is a directory, process its entries in lexicographic order
		policyFiles, err = e.processPolicyDirectory(ctx, policyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to process policy directory %s: %w", policyPath, err)
		}
	} else {
		// Hook type is a file, check if it's a valid Rego policy
		if e.isValidRegoPolicy(info.Name(), policyPath) {
			policyFiles = []string{policyPath}
		}
	}

	// Cache the discovered policies
	e.discovered[hookType] = policyFiles

	if e.logger != nil {
		e.logger.Debug("Discovered policies",
			zap.String("hook_type", string(hookType)),
			zap.String("policy_path", policyPath),
			zap.Bool("is_directory", info.IsDir()),
			zap.Int("count", len(policyFiles)),
			zap.Strings("policies", policyFiles),
		)
	}

	return policyFiles, nil
}

// processPolicyDirectory processes a policy directory and returns policies in lexicographic order.
func (e *RegoExecutor) processPolicyDirectory(ctx context.Context, dirPath string) ([]string, error) {
	var policies []string

	// Read directory entries
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// Process entries in lexicographic order (ReadDir already sorts by name)
	for _, entry := range entries {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if entry.IsDir() {
			// Skip subdirectories for now
			continue
		}

		entryPath := filepath.Join(dirPath, entry.Name())

		// Check if it's a valid Rego policy
		if e.isValidRegoPolicy(entry.Name(), entryPath) {
			policies = append(policies, entryPath)
		}
	}

	return policies, nil
}

// isValidRegoPolicy checks if a file is a valid Rego policy.
func (e *RegoExecutor) isValidRegoPolicy(filename, fullPath string) bool {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".rego" {
		return false
	}

	// Check if file exists and is readable
	info, err := os.Stat(fullPath)
	if err != nil {
		return false
	}

	// Check if it's a file (not directory)
	if info.IsDir() {
		return false
	}

	// Basic check for Rego content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return false
	}

	contentStr := string(content)

	// Check for basic Rego keywords
	regoKeywords := []string{"package", "import", "allow", "deny", "default"}
	hasKeyword := false
	for _, keyword := range regoKeywords {
		if strings.Contains(contentStr, keyword) {
			hasKeyword = true
			break
		}
	}

	return hasKeyword
}

// executePolicyFile executes a single Rego policy file.
func (e *RegoExecutor) executePolicyFile(ctx context.Context, policyFile string, data *Data) error {
	// Read policy file
	policyContent, err := os.ReadFile(policyFile)
	if err != nil {
		return fmt.Errorf("failed to read policy file: %w", err)
	}

	// For now, we'll just log the policy execution
	// In a real implementation, you would use a Rego engine like OPA
	if e.logger != nil {
		e.logger.Debug("Executing Rego policy",
			zap.String("policy_file", policyFile),
			zap.String("hook_type", string(data.Type)),
			zap.String("extension", data.Extension),
		)
	}

	// TODO: Implement actual Rego policy execution
	// This would involve:
	// 1. Loading the policy into an OPA instance
	// 2. Evaluating the policy with the hook data
	// 3. Handling the policy decision

	_ = policyContent // Avoid unused variable warning

	return nil
}

// discoverPolicies discovers Rego policy files for specific hook types or all hook types.
// If no hook types are provided, it discovers all policies. Otherwise, it discovers policies for the specified types.
func (e *RegoExecutor) discoverPolicies(ctx context.Context, hookTypes ...Type) (map[Type][]string, error) {
	allPolicies := make(map[Type][]string)

	if e.policyDir == "" {
		return allPolicies, nil
	}

	// If specific hook types requested, check cache first
	if len(hookTypes) > 0 {
		// Check if all requested types are in cache
		allCached := true
		for _, hookType := range hookTypes {
			if policies, exists := e.discovered[hookType]; exists {
				allPolicies[hookType] = policies
			} else {
				allCached = false
				break
			}
		}
		if allCached {
			return allPolicies, nil
		}
	}

	// Read the policy directory to find hook types
	entries, err := os.ReadDir(e.policyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy directory %s: %w", e.policyDir, err)
	}

	// Process each entry as a potential hook type
	for _, entry := range entries {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		entryName := entry.Name()

		// Skip hidden files and directories
		if strings.HasPrefix(entryName, ".") {
			continue
		}

		// Convert entry name to hook type
		currentHookType := Type(entryName)

		// If specific hook types requested, only process those types
		if len(hookTypes) > 0 {
			found := false
			for _, hookType := range hookTypes {
				if currentHookType == hookType {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Discover policies for this hook type
		policies, err := e.findPolicyFiles(ctx, currentHookType)
		if err != nil {
			if e.logger != nil {
				e.logger.Warn("Failed to discover policies for hook type",
					zap.String("hook_type", entryName),
					zap.Error(err),
				)
			}
			continue
		}

		if len(policies) > 0 {
			allPolicies[currentHookType] = policies
		}
	}

	// Update cache
	e.discovered = allPolicies

	if e.logger != nil {
		totalPolicies := 0
		for _, policies := range allPolicies {
			totalPolicies += len(policies)
		}
		e.logger.Debug("Discovered policies",
			zap.String("policy_dir", e.policyDir),
			zap.Strings("requested_types", hookTypesToStrings(hookTypes)),
			zap.Int("total_policies", totalPolicies),
			zap.Int("hook_types", len(allPolicies)),
		)
	}

	return allPolicies, nil
}

// hookTypesToStrings converts a slice of Type to a slice of string for logging.
func hookTypesToStrings(hookTypes []Type) []string {
	if len(hookTypes) == 0 {
		return []string{"all"}
	}
	strings := make([]string, len(hookTypes))
	for i, hookType := range hookTypes {
		strings[i] = string(hookType)
	}
	return strings
}

// listDiscoveredPolicies returns all discovered policies for a hook type.
func (e *RegoExecutor) listDiscoveredPolicies(hookType Type) []string {
	if policies, exists := e.discovered[hookType]; exists {
		return policies
	}
	return []string{}
}

// refreshDiscovery refreshes the policy discovery cache.
func (e *RegoExecutor) refreshDiscovery() error {
	_, err := e.discoverPolicies(context.Background())
	return err
}

// getPolicyDirectory returns the policy directory path.
func (e *RegoExecutor) getPolicyDirectory() string {
	return e.policyDir
}

// countDiscoveredPolicies returns the count of discovered policies for a hook type.
func (e *RegoExecutor) countDiscoveredPolicies(hookType Type) int {
	if policies, exists := e.discovered[hookType]; exists {
		return len(policies)
	}
	return 0
}
