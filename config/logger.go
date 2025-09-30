package config

import (
	"log"
	"os"
)

// LoggerOptions provides configuration for logger
type LoggerOptions struct {
	Level LogLevel
}

// NewLogger creates a new logger instance
func NewLogger(opts LoggerOptions) (Logger, error) {
	return &simpleLogger{
		level: opts.Level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}, nil
}

// simpleLogger implements a simple logger
type simpleLogger struct {
	level  LogLevel
	logger *log.Logger
}

func (l *simpleLogger) Debug(msg string, fields ...interface{}) {
	if l.level <= LogLevelDebug {
		l.logger.Printf("[DEBUG] %s %v", msg, fields)
	}
}

func (l *simpleLogger) Info(msg string, fields ...interface{}) {
	if l.level <= LogLevelInfo {
		l.logger.Printf("[INFO] %s %v", msg, fields)
	}
}

func (l *simpleLogger) Warn(msg string, fields ...interface{}) {
	if l.level <= LogLevelWarn {
		l.logger.Printf("[WARN] %s %v", msg, fields)
	}
}

func (l *simpleLogger) Error(msg string, fields ...interface{}) {
	if l.level <= LogLevelError {
		l.logger.Printf("[ERROR] %s %v", msg, fields)
	}
}