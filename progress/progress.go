package progress

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner represents a spinner
type Spinner struct {
	message string
	frames  []string
	index   int
	done    bool
	mu      sync.Mutex
}

// Bar represents a progress bar
type Bar struct {
	message   string
	total     int64
	current   int64
	width     int
	done      bool
	fillChar  string
	emptyChar string
	mu        sync.Mutex
}

// New creates a new spinner
func New() *Spinner {
	return &Spinner{
		message: "Loading...",
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// NewBar creates a new progress bar
func NewBar(total int64) *Bar {
	return &Bar{
		message:   "Progress",
		total:     total,
		width:     50,
		fillChar:  "█",
		emptyChar: "░",
	}
}

// WithMessage sets the spinner message
func (s *Spinner) WithMessage(message string) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
	return s
}

// WithContext runs a function with a spinner
func (s *Spinner) WithContext(ctx context.Context, message string, fn func() error) error {
	s.mu.Lock()
	s.message = message
	s.done = false
	s.mu.Unlock()

	// Create internal spinner model
	model := &spinnerModel{
		spinner: s,
		done:    false,
	}

	// Create a new program
	p := tea.NewProgram(model, tea.WithOutput(os.Stderr))

	// Run the function in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- fn()
		s.mu.Lock()
		s.done = true
		model.done = true
		s.mu.Unlock()
	}()

	// Start the spinner
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running spinner: %v\n", err)
		}
	}()

	// Wait for completion or context cancellation
	select {
	case err := <-errChan:
		s.mu.Lock()
		s.done = true
		model.done = true
		s.mu.Unlock()
		return err
	case <-ctx.Done():
		s.mu.Lock()
		s.done = true
		model.done = true
		s.mu.Unlock()
		return ctx.Err()
	}
}

// WithMessage sets the progress bar message
func (b *Bar) WithMessage(message string) *Bar {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.message = message
	return b
}

// Add adds to the current progress
func (b *Bar) Add(inc int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current += inc
	if b.current > b.total {
		b.current = b.total
	}
}

// SetCurrent sets the current progress
func (b *Bar) SetCurrent(current int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = current
	if b.current > b.total {
		b.current = b.total
	}
}

// WithContext runs a function with a progress bar
func (b *Bar) WithContext(ctx context.Context, message string, total int64, fn func(update func(int64)) error) error {
	b.mu.Lock()
	b.message = message
	b.total = total
	b.current = 0
	b.done = false
	b.mu.Unlock()

	// Create internal progress bar model
	model := &barModel{
		bar:  b,
		done: false,
	}

	// Create a new program
	p := tea.NewProgram(model, tea.WithOutput(os.Stderr))

	// Run the function in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- fn(func(inc int64) {
			b.Add(inc)
		})
		b.mu.Lock()
		b.done = true
		model.done = true
		b.mu.Unlock()
	}()

	// Start the progress bar
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running progress bar: %v\n", err)
		}
	}()

	// Wait for completion or context cancellation
	select {
	case err := <-errChan:
		b.mu.Lock()
		b.done = true
		model.done = true
		b.mu.Unlock()
		return err
	case <-ctx.Done():
		b.mu.Lock()
		b.done = true
		model.done = true
		b.mu.Unlock()
		return ctx.Err()
	}
}

// Internal models that implement tea.Model interface
// These are not exposed to the public API

type spinnerModel struct {
	spinner *Spinner
	done    bool
}

func (m *spinnerModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.done {
			return m, tea.Quit
		}
		m.spinner.mu.Lock()
		m.spinner.index = (m.spinner.index + 1) % len(m.spinner.frames)
		m.spinner.mu.Unlock()
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *spinnerModel) View() string {
	if m.done {
		return ""
	}
	m.spinner.mu.Lock()
	defer m.spinner.mu.Unlock()
	return fmt.Sprintf("%s %s", m.spinner.frames[m.spinner.index], m.spinner.message)
}

type barModel struct {
	bar  *Bar
	done bool
}

func (m *barModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *barModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.done {
			return m, tea.Quit
		}
		return m, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *barModel) View() string {
	if m.done {
		return ""
	}
	m.bar.mu.Lock()
	defer m.bar.mu.Unlock()

	if m.bar.total <= 0 {
		return fmt.Sprintf("%s: %d", m.bar.message, m.bar.current)
	}

	percent := float64(m.bar.current) / float64(m.bar.total)
	filled := int(float64(m.bar.width) * percent)
	empty := m.bar.width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += m.bar.fillChar
	}
	for i := 0; i < empty; i++ {
		bar += m.bar.emptyChar
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	return fmt.Sprintf("%s: %s %d/%d (%.1f%%)",
		m.bar.message,
		style.Render(bar),
		m.bar.current,
		m.bar.total,
		percent*100)
}

// tickMsg is a message sent on each tick
type tickMsg struct{}
