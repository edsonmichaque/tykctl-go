// Package fs provides filesystem abstraction for tykctl.
package fs

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// FS provides filesystem operations
type FS struct {
	fs afero.Fs
}

// New creates a new filesystem instance
func New() *FS {
	return &FS{
		fs: afero.NewOsFs(),
	}
}

// NewMem creates a new in-memory filesystem for testing
func NewMem() *FS {
	return &FS{
		fs: afero.NewMemMapFs(),
	}
}

// MkdirAll creates a directory and all parent directories
// Idempotent: succeeds if directory already exists
func (fs *FS) MkdirAll(ctx context.Context, path string, perm os.FileMode) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Check if directory already exists
		if exists, err := afero.Exists(fs.fs, path); err == nil && exists {
			// Directory exists, check if it's actually a directory
			if info, err := fs.fs.Stat(path); err == nil && info.IsDir() {
				return nil // Directory already exists, success
			}
		}
		return fs.fs.MkdirAll(path, perm)
	}
}

// Stat returns file info
func (fs *FS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return fs.fs.Stat(name)
	}
}

// ReadFile reads a file
func (fs *FS) ReadFile(ctx context.Context, name string) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return afero.ReadFile(fs.fs, name)
	}
}

// WriteFile writes data to a file
// Idempotent: overwrites existing file with same content
func (fs *FS) WriteFile(ctx context.Context, name string, data []byte, perm os.FileMode) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Check if file already exists and has same content
		if exists, err := afero.Exists(fs.fs, name); err == nil && exists {
			if existingData, err := afero.ReadFile(fs.fs, name); err == nil {
				// Compare content byte by byte
				if len(existingData) == len(data) {
					same := true
					for i, b := range data {
						if existingData[i] != b {
							same = false
							break
						}
					}
					if same {
						return nil // File already exists with same content
					}
				}
			}
		}
		return afero.WriteFile(fs.fs, name, data, perm)
	}
}

// Remove removes a file or directory
// Idempotent: succeeds if file/directory doesn't exist
func (fs *FS) Remove(ctx context.Context, name string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Check if file exists before removing
		if exists, err := afero.Exists(fs.fs, name); err != nil {
			return err
		} else if !exists {
			return nil // File doesn't exist, consider it removed
		}
		return fs.fs.Remove(name)
	}
}

// RemoveAll removes a path and all children
// Idempotent: succeeds if path doesn't exist
func (fs *FS) RemoveAll(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Check if path exists before removing
		if exists, err := afero.Exists(fs.fs, path); err != nil {
			return err
		} else if !exists {
			return nil // Path doesn't exist, consider it removed
		}
		return fs.fs.RemoveAll(path)
	}
}

// Exists checks if a path exists
func (fs *FS) Exists(ctx context.Context, path string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		return afero.Exists(fs.fs, path)
	}
}

// JoinPath joins path elements
func (fs *FS) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// GetTempDir returns the system temp directory
func (fs *FS) GetTempDir() string {
	return os.TempDir()
}

// EnsureDir ensures a directory exists, creating it if necessary
// Idempotent: succeeds if directory already exists
func (fs *FS) EnsureDir(ctx context.Context, path string, perm os.FileMode) error {
	return fs.MkdirAll(ctx, path, perm)
}

// EnsureFile ensures a file exists with the given content
// Idempotent: only writes if file doesn't exist or content differs
func (fs *FS) EnsureFile(ctx context.Context, path string, data []byte, perm os.FileMode) error {
	return fs.WriteFile(ctx, path, data, perm)
}

// Touch creates an empty file if it doesn't exist
// Idempotent: succeeds if file already exists
func (fs *FS) Touch(ctx context.Context, path string, perm os.FileMode) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Check if file already exists
		if exists, err := afero.Exists(fs.fs, path); err != nil {
			return err
		} else if exists {
			return nil // File already exists
		}
		// Create empty file
		return afero.WriteFile(fs.fs, path, []byte{}, perm)
	}
}

// RemoveIfExists removes a file or directory if it exists
// Idempotent: succeeds if file/directory doesn't exist
func (fs *FS) RemoveIfExists(ctx context.Context, path string) error {
	return fs.Remove(ctx, path)
}

// RemoveAllIfExists removes a path and all children if it exists
// Idempotent: succeeds if path doesn't exist
func (fs *FS) RemoveAllIfExists(ctx context.Context, path string) error {
	return fs.RemoveAll(ctx, path)
}

// CreateDirIfNotExists creates a directory if it doesn't exist
// Idempotent: succeeds if directory already exists
func (fs *FS) CreateDirIfNotExists(ctx context.Context, path string, perm os.FileMode) error {
	return fs.MkdirAll(ctx, path, perm)
}

// WriteFileIfDifferent writes a file only if content is different
// Idempotent: skips write if content is identical
func (fs *FS) WriteFileIfDifferent(ctx context.Context, path string, data []byte, perm os.FileMode) error {
	return fs.WriteFile(ctx, path, data, perm)
}
