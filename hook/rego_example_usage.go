package main

import (
	"context"
	"fmt"
	"log"

	"github.com/edsonmichaque/tykctl-go/hook"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	defer logger.Sync()

	// Create hook manager with logger
	manager := hook.NewWithLogger("/tmp/hooks", logger)
	ctx := context.Background()

	// Create example Rego hooks
	if err := hook.CreateExampleRegoHooks(ctx, manager); err != nil {
		log.Fatal("Failed to create example hooks:", err)
	}

	// Example 1: Authentication Hook
	fmt.Println("=== Authentication Hook Example ===")
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
		"required_permission": "extension:install",
	}

	authResult, err := manager.ExecuteRegoHook(ctx, "auth", authInput)
	if err != nil {
		log.Printf("Auth hook execution failed: %v", err)
	} else {
		fmt.Printf("Auth result: %s\n", authResult.String())
		fmt.Printf("Allowed: %t\n", authResult.IsAllowed())
		if reason := authResult.GetData("reason"); reason != nil {
			fmt.Printf("Reason: %v\n", reason)
		}
	}

	// Example 2: Validation Hook
	fmt.Println("\n=== Validation Hook Example ===")
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

	validationResult, err := manager.ExecuteRegoHook(ctx, "validation", validationInput)
	if err != nil {
		log.Printf("Validation hook execution failed: %v", err)
	} else {
		fmt.Printf("Validation result: %s\n", validationResult.String())
		fmt.Printf("Allowed: %t\n", validationResult.IsAllowed())
	}

	// Example 3: Rate Limiting Hook
	fmt.Println("\n=== Rate Limiting Hook Example ===")
	rateLimitInput := map[string]interface{}{
		"requests_count": 5,
		"rate_limit": map[string]interface{}{
			"max_requests":    10,
			"window_seconds":  60,
		},
		"time_window_remaining": 45,
		"current_time":          1234567890,
		"user": map[string]interface{}{
			"subscription": "basic",
		},
	}

	rateLimitResult, err := manager.ExecuteRegoHook(ctx, "rate_limiting", rateLimitInput)
	if err != nil {
		log.Printf("Rate limiting hook execution failed: %v", err)
	} else {
		fmt.Printf("Rate limiting result: %s\n", rateLimitResult.String())
		fmt.Printf("Allowed: %t\n", rateLimitResult.IsAllowed())
		if remaining := rateLimitResult.GetData("remaining_requests"); remaining != nil {
			fmt.Printf("Remaining requests: %v\n", remaining)
		}
	}

	// Example 4: Resource Access Hook
	fmt.Println("\n=== Resource Access Hook Example ===")
	resourceInput := map[string]interface{}{
		"user": map[string]interface{}{
			"id":     "user123",
			"role":   "user",
			"groups": []string{"developers", "testers"},
			"resources": []string{"resource1", "resource2"},
		},
		"resource": map[string]interface{}{
			"id":              "resource1",
			"visibility":      "private",
			"allowed_groups":  []string{"developers"},
			"restricted":      false,
		},
		"action": "read",
	}

	resourceResult, err := manager.ExecuteRegoHook(ctx, "resource_access", resourceInput)
	if err != nil {
		log.Printf("Resource access hook execution failed: %v", err)
	} else {
		fmt.Printf("Resource access result: %s\n", resourceResult.String())
		fmt.Printf("Allowed: %t\n", resourceResult.IsAllowed())
		if accessLevel := resourceResult.GetData("access_level"); accessLevel != nil {
			fmt.Printf("Access level: %v\n", accessLevel)
		}
	}

	// Example 5: Audit Hook
	fmt.Println("\n=== Audit Hook Example ===")
	auditInput := map[string]interface{}{
		"action": "create",
		"user": map[string]interface{}{
			"id": "user123",
		},
		"resource": map[string]interface{}{
			"type": "extension",
			"id":   "tykctl-example",
		},
	}

	auditResult, err := manager.ExecuteRegoHook(ctx, "audit", auditInput)
	if err != nil {
		log.Printf("Audit hook execution failed: %v", err)
	} else {
		fmt.Printf("Audit result: %s\n", auditResult.String())
		fmt.Printf("Allowed: %t\n", auditResult.IsAllowed())
		if auditLevel := auditResult.GetData("audit_level"); auditLevel != nil {
			fmt.Printf("Audit level: %v\n", auditLevel)
		}
		if shouldAudit := auditResult.GetData("should_audit"); shouldAudit != nil {
			fmt.Printf("Should audit: %v\n", shouldAudit)
		}
	}

	// Example 6: Compliance Hook
	fmt.Println("\n=== Compliance Hook Example ===")
	complianceInput := map[string]interface{}{
		"data": map[string]interface{}{
			"retention_days":     30,
			"encrypted":          true,
			"encryption_algorithm": "AES-256",
			"access_logged":      true,
			"log_retention_days": 90,
		},
		"policy": map[string]interface{}{
			"max_retention_days":     365,
			"min_log_retention_days": 30,
		},
	}

	complianceResult, err := manager.ExecuteRegoHook(ctx, "compliance", complianceInput)
	if err != nil {
		log.Printf("Compliance hook execution failed: %v", err)
	} else {
		fmt.Printf("Compliance result: %s\n", complianceResult.String())
		fmt.Printf("Allowed: %t\n", complianceResult.IsAllowed())
		if complianceChecks := complianceResult.GetData("compliance_checks"); complianceChecks != nil {
			fmt.Printf("Compliance checks: %v\n", complianceChecks)
		}
	}

	// List all registered Rego hooks
	fmt.Println("\n=== Registered Rego Hooks ===")
	hooks := manager.ListRegoHooks(ctx)
	for _, hook := range hooks {
		fmt.Printf("- %s: %s (enabled: %t)\n", hook.Name, hook.Description, hook.Enabled)
	}

	// Example 7: Custom Rego Hook
	fmt.Println("\n=== Custom Rego Hook Example ===")
	customHook := &hook.RegoHook{
		Name:        "custom_validation",
		Description: "Custom validation for specific use case",
		Policy: `
package policy

import rego.v1

default allow := false

# Allow if all conditions are met
allow if {
	input.extension.name != ""
	input.extension.version != ""
	input.extension.license in ["MIT", "Apache-2.0", "BSD-3-Clause"]
	input.extension.size_mb < 100
}

# Provide detailed validation results
validation_errors := [
	"extension name is required"
] if input.extension.name == ""

validation_errors := [
	"extension version is required"
] if input.extension.version == ""

validation_errors := [
	"unsupported license"
] if not input.extension.license in ["MIT", "Apache-2.0", "BSD-3-Clause"]

validation_errors := [
	"extension too large"
] if input.extension.size_mb >= 100
`,
		Query:   "data.policy.allow",
		Enabled: true,
		Timeout: 30,
	}

	if err := manager.RegisterRegoHook(ctx, customHook); err != nil {
		log.Printf("Failed to register custom hook: %v", err)
	} else {
		customInput := map[string]interface{}{
			"extension": map[string]interface{}{
				"name":     "tykctl-custom",
				"version":  "2.0.0",
				"license":  "MIT",
				"size_mb":  50,
			},
		}

		customResult, err := manager.ExecuteRegoHook(ctx, "custom_validation", customInput)
		if err != nil {
			log.Printf("Custom hook execution failed: %v", err)
		} else {
			fmt.Printf("Custom validation result: %s\n", customResult.String())
			fmt.Printf("Allowed: %t\n", customResult.IsAllowed())
			if errors := customResult.GetData("validation_errors"); errors != nil {
				fmt.Printf("Validation errors: %v\n", errors)
			}
		}
	}

	fmt.Println("\n=== Rego Hooks Example Completed ===")
}
