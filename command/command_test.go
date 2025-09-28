package command

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	if cmd == nil {
		t.Fatal("New() returned nil")
	}
	
	if cmd.Use != "test" {
		t.Errorf("Expected Use='test', got '%s'", cmd.Use)
	}
	
	if cmd.Short != "Test command" {
		t.Errorf("Expected Short='Test command', got '%s'", cmd.Short)
	}
}

func TestNewWithLong(t *testing.T) {
	longDesc := "This is a long description for the test command"
	cmd := NewWithLong("test", "Test command", longDesc, func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	if cmd == nil {
		t.Fatal("NewWithLong() returned nil")
	}
	
	if cmd.Use != "test" {
		t.Errorf("Expected Use='test', got '%s'", cmd.Use)
	}
	
	if cmd.Short != "Test command" {
		t.Errorf("Expected Short='Test command', got '%s'", cmd.Short)
	}
	
	if cmd.Long != longDesc {
		t.Errorf("Expected Long='%s', got '%s'", longDesc, cmd.Long)
	}
}

func TestSetLogger(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	result := cmd.SetLogger(logger)
	
	if result != cmd {
		t.Error("SetLogger should return the same command instance")
	}
	
	if cmd.logger != logger {
		t.Error("Logger was not set correctly")
	}
}

func TestSetContext(t *testing.T) {
	ctx := context.Background()
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	result := cmd.SetContext(ctx)
	
	if result != cmd {
		t.Error("SetContext should return the same command instance")
	}
	
	if cmd.context != ctx {
		t.Error("Context was not set correctly")
	}
}

func TestGetLogger(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	// Initially should be nil
	if cmd.GetLogger() != nil {
		t.Error("Initial logger should be nil")
	}
	
	// After setting logger
	cmd.SetLogger(logger)
	if cmd.GetLogger() != logger {
		t.Error("GetLogger should return the set logger")
	}
}

func TestGetContext(t *testing.T) {
	ctx := context.Background()
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	// Initially should be nil
	if cmd.GetContext() != nil {
		t.Error("Initial context should be nil")
	}
	
	// After setting context
	cmd.SetContext(ctx)
	if cmd.GetContext() != ctx {
		t.Error("GetContext should return the set context")
	}
}

func TestCommandExecution(t *testing.T) {
	executed := false
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		executed = true
		return nil
	})
	
	// Execute the command
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}
	
	if !executed {
		t.Error("Command run function was not executed")
	}
}

func TestCommandWithArgs(t *testing.T) {
	var receivedArgs []string
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		receivedArgs = args
		return nil
	})
	
	testArgs := []string{"arg1", "arg2", "arg3"}
	cmd.SetArgs(testArgs)
	err := cmd.Execute()
	
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}
	
	if len(receivedArgs) != len(testArgs) {
		t.Errorf("Expected %d args, got %d", len(testArgs), len(receivedArgs))
	}
	
	for i, arg := range testArgs {
		if receivedArgs[i] != arg {
			t.Errorf("Expected arg[%d]='%s', got '%s'", i, arg, receivedArgs[i])
		}
	}
}

func TestCommandWithLoggerAndContext(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var receivedLogger *zap.Logger
	var receivedContext context.Context
	
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		// Access logger and context from the command
		if cmdLogger := cmd.Context().Value("logger"); cmdLogger != nil {
			receivedLogger = cmdLogger.(*zap.Logger)
		}
		receivedContext = cmd.Context()
		return nil
	})
	
	cmd.SetLogger(logger)
	cmd.SetContext(ctx)
	
	// Add logger to context for testing
	ctxWithLogger := context.WithValue(ctx, "logger", logger)
	cmd.SetContext(ctxWithLogger)
	
	// Set the context on the cobra command
	cmd.Command.SetContext(ctxWithLogger)
	
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}
	
	if receivedLogger != logger {
		t.Error("Logger was not accessible in command execution")
	}
	
	if receivedContext != ctxWithLogger {
		t.Error("Context was not accessible in command execution")
	}
}

func TestCommandChaining(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()
	
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		return nil
	}).SetLogger(logger).SetContext(ctx)
	
	if cmd.logger != logger {
		t.Error("Chained SetLogger did not work")
	}
	
	if cmd.context != ctx {
		t.Error("Chained SetContext did not work")
	}
}

func TestCommandFlags(t *testing.T) {
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	// Add a flag
	cmd.Flags().String("test-flag", "default", "Test flag")
	
	// Set flag value
	cmd.SetArgs([]string{"--test-flag", "test-value"})
	
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}
	
	// Check flag value
	flagValue, err := cmd.Flags().GetString("test-flag")
	if err != nil {
		t.Errorf("Failed to get flag value: %v", err)
	}
	
	if flagValue != "test-value" {
		t.Errorf("Expected flag value 'test-value', got '%s'", flagValue)
	}
}

func TestCommandSubcommands(t *testing.T) {
	parentCmd := New("parent", "Parent command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	childCmd := New("child", "Child command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	parentCmd.AddCommand(childCmd.Command)
	
	// Execute child command
	parentCmd.SetArgs([]string{"child"})
	err := parentCmd.Execute()
	
	if err != nil {
		t.Errorf("Subcommand execution failed: %v", err)
	}
}

// Benchmark tests
func BenchmarkNewCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
			return nil
		})
		_ = cmd
	}
}

func BenchmarkCommandWithLogger(b *testing.B) {
	logger, _ := zap.NewDevelopment()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
			return nil
		})
		cmd.SetLogger(logger)
		_ = cmd
	}
}

func BenchmarkCommandExecution(b *testing.B) {
	cmd := New("test", "Test command", func(cmd *cobra.Command, args []string) error {
		return nil
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.SetArgs([]string{})
		_ = cmd.Execute()
	}
}