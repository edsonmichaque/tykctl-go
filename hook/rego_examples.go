package hook

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// ExampleRegoPolicies provides example Rego policies for common use cases
var ExampleRegoPolicies = map[string]string{
	"auth": `
package policy

import rego.v1

# Default deny
default allow := false

# Allow if user has admin role
allow if {
	input.user.role == "admin"
}

# Allow if user has specific permission
allow if {
	input.user.permissions[_] == input.required_permission
}

# Allow if user is owner of resource
allow if {
	input.user.id == input.resource.owner_id
}

# Deny if user is blocked
allow if {
	not input.user.blocked
}

# Provide reason for decision
reason := "access denied" if not allow
reason := "access granted" if allow
`,

	"validation": `
package policy

import rego.v1

# Validate extension installation
default allow := false

allow if {
	# Check if extension is from trusted source
	input.extension.source == "github.com"
	input.extension.verified == true
}

allow if {
	# Check if extension has required metadata
	input.extension.name != ""
	input.extension.version != ""
	input.extension.description != ""
}

# Validate configuration
allow if {
	input.config.timeout > 0
	input.config.timeout <= 300
}
`,

	"rate_limiting": `
package policy

import rego.v1

# Rate limiting policy
default allow := false

# Allow if within rate limits
allow if {
	input.requests_count < input.rate_limit.max_requests
	input.time_window_remaining > 0
}

# Allow if user has premium subscription
allow if {
	input.user.subscription == "premium"
}

# Calculate remaining requests
remaining_requests := input.rate_limit.max_requests - input.requests_count

# Calculate reset time
reset_time := input.current_time + input.rate_limit.window_seconds
`,

	"resource_access": `
package policy

import rego.v1

# Resource access control
default allow := false

# Allow if user has access to resource
allow if {
	input.user.resources[_] == input.resource.id
}

# Allow if resource is public
allow if {
	input.resource.visibility == "public"
}

# Allow if user is in resource's allowed groups
allow if {
	input.user.groups[_] in input.resource.allowed_groups
}

# Deny if resource is restricted
allow if {
	not input.resource.restricted
}

# Provide access level
access_level := "read" if allow and input.action == "read"
access_level := "write" if allow and input.action == "write"
access_level := "admin" if allow and input.user.role == "admin"
`,

	"audit": `
package policy

import rego.v1

# Audit logging policy
default allow := true

# Always allow audit logging
allow := true

# Determine audit level based on action
audit_level := "info" if input.action in ["read", "list"]
audit_level := "warn" if input.action in ["create", "update"]
audit_level := "error" if input.action in ["delete", "destroy"]

# Determine if action should be audited
should_audit := true if input.action in ["create", "update", "delete", "destroy"]
should_audit := false if input.action in ["read", "list"]

# Extract sensitive fields that should be masked
sensitive_fields := ["password", "token", "secret", "key", "credential"]
`,

	"compliance": `
package policy

import rego.v1

# Compliance checking policy
default allow := false

# Allow if all compliance checks pass
allow if {
	compliance_checks.passed
}

# Check data retention compliance
data_retention_compliant if {
	input.data.retention_days <= input.policy.max_retention_days
}

# Check data encryption compliance
encryption_compliant if {
	input.data.encrypted == true
	input.data.encryption_algorithm in ["AES-256", "ChaCha20"]
}

# Check access logging compliance
access_logging_compliant if {
	input.data.access_logged == true
	input.data.log_retention_days >= input.policy.min_log_retention_days
}

# Overall compliance check
compliance_checks := {
	"passed": data_retention_compliant and encryption_compliant and access_logging_compliant,
	"data_retention": data_retention_compliant,
	"encryption": encryption_compliant,
	"access_logging": access_logging_compliant,
}
`,
}

