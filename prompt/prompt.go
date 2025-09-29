package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/edsonmichaque/tykctl-go/terminal"
)

// Prompt represents an interactive prompt
type Prompt struct {
	terminal *terminal.Terminal
	reader   *bufio.Reader
}

// New creates a new prompt instance
func New() *Prompt {
	return &Prompt{
		terminal: terminal.New(),
		reader:   bufio.NewReader(os.Stdin),
	}
}

// AskString asks for a string input
func (p *Prompt) AskString(question string) (string, error) {
	fmt.Print(question + " ")
	text, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// AskStringWithDefault asks for a string input with a default value
func (p *Prompt) AskStringWithDefault(question, defaultValue string) (string, error) {
	if defaultValue != "" {
		question = fmt.Sprintf("%s [%s]", question, defaultValue)
	}

	answer, err := p.AskString(question)
	if err != nil {
		return "", err
	}

	if answer == "" {
		return defaultValue, nil
	}

	return answer, nil
}

// AskInt asks for an integer input
func (p *Prompt) AskInt(question string) (int, error) {
	answer, err := p.AskString(question)
	if err != nil {
		return 0, NewInputFailedError(question, answer, err)
	}

	value, err := strconv.Atoi(answer)
	if err != nil {
		return 0, NewInvalidNumberError(answer, question, err)
	}

	return value, nil
}

// AskIntWithDefault asks for an integer input with a default value
func (p *Prompt) AskIntWithDefault(question string, defaultValue int) (int, error) {
	question = fmt.Sprintf("%s [%d]", question, defaultValue)
	answer, err := p.AskString(question)
	if err != nil {
		return 0, NewInputFailedError(question, answer, err)
	}

	if answer == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(answer)
	if err != nil {
		return 0, NewInvalidNumberError(answer, question, err)
	}

	return value, nil
}

// AskBool asks for a boolean input (y/n)
func (p *Prompt) AskBool(question string) (bool, error) {
	answer, err := p.AskString(question + " (y/n)")
	if err != nil {
		return false, err
	}

	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

// AskBoolWithDefault asks for a boolean input with a default value
func (p *Prompt) AskBoolWithDefault(question string, defaultValue bool) (bool, error) {
	defaultStr := "n"
	if defaultValue {
		defaultStr = "y"
	}

	question = fmt.Sprintf("%s (y/n) [%s]", question, defaultStr)
	answer, err := p.AskString(question)
	if err != nil {
		return false, err
	}

	if answer == "" {
		return defaultValue, nil
	}

	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

// AskSelect asks for a selection from a list
func (p *Prompt) AskSelect(question string, options []string) (string, error) {
	if len(options) == 0 {
		return "", NewNoOptionsError(question)
	}

	fmt.Println(question)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}

	answer, err := p.AskString("Enter your choice (number):")
	if err != nil {
		return "", NewInputFailedError(question, answer, err)
	}

	choice, err := strconv.Atoi(answer)
	if err != nil {
		return "", NewInvalidChoiceError(0, len(options), options, question)
	}

	if choice < 1 || choice > len(options) {
		return "", NewChoiceOutOfRangeError(choice, len(options), options, question)
	}

	return options[choice-1], nil
}

// AskMultiSelect asks for multiple selections from a list
func (p *Prompt) AskMultiSelect(question string, options []string) ([]string, error) {
	if len(options) == 0 {
		return nil, NewNoOptionsError(question)
	}

	fmt.Println(question)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}

	answer, err := p.AskString("Enter your choices (comma-separated numbers):")
	if err != nil {
		return nil, err
	}

	choices := strings.Split(answer, ",")
	var selected []string

	for _, choiceStr := range choices {
		choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
		if err != nil {
			return nil, fmt.Errorf("invalid choice: %s", choiceStr)
		}

		if choice < 1 || choice > len(options) {
			return nil, fmt.Errorf("choice out of range: %d", choice)
		}

		selected = append(selected, options[choice-1])
	}

	return selected, nil
}

// AskPassword asks for a password input (hidden)
func (p *Prompt) AskPassword(question string) (string, error) {
	fmt.Print(question + " ")

	// For now, just read normally - in a real implementation,
	// you'd want to hide the input
	text, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

// AskConfirmation asks for confirmation
func (p *Prompt) AskConfirmation(question string) (bool, error) {
	return p.AskBool(question)
}

// AskConfirmationWithDefault asks for confirmation with a default value
func (p *Prompt) AskConfirmationWithDefault(question string, defaultValue bool) (bool, error) {
	return p.AskBoolWithDefault(question, defaultValue)
}

// AskYesNo asks for yes/no input
func (p *Prompt) AskYesNo(question string) (bool, error) {
	return p.AskBool(question)
}

// AskYesNoWithDefault asks for yes/no input with a default value
func (p *Prompt) AskYesNoWithDefault(question string, defaultValue bool) (bool, error) {
	return p.AskBoolWithDefault(question, defaultValue)
}

// SetTerminal sets the terminal instance
func (p *Prompt) SetTerminal(term *terminal.Terminal) {
	p.terminal = term
}

// GetTerminal returns the terminal instance
func (p *Prompt) GetTerminal() *terminal.Terminal {
	return p.terminal
}

// IsInteractive returns whether the prompt is interactive
func (p *Prompt) IsInteractive() bool {
	return p.terminal.IsInteractive()
}

// SetReader sets the reader
func (p *Prompt) SetReader(reader *bufio.Reader) {
	p.reader = reader
}

// GetReader returns the reader
func (p *Prompt) GetReader() *bufio.Reader {
	return p.reader
}
