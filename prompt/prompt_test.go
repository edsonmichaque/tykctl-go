package prompt

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	prompt := New()
	if prompt == nil {
		t.Error("New() returned nil")
	}
	
	if prompt.terminal == nil {
		t.Error("Terminal should not be nil")
	}
	
	if prompt.reader == nil {
		t.Error("Reader should not be nil")
	}
}

func TestAskString(t *testing.T) {
	// Create a prompt with a custom reader
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("test input\n"),
	}
	
	result, err := prompt.AskString("Enter something: ")
	if err != nil {
		t.Errorf("AskString failed: %v", err)
	}
	
	if result != "test input" {
		t.Errorf("Expected 'test input', got '%s'", result)
	}
}

func TestAskStringWithDefault(t *testing.T) {
	// Test with default value
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("\n"), // Empty input
	}
	
	result, err := prompt.AskStringWithDefault("Enter name", "default")
	if err != nil {
		t.Errorf("AskStringWithDefault failed: %v", err)
	}
	
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
	
	// Test with actual input
	prompt.reader = strings.NewReader("custom input\n")
	result, err = prompt.AskStringWithDefault("Enter name", "default")
	if err != nil {
		t.Errorf("AskStringWithDefault failed: %v", err)
	}
	
	if result != "custom input" {
		t.Errorf("Expected 'custom input', got '%s'", result)
	}
}

func TestAskInt(t *testing.T) {
	// Test valid integer
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("42\n"),
	}
	
	result, err := prompt.AskInt("Enter a number: ")
	if err != nil {
		t.Errorf("AskInt failed: %v", err)
	}
	
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
	
	// Test invalid integer
	prompt.reader = strings.NewReader("not a number\n")
	_, err = prompt.AskInt("Enter a number: ")
	if err == nil {
		t.Error("Expected error for invalid integer, got nil")
	}
}

func TestAskIntWithDefault(t *testing.T) {
	// Test with default value
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("\n"), // Empty input
	}
	
	result, err := prompt.AskIntWithDefault("Enter a number", 10)
	if err != nil {
		t.Errorf("AskIntWithDefault failed: %v", err)
	}
	
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}
	
	// Test with actual input
	prompt.reader = strings.NewReader("25\n")
	result, err = prompt.AskIntWithDefault("Enter a number", 10)
	if err != nil {
		t.Errorf("AskIntWithDefault failed: %v", err)
	}
	
	if result != 25 {
		t.Errorf("Expected 25, got %d", result)
	}
}

func TestAskBool(t *testing.T) {
	// Test yes
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("yes\n"),
	}
	
	result, err := prompt.AskBool("Continue? (y/n): ")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
	}
	
	if !result {
		t.Error("Expected true, got false")
	}
	
	// Test no
	prompt.reader = strings.NewReader("no\n")
	result, err = prompt.AskBool("Continue? (y/n): ")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
	}
	
	if result {
		t.Error("Expected false, got true")
	}
	
	// Test y
	prompt.reader = strings.NewReader("y\n")
	result, err = prompt.AskBool("Continue? (y/n): ")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
	}
	
	if !result {
		t.Error("Expected true, got false")
	}
	
	// Test n
	prompt.reader = strings.NewReader("n\n")
	result, err = prompt.AskBool("Continue? (y/n): ")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
	}
	
	if result {
		t.Error("Expected false, got true")
	}
}

func TestAskBoolWithDefault(t *testing.T) {
	// Test with default true
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("\n"), // Empty input
	}
	
	result, err := prompt.AskBoolWithDefault("Continue? (y/n)", true)
	if err != nil {
		t.Errorf("AskBoolWithDefault failed: %v", err)
	}
	
	if !result {
		t.Error("Expected true, got false")
	}
	
	// Test with default false
	result, err = prompt.AskBoolWithDefault("Continue? (y/n)", false)
	if err != nil {
		t.Errorf("AskBoolWithDefault failed: %v", err)
	}
	
	if result {
		t.Error("Expected false, got true")
	}
}

func TestAskChoice(t *testing.T) {
	choices := []string{"option1", "option2", "option3"}
	
	// Test valid choice
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("option2\n"),
	}
	
	result, err := prompt.AskChoice("Choose an option", choices)
	if err != nil {
		t.Errorf("AskChoice failed: %v", err)
	}
	
	if result != "option2" {
		t.Errorf("Expected 'option2', got '%s'", result)
	}
	
	// Test invalid choice
	prompt.reader = strings.NewReader("invalid\n")
	_, err = prompt.AskChoice("Choose an option", choices)
	if err == nil {
		t.Error("Expected error for invalid choice, got nil")
	}
}

func TestAskChoiceWithDefault(t *testing.T) {
	choices := []string{"option1", "option2", "option3"}
	
	// Test with default
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("\n"), // Empty input
	}
	
	result, err := prompt.AskChoiceWithDefault("Choose an option", choices, "option2")
	if err != nil {
		t.Errorf("AskChoiceWithDefault failed: %v", err)
	}
	
	if result != "option2" {
		t.Errorf("Expected 'option2', got '%s'", result)
	}
	
	// Test with actual input
	prompt.reader = strings.NewReader("option3\n")
	result, err = prompt.AskChoiceWithDefault("Choose an option", choices, "option2")
	if err != nil {
		t.Errorf("AskChoiceWithDefault failed: %v", err)
	}
	
	if result != "option3" {
		t.Errorf("Expected 'option3', got '%s'", result)
	}
}

func TestAskPassword(t *testing.T) {
	// This is hard to test without mocking os.Stdin
	// We'll just test that it doesn't panic
	prompt := New()
	
	// Create a mock reader that simulates password input
	prompt.reader = strings.NewReader("secretpassword\n")
	
	result, err := prompt.AskPassword("Enter password: ")
	if err != nil {
		t.Errorf("AskPassword failed: %v", err)
	}
	
	if result != "secretpassword" {
		t.Errorf("Expected 'secretpassword', got '%s'", result)
	}
}

func TestConfirm(t *testing.T) {
	// Test confirmation
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("yes\n"),
	}
	
	result, err := prompt.Confirm("Are you sure?")
	if err != nil {
		t.Errorf("Confirm failed: %v", err)
	}
	
	if !result {
		t.Error("Expected true, got false")
	}
	
	// Test rejection
	prompt.reader = strings.NewReader("no\n")
	result, err = prompt.Confirm("Are you sure?")
	if err != nil {
		t.Errorf("Confirm failed: %v", err)
	}
	
	if result {
		t.Error("Expected false, got true")
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prompt := New()
		_ = prompt
	}
}

func BenchmarkAskString(b *testing.B) {
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("test input\n"),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prompt.reader = strings.NewReader("test input\n")
		result, err := prompt.AskString("Enter something: ")
		if err != nil {
			b.Fatalf("AskString failed: %v", err)
		}
		_ = result
	}
}

func BenchmarkAskInt(b *testing.B) {
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   strings.NewReader("42\n"),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prompt.reader = strings.NewReader("42\n")
		result, err := prompt.AskInt("Enter a number: ")
		if err != nil {
			b.Fatalf("AskInt failed: %v", err)
		}
		_ = result
	}
}