// ExampleRegoHookUsage demonstrates how to use Rego hooks
func ExampleRegoHookUsage() {
	ctx := context.Background()
	logger := zap.NewNop() // Use a real logger in production

	// Create Rego hook manager
	regoManager := NewRegoHookManager(logger)

	// Register authentication hook
	authHook := &RegoHook{
		Name:        "auth",
		Description: "Authentication and authorization hook",
		Policy:      ExampleRegoPolicies["auth"],
		Query:       "data.policy.allow",
		Enabled:     true,
		Timeout:     30,
		Logger:      logger,
		Input: map[string]interface{}{
			"required_permission": "extension:install",
		},
	}

	if err := regoManager.RegisterRegoHook(ctx, authHook); err != nil {
		fmt.Printf("Failed to register auth hook: %v\n", err)
		return
	}

	// Register validation hook
	validationHook := &RegoHook{
		Name:        "validation",
		Description: "Extension validation hook",
		Policy:      ExampleRegoPolicies["validation"],
		Query:       "data.policy.allow",
		Enabled:     true,
		Timeout:     30,
		Logger:      logger,
	}

	if err := regoManager.RegisterRegoHook(ctx, validationHook); err != nil {
		fmt.Printf("Failed to register validation hook: %v\n", err)
		return
	}

	// Test authentication hook
	authInput := map[string]interface{}{
		"user": map[string]interface{}{
			"id":          "user123",
			"role":        "admin",
			"permissions": []string{"extension:install", "extension:delete"},
			"blocked":     false,
		},
		"resource": map[string]interface{}{
			"owner_id": "user123",
		},
	}

	authResult, err := regoManager.ExecuteRegoHook(ctx, "auth", authInput)
	if err != nil {
		fmt.Printf("Auth hook execution failed: %v\n", err)
		return
	}

	fmt.Printf("Auth result: %s\n", authResult.String())

	// Test validation hook
	validationInput := map[string]interface{}{
		"extension": map[string]interface{}{
			"name":        "tykctl-example",
			"version":     "1.0.0",
			"description": "Example extension",
			"source":      "github.com",
			"verified":    true,
		},
		"config": map[string]interface{}{
			"timeout": 30,
		},
	}

	validationResult, err := regoManager.ExecuteRegoHook(ctx, "validation", validationInput)
	if err != nil {
		fmt.Printf("Validation hook execution failed: %v\n", err)
		return
	}

	fmt.Printf("Validation result: %s\n", validationResult.String())

	// List all hooks
	hooks := regoManager.ListRegoHooks(ctx)
	fmt.Printf("Registered hooks: %d\n", len(hooks))
	for _, hook := range hooks {
		fmt.Printf("- %s: %s\n", hook.Name, hook.Description)
	}
}

// CreateExampleRegoHooks creates example Rego hooks for common scenarios
func CreateExampleRegoHooks(ctx context.Context, manager *Manager) error {
	// Authentication hook
	authHook := &RegoHook{
		Name:        "auth",
		Description: "Authentication and authorization",
		Policy:      ExampleRegoPolicies["auth"],
		Query:       "data.policy.allow",
		Enabled:     true,
		Timeout:     30,
	}

	if err := manager.RegisterRegoHook(ctx, authHook); err != nil {
		return fmt.Errorf("failed to register auth hook: %w", err)
	}

	// Validation hook
	validationHook := &RegoHook{
		Name:        "validation",
		Description: "Extension validation",
		Policy:      ExampleRegoPolicies["validation"],
		Query:       "data.policy.allow",
		Enabled:     true,
		Timeout:     30,
	}

	if err := manager.RegisterRegoHook(ctx, validationHook); err != nil {
		return fmt.Errorf("failed to register validation hook: %w", err)
	}

	// Rate limiting hook
	rateLimitHook := &RegoHook{
		Name:        "rate_limiting",
		Description: "Rate limiting and throttling",
		Policy:      ExampleRegoPolicies["rate_limiting"],
		Query:       "data.policy.allow",
		Enabled:     true,
		Timeout:     30,
	}

	if err := manager.RegisterRegoHook(ctx, rateLimitHook); err != nil {
		return fmt.Errorf("failed to register rate limiting hook: %w", err)
	}

	// Resource access hook
	resourceAccessHook := &RegoHook{
		Name:        "resource_access",
		Description: "Resource access control",
		Policy:      ExampleRegoPolicies["resource_access"],
		Query:       "data.policy.allow",
		Enabled:     true,
		Timeout:     30,
	}

	if err := manager.RegisterRegoHook(ctx, resourceAccessHook); err != nil {
		return fmt.Errorf("failed to register resource access hook: %w", err)
	}

	// Audit hook
	auditHook := &RegoHook{
		Name:        "audit",
		Description: "Audit logging and compliance",
		Policy:      ExampleRegoPolicies["audit"],
		Query:       "data.policy.allow",
		Enabled:     true,
		Timeout:     30,
	}

	if err := manager.RegisterRegoHook(ctx, auditHook); err != nil {
		return fmt.Errorf("failed to register audit hook: %w", err)
	}

	// Compliance hook
	complianceHook := &RegoHook{
		Name:        "compliance",
		Description: "Compliance checking",
		Policy:      ExampleRegoPolicies["compliance"],
		Query:       "data.policy.allow",
		Enabled:     true,
		Timeout:     30,
	}

	if err := manager.RegisterRegoHook(ctx, complianceHook); err != nil {
		return fmt.Errorf("failed to register compliance hook: %w", err)
	}

	return nil
}
