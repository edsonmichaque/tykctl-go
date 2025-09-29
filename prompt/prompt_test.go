package prompt

import (
	"bufio"
	"strings"
	"testing"

	"github.com/edsonmichaque/tykctl-go/terminal"
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
		reader:   bufio.NewReader(strings.NewReader("test input\n")),
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
		reader:   bufio.NewReader(strings.NewReader("\n")), // Empty input
	}

	result, err := prompt.AskStringWithDefault("Enter name", "default")
	if err != nil {
		t.Errorf("AskStringWithDefault failed: %v", err)
	}

	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}

	// Test with actual input
	prompt.reader = bufio.NewReader(strings.NewReader("custom input\n"))
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
		reader:   bufio.NewReader(strings.NewReader("42\n")),
	}

	result, err := prompt.AskInt("Enter a number: ")
	if err != nil {
		t.Errorf("AskInt failed: %v", err)
	}

	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test invalid integer
	prompt.reader = bufio.NewReader(strings.NewReader("not a number\n"))
	_, err = prompt.AskInt("Enter a number: ")
	if err == nil {
		t.Error("Expected error for invalid integer, got nil")
	}
}

func TestAskIntWithDefault(t *testing.T) {
	// Test with default value
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   bufio.NewReader(strings.NewReader("\n")), // Empty input
	}

	result, err := prompt.AskIntWithDefault("Enter a number", 10)
	if err != nil {
		t.Errorf("AskIntWithDefault failed: %v", err)
	}

	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}

	// Test with actual input
	prompt.reader = bufio.NewReader(strings.NewReader("25\n"))
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
		reader:   bufio.NewReader(strings.NewReader("yes\n")),
	}

	result, err := prompt.AskBool("Continue? (y/n): ")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
	}

	if !result {
		t.Error("Expected true, got false")
	}

	// Test no
	prompt.reader = bufio.NewReader(strings.NewReader("no\n"))
	result, err = prompt.AskBool("Continue? (y/n): ")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
	}

	if result {
		t.Error("Expected false, got true")
	}

	// Test y
	prompt.reader = bufio.NewReader(strings.NewReader("y\n"))
	result, err = prompt.AskBool("Continue? (y/n): ")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
	}

	if !result {
		t.Error("Expected true, got false")
	}

	// Test n
	prompt.reader = bufio.NewReader(strings.NewReader("n\n"))
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
		reader:   bufio.NewReader(strings.NewReader("\n")), // Empty input
	}

	result, err := prompt.AskBoolWithDefault("Continue? (y/n)", true)
	if err != nil {
		t.Errorf("AskBoolWithDefault failed: %v", err)
	}

	if !result {
		t.Error("Expected true, got false")
	}

	// Test with default false
	prompt2 := &Prompt{
		terminal: terminal.New(),
		reader:   bufio.NewReader(strings.NewReader("\n")), // Empty input
	}
	result, err = prompt2.AskBoolWithDefault("Continue? (y/n)", false)
	if err != nil {
		t.Errorf("AskBoolWithDefault failed: %v", err)
	}

	if result {
		t.Error("Expected false, got true")
	}
}

func TestAskSelect(t *testing.T) {
	choices := []string{"option1", "option2", "option3"}

	// Test valid choice
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   bufio.NewReader(strings.NewReader("2\n")),
	}

	result, err := prompt.AskSelect("Choose an option", choices)
	if err != nil {
		t.Errorf("AskSelect failed: %v", err)
	}

	if result != "option2" {
		t.Errorf("Expected 'option2', got '%s'", result)
	}

	// Test invalid choice
	prompt.reader = bufio.NewReader(strings.NewReader("invalid\n"))
	_, err = prompt.AskSelect("Choose an option", choices)
	if err == nil {
		t.Error("Expected error for invalid choice, got nil")
	}
}

func TestAskSelectWithDefault(t *testing.T) {
	choices := []string{"option1", "option2", "option3"}

	// Test with default
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   bufio.NewReader(strings.NewReader("1\n")), // Choose first option
	}

	result, err := prompt.AskSelect("Choose an option", choices)
	if err != nil {
		t.Errorf("AskSelect failed: %v", err)
	}

	if result != "option1" {
		t.Errorf("Expected 'option1', got '%s'", result)
	}

	// Test with actual input
	prompt.reader = bufio.NewReader(strings.NewReader("3\n"))
	result, err = prompt.AskSelect("Choose an option", choices)
	if err != nil {
		t.Errorf("AskSelect failed: %v", err)
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
	prompt.reader = bufio.NewReader(strings.NewReader("secretpassword\n"))

	result, err := prompt.AskPassword("Enter password: ")
	if err != nil {
		t.Errorf("AskPassword failed: %v", err)
	}

	if result != "secretpassword" {
		t.Errorf("Expected 'secretpassword', got '%s'", result)
	}
}

func TestAskBoolConfirmation(t *testing.T) {
	// Test confirmation
	prompt := &Prompt{
		terminal: terminal.New(),
		reader:   bufio.NewReader(strings.NewReader("yes\n")),
	}

	result, err := prompt.AskBool("Are you sure?")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
	}

	if !result {
		t.Error("Expected true, got false")
	}

	// Test rejection
	prompt.reader = bufio.NewReader(strings.NewReader("no\n"))
	result, err = prompt.AskBool("Are you sure?")
	if err != nil {
		t.Errorf("AskBool failed: %v", err)
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
		reader:   bufio.NewReader(strings.NewReader("test input\n")),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prompt.reader = bufio.NewReader(strings.NewReader("test input\n"))
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
		reader:   bufio.NewReader(strings.NewReader("42\n")),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prompt.reader = bufio.NewReader(strings.NewReader("42\n"))
		result, err := prompt.AskInt("Enter a number: ")
		if err != nil {
			b.Fatalf("AskInt failed: %v", err)
		}
		_ = result
	}
}
