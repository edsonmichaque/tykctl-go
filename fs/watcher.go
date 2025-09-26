// Package fs provides file watching functionality for tykctl.
package fs

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// WatchEvent represents a file system watch event
type WatchEvent struct {
	Name      string
	Op        fsnotify.Op
	Timestamp time.Time
}

// WatchHandler defines the interface for handling watch events
type WatchHandler interface {
	HandleEvent(ctx context.Context, event WatchEvent) error
}

// WatchHandlerFunc is a function type that implements WatchHandler
type WatchHandlerFunc func(ctx context.Context, event WatchEvent) error

// HandleEvent implements WatchHandler
func (f WatchHandlerFunc) HandleEvent(ctx context.Context, event WatchEvent) error {
	return f(ctx, event)
}

// Watcher provides file system watching capabilities
type Watcher struct {
	watcher  *fsnotify.Watcher
	handlers map[string][]WatchHandler
	mu       sync.RWMutex
	logger   *zap.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewWatcher creates a new file watcher
func NewWatcher() (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file watcher")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Watcher{
		watcher:  watcher,
		handlers: make(map[string][]WatchHandler),
		logger:   zap.NewNop(), // Default to no-op logger
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// NewWatcherWithLogger creates a new file watcher with a logger
func NewWatcherWithLogger(logger *zap.Logger) (*Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}
	w.logger = logger
	return w, nil
}

// AddHandler adds a handler for a specific path pattern
func (w *Watcher) AddHandler(pathPattern string, handler WatchHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.handlers[pathPattern] == nil {
		w.handlers[pathPattern] = make([]WatchHandler, 0)
	}
	w.handlers[pathPattern] = append(w.handlers[pathPattern], handler)
}

// AddHandlerFunc adds a handler function for a specific path pattern
func (w *Watcher) AddHandlerFunc(pathPattern string, handlerFunc func(ctx context.Context, event WatchEvent) error) {
	w.AddHandler(pathPattern, WatchHandlerFunc(handlerFunc))
}

// Watch starts watching a directory or file
func (w *Watcher) Watch(path string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path")
	}

	if err := w.watcher.Add(absPath); err != nil {
		return errors.Wrapf(err, "failed to watch path: %s", absPath)
	}

	w.logger.Info("Started watching path", zap.String("path", absPath))
	return nil
}

// WatchRecursive watches a directory and all its subdirectories
func (w *Watcher) WatchRecursive(rootPath string) error {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path")
	}

	// Walk the directory tree and add all directories
	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := w.watcher.Add(path); err != nil {
				w.logger.Warn("Failed to watch directory",
					zap.String("path", path),
					zap.Error(err))
			} else {
				w.logger.Debug("Watching directory", zap.String("path", path))
			}
		}
		return nil
	})

	if err != nil {
		return errors.Wrapf(err, "failed to walk directory tree: %s", absPath)
	}

	w.logger.Info("Started recursive watching", zap.String("root", absPath))
	return nil
}

// Start begins the file watching loop
func (w *Watcher) Start() {
	w.wg.Add(1)
	go w.watchLoop()
}

// Stop stops the file watcher
func (w *Watcher) Stop() {
	w.cancel()
	w.wg.Wait()
	w.watcher.Close()
}

// watchLoop is the main watching loop
func (w *Watcher) watchLoop() {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			w.logger.Info("File watcher stopped")
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				w.logger.Warn("File watcher events channel closed")
				return
			}
			w.handleEvent(event)
		case err, ok := <-w.watcher.Errors:
			if !ok {
				w.logger.Warn("File watcher errors channel closed")
				return
			}
			w.logger.Error("File watcher error", zap.Error(err))
		}
	}
}

// handleEvent processes a file system event
func (w *Watcher) handleEvent(event fsnotify.Event) {
	watchEvent := WatchEvent{
		Name:      event.Name,
		Op:        event.Op,
		Timestamp: time.Now(),
	}

	w.logger.Debug("File system event",
		zap.String("name", event.Name),
		zap.String("op", event.Op.String()))

	// Find matching handlers
	w.mu.RLock()
	handlers := w.findMatchingHandlers(event.Name)
	w.mu.RUnlock()

	// Execute handlers
	for _, handler := range handlers {
		if err := handler.HandleEvent(w.ctx, watchEvent); err != nil {
			w.logger.Error("Handler error",
				zap.String("path", event.Name),
				zap.Error(err))
		}
	}
}

// findMatchingHandlers finds handlers that match the given path
func (w *Watcher) findMatchingHandlers(path string) []WatchHandler {
	var matchingHandlers []WatchHandler

	for pattern, handlers := range w.handlers {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			matchingHandlers = append(matchingHandlers, handlers...)
		}
	}

	return matchingHandlers
}

// WatchConfigFile watches a configuration file for changes
func (w *Watcher) WatchConfigFile(configPath string, reloadFunc func() error) error {
	// Watch the directory containing the config file
	configDir := filepath.Dir(configPath)
	if err := w.Watch(configDir); err != nil {
		return errors.Wrap(err, "failed to watch config directory")
	}

	// Add handler for config file changes
	configFileName := filepath.Base(configPath)
	w.AddHandlerFunc(configFileName, func(ctx context.Context, event WatchEvent) error {
		if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
			w.logger.Info("Configuration file changed, reloading",
				zap.String("file", event.Name))
			return reloadFunc()
		}
		return nil
	})

	return nil
}

// WatchExtensions watches the extensions directory for changes
func (w *Watcher) WatchExtensions(extensionsDir string, onChangeFunc func() error) error {
	if err := w.WatchRecursive(extensionsDir); err != nil {
		return errors.Wrap(err, "failed to watch extensions directory")
	}

	// Add handler for extension changes
	w.AddHandlerFunc("tykctl-*", func(ctx context.Context, event WatchEvent) error {
		if event.Op&fsnotify.Write == fsnotify.Write ||
			event.Op&fsnotify.Create == fsnotify.Create ||
			event.Op&fsnotify.Remove == fsnotify.Remove {
			w.logger.Info("Extension changed",
				zap.String("path", event.Name),
				zap.String("op", event.Op.String()))
			return onChangeFunc()
		}
		return nil
	})

	return nil
}

// Global file watcher instance
var globalWatcher *Watcher
var watcherOnce sync.Once

// GetGlobalWatcher returns the global file watcher instance
func GetGlobalWatcher() *Watcher {
	watcherOnce.Do(func() {
		var err error
		globalWatcher, err = NewWatcher()
		if err != nil {
			// Use a no-op logger if we can't create the watcher
			zap.NewNop().Error("Failed to create global file watcher", zap.Error(err))
		}
	})
	return globalWatcher
}

// StartGlobalWatcher starts the global file watcher
func StartGlobalWatcher() {
	watcher := GetGlobalWatcher()
	if watcher != nil {
		watcher.Start()
	}
}

// StopGlobalWatcher stops the global file watcher
func StopGlobalWatcher() {
	if globalWatcher != nil {
		globalWatcher.Stop()
	}
}
