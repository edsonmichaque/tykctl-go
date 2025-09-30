package prompt

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/edsonmichaque/tykctl-go/terminal"
)

// Prompt represents an interactive prompt
type Prompt struct {
	terminal *terminal.Terminal
}

// New creates a new prompt instance
func New() *Prompt {
	return &Prompt{
		terminal: terminal.New(),
	}
}

// AskString asks for a string input
func (p *Prompt) AskString(question string) (string, error) {
	model := &inputModel{
		question: question,
		input:    "",
		done:     false,
	}
	
	program := tea.NewProgram(model)
	result, err := program.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get input: %w", err)
	}
	
	model = result.(*inputModel)
	return strings.TrimSpace(model.input), nil
}

// AskStringWithDefault asks for a string input with a default value
func (p *Prompt) AskStringWithDefault(question, defaultValue string) (string, error) {
	model := &inputModel{
		question: question,
		input:    defaultValue,
		done:     false,
	}
	
	program := tea.NewProgram(model)
	result, err := program.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get input: %w", err)
	}
	
	model = result.(*inputModel)
	answer := strings.TrimSpace(model.input)
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
	answer, err := p.AskStringWithDefault(question, strconv.Itoa(defaultValue))
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
	model := &confirmModel{
		question: question,
		choice:   -1, // -1 = no choice yet, 0 = no, 1 = yes
		done:     false,
	}
	
	program := tea.NewProgram(model)
	result, err := program.Run()
	if err != nil {
		return false, fmt.Errorf("failed to get confirmation: %w", err)
	}
	
	model = result.(*confirmModel)
	return model.choice == 1, nil
}

// AskBoolWithDefault asks for a boolean input with a default value
func (p *Prompt) AskBoolWithDefault(question string, defaultValue bool) (bool, error) {
	model := &confirmModel{
		question:     question,
		choice:       -1, // -1 = no choice yet, 0 = no, 1 = yes
		defaultValue: defaultValue,
		done:         false,
	}
	
	program := tea.NewProgram(model)
	result, err := program.Run()
	if err != nil {
		return false, fmt.Errorf("failed to get confirmation: %w", err)
	}
	
	model = result.(*confirmModel)
	if model.choice == -1 {
		return defaultValue, nil
	}
	return model.choice == 1, nil
}

// AskSelect asks for a selection from a list
func (p *Prompt) AskSelect(question string, options []string) (string, error) {
	if len(options) == 0 {
		return "", NewNoOptionsError(question)
	}
	
	model := &selectModel{
		question: question,
		options:  options,
		choice:   0,
		done:     false,
	}
	
	program := tea.NewProgram(model)
	result, err := program.Run()
	if err != nil {
		return "", NewInputFailedError(question, "", err)
	}
	
	model = result.(*selectModel)
	return model.options[model.choice], nil
}

// AskMultiSelect asks for multiple selections from a list
func (p *Prompt) AskMultiSelect(question string, options []string) ([]string, error) {
	if len(options) == 0 {
		return nil, NewNoOptionsError(question)
	}
	
	model := &multiSelectModel{
		question: question,
		options:  options,
		choices:  make([]bool, len(options)),
		cursor:   0,
		done:     false,
	}
	
	program := tea.NewProgram(model)
	result, err := program.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get multi-select: %w", err)
	}
	
	model = result.(*multiSelectModel)
	var selected []string
	for i, choice := range model.choices {
		if choice {
			selected = append(selected, model.options[i])
		}
	}
	return selected, nil
}

// AskPassword asks for a password input (hidden)
func (p *Prompt) AskPassword(question string) (string, error) {
	model := &passwordModel{
		question: question,
		input:    "",
		done:     false,
	}
	
	program := tea.NewProgram(model)
	result, err := program.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}
	
	model = result.(*passwordModel)
	return model.input, nil
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

// SetReader is kept for compatibility but does nothing (bubbletea handles input internally)
func (p *Prompt) SetReader(reader interface{}) {
	// bubbletea handles input internally, so this is a no-op
}

