// Package spinner provides spinner functionality for tykctl.
package spinner

import (
	"context"
	"time"

	"github.com/briandowns/spinner"
)

// Spinner manages spinner operations
type Spinner struct {
	spinner *spinner.Spinner
}

// New creates a new spinner
func New() *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	return &Spinner{spinner: s}
}

// NewWithCharSet creates a new spinner with specific character set
func NewWithCharSet(charSet int, duration time.Duration) *Spinner {
	s := spinner.New(spinner.CharSets[charSet], duration)
	return &Spinner{spinner: s}
}

// Start starts the spinner with a message
func (s *Spinner) Start(message string) {
	s.spinner.Suffix = " " + message
	s.spinner.Start()
}

// Update updates the spinner message
func (s *Spinner) Update(message string) {
	s.spinner.Suffix = " " + message
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.spinner.Stop()
}

// StopWithMessage stops the spinner and prints a final message
func (s *Spinner) StopWithMessage(message string) {
	s.spinner.FinalMSG = message + "\n"
	s.spinner.Stop()
}

// WithContext runs a function with a spinner that respects context cancellation
func (s *Spinner) WithContext(ctx context.Context, message string, fn func() error) error {
	s.Start(message)
	defer s.Stop()

	// Create a channel to receive the function result
	resultChan := make(chan error, 1)

	// Run the function in a goroutine
	go func() {
		resultChan <- fn()
	}()

	// Wait for either context cancellation or function completion
	select {
	case <-ctx.Done():
		s.StopWithMessage("Operation cancelled")
		return ctx.Err()
	case err := <-resultChan:
		if err != nil {
			s.StopWithMessage("Operation failed")
			return err
		}
		s.StopWithMessage("Operation completed")
		return nil
	}
}

// WithTimeout runs a function with a spinner and timeout
func (s *Spinner) WithTimeout(timeout time.Duration, message string, fn func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.WithContext(ctx, message, fn)
}
