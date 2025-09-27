// Package logger provides structured logging with Zap integration.
//
// Features:
//   - Structured Logging: JSON-formatted logs with structured fields
//   - Zap Integration: Built on the high-performance Zap logging library
//   - Multiple Levels: Debug, Info, Warn, Error, Fatal, Panic levels
//   - Context Support: Integration with context.Context for request tracing
//   - Configurable Output: Console and file output options
//   - Performance Optimized: High-performance logging with minimal allocation
//
// Example:
//   log := logger.New()
//   log.Info("Application started", zap.String("version", "1.0.0"))
//   log.Error("Operation failed", zap.Error(err))
package logger