// Package logger provides logging functionality for tykctl.
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional functionality
type Logger struct {
	*zap.Logger
}

// Config represents logger configuration
type Config struct {
	Debug   bool
	Verbose bool
	NoColor bool
}

// New creates a new logger with the given configuration
func New(config Config) *Logger {
	var zapConfig zap.Config

	if config.Debug || config.Verbose {
		// Development config for debug/verbose mode
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		// Production config for normal mode
		zapConfig = zap.NewProductionConfig()
	}

	// Set log level
	if config.Debug {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else if config.Verbose {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	}

	// Disable colors if requested
	if config.NoColor || os.Getenv("NO_COLOR") != "" {
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Set output to stderr
	zapConfig.OutputPaths = []string{"stderr"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}

	// Build the logger
	zapLogger, err := zapConfig.Build()
	if err != nil {
		// Fallback to a basic logger if config fails
		zapLogger, _ = zap.NewProduction()
	}

	return &Logger{Logger: zapLogger}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() {
	if l.Logger != nil {
		l.Logger.Sync()
	}
}

// Global logger instance for backward compatibility
var global *Logger

// InitGlobal initializes the global logger
func InitGlobal(config Config) {
	global = New(config)
}

// GetGlobal returns the global logger instance
func GetGlobal() *Logger {
	if global == nil {
		// Initialize with default config if not set
		global = New(Config{})
	}
	return global
}

// SyncGlobal flushes any buffered log entries from the global logger
func SyncGlobal() {
	if global != nil {
		global.Sync()
	}
}
