package command

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Command represents a command with additional functionality
type Command struct {
	*cobra.Command
	logger  *zap.Logger
	context context.Context
}

// New creates a new command
func New(use, short string, runE func(*cobra.Command, []string) error) *Command {
	return &Command{
		Command: &cobra.Command{
			Use:   use,
			Short: short,
			RunE:  runE,
		},
	}
}

// NewWithLong creates a new command with long description
func NewWithLong(use, short, long string, runE func(*cobra.Command, []string) error) *Command {
	return &Command{
		Command: &cobra.Command{
			Use:   use,
			Short: short,
			Long:  long,
			RunE:  runE,
		},
	}
}

// SetLogger sets the logger for the command
func (c *Command) SetLogger(logger *zap.Logger) *Command {
	c.logger = logger
	return c
}

// SetContext sets the context for the command
func (c *Command) SetContext(ctx context.Context) *Command {
	c.context = ctx
	return c
}

// GetLogger returns the logger
func (c *Command) GetLogger() *zap.Logger {
	return c.logger
}

// GetContext returns the context
func (c *Command) GetContext() context.Context {
	return c.context
}
