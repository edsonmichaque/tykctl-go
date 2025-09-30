package progress

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

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

	// Start spinner in background
	stop := make(chan struct{})
	go s.run(stop)

	// Run the function
	errChan := make(chan error, 1)
	go func() {
		errChan <- fn()
	}()

	// Wait for completion or context cancellation
	select {
	case err := <-errChan:
		s.mu.Lock()
		s.done = true
		s.mu.Unlock()
		close(stop)
		s.clearLine()
		return err
	case <-ctx.Done():
		s.mu.Lock()
		s.done = true
		s.mu.Unlock()
		close(stop)
		s.clearLine()
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

	// Start progress bar in background
	stop := make(chan struct{})
	go b.run(stop)

	// Run the function
	errChan := make(chan error, 1)
	go func() {
		errChan <- fn(func(inc int64) {
			b.Add(inc)
		})
	}()

	// Wait for completion or context cancellation
	select {
	case err := <-errChan:
		b.mu.Lock()
		b.done = true
		b.mu.Unlock()
		close(stop)
		b.clearLine()
		return err
	case <-ctx.Done():
		b.mu.Lock()
		b.done = true
		b.mu.Unlock()
		close(stop)
		b.clearLine()
		return ctx.Err()
	}
}

// run runs the spinner animation
func (s *Spinner) run(stop <-chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			s.mu.Lock()
			if s.done {
				s.mu.Unlock()
				return
			}
			s.index = (s.index + 1) % len(s.frames)
			s.mu.Unlock()
			s.render()
		}
	}
}

// run runs the progress bar animation
func (b *Bar) run(stop <-chan struct{}) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			b.mu.Lock()
			if b.done {
				b.mu.Unlock()
				return
			}
			b.mu.Unlock()
			b.render()
		}
	}
}

// render renders the spinner
func (s *Spinner) render() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Clear the line and print spinner
	fmt.Fprintf(os.Stderr, "\r%s %s", s.frames[s.index], s.message)
}

// render renders the progress bar
func (b *Bar) render() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.total <= 0 {
		fmt.Fprintf(os.Stderr, "\r%s: %d", b.message, b.current)
		return
	}

	percent := float64(b.current) / float64(b.total)
	filled := int(float64(b.width) * percent)
	empty := b.width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += b.fillChar
	}
	for i := 0; i < empty; i++ {
		bar += b.emptyChar
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	fmt.Fprintf(os.Stderr, "\r%s: %s %d/%d (%.1f%%)",
		b.message,
		style.Render(bar),
		b.current,
		b.total,
		percent*100)
}

// clearLine clears the current line
func (s *Spinner) clearLine() {
	fmt.Fprintf(os.Stderr, "\r%s", "\033[2K")
}

// clearLine clears the current line
func (b *Bar) clearLine() {
	fmt.Fprintf(os.Stderr, "\r%s", "\033[2K")
}