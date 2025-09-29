package prompt

import (
	"errors"
	"fmt"
	"log"
)

// ExampleErrorHandling demonstrates how to use the new error types
func ExampleErrorHandling() {
	p := New()

	// Example 1: Handle choice errors
	options := []string{"option1", "option2", "option3"}
	choice, err := p.AskSelect("Choose an option:", options)
	if err != nil {
		if IsChoiceError(err) {
			var choiceErr *ChoiceError
			if errors.As(err, &choiceErr) {
				fmt.Printf("Choice error: choice %d is out of range (max: %d)\n",
					choiceErr.Choice, choiceErr.MaxChoice)
			}
		} else if IsPromptError(err) {
			var promptErr *PromptError
			if errors.As(err, &promptErr) {
				fmt.Printf("Prompt error: %s\n", promptErr.Message)
			}
		}
		log.Fatal(err)
	}
	fmt.Printf("You chose: %s\n", choice)

	// Example 2: Handle input errors
	number, err := p.AskInt("Enter a number:")
	if err != nil {
		if IsInputError(err) {
			var inputErr *InputError
			if errors.As(err, &inputErr) {
				fmt.Printf("Input error: '%s' is not a valid %s\n",
					inputErr.Input, inputErr.Expected)
			}
		}
		log.Fatal(err)
	}
	fmt.Printf("You entered: %d\n", number)

	// Example 3: Handle validation errors
	text, err := p.AskString("Enter some text:")
	if err != nil {
		if IsValidationError(err) {
			var validationErr *ValidationError
			if errors.As(err, &validationErr) {
				fmt.Printf("Validation error: %s for field '%s'\n",
					validationErr.Message, validationErr.Field)
			}
		}
		log.Fatal(err)
	}
	fmt.Printf("You entered: %s\n", text)
}

// ExampleErrorTypes demonstrates the different error types
func ExampleErrorTypes() {
	// Standard errors
	fmt.Println("Standard errors:")
	fmt.Printf("ErrInputFailed: %v\n", ErrInputFailed)
	fmt.Printf("ErrInvalidChoice: %v\n", ErrInvalidChoice)
	fmt.Printf("ErrNoOptions: %v\n", ErrNoOptions)

	// Structured errors
	fmt.Println("\nStructured errors:")

	// PromptError
	promptErr := NewPromptError("test_type", "test message", "test question", "test input", nil)
	fmt.Printf("PromptError: %v\n", promptErr)

	// ValidationError
	validationErr := NewValidationError("field", "value", "rule", "message", "question", nil)
	fmt.Printf("ValidationError: %v\n", validationErr)

	// ChoiceError
	choiceErr := NewChoiceError(5, 3, []string{"a", "b", "c"}, "out of range", "choose", nil)
	fmt.Printf("ChoiceError: %v\n", choiceErr)

	// InputError
	inputErr := NewInputError("invalid", "valid", "invalid format", "enter", nil)
	fmt.Printf("InputError: %v\n", inputErr)
}

// ExampleErrorHelpers demonstrates error helper functions
func ExampleErrorHelpers() {
	// Create some errors
	promptErr := NewPromptError("test", "message", "question", "input", nil)
	validationErr := NewValidationError("field", "value", "rule", "message", "question", nil)
	choiceErr := NewChoiceError(1, 3, []string{"a", "b", "c"}, "message", "question", nil)
	inputErr := NewInputError("input", "expected", "message", "question", nil)

	// Test error type checking
	fmt.Printf("IsPromptError(promptErr): %t\n", IsPromptError(promptErr))
	fmt.Printf("IsValidationError(validationErr): %t\n", IsValidationError(validationErr))
	fmt.Printf("IsChoiceError(choiceErr): %t\n", IsChoiceError(choiceErr))
	fmt.Printf("IsInputError(inputErr): %t\n", IsInputError(inputErr))

	// Test error wrapping
	wrappedErr := WrapError(promptErr, "additional context")
	fmt.Printf("Wrapped error: %v\n", wrappedErr)

	wrappedErrF := WrapErrorf(promptErr, "additional context with %s", "formatting")
	fmt.Printf("Wrapped error with formatting: %v\n", wrappedErrF)
}
