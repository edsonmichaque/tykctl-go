package spinner

import (
	"context"

	"github.com/edsonmichaque/tykctl-go/progress"
)

// Spinner wraps the progress.Spinner to provide a compatible interface
type Spinner struct {
	*progress.Spinner
}

// New creates a new spinner instance using the progress package
func New() *Spinner {
	return &Spinner{
		Spinner: progress.New(),
	}
}

// WithContext runs a function with a spinner showing the given message
func (s *Spinner) WithContext(ctx context.Context, message string, fn func() error) error {
	return s.Spinner.WithContext(ctx, message, fn)
}