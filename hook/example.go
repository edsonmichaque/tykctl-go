package hook

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Example demonstrates how to use the hook system
func Example() {
	// Create a builtin processor
	processor := NewBuiltinProcessor(nil)
	ctx := context.Background()

	// Register some builtin hooks
	err := processor.Register("before-install", func(ctx context.Context, data *Data) error {
		log.Printf("Before install hook: Installing extension %s", data.Extension)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = processor.Register("after-install", func(ctx context.Context, data *Data) error {
		log.Printf("After install hook: Extension %s installed", data.Extension)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = processor.Register("before-run", func(ctx context.Context, data *Data) error {
		log.Printf("Before run hook: Starting extension %s", data.Extension)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = processor.Register("after-run", func(ctx context.Context, data *Data) error {
		log.Printf("After run hook: Extension %s completed", data.Extension)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Execute hooks
	hookData := NewData("before-install", "my-extension").
		WithMetadataMap(map[string]interface{}{
			"version": "1.0.0",
			"author":  "example",
		})

	// Execute before install hooks
	if err := processor.Execute(ctx, "before-install", hookData); err != nil {
		log.Printf("Before install hooks failed: %v", err)
		return
	}

	// Simulate installation work
	time.Sleep(100 * time.Millisecond)

	// Execute after install hooks
	hookData.Type = "after-install"
	if err := processor.Execute(ctx, "after-install", hookData); err != nil {
		log.Printf("After install hooks failed: %v", err)
		return
	}

	// Execute before run hooks
	hookData.Type = "before-run"
	if err := processor.Execute(ctx, "before-run", hookData); err != nil {
		log.Printf("Before run hooks failed: %v", err)
		return
	}

	// Simulate running work
	time.Sleep(100 * time.Millisecond)

	// Execute after run hooks
	hookData.Type = "after-run"
	if err := processor.Execute(ctx, "after-run", hookData); err != nil {
		log.Printf("After run hooks failed: %v", err)
		return
	}

	fmt.Println("All hooks executed successfully!")
}

// ExamplePolicyProcessor demonstrates how to use the policy processor
func ExamplePolicyProcessor() {
	// Create a policy processor
	processor := NewPolicyProcessor(nil, "/path/to/policies")
	ctx := context.Background()

	// Execute policy hooks
	hookData := NewData("before-install", "my-extension").
		WithMetadata("version", "1.0.0")

	err := processor.Execute(ctx, "before-install", hookData)
	if err != nil {
		log.Printf("Policy execution failed: %v", err)
	}

	// Get the underlying Rego executor for advanced usage
	regoExecutor := processor.GetRegoExecutor()
	if regoExecutor != nil {
		log.Println("Rego executor is available")
	}
}

// ExampleScriptProcessor demonstrates how to use the script processor
func ExampleScriptProcessor() {
	// Create a script processor
	processor := NewScriptProcessor(nil, "/path/to/scripts")
	ctx := context.Background()

	// Execute script hooks
	hookData := NewData("before-install", "my-extension").
		WithMetadata("version", "1.0.0")

	err := processor.Execute(ctx, "before-install", hookData)
	if err != nil {
		log.Printf("Script execution failed: %v", err)
	}

	// Get the underlying script executor for advanced usage
	scriptExecutor := processor.GetScriptExecutor()
	if scriptExecutor != nil {
		log.Println("Script executor is available")
	}
}