// GetReader is kept for compatibility but returns nil (bubbletea handles input internally)
func (p *Prompt) GetReader() interface{} {
	// bubbletea handles input internally, so we return nil
	return nil
}

// inputModel handles text input prompts
type inputModel struct {
	question string
	input    string
	done     bool
}

func (m *inputModel) Init() tea.Cmd {
	return nil
}

func (m *inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}
	}
	return m, nil
}

func (m *inputModel) View() string {
	if m.done {
		return ""
	}
	
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return fmt.Sprintf("%s %s", style.Render("?"), m.question+" "+m.input+"_")
}

// confirmModel handles yes/no confirmation prompts
type confirmModel struct {
	question     string
	choice       int // -1 = no choice, 0 = no, 1 = yes
	defaultValue bool
	done         bool
}

func (m *confirmModel) Init() tea.Cmd {
	return nil
}

func (m *confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.choice == -1 {
				m.choice = 0
				if m.defaultValue {
					m.choice = 1
				}
			}
			m.done = true
			return m, tea.Quit
		case "y", "Y":
			m.choice = 1
		case "n", "N":
			m.choice = 0
		}
	}
	return m, nil
}

func (m *confirmModel) View() string {
	if m.done {
		return ""
	}
	
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	yesStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	noStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	
	yes := "Yes"
	no := "No"
	
	if m.choice == 1 {
		yes = yesStyle.Render("Yes")
	} else if m.choice == 0 {
		no = noStyle.Render("No")
	}
	
	return fmt.Sprintf("%s %s (%s/%s)", style.Render("?"), m.question, yes, no)
}

// selectModel handles single selection prompts
type selectModel struct {
	question string
	options  []string
	choice   int
	done     bool
}

func (m *selectModel) Init() tea.Cmd {
	return nil
}

func (m *selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "up", "k":
			if m.choice > 0 {
				m.choice--
			}
		case "down", "j":
			if m.choice < len(m.options)-1 {
				m.choice++
			}
		}
	}
	return m, nil
}

func (m *selectModel) View() string {
	if m.done {
		return ""
	}
	
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	
	var lines []string
	lines = append(lines, fmt.Sprintf("%s %s", style.Render("?"), m.question))
	
	for i, option := range m.options {
		if i == m.choice {
			lines = append(lines, fmt.Sprintf("  %s %s", selectedStyle.Render(">"), option))
		} else {
			lines = append(lines, fmt.Sprintf("    %s", option))
		}
	}
	
	return strings.Join(lines, "\n")
}

// multiSelectModel handles multiple selection prompts
type multiSelectModel struct {
	question string
	options  []string
	choices  []bool
	cursor   int
	done     bool
}

func (m *multiSelectModel) Init() tea.Cmd {
	return nil
}

func (m *multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case " ":
			m.choices[m.cursor] = !m.choices[m.cursor]
		}
	}
	return m, nil
}

func (m *multiSelectModel) View() string {
	if m.done {
		return ""
	}
	
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	checkedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	
	var lines []string
	lines = append(lines, fmt.Sprintf("%s %s", style.Render("?"), m.question))
	
	for i, option := range m.options {
		var prefix string
		if i == m.cursor {
			prefix = selectedStyle.Render(">")
		} else {
			prefix = " "
		}
		
		var checkbox string
		if m.choices[i] {
			checkbox = checkedStyle.Render("✓")
		} else {
			checkbox = " "
		}
		
		lines = append(lines, fmt.Sprintf("  %s [%s] %s", prefix, checkbox, option))
	}
	
	return strings.Join(lines, "\n")
}

// passwordModel handles password input prompts
type passwordModel struct {
	question string
	input    string
	done     bool
}

func (m *passwordModel) Init() tea.Cmd {
	return nil
}

func (m *passwordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}
	}
	return m, nil
}

func (m *passwordModel) View() string {
	if m.done {
		return ""
	}
	
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	hidden := strings.Repeat("*", len(m.input))
	return fmt.Sprintf("%s %s %s_", style.Render("?"), m.question, hidden)
}