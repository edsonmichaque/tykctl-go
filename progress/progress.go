package progress

import (
	"context"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// Spinner represents a spinner
type Spinner struct {
	spinner *spinner.Spinner
	mu      sync.Mutex
}

// Bar represents a progress bar
type Bar struct {
	bar     *mpb.Bar
	message string
	total   int64
	current int64
	mu      sync.Mutex
}

// New creates a new spinner
func New() *Spinner {
	// Custom spinner frames: ⠋⠙⠹⠸⠼
	customFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼"}
	s := spinner.New(customFrames, 100*time.Millisecond)
	s.Suffix = " "
	s.FinalMSG = "✓ Complete!\n"
	return &Spinner{
		spinner: s,
	}
}

// NewBar creates a new progress bar
func NewBar(total int64) *Bar {
	return &Bar{
		total: total,
	}
}

// WithMessage sets the spinner message
func (s *Spinner) WithMessage(message string) *Spinner {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.spinner != nil {
		s.spinner.Suffix = " " + message
	}
	return s
}

// WithContext runs a function with a spinner
func (s *Spinner) WithContext(ctx context.Context, message string, fn func() error) error {
	s.mu.Lock()
	if s.spinner != nil {
		s.spinner.Suffix = " " + message
		s.spinner.Start()
	}
	s.mu.Unlock()

	// Run the function in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- fn()
	}()

	// Wait for completion or context cancellation
	select {
	case err := <-errChan:
		s.mu.Lock()
		if s.spinner != nil {
			s.spinner.Stop()
		}
		s.mu.Unlock()
		return err
	case <-ctx.Done():
		s.mu.Lock()
		if s.spinner != nil {
			s.spinner.Stop()
		}
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
	if b.bar != nil {
		b.bar.IncrBy(int(inc))
	}
	b.current += inc
	if b.current > b.total {
		b.current = b.total
	}
}

// SetCurrent sets the current progress
func (b *Bar) SetCurrent(current int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.bar != nil {
		b.bar.SetCurrent(current)
	}
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
	b.mu.Unlock()

	// Create progress container
	p := mpb.New(mpb.WithWidth(64), mpb.WithRefreshRate(50*time.Millisecond))

	// Create progress bar
	bar := p.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(message, decor.WC{W: len(message) + 1, C: decor.DindentRight}),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WC{W: 5}),
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WC{W: 4}), "done",
			),
		),
	)

	b.mu.Lock()
	b.bar = bar
	b.mu.Unlock()

	// Run the function in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- fn(func(inc int64) {
			b.Add(inc)
		})
	}()

	// Wait for completion or context cancellation
	select {
	case err := <-errChan:
		p.Wait()
		return err
	case <-ctx.Done():
		p.Wait()
		return ctx.Err()
	}
}
