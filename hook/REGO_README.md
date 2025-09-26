# OPA/Rego Hooks

The `tykctl-go` hook system includes comprehensive support for Open Policy Agent (OPA) with Rego policies. This allows you to implement sophisticated authorization, validation, and decision-making logic using Rego policies.

## Features

- **Rego Policy Execution**: Execute OPA Rego policies as hooks
- **Policy Management**: Register, unregister, enable/disable Rego hooks
- **Directory Loading**: Load Rego policies from directories
- **Rich Results**: Detailed policy execution results with reasoning
- **Integration**: Seamless integration with builtin and external hooks
- **Example Policies**: Pre-built policies for common use cases

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/edsonmichaque/tykctl/pkg/tykctl-go/hook"
)

func main() {
    // Create hook manager
    manager := hook.New()
    ctx := context.Background()
    
    // Register a Rego hook
    authHook := &hook.RegoHook{
        Name:        "auth",
        Description: "Authentication and authorization",
        Policy: `
package policy

default allow := false

allow if {
    input.user.role == "admin"
}

allow if {
    input.user.permissions[_] == input.required_permission
}
`,
        Query:   "data.policy.allow",
        Enabled: true,
        Timeout: 30,
    }
    
    if err := manager.RegisterRegoHook(ctx, authHook); err != nil {
        panic(err)
    }
    
    // Execute the hook
    input := map[string]interface{}{
        "user": map[string]interface{}{
            "role":        "admin",
            "permissions": []string{"extension:install"},
        },
        "required_permission": "extension:install",
    }
    
    result, err := manager.ExecuteRegoHook(ctx, "auth", input)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Result: %s\n", result.String())
    fmt.Printf("Allowed: %t\n", result.IsAllowed())
}
```

## RegoHook Structure

```go
type RegoHook struct {
    Name        string                 // Hook name (required)
    Description string                 // Hook description
    Policy      string                 // Rego policy (required)
    Input       map[string]interface{} // Default input data
    Query       string                 // Rego query (default: "data.policy.allow")
    Enabled     bool                   // Whether hook is enabled
    Timeout     int                   // Execution timeout in seconds
    Logger      *zap.Logger           // Logger instance
}
```

## Example Policies

### Authentication & Authorization

```rego
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
```

### Extension Validation

```rego
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

# Provide detailed validation results
validation_errors := [
    "extension source not trusted"
] if input.extension.source != "github.com"

validation_errors := [
    "extension not verified"
] if not input.extension.verified
```

### Rate Limiting

```rego
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
```

### Resource Access Control

```rego
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
```

### Audit Logging

```rego
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
```

### Compliance Checking

```rego
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
```

## API Reference

### Registering Rego Hooks

```go
// Register a single Rego hook
err := manager.RegisterRegoHook(ctx, &hook.RegoHook{
    Name:        "my_hook",
    Description: "My custom Rego hook",
    Policy:      regoPolicy,
    Query:       "data.policy.allow",
    Enabled:     true,
    Timeout:     30,
})

// Load Rego hooks from directory
err := manager.LoadRegoHooksFromDirectory(ctx, "/path/to/rego/policies")
```

### Executing Rego Hooks

```go
// Execute a Rego hook
result, err := manager.ExecuteRegoHook(ctx, "my_hook", inputData)
if err != nil {
    // Handle error
}

// Check result
if result.IsAllowed() {
    // Action is allowed
} else {
    // Action is denied
    fmt.Printf("Reason: %s\n", result.Reason)
}

// Get additional data
if data := result.GetData("custom_field"); data != nil {
    // Use custom data
}
```

### Managing Rego Hooks

```go
// List all Rego hooks
hooks := manager.ListRegoHooks(ctx)

// Get specific Rego hook
hook, err := manager.GetRegoHook(ctx, "my_hook")

// Enable/disable Rego hook
err := manager.EnableRegoHook(ctx, "my_hook")
err := manager.DisableRegoHook(ctx, "my_hook")

// Unregister Rego hook
err := manager.UnregisterRegoHook(ctx, "my_hook")
```

## RegoResult Structure

```go
type RegoResult struct {
    Allowed bool                   `json:"allowed"`
    Reason  string                 `json:"reason"`
    Data    map[string]interface{} `json:"data,omitempty"`
}

// Methods
func (r *RegoResult) IsAllowed() bool
func (r *RegoResult) String() string
func (r *RegoResult) ToJSON() ([]byte, error)
func (r *RegoResult) GetData(key string) interface{}
func (r *RegoResult) SetData(key string, value interface{})
```

## Integration with Extension Lifecycle

Rego hooks can be integrated into the extension lifecycle:

```go
// Before installing an extension
func (i *Installer) InstallExtension(ctx context.Context, name string) error {
    // Execute validation Rego hook
    input := map[string]interface{}{
        "extension": map[string]interface{}{
            "name":    name,
            "source":  "github.com",
            "verified": true,
        },
    }
    
    result, err := i.hooks.ExecuteRegoHook(ctx, "validation", input)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    if !result.IsAllowed() {
        return fmt.Errorf("extension validation failed: %s", result.Reason)
    }
    
    // Continue with installation...
}
```

## Best Practices

### 1. Policy Design
- Use clear, descriptive package names
- Implement default deny policies
- Provide detailed reasoning in results
- Use structured input data

### 2. Error Handling
- Always check for execution errors
- Handle timeout scenarios
- Validate input data before execution

### 3. Performance
- Keep policies simple and efficient
- Use appropriate timeouts
- Cache frequently used policies

### 4. Security
- Validate all input data
- Use least privilege principles
- Audit policy changes

## Examples

See `rego_example_usage.go` for comprehensive usage examples including:
- Authentication and authorization
- Extension validation
- Rate limiting
- Resource access control
- Audit logging
- Compliance checking
- Custom policy creation

## Dependencies

The Rego hook system requires the OPA Go SDK:
```go
import "github.com/open-policy-agent/opa/rego"
```

This is automatically included when you import the hook package.
